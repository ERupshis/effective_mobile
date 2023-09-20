package main

import (
	"fmt"
	"log"

	"github.com/erupshis/effective_mobile/internal/agent/config"
	"github.com/erupshis/effective_mobile/internal/helpers"
	"github.com/erupshis/effective_mobile/internal/msgbroker"
)

func main() {
	cfg := config.Parse()

	writer := msgbroker.CreateKafkaProducer(cfg.BrokerAddr, cfg.Topic, 0)
	defer helpers.ExecuteWithLogError(writer.Close)

	//TODO: need to add goroutine to generate random data.
	//TODO: need to wrap data in json.
	// Send a sample message
	key := "message-key"
	value := "Hello, Kafka!"
	err := writer.SendMessage(key, value)
	if err != nil {
		log.Fatalf("Error sending message to Kafka: %v", err)
	}

	fmt.Printf("Message sent: Key=%s, Value=%s\n", key, value)

	ch := make(chan struct{})
	<-ch
}
