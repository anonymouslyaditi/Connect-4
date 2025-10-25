package main

import (
	"time"
)

// We'll keep models small and JSON friendly

type Game struct {
	Rows    int         `json:"rows"`
	Cols    int         `json:"cols"`
	Board   [][]int     `json:"board"` // [row][col] 0-empty,1,2
	Turn    int         `json:"turn"`  // which player's turn (1 or 2)
	Started time.Time   `json:"started"`
}

type GameRecord struct {
	ID        string    `json:"id"`
	Player1   string    `json:"player1"`
	Player2   string    `json:"player2"`
	Winner    string    `json:"winner"` // "draw" or username
	Duration  int64     `json:"duration_seconds"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
}

type Leaderboard map[string]int
