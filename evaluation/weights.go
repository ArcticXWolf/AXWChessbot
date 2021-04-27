package evaluation

import (
	"github.com/dylhunn/dragontoothmg"
	"go.janniklasrichter.de/axwchessbot/game"
)

type AdditionalModifier struct {
	BishopPairModifier        int
	KnightPairModifier        int
	RookPairModifier          int
	OpenRookModifier          int
	HalfRookModifier          int
	TempoModifier             int
	KingShieldRank2Modifier   int
	KingShieldRank3Modifier   int
	RookBlockedByKingModifier int
	DiagonalMobilityModifier  int
	LinearMobilityModifier    int
}

type GamephaseWeights struct {
	Material           map[dragontoothmg.Piece]int
	PassedPawnModifier [64]int
	PieceSquareTables  map[dragontoothmg.Piece][64]int
}

type Weights struct {
	Midgame            GamephaseWeights
	Endgame            GamephaseWeights
	AdditionalModifier AdditionalModifier
}

var (
	additionalModifiers = AdditionalModifier{
		BishopPairModifier:        30,
		KnightPairModifier:        -8,
		RookPairModifier:          -16,
		OpenRookModifier:          10,
		HalfRookModifier:          5,
		TempoModifier:             10,
		KingShieldRank2Modifier:   10,
		KingShieldRank3Modifier:   5,
		RookBlockedByKingModifier: 24,
		DiagonalMobilityModifier:  2,
		LinearMobilityModifier:    4,
	}
	weightsForAllPhases = GamephaseWeights{
		Material: map[dragontoothmg.Piece]int{
			dragontoothmg.Pawn:   100,
			dragontoothmg.Knight: 320,
			dragontoothmg.Bishop: 330,
			dragontoothmg.Rook:   500,
			dragontoothmg.Queen:  900,
			dragontoothmg.King:   0,
		},
		PassedPawnModifier: [64]int{
			0, 0, 0, 0, 0, 0, 0, 0,
			10, 10, 10, 10, 10, 10, 10, 10,
			20, 20, 20, 20, 20, 20, 20, 20,
			40, 40, 40, 40, 40, 40, 40, 40,
			60, 60, 60, 60, 60, 60, 60, 60,
			80, 80, 80, 80, 80, 80, 80, 80,
			100, 100, 100, 100, 100, 100, 100, 100,
			0, 0, 0, 0, 0, 0, 0, 0,
		},
		PieceSquareTables: map[dragontoothmg.Piece][64]int{
			dragontoothmg.Pawn: [64]int{
				0, 0, 0, 0, 0, 0, 0, 0,
				-6, -4, 1, -24, -24, 1, -4, -6,
				-4, -4, 1, 5, 5, 1, -4, -4,
				-6, -4, 5, 10, 10, 5, -4, -6,
				-6, -4, 2, 8, 8, 2, -4, -6,
				-6, -4, 1, 2, 2, 1, -4, -6,
				-6, -4, 1, 1, 1, 1, -4, -6,
				0, 0, 0, 0, 0, 0, 0, 0,
			},
			dragontoothmg.Knight: [64]int{
				-8, -12, -8, -8, -8, -8, -12, -8,
				-8, 0, 1, 2, 2, 1, 0, -8,
				-8, 0, 4, 4, 4, 4, 0, -8,
				-8, 0, 4, 8, 8, 4, 0, -8,
				-8, 0, 4, 8, 8, 4, 0, -8,
				-8, 0, 4, 4, 4, 4, 0, -8,
				-8, 0, 0, 0, 0, 0, 0, -8,
				-8, -8, -8, -8, -8, -8, -8, -8,
			},
			dragontoothmg.Bishop: [64]int{
				-4, -4, -12, -4, -4, -12, -4, -4,
				-4, 2, 1, 1, 1, 1, 2, -4,
				-4, 1, 2, 4, 4, 2, 1, -4,
				-4, 0, 4, 6, 6, 4, 0, -4,
				-4, 0, 4, 6, 6, 4, 0, -4,
				-4, 0, 2, 4, 4, 2, 0, -4,
				-4, 0, 0, 0, 0, 0, 0, -4,
				-4, -4, -4, -4, -4, -4, -4, -4,
			},
			dragontoothmg.Rook: [64]int{
				0, 0, 0, 2, 2, 0, 0, 0,
				-5, 0, 0, 0, 0, 0, 0, -5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				-5, 0, 0, 0, 0, 0, 0, -5,
				20, 20, 20, 20, 20, 20, 20, 20,
				5, 5, 5, 5, 5, 5, 5, 5,
			},
			dragontoothmg.Queen: [64]int{
				-5, -5, -5, -5, -5, -5, -5, -5,
				0, 0, 1, 1, 1, 1, 0, 0,
				0, 0, 1, 2, 2, 1, 0, 0,
				0, 0, 2, 3, 3, 2, 0, 0,
				0, 0, 2, 3, 3, 2, 0, 0,
				0, 0, 1, 2, 2, 1, 0, 0,
				0, 0, 1, 1, 1, 1, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
			},
			dragontoothmg.King: [64]int{
				40, 50, 30, 10, 10, 30, 50, 40,
				30, 40, 20, 0, 0, 20, 40, 30,
				10, 20, 0, -20, -20, 0, 20, 10,
				0, 10, -10, -30, -30, -10, 10, 0,
				-10, 0, -20, -40, -40, -20, 0, -10,
				-20, -10, -30, -50, -50, -30, -10, -20,
				-30, -20, -40, -60, -60, -40, -20, -30,
				-40, -30, -50, -70, -70, -50, -30, -40,
			},
		},
	}
)

