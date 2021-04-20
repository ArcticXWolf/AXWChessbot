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
	NodesTraversed    uint
	QNodesTraversed   uint
	MaxDepthCompleted uint
	CacheHits         uint
	CacheUse          uint
	TotalSearchTime   time.Duration
}

type Search struct {
	Game                   *game.Game
	logger                 *log.Logger
	evaluationProvider     evaluation_provider.EvaluationProvider
	transpositionTable     *TranspositionTable
	SearchDone             bool
	MaximumDepthAlphaBeta  uint
	MaximumDepthQuiescence uint
	SearchInfo             SearchInfo
}

func New(game *game.Game, logger *log.Logger, transpositionTable *TranspositionTable, evaluationProvider evaluation_provider.EvaluationProvider, maxABDepth, maxQDepth uint) *Search {
	return &Search{
		Game:                   game,
		logger:                 logger,
		evaluationProvider:     evaluationProvider,
		transpositionTable:     transpositionTable,
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
	// s.logger.Printf("SearchInfo: %v", s.SearchInfo)

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
	var alphaOriginal int = alpha
	var moves, newMoves []dragontoothmg.Move
	var cancelled, newCancelled bool = false, false
	var score, newScore int

	if depthLeft <= 0 || s.Game.Result != game.GameNotOver {
		// add leaf to nodecount, but do not count it in qnodecount (prevent overlap in both counts)
		s.SearchInfo.NodesTraversed++
		s.SearchInfo.QNodesTraversed--
		score = s.quiescenceSearch(alpha, beta, int(s.MaximumDepthQuiescence))
		if s.Game.Result == game.BlackWon || s.Game.Result == game.WhiteWon {
			if score > 0 {
				score = 1000000 - (int(s.MaximumDepthAlphaBeta) - depthLeft) // game won, minimize path to victory
			} else {
				score = -1000000 + (int(s.MaximumDepthAlphaBeta) - depthLeft) // game lost, maximize path for enemy
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
				moves = append(moves, cacheMove)
				return moves, cacheScore, cancelled
			}
		}
	}

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

			newMoves, newScore, newCancelled = s.alphaBeta(ctx, depthLeft-1, -beta, -alpha, m, previousMoves)
			newScore = -newScore
			cancelled = cancelled || newCancelled
			// s.logger.Printf("Conc:\t%s%d %v\n", strings.Repeat("\t", int(s.MaximumDepthAlphaBeta)-depthLeft), newScore, &m)

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
	// s.logger.Printf("Node:\t%s%d %d %d %v\n", strings.Repeat("\t", int(s.MaximumDepthAlphaBeta)-depthLeft), bestScore, alpha, beta, &moves)
	return moves, bestScore, cancelled
}

func (s *Search) quiescenceSearch(alpha, beta, depthLeft int) int {
	var standPat, score int
	s.SearchInfo.QNodesTraversed++

	standPat = s.evaluationProvider.CalculateEvaluation(s.Game)
	if standPat >= beta {
		return beta
	}
	if standPat > alpha {
		alpha = standPat
	}

	if depthLeft > 0 && s.Game.Result == game.GameNotOver {
		for _, move := range s.getCapturesInOrder() {
			s.Game.PushMove(move)
			score = -s.quiescenceSearch(-beta, -alpha, depthLeft-1)
			s.Game.PopMove()

			if score >= beta {
				return beta
			}
			if score > alpha {
				alpha = score
			}
		}
	}

	return alpha
}

func (s *Search) getMovesInOrder(depthLeft int, previousMoves []dragontoothmg.Move) []dragontoothmg.Move {
	var legal_moves []dragontoothmg.Move = s.Game.Position.GenerateLegalMoves()

	if len(previousMoves) > depthLeft {
		// Find previous move in move list
		var index int = -1
		for i, m := range legal_moves {
			if previousMoves[depthLeft-1] == m {
				index = i
			}
		}

		// swap if found
		if index > 0 {
			legal_moves[0], legal_moves[index] = legal_moves[index], legal_moves[0]
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
		if bitboardsOpponent.All&(1<<move.To()) > 0 {
			captures = append(captures, move)
		}
	}

	sort.Slice(captures, func(i, j int) bool {
		return s.getCaptureMVVLVA(captures[i], bitboardsOwn, bitboardsOpponent) < s.getCaptureMVVLVA(captures[j], bitboardsOwn, bitboardsOpponent)
	})

	return captures
}

func (s *Search) getCaptureMVVLVA(move dragontoothmg.Move, bitboardsOwn dragontoothmg.Bitboards, bitboardsOpponent dragontoothmg.Bitboards) (score int) {
	pieceTypeFrom, _ := s.getPieceTypeAtPosition(move.From(), bitboardsOwn)
	pieceTypeTo, _ := s.getPieceTypeAtPosition(move.To(), bitboardsOpponent)

	return (1200 - s.evaluationProvider.GetPieceTypeValue(pieceTypeTo)) + int(pieceTypeFrom)
}

func (s *Search) getPieceTypeAtPosition(position uint8, bitboards dragontoothmg.Bitboards) (pieceType dragontoothmg.Piece, occupied bool) {
	if bitboards.Pawns&(1<<position) > 0 {
		return dragontoothmg.Pawn, true
	} else if bitboards.Knights&(1<<position) > 0 {
		return dragontoothmg.Knight, true
	} else if bitboards.Bishops&(1<<position) > 0 {
		return dragontoothmg.Bishop, true
	} else if bitboards.Rooks&(1<<position) > 0 {
		return dragontoothmg.Rook, true
	} else if bitboards.Queens&(1<<position) > 0 {
		return dragontoothmg.Queen, true
	} else if bitboards.Kings&(1<<position) > 0 {
		return dragontoothmg.King, true
	}
	return
}
