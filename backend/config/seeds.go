package config

import (
	"errors"
	"time"

	"github.com/OctiDelFabro/ticketapp/backend/domain"
	"github.com/OctiDelFabro/ticketapp/backend/utils"
	"gorm.io/gorm"
)

const (
	demoPassword       = "123456"
	demoClientRole     = "CLIENT"
	demoAdminRole      = "ADMIN"
	ticketStatusActive = "ACTIVE"
	ticketStatusCancel = "CANCELLED"
)

func SeedEvents(db *gorm.DB) error {
	events := []domain.Event{
		{
			Title:           "Rock Nacional",
			Description:     "Evento de música rock nacional.",
			Category:        "Música",
			Location:        "Córdoba",
			StartDate:       time.Date(2026, time.September, 12, 21, 0, 0, 0, time.Local),
			DurationMinutes: 120,
			Capacity:        100,
			Price:           15000,
			Active:          true,
		},
		{
			Title:           "Stand Up Night",
			Description:     "Noche de comedia con shows de stand up.",
			Category:        "Otros",
			Location:        "Villa Allende",
			StartDate:       time.Date(2026, time.October, 3, 20, 30, 0, 0, time.Local),
			DurationMinutes: 90,
			Capacity:        80,
			Price:           8000,
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
			Price:           20000,
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
			Price:           10000,
			Active:          true,
		},
	}

	for _, event := range events {
		var existing domain.Event
		err := db.Where("title = ?", event.Title).First(&existing).Error
		if err == nil {
			if existing.Price == 0 && event.Price > 0 {
				if err := db.Model(&existing).Update("price", event.Price).Error; err != nil {
					return err
				}
			}
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		if err := db.Create(&event).Error; err != nil {
			return err
		}
	}

	return nil
}

func SeedDemoData(db *gorm.DB) error {
	if err := SeedDemoUsers(db); err != nil {
		return err
	}

	return SeedDemoTickets(db)
}

func SeedDemoUsers(db *gorm.DB) error {
	// Development-only demo users. Existing accounts are left untouched so
	// passwords are never overwritten by the local seed process.
	demoUsers := []domain.User{
		{Name: "Lorenzo", Email: "lorenzo@test.com", Role: demoClientRole},
		{Name: "Pablo", Email: "pablo@test.com", Role: demoClientRole},
		{Name: "Octavio", Email: "octavio@test.com", Role: demoClientRole},
		{Name: "Admin", Email: "admin@test.com", Role: demoAdminRole},
	}

	for _, demoUser := range demoUsers {
		var existing domain.User
		err := db.Where("email = ?", demoUser.Email).First(&existing).Error
		if err == nil {
			continue
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		passwordHash, err := utils.HashPassword(demoPassword)
		if err != nil {
			return err
		}

		demoUser.PasswordHash = passwordHash
		if err := db.Create(&demoUser).Error; err != nil {
			return err
		}
	}

	return nil
}

func SeedDemoTickets(db *gorm.DB) error {
	demoTickets := []struct {
		UserEmail  string
		EventTitle string
		Status     string
	}{
		{UserEmail: "octavio@test.com", EventTitle: "Rock Nacional", Status: ticketStatusActive},
		{UserEmail: "lorenzo@test.com", EventTitle: "Festival Tech", Status: ticketStatusActive},
		{UserEmail: "pablo@test.com", EventTitle: "Stand Up Night", Status: ticketStatusCancel},
	}

	return db.Transaction(func(tx *gorm.DB) error {
		for _, demoTicket := range demoTickets {
			var user domain.User
			if err := tx.Where("email = ?", demoTicket.UserEmail).First(&user).Error; err != nil {
				return err
			}

			var event domain.Event
			if err := tx.Where("title = ?", demoTicket.EventTitle).First(&event).Error; err != nil {
				return err
			}

			var existing domain.Ticket
			err := tx.Where("user_id = ? AND event_id = ?", user.ID, event.ID).First(&existing).Error
			if err == nil {
				continue
			}
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return err
			}

			if demoTicket.Status == ticketStatusActive {
				var activeTickets int64
				if err := tx.Model(&domain.Ticket{}).Where("event_id = ? AND status = ?", event.ID, ticketStatusActive).Count(&activeTickets).Error; err != nil {
					return err
				}
				if event.Capacity-int(activeTickets) <= 0 {
					continue
				}
			}

			ticket := domain.Ticket{
				UserID:       user.ID,
				EventID:      event.ID,
				Status:       demoTicket.Status,
				PurchaseDate: time.Now(),
			}
			if err := tx.Create(&ticket).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
