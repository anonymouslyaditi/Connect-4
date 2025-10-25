package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

// Database wraps a SQL database connection
type Database struct {
	db      *sql.DB
	enabled bool
}

// NewDatabase creates a new database connection
func NewDatabase(config *Config) *Database {
	if !config.DBEnabled {
		log.Println("Database is disabled, using file-based storage")
		return &Database{enabled: false}
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		config.DBHost, config.DBPort, config.DBUser, config.DBPassword, config.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Failed to connect to database: %v. Falling back to file storage.", err)
		return &Database{enabled: false}
	}

	// Test connection
	if err := db.Ping(); err != nil {
		log.Printf("Failed to ping database: %v. Falling back to file storage.", err)
		return &Database{enabled: false}
	}

	// Initialize schema
	if err := initSchema(db); err != nil {
		log.Printf("Failed to initialize schema: %v. Falling back to file storage.", err)
		return &Database{enabled: false}
	}

	log.Println("Database connection established successfully")
	return &Database{db: db, enabled: true}
}

// initSchema creates the necessary tables if they don't exist
func initSchema(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS games (
		id VARCHAR(255) PRIMARY KEY,
		player1 VARCHAR(255) NOT NULL,
		player2 VARCHAR(255) NOT NULL,
		winner VARCHAR(255) NOT NULL,
		duration_seconds BIGINT NOT NULL,
		started_at TIMESTAMP NOT NULL,
		ended_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_games_player1 ON games(player1);
	CREATE INDEX IF NOT EXISTS idx_games_player2 ON games(player2);
	CREATE INDEX IF NOT EXISTS idx_games_winner ON games(winner);
	CREATE INDEX IF NOT EXISTS idx_games_ended_at ON games(ended_at);

	CREATE TABLE IF NOT EXISTS leaderboard (
		username VARCHAR(255) PRIMARY KEY,
		wins INT NOT NULL DEFAULT 0,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(schema)
	return err
}

// SaveGame saves a completed game to the database
func (d *Database) SaveGame(rec GameRecord) error {
	if !d.enabled {
		return nil
	}

	query := `
		INSERT INTO games (id, player1, player2, winner, duration_seconds, started_at, ended_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := d.db.Exec(query, rec.ID, rec.Player1, rec.Player2, rec.Winner, rec.Duration, rec.StartedAt, rec.EndedAt)
	if err != nil {
		log.Printf("Failed to save game to database: %v", err)
		return err
	}

	return nil
}

// IncrementWinner increments the win count for a player
func (d *Database) IncrementWinner(username string) error {
	if !d.enabled || username == "draw" || username == "" {
		return nil
	}

	query := `
		INSERT INTO leaderboard (username, wins, updated_at)
		VALUES ($1, 1, $2)
		ON CONFLICT (username)
		DO UPDATE SET wins = leaderboard.wins + 1, updated_at = $2
	`

	_, err := d.db.Exec(query, username, time.Now())
	if err != nil {
		log.Printf("Failed to increment winner in database: %v", err)
		return err
	}

	return nil
}

// GetLeaderboard retrieves the leaderboard from the database
func (d *Database) GetLeaderboard() (Leaderboard, error) {
	if !d.enabled {
		return nil, fmt.Errorf("database not enabled")
	}

	query := `SELECT username, wins FROM leaderboard ORDER BY wins DESC`
	rows, err := d.db.Query(query)
	if err != nil {
		log.Printf("Failed to query leaderboard: %v", err)
		return nil, err
	}
	defer rows.Close()

	lb := make(Leaderboard)
	for rows.Next() {
		var username string
		var wins int
		if err := rows.Scan(&username, &wins); err != nil {
			continue
		}
		lb[username] = wins
	}

	return lb, nil
}

// GetGames retrieves all games from the database
func (d *Database) GetGames() ([]GameRecord, error) {
	if !d.enabled {
		return nil, fmt.Errorf("database not enabled")
	}

	query := `SELECT id, player1, player2, winner, duration_seconds, started_at, ended_at FROM games ORDER BY ended_at DESC`
	rows, err := d.db.Query(query)
	if err != nil {
		log.Printf("Failed to query games: %v", err)
		return nil, err
	}
	defer rows.Close()

	var games []GameRecord
	for rows.Next() {
		var rec GameRecord
		if err := rows.Scan(&rec.ID, &rec.Player1, &rec.Player2, &rec.Winner, &rec.Duration, &rec.StartedAt, &rec.EndedAt); err != nil {
			continue
		}
		games = append(games, rec)
	}

	return games, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	if d.enabled && d.db != nil {
		return d.db.Close()
	}
	return nil
}

