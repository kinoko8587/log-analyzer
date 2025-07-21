package domain

import "context"

type Message struct {
	Key   []byte
	Value []byte
}

type QueueSubscriber interface {
	Subscribe(ctx context.Context, handler func(Message) error) error
	Close() error
}