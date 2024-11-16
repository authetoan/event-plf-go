package mocks

import (
	"booking-service/internal/models"
	kafkaModels "booking-service/pkg/kafka/models"

	"github.com/stretchr/testify/mock"
)

type TicketServiceMock struct {
	mock.Mock
}

func (m *TicketServiceMock) CreateTicketsForEvent(eventID uint, numTickets int, price float64) ([]models.Ticket, error) {
	args := m.Called(eventID, numTickets, price)
	return args.Get(0).([]models.Ticket), args.Error(1)
}

func (m *TicketServiceMock) ReserveTicket(ticketID, userID, bookingID uint) error {
	args := m.Called(ticketID, userID, bookingID)
	return args.Error(0)
}

func (m *TicketServiceMock) ListAvailableTickets(eventID uint) ([]models.Ticket, error) {
	args := m.Called(eventID)
	return args.Get(0).([]models.Ticket), args.Error(1)
}

func (m *TicketServiceMock) ReleaseTicket(ticketID uint) error {
	args := m.Called(ticketID)
	return args.Error(0)
}

func (m *TicketServiceMock) MarkTicketAsSold(ticketID uint) error {
	args := m.Called(ticketID)
	return args.Error(0)
}

func (m *TicketServiceMock) HandleBookingEvent(eventType string, payload kafkaModels.BookingEvent) error {
	args := m.Called(eventType, payload)
	return args.Error(0)
}

func (m *TicketServiceMock) GetTicketByID(ticketID uint) (*models.Ticket, error) {
	args := m.Called(ticketID)
	return args.Get(0).(*models.Ticket), args.Error(1)
}

func (m *TicketServiceMock) UpdateTicketStatus(ticketID uint, status models.TicketStatus) error {
	args := m.Called(ticketID, status)
	return args.Error(0)
}
