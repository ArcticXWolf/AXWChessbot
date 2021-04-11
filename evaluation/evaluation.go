package evaluation

import (
	"math/bits"

	"github.com/dylhunn/dragontoothmg"
	"go.janniklasrichter.de/axwchessbot/game"
)

type EvaluationPart struct {
	GamePhase               uint8
	PieceScore              map[dragontoothmg.Piece]int
	PieceSquareScoreMidgame map[dragontoothmg.Piece]int
	PieceSquareScoreEndgame map[dragontoothmg.Piece]int
}

type Evaluation struct {
	Game                  *game.Game
	White                 EvaluationPart
	Black                 EvaluationPart
	GameOver              bool
	TotalScore            int
	TotalScorePerspective int
}

func (e *Evaluation) updateTotal() {
	if e.GameOver {
		if e.Game.Result == game.WhiteWon {
			e.TotalScore = int(^uint(0) >> 1) // MAX INT
		} else if e.Game.Result == game.BlackWon {
			e.TotalScore = -int(^uint(0)>>1) - 1 // MIN INT
		} else {
			e.TotalScore = 0 // DRAW
		}

		e.TotalScorePerspective = e.TotalScore
		if !e.Game.Position.Wtomove {
			e.TotalScorePerspective = -e.TotalScorePerspective
		}
		return
	}

	e.TotalScore = 0
	e.TotalScore += sumMapValues(e.White.PieceScore)
	e.TotalScore -= sumMapValues(e.Black.PieceScore)

	e.TotalScorePerspective = e.TotalScore
	if !e.Game.Position.Wtomove {
		e.TotalScorePerspective = -e.TotalScorePerspective
	}
}

func CalculateEvaluation(g *game.Game) *Evaluation {
	is_game_over := g.Result != game.GameNotOver

	whiteEval := EvaluationPart{}
	blackEval := EvaluationPart{}

	if !is_game_over {
		whiteEval = calculateEvaluationPart(g, game.White)
		blackEval = calculateEvaluationPart(g, game.Black)
	}

	eval := &Evaluation{
		Game:     g,
		White:    whiteEval,
		Black:    blackEval,
		GameOver: is_game_over,
	}
	eval.updateTotal()
	return eval
}

func calculateEvaluationPart(g *game.Game, color game.PlayerColor) EvaluationPart {
	evalPart := EvaluationPart{
		GamePhase:               calculateGamephase(g, color),
		PieceScore:              calculatePieceScore(g, color),
		PieceSquareScoreMidgame: make(map[dragontoothmg.Piece]int),
		PieceSquareScoreEndgame: make(map[dragontoothmg.Piece]int),
	}
	return evalPart
}

func calculateGamephase(g *game.Game, color game.PlayerColor) uint8 {
	bboards := g.Position.White
	if color == game.Black {
		bboards = g.Position.Black
	}

	return uint8(bits.OnesCount64(bboards.Knights)) +
		uint8(bits.OnesCount64(bboards.Bishops)) +
		2*uint8(bits.OnesCount64(bboards.Rooks)) +
		4*uint8(bits.OnesCount64(bboards.Queens))
}

func calculatePieceScore(g *game.Game, color game.PlayerColor) map[dragontoothmg.Piece]int {
	ps := make(map[dragontoothmg.Piece]int, 6)
	bboards := g.Position.White
	if color == game.Black {
		bboards = g.Position.Black
	}

	for i := 1; i <= 6; i++ {
		bitboard := getBitboardByPieceType(&bboards, dragontoothmg.Piece(i))
		count := bits.OnesCount64(bitboard)
		ps[dragontoothmg.Piece(i)] += count * GetWeights().Midgame.Material[dragontoothmg.Piece(i)]
	}

	return ps
}

func getBitboardByPieceType(bbs *dragontoothmg.Bitboards, pieceType dragontoothmg.Piece) uint64 {
	switch pieceType {
	case dragontoothmg.Pawn:
		return bbs.Pawns
	case dragontoothmg.Knight:
		return bbs.Knights
	case dragontoothmg.Bishop:
		return bbs.Bishops
	case dragontoothmg.Rook:
		return bbs.Rooks
	case dragontoothmg.Queen:
		return bbs.Queens
	default:
		return bbs.Kings
	}
}

func sumMapValues(mapToSum map[dragontoothmg.Piece]int) int {
	result := 0
	for _, value := range mapToSum {
		result += value
	}
	return result
}
