package services

import (
	"errors"
	"strings"
	"time"

	"github.com/OctiDelFabro/ticketapp/backend/dao"
	"github.com/OctiDelFabro/ticketapp/backend/domain"
	"gorm.io/gorm"
)

var (
	ErrEventNotFound        = errors.New("event not found")
	ErrInvalidEventRequest  = errors.New("invalid event request")
	ErrCapacityBelowTickets = errors.New("capacity cannot be lower than active tickets sold")
)

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

type CreateEventRequest struct {
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	ImageURL        string    `json:"image_url"`
	Category        string    `json:"category"`
	Location        string    `json:"location"`
	StartDate       time.Time `json:"start_date"`
	DurationMinutes int       `json:"duration_minutes"`
	Capacity        int       `json:"capacity"`
	Active          bool      `json:"active"`
}

type UpdateEventRequest struct {
	Title           *string    `json:"title"`
	Description     *string    `json:"description"`
	ImageURL        *string    `json:"image_url"`
	Category        *string    `json:"category"`
	Location        *string    `json:"location"`
	StartDate       *time.Time `json:"start_date"`
	DurationMinutes *int       `json:"duration_minutes"`
	Capacity        *int       `json:"capacity"`
	Active          *bool      `json:"active"`
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

func CreateEvent(db *gorm.DB, req CreateEventRequest) (*EventResponse, error) {
	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)
	req.Category = strings.TrimSpace(req.Category)
	req.Location = strings.TrimSpace(req.Location)
	req.ImageURL = strings.TrimSpace(req.ImageURL)

	if req.Title == "" || req.Description == "" || req.Category == "" || req.Location == "" || req.StartDate.IsZero() || req.DurationMinutes <= 0 || req.Capacity <= 0 {
		return nil, ErrInvalidEventRequest
	}

	event := domain.Event{
		Title:           req.Title,
		Description:     req.Description,
		ImageURL:        req.ImageURL,
		Category:        req.Category,
		Location:        req.Location,
		StartDate:       req.StartDate,
		DurationMinutes: req.DurationMinutes,
		Capacity:        req.Capacity,
		Active:          req.Active,
	}

	if err := dao.CreateEvent(db, &event); err != nil {
		return nil, err
	}

	response, err := buildEventResponse(db, event)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func UpdateEvent(db *gorm.DB, id uint, req UpdateEventRequest) (*EventResponse, error) {
	var updatedEvent domain.Event
	err := db.Transaction(func(tx *gorm.DB) error {
		event, err := dao.FindEventByIDIncludingInactive(tx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrEventNotFound
			}
			return err
		}

		if req.Title != nil {
			event.Title = strings.TrimSpace(*req.Title)
		}
		if req.Description != nil {
			event.Description = strings.TrimSpace(*req.Description)
		}
		if req.ImageURL != nil {
			event.ImageURL = strings.TrimSpace(*req.ImageURL)
		}
		if req.Category != nil {
			event.Category = strings.TrimSpace(*req.Category)
		}
		if req.Location != nil {
			event.Location = strings.TrimSpace(*req.Location)
		}
		if req.StartDate != nil {
			event.StartDate = *req.StartDate
		}
		if req.DurationMinutes != nil {
			if *req.DurationMinutes <= 0 {
				return ErrInvalidEventRequest
			}
			event.DurationMinutes = *req.DurationMinutes
		}
		if req.Capacity != nil {
			if *req.Capacity <= 0 {
				return ErrInvalidEventRequest
			}

			activeTickets, err := dao.CountActiveTicketsByEventID(tx, id)
			if err != nil {
				return err
			}
			if *req.Capacity < int(activeTickets) {
				return ErrCapacityBelowTickets
			}
			event.Capacity = *req.Capacity
		}
		if req.Active != nil {
			event.Active = *req.Active
		}

		if event.Title == "" || event.Description == "" || event.Category == "" || event.Location == "" || event.StartDate.IsZero() {
			return ErrInvalidEventRequest
		}

		if err := dao.UpdateEvent(tx, event); err != nil {
			return err
		}

		updatedEvent = *event
		return nil
	})
	if err != nil {
		return nil, err
	}

	response, err := buildEventResponse(db, updatedEvent)
	if err != nil {
		return nil, err
	}

	return &response, nil
}

func DisableEvent(db *gorm.DB, id uint) error {
	return db.Transaction(func(tx *gorm.DB) error {
		event, err := dao.FindEventByIDIncludingInactive(tx, id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrEventNotFound
			}
			return err
		}

		event.Active = false
		return dao.UpdateEvent(tx, event)
	})
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
