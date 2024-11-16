package kafka

import (
	"booking-service/pkg/kafka/models"
	"encoding/json"
	"fmt"
)

func ParseBookingEvent(message []byte) models.BookingEvent {
	var event models.BookingEvent
	if err := json.Unmarshal(message, &event); err != nil {
		fmt.Printf("Failed to parse Kafka message: %v\n", err)
	}
	return event
}
