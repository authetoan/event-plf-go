package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type consumerImpl struct {
	reader *kafka.Reader
}

func NewConsumer(broker, topic, groupID string) Consumer {
	return &consumerImpl{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{broker},
			Topic:   topic,
			GroupID: groupID,
		}),
	}
}

func (c *consumerImpl) Consume(handler func(message []byte) error) error {
	defer c.reader.Close()

	for {
		msg, err := c.reader.ReadMessage(context.Background())
		if err != nil {
			return fmt.Errorf("failed to read message: %w", err)
		}

		if err := handler(msg.Value); err != nil {
			return fmt.Errorf("failed to handle message: %w", err)
		}
	}
}
