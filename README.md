# 4 in a Row — Real-Time Connect Four Game Server

A production-ready, real-time multiplayer Connect-4 game server built with Go, featuring WebSocket communication, competitive AI bot, Kafka analytics, and PostgreSQL persistence.

## 🎮 Features

### Core Gameplay
- **Real-time multiplayer** via WebSockets
- **Smart matchmaking** with 10-second timeout (configurable)
- **Competitive AI bot** with strategic decision-making (blocks wins, creates opportunities)
- **Reconnection support** - Players can rejoin within 30 seconds (configurable)
- **Automatic forfeit** if player doesn't reconnect in time

### Backend Architecture
- **In-memory game state** for active games
- **Dual persistence** - PostgreSQL (primary) with JSON file fallback
- **Kafka integration** for real-time analytics events
- **Configurable via environment variables**
- **Separation of concerns** - Clean modular architecture

### Analytics
- **Kafka-based event streaming** for game analytics
- **Comprehensive metrics tracking**:
  - Total games, moves, duration statistics
  - Per-user stats (games played, win rate, average duration)
  - Games per hour/day tracking
  - Bot vs Player game statistics
  - Draw and forfeit tracking

### Frontend
- **Modern responsive UI** with gradient design
- **Real-time game updates** via WebSocket
- **Turn indicators** and player highlighting
- **Winner announcements** with animations
- **Live leaderboard** with auto-refresh
- **Visual feedback** for all game states

## 📁 Project Structure

```
Connect-4/
├── server/              # Main game server
│   ├── main.go         # HTTP server & WebSocket handler
│   ├── game.go         # Game logic & session management
│   ├── bot.go          # AI bot implementation
│   ├── ws.go           # WebSocket client handling
│   ├── models.go       # Data models
│   ├── store.go        # File-based storage
│   ├── database.go     # PostgreSQL persistence
│   ├── kafka.go        # Kafka producer
│   └── config.go       # Configuration management
├── cmd/analytics/      # Analytics consumer
│   └── consumer.go     # Kafka consumer with metrics
├── static/             # Frontend files
│   ├── index.html      # Game UI
│   └── app.js          # WebSocket client & game logic
├── data/               # Data directory (auto-created)
│   ├── games.json      # Completed games (fallback)
│   ├── leaderboard.json # Player wins (fallback)
│   └── events.jsonl    # Event log
├── go.mod              # Go dependencies
└── README.md           # This file
```

## 🚀 Quick Start

### Prerequisites
- Go 1.20 or higher
- (Optional) PostgreSQL 12+ for database persistence
- (Optional) Apache Kafka for analytics streaming

### 1. Install Dependencies

```bash
go mod download
```

### 2. Run the Game Server

**Basic mode (file-based storage):**
```bash
go run ./server
```

**With PostgreSQL:**
```bash
export DB_ENABLED=true
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=yourpassword
export DB_NAME=connect4
go run ./server
```

**With Kafka:**
```bash
export KAFKA_ENABLED=true
export KAFKA_BROKERS=localhost:9092
export KAFKA_TOPIC=game-analytics
go run ./server
```

**Full production mode:**
```bash
export DB_ENABLED=true
export KAFKA_ENABLED=true
export SERVER_PORT=8080
export MATCH_TIMEOUT=10
export RECONNECT_TIMEOUT=30
go run ./server
```

Server will start on `http://localhost:8080`

### 3. Play the Game

1. Open `http://localhost:8080` in your browser
2. Enter a username and click "Join Matchmaking"
3. Open another browser window/tab with a different username
4. If no opponent joins within 10 seconds, you'll play against the bot

### 4. Run Analytics Consumer (Optional)

Requires Kafka to be running:

```bash
go run ./cmd/analytics -brokers localhost:9092 -topic game-analytics -group analytics-consumer
```

The consumer will display real-time analytics including:
- Total games and moves
- Average game duration
- Top winners
- Games per day/hour
- Per-user statistics (win rate, games played, etc.)

## ⚙️ Configuration

All configuration is done via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `8080` | HTTP server port |
| `DATA_DIR` | `data` | Directory for file storage |
| `KAFKA_ENABLED` | `false` | Enable Kafka producer |
| `KAFKA_BROKERS` | `localhost:9092` | Kafka broker addresses (comma-separated) |
| `KAFKA_TOPIC` | `game-analytics` | Kafka topic for events |
| `DB_ENABLED` | `false` | Enable PostgreSQL |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | `postgres` | Database password |
| `DB_NAME` | `connect4` | Database name |
| `MATCH_TIMEOUT` | `10` | Seconds to wait for matchmaking |
| `RECONNECT_TIMEOUT` | `30` | Seconds to allow reconnection |

## 🗄️ Database Setup

If using PostgreSQL, the tables are created automatically on first run:

```sql
-- Games table
CREATE TABLE games (
    id VARCHAR(255) PRIMARY KEY,
    player1 VARCHAR(255) NOT NULL,
    player2 VARCHAR(255) NOT NULL,
    winner VARCHAR(255) NOT NULL,
    duration_seconds BIGINT NOT NULL,
    started_at TIMESTAMP NOT NULL,
    ended_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Leaderboard table
CREATE TABLE leaderboard (
    username VARCHAR(255) PRIMARY KEY,
    wins INT NOT NULL DEFAULT 0,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```


## 🎯 Game Rules

- **Board**: 7 columns × 6 rows
- **Players**: Two players take turns (Red vs Yellow)
- **Objective**: Connect 4 discs vertically, horizontally, or diagonally
- **Gameplay**: Click any column to drop your disc
- **Win**: First to connect 4 wins
- **Draw**: Board fills with no winner

## 🤖 Bot Strategy

The competitive bot uses a strategic algorithm:

