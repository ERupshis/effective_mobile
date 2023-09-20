package msgbroker

import (
	"context"
)

type Message struct {
	Key   []byte
	Value []byte
}

type Producer interface {
	SendMessage(key, value string) error
	Close() error
}

type Consumer interface {
	Listen(ctx context.Context, chMessages chan<- Message)
	ReadMessage(ctx context.Context) (Message, error)
	Close() error
}
