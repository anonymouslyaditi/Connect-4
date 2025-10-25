package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"sync"
	"time"
)

// GameSession handles a match between two players (or bot)

type GameSession struct {
	ID         string
	Player1    string
	Player2    string
	Players    map[string]int // username -> 1 or 2
	Game       *Game
	TurnMu     sync.Mutex
	State      string // playing, finished
	Result     string // "p1","p2","draw"
	StartedAt  time.Time
	FinishedAt time.Time
	IsBot      bool
	clients    map[string]*Client
}

func NewGameSession(p1, p2 string) *GameSession {
	id := fmt.Sprintf("g_%d", time.Now().UnixNano())
	rows := 6
	cols := 7
	board := make([][]int, rows)
	for r := 0; r < rows; r++ {
		board[r] = make([]int, cols)
	}
	g := &Game{Rows: rows, Cols: cols, Board: board, Turn: 1, Started: time.Now()}
	return &GameSession{ID: id, Player1: p1, Player2: p2, Players: map[string]int{p1: 1, p2: 2}, Game: g, State: "playing", StartedAt: time.Now(), clients: map[string]*Client{}}
}

func (s *GameSession) run() {
	// If bot and bot moves first when it's player 2, process it
	for s.State == "playing" {
		// if bot present and it's bot's turn, make a bot move
		if s.IsBot && s.Game.Turn == s.getBotPlayer() {
			col := BotNextMove(s.Game, s.getBotPlayer())
			s.applyMove(col)
		} else {
			// wait for moves via WebSocket (client readPump will call applyMove)
			time.Sleep(200 * time.Millisecond)
		}
	}
}

func (s *GameSession) getBotPlayer() int {
	if s.Player2 == "Bot" {
		return 2
	}
	return 1
}

func (s *GameSession) reconnect(username string, client *Client) {
	s.clients[username] = client
	client.SendJSON(map[string]interface{}{"type": "reconnected", "gameId": s.ID, "state": s.Game})
}

func (s *GameSession) applyMove(col int) {
	s.TurnMu.Lock()
	defer s.TurnMu.Unlock()
	if s.State != "playing" {
		return
	}
	if col < 0 || col >= s.Game.Cols {
		return
	}
	// drop
	for r := s.Game.Rows - 1; r >= 0; r-- {
		if s.Game.Board[r][col] == 0 {
			s.Game.Board[r][col] = s.Game.Turn
			// check win
			if checkWin(s.Game.Board, r, col, s.Game.Turn) {
				// finish
				winner := ""
				if s.Game.Turn == 1 {
					winner = s.Player1
				} else {
					winner = s.Player2
				}
				fmt.Printf("WIN DETECTED! Winner: %s\n", winner)
				s.State = "finished"
				s.Result = winner
				fmt.Printf("Set s.Result to: %s\n", s.Result)
				s.FinishedAt = time.Now()
				// persist
				rec := GameRecord{
					ID:        s.ID,
					Player1:   s.Player1,
					Player2:   s.Player2,
					Winner:    winner,
					StartedAt: s.StartedAt,
					EndedAt:   s.FinishedAt,
					Duration:  int64(s.FinishedAt.Sub(s.StartedAt).Seconds()),
				}

				// Save to database if enabled, otherwise use file store
				if database.enabled {
					database.SaveGame(rec)
					database.IncrementWinner(winner)
				} else {
					store.AppendGame(rec)
					store.IncrementWinner(winner)
				}

				// emit event to Kafka and file
				emitEvent(map[string]interface{}{
					"type": "game_finished",
					"game": rec,
				})
				broadcastState(s)
				return
			}
			// check draw: board is full AND no winning condition exists
			if boardFull(s.Game.Board) && !checkAnyWin(s.Game.Board) {
				fmt.Println("DRAW DETECTED!")
				s.State = "finished"
				s.Result = "draw"
				fmt.Printf("Set s.Result to: %s\n", s.Result)
				s.FinishedAt = time.Now()
				rec := GameRecord{
					ID:        s.ID,
					Player1:   s.Player1,
					Player2:   s.Player2,
					Winner:    "draw",
					StartedAt: s.StartedAt,
					EndedAt:   s.FinishedAt,
					Duration:  int64(s.FinishedAt.Sub(s.StartedAt).Seconds()),
				}

				// Save to database if enabled, otherwise use file store
				if database.enabled {
					database.SaveGame(rec)
				} else {
					store.AppendGame(rec)
				}

				// emit event to Kafka and file
				emitEvent(map[string]interface{}{
					"type": "game_finished",
					"game": rec,
				})
				broadcastState(s)
				return
			}
			// switch turn
			s.Game.Turn = 3 - s.Game.Turn
			// broadcast
			emitEvent(map[string]interface{}{
				"type":   "move",
				"gameId": s.ID,
				"col":    col,
				"player": s.Game.Board[r][col],
			})
			broadcastState(s)
			return
		}
	}
}

