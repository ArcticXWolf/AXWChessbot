package evaluation

import (
	"reflect"
	"testing"

	"github.com/dylhunn/dragontoothmg"
	"go.janniklasrichter.de/axwchessbot/game"
)

func Test_calculateGamephase(t *testing.T) {
	type args struct {
		g     *game.Game
		color game.PlayerColor
	}
	tests := []struct {
		name string
		args args
		want uint8
	}{
		{"GameStart White", args{game.NewFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"), game.White}, 12},
		{"GameStart Black", args{game.NewFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"), game.Black}, 12},
		{"KvKQ White", args{game.NewFromFen("1q6/2k5/8/8/8/8/3K4/8 w - - 0 1"), game.White}, 0},
		{"KvKQ Black", args{game.NewFromFen("1q6/2k5/8/8/8/8/3K4/8 w - - 0 1"), game.Black}, 4},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateGamephase(tt.args.g, tt.args.color); got != tt.want {
				t.Errorf("calculateGamephase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculatePieceScore(t *testing.T) {
	type args struct {
		g     *game.Game
		color game.PlayerColor
	}
	tests := []struct {
		name string
		args args
		want map[dragontoothmg.Piece]int
	}{
		{"GameStart White", args{game.NewFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"), game.White}, map[dragontoothmg.Piece]int{dragontoothmg.Pawn: 800, dragontoothmg.Knight: 640, dragontoothmg.Bishop: 660, dragontoothmg.Rook: 1000, dragontoothmg.Queen: 900, dragontoothmg.King: 0}},
		{"GameStart Black", args{game.NewFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"), game.Black}, map[dragontoothmg.Piece]int{dragontoothmg.Pawn: 800, dragontoothmg.Knight: 640, dragontoothmg.Bishop: 660, dragontoothmg.Rook: 1000, dragontoothmg.Queen: 900, dragontoothmg.King: 0}},
		{"KRNPvKQB White", args{game.NewFromFen("1q6/2k1b3/8/8/8/5P2/3KN3/7R w - - 0 1"), game.White}, map[dragontoothmg.Piece]int{dragontoothmg.Pawn: 100, dragontoothmg.Knight: 320, dragontoothmg.Bishop: 0, dragontoothmg.Rook: 500, dragontoothmg.Queen: 0, dragontoothmg.King: 0}},
		{"KRNPvKQB Black", args{game.NewFromFen("1q6/2k1b3/8/8/8/5P2/3KN3/7R w - - 0 1"), game.Black}, map[dragontoothmg.Piece]int{dragontoothmg.Pawn: 0, dragontoothmg.Knight: 0, dragontoothmg.Bishop: 330, dragontoothmg.Rook: 0, dragontoothmg.Queen: 900, dragontoothmg.King: 0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculatePieceScore(tt.args.g, tt.args.color); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("calculatePieceScore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvaluation_updateTotal(t *testing.T) {
	type fields struct {
		Game                  *game.Game
		White                 EvaluationPart
		Black                 EvaluationPart
		GameOver              bool
		TotalScore            int
		TotalScorePerspective int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			"GameStart",
			fields{
				Game: game.NewFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"),
				White: EvaluationPart{
					GamePhase:  12,
					PieceScore: map[dragontoothmg.Piece]int{dragontoothmg.Pawn: 800, dragontoothmg.Knight: 640, dragontoothmg.Bishop: 660, dragontoothmg.Rook: 1000, dragontoothmg.Queen: 900, dragontoothmg.King: 0},
				},
				Black: EvaluationPart{
					GamePhase:  12,
					PieceScore: map[dragontoothmg.Piece]int{dragontoothmg.Pawn: 800, dragontoothmg.Knight: 640, dragontoothmg.Bishop: 660, dragontoothmg.Rook: 1000, dragontoothmg.Queen: 900, dragontoothmg.King: 0},
				},
			},
			0,
		},
		{
			"KRNPvKQB",
			fields{
				Game: game.NewFromFen("1q6/2k1b3/8/8/8/5P2/3KN3/7R b - - 0 1"),
				White: EvaluationPart{
					GamePhase:  12,
					PieceScore: map[dragontoothmg.Piece]int{dragontoothmg.Pawn: 100, dragontoothmg.Knight: 320, dragontoothmg.Bishop: 0, dragontoothmg.Rook: 500, dragontoothmg.Queen: 0, dragontoothmg.King: 0},
				},
				Black: EvaluationPart{
					GamePhase:  12,
					PieceScore: map[dragontoothmg.Piece]int{dragontoothmg.Pawn: 0, dragontoothmg.Knight: 0, dragontoothmg.Bishop: 330, dragontoothmg.Rook: 0, dragontoothmg.Queen: 900, dragontoothmg.King: 0},
				},
			},
			-310,
		},
		{
			"KQvK Checkmate White",
			fields{
				Game: game.NewFromFen("7k/6Q1/5K2/8/8/8/8/8 b - - 0 1"),
				White: EvaluationPart{
					GamePhase:  4,
					PieceScore: map[dragontoothmg.Piece]int{dragontoothmg.Pawn: 0, dragontoothmg.Knight: 0, dragontoothmg.Bishop: 0, dragontoothmg.Rook: 0, dragontoothmg.Queen: 900, dragontoothmg.King: 0},
				},
				Black: EvaluationPart{
					GamePhase:  0,
					PieceScore: map[dragontoothmg.Piece]int{dragontoothmg.Pawn: 0, dragontoothmg.Knight: 0, dragontoothmg.Bishop: 0, dragontoothmg.Rook: 0, dragontoothmg.Queen: 0, dragontoothmg.King: 0},
				},
				GameOver: true,
			},
			int(^uint(0) >> 1),
		},
		{
			"KQvK Checkmate Black",
			fields{
				Game: game.NewFromFen("8/8/8/8/8/2k5/1q6/K7 w - - 0 1"),
				White: EvaluationPart{
					GamePhase:  0,
					PieceScore: map[dragontoothmg.Piece]int{dragontoothmg.Pawn: 0, dragontoothmg.Knight: 0, dragontoothmg.Bishop: 0, dragontoothmg.Rook: 0, dragontoothmg.Queen: 0, dragontoothmg.King: 0},
				},
				Black: EvaluationPart{
					GamePhase:  4,
					PieceScore: map[dragontoothmg.Piece]int{dragontoothmg.Pawn: 0, dragontoothmg.Knight: 0, dragontoothmg.Bishop: 0, dragontoothmg.Rook: 0, dragontoothmg.Queen: 900, dragontoothmg.King: 0},
				},
				GameOver: true,
			},
			-int(^uint(0)>>1) - 1,
		},
		{
			"KvKR Stalemate",
			fields{
				Game: game.NewFromFen("8/8/8/8/8/2k5/1r6/K7 w - - 0 1"),
				White: EvaluationPart{
					GamePhase:  0,
					PieceScore: map[dragontoothmg.Piece]int{dragontoothmg.Pawn: 0, dragontoothmg.Knight: 0, dragontoothmg.Bishop: 0, dragontoothmg.Rook: 0, dragontoothmg.Queen: 0, dragontoothmg.King: 0},
				},
				Black: EvaluationPart{
					GamePhase:  2,
					PieceScore: map[dragontoothmg.Piece]int{dragontoothmg.Pawn: 0, dragontoothmg.Knight: 0, dragontoothmg.Bishop: 0, dragontoothmg.Rook: 500, dragontoothmg.Queen: 0, dragontoothmg.King: 0},
				},
				GameOver: true,
			},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Evaluation{
				Game:                  tt.fields.Game,
				White:                 tt.fields.White,
				Black:                 tt.fields.Black,
				GameOver:              tt.fields.GameOver,
				TotalScore:            tt.fields.TotalScore,
				TotalScorePerspective: tt.fields.TotalScorePerspective,
			}
			e.updateTotal()
			if !reflect.DeepEqual(e.TotalScore, tt.want) {
				t.Errorf("Evaluation.updateTotal() = %v, want %v", e.TotalScore, tt.want)
			}
		})
	}
}

func TestCalculateEvaluation(t *testing.T) {
	game_gamestart := game.NewFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	tests := []struct {
		name string
		args *game.Game
		want *Evaluation
	}{
		{
			"GameStart",
			game_gamestart,
			&Evaluation{
				Game: game_gamestart,
				White: EvaluationPart{
					GamePhase:               12,
					PieceScore:              map[dragontoothmg.Piece]int{dragontoothmg.Pawn: 800, dragontoothmg.Knight: 640, dragontoothmg.Bishop: 660, dragontoothmg.Rook: 1000, dragontoothmg.Queen: 900, dragontoothmg.King: 0},
					PieceSquareScoreMidgame: make(map[dragontoothmg.Piece]int),
					PieceSquareScoreEndgame: make(map[dragontoothmg.Piece]int),
				},
				Black: EvaluationPart{
					GamePhase:               12,
					PieceScore:              map[dragontoothmg.Piece]int{dragontoothmg.Pawn: 800, dragontoothmg.Knight: 640, dragontoothmg.Bishop: 660, dragontoothmg.Rook: 1000, dragontoothmg.Queen: 900, dragontoothmg.King: 0},
					PieceSquareScoreMidgame: make(map[dragontoothmg.Piece]int),
					PieceSquareScoreEndgame: make(map[dragontoothmg.Piece]int),
				},
				GameOver:              false,
				TotalScore:            0,
				TotalScorePerspective: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalculateEvaluation(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CalculateEvaluation() = %v, want %v", got, tt.want)
			}
		})
	}
}
