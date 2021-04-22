package game

import (
	"reflect"
	"testing"
)

func TestGameOverDetection(t *testing.T) {
	type args struct {
		fen string
	}
	tests := []struct {
		name       string
		args       args
		wantResult GameResult
		wantReason DrawReason
	}{
		{"GameNotOver", args{"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"}, GameNotOver, NoDraw},
		{"Checkmate White", args{"rn1qkbnr/pbpp1Qpp/1p6/4p3/2B1P3/8/PPPP1PPP/RNB1K1NR b KQkq - 0 1"}, WhiteWon, NoDraw},
		{"Checkmate Black", args{"rnb1k1nr/pppp1ppp/8/2b1p3/2B1P3/2N5/PPPP1qPP/R1BQK2R w KQkq - 0 1"}, BlackWon, NoDraw},
		{"Draw Stalemate", args{"7k/6R1/5K2/8/8/8/8/8 b - - 0 1"}, Draw, Stalemate},
		{"Draw InsufficientMaterial KBvKB Whitefields", args{"8/2kb4/8/8/8/3B4/4K3/8 w - - 0 1"}, Draw, InsufficientMaterial},
		{"Draw InsufficientMaterial KBvKB Blackfields", args{"8/2k1b3/8/8/5B2/5K2/8/8 w - - 0 1"}, Draw, InsufficientMaterial},
		{"Draw InsufficientMaterial KvK", args{"8/2k5/8/8/8/8/3K4/8 w - - 0 1"}, Draw, InsufficientMaterial},
		{"Draw InsufficientMaterial KNvK", args{"8/2k5/8/8/4N3/5K2/8/8 w - - 0 1"}, Draw, InsufficientMaterial},
		{"Draw InsufficientMaterial KvKB", args{"8/2k5/3b4/8/8/5K2/8/8 w - - 0 1"}, Draw, InsufficientMaterial},
		{"Draw FiftyMoveRule", args{"rNB1k2r/ppp2ppp/3p3B/3np2Q/3NP2q/3P3b/PPP2PPP/Rnb1K2R w - - 100 43"}, Draw, FiftyMoveRule},
		{"GameNotOver KBvKB", args{"8/2k1b3/8/8/4B3/5K2/8/8 w - - 0 1"}, GameNotOver, NoDraw},
		{"King Missing", args{"8/8/8/8/8/5K2/8/8 w - - 0 1"}, GameNotOver, NoDraw},
		{"GameOver from #6", args{"3r2k1/rp1n3p/2pb2p1/p1n3P1/2P1P3/P1N5/1PQ1BB1q/3R1R1K w - - 6 36"}, BlackWon, NoDraw},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFromFen(tt.args.fen)
			if !reflect.DeepEqual(got.Result, tt.wantResult) {
				t.Errorf("Game Over Detection, Result = %v, want %v", got.Result, tt.wantResult)
			}
			if !reflect.DeepEqual(got.DrawReason, tt.wantReason) {
				t.Errorf("Game Over Detection, Reason = %v, want %v", got.DrawReason, tt.wantReason)
			}
		})
	}
}

func TestDrawByRepetition(t *testing.T) {
	game := NewFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	game.PushMoveStr("e2e4")
	game.PushMoveStr("e7e5")

	game.PushMoveStr("f1c4")
	game.PushMoveStr("f8c5")
	game.PushMoveStr("c4f1")
	game.PushMoveStr("c5f8")

	game.PushMoveStr("f1c4")
	game.PushMoveStr("f8c5")
	game.PushMoveStr("c4f1")
	game.PushMoveStr("c5f8")

	game.PushMoveStr("f1c4")
	game.PushMoveStr("f8c5")
	game.PushMoveStr("c4f1")
	game.PushMoveStr("c5f8")

	if game.Result != Draw {
		t.Errorf("Draw by Repetition, Result = %v, want %v", game.Result, Draw)
	}
	if game.DrawReason != ThreefoldRepetition {
		t.Errorf("Draw by Repetition, Result = %v, want %v", game.DrawReason, ThreefoldRepetition)
	}
}

func TestPushPopMove(t *testing.T) {
	game := New()
	unmodified_game := New()
	game.PushMoveStr("e2e4")

	if reflect.DeepEqual(game, unmodified_game) {
		t.Errorf("PushMove = %v is wrongfully equal to %v", game, unmodified_game)
	}

	game.PopMove()

	if !reflect.DeepEqual(game, unmodified_game) {
		t.Errorf("PushPopMove = %v, want %v", game, unmodified_game)
	}
}

func TestPushMoveError(t *testing.T) {
	game := New()

	err := game.PushMoveStr("e123e1")

	if err == nil {
		t.Errorf("PushMoveError = %v is not an error", err)
	}
}

func TestPopNoMoves(t *testing.T) {
	game := New()

	err := game.PopMove()

	if err == nil {
		t.Errorf("PopNoMoves = %v, is not an error", err)
	}
}
