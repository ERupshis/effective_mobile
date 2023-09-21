package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/erupshis/effective_mobile/internal/client"
	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/helpers"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/msgbroker"
	"github.com/erupshis/effective_mobile/internal/server/config"
	"github.com/erupshis/effective_mobile/internal/server/controllers/errorsctrl"
	"github.com/erupshis/effective_mobile/internal/server/controllers/extradatactrl"
	"github.com/erupshis/effective_mobile/internal/server/controllers/msgbrokerctrl"
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
	//strg := storage.CreateRamStorage()

	//channels for filling missing persons data.
	chPartialPersonData := make(chan datastructs.ExtraDataFilling, 10)
	chFullPersonData := make(chan datastructs.ExtraDataFilling, 10)

	//message broker controller.
	chErrorsBrokerCtrl := make(chan msgbroker.Message, 10)
	msgController := msgbrokerctrl.Create(chMessages, chErrorsBrokerCtrl, chPartialPersonData, log)
	go msgController.Run(ctxWithCancel)

	//extra data controller.
	chErrorsExtraCtrl := make(chan msgbroker.Message, 10)
	clientForRemoteAPI := client.CreateDefault(log)
	extraController := extradatactrl.Create(chPartialPersonData, chFullPersonData, chErrorsExtraCtrl, clientForRemoteAPI, log)
	go extraController.Run(ctxWithCancel)

	//errors controller.
	errorsController := errorsctrl.Create([]<-chan msgbroker.Message{chErrorsBrokerCtrl, chErrorsExtraCtrl}, chMessageErrors, log)
	go errorsController.Run(ctxWithCancel)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
}
