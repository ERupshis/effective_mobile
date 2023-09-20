package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/helpers"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/msgbroker"
	"github.com/erupshis/effective_mobile/internal/server/config"
	"github.com/erupshis/effective_mobile/internal/server/storage"
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

	//kafka.
	brokerReader := msgbroker.CreateKafkaConsumer(cfg.BrokerAddr, cfg.TopicIn, cfg.Group, log)
	defer helpers.ExecuteWithLogError(brokerReader.Close, log)

	brokerWriter := msgbroker.CreateKafkaProducer(cfg.BrokerAddr, cfg.TopicError, log)
	defer helpers.ExecuteWithLogError(brokerReader.Close, log)

	ctxWithCancel, cancel := context.WithCancel(context.Background())
	defer cancel()

	chMessages := make(chan msgbroker.Message, 10)
	go brokerReader.Listen(ctxWithCancel, chMessages)

	chMessageErrors := make(chan msgbroker.Message, 10)
	defer close(chMessageErrors)
	go brokerWriter.Listen(ctxWithCancel, chMessageErrors)

	//storage.
	strg := storage.CreateRamStorage()

	//msgbrokercontroller.

	for message := range chMessages {
		personData := datastructs.PersonData{}
		if err := json.Unmarshal(message.Value, &personData); err != nil {
			log.Info("failed to parse JSON message: %s, error: %v", message.Value, err)
			continue
		}

		if err = strg.SavePersonData(&personData); err != nil {
			log.Info("storage save value fail. message: %s, error: %v", message.Value, err)
		}
		log.Info("%s\n", message.Value)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
}
