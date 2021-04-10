package main

import (
	"fmt"

	"go.janniklasrichter.de/axwchessbot/evaluation"
	"go.janniklasrichter.de/axwchessbot/game"
)

func main() {
	fens := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RN2KBNR w KQkq - 0 1",
		"rn2kbnr/pppppppp/8/8/8/8/PPPPPPPP/RN2KBNR w KQkq - 0 1",
		"rn2kbnr/pp3ppp/8/8/8/8/PPPPPPPP/RN2KBNR w KQkq - 0 1",
		"rn2kbnr/pp3ppp/8/8/8/8/PPPPPPPP/RN2KBN1 w Qkq - 0 1",
	}

	for _, fen := range fens {
		game := game.NewFromFen(fen)
		eval := evaluation.CalculateEvaluation(game)
		fmt.Println("Board: ", game.Position.ToFen())
		fmt.Println("Evaluation: ", eval.TotalScore)
	}
}
