package model

import "gorm.io/gorm"

type TicketStatus string

const (
	Available TicketStatus = "available"
	Booked    TicketStatus = "booked"
	Cancelled TicketStatus = "cancelled"
)

type PaymentStatus string

const (
	Pending PaymentStatus = "waiting"
	Success PaymentStatus = "success"
	Cancel  PaymentStatus = "cancel"
)

type Ticket struct {
	gorm.Model
	EventID       uint          `gorm:"not null" json:"event_id"`
	Event         Event         `gorm:"foreignKey:EventID" json:"event"`
	UserID        uint          `gorm:"not null" json:"user_id"`
	User          User          `gorm:"foreignKey:UserID" json:"user"`
	Qty           int           `gorm:"not null" json:"qty"`
	SubTotal      float64       `json:"sub_total"`
	Status        TicketStatus  `gorm:"type:enum('available','booked','cancelled');default:'available'" json:"status"`
	PaymentStatus PaymentStatus `gorm:"type:enum('waiting','success','cancel');default:'waiting'" json:"role"`
	BookingDate   string        `gorm:"not null" json:"booking_date"` // Format: "2006-01-02 15:04:05"
}
