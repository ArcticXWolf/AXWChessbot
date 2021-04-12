package evaluation_provider

import "go.janniklasrichter.de/axwchessbot/game"

type EvaluationProvider interface {
	CalculateEvaluation(g *game.Game) int
}
