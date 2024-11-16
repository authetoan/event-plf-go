package services

import (
	"booking-service/internal/models"
	"booking-service/internal/repositories"
	"booking-service/pkg/kafka"
	"booking-service/utils"
	"fmt"
)

type BookingService interface {
	CreateBooking(userID, eventID uint, ticketIDs []uint) (*models.Booking, error)
	ConfirmBooking(bookingID uint) error
	CancelBooking(bookingID uint) error
	UpdateBookingStatus(bookingID uint, status models.BookingStatus) error
	ListBookingsByUserID(userID uint, page, pageSize int) ([]models.Booking, error)
	GetBookingByID(bookingID uint) (*models.Booking, error)
}

type bookingServiceImpl struct {
	BookingRepo   repositories.BookingRepository
	TicketService TicketService
	Logger        *utils.Logger
	KafkaProducer kafka.Producer
}

func NewBookingService(
	bookingRepo repositories.BookingRepository,
	ticketService TicketService,
	kafkaProducer kafka.Producer,
) BookingService {
	return &bookingServiceImpl{
		BookingRepo:   bookingRepo,
		TicketService: ticketService,
		Logger:        utils.NewLogger(),
		KafkaProducer: kafkaProducer,
	}
}

func (s *bookingServiceImpl) CreateBooking(userID, eventID uint, ticketIDs []uint) (*models.Booking, error) {
	s.Logger.Info(fmt.Sprintf("Creating booking for user %d and event %d", userID, eventID))

	tickets := make([]models.Ticket, 0, len(ticketIDs))
	for _, ticketID := range ticketIDs {
		ticket, err := s.TicketService.GetTicketByID(ticketID)
		if err != nil {
			appErr := utils.NewAppError(500, "Failed to retrieve ticket", err.Error())
			s.Logger.Error(appErr.Error())
			return nil, appErr
		}
		if ticket == nil || ticket.Status != models.TicketStatusAvailable {
			err := utils.NewAppError(400, "Ticket not available", fmt.Sprintf("Ticket %d is not available", ticketID))
			s.Logger.Warn(err.Error())
			return nil, err
		}
		tickets = append(tickets, *ticket)
	}

	var totalAmount float64
	for _, ticket := range tickets {
		totalAmount += ticket.Price
	}

	booking := &models.Booking{
		UserID:      userID,
		EventID:     eventID,
		TotalAmount: totalAmount,
		Status:      models.BookingStatusPending,
		Tickets:     tickets,
	}

	if err := s.BookingRepo.CreateBooking(booking); err != nil {
		appErr := utils.NewAppError(500, "Failed to create booking", err.Error())
		s.Logger.Error(appErr.Error())
		return nil, appErr
	}

	for _, ticket := range tickets {
		if err := s.TicketService.ReserveTicket(ticket.ID, userID, booking.ID); err != nil {
			appErr := utils.NewAppError(500, "Failed to reserve ticket", err.Error())
			s.Logger.Error(appErr.Error())
			return nil, appErr
		}
	}

	if err := s.publishEvent("booking.created", booking); err != nil {
		s.Logger.Error(fmt.Sprintf("Failed to publish booking.created event: %v", err))
	}

	s.Logger.Info(fmt.Sprintf("Booking created successfully: %+v", booking))
	return booking, nil
}

func (s *bookingServiceImpl) ConfirmBooking(bookingID uint) error {
	s.Logger.Info(fmt.Sprintf("Confirming booking %d", bookingID))

	if err := s.BookingRepo.UpdateBookingStatus(bookingID, models.BookingStatusConfirmed); err != nil {
		appErr := utils.NewAppError(500, "Failed to confirm booking", err.Error())
		s.Logger.Error(appErr.Error())
		return appErr
	}

	booking, err := s.BookingRepo.GetBookingByID(bookingID)
	if err != nil {
		appErr := utils.NewAppError(500, "Failed to retrieve booking", err.Error())
		s.Logger.Error(appErr.Error())
		return appErr
	}
	if booking == nil {
		err := utils.NewAppError(404, "Booking not found", fmt.Sprintf("Booking ID %d not found", bookingID))
		s.Logger.Warn(err.Error())
		return err
	}

	for _, ticket := range booking.Tickets {
		if err := s.TicketService.UpdateTicketStatus(ticket.ID, models.TicketStatusSold); err != nil {
			appErr := utils.NewAppError(500, fmt.Sprintf("Failed to mark ticket %d as sold", ticket.ID), err.Error())
			s.Logger.Error(appErr.Error())
		}
	}

	if err := s.publishEvent("booking.confirmed", booking); err != nil {
		s.Logger.Error(fmt.Sprintf("Failed to publish booking.confirmed event: %v", err))
	}

	s.Logger.Info(fmt.Sprintf("Booking %d confirmed successfully", bookingID))
	return nil
}

