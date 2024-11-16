package mocks

import (
	"booking-service/internal/models"
	"github.com/stretchr/testify/mock"
	"time"
)

type BookingRepositoryMock struct {
	mock.Mock
}

func (m *BookingRepositoryMock) GetPendingBookingsOlderThan(duration time.Duration) ([]models.Booking, error) {
	panic("implement me")
}

func (m *BookingRepositoryMock) CreateBooking(booking *models.Booking) error {
	args := m.Called(booking)
	return args.Error(0)
}

func (m *BookingRepositoryMock) GetBookingByID(bookingID uint) (*models.Booking, error) {
	args := m.Called(bookingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Booking), args.Error(1)
}

func (m *BookingRepositoryMock) UpdateBookingStatus(bookingID uint, status models.BookingStatus) error {
	args := m.Called(bookingID, status)
	return args.Error(0)
}

func (m *BookingRepositoryMock) DeleteBooking(bookingID uint) error {
	args := m.Called(bookingID)
	return args.Error(0)
}

func (m *BookingRepositoryMock) ListBookingsByUserID(userID uint, page, pageSize int) ([]models.Booking, error) {
	args := m.Called(userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]models.Booking), args.Error(1)
}
