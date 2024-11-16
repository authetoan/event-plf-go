package models

type BookingEvent struct {
	BookingID uint   `json:"booking_id"`
	UserID    uint   `json:"user_id"`
	EventID   uint   `json:"event_id"`
	TicketIDs []uint `json:"ticket_ids"`
}
