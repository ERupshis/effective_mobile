package msgbroker // import "github.com/erupshis/effective_mobile/internal/msgbroker"

type Consumer interface{ ... }
    func CreateKafkaConsumer(brokerAddr []string, topic string, groupID string, log logger.BaseLogger) Consumer
type KafkaConsumer struct{ ... }
type KafkaProducer struct{ ... }
type Message struct{ ... }
type Producer interface{ ... }
    func CreateKafkaProducer(brokerAddr []string, topic string, log logger.BaseLogger) Producer
