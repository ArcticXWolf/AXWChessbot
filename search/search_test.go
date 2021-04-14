package search

import (
	"context"
	"log"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/dylhunn/dragontoothmg"
	"go.janniklasrichter.de/axwchessbot/evaluation"
	"go.janniklasrichter.de/axwchessbot/evaluation_null"
	"go.janniklasrichter.de/axwchessbot/evaluation_provider"
	"go.janniklasrichter.de/axwchessbot/game"
)

func benchmarkSearchEvaluation(evaluator evaluation_provider.EvaluationProvider, abdepth uint, b *testing.B) {
	var nps float64
	var start time.Time
	logger := log.New(os.Stderr, "", log.LstdFlags)
	ctx := context.Background()
	transpositionTable := NewTranspositionTable(1000000)

	for n := 0; n < b.N; n++ {
		start = time.Now()
		searchObj := New(game.New(), logger, transpositionTable, evaluator, abdepth, 4)
		searchObj.SearchBestMove(ctx)
		nps += float64(searchObj.SearchInfo.NodesTraversed) / float64(time.Since(start).Seconds())
	}

	b.ReportMetric(float64(nps)/float64(b.N), "nodes/sec/op")
}

func benchmarkSearchNullEvaluation(abdepth uint, b *testing.B) {
	evaluator := evaluation_null.EvaluationNull{}
	benchmarkSearchEvaluation(&evaluator, abdepth, b)
}

func BenchmarkSearchNullEvaluation2(b *testing.B)  { benchmarkSearchNullEvaluation(2, b) }
func BenchmarkSearchNullEvaluation3(b *testing.B)  { benchmarkSearchNullEvaluation(3, b) }
func BenchmarkSearchNullEvaluation4(b *testing.B)  { benchmarkSearchNullEvaluation(4, b) }
func BenchmarkSearchNullEvaluation5(b *testing.B)  { benchmarkSearchNullEvaluation(5, b) }
func BenchmarkSearchNullEvaluation6(b *testing.B)  { benchmarkSearchNullEvaluation(6, b) }
func BenchmarkSearchNullEvaluation7(b *testing.B)  { benchmarkSearchNullEvaluation(7, b) }
func BenchmarkSearchNullEvaluation8(b *testing.B)  { benchmarkSearchNullEvaluation(8, b) }
func BenchmarkSearchNullEvaluation9(b *testing.B)  { benchmarkSearchNullEvaluation(9, b) }
func BenchmarkSearchNullEvaluation10(b *testing.B) { benchmarkSearchNullEvaluation(10, b) }
func BenchmarkSearchNullEvaluation11(b *testing.B) { benchmarkSearchNullEvaluation(11, b) }

func benchmarkSearchFullEvaluation(abdepth uint, b *testing.B) {
	evaluator := evaluation.Evaluation{}
	benchmarkSearchEvaluation(&evaluator, abdepth, b)
}

func BenchmarkSearchFullEvaluation2(b *testing.B)  { benchmarkSearchFullEvaluation(2, b) }
func BenchmarkSearchFullEvaluation3(b *testing.B)  { benchmarkSearchFullEvaluation(3, b) }
func BenchmarkSearchFullEvaluation4(b *testing.B)  { benchmarkSearchFullEvaluation(4, b) }
func BenchmarkSearchFullEvaluation5(b *testing.B)  { benchmarkSearchFullEvaluation(5, b) }
func BenchmarkSearchFullEvaluation6(b *testing.B)  { benchmarkSearchFullEvaluation(6, b) }
func BenchmarkSearchFullEvaluation7(b *testing.B)  { benchmarkSearchFullEvaluation(7, b) }
func BenchmarkSearchFullEvaluation8(b *testing.B)  { benchmarkSearchFullEvaluation(8, b) }
func BenchmarkSearchFullEvaluation9(b *testing.B)  { benchmarkSearchFullEvaluation(9, b) }
func BenchmarkSearchFullEvaluation10(b *testing.B) { benchmarkSearchFullEvaluation(10, b) }
func BenchmarkSearchFullEvaluation11(b *testing.B) { benchmarkSearchFullEvaluation(11, b) }

func TestSearch_SearchBestMove(t *testing.T) {
	type fields struct {
		Game                   *game.Game
		MaximumDepthAlphaBeta  uint
		MaximumDepthQuiescence uint
	}
	tests := []struct {
		name   string
		fields fields
		want   dragontoothmg.Move
		want1  int
	}{
		{"Mate in 1 - Black", fields{game.NewFromFen("1r4k1/p4p1p/p4p2/2pn4/K3b3/3q2n1/3r4/b7 b - - 11 38"), 4, 4}, 712, 999999},
	}
	logger := log.New(os.Stdout, "", log.LstdFlags)
	evaluator := evaluation.Evaluation{}
	transpositionTable := NewTranspositionTable(1000000)
	ctx := context.Background()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(tt.fields.Game, logger, transpositionTable, &evaluator, tt.fields.MaximumDepthAlphaBeta, tt.fields.MaximumDepthQuiescence)
			got, got1 := s.SearchBestMove(ctx)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Search.SearchBestMove() move = %v (%v), want %v (%v)", &got, got, &tt.want, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Search.SearchBestMove() score = %v, want %v", got1, tt.want1)
			}
		})
	}
}
