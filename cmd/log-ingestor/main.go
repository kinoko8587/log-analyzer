package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/pinjung/log-analyzer/internal/infrastructure/db"
	"github.com/pinjung/log-analyzer/internal/infrastructure/kafka"
)

func getKafkaConfig() (brokers []string, topic string, groupID string) {
	brokers = []string{os.Getenv("KAFKA_BROKERS")} // 支援多 broker 時可改 split
	topic = os.Getenv("KAFKA_TOPIC")
	groupID = os.Getenv("KAFKA_GROUP_ID")
	return
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  No .env file found, relying on system env vars")
	}

	bunDB := db.NewBunDB()
	defer bunDB.Close()

	repo := db.NewPostgresLogRepository(bunDB)

	kafkaBrokers, kafkaTopic, kafkaGroupID := getKafkaConfig()
	consumer := kafka.NewConsumer(
		kafkaBrokers,
		kafkaTopic,
		kafkaGroupID,
		repo,
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errChan := make(chan error, 1)
	go func() {
		if err := consumer.Start(ctx); err != nil {
			errChan <- err
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-stop:
		log.Println("🛑 Received shutdown signal...")
	case err := <-errChan:
		log.Printf("❌ Consumer error: %v", err)
	}

	log.Println("🛑 Shutting down log-ingestor...")
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	shutdownDone := make(chan struct{})
	go func() {
		if err := consumer.Close(); err != nil {
			log.Printf("❌ Error closing consumer: %v", err)
		}
		close(shutdownDone)
	}()

	select {
	case <-shutdownDone:
		log.Println("✅ Graceful shutdown completed")
	case <-shutdownCtx.Done():
		log.Println("⚠️  Shutdown timeout exceeded")
	}
}
