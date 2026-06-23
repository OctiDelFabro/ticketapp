package dao

import (
	"time"

	"github.com/OctiDelFabro/ticketapp/backend/domain"
	"gorm.io/gorm"
)

type UserStatsRow struct {
	TotalUsers  int64
	ClientUsers int64
	AdminUsers  int64
}

type EventStatsTotalsRow struct {
	TotalEvents       int64
	ActiveEvents      int64
	InactiveEvents    int64
	TotalCapacity     int64
	AvailableCapacity int64
}

type TicketStatsRow struct {
	TotalTickets     int64
	ActiveTickets    int64
	CancelledTickets int64
	EstimatedRevenue float64
}

type EventStatsRow struct {
	EventID          uint
	Title            string
	Description      string
	Category         string
	Location         string
	StartDate        time.Time
	DurationMinutes  int
	Active           bool
	Capacity         int
	Price            float64
	ActiveTickets    int64
	CancelledTickets int64
	TotalTickets     int64
}

func GetUserStats(db *gorm.DB) (UserStatsRow, error) {
	var stats UserStatsRow
	err := db.Model(&domain.User{}).
		Select(`
			COUNT(*) AS total_users,
			COALESCE(SUM(CASE WHEN role = ? THEN 1 ELSE 0 END), 0) AS client_users,
			COALESCE(SUM(CASE WHEN role = ? THEN 1 ELSE 0 END), 0) AS admin_users
		`, "CLIENT", "ADMIN").
		Scan(&stats).Error
	return stats, err
}

func GetEventStatsTotals(db *gorm.DB) (EventStatsTotalsRow, error) {
	var stats EventStatsTotalsRow
	err := db.Model(&domain.Event{}).
		Select(`
			COUNT(*) AS total_events,
			COALESCE(SUM(CASE WHEN active = ? THEN 1 ELSE 0 END), 0) AS active_events,
			COALESCE(SUM(CASE WHEN active = ? THEN 1 ELSE 0 END), 0) AS inactive_events,
			COALESCE(SUM(CASE WHEN active = ? THEN capacity ELSE 0 END), 0) AS total_capacity,
			COALESCE(SUM(CASE WHEN active = ? THEN CASE WHEN capacity - COALESCE(active_ticket_counts.active_tickets, 0) < 0 THEN 0 ELSE capacity - COALESCE(active_ticket_counts.active_tickets, 0) END ELSE 0 END), 0) AS available_capacity
		`, true, false, true, true).
		Joins("LEFT JOIN (SELECT event_id, COUNT(*) AS active_tickets FROM tickets WHERE status = ? GROUP BY event_id) AS active_ticket_counts ON active_ticket_counts.event_id = events.id", ActiveTicketStatus).
		Scan(&stats).Error
	return stats, err
}

func GetTicketStats(db *gorm.DB) (TicketStatsRow, error) {
	var stats TicketStatsRow
	err := db.Model(&domain.Ticket{}).
		Select(`
			COUNT(tickets.id) AS total_tickets,
			COALESCE(SUM(CASE WHEN tickets.status = ? THEN 1 ELSE 0 END), 0) AS active_tickets,
			COALESCE(SUM(CASE WHEN tickets.status = ? THEN 1 ELSE 0 END), 0) AS cancelled_tickets,
			COALESCE(SUM(CASE WHEN tickets.status = ? THEN events.price ELSE 0 END), 0) AS estimated_revenue
		`, ActiveTicketStatus, "CANCELLED", ActiveTicketStatus).
		Joins("LEFT JOIN events ON events.id = tickets.event_id").
		Scan(&stats).Error
	return stats, err
}

func FindEventStats(db *gorm.DB) ([]EventStatsRow, error) {
	var rows []EventStatsRow
	err := eventStatsQuery(db).Order("events.id ASC").Scan(&rows).Error
	return rows, err
}

func FindEventStatsByID(db *gorm.DB, eventID uint) (*EventStatsRow, error) {
	var row EventStatsRow
	if err := eventStatsQuery(db).Where("events.id = ?", eventID).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func eventStatsQuery(db *gorm.DB) *gorm.DB {
	return db.Model(&domain.Event{}).
		Select(`
			events.id AS event_id,
			events.title,
			events.description,
			events.category,
			events.location,
			events.start_date,
			events.duration_minutes,
			events.active,
			events.capacity,
			events.price,
			COALESCE(SUM(CASE WHEN tickets.status = ? THEN 1 ELSE 0 END), 0) AS active_tickets,
			COALESCE(SUM(CASE WHEN tickets.status = ? THEN 1 ELSE 0 END), 0) AS cancelled_tickets,
			COUNT(tickets.id) AS total_tickets
		`, ActiveTicketStatus, "CANCELLED").
		Joins("LEFT JOIN tickets ON tickets.event_id = events.id").
		Group("events.id")
}
