package mocks

import (
	"booking-service/internal/models"
	"github.com/stretchr/testify/mock"
)

type TicketRepositoryMock struct {
	mock.Mock
}

func (m *TicketRepositoryMock) CreateTicket(ticket *models.Ticket) error {
	panic("implement me")
}

func (m *TicketRepositoryMock) CreateTicketsBatch(tickets []models.Ticket) error {
	args := m.Called(tickets)
	return args.Error(0)
}

func (m *TicketRepositoryMock) GetTicketByID(ticketID uint) (*models.Ticket, error) {
	args := m.Called(ticketID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Ticket), args.Error(1)
}

func (m *TicketRepositoryMock) ReserveTicket(ticketID, userID, bookingID uint) error {
	args := m.Called(ticketID, userID, bookingID)
	return args.Error(0)
}

func (m *TicketRepositoryMock) UpdateTicketStatus(ticketID uint, status models.TicketStatus) error {
	args := m.Called(ticketID, status)
	return args.Error(0)
}

func (m *TicketRepositoryMock) ListAvailableTickets(eventID uint) ([]models.Ticket, error) {
	args := m.Called(eventID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Ticket), args.Error(1)
}

func (m *TicketRepositoryMock) DeleteTicket(ticketID uint) error {
	args := m.Called(ticketID)
	return args.Error(0)
}
