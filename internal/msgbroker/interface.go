// Package msgbroker implements message transactions.
package msgbroker

import (
	"context"
)

// Message structure of message.
type Message struct {
	Key   []byte
	Value []byte
}

// Producer interface of messages sender.
//
//go:generate mockgen -destination=../../mocks/mock_Producer.go -package=mocks github.com/erupshis/effective_mobile/internal/msgbroker Producer
type Producer interface {
	// Listen goroutine method for listening response messages.
	Listen(ctx context.Context, chMessages <-chan Message)

	// SendMessage method for sending messages.
	SendMessage(ctx context.Context, key, value string) error

	// Close flushes pending sends and waits for all writes to complete before returning.
	Close() error
}

// Consumer interface of messages reader.
//
//go:generate mockgen -destination=../../mocks/mock_Consumer.go -package=mocks github.com/erupshis/effective_mobile/internal/msgbroker Consumer
type Consumer interface {
	// Listen goroutine method for listening to send response messages.
	Listen(ctx context.Context, chMessages chan<- Message)

	// ReadMessage method for reading incoming messages.
	ReadMessage(ctx context.Context) (Message, error)

	// Close closes the stream, preventing the program from reading any more messages from it.
	Close() error
}
