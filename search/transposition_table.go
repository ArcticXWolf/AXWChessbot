package search

import "github.com/dylhunn/dragontoothmg"

type transpositionTableEntry struct {
	lock  uint8
	hash  uint64
	move  dragontoothmg.Move
	score int32
	depth uint8
}
