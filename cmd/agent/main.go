package main

import (
	"fmt"
	"log"

	"github.com/ERupshis/effective_mobile/internal/helpers"
	"github.com/ERupshis/effective_mobile/internal/msgbroker"
)

func main() {
	//TODO envs config.
	brokerAddr := []string{"localhost:9092"}
	topic := "FIO"

	writer := msgbroker.CreateKafkaProducer(brokerAddr, topic, 0)
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
