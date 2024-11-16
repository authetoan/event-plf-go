package services_test

import (
	"booking-service/internal/models"
	"booking-service/internal/services"
	kafkaModels "booking-service/pkg/kafka/models"
	"booking-service/test/mocks"
	"booking-service/utils"
	"errors"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateTicketsForEvent_Success(t *testing.T) {
	ticketRepoMock := new(mocks.TicketRepositoryMock)
	ticketService := services.NewTicketService(ticketRepoMock)

	eventID := uint(1)
	numTickets := 10
	price := 50.0

	tickets := make([]models.Ticket, numTickets)
	for i := 0; i < numTickets; i++ {
		tickets[i] = models.Ticket{
			EventID: eventID,
			Price:   price,
			Status:  models.TicketStatusAvailable,
		}
	}

	ticketRepoMock.On("CreateTicketsBatch", mock.Anything).Return(nil)

	result, err := ticketService.CreateTicketsForEvent(eventID, numTickets, price)

	assert.NoError(t, err)
	assert.Len(t, result, numTickets)
	ticketRepoMock.AssertExpectations(t)
}

func TestCreateTicketsForEvent_InvalidInput(t *testing.T) {
	ticketRepoMock := new(mocks.TicketRepositoryMock)
	ticketService := services.NewTicketService(ticketRepoMock)

	_, err := ticketService.CreateTicketsForEvent(1, -1, 50.0)

	assert.Error(t, err)
	var appErr *utils.AppError
	ok := errors.As(err, &appErr)
	assert.True(t, ok)
	assert.Equal(t, 400, appErr.Code)
	ticketRepoMock.AssertExpectations(t)
}

func TestReserveTicket_Success(t *testing.T) {
	ticketRepoMock := new(mocks.TicketRepositoryMock)
	ticketService := services.NewTicketService(ticketRepoMock)

	ticketID := uint(1)
	userID := uint(1)
	bookingID := uint(1)
	ticket := &models.Ticket{
		ID:     ticketID,
		Status: models.TicketStatusAvailable,
	}

	ticketRepoMock.On("GetTicketByID", ticketID).Return(ticket, nil)
	ticketRepoMock.On("ReserveTicket", ticketID, userID, bookingID).Return(nil)

	err := ticketService.ReserveTicket(ticketID, userID, bookingID)

	assert.NoError(t, err)
	ticketRepoMock.AssertExpectations(t)
}

func TestReserveTicket_NotAvailable(t *testing.T) {
	ticketRepoMock := new(mocks.TicketRepositoryMock)
	ticketService := services.NewTicketService(ticketRepoMock)

	ticketID := uint(1)
	userID := uint(1)
	bookingID := uint(1)
	ticket := &models.Ticket{
		ID:     ticketID,
		Status: models.TicketStatusSold,
	}

	ticketRepoMock.On("GetTicketByID", ticketID).Return(ticket, nil)

	err := ticketService.ReserveTicket(ticketID, userID, bookingID)

	assert.Error(t, err)
	var appErr *utils.AppError
	ok := errors.As(err, &appErr)
	assert.True(t, ok)
	assert.Equal(t, 400, appErr.Code)
	ticketRepoMock.AssertExpectations(t)
}

func TestHandleBookingEvent_BookingCanceled(t *testing.T) {
	ticketRepoMock := new(mocks.TicketRepositoryMock)
	ticketService := services.NewTicketService(ticketRepoMock)

	eventType := "booking.canceled"
	payload := kafkaModels.BookingEvent{
		TicketIDs: []uint{1, 2},
	}

	// Mock GetTicketByID for each ticket
	ticketRepoMock.On("GetTicketByID", uint(1)).Return(&models.Ticket{
		ID:     1,
		Status: models.TicketStatusReserved,
	}, nil)
	ticketRepoMock.On("GetTicketByID", uint(2)).Return(&models.Ticket{
		ID:     2,
		Status: models.TicketStatusReserved,
	}, nil)

	// Mock UpdateTicketStatus for each ticket
	ticketRepoMock.On("UpdateTicketStatus", uint(1), models.TicketStatusAvailable).Return(nil)
	ticketRepoMock.On("UpdateTicketStatus", uint(2), models.TicketStatusAvailable).Return(nil)

	// Execute the test
	err := ticketService.HandleBookingEvent(eventType, payload)

	// Assertions
	assert.NoError(t, err)
	ticketRepoMock.AssertExpectations(t)
}

func TestHandleBookingEvent_BookingConfirmed(t *testing.T) {
	ticketRepoMock := new(mocks.TicketRepositoryMock)
	ticketService := services.NewTicketService(ticketRepoMock)

	eventType := "booking.confirmed"
	payload := kafkaModels.BookingEvent{
		TicketIDs: []uint{1, 2},
	}

	// Mock GetTicketByID for each ticket
	ticketRepoMock.On("GetTicketByID", uint(1)).Return(&models.Ticket{
		ID:     1,
		Status: models.TicketStatusReserved,
	}, nil)
	ticketRepoMock.On("GetTicketByID", uint(2)).Return(&models.Ticket{
		ID:     2,
		Status: models.TicketStatusReserved,
	}, nil)

	// Mock UpdateTicketStatus for each ticket
	ticketRepoMock.On("UpdateTicketStatus", uint(1), models.TicketStatusSold).Return(nil)
	ticketRepoMock.On("UpdateTicketStatus", uint(2), models.TicketStatusSold).Return(nil)

	// Execute the test
	err := ticketService.HandleBookingEvent(eventType, payload)

	// Assertions
	assert.NoError(t, err)
	ticketRepoMock.AssertExpectations(t)
}

func TestHandleBookingEvent_UnhandledEvent(t *testing.T) {
	ticketRepoMock := new(mocks.TicketRepositoryMock)
	ticketService := services.NewTicketService(ticketRepoMock)

	eventType := "unknown.event"
	payload := kafkaModels.BookingEvent{
		TicketIDs: []uint{1, 2},
	}

	err := ticketService.HandleBookingEvent(eventType, payload)

	assert.NoError(t, err)
	ticketRepoMock.AssertExpectations(t)
}
