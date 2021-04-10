package game

import (
	"testing"
)

func TestCheckmate(t *testing.T) {
	fenStr := "rn1qkbnr/pbpp1ppp/1p6/4p3/2B1P3/5Q2/PPPP1PPP/RNB1K1NR w KQkq - 0 1"
	game := NewFromFen(fenStr)
	game.pushMove("f3f7")
	if game.Result != WhiteWon {
		t.Fatalf("expected result %d but got %d", WhiteWon, game.Result)
	}
}
