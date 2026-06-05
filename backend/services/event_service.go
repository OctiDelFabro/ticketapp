package services

import (
	"errors"
	"time"

	"github.com/OctiDelFabro/ticketapp/backend/dao"
	"github.com/OctiDelFabro/ticketapp/backend/domain"
	"gorm.io/gorm"
)

var ErrEventNotFound = errors.New("event not found")

type EventResponse struct {
	ID                uint      `json:"id"`
	Title             string    `json:"title"`
	Description       string    `json:"description"`
	ImageURL          string    `json:"image_url"`
	Category          string    `json:"category"`
	Location          string    `json:"location"`
	StartDate         time.Time `json:"start_date"`
	DurationMinutes   int       `json:"duration_minutes"`
	Capacity          int       `json:"capacity"`
	AvailableCapacity int       `json:"available_capacity"`
	Active            bool      `json:"active"`
}

func ListEvents(db *gorm.DB, filters dao.EventFilters) ([]EventResponse, error) {
	events, err := dao.FindEvents(db, filters)
	if err != nil {
		return nil, err
	}

	responses := make([]EventResponse, 0, len(events))
	for _, event := range events {
		response, err := buildEventResponse(db, event)
		if err != nil {
			return nil, err
		}

		responses = append(responses, response)
	}

	return responses, nil
}

func GetEventByID(db *gorm.DB, id uint) (*EventResponse, error) {
	event, err := dao.FindEventByID(db, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEventNotFound
		}
		return nil, err
	}

	response, err := buildEventResponse(db, *event)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func buildEventResponse(db *gorm.DB, event domain.Event) (EventResponse, error) {
	availableCapacity, err := calculateAvailableCapacity(db, event)
	if err != nil {
		return EventResponse{}, err
	}

	return EventResponse{
		ID:                event.ID,
		Title:             event.Title,
		Description:       event.Description,
		ImageURL:          event.ImageURL,
		Category:          event.Category,
		Location:          event.Location,
		StartDate:         event.StartDate,
		DurationMinutes:   event.DurationMinutes,
		Capacity:          event.Capacity,
		AvailableCapacity: availableCapacity,
		Active:            event.Active,
	}, nil
}

func calculateAvailableCapacity(db *gorm.DB, event domain.Event) (int, error) {
	var activeTickets int64
	if err := db.Model(&domain.Ticket{}).
		Where("event_id = ? AND status = ?", event.ID, "ACTIVE").
		Count(&activeTickets).Error; err != nil {
		return 0, err
	}

	availableCapacity := event.Capacity - int(activeTickets)
	if availableCapacity < 0 {
		return 0, nil
	}

	return availableCapacity, nil
}
