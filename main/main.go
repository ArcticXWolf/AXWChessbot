package main

import (
	"fmt"

	"github.com/dylhunn/dragontoothmg"
	"go.janniklasrichter.de/axwchessbot/evaluation"
)

func main() {
	board := dragontoothmg.ParseFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	eval := evaluation.CalculateEvaluation(board)
	fmt.Println("Board: ", board.ToFen())
	fmt.Println("Evaluation: ", eval.Total)
}
