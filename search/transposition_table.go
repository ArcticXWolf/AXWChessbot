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
	MaxSizeInEntries int
	Entries          map[uint64]transpositionTableEntry
}

func NewTranspositionTable(maxSizeInBytes int) *TranspositionTable {
	maxSizeInEntries := maxSizeInBytes / 100 // Just a guesstimate, I still have not found how to calculate the bytesize of a map
	return &TranspositionTable{
		MaxSizeInEntries: maxSizeInEntries,
		Entries:          make(map[uint64]transpositionTableEntry, maxSizeInEntries),
	}
}

func (tt *TranspositionTable) Empty() {
	tt.Entries = make(map[uint64]transpositionTableEntry, tt.MaxSizeInEntries)
}

func (tt *TranspositionTable) InsertIfNeeded(hash uint64, move dragontoothmg.Move, score int, depth int, bound alphaBetaBound) {
	if len(tt.Entries) >= tt.MaxSizeInEntries {
		tt.Empty()
	}

	entry, found := tt.Entries[hash]

	if !found {
		tt.Entries[hash] = transpositionTableEntry{
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
	entry, found := tt.Entries[hash]

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
