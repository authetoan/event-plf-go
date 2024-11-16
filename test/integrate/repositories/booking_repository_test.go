package repositories_test

import (
	"booking-service/internal/models"
	"booking-service/internal/repositories"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to in-memory database: %v", err)
	}
	// Migrate schema
	err = db.AutoMigrate(&models.Booking{}, &models.Ticket{}, &models.Event{})
	if err != nil {
		t.Fatalf("Failed to migrate schema: %v", err)
	}
	return db
}

func tearDownTestDB(db *gorm.DB, t *testing.T) {
	err := db.Exec("DELETE FROM bookings").Error
	if err != nil {
		t.Fatalf("Failed to clear database: %v", err)
	}
}

func TestCreateBooking(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.NewBookingRepository(db)

	booking := &models.Booking{
		UserID:      1,
		EventID:     1,
		TotalAmount: 100.0,
		Status:      models.BookingStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := repo.CreateBooking(booking)
	assert.NoError(t, err)
	assert.NotZero(t, booking.ID) // Ensure the ID is auto-generated

	var storedBooking models.Booking
	err = db.First(&storedBooking, "id = ?", booking.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, uint(1), storedBooking.UserID)
	assert.Equal(t, uint(1), storedBooking.EventID)
}

func TestGetBookingByID(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.NewBookingRepository(db)

	booking := &models.Booking{
		UserID:      1,
		EventID:     1,
		TotalAmount: 100.0,
		Status:      models.BookingStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := db.Create(booking).Error
	assert.NoError(t, err)

	result, err := repo.GetBookingByID(booking.ID)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, booking.ID, result.ID)
	t.Cleanup(func() { tearDownTestDB(db, t) })
}

func TestUpdateBookingStatus(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.NewBookingRepository(db)

	booking := &models.Booking{
		UserID:      1,
		EventID:     1,
		TotalAmount: 100.0,
		Status:      models.BookingStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := db.Create(booking).Error
	assert.NoError(t, err)

	err = repo.UpdateBookingStatus(booking.ID, models.BookingStatusConfirmed)
	assert.NoError(t, err)

	var updatedBooking models.Booking
	err = db.First(&updatedBooking, "id = ?", booking.ID).Error
	assert.NoError(t, err)
	assert.Equal(t, models.BookingStatusConfirmed, updatedBooking.Status)
	t.Cleanup(func() { tearDownTestDB(db, t) })
}

func TestDeleteBooking(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.NewBookingRepository(db)

	booking := &models.Booking{
		UserID:      1,
		EventID:     1,
		TotalAmount: 100.0,
		Status:      models.BookingStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := db.Create(booking).Error
	assert.NoError(t, err)

	err = repo.DeleteBooking(booking.ID)
	assert.NoError(t, err)

	var deletedBooking models.Booking
	err = db.First(&deletedBooking, "id = ?", booking.ID).Error
	assert.Error(t, err)
	assert.Equal(t, gorm.ErrRecordNotFound, err)
	t.Cleanup(func() { tearDownTestDB(db, t) })
}

func TestGetPendingBookingsOlderThan(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.NewBookingRepository(db)

	oldBooking := &models.Booking{
		UserID:      1,
		EventID:     1,
		TotalAmount: 100.0,
		Status:      models.BookingStatusPending,
		CreatedAt:   time.Now().Add(-2 * time.Hour),
		UpdatedAt:   time.Now().Add(-2 * time.Hour),
	}

	newBooking := &models.Booking{
		UserID:      2,
		EventID:     2,
		TotalAmount: 200.0,
		Status:      models.BookingStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := db.Create(oldBooking).Error
	assert.NoError(t, err)

	err = db.Create(newBooking).Error
	assert.NoError(t, err)

	result, err := repo.GetPendingBookingsOlderThan(1 * time.Hour)
	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, oldBooking.ID, result[0].ID)
}

func TestListBookingsByUserID(t *testing.T) {
	db := setupTestDB(t)
	repo := repositories.NewBookingRepository(db)

	booking1 := &models.Booking{
		UserID:      1,
		EventID:     1,
		TotalAmount: 100.0,
		Status:      models.BookingStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	booking2 := &models.Booking{
		UserID:      1,
		EventID:     2,
		TotalAmount: 200.0,
		Status:      models.BookingStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := db.Create(booking1).Error
	assert.NoError(t, err)

	err = db.Create(booking2).Error
	assert.NoError(t, err)

	result, err := repo.ListBookingsByUserID(1, 1, 10)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, uint(1), result[0].EventID)
	assert.Equal(t, uint(2), result[1].EventID)
	t.Cleanup(func() { tearDownTestDB(db, t) })
}
