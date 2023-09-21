package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/helpers"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/msgbroker"
	"github.com/erupshis/effective_mobile/internal/server/config"
	"github.com/erupshis/effective_mobile/internal/server/controllers/extradatactrl"
	"github.com/erupshis/effective_mobile/internal/server/controllers/msgbrokerctrl"
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
	go brokerWriter.Listen(ctxWithCancel, chMessageErrors)

	//storage.
	strg := storage.CreateRamStorage()

	//channels for filling missing persons data.
	chPartialPersonData := make(chan datastructs.PersonData, 10)
	chFullPersonData := make(chan datastructs.PersonData, 10)

	//message broker controller.
	msgController := msgbrokerctrl.Create(chMessages, chMessageErrors, chPartialPersonData, log)
	go msgController.Run(ctxWithCancel)

	//extra data controller.
	extraController := extradatactrl.Create(chPartialPersonData, chFullPersonData, strg, log)
	go extraController.Run(ctxWithCancel)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
}
