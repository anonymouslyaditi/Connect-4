package main

// Bot aims to be competitive: block immediate wins, take immediate wins, else pick center/near-center.

func BotNextMove(g *Game, botPlayer int) int {
	opp := 3 - botPlayer
	cols := g.Cols
	// 1) Winning move for bot
	for c := 0; c < cols; c++ {
		r := firstEmptyRow(g.Board, c)
		if r == -1 {
			continue
		}
		b := cloneBoard(g.Board)
		b[r][c] = botPlayer
		if checkWin(b, r, c, botPlayer) {
			return c
		}
	}
	// 2) Block opponent immediate win
	for c := 0; c < cols; c++ {
		r := firstEmptyRow(g.Board, c)
		if r == -1 {
			continue
		}
		b := cloneBoard(g.Board)
		b[r][c] = opp
		if checkWin(b, r, c, opp) {
			return c
		}
	}
	// 3) Try to create a 3-in-a-row (simple heuristic): prefer center columns
	center := cols / 2
	order := []int{center}
	for i := 1; i <= cols; i++ {
		if center-i >= 0 {
			order = append(order, center-i)
		}
		if center+i < cols {
			order = append(order, center+i)
		}
	}
	for _, c := range order {
		if firstEmptyRow(g.Board, c) != -1 {
			return c
		}
	}
	// fallback
	for c := 0; c < cols; c++ {
		if firstEmptyRow(g.Board, c) != -1 {
			return c
		}
	}
	return 0
}

func firstEmptyRow(b [][]int, col int) int {
	for r := 0; r < len(b); r++ {
		if b[r][col] != 0 {
			continue
		}
	}
	// need bottom-most empty
	for r := len(b) - 1; r >= 0; r-- {
		if b[r][col] == 0 {
			return r
		}
	}
	return -1
}
