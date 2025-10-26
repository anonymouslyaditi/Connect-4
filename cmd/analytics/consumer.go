package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

// Simple Kafka consumer that reads messages from the 'game-analytics' topic and prints basic metrics.
func main() {
	brokers := flag.String("brokers", "localhost:9092", "kafka brokers comma separated")
	topic := flag.String("topic", "game-analytics", "kafka topic")
	group := flag.String("group", "analytics-consumer", "consumer group")
	flag.Parse()

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{*brokers},
		GroupID:  *group,
		Topic:    *topic,
		MinBytes: 1,
		MaxBytes: 10e6,
	})
	defer r.Close()

	log.Printf("Starting analytics consumer for topic %s (brokers=%s)", *topic, *brokers)
	metrics := NewMetrics()
	ctx := context.Background()
	for {
		m, err := r.ReadMessage(ctx)
		if err != nil {
			log.Println("read error:", err)
			time.Sleep(2 * time.Second)
			continue
		}
		var obj map[string]interface{}
		if err := json.Unmarshal(m.Value, &obj); err != nil {
			log.Println("invalid json in message")
			continue
		}
		metrics.ProcessEvent(obj)
		metrics.LogSnapshot()
	}
}

// Metrics stores comprehensive game analytics
type Metrics struct {
	TotalGames      int
	TotalMoves      int
	ByWinner        map[string]int
	ByHour          map[string]int
	ByDay           map[string]int
	UserStats       map[string]*UserStats
	TotalDuration   int64
	AverageDuration float64
	BotGames        int
	PlayerGames     int
	Draws           int
	Forfeits        int
}

// UserStats tracks per-user metrics
type UserStats struct {
	Username      string
	GamesPlayed   int
	Wins          int
	Losses        int
	Draws         int
	TotalDuration int64
	AvgDuration   float64
	WinRate       float64
}

func NewMetrics() *Metrics {
	return &Metrics{
		ByWinner:  map[string]int{},
		ByHour:    map[string]int{},
		ByDay:     map[string]int{},
		UserStats: map[string]*UserStats{},
	}
}

func (m *Metrics) ProcessEvent(e map[string]interface{}) {
	t := e["type"]
	switch t {
	case "game_finished":
		if g, ok := e["game"].(map[string]interface{}); ok {
			m.TotalGames++

			// Extract game data
			winner, _ := g["winner"].(string)
			player1, _ := g["player1"].(string)
			player2, _ := g["player2"].(string)
			duration, _ := g["duration_seconds"].(float64)
			durationInt := int64(duration)

			// Track winner
			m.ByWinner[winner]++

			// Track draws and forfeits
			if winner == "draw" {
				m.Draws++
			}
			if reason, ok := e["reason"].(string); ok && reason == "forfeit" {
				m.Forfeits++
			}

			// Track bot vs player games
			if player2 == "Bot" || player1 == "Bot" {
				m.BotGames++
			} else {
				m.PlayerGames++
			}

			// Track duration
			m.TotalDuration += durationInt
			if m.TotalGames > 0 {
				m.AverageDuration = float64(m.TotalDuration) / float64(m.TotalGames)
			}

			// Track by time
			if ended, ok := g["ended_at"].(string); ok {
				if tm, err := time.Parse(time.RFC3339, ended); err == nil {
					h := tm.Format("2006-01-02 15")
					m.ByHour[h]++
					d := tm.Format("2006-01-02")
					m.ByDay[d]++
				}
			}

			// Update user stats
			m.updateUserStats(player1, player2, winner, durationInt)
		}
	case "move":
		m.TotalMoves++
	default:
		// ignore other events
	}
}

func (m *Metrics) updateUserStats(player1, player2, winner string, duration int64) {
	// Skip bot
	players := []string{}
	if player1 != "Bot" {
		players = append(players, player1)
	}
	if player2 != "Bot" {
		players = append(players, player2)
	}

	for _, p := range players {
		if m.UserStats[p] == nil {
			m.UserStats[p] = &UserStats{Username: p}
		}
		stats := m.UserStats[p]
		stats.GamesPlayed++
		stats.TotalDuration += duration
		stats.AvgDuration = float64(stats.TotalDuration) / float64(stats.GamesPlayed)

		if winner == p {
			stats.Wins++
		} else if winner == "draw" {
			stats.Draws++
		} else {
			stats.Losses++
		}

		// Calculate win rate
		if stats.GamesPlayed > 0 {
			stats.WinRate = float64(stats.Wins) / float64(stats.GamesPlayed) * 100
		}
	}
}

func (m *Metrics) LogSnapshot() {
	fmt.Println("\n========== GAME ANALYTICS SNAPSHOT ==========")
	fmt.Printf("Total Games: %d\n", m.TotalGames)
	fmt.Printf("Total Moves: %d\n", m.TotalMoves)
	fmt.Printf("Average Game Duration: %.2f seconds\n", m.AverageDuration)
	fmt.Printf("Bot Games: %d | Player vs Player: %d\n", m.BotGames, m.PlayerGames)
	fmt.Printf("Draws: %d | Forfeits: %d\n", m.Draws, m.Forfeits)

	fmt.Println("\n--- Top Winners ---")
	for winner, count := range m.ByWinner {
		fmt.Printf("  %s: %d wins\n", winner, count)
	}

	fmt.Println("\n--- Games Per Day ---")
	for day, count := range m.ByDay {
		fmt.Printf("  %s: %d games\n", day, count)
	}

	fmt.Println("\n--- User Statistics ---")
	for username, stats := range m.UserStats {
		fmt.Printf("  %s: Played=%d, Wins=%d, Losses=%d, Draws=%d, WinRate=%.1f%%, AvgDuration=%.1fs\n",
			username, stats.GamesPlayed, stats.Wins, stats.Losses, stats.Draws, stats.WinRate, stats.AvgDuration)
	}

	fmt.Println("=============================================\n")
}
