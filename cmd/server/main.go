package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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
	reader := msgbroker.CreateKafkaConsumer(cfg.BrokerAddr, cfg.Topic, cfg.Group, log)
	defer helpers.ExecuteWithLogError(reader.Close, log)

	//kafka message reader.
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	defer cancel()

	chMessages := make(chan msgbroker.Message, 10)
	go reader.Listen(ctxWithCancel, chMessages)
	for message := range chMessages {
		fmt.Printf("%s\n", message.Value)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
}
