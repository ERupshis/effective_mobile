package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/erupshis/effective_mobile/internal/agent/config"
	"github.com/erupshis/effective_mobile/internal/agent/msggenerator"
	"github.com/erupshis/effective_mobile/internal/helpers"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/msgbroker"
)

func main() {
	//config.
	cfg := config.Parse()

	//log.
	log, err := logger.CreateZapLogger("info")
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to create logger: %v", err)
	}

	//kafka.
	writer := msgbroker.CreateKafkaProducer(cfg.BrokerAddr, cfg.Topic, log)
	defer helpers.ExecuteWithLogError(writer.Close, log)

	//random names generator.
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	defer cancel()

	go msggenerator.Run(ctxWithCancel, writer, log)

	// Create a channel to wait for signals (e.g., Ctrl+C) to gracefully exit.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
}
