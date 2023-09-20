package config

import (
	"flag"
	"log"
	"strings"

	"github.com/caarlos0/env"
	"github.com/erupshis/effective_mobile/internal/confighelper"
)

type Config struct {
	BrokerAddr []string
	Topic      string
}

func Parse() Config {
	var config = Config{}
	checkFlags(&config)
	checkEnvironments(&config)
	return config
}

// FLAGS PARSING.
const (
	flagBrokers = "br"
	flagTopic   = "tin"
)

func checkFlags(config *Config) {
	var brokers string
	flag.StringVar(&brokers, flagBrokers, "localhost:9092", "kafka brokers with ',' separator between")
	config.BrokerAddr = strings.Split(brokers, ",")

	flag.StringVar(&config.Topic, flagTopic, "FIO", "kafka producer topic")
	flag.Parse()
}

// ENVIRONMENTS PARSING.
type envConfig struct {
	BrokerAddr string `env:"BROKERS"`
	Topic      string `env:"TOPIC"`
}

func checkEnvironments(config *Config) {
	var envs = envConfig{}
	err := env.Parse(&envs)
	if err != nil {
		log.Fatal(err)
	}

	confighelper.SetEnvToParamIfNeed(&config.BrokerAddr, envs.BrokerAddr)
	confighelper.SetEnvToParamIfNeed(&config.Topic, envs.Topic)
}
