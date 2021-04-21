package search

import "github.com/dylhunn/dragontoothmg"

type killerMoveTable struct {
	entries [][2]dragontoothmg.Move
}

func newKillerMoveTable(initialStackSize int) *killerMoveTable {
	return &killerMoveTable{
		entries: make([][2]dragontoothmg.Move, initialStackSize),
	}
}

func (k *killerMoveTable) clear() {
	k.entries = nil
}

func (k *killerMoveTable) clearLevel(level int) {
	k.entries[level][0] = 0
	k.entries[level][1] = 0
}

func (k *killerMoveTable) update(ply int, move dragontoothmg.Move) {
	if ply >= len(k.entries)+1 {
		// something went wrong and we skipped a whole level
		// return and dont use killers
		return
	}

	if ply >= len(k.entries) {
		k.entries = append(k.entries, [2]dragontoothmg.Move{})
	}

	if k.entries[ply][0] != move {
		k.entries[ply][1], k.entries[ply][0] = k.entries[ply][0], move
	}
}

func (k *killerMoveTable) fetch(ply int) (move1 dragontoothmg.Move, found1 bool, move2 dragontoothmg.Move, found2 bool) {
	if ply >= len(k.entries) {
		return
	}
	return k.entries[ply][0], k.entries[ply][0] != 0, k.entries[ply][1], k.entries[ply][1] != 0
}
