package game

import (
	"errors"
	"math/bits"

	"github.com/dylhunn/dragontoothmg"
)

type PlayerColor bool

const (
	Black = false
	White = true
)

type GameResult uint8

const (
	GameNotOver GameResult = iota
	BlackWon
	Draw
	WhiteWon
)

type DrawReason uint8

const (
	NoDraw DrawReason = iota
	Stalemate
	ThreefoldRepetition
	FivefoldRepetition // Not used, because Threefold triggers first
	FiftyMoveRule
	SeventyFiveMoveRule
	InsufficientMaterial
)

type UnapplyGameState struct {
	Result     GameResult
	DrawReason DrawReason
}

type Game struct {
	Moves                []dragontoothmg.Move
	UnapplyMoveFunctions []func()
	UnapplyGameStates    []UnapplyGameState
	HashHistory          []uint64
	Position             *dragontoothmg.Board
	Result               GameResult
	DrawReason           DrawReason
}

func New() *Game {
	board := dragontoothmg.ParseFen(dragontoothmg.Startpos)
	game := &Game{
		Moves:                []dragontoothmg.Move{},
		UnapplyMoveFunctions: []func(){},
		UnapplyGameStates:    []UnapplyGameState{},
		HashHistory:          []uint64{board.Hash()},
		Position:             &board,
		Result:               GameNotOver,
		DrawReason:           NoDraw,
	}
	game.updateGameState()

	return game
}

func NewFromFen(fen string) *Game {
	board := dragontoothmg.ParseFen(fen)
	game := &Game{
		Moves:                []dragontoothmg.Move{},
		UnapplyMoveFunctions: []func(){},
		UnapplyGameStates:    []UnapplyGameState{},
		HashHistory:          []uint64{board.Hash()},
		Position:             &board,
		Result:               GameNotOver,
		DrawReason:           NoDraw,
	}
	game.updateGameState()

	return game
}

func (g *Game) updateGameState() {
	g.updateGameOverState()
}

func (g *Game) updateGameOverState() {
	g.Result = GameNotOver
	g.DrawReason = NoDraw
	moves := g.Position.GenerateLegalMoves()

	if len(moves) == 0 {
		if g.Position.UnderDirectAttack(true, uint8(bits.TrailingZeros64(g.Position.White.Kings))) {
			g.Result = BlackWon
		} else if g.Position.UnderDirectAttack(false, uint8(bits.TrailingZeros64(g.Position.Black.Kings))) {
			g.Result = WhiteWon
		} else {
			g.Result = Draw
			g.DrawReason = Stalemate
		}
		return
	}

	num_rep := g.numRepetitionsOfCurrentPosition()
	if num_rep >= 3 {
		g.Result = Draw
		g.DrawReason = ThreefoldRepetition
	}

	if g.Position.Halfmoveclock >= 100 {
		g.Result = Draw
		g.DrawReason = FiftyMoveRule
	}

	if g.isDrawnByInsufficientMaterial() {
		g.Result = Draw
		g.DrawReason = InsufficientMaterial
	}
}

func (g *Game) numRepetitionsOfCurrentPosition() int {
	var count int

	for _, hash := range g.HashHistory {
		if hash == g.Position.Hash() {
			count++
		}
	}

	return count
}

func (g *Game) isDrawnByInsufficientMaterial() bool {
	if (g.Position.White.Queens |
		g.Position.White.Rooks |
		g.Position.White.Pawns |
		g.Position.Black.Queens |
		g.Position.Black.Rooks |
		g.Position.Black.Pawns) > 0 {
		return false
	}

	if g.Position.White.Kings == 0 || g.Position.Black.Kings == 0 {
		return false
	}

	knightCount := bits.OnesCount64(g.Position.White.Knights) + bits.OnesCount64(g.Position.Black.Knights)
	bishopCount := bits.OnesCount64(g.Position.White.Bishops) + bits.OnesCount64(g.Position.Black.Bishops)

	if knightCount == 0 && bishopCount == 0 {
		return true
	}

	if knightCount == 1 && bishopCount == 0 {
		return true
	}

	if knightCount == 0 && bishopCount == 1 {
		return true
	}

	if knightCount == 0 {
		bishops := g.Position.White.Bishops | g.Position.Black.Bishops
		whiteSquareCount := 0
		blackSquareCount := 0

		for bishops != 0 {
			bishop := uint8(bits.TrailingZeros64(bishops))
			color := uint64(0xAA55AA55AA55AA55>>bishop) & 1
			bishops &= bishops - 1
			if color == 0 {
				whiteSquareCount++
			} else {
				blackSquareCount++
			}
		}

		if whiteSquareCount == 0 || blackSquareCount == 0 {
			return true
		}
	}

	return false
}

func (g *Game) PushMoveStr(moveString string) error {
	move, err := dragontoothmg.ParseMove(moveString)

	if err != nil {
		return err
	}

	return g.PushMove(move)
}

func (g *Game) PushMove(move dragontoothmg.Move) error {
	unapply := g.Position.Apply(move)
	g.Moves = append(g.Moves, move)
	g.UnapplyMoveFunctions = append(g.UnapplyMoveFunctions, unapply)
	g.UnapplyGameStates = append(g.UnapplyGameStates, UnapplyGameState{Result: g.Result, DrawReason: g.DrawReason})
	g.HashHistory = append(g.HashHistory, g.Position.Hash())

	g.updateGameState()

	return nil
}

func (g *Game) PopMove() error {
	if len(g.Moves) <= 0 {
		return errors.New("no move to pop in move history")
	}

	// pop move and unapply
	var unapply func()
	var unapplyGameState UnapplyGameState
	g.Moves = g.Moves[:len(g.Moves)-1]
	unapply, g.UnapplyMoveFunctions = g.UnapplyMoveFunctions[len(g.UnapplyMoveFunctions)-1], g.UnapplyMoveFunctions[:len(g.UnapplyMoveFunctions)-1]
	unapplyGameState, g.UnapplyGameStates = g.UnapplyGameStates[len(g.UnapplyGameStates)-1], g.UnapplyGameStates[:len(g.UnapplyGameStates)-1]
	g.HashHistory = g.HashHistory[:len(g.HashHistory)-1]
	unapply()

	g.Result = unapplyGameState.Result
	g.DrawReason = unapplyGameState.DrawReason

	return nil
}