func broadcastState(s *GameSession) {
	fmt.Printf("Broadcasting state: status=%s, result=%s\n", s.State, s.Result)
	for uname, cl := range s.clients {
		msg := map[string]interface{}{"type": "state", "gameId": s.ID, "state": s.Game, "you": s.Players[uname], "status": s.State, "result": s.Result}
		fmt.Printf("Sending to %s: %+v\n", uname, msg)
		_ = cl.SendJSON(msg)
	}
}

func boardFull(b [][]int) bool {
	for r := 0; r < len(b); r++ {
		for c := 0; c < len(b[0]); c++ {
			if b[r][c] == 0 {
				return false
			}
		}
	}
	return true
}

// checkAnyWin checks if there's any winning condition on the board for any player
func checkAnyWin(b [][]int) bool {
	rows := len(b)
	cols := len(b[0])
	// Check for any player (1 or 2) if they have 4 in a row
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			if b[r][c] != 0 {
				if checkWin(b, r, c, b[r][c]) {
					return true
				}
			}
		}
	}
	return false
}

// checkWin checks if placing at r,c by player p created a win
func checkWin(b [][]int, r, c, p int) bool {
	dirs := [][2]int{{0, 1}, {1, 0}, {1, 1}, {1, -1}}
	rows := len(b)
	cols := len(b[0])
	for _, d := range dirs {
		cnt := 1
		// Check in positive direction
		for step := 1; step <= 6; step++ {
			r2 := r + d[0]*step
			c2 := c + d[1]*step
			if r2 < 0 || r2 >= rows || c2 < 0 || c2 >= cols || b[r2][c2] != p {
				break
			}
			cnt++
		}
		// Check in negative direction
		for step := 1; step <= 6; step++ {
			r2 := r - d[0]*step
			c2 := c - d[1]*step
			if r2 < 0 || r2 >= rows || c2 < 0 || c2 >= cols || b[r2][c2] != p {
				break
			}
			cnt++
		}
		if cnt >= 4 {
			return true
		}
	}
	return false
}

// small helper to copy board
func cloneBoard(b [][]int) [][]int {
	rows := len(b)
	cols := len(b[0])
	n := make([][]int, rows)
	for r := 0; r < rows; r++ {
		n[r] = make([]int, cols)
		copy(n[r], b[r])
	}
	return n
}

// emitEvent writes to disk and sends to Kafka if enabled
func emitEvent(e map[string]interface{}) {
	b, _ := json.Marshal(e)

	// append to events file
	f, _ := openAppend(config.DataDir + "/events.jsonl")
	if f != nil {
		defer f.Close()
		f.Write(append(b, '\n'))
	}

	// Send to Kafka if enabled
	if kafkaProducer != nil {
		kafkaProducer.SendEvent(e)
	}
}

func openAppend(path string) (*os.File, error) {
	return os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
}

// small random helper used by bot
func randomInt(n int) int { return rand.Intn(n) }
