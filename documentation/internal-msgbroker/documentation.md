package msgbroker // import "github.com/erupshis/effective_mobile/internal/msgbroker"

Package msgbroker implements message transactions.

TYPES

type Consumer interface {
	// Listen goroutine method for listening to send response messages.
	Listen(ctx context.Context, chMessages chan<- Message)

	// ReadMessage method for reading incoming messages.
	ReadMessage(ctx context.Context) (Message, error)

	// Close closes the stream, preventing the program from reading any more messages from it.
	Close() error
}
    Consumer interface of messages reader.

func CreateKafkaConsumer(brokerAddr []string, topic string, groupID string, log logger.BaseLogger) Consumer
    CreateKafkaConsumer Create reader.

type KafkaConsumer struct {
	*kafka.Reader
	// Has unexported fields.
}
    KafkaConsumer Reader wrapper.

func (c *KafkaConsumer) Close() error
    Close closes the stream, preventing the program from reading any more
    messages from it.

func (c *KafkaConsumer) Listen(ctx context.Context, chMessages chan<- Message)
    Listen goroutine method for listening to send response messages.

func (c *KafkaConsumer) ReadMessage(ctx context.Context) (Message, error)
    ReadMessage method for reading incoming messages.

type KafkaProducer struct {
	kafka.Writer
	// Has unexported fields.
}
    KafkaProducer writer wrapper.

func (p *KafkaProducer) Close() error
    Close flushes pending sends and waits for all writes to complete before
    returning.

func (p *KafkaProducer) Listen(ctx context.Context, chMessages <-chan Message)
    Listen goroutine method for listening response messages.

func (p *KafkaProducer) SendMessage(ctx context.Context, key, value string) error
    SendMessage method for sending messages.

type Message struct {
	Key   []byte
	Value []byte
}
    Message structure of message.

type Producer interface {
	// Listen goroutine method for listening response messages.
	Listen(ctx context.Context, chMessages <-chan Message)

	// SendMessage method for sending messages.
	SendMessage(ctx context.Context, key, value string) error

	// Close flushes pending sends and waits for all writes to complete before returning.
	Close() error
}
    Producer interface of messages sender.

func CreateKafkaProducer(brokerAddr []string, topic string, log logger.BaseLogger) Producer
    CreateKafkaProducer Create writer.

