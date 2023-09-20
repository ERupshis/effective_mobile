package msgbroker

import (
	"context"
)

//go:generate easyjson -all interface.go
type MessageBody struct {
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
}

type Message struct {
	Key   []byte
	Value []byte
}

type Producer interface {
	SendMessage(key, value string) error
	Close() error
}

type Consumer interface {
	ReadMessage(ctx context.Context) (Message, error)
	Close() error
}
