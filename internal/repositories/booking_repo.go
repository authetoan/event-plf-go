package repositories

import (
	"booking-service/internal/models"
	"time"

	"gorm.io/gorm"
)

type BookingRepository interface {
	CreateBooking(booking *models.Booking) error
	GetBookingByID(bookingID uint) (*models.Booking, error)
	UpdateBookingStatus(bookingID uint, status models.BookingStatus) error
	DeleteBooking(bookingID uint) error
	ListBookingsByUserID(userID uint, page, pageSize int) ([]models.Booking, error)
	GetPendingBookingsOlderThan(duration time.Duration) ([]models.Booking, error)
}

type bookingRepositoryImpl struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepositoryImpl{
		db: db,
	}
}

func (r *bookingRepositoryImpl) CreateBooking(booking *models.Booking) error {
	if err := r.db.Create(booking).Error; err != nil {
		return err
	}
	return nil
}

func (r *bookingRepositoryImpl) GetBookingByID(bookingID uint) (*models.Booking, error) {
	var booking models.Booking
	if err := r.db.Preload("Tickets").First(&booking, "id = ?", bookingID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepositoryImpl) UpdateBookingStatus(bookingID uint, status models.BookingStatus) error {
	if err := r.db.Model(&models.Booking{}).Where("id = ?", bookingID).Update("status", status).Error; err != nil {
		return err
	}
	return nil
}

func (r *bookingRepositoryImpl) DeleteBooking(bookingID uint) error {
	if err := r.db.Delete(&models.Booking{}, "id = ?", bookingID).Error; err != nil {
		return err
	}
	return nil
}

func (r *bookingRepositoryImpl) ListBookingsByUserID(userID uint, page, pageSize int) ([]models.Booking, error) {
	var bookings []models.Booking
	offset := (page - 1) * pageSize
	if err := r.db.Where("user_id = ?", userID).Offset(offset).Limit(pageSize).Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}

func (r *bookingRepositoryImpl) GetPendingBookingsOlderThan(duration time.Duration) ([]models.Booking, error) {
	var bookings []models.Booking
	cutoffTime := time.Now().Add(-duration)

	if err := r.db.Where("status = ? AND created_at < ?", models.BookingStatusPending, cutoffTime).Find(&bookings).Error; err != nil {
		return nil, err
	}
	return bookings, nil
}
