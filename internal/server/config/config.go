package config

import (
	"flag"
	"log"
	"strings"

	"github.com/caarlos0/env"
	"github.com/erupshis/effective_mobile/internal/confighelper"
)

type Config struct {
	BrokerAddr  []string
	DatabaseDSN string
	Group       string
	Host        string
	TopicIn     string
	TopicError  string
}

func Parse() Config {
	var config = Config{}
	checkFlags(&config)
	checkEnvironments(&config)
	return config
}

// FLAGS PARSING.
const (
	flagAddress     = "a"
	flagBrokers     = "br"
	flagDatabaseDSN = "d"
	flagGroup       = "g"
	flagTopicIn     = "tin"
	flagTopicError  = "terr"
)

func checkFlags(config *Config) {
	var brokers string
	flag.StringVar(&brokers, flagBrokers, "localhost:9092", "kafka brokers with ',' separator between")
	config.BrokerAddr = strings.Split(brokers, ",")

	flag.StringVar(&config.Group, flagGroup, "groupServer", "kafka consumer group")
	flag.StringVar(&config.Host, flagAddress, "localhost:8080", "server endpoint")
	flag.StringVar(&config.DatabaseDSN, flagDatabaseDSN, "postgres://postgres:postgres@localhost:5432/effective_mobile_db?sslmode=disable", "database DSN")
	flag.StringVar(&config.TopicIn, flagTopicIn, "FIO", "kafka consumer topic")
	flag.StringVar(&config.TopicError, flagTopicError, "FIO_FAILED", "kafka producer topic(response in case of errors)")

	flag.Parse()
}

// ENVIRONMENTS PARSING.
type envConfig struct {
	BrokerAddr  string `env:"BROKERS"`
	DatabaseDSN string `env:"DATABASE_DSN"`
	Group       string `env:"GROUP"`
	Host        string `env:"ADDRESS"`
	TopicIn     string `env:"TOPIC_IN"`
	TopicError  string `env:"TOPIC_ERROR"`
}

func checkEnvironments(config *Config) {
	var envs = envConfig{}
	err := env.Parse(&envs)
	if err != nil {
		log.Fatal(err)
	}

	confighelper.SetEnvToParamIfNeed(&config.BrokerAddr, envs.BrokerAddr)
	confighelper.SetEnvToParamIfNeed(&config.DatabaseDSN, envs.DatabaseDSN)
	confighelper.SetEnvToParamIfNeed(&config.Group, envs.Group)
	confighelper.SetEnvToParamIfNeed(&config.Host, envs.Host)
	confighelper.SetEnvToParamIfNeed(&config.TopicIn, envs.TopicIn)
	confighelper.SetEnvToParamIfNeed(&config.TopicError, envs.TopicError)
}
