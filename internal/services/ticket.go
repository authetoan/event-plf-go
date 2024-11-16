package services

import (
	"booking-service/internal/models"
	"booking-service/internal/repositories"
	kafkaModels "booking-service/pkg/kafka/models"
	"booking-service/utils"
	"fmt"
)

type TicketService interface {
	CreateTicketsForEvent(eventID uint, numTickets int, price float64) ([]models.Ticket, error)
	GetTicketByID(ticketID uint) (*models.Ticket, error)
	ReserveTicket(ticketID, userID, bookingID uint) error
	ListAvailableTickets(eventID uint) ([]models.Ticket, error)
	UpdateTicketStatus(ticketID uint, status models.TicketStatus) error
	ReleaseTicket(ticketID uint) error
	MarkTicketAsSold(ticketID uint) error
	HandleBookingEvent(eventType string, payload kafkaModels.BookingEvent) error
}

type ticketServiceImpl struct {
	TicketRepo repositories.TicketRepository
	Logger     *utils.Logger
}

func NewTicketService(ticketRepo repositories.TicketRepository) TicketService {
	return &ticketServiceImpl{
		TicketRepo: ticketRepo,
		Logger:     utils.NewLogger(),
	}
}

func (s *ticketServiceImpl) CreateTicketsForEvent(eventID uint, numTickets int, price float64) ([]models.Ticket, error) {
	if numTickets <= 0 || price <= 0 {
		return nil, utils.NewAppError(400, "Invalid input", "Number of tickets or price cannot be zero or negative")
	}

	tickets := make([]models.Ticket, numTickets)
	for i := 0; i < numTickets; i++ {
		tickets[i] = models.Ticket{
			EventID: eventID,
			Price:   price,
			Status:  models.TicketStatusAvailable,
		}
	}

	if err := s.TicketRepo.CreateTicketsBatch(tickets); err != nil {
		return nil, utils.NewAppError(500, "Failed to create tickets", err.Error())
	}

	return tickets, nil
}

func (s *ticketServiceImpl) GetTicketByID(ticketID uint) (*models.Ticket, error) {
	ticket, err := s.TicketRepo.GetTicketByID(ticketID)
	if err != nil {
		return nil, utils.NewAppError(500, "Failed to retrieve ticket", err.Error())
	}
	if ticket == nil {
		return nil, utils.NewAppError(404, "Ticket not found", fmt.Sprintf("Ticket %d not found", ticketID))
	}
	return ticket, nil
}

func (s *ticketServiceImpl) ReserveTicket(ticketID, userID, bookingID uint) error {
	ticket, err := s.GetTicketByID(ticketID)
	if err != nil {
		return err
	}
	if ticket.Status != models.TicketStatusAvailable {
		return utils.NewAppError(400, "Ticket not available", fmt.Sprintf("Ticket %d is not available", ticketID))
	}

	return s.TicketRepo.ReserveTicket(ticketID, userID, bookingID)
}

func (s *ticketServiceImpl) ListAvailableTickets(eventID uint) ([]models.Ticket, error) {
	tickets, err := s.TicketRepo.ListAvailableTickets(eventID)
	if err != nil {
		return nil, utils.NewAppError(500, "Failed to list available tickets", err.Error())
	}
	return tickets, nil
}

func (s *ticketServiceImpl) UpdateTicketStatus(ticketID uint, status models.TicketStatus) error {
	ticket, err := s.GetTicketByID(ticketID)
	if err != nil {
		return err
	}
	if ticket == nil {
		return utils.NewAppError(404, "Ticket not found", fmt.Sprintf("Ticket %d not found", ticketID))
	}

	return s.TicketRepo.UpdateTicketStatus(ticketID, status)
}

func (s *ticketServiceImpl) ReleaseTicket(ticketID uint) error {
	ticket, err := s.GetTicketByID(ticketID)
	if err != nil {
		return err
	}
	if ticket.Status != models.TicketStatusReserved {
		return utils.NewAppError(400, "Ticket not reserved", fmt.Sprintf("Ticket %d is not reserved", ticketID))
	}

	return s.UpdateTicketStatus(ticketID, models.TicketStatusAvailable)
}

func (s *ticketServiceImpl) MarkTicketAsSold(ticketID uint) error {
	ticket, err := s.GetTicketByID(ticketID)
	if err != nil {
		return err
	}
	if ticket.Status != models.TicketStatusReserved {
		return utils.NewAppError(400, "Ticket not reserved", fmt.Sprintf("Ticket %d is not reserved", ticketID))
	}

	return s.UpdateTicketStatus(ticketID, models.TicketStatusSold)
}

func (s *ticketServiceImpl) HandleBookingEvent(eventType string, payload kafkaModels.BookingEvent) error {
	s.Logger.Info(fmt.Sprintf("Handling booking event: %s", eventType))

	switch eventType {
	case "booking.canceled":
		for _, ticketID := range payload.TicketIDs {
			if err := s.ReleaseTicket(ticketID); err != nil {
				s.Logger.Error(fmt.Sprintf("Failed to release ticket %d: %v", ticketID, err))
				return err
			}
		}
	case "booking.confirmed":
		for _, ticketID := range payload.TicketIDs {
			if err := s.MarkTicketAsSold(ticketID); err != nil {
				s.Logger.Error(fmt.Sprintf("Failed to mark ticket %d as sold: %v", ticketID, err))
				return err
			}
		}
	default:
		s.Logger.Warn(fmt.Sprintf("Unhandled event type: %s", eventType))
	}
	return nil
}
