package kafka

type Producer interface {
	Publish(topic string, message interface{}) error
}

type Consumer interface {
	Consume(handler func(message []byte) error) error
}
