package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/pinjung/log-analyzer/internal/infrastructure/db"
	"github.com/pinjung/log-analyzer/internal/infrastructure/kafka"
)

func getPostgresDSN() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_SSL"),
	)
}

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

	connStr := getPostgresDSN()
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("DB connect error: %v", err)
	}
	defer dbConn.Close()

	repo := db.NewPostgresLogRepository(dbConn)

	kafkaBrokers, kafkaTopic, kafkaGroupID := getKafkaConfig()
	consumer := kafka.NewConsumer(
		kafkaBrokers,
		kafkaTopic,
		kafkaGroupID,
		repo,
	)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Fatalf("Consumer failed: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Println("🛑 Shutting down log-ingestor...")
	cancel()
}
