# Setup Guide - 4 in a Row Game Server

This guide will help you set up and run the Connect-4 game server with all its components.

## Prerequisites

### Required
- **Go 1.20+** - [Download](https://golang.org/dl/)
- **Web Browser** - Chrome, Firefox, Safari, or Edge

### Optional (for full features)
- **PostgreSQL 12+** - [Download](https://www.postgresql.org/download/)
- **Apache Kafka** - [Download](https://kafka.apache.org/downloads)

## Installation Steps

### 1. Install Go

**Windows:**
1. Download Go installer from https://golang.org/dl/
2. Run the installer
3. Verify installation: Open PowerShell and run `go version`

**macOS:**
```bash
brew install go
go version
```

**Linux:**
```bash
wget https://go.dev/dl/go1.20.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.20.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin
go version
```

### 2. Clone/Download the Project

```bash
cd c:\Users\Aditi\Desktop\Connect-4
```

### 3. Install Go Dependencies

```bash
go mod download
```

This will download:
- `github.com/gorilla/websocket` - WebSocket support
- `github.com/segmentio/kafka-go` - Kafka client
- `github.com/lib/pq` - PostgreSQL driver

## Running the Server

### Option 1: Basic Mode (File Storage Only)

This is the simplest way to run the server. No database or Kafka required.

**Windows PowerShell:**
```powershell
go run ./server
```

**macOS/Linux:**
```bash
go run ./server
```

The server will:
- Start on http://localhost:8080
- Store data in `data/` directory (auto-created)
- Use JSON files for persistence

### Option 2: With PostgreSQL

#### Setup PostgreSQL

1. **Install PostgreSQL** (if not already installed)

2. **Create Database:**
```sql
CREATE DATABASE connect4;
```

3. **Set Environment Variables:**

**Windows PowerShell:**
```powershell
$env:DB_ENABLED="true"
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:DB_USER="postgres"
$env:DB_PASSWORD="your_password"
$env:DB_NAME="connect4"
go run ./server
```

**macOS/Linux:**
```bash
export DB_ENABLED=true
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=connect4
go run ./server
```

The server will automatically create the required tables on first run.

### Option 3: With Kafka Analytics

#### Setup Kafka

1. **Download Kafka** from https://kafka.apache.org/downloads

2. **Start Zookeeper:**
```bash
# Windows
bin\windows\zookeeper-server-start.bat config\zookeeper.properties

# macOS/Linux
bin/zookeeper-server-start.sh config/zookeeper.properties
```

3. **Start Kafka Broker:**
```bash
# Windows
bin\windows\kafka-server-start.bat config\server.properties

# macOS/Linux
bin/kafka-server-start.sh config/server.properties
```

4. **Create Topic:**
```bash
# Windows
bin\windows\kafka-topics.bat --create --topic game-analytics --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1

# macOS/Linux
bin/kafka-topics.sh --create --topic game-analytics --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1
```

5. **Run Server with Kafka:**

**Windows PowerShell:**
```powershell
$env:KAFKA_ENABLED="true"
$env:KAFKA_BROKERS="localhost:9092"
$env:KAFKA_TOPIC="game-analytics"
go run ./server
```

**macOS/Linux:**
```bash
export KAFKA_ENABLED=true
export KAFKA_BROKERS=localhost:9092
export KAFKA_TOPIC=game-analytics
go run ./server
```

6. **Run Analytics Consumer:**

Open a new terminal:
```bash
go run ./cmd/analytics -brokers localhost:9092 -topic game-analytics
```

### Option 4: Full Production Mode

Combine PostgreSQL and Kafka:

**Windows PowerShell:**
```powershell
$env:SERVER_PORT="8080"
$env:DB_ENABLED="true"
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:DB_USER="postgres"
$env:DB_PASSWORD="your_password"
$env:DB_NAME="connect4"
$env:KAFKA_ENABLED="true"
$env:KAFKA_BROKERS="localhost:9092"
$env:KAFKA_TOPIC="game-analytics"
$env:MATCH_TIMEOUT="10"
$env:RECONNECT_TIMEOUT="30"
go run ./server
```

**macOS/Linux:**
```bash
export SERVER_PORT=8080
export DB_ENABLED=true
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_password
export DB_NAME=connect4
export KAFKA_ENABLED=true
export KAFKA_BROKERS=localhost:9092
export KAFKA_TOPIC=game-analytics
export MATCH_TIMEOUT=10
export RECONNECT_TIMEOUT=30
go run ./server
```

## Playing the Game

1. **Open Browser:** Navigate to http://localhost:8080

2. **Enter Username:** Type your username in the input field

3. **Join Game:** Click "Join Matchmaking"

4. **Wait for Opponent:**
   - If another player joins within 10 seconds, you'll play against them
   - Otherwise, you'll play against the bot

5. **Play:** Click on any column to drop your disc

## Testing Different Scenarios

### Test 1: Player vs Bot
1. Open one browser window
2. Enter username and join
3. Wait 10 seconds
4. Bot will start playing

### Test 2: Player vs Player
1. Open two browser windows (or use incognito mode)
2. In first window: Enter "Alice" and join
3. In second window: Enter "Bob" and join (within 10 seconds)
4. Play against each other

### Test 3: Reconnection
1. Start a game
2. Close the browser tab (or refresh)
3. Within 30 seconds, reopen and enter the same username
4. You should reconnect to the same game

### Test 4: Leaderboard
1. Complete several games
2. Check the leaderboard at the bottom of the page
3. Winners should be ranked by number of wins

## Monitoring

### View Events (File Mode)
```bash
# View all events
cat data/events.jsonl

# View games
cat data/games.json

# View leaderboard
cat data/leaderboard.json
```

### View Analytics (Kafka Mode)
The analytics consumer will print real-time statistics including:
- Total games and moves
- Average game duration
- Top winners
- Games per day/hour
- Per-user statistics

### Database Queries (PostgreSQL Mode)
```sql
-- Connect to database
psql -U postgres -d connect4

-- View recent games
SELECT * FROM games ORDER BY ended_at DESC LIMIT 10;

-- View leaderboard
SELECT * FROM leaderboard ORDER BY wins DESC;

-- Average game duration
SELECT AVG(duration_seconds) as avg_duration FROM games;
```

## Troubleshooting

### "go: command not found"
- Go is not installed or not in PATH
- Install Go and restart terminal

### "port 8080 already in use"
- Another application is using port 8080
- Change port: `export SERVER_PORT=8081`
- Or stop the other application

### "database connection failed"
- PostgreSQL is not running
- Check credentials
- Server will fallback to file storage automatically

### "kafka connection failed"
- Kafka is not running
- Check broker address
- Server will continue without Kafka

### WebSocket connection fails
- Check browser console (F12)
- Verify server is running
- Try different browser

## Building for Production

### Build Binary
```bash
# Build for current platform
go build -o connect4-server ./server

# Run the binary
./connect4-server
```

### Build for Different Platforms
```bash
# Windows
GOOS=windows GOARCH=amd64 go build -o connect4-server.exe ./server

# Linux
GOOS=linux GOARCH=amd64 go build -o connect4-server-linux ./server

# macOS
GOOS=darwin GOARCH=amd64 go build -o connect4-server-mac ./server
```

## Next Steps

1. **Customize Configuration:** Edit environment variables to match your setup
2. **Add SSL/TLS:** Use nginx or similar reverse proxy
3. **Scale:** Deploy multiple instances behind a load balancer
4. **Monitor:** Set up logging and metrics collection
5. **Backup:** Regular database backups

## Support

For issues or questions:
1. Check the main README.md
2. Review the troubleshooting section
3. Check server logs for error messages

---

**Happy Gaming! 🎮**

