# Quick Start Guide - 4 in a Row

Get up and running in 5 minutes!

## Prerequisites

You need to install Go first. Choose one method:

### Method 1: Download Installer (Easiest)
1. Go to: https://go.dev/dl/
2. Download: `go1.21.x.windows-amd64.msi`
3. Run the installer
4. **Close and reopen PowerShell**

### Method 2: Using Chocolatey
```powershell
choco install golang
```

### Method 3: Using winget
```powershell
winget install GoLang.Go
```

## Quick Start (Basic Mode)

### Step 1: Open PowerShell
```powershell
cd C:\Users\Aditi\Desktop\Connect-4
```

### Step 2: Download Dependencies
```powershell
go mod download
```

### Step 3: Start the Server
```powershell
.\start.ps1
```

Or manually:
```powershell
go run ./server
```

### Step 4: Play!
Open your browser to: **http://localhost:8080**

That's it! 🎉

## Playing the Game

### Solo (vs Bot)
1. Enter your username
2. Click "Join Matchmaking"
3. Wait 10 seconds
4. Bot will join automatically

### Multiplayer (vs Friend)
1. Open two browser windows
2. Window 1: Enter "Alice" → Join
3. Window 2: Enter "Bob" → Join (within 10 seconds)
4. Play against each other!

### How to Play
- Click any column to drop your disc
- Connect 4 discs horizontally, vertically, or diagonally to win
- Red (Player 1) goes first

## Advanced Setup (Optional)

### With Docker (Kafka + PostgreSQL)

1. **Install Docker Desktop**: https://www.docker.com/products/docker-desktop

2. **Start Services**:
```powershell
docker-compose up -d
```

3. **Run Server with Full Features**:
```powershell
$env:DB_ENABLED="true"
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:DB_USER="postgres"
$env:DB_PASSWORD="postgres"
$env:DB_NAME="connect4"
$env:KAFKA_ENABLED="true"
$env:KAFKA_BROKERS="localhost:9092"
go run ./server
```

4. **Run Analytics Consumer** (in new terminal):
```powershell
go run ./cmd/analytics -brokers localhost:9092 -topic game-analytics
```

5. **View Kafka UI**: http://localhost:8090

## Troubleshooting

### "go is not recognized"
- Go is not installed or not in PATH
- Close and reopen PowerShell after installing Go
- Verify: `go version`

### "cannot find package"
Run:
```powershell
go mod tidy
go mod download
```

### Port 8080 in use
Change port:
```powershell
$env:SERVER_PORT="8081"
go run ./server
```

### Server won't start
1. Check if Go is installed: `go version`
2. Check if you're in the right directory: `pwd`
3. Check for errors in the terminal output

## What's Next?

- ✅ Play some games
- ✅ Check the leaderboard
- ✅ Try reconnecting (close tab, reopen within 30s)
- ⬜ Set up PostgreSQL for persistent storage
- ⬜ Set up Kafka for analytics
- ⬜ Build production executable: `go build -o connect4.exe ./server`

## File Structure

```
Connect-4/
├── server/          # Backend code
├── static/          # Frontend (HTML/CSS/JS)
├── cmd/analytics/   # Analytics consumer
├── data/            # Game data (auto-created)
├── start.ps1        # Quick start script
└── README.md        # Full documentation
```

## Commands Cheat Sheet

```powershell
# Start server (basic)
go run ./server

# Start server (with script)
.\start.ps1

# Download dependencies
go mod download

# Build executable
go build -o connect4.exe ./server

# Run executable
.\connect4.exe

# Start Docker services
docker-compose up -d

# Stop Docker services
docker-compose down

# View logs
docker-compose logs -f

# Run analytics consumer
go run ./cmd/analytics -brokers localhost:9092 -topic game-analytics
```

## Environment Variables

```powershell
# Server
$env:SERVER_PORT="8080"
$env:MATCH_TIMEOUT="10"
$env:RECONNECT_TIMEOUT="30"

# Database
$env:DB_ENABLED="true"
$env:DB_HOST="localhost"
$env:DB_USER="postgres"
$env:DB_PASSWORD="postgres"
$env:DB_NAME="connect4"

# Kafka
$env:KAFKA_ENABLED="true"
$env:KAFKA_BROKERS="localhost:9092"
$env:KAFKA_TOPIC="game-analytics"
```

## URLs

- **Game**: http://localhost:8080
- **Leaderboard API**: http://localhost:8080/leaderboard
- **Kafka UI** (if using Docker): http://localhost:8090

## Support

For detailed documentation:
- **Installation**: See `INSTALL.md`
- **Setup Guide**: See `SETUP.md`
- **API Documentation**: See `API.md`
- **Full README**: See `README.md`

---

**Have fun playing! 🎮🎉**

