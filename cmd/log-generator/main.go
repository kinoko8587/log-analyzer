package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/pinjung/log-analyzer/internal/domain"
	"github.com/segmentio/kafka-go"
)

var (
	logLevels = []string{"INFO", "WARN", "ERROR", "DEBUG"}
	messages  = []string{
		"User login successful",
		"Database connection established",
		"API request processed",
		"Cache hit for key: user_123",
		"Failed to connect to external service",
		"Request timeout after 30s",
		"Invalid authentication token",
		"Data processing completed",
		"Background job started",
		"Memory usage exceeded threshold",
	}
)

func generateRandomLog() domain.Log {
	return domain.Log{
		UserID:    fmt.Sprintf("user_%d", rand.Intn(1000)+1),
		Timestamp: time.Now(),
		Level:     logLevels[rand.Intn(len(logLevels))],
		Message:   messages[rand.Intn(len(messages))],
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, relying on system env vars")
	}

	brokers := []string{os.Getenv("KAFKA_BROKERS")}
	topic := os.Getenv("KAFKA_TOPIC")

	if brokers[0] == "" || topic == "" {
		log.Fatal("❌ KAFKA_BROKERS and KAFKA_TOPIC must be set")
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	defer writer.Close()

	log.Printf("🚀 Starting log generator...")
	log.Printf("📡 Kafka brokers: %v", brokers)
	log.Printf("📬 Topic: %s", topic)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	logCount := 0

	for {
		select {
		case <-stop:
			log.Printf("🛑 Shutting down... Total logs sent: %d", logCount)
			return
		case <-ticker.C:
			batchSize := rand.Intn(5) + 1
			var messages []kafka.Message

			for i := 0; i < batchSize; i++ {
				logEntry := generateRandomLog()
				data, err := json.Marshal(logEntry)
				if err != nil {
					log.Printf("❌ Failed to marshal log: %v", err)
					continue
				}

				messages = append(messages, kafka.Message{
					Key:   []byte(logEntry.UserID),
					Value: data,
				})
			}

			err := writer.WriteMessages(ctx, messages...)
			if err != nil {
				log.Printf("❌ Failed to write messages: %v", err)
			} else {
				logCount += len(messages)
				log.Printf("✅ Sent %d logs (Total: %d)", len(messages), logCount)
			}
		}
	}
}