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
	TopicIn    string
	TopicError string
	Group      string
}

func Parse() Config {
	var config = Config{}
	checkFlags(&config)
	checkEnvironments(&config)
	return config
}

// FLAGS PARSING.
const (
	flagBrokers    = "br"
	flagTopicIn    = "tin"
	flagTopicError = "terr"
	flagGroup      = "g"
)

func checkFlags(config *Config) {
	var brokers string
	flag.StringVar(&brokers, flagBrokers, "localhost:9092", "kafka brokers with ',' separator between")
	config.BrokerAddr = strings.Split(brokers, ",")

	flag.StringVar(&config.TopicIn, flagTopicIn, "FIO", "kafka consumer topic")
	flag.StringVar(&config.TopicError, flagTopicError, "FIO_FAILED", "kafka producer topic(response in case of errors)")
	flag.StringVar(&config.Group, flagGroup, "groupServer", "kafka consumer group")
	flag.Parse()
}

// ENVIRONMENTS PARSING.
type envConfig struct {
	BrokerAddr string `env:"BROKERS"`
	TopicIn    string `env:"TOPIC_IN"`
	TopicError string `env:"TOPIC_ERROR"`
	Group      string `env:"GROUP"`
}

func checkEnvironments(config *Config) {
	var envs = envConfig{}
	err := env.Parse(&envs)
	if err != nil {
		log.Fatal(err)
	}

	confighelper.SetEnvToParamIfNeed(&config.BrokerAddr, envs.BrokerAddr)
	confighelper.SetEnvToParamIfNeed(&config.TopicIn, envs.TopicIn)
	confighelper.SetEnvToParamIfNeed(&config.TopicError, envs.TopicError)
	confighelper.SetEnvToParamIfNeed(&config.Group, envs.Group)
}
