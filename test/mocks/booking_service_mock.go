package mocks

import (
	"booking-service/internal/models"
	"github.com/stretchr/testify/mock"
)

type BookingServiceMock struct {
	mock.Mock
}

func (m *BookingServiceMock) CreateBooking(userID, eventID uint, ticketIDs []uint) (*models.Booking, error) {
	args := m.Called(userID, eventID, ticketIDs)
	return args.Get(0).(*models.Booking), args.Error(1)
}

func (m *BookingServiceMock) ConfirmBooking(bookingID uint) error {
	args := m.Called(bookingID)
	return args.Error(0)
}

func (m *BookingServiceMock) CancelBooking(bookingID uint) error {
	args := m.Called(bookingID)
	return args.Error(0)
}

func (m *BookingServiceMock) UpdateBookingStatus(bookingID uint, status models.BookingStatus) error {
	args := m.Called(bookingID, status)
	return args.Error(0)
}

func (m *BookingServiceMock) ListBookingsByUserID(userID uint, page, pageSize int) ([]models.Booking, error) {
	args := m.Called(userID, page, pageSize)
	return args.Get(0).([]models.Booking), args.Error(1)
}

func (m *BookingServiceMock) GetBookingByID(bookingID uint) (*models.Booking, error) {
	args := m.Called(bookingID)

	// Check if the first argument is nil
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Booking), args.Error(1)
}
