package domain

import "time"

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Name         string `gorm:"size:100;not null"`
	Email        string `gorm:"size:150;uniqueIndex;not null"`
	PasswordHash string `gorm:"size:255;not null"`
	Role         string `gorm:"size:20;not null;default:CLIENT"`
	Tickets      []Ticket
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
