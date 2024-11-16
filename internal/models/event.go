package models

import "time"

type Event struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Date      time.Time `gorm:"not null" json:"date"`
	Location  string    `gorm:"not null" json:"location"`
	Capacity  int       `gorm:"not null" json:"capacity"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	Tickets  []Ticket  `gorm:"foreignKey:EventID" json:"tickets"`  // One-to-many relationship with Ticket
	Bookings []Booking `gorm:"foreignKey:EventID" json:"bookings"` // One-to-many relationship with Booking
}
