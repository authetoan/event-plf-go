package controllers_test

import (
	"booking-service/internal/controllers"
	"booking-service/internal/models"
	"booking-service/test/mocks"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(controller controllers.BookingController) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	apiRoutes := router.Group("/api")
	controllers.RegisterBookingRoutes(apiRoutes, controller)
	return router
}

func TestCreateBooking_Success(t *testing.T) {
	mockBookingService := new(mocks.BookingServiceMock)
	bookingController := controllers.NewBookingController(mockBookingService)
	router := setupRouter(bookingController)

	bookingRequest := map[string]interface{}{
		"user_id":    uint(1),
		"event_id":   uint(1),
		"ticket_ids": []uint{1, 2},
	}

	mockBooking := &models.Booking{
		ID:          123,
		UserID:      1,
		EventID:     1,
		TotalAmount: 200.0,
		Status:      models.BookingStatusPending,
	}

	mockBookingService.On("CreateBooking", uint(1), uint(1), []uint{1, 2}).Return(mockBooking, nil)

	requestBody, _ := json.Marshal(bookingRequest)
	req := httptest.NewRequest(http.MethodPost, "/api/bookings", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var response models.Booking
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, mockBooking.ID, response.ID)
	assert.Equal(t, mockBooking.UserID, response.UserID)
	assert.Equal(t, mockBooking.EventID, response.EventID)

	mockBookingService.AssertExpectations(t)
}

func TestCreateBooking_ValidationError(t *testing.T) {
	mockBookingService := new(mocks.BookingServiceMock)
	bookingController := controllers.NewBookingController(mockBookingService)
	router := setupRouter(bookingController)

	invalidRequest := map[string]interface{}{
		"event_id":   uint(1),
		"ticket_ids": []uint{1, 2},
	}

	requestBody, _ := json.Marshal(invalidRequest)
	req := httptest.NewRequest(http.MethodPost, "/api/bookings", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Invalid request payload", response["error"])
}

func TestGetBookingByID_Success(t *testing.T) {
	mockBookingService := new(mocks.BookingServiceMock)
	bookingController := controllers.NewBookingController(mockBookingService)
	router := setupRouter(bookingController)

	mockBooking := &models.Booking{
		ID:          123,
		UserID:      1,
		EventID:     1,
		TotalAmount: 200.0,
		Status:      models.BookingStatusConfirmed,
	}

	mockBookingService.On("GetBookingByID", uint(123)).Return(mockBooking, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/bookings/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response models.Booking
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, mockBooking.ID, response.ID)
	assert.Equal(t, mockBooking.UserID, response.UserID)
	assert.Equal(t, mockBooking.EventID, response.EventID)

	mockBookingService.AssertExpectations(t)
}

func TestGetBookingByID_NotFound(t *testing.T) {
	mockBookingService := new(mocks.BookingServiceMock)
	bookingController := controllers.NewBookingController(mockBookingService)
	router := setupRouter(bookingController)

	mockBookingService.On("GetBookingByID", uint(123)).Return(nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/bookings/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Booking not found", response["error"])
}

func TestCancelBooking_Success(t *testing.T) {
	mockBookingService := new(mocks.BookingServiceMock)
	bookingController := controllers.NewBookingController(mockBookingService)
	router := setupRouter(bookingController)

	mockBookingService.On("CancelBooking", uint(123)).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/bookings/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Booking cancelled successfully", response["message"])
}

func TestCancelBooking_Error(t *testing.T) {
	mockBookingService := new(mocks.BookingServiceMock)
	bookingController := controllers.NewBookingController(mockBookingService)
	router := setupRouter(bookingController)

	mockBookingService.On("CancelBooking", uint(123)).Return(assert.AnError)

	req := httptest.NewRequest(http.MethodDelete, "/api/bookings/123", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Failed to cancel booking", response["error"])
}