func (s *bookingServiceImpl) CancelBooking(bookingID uint) error {
	s.Logger.Info(fmt.Sprintf("Cancelling booking %d", bookingID))

	if err := s.BookingRepo.UpdateBookingStatus(bookingID, models.BookingStatusCanceled); err != nil {
		appErr := utils.NewAppError(500, "Failed to cancel booking", err.Error())
		s.Logger.Error(appErr.Error())
		return appErr
	}

	booking, err := s.BookingRepo.GetBookingByID(bookingID)
	if err != nil {
		appErr := utils.NewAppError(500, "Failed to retrieve booking", err.Error())
		s.Logger.Error(appErr.Error())
		return appErr
	}
	if booking == nil {
		err := utils.NewAppError(404, "Booking not found", fmt.Sprintf("Booking ID %d not found", bookingID))
		s.Logger.Warn(err.Error())
		return err
	}

	for _, ticket := range booking.Tickets {
		if err := s.TicketService.UpdateTicketStatus(ticket.ID, models.TicketStatusAvailable); err != nil {
			appErr := utils.NewAppError(500, fmt.Sprintf("Failed to release ticket %d", ticket.ID), err.Error())
			s.Logger.Error(appErr.Error())
		}
	}

	if err := s.publishEvent("booking.canceled", booking); err != nil {
		s.Logger.Error(fmt.Sprintf("Failed to publish booking.canceled event: %v", err))
	}

	s.Logger.Info(fmt.Sprintf("Booking %d cancelled successfully", bookingID))
	return nil
}

func (s *bookingServiceImpl) UpdateBookingStatus(bookingID uint, status models.BookingStatus) error {
	s.Logger.Info(fmt.Sprintf("Updating booking status for %d to %s", bookingID, status))
	if err := s.BookingRepo.UpdateBookingStatus(bookingID, status); err != nil {
		appErr := utils.NewAppError(500, "Failed to update booking status", err.Error())
		s.Logger.Error(appErr.Error())
		return appErr
	}

	s.Logger.Info(fmt.Sprintf("Booking status updated successfully: %d", bookingID))
	return nil
}

func (s *bookingServiceImpl) ListBookingsByUserID(userID uint, page, pageSize int) ([]models.Booking, error) {
	s.Logger.Info(fmt.Sprintf("Listing bookings for user %d", userID))
	bookings, err := s.BookingRepo.ListBookingsByUserID(userID, page, pageSize)
	if err != nil {
		appErr := utils.NewAppError(500, "Failed to list bookings", err.Error())
		s.Logger.Error(appErr.Error())
		return nil, appErr
	}
	s.Logger.Info(fmt.Sprintf("Retrieved %d bookings for user %d", len(bookings), userID))
	return bookings, nil
}

func (s *bookingServiceImpl) GetBookingByID(bookingID uint) (*models.Booking, error) {
	s.Logger.Info(fmt.Sprintf("Fetching booking by ID %d", bookingID))
	booking, err := s.BookingRepo.GetBookingByID(bookingID)
	if err != nil {
		appErr := utils.NewAppError(500, "Failed to retrieve booking", err.Error())
		s.Logger.Error(appErr.Error())
		return nil, appErr
	}
	if booking == nil {
		err := utils.NewAppError(404, "Booking not found", fmt.Sprintf("Booking ID %d not found", bookingID))
		s.Logger.Warn(err.Error())
		return nil, err
	}
	s.Logger.Info(fmt.Sprintf("Booking retrieved successfully: %+v", booking))
	return booking, nil
}

func (s *bookingServiceImpl) publishEvent(eventType string, payload interface{}) error {
	topic := fmt.Sprintf("booking.%s", eventType)
	return s.KafkaProducer.Publish(topic, payload)
}
