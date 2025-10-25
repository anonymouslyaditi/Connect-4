# Connect-4 Game Server Startup Script
# This script sets up environment variables and starts the server

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  4 in a Row - Game Server Startup" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Check if Go is installed
Write-Host "Checking for Go installation..." -ForegroundColor Yellow
$goVersion = & go version 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host "ERROR: Go is not installed or not in PATH!" -ForegroundColor Red
    Write-Host ""
    Write-Host "Please install Go from: https://go.dev/dl/" -ForegroundColor Yellow
    Write-Host "Or use Chocolatey: choco install golang" -ForegroundColor Yellow
    Write-Host "Or use winget: winget install GoLang.Go" -ForegroundColor Yellow
    Write-Host ""
    Write-Host "After installation, close and reopen PowerShell, then run this script again." -ForegroundColor Yellow
    Write-Host ""
    Read-Host "Press Enter to exit"
    exit 1
}

Write-Host "✓ Go is installed: $goVersion" -ForegroundColor Green
Write-Host ""

# Check if dependencies are downloaded
Write-Host "Checking dependencies..." -ForegroundColor Yellow
if (-not (Test-Path "go.sum")) {
    Write-Host "Downloading dependencies..." -ForegroundColor Yellow
    go mod download
    if ($LASTEXITCODE -ne 0) {
        Write-Host "ERROR: Failed to download dependencies!" -ForegroundColor Red
        Read-Host "Press Enter to exit"
        exit 1
    }
    Write-Host "✓ Dependencies downloaded" -ForegroundColor Green
} else {
    Write-Host "✓ Dependencies already present" -ForegroundColor Green
}
Write-Host ""

# Set environment variables
Write-Host "Configuring server..." -ForegroundColor Yellow

# Basic configuration
$env:SERVER_PORT = "8080"
$env:DATA_DIR = "data"
$env:MATCH_TIMEOUT = "10"
$env:RECONNECT_TIMEOUT = "30"

# Database configuration (disabled by default)
$env:DB_ENABLED = "false"
# Uncomment and configure if you have PostgreSQL:
# $env:DB_ENABLED = "true"
# $env:DB_HOST = "localhost"
# $env:DB_PORT = "5432"
# $env:DB_USER = "postgres"
# $env:DB_PASSWORD = "your_password"
# $env:DB_NAME = "connect4"

# Kafka configuration (disabled by default)
$env:KAFKA_ENABLED = "false"
# Uncomment and configure if you have Kafka:
# $env:KAFKA_ENABLED = "true"
# $env:KAFKA_BROKERS = "localhost:9092"
# $env:KAFKA_TOPIC = "game-analytics"

Write-Host "✓ Configuration loaded" -ForegroundColor Green
Write-Host ""

# Display configuration
Write-Host "Server Configuration:" -ForegroundColor Cyan
Write-Host "  Port: $env:SERVER_PORT" -ForegroundColor White
Write-Host "  Data Directory: $env:DATA_DIR" -ForegroundColor White
Write-Host "  Match Timeout: $env:MATCH_TIMEOUT seconds" -ForegroundColor White
Write-Host "  Reconnect Timeout: $env:RECONNECT_TIMEOUT seconds" -ForegroundColor White
Write-Host "  Database: $env:DB_ENABLED" -ForegroundColor White
Write-Host "  Kafka: $env:KAFKA_ENABLED" -ForegroundColor White
Write-Host ""

# Create data directory if it doesn't exist
if (-not (Test-Path $env:DATA_DIR)) {
    New-Item -ItemType Directory -Path $env:DATA_DIR | Out-Null
    Write-Host "✓ Created data directory" -ForegroundColor Green
}

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Starting server..." -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Server will be available at: http://localhost:$env:SERVER_PORT" -ForegroundColor Yellow
Write-Host "Press Ctrl+C to stop the server" -ForegroundColor Yellow
Write-Host ""

# Start the server
go run ./server

# If server exits
Write-Host ""
Write-Host "Server stopped." -ForegroundColor Yellow

