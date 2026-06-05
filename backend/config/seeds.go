package config

import (
	"time"

	"github.com/OctiDelFabro/ticketapp/backend/domain"
	"gorm.io/gorm"
)

func SeedEvents(db *gorm.DB) error {
	var count int64
	if err := db.Model(&domain.Event{}).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	events := []domain.Event{
		{
			Title:           "Rock Nacional",
			Description:     "Evento de música rock nacional.",
			Category:        "Música",
			Location:        "Córdoba",
			StartDate:       time.Date(2026, time.September, 12, 21, 0, 0, 0, time.Local),
			DurationMinutes: 120,
			Capacity:        100,
			Active:          true,
		},
		{
			Title:           "Stand Up Night",
			Description:     "Noche de comedia con shows de stand up.",
			Category:        "Comedia",
			Location:        "Villa Allende",
			StartDate:       time.Date(2026, time.October, 3, 20, 30, 0, 0, time.Local),
			DurationMinutes: 90,
			Capacity:        80,
			Active:          true,
		},
		{
			Title:           "Festival Tech",
			Description:     "Festival de tecnología con charlas y actividades.",
			Category:        "Tecnología",
			Location:        "Córdoba",
			StartDate:       time.Date(2026, time.November, 7, 15, 0, 0, 0, time.Local),
			DurationMinutes: 180,
			Capacity:        150,
			Active:          true,
		},
		{
			Title:           "Obra de Teatro",
			Description:     "Función teatral para todo público.",
			Category:        "Teatro",
			Location:        "Córdoba",
			StartDate:       time.Date(2026, time.December, 5, 19, 0, 0, 0, time.Local),
			DurationMinutes: 100,
			Capacity:        60,
			Active:          true,
		},
	}

	return db.Create(&events).Error
}
