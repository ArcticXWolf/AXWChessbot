package evaluation_null

import (
	"go.janniklasrichter.de/axwchessbot/game"
)

type EvaluationNull struct {
	Game       *game.Game
	TotalScore int
}

func (e *EvaluationNull) CalculateEvaluation(g *game.Game) int {
	e.Game = g
	e.TotalScore = 0
	return e.TotalScore
}
