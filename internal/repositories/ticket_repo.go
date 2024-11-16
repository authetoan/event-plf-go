package repositories

import (
	"booking-service/internal/models"
	"booking-service/utils"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type TicketRepository interface {
	CreateTicket(ticket *models.Ticket) error
	CreateTicketsBatch(tickets []models.Ticket) error
	GetTicketByID(ticketID uint) (*models.Ticket, error)
	UpdateTicketStatus(ticketID uint, status models.TicketStatus) error
	ReserveTicket(ticketID, userID, bookingID uint) error
	ListAvailableTickets(eventID uint) ([]models.Ticket, error)
	DeleteTicket(ticketID uint) error
}

type ticketRepositoryImpl struct {
	db     *gorm.DB
	logger *utils.Logger
}

func NewTicketRepository(db *gorm.DB) TicketRepository {
	return &ticketRepositoryImpl{
		db:     db,
		logger: utils.NewLogger(),
	}
}

func (r *ticketRepositoryImpl) CreateTicket(ticket *models.Ticket) error {
	r.logger.Info(fmt.Sprintf("Creating ticket: %+v", ticket))
	if err := r.db.Create(ticket).Error; err != nil {
		appErr := utils.NewAppError(500, "Failed to create ticket", err.Error())
		r.logger.Error(appErr.Error())
		return appErr
	}
	r.logger.Info(fmt.Sprintf("Ticket created successfully: %+v", ticket))
	return nil
}

func (r *ticketRepositoryImpl) CreateTicketsBatch(tickets []models.Ticket) error {
	r.logger.Info(fmt.Sprintf("Creating batch of %d tickets", len(tickets)))
	if err := r.db.Create(&tickets).Error; err != nil {
		appErr := utils.NewAppError(500, "Failed to create tickets batch", err.Error())
		r.logger.Error(appErr.Error())
		return appErr
	}
	r.logger.Info(fmt.Sprintf("Batch of %d tickets created successfully", len(tickets)))
	return nil
}

func (r *ticketRepositoryImpl) GetTicketByID(ticketID uint) (*models.Ticket, error) {
	r.logger.Info(fmt.Sprintf("Fetching ticket by ID: %s", ticketID))
	var ticket models.Ticket
	if err := r.db.First(&ticket, "id = ?", ticketID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.logger.Warn(fmt.Sprintf("Ticket not found: %s", ticketID))
			return nil, nil
		}
		appErr := utils.NewAppError(500, "Failed to retrieve ticket", err.Error())
		r.logger.Error(appErr.Error())
		return nil, appErr
	}
	r.logger.Info(fmt.Sprintf("Ticket retrieved successfully: %+v", ticket))
	return &ticket, nil
}

func (r *ticketRepositoryImpl) UpdateTicketStatus(ticketID uint, status models.TicketStatus) error {
	r.logger.Info(fmt.Sprintf("Updating status of ticket %s to %s", ticketID, status))
	if err := r.db.Model(&models.Ticket{}).Where("id = ?", ticketID).Update("status", status).Error; err != nil {
		appErr := utils.NewAppError(500, "Failed to update ticket status", err.Error())
		r.logger.Error(appErr.Error())
		return appErr
	}
	r.logger.Info(fmt.Sprintf("Ticket status updated successfully: %s", ticketID))
	return nil
}

func (r *ticketRepositoryImpl) ReserveTicket(ticketID, userID, bookingID uint) error {
	r.logger.Info(fmt.Sprintf("Reserving ticket %s for user %s and booking %s", ticketID, userID, bookingID))
	if err := r.db.Model(&models.Ticket{}).Where("id = ? AND status = ?", ticketID, models.TicketStatusAvailable).
		Updates(map[string]interface{}{
			"status":     models.TicketStatusReserved,
			"user_id":    userID,
			"booking_id": bookingID,
		}).Error; err != nil {
		appErr := utils.NewAppError(500, "Failed to reserve ticket", err.Error())
		r.logger.Error(appErr.Error())
		return appErr
	}
	r.logger.Info(fmt.Sprintf("Ticket reserved successfully: %s", ticketID))
	return nil
}

func (r *ticketRepositoryImpl) ListAvailableTickets(eventID uint) ([]models.Ticket, error) {
	r.logger.Info(fmt.Sprintf("Listing available tickets for event %s", eventID))
	var tickets []models.Ticket
	if err := r.db.Where("event_id = ? AND status = ?", eventID, models.TicketStatusAvailable).Find(&tickets).Error; err != nil {
		appErr := utils.NewAppError(500, "Failed to list available tickets", err.Error())
		r.logger.Error(appErr.Error())
		return nil, appErr
	}
	r.logger.Info(fmt.Sprintf("Found %d available tickets for event %s", len(tickets), eventID))
	return tickets, nil
}

func (r *ticketRepositoryImpl) DeleteTicket(ticketID uint) error {
	r.logger.Info(fmt.Sprintf("Deleting ticket %s", ticketID))
	if err := r.db.Delete(&models.Ticket{}, "id = ?", ticketID).Error; err != nil {
		appErr := utils.NewAppError(500, "Failed to delete ticket", err.Error())
		r.logger.Error(appErr.Error())
		return appErr
	}
	r.logger.Info(fmt.Sprintf("Ticket deleted successfully: %s", ticketID))
	return nil
}
