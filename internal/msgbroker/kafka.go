package msgbroker

import (
	"context"
	"fmt"
	"time"

	"github.com/erupshis/effective_mobile/internal/logger"
	"github.com/segmentio/kafka-go"
)

// KafkaProducer WRITER.
type KafkaProducer struct {
	kafka.Writer
	log logger.BaseLogger
}

// CreateKafkaProducer Create writer.
func CreateKafkaProducer(brokerAddr []string, topic string, log logger.BaseLogger) Producer {
	producer := &KafkaProducer{
		Writer: kafka.Writer{
			Addr:                   kafka.TCP(brokerAddr...),
			Topic:                  topic,
			Balancer:               &kafka.LeastBytes{},
			AllowAutoTopicCreation: true,
			Logger:                 log,
		},
		log: log,
	}
	return producer
}

func (p *KafkaProducer) Listen(ctx context.Context, chMessages <-chan Message) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-chMessages:
			if !ok {
				//channel was closed.
				return
			}

			err := p.SendMessage(ctx, string(msg.Key), string(msg.Value))
			if err != nil {
				p.log.Info("[KafkaProducer:Listen] send message '%v' finished with error: %v.", msg, err)
				time.Sleep(time.Second)
			}
			p.log.Info("[KafkaProducer:Listen] message sent: %s = %s\n", string(msg.Key), string(msg.Value))
		}
	}
}

func (p *KafkaProducer) SendMessage(ctx context.Context, key, value string) error {
	message := kafka.Message{
		Key:   []byte(key),
		Value: []byte(value),
	}

	if err := p.WriteMessages(ctx, message); err != nil {
		return fmt.Errorf("failed to send kafka message: %w", err)
	}

	return nil
}

func (p *KafkaProducer) Close() error {
	return p.Writer.Close()
}

type KafkaConsumer struct {
	*kafka.Reader
	log logger.BaseLogger
}

// CreateKafkaConsumer Create reader.
func CreateKafkaConsumer(brokerAddr []string, topic string, groupID string, log logger.BaseLogger) Consumer {
	readerConfig := kafka.ReaderConfig{
		Brokers:     brokerAddr,
		Topic:       topic,
		GroupID:     groupID,
		Partition:   0,
		MaxBytes:    10e6, // 10MB
		StartOffset: kafka.FirstOffset,
	}

	reader := kafka.NewReader(readerConfig)
	return &KafkaConsumer{Reader: reader, log: log}
}

func (c *KafkaConsumer) Listen(ctx context.Context, chMessages chan<- Message) {
	for {
		select {
		case <-ctx.Done():
			close(chMessages)
			return
		default:
			msg, err := c.ReadMessage(ctx)
			if err != nil {
				c.log.Info("[KafkaConsumer:Listen] read message finished with error: %v. Sleep to retry.", err)
				time.Sleep(time.Second)
			}
			c.log.Info("[KafkaConsumer:Listen] message received: %s = %s\n", string(msg.Key), string(msg.Value))
			chMessages <- Message{Key: msg.Key, Value: msg.Value}
		}
	}
}

func (c *KafkaConsumer) ReadMessage(ctx context.Context) (Message, error) {
	rawMessage, err := c.Reader.ReadMessage(ctx)
	if err != nil {
		return Message{}, fmt.Errorf("failed to read kafka message: %w", err)
	}
	message := Message{
		Key:   rawMessage.Key,
		Value: rawMessage.Value,
	}
	return message, nil
}

func (c *KafkaConsumer) Close() error {
	return c.Reader.Close()
}
