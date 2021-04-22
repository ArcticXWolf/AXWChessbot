package search

import (
	"context"
	"log"
	"sort"
	"time"

	"github.com/dylhunn/dragontoothmg"
	"go.janniklasrichter.de/axwchessbot/evaluation_provider"
	"go.janniklasrichter.de/axwchessbot/game"
)

type SearchInfo struct {
	NodesTraversed     uint
	QNodesTraversed    uint
	MaxDepthCompleted  uint
	MaxQDepthCompleted uint
	CacheHits          uint
	CacheUse           uint
	SearchTimeStart    time.Time
	SearchDuration     time.Duration
}

type ProtocolLogger interface {
	SendInfo(depth, score, nodes, nps int, time time.Duration, pv []dragontoothmg.Move)
}

type Search struct {
	Game                   *game.Game
	protocol               ProtocolLogger
	logger                 *log.Logger
	evaluationProvider     evaluation_provider.EvaluationProvider
	transpositionTable     *TranspositionTable
	killerMoveTable        *killerMoveTable
	SearchDone             bool
	MaximumDepthAlphaBeta  uint
	MaximumDepthQuiescence uint
	SearchInfo             SearchInfo
}

func New(game *game.Game, protocol ProtocolLogger, logger *log.Logger, transpositionTable *TranspositionTable, evaluationProvider evaluation_provider.EvaluationProvider, maxABDepth, maxQDepth uint) *Search {
	return &Search{
		Game:                   game,
		protocol:               protocol,
		logger:                 logger,
		evaluationProvider:     evaluationProvider,
		transpositionTable:     transpositionTable,
		killerMoveTable:        newKillerMoveTable(int(maxABDepth)),
		SearchDone:             false,
		MaximumDepthAlphaBeta:  maxABDepth,
		MaximumDepthQuiescence: maxQDepth,
		SearchInfo:             SearchInfo{},
	}
}

func (s *Search) SearchBestMove(ctx context.Context) (dragontoothmg.Move, int) {
	s.SearchInfo.SearchTimeStart = time.Now()

	moves, score := s.iterativeDeepening(ctx)

	s.SearchDone = true
	s.SearchInfo.SearchDuration = time.Since(s.SearchInfo.SearchTimeStart)
	s.protocol.SendInfo(
		int(s.SearchInfo.MaxDepthCompleted),
		score,
		int(s.SearchInfo.NodesTraversed+s.SearchInfo.QNodesTraversed),
		int(float64(s.SearchInfo.NodesTraversed+s.SearchInfo.QNodesTraversed)/s.SearchInfo.SearchDuration.Seconds()),
		time.Duration(s.SearchInfo.SearchDuration.Milliseconds()),
		moves,
	)

	return moves[len(moves)-1], score
}

func (s *Search) iterativeDeepening(ctx context.Context) ([]dragontoothmg.Move, int) {
	var moves, movesNew []dragontoothmg.Move
	var score, scoreNew int
	var cancelled bool
	var qdepth int = 0
	moves, score, _ = s.alphaBetaRoot(ctx, 1, qdepth, moves)
	s.SearchInfo.MaxDepthCompleted = 1
	s.SearchInfo.SearchDuration = time.Since(s.SearchInfo.SearchTimeStart)
	s.protocol.SendInfo(
		int(s.SearchInfo.MaxDepthCompleted),
		score,
		int(s.SearchInfo.NodesTraversed),
		int(float64(s.SearchInfo.NodesTraversed)/s.SearchInfo.SearchDuration.Seconds()),
		time.Duration(s.SearchInfo.SearchDuration.Milliseconds()),
		moves,
	)

	for i := 2; uint(i) <= s.MaximumDepthAlphaBeta; i++ {
		select {
		case <-ctx.Done():
			return moves, score
		default:
			qdepth = int(i / 2)
			if qdepth > int(s.MaximumDepthQuiescence) {
				qdepth = int(s.MaximumDepthQuiescence)
			}
			movesNew, scoreNew, cancelled = s.alphaBetaRoot(ctx, i, qdepth, moves)
			if !cancelled {
				moves, score = movesNew, scoreNew
				s.SearchInfo.MaxDepthCompleted = uint(i)
				s.SearchInfo.SearchDuration = time.Since(s.SearchInfo.SearchTimeStart)
				s.protocol.SendInfo(
					int(s.SearchInfo.MaxDepthCompleted),
					score,
					int(s.SearchInfo.NodesTraversed),
					int(float64(s.SearchInfo.NodesTraversed)/s.SearchInfo.SearchDuration.Seconds()),
					time.Duration(s.SearchInfo.SearchDuration.Milliseconds()),
					moves,
				)
			}
		}
	}

	return moves, score
}

