package tests

import (
	"testing"
	"time"

	"github.com/OctiDelFabro/ticketapp/backend/dao"
	"github.com/OctiDelFabro/ticketapp/backend/domain"
	"github.com/OctiDelFabro/ticketapp/backend/services"
)

func TestListEventsReturnsOnlyActiveEvents(t *testing.T) {
	db := newTestDB(t)
	active := domain.Event{
		Title:           "Active Concert",
		Description:     "An active event description",
		ImageURL:        "https://example.com/active-event.jpg",
		Category:        "Music",
		Location:        "Buenos Aires",
		StartDate:       time.Now().Add(48 * time.Hour).UTC().Truncate(time.Second),
		DurationMinutes: 120,
		Capacity:        10,
		Active:          true,
	}
	inactive := domain.Event{
		Title:           "Inactive Concert",
		Description:     "An inactive event description",
		ImageURL:        "https://example.com/inactive-event.jpg",
		Category:        "Music",
		Location:        "Buenos Aires",
		StartDate:       time.Now().Add(48 * time.Hour).UTC().Truncate(time.Second),
		DurationMinutes: 120,
		Capacity:        10,
		Active:          true,
	}
	mustNoError(t, db.Create(&active).Error)
	mustNoError(t, db.Create(&inactive).Error)
	mustNoError(t, db.Model(&inactive).Update("active", false).Error)

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
	event := createTestEvent(t, db)
	event.Active = false
	mustNoError(t, db.Save(&event).Error)

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

func validCreateEventRequest() services.CreateEventRequest {
	return services.CreateEventRequest{
		Title:           "Nuevo evento",
		Description:     "Descripción del evento",
		ImageURL:        "https://example.com/nuevo.jpg",
		Category:        "Música",
		Location:        "Buenos Aires",
		StartDate:       time.Now().Add(72 * time.Hour).UTC().Truncate(time.Second),
		DurationMinutes: 90,
		Capacity:        50,
		Price:           25.5,
	}
}

func TestCreateEventCreatesValidEvent(t *testing.T) {
	db := newTestDB(t)
	req := validCreateEventRequest()

	response, err := services.CreateEvent(db, req)

	mustNoError(t, err)
	mustNotNil(t, response)
	mustEqual(t, "Nuevo evento", response.Title)
	mustEqual(t, "Música", response.Category)
	mustEqual(t, 50, response.Capacity)
	mustEqual(t, 50, response.AvailableCapacity)
	mustEqual(t, 0, response.TicketsSold)
	mustTrue(t, response.Active)
}

func TestCreateEventCategoryValidation(t *testing.T) {
	db := newTestDB(t)
	req := validCreateEventRequest()
	req.Category = "   "
	_, err := services.CreateEvent(db, req)
	mustErrorIs(t, err, services.ErrInvalidEventRequest)

	req = validCreateEventRequest()
	req.Category = "Music"
	_, err = services.CreateEvent(db, req)
	mustErrorIs(t, err, services.ErrInvalidEventCategory)

	for _, category := range []string{"Música", "Teatro", "Deportes", "Tecnología", "Otros"} {
		t.Run(category, func(t *testing.T) {
			req := validCreateEventRequest()
			req.Title = "Evento " + category
			req.Category = category
			response, err := services.CreateEvent(db, req)
			mustNoError(t, err)
			mustEqual(t, category, response.Category)
		})
	}
}

func TestUpdateEventCategoryAndCapacityValidation(t *testing.T) {
	db := newTestDB(t)
	userOne := createTestUser(t, db, func(user *domain.User) { user.Email = "event-user-one@example.com" })
	userTwo := createTestUser(t, db, func(user *domain.User) { user.Email = "event-user-two@example.com" })
	event := createTestEvent(t, db, func(event *domain.Event) { event.Capacity = 3 })
	createTestTicket(t, db, userOne.ID, event.ID, "ACTIVE")
	createTestTicket(t, db, userTwo.ID, event.ID, "ACTIVE")
	createTestTicket(t, db, userTwo.ID, event.ID, "CANCELLED")

	validCategory := "Teatro"
	response, err := services.UpdateEvent(db, event.ID, services.UpdateEventRequest{Category: &validCategory})
	mustNoError(t, err)
	mustEqual(t, "Teatro", response.Category)

	invalidCategory := "Cine"
	_, err = services.UpdateEvent(db, event.ID, services.UpdateEventRequest{Category: &invalidCategory})
	mustErrorIs(t, err, services.ErrInvalidEventCategory)

	capacityTooLow := 1
	_, err = services.UpdateEvent(db, event.ID, services.UpdateEventRequest{Capacity: &capacityTooLow})
	mustErrorIs(t, err, services.ErrCapacityBelowTickets)
}

func TestEventResponseCountsOnlyActiveTicketsForSoldAndAvailability(t *testing.T) {
	db := newTestDB(t)
	userOne := createTestUser(t, db, func(user *domain.User) { user.Email = "sold-one@example.com" })
	userTwo := createTestUser(t, db, func(user *domain.User) { user.Email = "sold-two@example.com" })
	event := createTestEvent(t, db, func(event *domain.Event) { event.Capacity = 4 })
	createTestTicket(t, db, userOne.ID, event.ID, "ACTIVE")
	createTestTicket(t, db, userTwo.ID, event.ID, "CANCELLED")

	response, err := services.GetEventByID(db, event.ID)

	mustNoError(t, err)
	mustEqual(t, 1, response.TicketsSold)
	mustEqual(t, 3, response.AvailableCapacity)
}
