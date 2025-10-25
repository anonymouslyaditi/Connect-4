# Installation Guide for Windows

## Step 1: Install Go

### Option A: Using Chocolatey (Recommended)
If you have Chocolatey installed:

```powershell
choco install golang
```

### Option B: Manual Installation

1. **Download Go:**
   - Visit: https://go.dev/dl/
   - Download: `go1.21.x.windows-amd64.msi` (latest version)

2. **Run the Installer:**
   - Double-click the downloaded `.msi` file
   - Follow the installation wizard
   - Default installation path: `C:\Program Files\Go`

3. **Verify Installation:**
   Open a **NEW** PowerShell window and run:
   ```powershell
   go version
   ```
   
   You should see something like: `go version go1.21.x windows/amd64`

### Option C: Using winget (Windows 10/11)

```powershell
winget install GoLang.Go
```

## Step 2: Install Project Dependencies

After Go is installed, open PowerShell in the project directory:

```powershell
cd C:\Users\Aditi\Desktop\Connect-4
```

Then run:

```powershell
go mod download
```

This will download:
- `github.com/gorilla/websocket` - WebSocket support
- `github.com/segmentio/kafka-go` - Kafka client
- `github.com/lib/pq` - PostgreSQL driver

## Step 3: Verify Installation

Run this command to ensure everything is set up:

```powershell
go mod verify
```

## Step 4: Run the Server

### Basic Mode (No Database/Kafka)

```powershell
go run ./server
```

The server will start on http://localhost:8080

### With Environment Variables

```powershell
$env:SERVER_PORT="8080"
$env:MATCH_TIMEOUT="10"
$env:RECONNECT_TIMEOUT="30"
go run ./server
```

## Step 5: Test the Application

1. Open your browser to: http://localhost:8080
2. Enter a username and click "Join Matchmaking"
3. Wait 10 seconds for the bot to join, or open another browser window to play against yourself

## Optional: Install PostgreSQL

### Using Chocolatey:
```powershell
choco install postgresql
```

### Manual Installation:
1. Download from: https://www.postgresql.org/download/windows/
2. Run the installer
3. Remember the password you set for the `postgres` user

### Setup Database:
```powershell
# Connect to PostgreSQL
psql -U postgres

# Create database
CREATE DATABASE connect4;

# Exit
\q
```

### Run Server with Database:
```powershell
$env:DB_ENABLED="true"
$env:DB_HOST="localhost"
$env:DB_PORT="5432"
$env:DB_USER="postgres"
$env:DB_PASSWORD="your_password"
$env:DB_NAME="connect4"
go run ./server
```

## Optional: Install Kafka

### Using Docker (Easiest):
```powershell
# Install Docker Desktop from: https://www.docker.com/products/docker-desktop

# Run Kafka with Docker Compose
docker-compose up -d
```

Create a `docker-compose.yml` file in the project root:

```yaml
version: '3'
services:
  zookeeper:
    image: confluentinc/cp-zookeeper:latest
    environment:
      ZOOKEEPER_CLIENT_PORT: 2181
      ZOOKEEPER_TICK_TIME: 2000
    ports:
      - 2181:2181

  kafka:
    image: confluentinc/cp-kafka:latest
    depends_on:
      - zookeeper
    ports:
      - 9092:9092
    environment:
      KAFKA_BROKER_ID: 1
      KAFKA_ZOOKEEPER_CONNECT: zookeeper:2181
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
```

### Run Server with Kafka:
```powershell
$env:KAFKA_ENABLED="true"
$env:KAFKA_BROKERS="localhost:9092"
$env:KAFKA_TOPIC="game-analytics"
go run ./server
```

### Run Analytics Consumer:
```powershell
go run ./cmd/analytics -brokers localhost:9092 -topic game-analytics
```

## Troubleshooting

### "go: command not found" or "go is not recognized"
- Close and reopen PowerShell after installing Go
- Check if Go is in PATH: `$env:PATH`
- Manually add to PATH if needed:
  ```powershell
  $env:PATH += ";C:\Program Files\Go\bin"
  ```

### "cannot find package"
- Run: `go mod tidy`
- Then: `go mod download`

### Port 8080 already in use
- Change port: `$env:SERVER_PORT="8081"`
- Or find and stop the process using port 8080

### Permission Denied
- Run PowerShell as Administrator
- Or change the data directory to a user-writable location

## Quick Start Script

Save this as `start.ps1`:

```powershell
# Quick start script for Connect-4 server

Write-Host "Starting Connect-4 Game Server..." -ForegroundColor Green

# Set environment variables
$env:SERVER_PORT="8080"
$env:DATA_DIR="data"
$env:MATCH_TIMEOUT="10"
$env:RECONNECT_TIMEOUT="30"

# Optional: Enable database (uncomment if you have PostgreSQL)
# $env:DB_ENABLED="true"
# $env:DB_HOST="localhost"
# $env:DB_PORT="5432"
# $env:DB_USER="postgres"
# $env:DB_PASSWORD="your_password"
# $env:DB_NAME="connect4"

# Optional: Enable Kafka (uncomment if you have Kafka)
# $env:KAFKA_ENABLED="true"
# $env:KAFKA_BROKERS="localhost:9092"
# $env:KAFKA_TOPIC="game-analytics"

Write-Host "Configuration loaded" -ForegroundColor Yellow
Write-Host "Server will start on http://localhost:$env:SERVER_PORT" -ForegroundColor Yellow

# Run the server
go run ./server
```

Run it with:
```powershell
.\start.ps1
```

## Building Executable

To create a standalone executable:

```powershell
# Build for Windows
go build -o connect4-server.exe ./server

# Run the executable
.\connect4-server.exe
```

The executable can be distributed and run without Go installed.

## Next Steps

1. ✅ Install Go
2. ✅ Download dependencies
3. ✅ Run the server
4. ✅ Test in browser
5. ⬜ (Optional) Install PostgreSQL
6. ⬜ (Optional) Install Kafka
7. ⬜ (Optional) Build production executable

## Support

If you encounter issues:
1. Check the SETUP.md file
2. Review the README.md
3. Check server logs for errors
4. Ensure all prerequisites are installed

---

**Ready to play! 🎮**

