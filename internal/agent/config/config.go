// Package config agent's setting parser. Applies flags and environments. Environments are prioritized.
package config

import (
	"flag"
	"log"
	"strings"

	"github.com/caarlos0/env"
	"github.com/erupshis/effective_mobile/internal/confighelper"
)

// Config agent's settings.
type Config struct {
	BrokerAddr []string
	TopicIn    string
	TopicOut   string
	Group      string
}

// Parse main func to parse variables.
func Parse() Config {
	var config = Config{}
	checkFlags(&config)
	checkEnvironments(&config)
	return config
}

// FLAGS PARSING.
const (
	flagBrokers  = "br"
	flagTopicIn  = "tin"
	flagTopicOut = "tout"
	flagGroup    = "g"
)

// checkFlags checks flags of app's launch.
func checkFlags(config *Config) {
	//kafka.
	var brokers string
	flag.StringVar(&brokers, flagBrokers, "localhost:9092", "kafka brokers with ',' separator between")
	config.BrokerAddr = strings.Split(brokers, ",")

	flag.StringVar(&config.TopicIn, flagTopicIn, "FIO_FAILED", "kafka consumer topic for incoming errors")
	flag.StringVar(&config.TopicOut, flagTopicOut, "FIO", "kafka producer topic")

	flag.StringVar(&config.Group, flagGroup, "groupAgent", "kafka consumer group")
	flag.Parse()
}

// ENVIRONMENTS PARSING.
// envConfig struct of environments suitable for agent.
type envConfig struct {
	BrokerAddr string `env:"BROKERS"`
	TopicError string `env:"TOPIC_ERROR"`
	TopicOut   string `env:"TOPIC_OUT"`
	Group      string `env:"GROUP"`
}

// checkEnvironments checks environments suitable for agent.
func checkEnvironments(config *Config) {
	var envs = envConfig{}
	err := env.Parse(&envs)
	if err != nil {
		log.Fatal(err)
	}

	//kafka.
	_ = confighelper.SetEnvToParamIfNeed(&config.BrokerAddr, envs.BrokerAddr)
	_ = confighelper.SetEnvToParamIfNeed(&config.TopicIn, envs.TopicError)
	_ = confighelper.SetEnvToParamIfNeed(&config.TopicOut, envs.TopicOut)
	_ = confighelper.SetEnvToParamIfNeed(&config.Group, envs.Group)
}
