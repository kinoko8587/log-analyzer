package kafka

import (
	"context"
	"encoding/json"
	"log"

	"github.com/pinjung/log-analyzer/internal/domain"
	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	Reader     *kafka.Reader
	Repository domain.LogRepository
}

func NewConsumer(brokers []string, topic string, groupID string, repo domain.LogRepository) *Consumer {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})
	return &Consumer{Reader: r, Repository: repo}
}

func (c *Consumer) Start(ctx context.Context) error {
	return c.Subscribe(ctx, func(msg domain.Message) error {
		var l domain.Log
		if err := json.Unmarshal(msg.Value, &l); err != nil {
			log.Printf("❌ Failed to unmarshal log: %v", err)
			return err
		}
		log.Printf("📥 Received log: UserID=%s, Level=%s, Message=%s", l.UserID, l.Level, l.Message)
		if err := c.Repository.SaveLog(&l); err != nil {
			log.Printf("❌ Failed to save log: %v", err)
			return err
		}
		log.Printf("✅ Saved log to database")
		return nil
	})
}

func (c *Consumer) Subscribe(ctx context.Context, handler func(domain.Message) error) error {
	log.Println("🟢 Kafka consumer started...")

	for {
		select {
		case <-ctx.Done():
			log.Println("🛑 Context cancelled, stopping consumer...")
			return ctx.Err()
		default:
			m, err := c.Reader.ReadMessage(ctx)
			if err != nil {
				log.Printf("❌ read error: %v", err)
				continue
			}

			msg := domain.Message{
				Key:   m.Key,
				Value: m.Value,
			}

			if err := handler(msg); err != nil {
				log.Printf("❌ handler error: %v", err)
			}
		}
	}
}

func (c *Consumer) Close() error {
	return c.Reader.Close()
}
