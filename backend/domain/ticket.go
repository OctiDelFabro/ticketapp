package domain

import "time"

type Ticket struct {
	ID           uint      `gorm:"primaryKey"`
	UserID       uint      `gorm:"not null;index"`
	EventID      uint      `gorm:"not null;index"`
	User         User      `gorm:"foreignKey:UserID"`
	Event        Event     `gorm:"foreignKey:EventID"`
	Status       string    `gorm:"size:20;not null;default:ACTIVE"`
	PurchaseDate time.Time `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
