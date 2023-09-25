package msgbroker // import "github.com/erupshis/effective_mobile/internal/msgbroker"


TYPES

type Consumer interface {
	Listen(ctx context.Context, chMessages chan<- Message)
	ReadMessage(ctx context.Context) (Message, error)
	Close() error
}

func CreateKafkaConsumer(brokerAddr []string, topic string, groupID string, log logger.BaseLogger) Consumer
    CreateKafkaConsumer Create reader.

type KafkaConsumer struct {
	*kafka.Reader
	// Has unexported fields.
}

func (c *KafkaConsumer) Close() error

func (c *KafkaConsumer) Listen(ctx context.Context, chMessages chan<- Message)

func (c *KafkaConsumer) ReadMessage(ctx context.Context) (Message, error)

type KafkaProducer struct {
	kafka.Writer
	// Has unexported fields.
}
    KafkaProducer WRITER.

func (p *KafkaProducer) Close() error

func (p *KafkaProducer) Listen(ctx context.Context, chMessages <-chan Message)

func (p *KafkaProducer) SendMessage(ctx context.Context, key, value string) error

type Message struct {
	Key   []byte
	Value []byte
}

type Producer interface {
	Listen(ctx context.Context, chMessages <-chan Message)
	SendMessage(ctx context.Context, key, value string) error
	Close() error
}

func CreateKafkaProducer(brokerAddr []string, topic string, log logger.BaseLogger) Producer
    CreateKafkaProducer Create writer.

