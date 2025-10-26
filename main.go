package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}

// Global configuration and services
var (
	config        *Config
	kafkaProducer *KafkaProducer
	database      *Database
)

// in-memory state
var (
	gamesMu sync.Mutex
	games   = map[string]*GameSession{} // gameId -> session
	waitMu  sync.Mutex
	waiting = []string{}           // usernames waiting
	clients = map[string]*Client{} // username -> client
	store   *FileStore
	roomsMu sync.Mutex
	rooms   = map[string]*Room{} // roomId -> room
)

func main() {
	// Load configuration
	config = LoadConfig()
	log.Printf("Configuration loaded: Port=%s, KafkaEnabled=%v, DBEnabled=%v",
		config.ServerPort, config.KafkaEnabled, config.DBEnabled)

	// Initialize services
	os.MkdirAll("static", 0755)
	os.MkdirAll(config.DataDir, 0755)

	// Initialize Kafka producer
	kafkaProducer = NewKafkaProducer(config)
	defer kafkaProducer.Close()

	// Initialize database
	database = NewDatabase(config)
	defer database.Close()

	// Initialize file store (fallback or primary storage)
	store = NewFileStore(
		config.DataDir+"/games.json",
		config.DataDir+"/leaderboard.json",
	)

	// Setup HTTP handlers
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/leaderboard", leaderboardHandler)
	http.HandleFunc("/rooms", roomsHandler)

	addr := ":" + config.ServerPort
	log.Printf("Server starting on %s", addr)
	go reaper()
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}

