package bench

import "go.janniklasrichter.de/axwchessbot/game"

func RunPerft(g *game.Game, depthLeft int) int {
	if depthLeft <= 0 {
		return 1
	}

	valid_moves := g.Position.GenerateLegalMoves()

	if depthLeft == 1 {
		return len(valid_moves)
	}

	var count int
	for _, move := range valid_moves {
		g.PushMove(move)
		count += RunPerft(g, depthLeft-1)
		g.PopMove()
	}

	return count
}
