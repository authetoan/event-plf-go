package kafka

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

type producerImpl struct {
	writer *kafka.Writer
}

func NewProducer(broker string) Producer {
	return &producerImpl{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(broker),
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *producerImpl) Publish(topic string, message interface{}) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to serialize message: %w", err)
	}

	err = p.writer.WriteMessages(
		context.Background(),
		kafka.Message{
			Topic: topic,
			Value: payload,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to publish message to topic %s: %w", topic, err)
	}

	return nil
}