func (s *Search) alphaBetaRoot(ctx context.Context, depth, qdepth int, previousMoves []dragontoothmg.Move) ([]dragontoothmg.Move, int, bool) {
	beta := 1000000000
	alpha := -1000000000
	var move dragontoothmg.Move

	s.killerMoveTable.clear()
	resultMove, resultScore, cancelled := s.alphaBeta(ctx, depth, qdepth, 0, alpha, beta, move, previousMoves)
	return resultMove, resultScore, cancelled
}

func (s *Search) alphaBeta(ctx context.Context, depthLeft, qdepth, ply, alpha, beta int, move dragontoothmg.Move, previousMoves []dragontoothmg.Move) ([]dragontoothmg.Move, int, bool) {
	var bestScore int = -1000000000
	var bestMove dragontoothmg.Move
	var alphaOriginal int = alpha
	var moves, newMoves []dragontoothmg.Move
	var cancelled, newCancelled bool = false, false
	var score, newScore int

	if depthLeft <= 0 || s.Game.Result != game.GameNotOver {
		// add leaf to nodecount, but do not count it in qnodecount (prevent overlap in both counts)
		s.SearchInfo.NodesTraversed++
		s.SearchInfo.QNodesTraversed--
		score = s.quiescenceSearch(alpha, beta, qdepth, ply)
		if s.Game.Result == game.BlackWon || s.Game.Result == game.WhiteWon {
			if score > 0 {
				score = 1000000 - ply // game won, minimize path to victory
			} else {
				score = -1000000 + ply // game lost, maximize path for enemy
			}
		}
		return moves, score, cancelled
	}

	s.SearchInfo.NodesTraversed++

	cacheMove, cacheScore, cacheDepth, cacheBound, cacheFound := s.transpositionTable.Load(s.Game.Position.Hash())
	if cacheFound {
		s.SearchInfo.CacheHits++
		if cacheDepth >= depthLeft {
			if cacheBound == alphaBetaBoundExact {
				s.SearchInfo.CacheUse++
				if !isCaptureOrPromotionMove(s.Game, cacheMove) {
					s.killerMoveTable.update(ply, cacheMove)
				}
				moves = append(moves, cacheMove)
				return moves, cacheScore, cancelled
			} else if cacheBound == alphaBetaBoundLower {
				if cacheScore > alpha {
					alpha = cacheScore
				}
			} else if cacheBound == alphaBetaBoundUpper {
				if cacheScore < beta {
					beta = cacheScore
				}
			}
			if alpha >= beta {
				s.SearchInfo.CacheUse++
				if !isCaptureOrPromotionMove(s.Game, cacheMove) {
					s.killerMoveTable.update(ply, cacheMove)
				}
				moves = append(moves, cacheMove)
				return moves, cacheScore, cancelled
			}
		}
	}

	legal_moves := s.getMovesInOrder(ply, depthLeft, previousMoves)

	var lastMove dragontoothmg.Move
moveIterator:
	for _, m := range legal_moves {
		lastMove = m
		select {
		case <-ctx.Done():
			cancelled = true
			break moveIterator
		default:
			s.Game.PushMove(m)

			newMoves, newScore, newCancelled = s.alphaBeta(ctx, depthLeft-1, qdepth, ply+1, -beta, -alpha, m, previousMoves)
			newScore = -newScore
			cancelled = cancelled || newCancelled
			//s.logger.Printf("Conc: %11d a:%11d b:%11d \t%s%v (%d)\n", bestScore, alpha, beta, strings.Repeat("\t", ply), &m, newScore)

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
				if !isCaptureOrPromotionMove(s.Game, m) {
					s.killerMoveTable.update(ply, m)
				}
				break moveIterator
			}
		}

	}

	if bestMove == 0 {
		bestMove = lastMove
	}

	if !cancelled {
		var bound alphaBetaBound = alphaBetaBoundExact
		if bestScore <= alphaOriginal {
			bound = alphaBetaBoundUpper
		} else if bestScore >= beta {
			bound = alphaBetaBoundLower
		}
		s.transpositionTable.InsertIfNeeded(s.Game.Position.Hash(), bestMove, bestScore, depthLeft, bound)
	}

	moves = append(moves, bestMove)
	//s.logger.Printf("Node: %11d a:%11d b:%11d \t%s%v\n", bestScore, alpha, beta, strings.Repeat("\t", ply), &moves)
	return moves, bestScore, cancelled
}

