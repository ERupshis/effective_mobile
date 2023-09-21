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
	brokerWriter := msgbroker.CreateKafkaProducer(cfg.BrokerAddr, cfg.TopicOut, log)
	defer helpers.ExecuteWithLogError(brokerWriter.Close, log)

	brokerReader := msgbroker.CreateKafkaConsumer(cfg.BrokerAddr, cfg.TopicIn, cfg.Group, log)
	defer helpers.ExecuteWithLogError(brokerReader.Close, log)

	ctxWithCancel, cancel := context.WithCancel(context.Background())
	defer cancel()

	chMessageErrors := make(chan msgbroker.Message, 10)
	go brokerReader.Listen(ctxWithCancel, chMessageErrors)
	go func(ctx context.Context, chOut <-chan msgbroker.Message, log logger.BaseLogger) {
		for {
			select {
			case <-ctx.Done():
				return
			case _, ok := <-chOut:
				if !ok {
					return
				}
				//some work.
			}
		}
	}(ctxWithCancel, chMessageErrors, log)

	//random names generator.
	go msggenerator.Run(ctxWithCancel, brokerWriter, log)

	// Create a channel to wait for signals (e.g., Ctrl+C) to gracefully exit.
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
}
