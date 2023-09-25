package config // import "github.com/erupshis/effective_mobile/internal/agent/config"


TYPES

type Config struct {
	BrokerAddr []string
	TopicIn    string
	TopicOut   string
	Group      string
}

func Parse() Config

