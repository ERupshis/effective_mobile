package main

import (
	"fmt"
	"os"

	"github.com/erupshis/effective_mobile/internal/agent/config"
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
	defer log.Sync()

	//kafka.
	writer := msgbroker.CreateKafkaProducer(cfg.BrokerAddr, cfg.Topic, log)
	defer helpers.ExecuteWithLogError(writer.Close, log)

	//TODO: need to add goroutine to generate random data.
	//TODO: need to wrap data in json.
	// Send a sample message
	key := "message-key"
	value := "Hello, Kafka!"
	err = writer.SendMessage(key, value)
	if err != nil {
		log.Info("send message failed: %v", err)
	}

	fmt.Printf("Message sent: Key=%s, Value=%s\n", key, value)

	ch := make(chan struct{})
	<-ch
}
