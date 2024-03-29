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
	PairModifier            int
	TempoModifier           int
	RookFileModifier        int
	BlockedPiecesModifier   int
	KingSafetyModifier      int
	PassedPawnModifier      int
	MobilityModifier        int
}

type Evaluation struct {
	Game                  *game.Game
	White                 EvaluationPart
	Black                 EvaluationPart
	GameOver              bool
	TotalScore            int
	TotalScorePerspective int
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

func (e *Evaluation) GetPieceTypeValue(pieceType dragontoothmg.Piece) int {
	return weights[game.White].Midgame.Material[pieceType]
}

func (e *Evaluation) updateTotal() {
	if e.GameOver {
		if e.Game.Result == game.WhiteWon {
			e.TotalScore = 100000
		} else if e.Game.Result == game.BlackWon {
			e.TotalScore = -100000
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

	e.TotalScore += e.White.PairModifier
	e.TotalScore -= e.Black.PairModifier
	e.TotalScore += e.White.TempoModifier
	e.TotalScore -= e.Black.TempoModifier
	e.TotalScore += e.White.RookFileModifier
	e.TotalScore -= e.Black.RookFileModifier
	e.TotalScore += e.White.BlockedPiecesModifier
	e.TotalScore -= e.Black.BlockedPiecesModifier
	e.TotalScore += e.White.KingSafetyModifier
	e.TotalScore -= e.Black.KingSafetyModifier
	e.TotalScore += e.White.PassedPawnModifier
	e.TotalScore -= e.Black.PassedPawnModifier
	e.TotalScore += e.White.MobilityModifier
	e.TotalScore -= e.Black.MobilityModifier

	e.TotalScorePerspective = e.TotalScore
	if !e.Game.Position.Wtomove {
		e.TotalScorePerspective = -e.TotalScore
	}
}

func calculateEvaluationPart(g *game.Game, color game.PlayerColor) EvaluationPart {
	ps, pstMid, pstEnd := calculateMaterialScore(g, color)
	evalPart := EvaluationPart{
		GamePhase:               calculateGamephase(g, color),
		PieceScore:              ps,
		PieceSquareScoreMidgame: pstMid,
		PieceSquareScoreEndgame: pstEnd,
		PairModifier:            calculatePairModifier(g, color),
		TempoModifier:           calculateTempoModifier(g, color),
		RookFileModifier:        calculateRookModifier(g, color),
		PassedPawnModifier:      calculatePassedPawns(g, color),
		MobilityModifier:        calculateMobilityModifier(g, color),
		KingSafetyModifier:      calculateKingSafety(g, color),
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

	// Thanks to https://github.com/dylhunn/dragontooth for the extract pieces from bitboard pattern
	for x = bitboard; x != 0; x &= x - 1 {
		square := bits.TrailingZeros64(x)

		ps += weights[color].Midgame.Material[pieceType]
		pstMid += weights[color].Midgame.PieceSquareTables[pieceType][square]
		pstEnd += weights[color].Endgame.PieceSquareTables[pieceType][square]
	}
	return ps, pstMid, pstEnd
}

func calculatePairModifier(g *game.Game, color game.PlayerColor) (result int) {
	bboards := g.Position.White
	if color == game.Black {
		bboards = g.Position.Black
	}

	if bits.OnesCount64(bboards.Bishops) >= 2 {
		result += weights[color].AdditionalModifier.BishopPairModifier
	}
	if bits.OnesCount64(bboards.Knights) >= 2 {
		result += weights[color].AdditionalModifier.KnightPairModifier
	}
	if bits.OnesCount64(bboards.Rooks) >= 2 {
		result += weights[color].AdditionalModifier.RookPairModifier
	}

	return
}

func calculateTempoModifier(g *game.Game, color game.PlayerColor) (result int) {
	if g.Position.Wtomove == bool(color) {
		result += weights[color].AdditionalModifier.TempoModifier
	}
	return
}

func calculateRookModifier(g *game.Game, color game.PlayerColor) (result int) {
	bboardsOwn := g.Position.White
	bboardsOther := g.Position.Black
	if color == game.Black {
		bboardsOwn = g.Position.Black
		bboardsOther = g.Position.White
	}

	pawnFillOwn := calculatePawnFileFill(bboardsOwn.Pawns)
	pawnFillOther := calculatePawnFileFill(bboardsOther.Pawns)

	openFiles := ^pawnFillOwn & ^pawnFillOther
	halfOpenFiles := ^pawnFillOwn ^ openFiles

	rooksOnOpenFiles := bits.OnesCount64(bboardsOwn.Rooks & openFiles)
	rooksOnHalfOpenFiles := bits.OnesCount64(bboardsOwn.Rooks & halfOpenFiles)

	return rooksOnOpenFiles*weights[color].AdditionalModifier.OpenRookModifier + rooksOnHalfOpenFiles*weights[color].AdditionalModifier.HalfRookModifier
}

func calculatePassedPawns(g *game.Game, color game.PlayerColor) (result int) {
	if color == game.White {
		frontSpansBlack := calculatePawnSouthFill(g.Position.Black.Pawns) & ^g.Position.Black.Pawns
		attackingSpansBlack := frontSpansBlack
		attackingSpansBlack |= (frontSpansBlack << 1) & ^bitboardFileA //shift everything east and care for wraps
		attackingSpansBlack |= (frontSpansBlack >> 1) & ^bitboardFileH //shift everything west and care for wraps
		whitePassedPawns := g.Position.White.Pawns & ^attackingSpansBlack
		for x := whitePassedPawns; x != 0; x &= x - 1 {
			square := bits.TrailingZeros64(x)
			result += weights[color].Midgame.PassedPawnModifier[square]
		}
		return
	}

	//black
	frontSpansWhite := calculatePawnNorthFill(g.Position.White.Pawns) & ^g.Position.White.Pawns
	attackingSpansWhite := frontSpansWhite
	attackingSpansWhite |= (frontSpansWhite << 1) & ^bitboardFileA //shift everything east and care for wraps
	attackingSpansWhite |= (frontSpansWhite >> 1) & ^bitboardFileH //shift everything west and care for wraps
	blackPassedPawns := g.Position.Black.Pawns & ^attackingSpansWhite
	for x := blackPassedPawns; x != 0; x &= x - 1 {
		square := bits.TrailingZeros64(x)
		result += weights[color].Midgame.PassedPawnModifier[square]
	}
	return
}

func calculateMobilityModifier(g *game.Game, color game.PlayerColor) (result int) {
	return calculateDiagonalMobilityModifier(g, color) + calculateLinearMobilityModifier(g, color)
}

func calculateDiagonalMobilityModifier(g *game.Game, color game.PlayerColor) (result int) {
	diagonalBB := g.Position.White.Bishops | g.Position.White.Queens
	ownBB := g.Position.White.All
	if color == game.Black {
		diagonalBB = g.Position.Black.Bishops | g.Position.Black.Queens
		ownBB = g.Position.Black.All
	}
	allBB := g.Position.White.All | g.Position.Black.All

	for x := diagonalBB; x != 0; x &= x - 1 {
		square := bits.TrailingZeros64(x)
		movableSquares := dragontoothmg.CalculateBishopMoveBitboard(uint8(square), allBB) & ^ownBB
		result += bits.OnesCount64(movableSquares) * weights[color].AdditionalModifier.DiagonalMobilityModifier
	}

	return
}

func calculateLinearMobilityModifier(g *game.Game, color game.PlayerColor) (result int) {
	linearBB := g.Position.White.Rooks | g.Position.White.Queens
	ownBB := g.Position.White.All
	if color == game.Black {
		linearBB = g.Position.Black.Rooks | g.Position.Black.Queens
		ownBB = g.Position.Black.All
	}
	allBB := g.Position.White.All | g.Position.Black.All

	for x := linearBB; x != 0; x &= x - 1 {
		square := bits.TrailingZeros64(x)
		movableSquares := dragontoothmg.CalculateRookMoveBitboard(uint8(square), allBB) & ^ownBB
		result += bits.OnesCount64(movableSquares) * weights[color].AdditionalModifier.LinearMobilityModifier
	}

	return
}

func calculateKingSafety(g *game.Game, color game.PlayerColor) (result int) {
	attackUnitCounts := 0
	kingZoneBB := g.Position.White.Kings
	linearAttackersBB := g.Position.Black.Rooks | g.Position.Black.Queens
	diagonalAttackersBB := g.Position.Black.Bishops | g.Position.Black.Queens
	ownBB := g.Position.White.All
	if color == game.Black {
		kingZoneBB = g.Position.Black.Kings
		linearAttackersBB = g.Position.White.Rooks | g.Position.White.Queens
		diagonalAttackersBB = g.Position.White.Bishops | g.Position.White.Queens
		ownBB = g.Position.Black.All
	}
	allBB := g.Position.White.All | g.Position.Black.All

	// up + down
	kingZoneBB = kingZoneBB | (kingZoneBB << 8) | (kingZoneBB >> 8)
	// left + right <- This also includes the up+down squares of above
	kingZoneBB = kingZoneBB | ((kingZoneBB & ^bitboardFileA) >> 1) | ((kingZoneBB & ^bitboardFileH) << 1)

	for x := linearAttackersBB; x != 0; x &= x - 1 {
		square := bits.TrailingZeros64(x)
		movableSquares := dragontoothmg.CalculateRookMoveBitboard(uint8(square), allBB) & ^ownBB
		if movableSquares&kingZoneBB != 0 {
			attackUnitCounts += weights[color].AdditionalModifier.LinearAttackUnit
		}
	}
	for x := diagonalAttackersBB; x != 0; x &= x - 1 {
		square := bits.TrailingZeros64(x)
		movableSquares := dragontoothmg.CalculateBishopMoveBitboard(uint8(square), allBB) & ^ownBB
		if movableSquares&kingZoneBB != 0 {
			attackUnitCounts += weights[color].AdditionalModifier.DiagonalAttackUnit
		}
	}

	if attackUnitCounts >= len(weights[color].AdditionalModifier.KingSafetyTable) {
		return weights[color].AdditionalModifier.KingSafetyTable[len(weights[color].AdditionalModifier.KingSafetyTable)-1]
	}

	return -weights[color].AdditionalModifier.KingSafetyTable[attackUnitCounts]
}