var (
	midgameWeights = GamephaseWeights{
		Material:           weightsForAllPhases.Material,
		PassedPawnModifier: weightsForAllPhases.PassedPawnModifier,
		PieceSquareTables:  weightsForAllPhases.PieceSquareTables,
	}
	endgameWeights = GamephaseWeights{
		Material:           weightsForAllPhases.Material,
		PassedPawnModifier: weightsForAllPhases.PassedPawnModifier,
		PieceSquareTables: map[dragontoothmg.Piece][64]int{
			dragontoothmg.Pawn:   weightsForAllPhases.PieceSquareTables[dragontoothmg.Pawn],
			dragontoothmg.Knight: weightsForAllPhases.PieceSquareTables[dragontoothmg.Knight],
			dragontoothmg.Bishop: weightsForAllPhases.PieceSquareTables[dragontoothmg.Bishop],
			dragontoothmg.Rook:   weightsForAllPhases.PieceSquareTables[dragontoothmg.Rook],
			dragontoothmg.Queen:  weightsForAllPhases.PieceSquareTables[dragontoothmg.Queen],
			dragontoothmg.King: [64]int{
				-72, -48, -36, -24, -24, -36, -48, -72,
				-48, -24, -12, 0, 0, -12, -24, -48,
				-36, -12, 0, 12, 12, 0, -12, -36,
				-24, 0, 12, 24, 24, 12, 0, -24,
				-24, 0, 12, 24, 24, 12, 0, -24,
				-36, -12, 0, 12, 12, 0, -12, -36,
				-48, -24, -12, 0, 0, -12, -24, -48,
				-72, -48, -36, -24, -24, -36, -48, -72,
			},
		},
	}
)

func flipPstArrayVertically(pst [64]int) [64]int {
	var pstNew [64]int

	for i := 0; i < len(pst); i++ {
		pstNew[i^0x38] = pst[i]
	}

	return pstNew
}

var (
	weights = map[game.PlayerColor]Weights{
		game.White: Weights{
			Midgame:            midgameWeights,
			Endgame:            endgameWeights,
			AdditionalModifier: additionalModifiers,
		},
		game.Black: Weights{
			Midgame: GamephaseWeights{
				Material:           midgameWeights.Material,
				PassedPawnModifier: flipPstArrayVertically(midgameWeights.PassedPawnModifier),
				PieceSquareTables: map[dragontoothmg.Piece][64]int{
					dragontoothmg.Pawn:   flipPstArrayVertically(midgameWeights.PieceSquareTables[dragontoothmg.Pawn]),
					dragontoothmg.Knight: flipPstArrayVertically(midgameWeights.PieceSquareTables[dragontoothmg.Knight]),
					dragontoothmg.Bishop: flipPstArrayVertically(midgameWeights.PieceSquareTables[dragontoothmg.Bishop]),
					dragontoothmg.Rook:   flipPstArrayVertically(midgameWeights.PieceSquareTables[dragontoothmg.Rook]),
					dragontoothmg.Queen:  flipPstArrayVertically(midgameWeights.PieceSquareTables[dragontoothmg.Queen]),
					dragontoothmg.King:   flipPstArrayVertically(midgameWeights.PieceSquareTables[dragontoothmg.King]),
				},
			},
			Endgame: GamephaseWeights{
				Material:           endgameWeights.Material,
				PassedPawnModifier: flipPstArrayVertically(endgameWeights.PassedPawnModifier),
				PieceSquareTables: map[dragontoothmg.Piece][64]int{
					dragontoothmg.Pawn:   flipPstArrayVertically(endgameWeights.PieceSquareTables[dragontoothmg.Pawn]),
					dragontoothmg.Knight: flipPstArrayVertically(endgameWeights.PieceSquareTables[dragontoothmg.Knight]),
					dragontoothmg.Bishop: flipPstArrayVertically(endgameWeights.PieceSquareTables[dragontoothmg.Bishop]),
					dragontoothmg.Rook:   flipPstArrayVertically(endgameWeights.PieceSquareTables[dragontoothmg.Rook]),
					dragontoothmg.Queen:  flipPstArrayVertically(endgameWeights.PieceSquareTables[dragontoothmg.Queen]),
					dragontoothmg.King:   flipPstArrayVertically(endgameWeights.PieceSquareTables[dragontoothmg.King]),
				},
			},
			AdditionalModifier: additionalModifiers,
		},
	}
)
