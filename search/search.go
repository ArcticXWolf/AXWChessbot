package search

import (
	"log"
	"strings"
	"time"

	"github.com/dylhunn/dragontoothmg"
	"go.janniklasrichter.de/axwchessbot/evaluation_provider"
	"go.janniklasrichter.de/axwchessbot/game"
)

type SearchInfo struct {
	NodesTraversed    uint
	QNodesTraversed   uint
	MaxDepthCompleted uint
	TotalSearchTime   time.Duration
}

type Search struct {
	Game                   *game.Game
	logger                 *log.Logger
	evaluationProvider     evaluation_provider.EvaluationProvider
	SearchDone             bool
	MaximumDepthAlphaBeta  uint
	MaximumDepthQuiescence uint
	SearchInfo             SearchInfo
}

func New(game *game.Game, logger *log.Logger, evaluationProvider evaluation_provider.EvaluationProvider, maxABDepth, maxQDepth uint) *Search {
	return &Search{
		Game:                   game,
		logger:                 logger,
		evaluationProvider:     evaluationProvider,
		SearchDone:             false,
		MaximumDepthAlphaBeta:  maxABDepth,
		MaximumDepthQuiescence: maxQDepth,
		SearchInfo:             SearchInfo{},
	}
}

func (s *Search) SearchBestMove() (dragontoothmg.Move, int) {
	start := time.Now()

	moves, score := s.iterativeDeepening()

	s.SearchDone = true
	s.SearchInfo.TotalSearchTime = time.Since(start)

	return moves[len(moves)-1], score
}

func (s *Search) iterativeDeepening() ([]dragontoothmg.Move, int) {
	moves, score := s.alphaBetaRoot(1)
	s.SearchInfo.MaxDepthCompleted = 1

	for i := 2; uint(i) <= s.MaximumDepthAlphaBeta; i++ {
		moves, score = s.alphaBetaRoot(i)
		s.SearchInfo.MaxDepthCompleted = uint(i)
	}

	return moves, score
}

func (s *Search) alphaBetaRoot(depth int) ([]dragontoothmg.Move, int) {
	beta := int(^uint(0) >> 1) // MAX INT
	alpha := -beta - 1         // MIN INT
	var move dragontoothmg.Move
	var previous_moves []dragontoothmg.Move

	resultMove, resultScore := s.alphaBeta(depth, alpha, beta, move, previous_moves)
	s.logger.Printf("Result: %v  %d\n", &resultMove[len(resultMove)-1], resultScore)
	return resultMove, resultScore
}

func (s *Search) alphaBeta(depthLeft, alpha, beta int, move dragontoothmg.Move, previous_moves []dragontoothmg.Move) ([]dragontoothmg.Move, int) {
	var bestScore int = -int(^uint(0)>>1) - 1
	var bestMove dragontoothmg.Move
	var moves []dragontoothmg.Move

	if depthLeft <= 0 || s.Game.Result != game.GameNotOver {
		score := s.evaluationProvider.CalculateEvaluation(s.Game)
		return moves, score
	}

	s.SearchInfo.NodesTraversed++

	legal_moves := s.Game.Position.GenerateLegalMoves()

	var lastMove dragontoothmg.Move
	for _, m := range legal_moves {
		lastMove = m
		s.Game.PushMove(m)

		newMoves, newScore := s.alphaBeta(depthLeft-1, -beta, -alpha, m, previous_moves)
		newScore = -newScore
		s.logger.Printf("Conc:%s%d %v\n", strings.Repeat("\t", 2-depthLeft), newScore, &m)

		s.Game.PopMove()

		if newScore > bestScore {
			moves = newMoves
			bestScore = newScore
			bestMove = m
		}
		if newScore > alpha {
			alpha = newScore
		}
		if alpha >= beta {
			break
		}
	}

	if bestMove == 0 {
		bestMove = lastMove
	}
	moves = append(moves, bestMove)
	s.logger.Printf("Node:%s%d %v\n", strings.Repeat("\t", 2-depthLeft), bestScore, &moves)
	return moves, bestScore
}
