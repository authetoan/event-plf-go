package services_test

import (
	"booking-service/internal/models"
	"booking-service/internal/services"
	"booking-service/test/mocks"
	"booking-service/utils"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupMocks() (*mocks.BookingRepositoryMock, *mocks.TicketServiceMock, *mocks.KafkaProducerMock, services.BookingService) {
	bookingRepoMock := new(mocks.BookingRepositoryMock)
	ticketServiceMock := new(mocks.TicketServiceMock)
	kafkaProducerMock := new(mocks.KafkaProducerMock)
	bookingService := services.NewBookingService(bookingRepoMock, ticketServiceMock, kafkaProducerMock)
	return bookingRepoMock, ticketServiceMock, kafkaProducerMock, bookingService
}

func TestCreateBooking_Success(t *testing.T) {
	bookingRepoMock, ticketServiceMock, kafkaProducerMock, bookingService := setupMocks()

	userID := uint(1)
	eventID := uint(1)
	ticketIDs := []uint{1, 2}
	mockTickets := []models.Ticket{
		{ID: 1, Price: 100, Status: models.TicketStatusAvailable},
		{ID: 2, Price: 150, Status: models.TicketStatusAvailable},
	}
	mockBooking := &models.Booking{
		ID:          1,
		UserID:      userID,
		EventID:     eventID,
		TotalAmount: 250,
		Status:      models.BookingStatusPending,
		Tickets:     mockTickets,
	}

	for _, ticket := range mockTickets {
		ticketServiceMock.On("GetTicketByID", ticket.ID).Return(&ticket, nil)
		ticketServiceMock.On("ReserveTicket", ticket.ID, userID, mock.AnythingOfType("uint")).Return(nil)
	}

	bookingRepoMock.On("CreateBooking", mock.Anything).Return(nil).Run(func(args mock.Arguments) {
		booking := args.Get(0).(*models.Booking)
		booking.ID = 1
	})
	kafkaProducerMock.On("Publish", mock.Anything, mock.Anything).Return(nil)

	result, err := bookingService.CreateBooking(userID, eventID, ticketIDs)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockBooking.UserID, result.UserID)
	assert.Equal(t, mockBooking.EventID, result.EventID)
	assert.Equal(t, mockBooking.TotalAmount, result.TotalAmount)

	bookingRepoMock.AssertExpectations(t)
	ticketServiceMock.AssertExpectations(t)
	kafkaProducerMock.AssertExpectations(t)
}

func TestCreateBooking_TicketNotAvailable(t *testing.T) {
	_, ticketServiceMock, _, bookingService := setupMocks()

	userID := uint(1)
	eventID := uint(1)
	ticketIDs := []uint{1}
	ticketServiceMock.On("GetTicketByID", uint(1)).Return(&models.Ticket{
		ID: 1, Status: models.TicketStatusSold,
	}, nil)

	result, err := bookingService.CreateBooking(userID, eventID, ticketIDs)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "Ticket not available")

	ticketServiceMock.AssertExpectations(t)
}

func TestConfirmBooking_Success(t *testing.T) {
	bookingRepoMock := new(mocks.BookingRepositoryMock)
	ticketServiceMock := new(mocks.TicketServiceMock)
	kafkaProducerMock := new(mocks.KafkaProducerMock)
	bookingService := services.NewBookingService(bookingRepoMock, ticketServiceMock, kafkaProducerMock)

	bookingID := uint(1)
	mockBooking := &models.Booking{
		ID: bookingID,
		Tickets: []models.Ticket{
			{ID: 1},
			{ID: 2},
		},
	}

	// Mock repository and service calls
	bookingRepoMock.On("GetBookingByID", bookingID).Return(mockBooking, nil)
	bookingRepoMock.On("UpdateBookingStatus", bookingID, models.BookingStatusConfirmed).Return(nil)

	// Mock ticket service calls for UpdateTicketStatus
	ticketServiceMock.On("UpdateTicketStatus", uint(1), models.TicketStatusSold).Return(nil)
	ticketServiceMock.On("UpdateTicketStatus", uint(2), models.TicketStatusSold).Return(nil)

	// Mock Kafka publish call
	kafkaProducerMock.On("Publish", mock.Anything, mock.Anything).Return(nil)

	// Run the test
	err := bookingService.ConfirmBooking(bookingID)

	// Assertions
	assert.NoError(t, err)
	bookingRepoMock.AssertExpectations(t)
	ticketServiceMock.AssertExpectations(t)
	kafkaProducerMock.AssertExpectations(t)
}

func TestCancelBooking_Success(t *testing.T) {
	bookingRepoMock := new(mocks.BookingRepositoryMock)
	ticketServiceMock := new(mocks.TicketServiceMock)
	kafkaProducerMock := new(mocks.KafkaProducerMock)
	bookingService := services.NewBookingService(bookingRepoMock, ticketServiceMock, kafkaProducerMock)

	bookingID := uint(1)
	mockBooking := &models.Booking{
		ID: bookingID,
		Tickets: []models.Ticket{
			{ID: 1},
			{ID: 2},
		},
	}

	// Mock repository and service calls
	bookingRepoMock.On("GetBookingByID", bookingID).Return(mockBooking, nil)
	bookingRepoMock.On("UpdateBookingStatus", bookingID, models.BookingStatusCanceled).Return(nil)

	// Mock ticket service calls for UpdateTicketStatus
	ticketServiceMock.On("UpdateTicketStatus", uint(1), models.TicketStatusAvailable).Return(nil)
	ticketServiceMock.On("UpdateTicketStatus", uint(2), models.TicketStatusAvailable).Return(nil)

	// Mock Kafka publish call
	kafkaProducerMock.On("Publish", mock.Anything, mock.Anything).Return(nil)

	// Run the test
	err := bookingService.CancelBooking(bookingID)

	// Assertions
	assert.NoError(t, err)
	bookingRepoMock.AssertExpectations(t)
	ticketServiceMock.AssertExpectations(t)
	kafkaProducerMock.AssertExpectations(t)
}

func TestGetBookingByID_Success(t *testing.T) {
	bookingRepoMock, _, _, bookingService := setupMocks()

	bookingID := uint(1)
	mockBooking := &models.Booking{
		ID: bookingID,
	}

	bookingRepoMock.On("GetBookingByID", bookingID).Return(mockBooking, nil)

	result, err := bookingService.GetBookingByID(bookingID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, mockBooking.ID, result.ID)

	bookingRepoMock.AssertExpectations(t)
}

func TestGetBookingByID_NotFound(t *testing.T) {
	bookingRepoMock, _, _, bookingService := setupMocks()

	bookingID := uint(1)
	bookingRepoMock.On("GetBookingByID", bookingID).Return(nil, nil)

	result, err := bookingService.GetBookingByID(bookingID)

	assert.Error(t, err)
	assert.Nil(t, result)

	appErr, ok := err.(*utils.AppError)
	assert.True(t, ok)
	assert.Equal(t, 404, appErr.Code)
	assert.Equal(t, "Booking not found", appErr.Message)

	bookingRepoMock.AssertExpectations(t)
}
