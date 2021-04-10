package evaluation

import "github.com/dylhunn/dragontoothmg"

type GamephaseWeights struct {
	Material          map[dragontoothmg.Piece]int
	PieceSquareTables map[dragontoothmg.Piece][64]int
}

type Weights struct {
	Midgame GamephaseWeights
	Endgame GamephaseWeights
}

func getWeightsForAllPhases() GamephaseWeights {
	return GamephaseWeights{
		Material: map[dragontoothmg.Piece]int{
			dragontoothmg.Pawn:   100,
			dragontoothmg.Knight: 300,
			dragontoothmg.Bishop: 300,
			dragontoothmg.Rook:   500,
			dragontoothmg.Queen:  900,
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
}

func getMidgameWeights() GamephaseWeights {
	return getWeightsForAllPhases()
}

func getEndgameWeights() GamephaseWeights {
	gpw := getWeightsForAllPhases()

	gpw.PieceSquareTables[dragontoothmg.King] = [64]int{
		40, 50, 30, 10, 10, 30, 50, 40,
		30, 40, 20, 0, 0, 20, 40, 30,
		10, 20, 0, -20, -20, 0, 20, 10,
		0, 10, -10, -30, -30, -10, 10, 0,
		-10, 0, -20, -40, -40, -20, 0, -10,
		-20, -10, -30, -50, -50, -30, -10, -20,
		-30, -20, -40, -60, -60, -40, -20, -30,
		-40, -30, -50, -70, -70, -50, -30, -40,
	}

	return gpw
}

func GetWeights() Weights {
	return Weights{
		Midgame: getMidgameWeights(),
		Endgame: getEndgameWeights(),
	}
}
