package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// Client wraps a websocket connection and username
type Client struct {
	Username string
	Conn     *websocket.Conn
	sendMu   chan struct{}
}

func (c *Client) SendJSON(v any) error {
	c.Conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	return c.Conn.WriteJSON(v)
}

// readPump listens for incoming messages from a client and routes them
func (c *Client) readPump(sess *GameSession) {
	defer c.Conn.Close()

	// Set up ping/pong to keep connection alive
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	// Start ping ticker to keep connection alive
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	go func() {
		for range ticker.C {
			c.Conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}()

	for {
		var m map[string]any
		if err := c.Conn.ReadJSON(&m); err != nil {
			log.Printf("readPump read error for user %s: %v", c.Username, err)
			// if client disconnected, allow reconnection timeout
			if sess != nil && sess.State == "playing" {
				// Mark client as disconnected
				delete(sess.clients, c.Username)

				// Start reconnection timer
				go func(u string, s *GameSession) {
					reconnectTimeout := time.Duration(config.ReconnectTimeout) * time.Second
					time.Sleep(reconnectTimeout)

					// Check if still not reconnected and game still playing
					s.TurnMu.Lock()
					defer s.TurnMu.Unlock()

					if s.State != "playing" {
						return // Game already finished
					}

					if _, ok := s.clients[u]; !ok {
						// Player didn't reconnect, declare forfeit
						log.Printf("Player %s forfeited due to disconnect", u)

						// Determine winner (the other player)
						winner := s.Player1
						if u == s.Player1 {
							winner = s.Player2
						}

						s.State = "finished"
						s.Result = winner
						s.FinishedAt = time.Now()

						// Save game result
						rec := GameRecord{
							ID:        s.ID,
							Player1:   s.Player1,
							Player2:   s.Player2,
							Winner:    winner,
							StartedAt: s.StartedAt,
							EndedAt:   s.FinishedAt,
							Duration:  int64(s.FinishedAt.Sub(s.StartedAt).Seconds()),
						}

						if database.enabled {
							database.SaveGame(rec)
							database.IncrementWinner(winner)
						} else {
							store.AppendGame(rec)
							store.IncrementWinner(winner)
						}

						emitEvent(map[string]any{
							"type":   "game_finished",
							"reason": "forfeit",
							"game":   rec,
						})

						broadcastState(s)
					}
				}(c.Username, sess)
			}
			return
		}
		// handle types
		typ, _ := m["type"].(string)
		switch typ {
		case "move":
			if sess == nil {
				// find session by id
				if gid, ok := m["gameId"].(string); ok {
					gamesMu.Lock()
					sess = games[gid]
					gamesMu.Unlock()
				}
			}
			if sess == nil {
				continue
			}
			colf, _ := m["col"].(float64)
			col := int(colf)
			sess.applyMove(col)
		case "join":
			// ignored here
		}
	}
}
