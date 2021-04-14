package evaluation_provider

import (
	"github.com/dylhunn/dragontoothmg"
	"go.janniklasrichter.de/axwchessbot/game"
)

type EvaluationProvider interface {
	CalculateEvaluation(g *game.Game) int
	GetPieceTypeValue(pieceType dragontoothmg.Piece) int
}