func leaderboardHandler(w http.ResponseWriter, r *http.Request) {
	var lb Leaderboard

	// Try database first, fallback to file store
	if database.enabled {
		var err error
		lb, err = database.GetLeaderboard()
		if err != nil {
			lb = store.LoadLeaderboard()
		}
	} else {
		lb = store.LoadLeaderboard()
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(lb)
}

func roomsHandler(w http.ResponseWriter, r *http.Request) {
	roomsMu.Lock()
	defer roomsMu.Unlock()

	// Build list of available rooms
	roomList := []RoomInfo{}
	for _, room := range rooms {
		// Only show rooms that are waiting for players
		if room.Status == "waiting" {
			playerCount := 0
			if room.Player1 != "" {
				playerCount++
			}
			if room.Player2 != "" {
				playerCount++
			}
			roomList = append(roomList, RoomInfo{
				ID:         room.ID,
				Name:       room.Name,
				Creator:    room.Creator,
				Players:    playerCount,
				MaxPlayers: 2,
				Status:     room.Status,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roomList)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer c.Close()

	// expect first message to be join, create_room, or join_room
	var msg map[string]interface{}
	c.SetReadDeadline(time.Now().Add(30 * time.Second))
	if err := c.ReadJSON(&msg); err != nil {
		log.Println("read join err:", err)
		return
	}

	type JoinMsg struct {
		Type     string `json:"type"`
		Username string `json:"username"`
		GameID   string `json:"gameId,omitempty"`
		RoomID   string `json:"roomId,omitempty"`
		RoomName string `json:"roomName,omitempty"`
	}

	msgType, _ := msg["type"].(string)
	if msgType != "join" && msgType != "create_room" && msgType != "join_room" {
		c.WriteJSON(map[string]string{"error": "first message must be join, create_room, or join_room"})
		return
	}

	join := JoinMsg{}
	b, _ := json.Marshal(msg)
	json.Unmarshal(b, &join)
	username := join.Username

	client := &Client{Username: username, Conn: c}
	clients[username] = client
	defer func() { delete(clients, username) }()

	// Handle different message types
	switch msgType {
	case "create_room":
		// Create a new room
		room := createRoom(username, join.RoomName)
		c.WriteJSON(map[string]interface{}{
			"type":   "room_created",
			"roomId": room.ID,
			"room":   room,
		})
		log.Printf("Player %s created room %s (%s)", username, room.Name, room.ID)
		client.readPump(nil)
		return

	case "join_room":
		// Join an existing room
		roomsMu.Lock()
		room, ok := rooms[join.RoomID]
		roomsMu.Unlock()

		if !ok {
			c.WriteJSON(map[string]string{"error": "room not found"})
			return
		}

		if room.Status != "waiting" {
			c.WriteJSON(map[string]string{"error": "room is not available"})
			return
		}

		// Add player to room
		roomsMu.Lock()
		if room.Player1 == "" {
			room.Player1 = username
		} else if room.Player2 == "" {
			room.Player2 = username
		} else {
			roomsMu.Unlock()
			c.WriteJSON(map[string]string{"error": "room is full"})
			return
		}

		// If room now has 2 players, start the game
		if room.Player1 != "" && room.Player2 != "" {
			room.Status = "playing"
			p1 := room.Player1
			p2 := room.Player2
			roomsMu.Unlock()

			log.Printf("Room %s is full, starting game: %s vs %s", room.ID, p1, p2)
			go startGameFromRoom(room.ID, p1, p2)
		} else {
			roomsMu.Unlock()
			c.WriteJSON(map[string]interface{}{
				"type":   "room_joined",
				"roomId": room.ID,
				"room":   room,
			})
		}

		client.readPump(nil)
		return

	case "join":
		// if GameID provided, try to reconnect
		if join.GameID != "" {
			gamesMu.Lock()
			sess, ok := games[join.GameID]
			gamesMu.Unlock()
			if ok {
				sess.reconnect(username, client)
				// keep reading messages
				client.readPump(sess)
				return
			}
		}

		// otherwise join matchmaking
		enqueueWaiting(username)
		// notify client that they're waiting (always 15 seconds)
		c.WriteJSON(map[string]interface{}{"type": "waiting", "timeout": 15})

		// keep reading messages until connection closed
		client.readPump(nil)
	}
}

// reaper cleans up timed-out games and old rooms
func reaper() {
	for range time.NewTicker(5 * time.Second).C {
		gamesMu.Lock()
		for id, g := range games {
			if g.State == "finished" {
				// remove after some time
				if time.Since(g.FinishedAt) > 10*time.Minute {
					delete(games, id)
				}
			}
		}
		gamesMu.Unlock()

		// Clean up old rooms
		roomsMu.Lock()
		for id, room := range rooms {
			// Remove rooms that have been waiting for more than 10 minutes
			if room.Status == "waiting" && time.Since(room.CreatedAt) > 10*time.Minute {
				log.Printf("Removing stale room %s (%s)", room.Name, room.ID)
				delete(rooms, id)
			}
			// Remove finished rooms after 5 minutes
			if room.Status == "finished" {
				// Check if associated game is also cleaned up
				gamesMu.Lock()
				if _, exists := games[room.GameID]; !exists {
					delete(rooms, id)
				}
				gamesMu.Unlock()
			}
		}
		roomsMu.Unlock()
	}
}

func enqueueWaiting(username string) {
	waitMu.Lock()
	waiting = append(waiting, username)
	waitMu.Unlock()

	log.Printf("Player %s joined matchmaking queue, waiting 15 seconds...", username)

	// Always wait 15 seconds before starting a game
	// This ensures players have enough time to find human opponents
	go func() {
		// Wait for exactly 15 seconds
		time.Sleep(15 * time.Second)

		// After 15 seconds, check if another player is available
		waitMu.Lock()
		defer waitMu.Unlock()

		// Find this player in the waiting queue
		playerIndex := -1
		for i, name := range waiting {
			if name == username {
				playerIndex = i
				break
			}
		}

		// If player is no longer in queue (already matched), return
		if playerIndex == -1 {
			log.Printf("Player %s already matched, skipping", username)
			return
		}

		// Check if there's another player waiting
		if len(waiting) >= 2 {
			// Match with another player
			p1 := waiting[0]
			p2 := waiting[1]
			waiting = waiting[2:]
			log.Printf("Matching %s with %s after 15 second wait", p1, p2)
			go startGame(p1, p2)
		} else {
			// No other player available, start game with bot
			waiting = waiting[1:]
			log.Printf("No opponent found for %s after 15 seconds, starting bot game", username)
			go startGameWithBot(username)
		}
	}()
}

func startGame(p1, p2 string) {
	log.Printf("Starting game: %s vs %s", p1, p2)
	g := NewGameSession(p1, p2)
	gamesMu.Lock()
	games[g.ID] = g
	gamesMu.Unlock()
	// register connected clients in the session and send initial state
	if c1, ok := clients[p1]; ok {
		g.clients[p1] = c1
		c1.SendJSON(map[string]interface{}{"type": "start", "gameId": g.ID, "you": 1, "opponent": p2, "state": g.Game})
	}
	if c2, ok := clients[p2]; ok {
		g.clients[p2] = c2
		c2.SendJSON(map[string]interface{}{"type": "start", "gameId": g.ID, "you": 2, "opponent": p1, "state": g.Game})
	}
	// start goroutine to process game moves
	go g.run()
}

func startGameWithBot(player string) {
	botName := "Bot"
	log.Printf("Starting game: %s vs BOT", player)
	g := NewGameSession(player, botName)
	g.IsBot = true
	gamesMu.Lock()
	games[g.ID] = g
	gamesMu.Unlock()
	// register player client and send initial state
	if c1, ok := clients[player]; ok {
		g.clients[player] = c1
		c1.SendJSON(map[string]interface{}{"type": "start", "gameId": g.ID, "you": 1, "opponent": botName, "state": g.Game})
	}
	go g.run()
}

// createRoom creates a new game room
func createRoom(creator, roomName string) *Room {
	roomsMu.Lock()
	defer roomsMu.Unlock()

	if roomName == "" {
		roomName = creator + "'s room"
	}

	room := &Room{
		ID:        fmt.Sprintf("r_%d", time.Now().UnixNano()),
		Name:      roomName,
		Creator:   creator,
		Player1:   creator,
		Status:    "waiting",
		CreatedAt: time.Now(),
	}

	rooms[room.ID] = room
	return room
}

// startGameFromRoom starts a game from a room
func startGameFromRoom(roomID, p1, p2 string) {
	log.Printf("Starting game from room %s: %s vs %s", roomID, p1, p2)
	g := NewGameSession(p1, p2)

	gamesMu.Lock()
	games[g.ID] = g
	gamesMu.Unlock()

	// Update room with game ID
	roomsMu.Lock()
	if room, ok := rooms[roomID]; ok {
		room.GameID = g.ID
	}
	roomsMu.Unlock()

	// register connected clients in the session and send initial state
	if c1, ok := clients[p1]; ok {
		g.clients[p1] = c1
		c1.SendJSON(map[string]interface{}{"type": "start", "gameId": g.ID, "you": 1, "opponent": p2, "state": g.Game})
	}
	if c2, ok := clients[p2]; ok {
		g.clients[p2] = c2
		c2.SendJSON(map[string]interface{}{"type": "start", "gameId": g.ID, "you": 2, "opponent": p1, "state": g.Game})
	}
	// start goroutine to process game moves
	go g.run()
}

// After scaffolding we will add additional files implementing GameSession, Client, store, bot, etc.
