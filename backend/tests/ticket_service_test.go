package tests

import (
	"testing"

	"github.com/OctiDelFabro/ticketapp/backend/domain"
	"github.com/OctiDelFabro/ticketapp/backend/services"
)

func TestPurchaseTicketSuccessCreatesActiveTicket(t *testing.T) {
	db := newTestDB(t)
	user := createTestUser(t, db)
	event := createTestEvent(t, db)

	response, err := services.PurchaseTicket(db, user.ID, services.PurchaseTicketRequest{EventID: event.ID})

	mustNoError(t, err)
	mustNotNil(t, response)
	mustEqual(t, user.ID, response.UserID)
	mustEqual(t, event.ID, response.EventID)
	mustEqual(t, "ACTIVE", response.Status)

	var ticket domain.Ticket
	mustNoError(t, db.First(&ticket, response.ID).Error)
	mustEqual(t, user.ID, ticket.UserID)
	mustEqual(t, event.ID, ticket.EventID)
	mustEqual(t, "ACTIVE", ticket.Status)
}

func TestPurchaseTicketWithoutCapacityReturnsError(t *testing.T) {
	db := newTestDB(t)
	firstUser := createTestUser(t, db, func(user *domain.User) { user.Email = "first@example.com" })
	secondUser := createTestUser(t, db, func(user *domain.User) { user.Email = "second@example.com" })
	event := createTestEvent(t, db, func(event *domain.Event) { event.Capacity = 1 })
	createTestTicket(t, db, firstUser.ID, event.ID, "ACTIVE")

	_, err := services.PurchaseTicket(db, secondUser.ID, services.PurchaseTicketRequest{EventID: event.ID})

	mustErrorIs(t, err, services.ErrNoTicketCapacity)
}

func TestPurchaseTicketDuplicateReturnsError(t *testing.T) {
	db := newTestDB(t)
	user := createTestUser(t, db)
	event := createTestEvent(t, db)
	createTestTicket(t, db, user.ID, event.ID, "ACTIVE")

	_, err := services.PurchaseTicket(db, user.ID, services.PurchaseTicketRequest{EventID: event.ID})

	mustErrorIs(t, err, services.ErrTicketAlreadyExists)
}

func TestGetMyTicketsReturnsOnlyAuthenticatedUserTickets(t *testing.T) {
	db := newTestDB(t)
	user := createTestUser(t, db, func(user *domain.User) { user.Email = "owner@example.com" })
	otherUser := createTestUser(t, db, func(user *domain.User) { user.Email = "other@example.com" })
	eventOne := createTestEvent(t, db, func(event *domain.Event) { event.Title = "Owner Event" })
	eventTwo := createTestEvent(t, db, func(event *domain.Event) { event.Title = "Other Event" })
	owned := createTestTicket(t, db, user.ID, eventOne.ID, "ACTIVE")
	createTestTicket(t, db, otherUser.ID, eventTwo.ID, "ACTIVE")

	tickets, err := services.GetMyTickets(db, user.ID)

	mustNoError(t, err)
	mustLen(t, tickets, 1)
	mustEqual(t, owned.ID, tickets[0].ID)
	mustEqual(t, user.ID, tickets[0].UserID)
}

func TestCancelTicketSuccessChangesStatusToCancelled(t *testing.T) {
	db := newTestDB(t)
	user := createTestUser(t, db)
	event := createTestEvent(t, db)
	ticket := createTestTicket(t, db, user.ID, event.ID, "ACTIVE")

	response, err := services.CancelTicket(db, user.ID, ticket.ID)

	mustNoError(t, err)
	mustNotNil(t, response)
	mustEqual(t, "CANCELLED", response.Status)

	var updated domain.Ticket
	mustNoError(t, db.First(&updated, ticket.ID).Error)
	mustEqual(t, "CANCELLED", updated.Status)
}

func TestCancelTicketOwnedByOtherUserReturnsError(t *testing.T) {
	db := newTestDB(t)
	owner := createTestUser(t, db, func(user *domain.User) { user.Email = "owner@example.com" })
	otherUser := createTestUser(t, db, func(user *domain.User) { user.Email = "other@example.com" })
	event := createTestEvent(t, db)
	ticket := createTestTicket(t, db, owner.ID, event.ID, "ACTIVE")

	_, err := services.CancelTicket(db, otherUser.ID, ticket.ID)

	mustErrorIs(t, err, services.ErrTicketForbidden)
}

func TestTransferTicketSuccessChangesOwnerAndKeepsActive(t *testing.T) {
	db := newTestDB(t)
	owner := createTestUser(t, db, func(user *domain.User) { user.Email = "owner@example.com" })
	target := createTestUser(t, db, func(user *domain.User) { user.Email = "target@example.com" })
	event := createTestEvent(t, db)
	ticket := createTestTicket(t, db, owner.ID, event.ID, "ACTIVE")

	response, err := services.TransferTicket(db, owner.ID, ticket.ID, services.TransferTicketRequest{TargetEmail: target.Email})

	mustNoError(t, err)
	mustNotNil(t, response)
	mustEqual(t, target.ID, response.UserID)
	mustEqual(t, "ACTIVE", response.Status)

	var updated domain.Ticket
	mustNoError(t, db.First(&updated, ticket.ID).Error)
	mustEqual(t, target.ID, updated.UserID)
	mustEqual(t, "ACTIVE", updated.Status)
}

func TestTransferTicketToMissingUserReturnsError(t *testing.T) {
	db := newTestDB(t)
	owner := createTestUser(t, db)
	event := createTestEvent(t, db)
	ticket := createTestTicket(t, db, owner.ID, event.ID, "ACTIVE")

	_, err := services.TransferTicket(db, owner.ID, ticket.ID, services.TransferTicketRequest{TargetEmail: "missing@example.com"})

	mustErrorIs(t, err, services.ErrTargetUserNotFound)
}

func TestTransferTicketToUserWithActiveTicketForSameEventReturnsError(t *testing.T) {
	db := newTestDB(t)
	owner := createTestUser(t, db, func(user *domain.User) { user.Email = "owner@example.com" })
	target := createTestUser(t, db, func(user *domain.User) { user.Email = "target@example.com" })
	event := createTestEvent(t, db, func(event *domain.Event) { event.Capacity = 2 })
	ownedTicket := createTestTicket(t, db, owner.ID, event.ID, "ACTIVE")
	createTestTicket(t, db, target.ID, event.ID, "ACTIVE")

	_, err := services.TransferTicket(db, owner.ID, ownedTicket.ID, services.TransferTicketRequest{TargetEmail: target.Email})

	mustErrorIs(t, err, services.ErrTargetUserAlreadyHasSeat)
}
