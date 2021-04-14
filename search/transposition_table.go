package search

import (
	"sync/atomic"

	"github.com/dylhunn/dragontoothmg"
)

type alphaBetaBound uint8

const (
	alphaBetaBoundLower = iota
	alphaBetaBoundExact
	alphaBetaBoundUpper
)

type transpositionTableEntry struct {
	lock  int32
	move  dragontoothmg.Move
	score int32
	depth uint8
	bound alphaBetaBound
}

type TranspositionTable struct {
	maxSize int
	entries map[uint64]transpositionTableEntry
}

// 1 entry is 20 bytes
func NewTranspositionTable(maxSize int) *TranspositionTable {
	return &TranspositionTable{
		maxSize: maxSize,
		entries: make(map[uint64]transpositionTableEntry, maxSize),
	}
}

func (tt *TranspositionTable) Empty() {
	tt.entries = make(map[uint64]transpositionTableEntry, tt.maxSize)
}

func (tt *TranspositionTable) InsertIfNeeded(hash uint64, move dragontoothmg.Move, score int, depth int, bound alphaBetaBound) {
	entry, found := tt.entries[hash]

	if !found {
		tt.entries[hash] = transpositionTableEntry{
			lock:  0,
			move:  move,
			score: int32(score),
			depth: uint8(depth),
			bound: bound,
		}
		return
	}

	if atomic.CompareAndSwapInt32(&entry.lock, 0, 1) {
		if depth >= int(entry.depth) {
			entry.move = move
			entry.score = int32(score)
			entry.depth = uint8(depth)
			entry.bound = bound
		}
		atomic.StoreInt32(&entry.lock, 0)
	}
}

func (tt *TranspositionTable) Load(hash uint64) (move dragontoothmg.Move, score int, depth int, bound alphaBetaBound, ok bool) {
	ok = false
	entry, found := tt.entries[hash]

	if !found {
		return
	}

	if atomic.CompareAndSwapInt32(&entry.lock, 0, 1) {
		move = entry.move
		score = int(entry.score)
		depth = int(entry.depth)
		bound = entry.bound
		ok = true
		atomic.StoreInt32(&entry.lock, 0)
	}

	return
}
