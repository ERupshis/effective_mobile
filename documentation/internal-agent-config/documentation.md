package config // import "github.com/erupshis/effective_mobile/internal/agent/config"

Package config agent's setting parser. Applies flags and environments.
Environments are prioritized.

TYPES

type Config struct {
	BrokerAddr []string
	TopicIn    string
	TopicOut   string
	Group      string
}
    Config agent's settings.

func Parse() Config
    Parse main func to parse variables.

