package dao

import (
	"strings"

	"github.com/OctiDelFabro/ticketapp/backend/domain"
	"gorm.io/gorm"
)

type EventFilters struct {
	Search        string
	Category      string
	AvailableOnly bool
}

func FindEvents(db *gorm.DB, filters EventFilters) ([]domain.Event, error) {
	var events []domain.Event

	query := db.Where("active = ?", true)

	if search := strings.TrimSpace(filters.Search); search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}

	if category := strings.TrimSpace(filters.Category); category != "" {
		query = query.Where("category = ?", category)
	}

	if filters.AvailableOnly {
		activeTicketsSubquery := db.Model(&domain.Ticket{}).
			Select("COUNT(*)").
			Where("tickets.event_id = events.id").
			Where("tickets.status = ?", "ACTIVE")

		query = query.Where("capacity > (?)", activeTicketsSubquery)
	}

	if err := query.Find(&events).Error; err != nil {
		return nil, err
	}

	return events, nil
}

func FindEventByID(db *gorm.DB, id uint) (*domain.Event, error) {
	var event domain.Event
	if err := db.Where("id = ? AND active = ?", id, true).First(&event).Error; err != nil {
		return nil, err
	}

	return &event, nil
}

func FindEventByIDIncludingInactive(db *gorm.DB, id uint) (*domain.Event, error) {
	var event domain.Event
	if err := db.First(&event, id).Error; err != nil {
		return nil, err
	}

	return &event, nil
}

func CreateEvent(db *gorm.DB, event *domain.Event) error {
	return db.Create(event).Error
}

func UpdateEvent(db *gorm.DB, event *domain.Event) error {
	return db.Save(event).Error
}
