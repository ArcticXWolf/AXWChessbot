package search

import (
	"testing"
	"time"

	"go.janniklasrichter.de/axwchessbot/evaluation"
	"go.janniklasrichter.de/axwchessbot/evaluation_null"
	"go.janniklasrichter.de/axwchessbot/evaluation_provider"
	"go.janniklasrichter.de/axwchessbot/game"
)

func benchmarkSearchEvaluation(evaluator evaluation_provider.EvaluationProvider, abdepth uint, b *testing.B) {
	var nps float64
	var start time.Time

	for n := 0; n < b.N; n++ {
		start = time.Now()
		searchObj := New(game.New(), evaluator, abdepth, 4)
		searchObj.SearchBestMove()
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
