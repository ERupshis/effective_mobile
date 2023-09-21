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
	"github.com/erupshis/effective_mobile/internal/server/controllers/msgsavectrl"
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

	//storage.
	strg, err := storage.CreatePostgreDB(ctxWithCancel, cfg, log)
	if err != nil {
		log.Info("failed to connect to storage: %v", err)
		return
	}

	//save messages controller.
	chErrorsSaveCtrl := make(chan msgbroker.Message, 10)
	saveController := msgsavectrl.Create(chFullPersonData, chErrorsSaveCtrl, strg, log)
	go saveController.Run(ctxWithCancel)

	//errors controller.
	errorsController := errorsctrl.Create([]<-chan msgbroker.Message{chErrorsBrokerCtrl, chErrorsExtraCtrl, chErrorsSaveCtrl}, chMessageErrors, log)
	go errorsController.Run(ctxWithCancel)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh
}
