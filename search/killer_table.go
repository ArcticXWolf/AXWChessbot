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
