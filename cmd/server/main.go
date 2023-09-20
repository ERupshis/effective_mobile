package main

import (
	"context"
	"fmt"

	"github.com/ERupshis/effective_mobile/internal/helpers"
	"github.com/ERupshis/effective_mobile/internal/msgbroker"
)

func main() {
	//TODO envs config.
	brokerAddr := []string{"localhost:9092"}
	topic := "FIO"
	groupID := "group"

	reader := msgbroker.CreateKafkaConsumer(brokerAddr, topic, groupID)
	defer helpers.ExecuteWithLogError(reader.Close)

	for {
		m, err := reader.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message: %s = %s\n", string(m.Key), string(m.Value))
	}

	//if err := r.Close(); err != nil {
	//	log.Fatal("failed to close reader:", err)
	//}
}
