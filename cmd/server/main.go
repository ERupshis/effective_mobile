package main

import (
	"context"
	"fmt"
	"os"

	"github.com/erupshis/effective_mobile/internal/helpers"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/msgbroker"
	"github.com/erupshis/effective_mobile/internal/server/config"
)

func main() {
	//config.
	cfg := config.Parse()

	//log.
	log, err := logger.CreateZapLogger("info")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create logger: %v", err)
	}
	defer log.Sync()

	//kafka
	reader := msgbroker.CreateKafkaConsumer(cfg.BrokerAddr, cfg.Topic, cfg.Group)
	defer helpers.ExecuteWithLogError(reader.Close, log)

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
