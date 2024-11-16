package models

import "time"

type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "PENDING"
	BookingStatusConfirmed BookingStatus = "CONFIRMED"
	BookingStatusCanceled  BookingStatus = "CANCELED"
)

type Booking struct {
	ID          uint          `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint          `gorm:"not null" json:"user_id"`
	EventID     uint          `gorm:"not null" json:"event_id"`
	TotalAmount float64       `gorm:"not null" json:"total_amount"`
	Status      BookingStatus `gorm:"not null" json:"status"`
	CreatedAt   time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time     `gorm:"autoUpdateTime" json:"updated_at"`

	Tickets []Ticket `gorm:"foreignKey:BookingID" json:"tickets"`
	Event   Event    `gorm:"foreignKey:EventID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"event"` // Add relationship with Event
}
