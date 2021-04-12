package bench

import (
	"testing"
	"time"

	"go.janniklasrichter.de/axwchessbot/game"
)

func Test_runPerft(t *testing.T) {
	type args struct {
		g         *game.Game
		depthLeft int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"StartingPosition n=1", args{g: game.New(), depthLeft: 1}, 20},
		{"StartingPosition n=2", args{g: game.New(), depthLeft: 2}, 400},
		{"StartingPosition n=3", args{g: game.New(), depthLeft: 3}, 8902},
		{"StartingPosition n=4", args{g: game.New(), depthLeft: 4}, 197281},
		{"StartingPosition n=5", args{g: game.New(), depthLeft: 5}, 4865609},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := runPerft(tt.args.g, tt.args.depthLeft); got != tt.want {
				t.Errorf("runPerft() = %v, want %v", got, tt.want)
			}
		})
	}
}

func benchmarkRunPerft(i int, b *testing.B) {
	var nodes int
	var nps float64
	var start time.Time

	for n := 0; n < b.N; n++ {
		start = time.Now()
		nodes = runPerft(game.New(), i)
		nps += float64(nodes) / float64(time.Since(start).Seconds())
	}

	b.ReportMetric(float64(nps)/float64(b.N), "nodes/sec/op")
}

func Benchmark_runPerft4(b *testing.B) { benchmarkRunPerft(4, b) }
func Benchmark_runPerft5(b *testing.B) { benchmarkRunPerft(5, b) }
func Benchmark_runPerft6(b *testing.B) { benchmarkRunPerft(6, b) }
func Benchmark_runPerft7(b *testing.B) { benchmarkRunPerft(7, b) }
