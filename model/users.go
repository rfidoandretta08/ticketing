package model

import "gorm.io/gorm"

type Role string

const (
	Admin Role = "admin"
	Users Role = "user"
)

type User struct {
	gorm.Model
	Name     string   `gorm:"not null" json:"name"`
	Password string   `gorm:"not null" json:"-"`
	Email    string   `gorm:"unique;not null" json:"email"`
	Role     Role     `gorm:"type:enum('admin','user');default:'user'" json:"role"`
	Tickets  []Ticket `json:"tickets,omitempty"`
}
