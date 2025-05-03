package model

import "gorm.io/gorm"

type TicketStatus string

const (
	Available TicketStatus = "available"
	Booked    TicketStatus = "booked"
	Cancelled TicketStatus = "cancelled"
)

type Ticket struct {
	gorm.Model
	EventID     uint         `gorm:"not null" json:"event_id"`
	Event       Event        `gorm:"foreignKey:EventID" json:"event"`
	UserID      uint         `gorm:"not null" json:"user_id"`
	User        User         `gorm:"foreignKey:UserID" json:"user"`
	Status      TicketStatus `gorm:"type:enum('available','booked','cancelled');default:'available'" json:"status"`
	BookingDate string       `gorm:"not null" json:"booking_date"` // Format: "2006-01-02 15:04:05"
}
