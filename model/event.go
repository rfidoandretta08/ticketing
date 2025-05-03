package model

import "gorm.io/gorm"

type EventStatus string

const (
	Upcoming  EventStatus = "upcoming"
	Ongoing   EventStatus = "ongoing"
	Completed EventStatus = "completed"
)

type Event struct {
	gorm.Model
	Name        string      `gorm:"unique;not null" json:"name"`
	Description string      `gorm:"not null" json:"description"`
	Location    string      `gorm:"not null" json:"location"`
	DateTime    string      `gorm:"not null" json:"date_time"` // Format: "2006-01-02 15:04:05"
	Capacity    int         `gorm:"not null;check:capacity > 0" json:"capacity"`
	Price       float64     `gorm:"not null;check:price >= 0" json:"price"`
	Status      EventStatus `gorm:"type:enum('upcoming','ongoing','completed');default:'upcoming'" json:"status"`
	Tickets     []Ticket    `json:"tickets,omitempty"`
}
