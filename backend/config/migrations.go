package config

import (
	"github.com/OctiDelFabro/ticketapp/backend/domain"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	return db.AutoMigrate(&domain.User{}, &domain.Event{}, &domain.Ticket{})
}
