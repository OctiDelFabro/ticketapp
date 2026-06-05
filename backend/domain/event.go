package domain

import "time"

type Event struct {
	ID              uint      `gorm:"primaryKey"`
	Title           string    `gorm:"size:150;not null"`
	Description     string    `gorm:"type:text;not null"`
	ImageURL        string    `gorm:"size:255"`
	Category        string    `gorm:"size:100;not null"`
	Location        string    `gorm:"size:150;not null"`
	StartDate       time.Time `gorm:"not null"`
	DurationMinutes int       `gorm:"not null"`
	Capacity        int       `gorm:"not null"`
	Active          bool      `gorm:"not null;default:true"`
	Tickets         []Ticket
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
