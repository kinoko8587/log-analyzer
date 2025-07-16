package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	"github.com/segmentio/kafka-go"
)

// Log defines the mock log format
type Log struct {
	UserID    string    `json:"user_id"`
	Timestamp time.Time `json:"timestamp"`
	Level     string    `json:"level"` // info, warn, error
	Message   string    `json:"message"`
}

var levels = []string{"info", "warn", "error"}
var messages = []string{
	"User login",
	"Payment timeout",
	"Internal server error",
	"Request completed",
	"DB query slow",
	"Unauthorized access",
}

func generateLog() Log {
	return Log{
		UserID:    randomUser(),
		Timestamp: time.Now(),
		Level:     levels[rand.Intn(len(levels))],
		Message:   messages[rand.Intn(len(messages))],
	}
}

func randomUser() string {
	return "user_" + string(rune(rand.Intn(100)+65)) // user_A ~ user_ZZ
}

func main() {
	rand.Seed(time.Now().UnixNano())

	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "logs.ingest",
		Async:   true,
	})
	defer kafkaWriter.Close()

	log.Println("🚀 Log generator started...")
	ticker := time.NewTicker(10 * time.Millisecond) // ~100 logs/sec
	defer ticker.Stop()

	for {
		<-ticker.C
		logData := generateLog()

		value, _ := json.Marshal(logData)

		err := kafkaWriter.WriteMessages(context.Background(), kafka.Message{
			Key:   []byte(logData.UserID),
			Value: value,
		})

		if err != nil {
			log.Printf("❌ failed to write: %v\n", err)
		}
	}
}
