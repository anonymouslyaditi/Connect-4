package main

import (
	"encoding/json"
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

func wsHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer c.Close()

	// expect first message to be join
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
	}

	if msg["type"] != "join" {
		c.WriteJSON(map[string]string{"error": "first message must be join"})
		return
	}
	join := JoinMsg{}
	b, _ := json.Marshal(msg)
	json.Unmarshal(b, &join)
	username := join.Username

	client := &Client{Username: username, Conn: c}
	clients[username] = client
	defer func() { delete(clients, username) }()

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
	// notify client that they're waiting
	c.WriteJSON(map[string]interface{}{"type": "waiting", "timeout": config.MatchTimeout})

	// keep reading messages until connection closed
	client.readPump(nil)
}

// reaper cleans up timed-out games (it also enforces reconnect timeouts)
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
	}
}

func enqueueWaiting(username string) {
	waitMu.Lock()
	waiting = append(waiting, username)
	waitMu.Unlock()
	// try to match
	go func() {
		start := time.Now()
		timeout := time.Duration(config.MatchTimeout) * time.Second
		for {
			// if another player available, match immediately
			waitMu.Lock()
			if len(waiting) >= 2 {
				p1 := waiting[0]
				p2 := waiting[1]
				waiting = waiting[2:]
				waitMu.Unlock()
				startGame(p1, p2)
				return
			}
			waitMu.Unlock()
			// if timeout elapsed, start a bot game with first player
			if time.Since(start) >= timeout {
				waitMu.Lock()
				if len(waiting) >= 1 && waiting[0] == username {
					waiting = waiting[1:]
					waitMu.Unlock()
					startGameWithBot(username)
					return
				}
				waitMu.Unlock()
			}
			time.Sleep(300 * time.Millisecond)
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

// After scaffolding we will add additional files implementing GameSession, Client, store, bot, etc.
