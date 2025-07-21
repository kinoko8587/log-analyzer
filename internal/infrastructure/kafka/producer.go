package kafka

import (
	"context"
	"encoding/json"

	"github.com/pinjung/log-analyzer/internal/domain"
	"github.com/segmentio/kafka-go"
)

type Producer struct {
	Writer *kafka.Writer
}

func NewProducer(brokers []string, topic string) *Producer {
	w := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	})
	return &Producer{Writer: w}
}

func (p *Producer) Publish(ctx context.Context, logs []domain.Log) error {
	messages := make([]kafka.Message, 0, len(logs))
	
	for _, log := range logs {
		data, err := json.Marshal(log)
		if err != nil {
			return err
		}
		
		messages = append(messages, kafka.Message{
			Key:   []byte(log.UserID),
			Value: data,
		})
	}
	
	return p.Writer.WriteMessages(ctx, messages...)
}

func (p *Producer) Close() error {
	return p.Writer.Close()
}