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

func Test_calculateMaterialScoreForPieceType(t *testing.T) {
	type args struct {
		g         *game.Game
		color     game.PlayerColor
		pieceType dragontoothmg.Piece
	}
	tests := []struct {
		name       string
		args       args
		wantPs     int
		wantPstMid int
		wantPstEnd int
	}{
		{"GameStart White", args{game.NewFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"), game.White, dragontoothmg.Pawn}, 800, -66, -66},
		{"GameStart Black", args{game.NewFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"), game.Black, dragontoothmg.Pawn}, 800, -66, -66},
		{"KRNPvKQB White P", args{game.NewFromFen("1q6/2k1b3/8/8/8/5P2/3KN3/7R w - - 0 1"), game.White, dragontoothmg.Pawn}, 100, 1, 1},
		{"KRNPvKQB White N", args{game.NewFromFen("1q6/2k1b3/8/8/8/5P2/3KN3/7R w - - 0 1"), game.White, dragontoothmg.Knight}, 320, 2, 2},
		{"KRNPvKQB White B", args{game.NewFromFen("1q6/2k1b3/8/8/8/5P2/3KN3/7R w - - 0 1"), game.White, dragontoothmg.Bishop}, 0, 0, 0},
		{"KRNPvKQB White R", args{game.NewFromFen("1q6/2k1b3/8/8/8/5P2/3KN3/7R w - - 0 1"), game.White, dragontoothmg.Rook}, 500, 0, 0},
		{"KRNPvKQB White Q", args{game.NewFromFen("1q6/2k1b3/8/8/8/5P2/3KN3/7R w - - 0 1"), game.White, dragontoothmg.Queen}, 0, 0, 0},
		{"KRNPvKQB White K", args{game.NewFromFen("1q6/2k1b3/8/8/8/5P2/3KN3/7R w - - 0 1"), game.White, dragontoothmg.King}, 0, 0, 0},
		{"KRNPvKQB Black P", args{game.NewFromFen("1q6/2k1b3/8/8/8/5P2/3KN3/7R w - - 0 1"), game.Black, dragontoothmg.Pawn}, 0, 0, 0},
		{"KRNPvKQB Black N", args{game.NewFromFen("1q6/2k1b3/8/8/8/5P2/3KN3/7R w - - 0 1"), game.Black, dragontoothmg.Knight}, 0, 0, 0},
		{"KRNPvKQB Black B", args{game.NewFromFen("1q6/2k1b3/8/8/8/5P2/3KN3/7R w - - 0 1"), game.Black, dragontoothmg.Bishop}, 330, 1, 1},
		{"KRNPvKQB Black R", args{game.NewFromFen("1q6/2k1b3/8/8/8/5P2/3KN3/7R w - - 0 1"), game.Black, dragontoothmg.Rook}, 0, 0, 0},
		{"KRNPvKQB Black Q", args{game.NewFromFen("1q6/2k1b3/8/8/8/5P2/3KN3/7R w - - 0 1"), game.Black, dragontoothmg.Queen}, 900, -5, -5},
		{"KRNPvKQB Black K", args{game.NewFromFen("1q6/2k1b3/8/8/8/5P2/3KN3/7R w - - 0 1"), game.Black, dragontoothmg.King}, 0, 20, -12},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bitboard := getBitboardByPieceType(&tt.args.g.Position.White, tt.args.pieceType)
			if tt.args.color == game.Black {
				bitboard = getBitboardByPieceType(&tt.args.g.Position.Black, tt.args.pieceType)
			}
			gotPs, gotPstMid, gotPstEnd := calculateMaterialScoreForPieceType(tt.args.g, tt.args.color, tt.args.pieceType, bitboard)
			if gotPs != tt.wantPs {
				t.Errorf("calculateMaterialScoreForPieceType() ps = %v, want %v", gotPs, tt.wantPs)
			}
			if gotPstMid != tt.wantPstMid {
				t.Errorf("calculateMaterialScoreForPieceType() pstMid = %v, want %v", gotPstMid, tt.wantPstMid)
			}
			if gotPstEnd != tt.wantPstEnd {
				t.Errorf("calculateMaterialScoreForPieceType() pstEnd = %v, want %v", gotPstEnd, tt.wantPstEnd)
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
					PieceScore: 4000,
				},
				Black: EvaluationPart{
					GamePhase:  12,
					PieceScore: 4000,
				},
			},
			0,
		},
		{
			"KRNPvKQB",
			fields{
				Game: game.NewFromFen("1q6/2k1b3/8/8/8/5P2/3KN3/7R b - - 0 1"),
				White: EvaluationPart{
					GamePhase:               12,
					PieceScore:              920,
					PieceSquareScoreMidgame: 3,
					PieceSquareScoreEndgame: 3,
				},
				Black: EvaluationPart{
					GamePhase:               12,
					PieceScore:              1230,
					PieceSquareScoreMidgame: 16,
					PieceSquareScoreEndgame: 16,
				},
			},
			-323,
		},
		{
			"KQvK Checkmate White",
			fields{
				Game: game.NewFromFen("7k/6Q1/5K2/8/8/8/8/8 b - - 0 1"),
				White: EvaluationPart{
					GamePhase:  4,
					PieceScore: 0,
				},
				Black: EvaluationPart{
					GamePhase:  0,
					PieceScore: 900,
				},
				GameOver: true,
			},
			1000000,
		},
		{
			"KQvK Checkmate Black",
			fields{
				Game: game.NewFromFen("8/8/8/8/8/2k5/1q6/K7 w - - 0 1"),
				White: EvaluationPart{
					GamePhase:  0,
					PieceScore: 900,
				},
				Black: EvaluationPart{
					GamePhase:  4,
					PieceScore: 0,
				},
				GameOver: true,
			},
			-1000000,
		},
		{
			"KvKR Stalemate",
			fields{
				Game: game.NewFromFen("8/8/8/8/8/2k5/1r6/K7 w - - 0 1"),
				White: EvaluationPart{
					GamePhase:  0,
					PieceScore: 0,
				},
				Black: EvaluationPart{
					GamePhase:  2,
					PieceScore: 500,
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
		want int
	}{
		{
			"GameStart",
			game_gamestart,
			10,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			eval := Evaluation{}
			if got := eval.CalculateEvaluation(tt.args); got != tt.want {
				t.Errorf("CalculateEvaluation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculatePairModifier(t *testing.T) {
	type args struct {
		g     *game.Game
		color game.PlayerColor
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"GameStart White", args{game.NewFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"), game.White}, 6},
		{"GameStart Black", args{game.NewFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"), game.Black}, 6},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculatePairModifier(tt.args.g, tt.args.color); got != tt.want {
				t.Errorf("calculatePairModifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateTempoModifier(t *testing.T) {
	type args struct {
		g     *game.Game
		color game.PlayerColor
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"GameStart White", args{game.NewFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"), game.White}, 10},
		{"GameStart Black", args{game.NewFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"), game.Black}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateTempoModifier(tt.args.g, tt.args.color); got != tt.want {
				t.Errorf("calculateTempoModifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateRookModifier(t *testing.T) {
	type args struct {
		g     *game.Game
		color game.PlayerColor
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"GameStart White", args{game.NewFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"), game.White}, 0},
		{"GameStart Black", args{game.NewFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"), game.Black}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateRookModifier(tt.args.g, tt.args.color); got != tt.want {
				t.Errorf("calculateRookModifier() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculatePassedPawns(t *testing.T) {
	type args struct {
		g     *game.Game
		color game.PlayerColor
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"GameStart White", args{game.NewFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"), game.White}, 0},
		{"GameStart Black", args{game.NewFromFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"), game.Black}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculatePassedPawns(tt.args.g, tt.args.color); got != tt.want {
				t.Errorf("calculatePassedPawns() = %v, want %v", got, tt.want)
			}
		})
	}
}
