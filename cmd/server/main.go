package main

import (
	"context"
	"fmt"

	"github.com/erupshis/effective_mobile/internal/helpers"
	"github.com/erupshis/effective_mobile/internal/msgbroker"
	"github.com/erupshis/effective_mobile/internal/server/config"
)

func main() {
	cfg := config.Parse()

	reader := msgbroker.CreateKafkaConsumer(cfg.BrokerAddr, cfg.Topic, cfg.Group)
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
