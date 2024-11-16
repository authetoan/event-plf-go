package controllers

import (
	"booking-service/internal/models"
	"booking-service/internal/services"
	"booking-service/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookingController interface {
	CreateBooking(c *gin.Context)
	GetBookingByID(c *gin.Context)
	UpdateBookingStatus(c *gin.Context)
	CancelBooking(c *gin.Context)
	ListBookingsByUserID(c *gin.Context)
}

type bookingControllerImpl struct {
	BookingService services.BookingService
	Logger         *utils.Logger
}

func NewBookingController(bookingService services.BookingService) BookingController {
	return &bookingControllerImpl{
		BookingService: bookingService,
		Logger:         utils.NewLogger(),
	}
}

func (bc *bookingControllerImpl) CreateBooking(c *gin.Context) {
	var request struct {
		UserID    uint   `json:"user_id" binding:"required"`
		EventID   uint   `json:"event_id" binding:"required"`
		TicketIDs []uint `json:"ticket_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		bc.Logger.Warn("Invalid request payload: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	booking, err := bc.BookingService.CreateBooking(request.UserID, request.EventID, request.TicketIDs)
	if err != nil {
		bc.Logger.Error("Failed to create booking: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create booking", "details": err.Error()})
		return
	}

	bc.Logger.Info("Booking created successfully: " + strconv.Itoa(int(booking.ID)))
	c.JSON(http.StatusCreated, booking)
}

func (bc *bookingControllerImpl) GetBookingByID(c *gin.Context) {
	bookingIDStr := c.Param("id")
	bookingID, err := strconv.Atoi(bookingIDStr)
	if err != nil {
		bc.Logger.Warn("Invalid booking ID: " + bookingIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}
	booking, err := bc.BookingService.GetBookingByID(uint(bookingID))
	if err != nil {
		bc.Logger.Error("Failed to retrieve booking: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve booking", "details": err.Error()})
		return
	}

	if booking == nil {
		bc.Logger.Warn("Booking not found: " + bookingIDStr)
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking not found"})
		return
	}

	bc.Logger.Info("Booking retrieved successfully: " + bookingIDStr)
	c.JSON(http.StatusOK, booking)
}

func (bc *bookingControllerImpl) UpdateBookingStatus(c *gin.Context) {
	bookingIDStr := c.Param("id")
	bookingID, err := strconv.Atoi(bookingIDStr)
	if err != nil {
		bc.Logger.Warn("Invalid booking ID: " + bookingIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}
	var request struct {
		Status models.BookingStatus `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		bc.Logger.Warn("Invalid request payload: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	if err := bc.BookingService.UpdateBookingStatus(uint(bookingID), request.Status); err != nil {
		bc.Logger.Error("Failed to update booking status: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update booking status", "details": err.Error()})
		return
	}

	bc.Logger.Info("Booking status updated successfully: " + bookingIDStr)
	c.JSON(http.StatusOK, gin.H{"message": "Booking status updated successfully"})
}

func (bc *bookingControllerImpl) CancelBooking(c *gin.Context) {
	bookingIDStr := c.Param("id")
	bookingID, err := strconv.Atoi(bookingIDStr)
	if err != nil {
		bc.Logger.Warn("Invalid booking ID: " + bookingIDStr)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid booking ID"})
		return
	}
	if err := bc.BookingService.CancelBooking(uint(bookingID)); err != nil {
		bc.Logger.Error("Failed to cancel booking: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel booking", "details": err.Error()})
		return
	}

	bc.Logger.Info("Booking cancelled successfully: " + bookingIDStr)
	c.JSON(http.StatusOK, gin.H{"message": "Booking cancelled successfully"})
}

func (bc *bookingControllerImpl) ListBookingsByUserID(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	page := utils.ParseQueryParamAsInt(c, "page", 1)
	pageSize := utils.ParseQueryParamAsInt(c, "page_size", 10)

	bookings, err := bc.BookingService.ListBookingsByUserID(uint(userID), page, pageSize)
	if err != nil {
		bc.Logger.Error("Failed to list bookings for user: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list bookings", "details": err.Error()})
		return
	}

	bc.Logger.Info("Bookings retrieved successfully for user: " + userIDStr)
	c.JSON(http.StatusOK, bookings)
}

func RegisterBookingRoutes(router *gin.RouterGroup, controller BookingController) {
	bookingRoutes := router.Group("/bookings")
	{
		bookingRoutes.POST("", controller.CreateBooking)
		bookingRoutes.GET("/:id", controller.GetBookingByID)
		bookingRoutes.PUT("/:id/status", controller.UpdateBookingStatus)
		bookingRoutes.DELETE("/:id", controller.CancelBooking)
		bookingRoutes.GET("/user/:user_id", controller.ListBookingsByUserID)
	}
}