1. **Win immediately** - Takes winning move if available
2. **Block opponent** - Blocks opponent's winning move
3. **Strategic positioning** - Prefers center columns for better opportunities
4. **Fallback** - Makes valid moves when no immediate threats/opportunities

## 📊 Analytics Events

The system emits the following events to Kafka:

### `game_finished`
```json
{
  "type": "game_finished",
  "timestamp": "2025-10-24T10:30:00Z",
  "game": {
    "id": "g_1729765800000000000",
    "player1": "alice",
    "player2": "bob",
    "winner": "alice",
    "duration_seconds": 120,
    "started_at": "2025-10-24T10:28:00Z",
    "ended_at": "2025-10-24T10:30:00Z"
  },
  "reason": "win" // or "draw" or "forfeit"
}
```

### `move`
```json
{
  "type": "move",
  "timestamp": "2025-10-24T10:29:00Z",
  "gameId": "g_1729765800000000000",
  "col": 3,
  "player": 1
}
```

## 🔌 WebSocket API

### Client → Server

**Join Game:**
```json
{
  "type": "join",
  "username": "player1",
  "gameId": "g_xxx" // optional, for reconnection
}
```

**Make Move:**
```json
{
  "type": "move",
  "gameId": "g_xxx",
  "col": 3
}
```

### Server → Client

**Waiting for Opponent:**
```json
{
  "type": "waiting",
  "timeout": 10
}
```

**Game Started:**
```json
{
  "type": "start",
  "gameId": "g_xxx",
  "you": 1,
  "opponent": "player2",
  "state": {
    "rows": 6,
    "cols": 7,
    "board": [[0,0,0,0,0,0,0], ...],
    "turn": 1
  }
}
```

**Game State Update:**
```json
{
  "type": "state",
  "gameId": "g_xxx",
  "state": { ... },
  "you": 1,
  "status": "playing" // or "finished"
}
```

**Reconnected:**
```json
{
  "type": "reconnected",
  "gameId": "g_xxx",
  "state": { ... }
}
```

## 🧪 Testing

### Manual Testing

1. **Single Player vs Bot:**
   - Open browser, enter username, wait 10 seconds
   - Bot should start playing

2. **Two Players:**
   - Open two browser windows
   - Both join within 10 seconds
   - Play against each other

3. **Reconnection:**
   - Start a game
   - Close one browser tab
   - Reopen within 30 seconds with same username
   - Should reconnect to same game

4. **Leaderboard:**
   - Complete several games
   - Check leaderboard updates
   - Verify win counts

### Load Testing

```bash
# Install dependencies
go get -u github.com/gorilla/websocket

# Create a simple load test script or use existing tools
```

## 📈 Monitoring

### File-based Logs
- `data/events.jsonl` - All game events
- `data/games.json` - Completed games
- `data/leaderboard.json` - Player rankings

### Database Queries

```sql
-- Top 10 players
SELECT username, wins FROM leaderboard ORDER BY wins DESC LIMIT 10;

-- Recent games
SELECT * FROM games ORDER BY ended_at DESC LIMIT 20;

-- Average game duration
SELECT AVG(duration_seconds) FROM games;

-- Games per day
SELECT DATE(ended_at), COUNT(*) FROM games GROUP BY DATE(ended_at);
```

## 🐛 Troubleshooting

### Server won't start
- Check if port 8080 is available
- Verify Go version: `go version`
- Run `go mod tidy` to fix dependencies

### Database connection fails
- Verify PostgreSQL is running
- Check credentials in environment variables
- Server will fallback to file storage automatically

### Kafka connection fails
- Verify Kafka is running: `kafka-topics.sh --list --bootstrap-server localhost:9092`
- Check broker address in environment variables
- Server will continue without Kafka (events written to file only)

### WebSocket connection fails
- Check browser console for errors
- Verify server is running
- Check firewall settings

## 🚀 Production Deployment

### Recommended Setup

1. **Use PostgreSQL** for persistence
2. **Enable Kafka** for analytics
3. **Set appropriate timeouts** based on your needs
4. **Use reverse proxy** (nginx) for SSL/TLS
5. **Monitor logs** and metrics
6. **Set up backups** for database

### Docker Deployment (Example)

```dockerfile
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o server ./server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/server .
COPY --from=builder /app/static ./static
EXPOSE 8080
CMD ["./server"]
```

### Environment Variables for Production

```bash
SERVER_PORT=8080
DB_ENABLED=true
DB_HOST=postgres-host
DB_PORT=5432
DB_USER=connect4_user
DB_PASSWORD=secure_password
DB_NAME=connect4_prod
KAFKA_ENABLED=true
KAFKA_BROKERS=kafka1:9092,kafka2:9092,kafka3:9092
KAFKA_TOPIC=game-analytics
MATCH_TIMEOUT=10
RECONNECT_TIMEOUT=30
```

## 📝 API Endpoints

- `GET /` - Serve frontend
- `GET /ws` - WebSocket endpoint for game
- `GET /leaderboard` - Get current leaderboard (JSON)

## 🤝 Contributing

This is an assignment project, but improvements are welcome:

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## 📄 License

This project is created as a backend engineering intern assignment.

## 👨‍💻 Author

Built with ❤️ for the Backend Engineering Intern Assignment

---

**Note**: This implementation follows all requirements from the assignment:
- ✅ Real-time multiplayer with WebSockets
- ✅ 10-second matchmaking with bot fallback
- ✅ Competitive bot (not random)
- ✅ 30-second reconnection support
- ✅ In-memory + persistent storage
- ✅ Leaderboard tracking
- ✅ Simple functional frontend
- ✅ Kafka analytics integration
- ✅ Separation of concerns
- ✅ PostgreSQL support
- ✅ Comprehensive metrics tracking


