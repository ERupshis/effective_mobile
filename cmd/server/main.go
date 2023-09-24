package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/erupshis/effective_mobile/internal/client"
	"github.com/erupshis/effective_mobile/internal/datastructs"
	"github.com/erupshis/effective_mobile/internal/helpers"
	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/erupshis/effective_mobile/internal/msgbroker"
	"github.com/erupshis/effective_mobile/internal/server/config"
	"github.com/erupshis/effective_mobile/internal/server/controllers/errorsctrl"
	"github.com/erupshis/effective_mobile/internal/server/controllers/extradatactrl"
	"github.com/erupshis/effective_mobile/internal/server/controllers/httpctrl"
	"github.com/erupshis/effective_mobile/internal/server/controllers/msgbrokerctrl"
	"github.com/erupshis/effective_mobile/internal/server/controllers/msgsavectrl"
	"github.com/erupshis/effective_mobile/internal/server/controllers/qraphqlctrl"
	"github.com/erupshis/effective_mobile/internal/server/storage"
	"github.com/erupshis/effective_mobile/internal/server/storage/cache"
	"github.com/erupshis/effective_mobile/internal/server/storage/cache/manager"
	"github.com/erupshis/effective_mobile/internal/server/storage/managers/postgresql"
	"github.com/go-chi/chi/v5"
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
	//manager.
	queriesHandler := postgresql.CreateHandler(log)
	storageManager, err := postgresql.CreatePostgreDB(ctxWithCancel, cfg, queriesHandler, log)
	if err != nil {
		log.Info("failed to connect to storage: %v", err)
		return
	}
	defer helpers.ExecuteWithLogError(storageManager.Close, log)

	mainStorage := storage.Create(storageManager, log)

	//cacheManager.
	cacheManager, err := manager.CreateRedis(ctxWithCancel, cfg.CacheDSN, log)
	if err != nil {
		log.Info("failed to create cache: %v", err)
		return
	}
	defer helpers.ExecuteWithLogError(cacheManager.Close, log)

	//storage with cache as Decorator.
	cachedStorage := cache.Create(cacheManager, mainStorage, log)

	//save messages controller.
	chErrorsSaveCtrl := make(chan msgbroker.Message, 10)
	saveController := msgsavectrl.Create(chFullPersonData, chErrorsSaveCtrl, cachedStorage, log)
	go saveController.Run(ctxWithCancel)

	//errors controller.
	errorsController := errorsctrl.Create([]<-chan msgbroker.Message{chErrorsBrokerCtrl, chErrorsExtraCtrl, chErrorsSaveCtrl}, chMessageErrors, log)
	go errorsController.Run(ctxWithCancel)

	//http controller.
	httpController := httpctrl.Create(cachedStorage, log)

	//graphQL controller.
	graphqlController := qraphqlctrl.Create(cachedStorage, log)

	//rest routing.
	router := chi.NewRouter()
	router.Mount("/", httpController.Route())
	router.Mount("/graphql", graphqlController.Route())

	log.Info("server is launching with Host setting: %s", cfg.Host)
	if err := http.ListenAndServe(cfg.Host, router); err != nil {
		log.Info("server refused to start with error: %v", err)
		panic(err)
	}
}
