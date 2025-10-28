# 4 in a Row — Real-Time Connect Four Game Server

A production-ready, real-time multiplayer Connect-4 game server built with Go, featuring WebSocket communication, competitive AI bot, Kafka analytics, and PostgreSQL persistence.

---

##  Table of Contents
- [Setup Instructions](#-setup-instructions)
- [Game Features](#-game-features)
- [How to Run](#-how-to-run)
- [Configuration](#️-configuration)
- [Project Structure](#-project-structure)
- [API Documentation](#-api-documentation)

---

##  Setup Instructions

### Prerequisites

Before setting up the project, ensure you have the following installed:

- **Go 1.20 or higher** - [Download Go](https://golang.org/dl/)
- **Apache Kafka** - For real-time analytics streaming
- **(Optional) PostgreSQL 12+** - For database persistence
- **(Optional) Docker & Docker Compose** - For running services in containers

### Step 1: Clone the Repository

```bash
git clone <repository-url>
cd Connect-4
```

### Step 2: Install Go Dependencies

```bash
go mod download
```

This will download all required Go packages including:
- `github.com/gorilla/websocket` - WebSocket support
- `github.com/segmentio/kafka-go` - Kafka client
- `github.com/lib/pq` - PostgreSQL driver

### Step 3: Environment Configuration

The project uses environment variables for configuration. You can either:

**Option A: Use the provided `.env.example` file**

```bash
# Copy the example file
cp .env.example .env

# Edit .env with your preferred settings
```

**Option B: Set environment variables manually**

See the [Configuration](#️-configuration) section below for all available variables.

### Step 4: Optional Services Setup

#### PostgreSQL Setup (Optional)

If you want to use PostgreSQL for persistence:

**Using Docker:**
```bash
docker-compose up -d postgres
```

**Manual Installation:**
1. Install PostgreSQL 12+
2. Create a database:
   ```sql
   CREATE DATABASE connect4;
   ```
3. Set environment variables:
   ```bash
   export DB_ENABLED=true
   export DB_HOST=localhost
   export DB_PORT=5432
   export DB_USER=postgres
   export DB_PASSWORD=yourpassword
   export DB_NAME=connect4
   ```

The required tables will be created automatically on first run.

#### Kafka Setup (Optional)

If you want to enable real-time analytics:

**Using Docker:**
```bash
docker-compose up -d zookeeper kafka
```

**Manual Installation:**
1. Install Apache Kafka
2. Start Zookeeper and Kafka broker
3. Set environment variables:
   ```bash
   export KAFKA_ENABLED=true
   export KAFKA_BROKERS=localhost:9092
   export KAFKA_TOPIC=game-analytics
   ```

#### Using Docker Compose for All Services

To run PostgreSQL, Kafka, and Zookeeper together:

```bash
docker-compose up -d
```

This will start:
- PostgreSQL on port 5432
- Kafka on port 9092
- Zookeeper on port 2181
- Kafka UI on port 8090 (for monitoring)

---

## Game Features

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

### Game Rules
- **Board**: 7 columns × 6 rows
- **Players**: Two players take turns (Red vs Yellow)
- **Objective**: Connect 4 discs vertically, horizontally, or diagonally
- **Gameplay**: Click any column to drop your disc
- **Win**: First to connect 4 wins
- **Draw**: Board fills with no winner

### Bot Strategy

The competitive bot uses a strategic algorithm:

1. **Win immediately** - Takes winning move if available
2. **Block opponent** - Blocks opponent's winning move
3. **Strategic positioning** - Prefers center columns for better opportunities
4. **Fallback** - Makes valid moves when no immediate threats/opportunities

---

## How to Run

### Method 1: Quick Start (Basic Mode)

Run the server with file-based storage (no database or Kafka required):

```bash
go run ./server
```

The server will start on `http://localhost:8080`

### Method 2: Using PowerShell Script (Windows)

For Windows users, use the provided startup script:

```powershell
.\start.ps1
```

This script will:
- Check if Go is installed
- Download dependencies if needed
- Set default environment variables
- Create the data directory
- Start the server

### Method 3: With PostgreSQL

Run with database persistence:

```bash
# Set environment variables
export DB_ENABLED=true
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=yourpassword
export DB_NAME=connect4

# Start the server
go run ./server
```

### Method 4: With Kafka Analytics

Run with real-time analytics:

```bash
# Set environment variables
export KAFKA_ENABLED=true
export KAFKA_BROKERS=localhost:9092
export KAFKA_TOPIC=game-analytics

# Start the server
go run ./server
```

### Method 5: Full Production Mode

Run with all features enabled:

```bash
# Set all environment variables
export SERVER_PORT=8080
export DB_ENABLED=true
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=yourpassword
export DB_NAME=connect4
export KAFKA_ENABLED=true
export KAFKA_BROKERS=localhost:9092
export KAFKA_TOPIC=game-analytics
export MATCH_TIMEOUT=10
export RECONNECT_TIMEOUT=30

# Start the server
go run ./server
```

### Method 6: Using Docker

Build and run with Docker:

```bash
# Build the image
docker build -t connect4-server .

# Run the container
docker run -p 8080:8080 connect4-server
```

### Playing the Game

Once the server is running:

1. **Open your browser** and navigate to `http://localhost:8080`
2. **Enter a username** in the input field
3. **Click "Join Matchmaking"** to start searching for an opponent
4. **Wait for matchmaking**:
   - If another player joins within 10 seconds, you'll play against them
   - If no one joins, you'll automatically play against the AI bot
5. **Play the game**:
   - Click on any column to drop your disc
   - Take turns with your opponent
   - First to connect 4 wins!
6. **Check the leaderboard** to see top players

### Running the Analytics Consumer (Optional)

If you have Kafka enabled, you can run the analytics consumer to see real-time game statistics:

```bash
go run ./cmd/analytics -brokers localhost:9092 -topic game-analytics -group analytics-consumer
```

The consumer will display:
- Total games and moves
- Average game duration
- Top winners
- Games per day/hour
- Per-user statistics (win rate, games played, etc.)

### Testing Multiplayer

To test multiplayer functionality:

1. **Open two browser windows** (or use incognito mode)
2. **Enter different usernames** in each window
3. **Click "Join Matchmaking"** in both windows within 10 seconds
4. **Play against each other** in real-time

### Testing Reconnection

To test the reconnection feature:

1. **Start a game** with a username
2. **Close the browser tab** during the game
3. **Reopen the browser** within 30 seconds
4. **Enter the same username** and click "Join Matchmaking"
5. **You should reconnect** to your ongoing game

---

## Configuration

All configuration is done via environment variables. You can set these in your shell, in a `.env` file, or in the `start.ps1` script.

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `8080` | HTTP server port |
| `DATA_DIR` | `data` | Directory for file storage |
| `MATCH_TIMEOUT` | `10` | Seconds to wait for matchmaking |
| `RECONNECT_TIMEOUT` | `30` | Seconds to allow reconnection |
| `KAFKA_ENABLED` | `false` | Enable Kafka producer |
| `KAFKA_BROKERS` | `localhost:9092` | Kafka broker addresses (comma-separated) |
| `KAFKA_TOPIC` | `game-analytics` | Kafka topic for events |
| `DB_ENABLED` | `false` | Enable PostgreSQL |
| `DB_HOST` | `localhost` | PostgreSQL host |
| `DB_PORT` | `5432` | PostgreSQL port |
| `DB_USER` | `postgres` | Database user |
| `DB_PASSWORD` | `postgres` | Database password |
| `DB_NAME` | `connect4` | Database name |

### Example Configuration Files

**`.env` file example:**
```bash
# Server Configuration
SERVER_PORT=8080
DATA_DIR=data

# Matchmaking & Game Settings
MATCH_TIMEOUT=10
RECONNECT_TIMEOUT=30

# Kafka Configuration (Optional)
KAFKA_ENABLED=true
KAFKA_BROKERS=localhost:9092
KAFKA_TOPIC=game-analytics

# Database Configuration (Optional)
DB_ENABLED=true
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=connect4
```

---

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
├── .env.example        # Example environment variables
├── start.ps1           # Windows startup script
├── docker-compose.yml  # Docker services configuration
├── go.mod              # Go dependencies
└── README.md           # This file
```

---

## � API Documentation

### HTTP Endpoints

- `GET /` - Serve the game frontend (HTML/CSS/JS)
- `GET /ws` - WebSocket endpoint for real-time game communication
- `GET /leaderboard` - Get current leaderboard (JSON format)

### Database Schema

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

### WebSocket API

#### Client → Server Messages

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

#### Server → Client Messages

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

### Kafka Analytics Events

The system emits the following events to Kafka when enabled:

**Game Finished Event:**
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

**Move Event:**
```json
{
  "type": "move",
  "timestamp": "2025-10-24T10:29:00Z",
  "gameId": "g_1729765800000000000",
  "col": 3,
  "player": 1
}
```

---

##  Monitoring & Analytics

### File-based Storage

When running without PostgreSQL, data is stored in JSON files:

- `data/events.jsonl` - All game events (append-only log)
- `data/games.json` - Completed games history
- `data/leaderboard.json` - Player rankings and win counts

### Database Queries

If using PostgreSQL, you can query game data:

```sql
-- Top 10 players
SELECT username, wins FROM leaderboard ORDER BY wins DESC LIMIT 10;

-- Recent games
SELECT * FROM games ORDER BY ended_at DESC LIMIT 20;

-- Average game duration
SELECT AVG(duration_seconds) FROM games;

-- Games per day
SELECT DATE(ended_at), COUNT(*) FROM games GROUP BY DATE(ended_at);

-- Win rate by player
SELECT username, wins,
       (wins::float / (SELECT COUNT(*) FROM games WHERE player1 = username OR player2 = username)) as win_rate
FROM leaderboard;
```

### Kafka Monitoring

If using Kafka, you can monitor events using:

1. **Kafka UI** (included in docker-compose):
   - Access at `http://localhost:8090`
   - View topics, messages, and consumer groups

2. **Analytics Consumer**:
   - Run `go run ./cmd/analytics` to see real-time statistics
   - Displays aggregated metrics and per-user stats

---


