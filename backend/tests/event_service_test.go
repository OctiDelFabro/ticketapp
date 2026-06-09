package tests

import (
	"testing"

	"github.com/OctiDelFabro/ticketapp/backend/dao"
	"github.com/OctiDelFabro/ticketapp/backend/domain"
	"github.com/OctiDelFabro/ticketapp/backend/services"
)

func TestListEventsReturnsOnlyActiveEvents(t *testing.T) {
	db := newTestDB(t)
	active := createTestEvent(t, db, func(event *domain.Event) { event.Title = "Active Concert" })
	createTestEvent(t, db, func(event *domain.Event) {
		event.Title = "Inactive Concert"
		event.Active = false
	})

	events, err := services.ListEvents(db, dao.EventFilters{})

	mustNoError(t, err)
	mustLen(t, events, 1)
	mustEqual(t, active.ID, events[0].ID)
	mustTrue(t, events[0].Active)
}

func TestListEventsFiltersBySearch(t *testing.T) {
	db := newTestDB(t)
	matched := createTestEvent(t, db, func(event *domain.Event) { event.Title = "Rock Festival" })
	createTestEvent(t, db, func(event *domain.Event) { event.Title = "Tech Meetup" })

	events, err := services.ListEvents(db, dao.EventFilters{Search: "Rock"})

	mustNoError(t, err)
	mustLen(t, events, 1)
	mustEqual(t, matched.ID, events[0].ID)
}

func TestListEventsFiltersByCategory(t *testing.T) {
	db := newTestDB(t)
	matched := createTestEvent(t, db, func(event *domain.Event) { event.Category = "Sports" })
	createTestEvent(t, db, func(event *domain.Event) { event.Category = "Music" })

	events, err := services.ListEvents(db, dao.EventFilters{Category: "Sports"})

	mustNoError(t, err)
	mustLen(t, events, 1)
	mustEqual(t, matched.ID, events[0].ID)
	mustEqual(t, "Sports", events[0].Category)
}

func TestListEventsAvailableOnlyReturnsEventsWithCapacity(t *testing.T) {
	db := newTestDB(t)
	user := createTestUser(t, db)
	available := createTestEvent(t, db, func(event *domain.Event) {
		event.Title = "Available Event"
		event.Capacity = 2
	})
	full := createTestEvent(t, db, func(event *domain.Event) {
		event.Title = "Full Event"
		event.Capacity = 1
	})
	createTestTicket(t, db, user.ID, full.ID, "ACTIVE")

	events, err := services.ListEvents(db, dao.EventFilters{AvailableOnly: true})

	mustNoError(t, err)
	mustLen(t, events, 1)
	mustEqual(t, available.ID, events[0].ID)
}

func TestGetEventByIDReturnsExistingActiveEvent(t *testing.T) {
	db := newTestDB(t)
	event := createTestEvent(t, db, func(event *domain.Event) { event.Title = "Active Event" })

	response, err := services.GetEventByID(db, event.ID)

	mustNoError(t, err)
	mustNotNil(t, response)
	mustEqual(t, event.ID, response.ID)
	mustEqual(t, "Active Event", response.Title)
}

func TestGetEventByIDReturnsErrorWhenMissing(t *testing.T) {
	db := newTestDB(t)

	_, err := services.GetEventByID(db, 999)

	mustErrorIs(t, err, services.ErrEventNotFound)
}

func TestGetEventByIDReturnsErrorWhenInactive(t *testing.T) {
	db := newTestDB(t)
	event := createTestEvent(t, db, func(event *domain.Event) { event.Active = false })

	_, err := services.GetEventByID(db, event.ID)

	mustErrorIs(t, err, services.ErrEventNotFound)
}

func TestAvailableCapacityCountsOnlyActiveTickets(t *testing.T) {
	db := newTestDB(t)
	userOne := createTestUser(t, db, func(user *domain.User) { user.Email = "one@example.com" })
	userTwo := createTestUser(t, db, func(user *domain.User) { user.Email = "two@example.com" })
	userThree := createTestUser(t, db, func(user *domain.User) { user.Email = "three@example.com" })
	event := createTestEvent(t, db, func(event *domain.Event) { event.Capacity = 3 })
	createTestTicket(t, db, userOne.ID, event.ID, "ACTIVE")
	createTestTicket(t, db, userTwo.ID, event.ID, "ACTIVE")
	createTestTicket(t, db, userThree.ID, event.ID, "CANCELLED")

	response, err := services.GetEventByID(db, event.ID)

	mustNoError(t, err)
	mustNotNil(t, response)
	mustEqual(t, 1, response.AvailableCapacity)
}
