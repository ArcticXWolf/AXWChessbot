package search

import (
	"testing"

	"github.com/dylhunn/dragontoothmg"
	"go.janniklasrichter.de/axwchessbot/game"
)

func Test_isCaptureOrPromotionMove(t *testing.T) {
	type args struct {
		game *game.Game
		move dragontoothmg.Move
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"En Passant", args{game.NewFromFen("rnbqkbnr/pp3ppp/3p4/3Pp3/1Pp1P3/5P2/P1P3PP/RNBQKBNR b KQkq b3 0 5"), getMove("c4b3")}, true},
		{"Missed En Passant", args{game.NewFromFen("rnbqkbnr/pp4pp/3p1p2/3Pp3/1Pp1P3/5P1N/P1P3PP/RNBQKB1R b KQkq - 1 6"), getMove("c4b3")}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isCaptureOrPromotionMove(tt.args.game, tt.args.move); got != tt.want {
				t.Errorf("isCaptureOrPromotionMove() = %v, want %v", got, tt.want)
			}
		})
	}
}