func (s *Search) quiescenceSearch(alpha, beta, depthLeft, ply int) int {
	var standPat, score int
	s.SearchInfo.QNodesTraversed++

	standPat = s.evaluationProvider.CalculateEvaluation(s.Game)
	if depthLeft <= 0 {
		return standPat
	}
	if standPat >= beta {
		return beta
	}
	if standPat > alpha {
		alpha = standPat
	}

	if s.Game.Result == game.BlackWon || s.Game.Result == game.WhiteWon {
		if score > 0 {
			return 1000000 - ply // game won, minimize path to victory
		} else {
			return -1000000 + ply // game lost, maximize path for enemy
		}
	}

	for _, move := range s.getCapturesInOrder() {
		s.Game.PushMove(move)
		score = -s.quiescenceSearch(-beta, -alpha, depthLeft-1, ply+1)
		//s.logger.Printf("Qonc: %11s a:%11d b:%11d \t%s%v (%d)\n", "-------", alpha, beta, strings.Repeat("\t", ply+1), &move, score)
		s.Game.PopMove()

		if score >= beta {
			return beta
		}
		if score > alpha {
			alpha = score
		}
	}

	//s.logger.Printf("Node: %11s a:%11d b:%11d \t%s%d\n", "-------", alpha, beta, strings.Repeat("\t", ply+1), alpha)
	return alpha
}

func (s *Search) getMovesInOrder(ply int, depthLeft int, previousMoves []dragontoothmg.Move) []dragontoothmg.Move {
	var legal_moves []dragontoothmg.Move = s.Game.Position.GenerateLegalMoves()

	if len(previousMoves) > depthLeft {
		var index int = -1
		for i, m := range legal_moves {
			if previousMoves[depthLeft-1] == m {
				index = i
			}
		}

		if index > 0 {
			legal_moves[0], legal_moves[index] = legal_moves[index], legal_moves[0]
		}
	}

	move1, found1, move2, found2 := s.killerMoveTable.fetch(ply)
	if found1 {
		var index int = -1
		for i, m := range legal_moves {
			if move1 == m {
				index = i
			}
		}

		if index > 1 {
			legal_moves[1], legal_moves[index] = legal_moves[index], legal_moves[1]
		}
	}
	if found2 {
		var index int = -1
		for i, m := range legal_moves {
			if move2 == m {
				index = i
			}
		}

		if index > 2 {
			legal_moves[2], legal_moves[index] = legal_moves[index], legal_moves[2]
		}
	}

	return legal_moves
}

func (s *Search) getCapturesInOrder() []dragontoothmg.Move {
	var captures []dragontoothmg.Move = []dragontoothmg.Move{}

	bitboardsOwn := s.Game.Position.White
	bitboardsOpponent := s.Game.Position.Black
	if !s.Game.Position.Wtomove {
		bitboardsOwn = s.Game.Position.Black
		bitboardsOpponent = s.Game.Position.White
	}

	for _, move := range s.Game.Position.GenerateLegalMoves() {
		if isCaptureOrPromotionMove(s.Game, move) {
			captures = append(captures, move)
		}
	}

	sort.Slice(captures, func(i, j int) bool {
		return s.getCaptureMVVLVA(captures[i], bitboardsOwn, bitboardsOpponent) < s.getCaptureMVVLVA(captures[j], bitboardsOwn, bitboardsOpponent)
	})

	return captures
}
