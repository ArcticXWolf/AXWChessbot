package search

import (
	"github.com/dylhunn/dragontoothmg"
	"go.janniklasrichter.de/axwchessbot/game"
)

func getMove(moveStr string) dragontoothmg.Move {
	move, _ := dragontoothmg.ParseMove(moveStr)
	return move
}

func (s *Search) getCaptureMVVLVA(move dragontoothmg.Move, bitboardsOwn dragontoothmg.Bitboards, bitboardsOpponent dragontoothmg.Bitboards) (score int) {
	pieceTypeFrom, _ := getPieceTypeAtPosition(move.From(), bitboardsOwn)
	pieceTypeTo, _ := getPieceTypeAtPosition(move.To(), bitboardsOpponent)

	return (1200 - s.evaluationProvider.GetPieceTypeValue(pieceTypeTo)) + int(pieceTypeFrom)
}

func isCaptureOrPromotionMove(game *game.Game, move dragontoothmg.Move) bool {
	bitboardsOwn := game.Position.White
	bitboardsOpponent := game.Position.Black
	if !game.Position.Wtomove {
		bitboardsOwn = game.Position.Black
		bitboardsOpponent = game.Position.White
	}

	return bitboardsOpponent.All&(1<<move.To()) > 0 || move.Promote() > 0 || (bitboardsOwn.Pawns&(1<<move.From()) > 0 && game.GetEnPassentSquare() == move.To())
}

func getPieceTypeAtPosition(position uint8, bitboards dragontoothmg.Bitboards) (pieceType dragontoothmg.Piece, occupied bool) {
	if bitboards.Pawns&(1<<position) > 0 {
		return dragontoothmg.Pawn, true
	} else if bitboards.Knights&(1<<position) > 0 {
		return dragontoothmg.Knight, true
	} else if bitboards.Bishops&(1<<position) > 0 {
		return dragontoothmg.Bishop, true
	} else if bitboards.Rooks&(1<<position) > 0 {
		return dragontoothmg.Rook, true
	} else if bitboards.Queens&(1<<position) > 0 {
		return dragontoothmg.Queen, true
	} else if bitboards.Kings&(1<<position) > 0 {
		return dragontoothmg.King, true
	}
	return
}
