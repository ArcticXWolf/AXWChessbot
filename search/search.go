package search

import (
	"context"
	"log"
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

func (s *Search) SearchBestMove(ctx context.Context) (dragontoothmg.Move, int) {
	start := time.Now()

	moves, score := s.iterativeDeepening(ctx)

	s.SearchDone = true
	s.SearchInfo.TotalSearchTime = time.Since(start)

	return moves[len(moves)-1], score
}

func (s *Search) iterativeDeepening(ctx context.Context) ([]dragontoothmg.Move, int) {
	var moves, movesNew []dragontoothmg.Move
	var score, scoreNew int
	var cancelled bool
	moves, score, _ = s.alphaBetaRoot(ctx, 1, moves)
	s.SearchInfo.MaxDepthCompleted = 1

	for i := 2; uint(i) <= s.MaximumDepthAlphaBeta; i++ {
		select {
		case <-ctx.Done():
			return moves, score
		default:
			movesNew, scoreNew, cancelled = s.alphaBetaRoot(ctx, i, moves)
			if !cancelled {
				moves, score = movesNew, scoreNew
				s.SearchInfo.MaxDepthCompleted = uint(i)
			}
		}
	}

	return moves, score
}

func (s *Search) alphaBetaRoot(ctx context.Context, depth int, previousMoves []dragontoothmg.Move) ([]dragontoothmg.Move, int, bool) {
	beta := 1000000000
	alpha := -1000000000
	var move dragontoothmg.Move

	resultMove, resultScore, cancelled := s.alphaBeta(ctx, depth, alpha, beta, move, previousMoves)
	return resultMove, resultScore, cancelled
}

func (s *Search) alphaBeta(ctx context.Context, depthLeft, alpha, beta int, move dragontoothmg.Move, previousMoves []dragontoothmg.Move) ([]dragontoothmg.Move, int, bool) {
	var bestScore int = -1000000000
	var bestMove dragontoothmg.Move
	var moves []dragontoothmg.Move
	var cancelled bool = false

	if depthLeft <= 0 || s.Game.Result != game.GameNotOver {
		score := s.evaluationProvider.CalculateEvaluation(s.Game)
		if s.Game.Result == game.BlackWon || s.Game.Result == game.WhiteWon {
			if score > 0 {
				score -= int(s.MaximumDepthAlphaBeta) - depthLeft // game won, minimize path to victory
			} else {
				score += int(s.MaximumDepthAlphaBeta) - depthLeft // game lost, maximize path for enemy
			}
		}
		return moves, score, cancelled
	}

	s.SearchInfo.NodesTraversed++

	legal_moves := s.getMovesInOrder(depthLeft, previousMoves)

	var lastMove dragontoothmg.Move
moveIterator:
	for _, m := range legal_moves {
		select {
		case <-ctx.Done():
			cancelled = true
			break moveIterator
		default:
			lastMove = m
			s.Game.PushMove(m)

			newMoves, newScore, newCancelled := s.alphaBeta(ctx, depthLeft-1, -beta, -alpha, m, previousMoves)
			newScore = -newScore
			cancelled = cancelled || newCancelled
			//s.logger.Printf("Conc:\t%s%d %d %d %v\n", strings.Repeat("\t", int(s.MaximumDepthAlphaBeta)-depthLeft), newScore, alpha, beta, &m)

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
				break moveIterator
			}
		}

	}

	if bestMove == 0 {
		bestMove = lastMove
	}
	moves = append(moves, bestMove)
	//s.logger.Printf("Node:\t%s%d %d %d %v\n", strings.Repeat("\t", int(s.MaximumDepthAlphaBeta)-depthLeft), bestScore, alpha, beta, &moves)
	return moves, bestScore, cancelled
}

func (s *Search) getMovesInOrder(depthLeft int, previousMoves []dragontoothmg.Move) []dragontoothmg.Move {
	legal_moves := s.Game.Position.GenerateLegalMoves()

	// if len(previousMoves) > depthLeft {
	// 	// Find previous move in move list
	// 	var index int = -1
	// 	for i, m := range legal_moves {
	// 		if previousMoves[depthLeft-1] == m {
	// 			index = i
	// 		}
	// 	}

	// 	// swap if found
	// 	if index > 0 {
	// 		oldFirstMove := legal_moves[0]
	// 		legal_moves[0] = legal_moves[index]
	// 		legal_moves[index] = oldFirstMove
	// 	}
	// }

	return legal_moves
}
