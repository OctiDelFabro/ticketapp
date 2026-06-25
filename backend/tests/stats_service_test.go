package tests

import (
	"testing"

	"github.com/OctiDelFabro/ticketapp/backend/domain"
	"github.com/OctiDelFabro/ticketapp/backend/services"
)

func TestStatsSummaryCountsActiveTicketsAndRevenueOnly(t *testing.T) {
	db := newTestDB(t)
	user := createTestUser(t, db)
	event := createTestEvent(t, db, func(event *domain.Event) { event.Capacity = 3; event.Price = 10 })
	createTestTicket(t, db, user.ID, event.ID, "ACTIVE")
	createTestTicket(t, db, user.ID, event.ID, "CANCELLED")

	summary, err := services.GetStatsSummary(db)

	mustNoError(t, err)
	mustEqual(t, int64(2), summary.TotalTickets)
	mustEqual(t, int64(1), summary.ActiveTickets)
	mustEqual(t, int64(1), summary.CancelledTickets)
	mustEqual(t, float64(10), summary.EstimatedRevenue)
	mustEqual(t, int64(2), summary.AvailableCapacity)
}

func TestEventStatsReflectSoldAvailableAndRevenue(t *testing.T) {
	db := newTestDB(t)
	user := createTestUser(t, db)
	event := createTestEvent(t, db, func(event *domain.Event) { event.Capacity = 4; event.Price = 12.5 })
	createTestTicket(t, db, user.ID, event.ID, "ACTIVE")
	createTestTicket(t, db, user.ID, event.ID, "ACTIVE")
	createTestTicket(t, db, user.ID, event.ID, "CANCELLED")

	stats, err := services.GetEventStats(db)

	mustNoError(t, err)
	mustLen(t, stats, 1)
	mustEqual(t, event.ID, stats[0].EventID)
	mustEqual(t, int64(2), stats[0].ActiveTickets)
	mustEqual(t, int64(1), stats[0].CancelledTickets)
	mustEqual(t, 2, stats[0].AvailableCapacity)
	mustEqual(t, float64(25), stats[0].EstimatedRevenue)
}
