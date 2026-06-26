package domain

import "time"

type Ticket struct {
	ID           uint   `gorm:"primaryKey"`
	UserID       uint   `gorm:"not null;index"`
	EventID      uint   `gorm:"not null;index"`
	User         User   `gorm:"foreignKey:UserID"`
	Event        Event  `gorm:"foreignKey:EventID"`
	GiftedByID   *uint  `gorm:"index"`
	GiftedBy     *User  `gorm:"foreignKey:GiftedByID"`
	GiftMessage  string `gorm:"size:250"`
	GiftedAt     *time.Time
	Status       string    `gorm:"size:20;not null;default:ACTIVE"`
	PurchaseDate time.Time `gorm:"not null"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
