package models

import "time"

type TicketStatus string

const (
	TicketStatusAvailable TicketStatus = "AVAILABLE"
	TicketStatusReserved  TicketStatus = "RESERVED"
	TicketStatusSold      TicketStatus = "SOLD"
)

type Ticket struct {
	ID        uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	EventID   uint         `gorm:"not null" json:"event_id"`
	Price     float64      `gorm:"not null" json:"price"`
	Status    TicketStatus `gorm:"not null" json:"status"`
	UserID    *uint        `json:"user_id,omitempty"`    // Nullable, only set when reserved
	BookingID *uint        `json:"booking_id,omitempty"` // Nullable, only set when booked
	CreatedAt time.Time    `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time    `gorm:"autoUpdateTime" json:"updated_at"`

	Event Event `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"event"` // Add relationship with Event
}
