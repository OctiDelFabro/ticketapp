package services

import (
	"errors"
	"math"
	"time"

	"github.com/OctiDelFabro/ticketapp/backend/dao"
	"gorm.io/gorm"
)

type StatsSummaryResponse struct {
	TotalUsers           int64   `json:"total_users"`
	ClientUsers          int64   `json:"client_users"`
	AdminUsers           int64   `json:"admin_users"`
	TotalEvents          int64   `json:"total_events"`
	ActiveEvents         int64   `json:"active_events"`
	InactiveEvents       int64   `json:"inactive_events"`
	TotalTickets         int64   `json:"total_tickets"`
	ActiveTickets        int64   `json:"active_tickets"`
	CancelledTickets     int64   `json:"cancelled_tickets"`
	TotalCapacity        int64   `json:"total_capacity"`
	AvailableCapacity    int64   `json:"available_capacity"`
	OccupancyRatePercent float64 `json:"occupancy_rate_percent"`
	EstimatedRevenue     float64 `json:"estimated_revenue"`
}

type EventStatsResponse struct {
	EventID              uint    `json:"event_id"`
	Title                string  `json:"title"`
	Category             string  `json:"category"`
	Location             string  `json:"location"`
	Active               bool    `json:"active"`
	Capacity             int     `json:"capacity"`
	Price                float64 `json:"price"`
	ActiveTickets        int64   `json:"active_tickets"`
	CancelledTickets     int64   `json:"cancelled_tickets"`
	TotalTickets         int64   `json:"total_tickets"`
	AvailableCapacity    int     `json:"available_capacity"`
	OccupancyRatePercent float64 `json:"occupancy_rate_percent"`
	EstimatedRevenue     float64 `json:"estimated_revenue"`
}

type EventReportResponse struct {
	EventID              uint      `json:"event_id"`
	Title                string    `json:"title"`
	Description          string    `json:"description"`
	Category             string    `json:"category"`
	Location             string    `json:"location"`
	StartDate            time.Time `json:"start_date"`
	DurationMinutes      int       `json:"duration_minutes"`
	Active               bool      `json:"active"`
	Capacity             int       `json:"capacity"`
	Price                float64   `json:"price"`
	ActiveTickets        int64     `json:"active_tickets"`
	CancelledTickets     int64     `json:"cancelled_tickets"`
	TotalTickets         int64     `json:"total_tickets"`
	AvailableCapacity    int       `json:"available_capacity"`
	OccupancyRatePercent float64   `json:"occupancy_rate_percent"`
	EstimatedRevenue     float64   `json:"estimated_revenue"`
}

func GetStatsSummary(db *gorm.DB) (*StatsSummaryResponse, error) {
	userStats, err := dao.GetUserStats(db)
	if err != nil {
		return nil, err
	}
	eventStats, err := dao.GetEventStatsTotals(db)
	if err != nil {
		return nil, err
	}
	ticketStats, err := dao.GetTicketStats(db)
	if err != nil {
		return nil, err
	}

	availableCapacity := eventStats.AvailableCapacity
	if availableCapacity < 0 {
		availableCapacity = 0
	}

	return &StatsSummaryResponse{
		TotalUsers:           userStats.TotalUsers,
		ClientUsers:          userStats.ClientUsers,
		AdminUsers:           userStats.AdminUsers,
		TotalEvents:          eventStats.TotalEvents,
		ActiveEvents:         eventStats.ActiveEvents,
		InactiveEvents:       eventStats.InactiveEvents,
		TotalTickets:         ticketStats.TotalTickets,
		ActiveTickets:        ticketStats.ActiveTickets,
		CancelledTickets:     ticketStats.CancelledTickets,
		TotalCapacity:        eventStats.TotalCapacity,
		AvailableCapacity:    availableCapacity,
		OccupancyRatePercent: percentage(ticketStats.ActiveTickets, eventStats.TotalCapacity),
		EstimatedRevenue:     ticketStats.EstimatedRevenue,
	}, nil
}

func GetEventStats(db *gorm.DB) ([]EventStatsResponse, error) {
	rows, err := dao.FindEventStats(db)
	if err != nil {
		return nil, err
	}

	responses := make([]EventStatsResponse, 0, len(rows))
	for _, row := range rows {
		responses = append(responses, buildEventStatsResponse(row))
	}
	return responses, nil
}

func GetEventReport(db *gorm.DB, eventID uint) (*EventReportResponse, error) {
	row, err := dao.FindEventStatsByID(db, eventID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrEventNotFound
		}
		return nil, err
	}

	stats := buildEventStatsResponse(*row)
	return &EventReportResponse{
		EventID:              stats.EventID,
		Title:                stats.Title,
		Description:          row.Description,
		Category:             stats.Category,
		Location:             stats.Location,
		StartDate:            row.StartDate,
		DurationMinutes:      row.DurationMinutes,
		Active:               stats.Active,
		Capacity:             stats.Capacity,
		Price:                stats.Price,
		ActiveTickets:        stats.ActiveTickets,
		CancelledTickets:     stats.CancelledTickets,
		TotalTickets:         stats.TotalTickets,
		AvailableCapacity:    stats.AvailableCapacity,
		OccupancyRatePercent: stats.OccupancyRatePercent,
		EstimatedRevenue:     stats.EstimatedRevenue,
	}, nil
}

func buildEventStatsResponse(row dao.EventStatsRow) EventStatsResponse {
	availableCapacity := row.Capacity - int(row.ActiveTickets)
	if availableCapacity < 0 {
		availableCapacity = 0
	}

	return EventStatsResponse{
		EventID:              row.EventID,
		Title:                row.Title,
		Category:             row.Category,
		Location:             row.Location,
		Active:               row.Active,
		Capacity:             row.Capacity,
		Price:                row.Price,
		ActiveTickets:        row.ActiveTickets,
		CancelledTickets:     row.CancelledTickets,
		TotalTickets:         row.TotalTickets,
		AvailableCapacity:    availableCapacity,
		OccupancyRatePercent: percentage(row.ActiveTickets, int64(row.Capacity)),
		EstimatedRevenue:     float64(row.ActiveTickets) * row.Price,
	}
}

func percentage(numerator int64, denominator int64) float64 {
	if denominator <= 0 {
		return 0
	}
	return math.Round((float64(numerator)/float64(denominator))*10000) / 100
}
