package evaluation

import (
	"github.com/dylhunn/dragontoothmg"
	"go.janniklasrichter.de/axwchessbot"
	"go.janniklasrichter.de/axwchessbot/game"
)

type EvaluationPart struct {
	GamePhase               int
	PieceScore              map[dragontoothmg.Piece]int
	PieceSquareScoreMidgame map[dragontoothmg.Piece]int
	PieceSquareScoreEndgame map[dragontoothmg.Piece]int
}

type Evaluation struct {
	Game     *game.Game
	White    EvaluationPart
	Black    EvaluationPart
	GameOver bool
	Total    int
}

func CalculateEvaluation(g *game.Game) Evaluation {
	whiteEval := calculateEvaluationPart(g, axwchessbot.White)
	blackEval := calculateEvaluationPart(g, axwchessbot.Black)
	eval := Evaluation{
		Game:  g,
		White: whiteEval,
		Black: blackEval,
	}
	return eval
}

func calculateEvaluationPart(g *game.Game, color axwchessbot.PlayerColor) EvaluationPart {
	evalPart := EvaluationPart{
		GamePhase: 0,
	}
	return evalPart
}
