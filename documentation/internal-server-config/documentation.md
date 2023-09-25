package config // import "github.com/erupshis/effective_mobile/internal/server/config"

Package config server's setting parser. Applies flags and environments.
Environments are prioritized.

TYPES

type Config struct {
	BrokerAddr  []string // BrokerAddr kafka's broker.
	DatabaseDSN string   // DatabaseDSN PostgreSQL data source name.
	CacheDSN    string   // CacheDSN Redis data source name.
	Group       string   // Group kafka's group.
	Host        string   // Host server's address.
	TopicIn     string   // TopicIn kafka's incoming message's topic.
	TopicError  string   // TopicIn kafka's out coming errors message's topic.
}
    Config server's settings.

func Parse() Config
    Parse main func to parse variables.

