package controllers

import (
	"booking-service/internal/services"
	"booking-service/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TicketController interface {
	CreateTicketsForEvent(c *gin.Context)
}

type ticketControllerImpl struct {
	TicketService services.TicketService
	Logger        *utils.Logger
}

func NewTicketController(ticketService services.TicketService) TicketController {
	return &ticketControllerImpl{
		TicketService: ticketService,
		Logger:        utils.NewLogger(),
	}
}

func (tc *ticketControllerImpl) CreateTicketsForEvent(c *gin.Context) {
	var request struct {
		NumTickets int     `json:"num_tickets" binding:"required"`
		Price      float64 `json:"price" binding:"required"`
	}

	eventIDStr := c.Param("event_id")
	eventID, err := strconv.ParseUint(eventIDStr, 10, 64)
	if err != nil {
		tc.Logger.Warn("Invalid event ID: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID", "details": err.Error()})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		tc.Logger.Warn("Invalid request payload: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload", "details": err.Error()})
		return
	}

	tickets, err := tc.TicketService.CreateTicketsForEvent(uint(eventID), request.NumTickets, request.Price)
	if err != nil {
		tc.Logger.Error("Failed to create tickets: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create tickets", "details": err.Error()})
		return
	}

	tc.Logger.Info("Tickets created successfully for event " + eventIDStr)
	c.JSON(http.StatusCreated, tickets)
}

func RegisterTicketRoutes(router *gin.RouterGroup, controller TicketController) {
	ticketRoutes := router.Group("/tickets")
	{
		ticketRoutes.POST("/:event_id", controller.CreateTicketsForEvent)
	}
}
