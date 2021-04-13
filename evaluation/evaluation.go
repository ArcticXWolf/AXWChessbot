package evaluation

import (
	"math/bits"

	"github.com/dylhunn/dragontoothmg"
	"go.janniklasrichter.de/axwchessbot/game"
)

type EvaluationPart struct {
	GamePhase               uint8
	PieceScore              int
	PieceSquareScoreMidgame int
	PieceSquareScoreEndgame int
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
			e.TotalScore = 1000000
		} else if e.Game.Result == game.BlackWon {
			e.TotalScore = -1000000
		} else {
			e.TotalScore = 0 // DRAW
		}

		e.TotalScorePerspective = e.TotalScore
		if !e.Game.Position.Wtomove {
			e.TotalScorePerspective = -e.TotalScore
		}
		return
	}

	gamePhase := int(e.White.GamePhase + e.Black.GamePhase)
	e.TotalScore = 0
	e.TotalScore += e.White.PieceScore
	e.TotalScore -= e.Black.PieceScore

	e.TotalScore += (gamePhase*e.White.PieceSquareScoreMidgame + (24-gamePhase)*e.White.PieceSquareScoreEndgame) / 24
	e.TotalScore -= (gamePhase*e.Black.PieceSquareScoreMidgame + (24-gamePhase)*e.Black.PieceSquareScoreEndgame) / 24

	e.TotalScorePerspective = e.TotalScore
	if !e.Game.Position.Wtomove {
		e.TotalScorePerspective = -e.TotalScore
	}
}

func (e *Evaluation) CalculateEvaluation(g *game.Game) int {
	e.Game = g
	e.GameOver = e.Game.Result != game.GameNotOver

	e.White = EvaluationPart{}
	e.Black = EvaluationPart{}

	if !e.GameOver {
		e.White = calculateEvaluationPart(g, game.White)
		e.Black = calculateEvaluationPart(g, game.Black)
	}

	e.updateTotal()
	return e.TotalScorePerspective
}

func calculateEvaluationPart(g *game.Game, color game.PlayerColor) EvaluationPart {
	ps, pstMid, pstEnd := calculateMaterialScore(g, color)
	evalPart := EvaluationPart{
		GamePhase:               calculateGamephase(g, color),
		PieceScore:              ps,
		PieceSquareScoreMidgame: pstMid,
		PieceSquareScoreEndgame: pstEnd,
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

func calculateMaterialScore(g *game.Game, color game.PlayerColor) (int, int, int) {
	var ps, pstMid, pstEnd, newPs, newPstMid, newPstEnd int
	var bitboard uint64

	bboards := g.Position.White
	if color == game.Black {
		bboards = g.Position.Black
	}

	for i := 1; i <= 6; i++ {
		bitboard = getBitboardByPieceType(&bboards, dragontoothmg.Piece(i))
		newPs, newPstMid, newPstEnd = calculateMaterialScoreForPieceType(g, color, dragontoothmg.Piece(i), bitboard)
		ps += newPs
		pstMid += newPstMid
		pstEnd += newPstEnd
	}

	return ps, pstMid, pstEnd
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

func calculateMaterialScoreForPieceType(g *game.Game, color game.PlayerColor, pieceType dragontoothmg.Piece, bitboard uint64) (int, int, int) {
	var x uint64
	ps, pstMid, pstEnd := 0, 0, 0
	for x = bitboard; x != 0; x &= x - 1 {
		square := bits.TrailingZeros64(x)

		ps += weights[color].Midgame.Material[pieceType]
		pstMid += weights[color].Midgame.PieceSquareTables[pieceType][square]
		pstEnd += weights[color].Endgame.PieceSquareTables[pieceType][square]
	}
	return ps, pstMid, pstEnd
}
