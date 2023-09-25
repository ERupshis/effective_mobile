package config // import "github.com/erupshis/effective_mobile/internal/server/config"


TYPES

type Config struct {
	BrokerAddr  []string
	DatabaseDSN string
	CacheDSN    string
	Group       string
	Host        string
	TopicIn     string
	TopicError  string
}

func Parse() Config

