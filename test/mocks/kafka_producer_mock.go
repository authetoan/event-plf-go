package mocks

import (
	"github.com/stretchr/testify/mock"
)

type KafkaProducerMock struct {
	mock.Mock
}

func (m *KafkaProducerMock) Publish(topic string, message interface{}) error {
	args := m.Called(topic, message)
	return args.Error(0)
}
