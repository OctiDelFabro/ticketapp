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

func TestPurchaseTicketsWithQuantityAboveCapacityReturnsError(t *testing.T) {
	db := newTestDB(t)
	user := createTestUser(t, db)
	event := createTestEvent(t, db, func(event *domain.Event) { event.Capacity = 1 })
	quantity := 2

	_, err := services.PurchaseTickets(db, user.ID, services.PurchaseTicketRequest{EventID: event.ID, Quantity: &quantity})

	mustErrorIs(t, err, services.ErrNoTicketCapacity)
}

func TestPurchaseTicketsAllowsSameUserToBuyMultipleTicketsForSameEvent(t *testing.T) {
	db := newTestDB(t)
	user := createTestUser(t, db)
	event := createTestEvent(t, db, func(event *domain.Event) { event.Capacity = 3 })
	quantity := 2

	response, err := services.PurchaseTickets(db, user.ID, services.PurchaseTicketRequest{EventID: event.ID, Quantity: &quantity})

	mustNoError(t, err)
	mustNotNil(t, response)
	mustEqual(t, 2, response.Quantity)
	mustLen(t, response.Tickets, 2)
	mustNotEqual(t, response.Tickets[0].ID, response.Tickets[1].ID)

	var count int64
	mustNoError(t, db.Model(&domain.Ticket{}).Where("user_id = ? AND event_id = ? AND status = ?", user.ID, event.ID, "ACTIVE").Count(&count).Error)
	mustEqual(t, int64(2), count)
}

func TestPurchaseTicketsWithInvalidQuantityReturnsError(t *testing.T) {
	db := newTestDB(t)
	user := createTestUser(t, db)
	event := createTestEvent(t, db)
	quantity := 0

	_, err := services.PurchaseTickets(db, user.ID, services.PurchaseTicketRequest{EventID: event.ID, Quantity: &quantity})

	mustErrorIs(t, err, services.ErrInvalidTicketQuantity)
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

func TestTransferTicketToUserWithActiveTicketForSameEventSucceeds(t *testing.T) {
	db := newTestDB(t)
	owner := createTestUser(t, db, func(user *domain.User) { user.Email = "owner@example.com" })
	target := createTestUser(t, db, func(user *domain.User) { user.Email = "target@example.com" })
	event := createTestEvent(t, db, func(event *domain.Event) { event.Capacity = 2 })
	ownedTicket := createTestTicket(t, db, owner.ID, event.ID, "ACTIVE")
	createTestTicket(t, db, target.ID, event.ID, "ACTIVE")

	response, err := services.TransferTicket(db, owner.ID, ownedTicket.ID, services.TransferTicketRequest{TargetEmail: target.Email})

	mustNoError(t, err)
	mustEqual(t, target.ID, response.UserID)
}

func TestPurchaseTicketWithZeroEventIDReturnsError(t *testing.T) {
	db := newTestDB(t)
	user := createTestUser(t, db)

	_, err := services.PurchaseTicket(db, user.ID, services.PurchaseTicketRequest{EventID: 0})

	mustErrorIs(t, err, services.ErrInvalidRequest)
}

func TestPurchaseTicketForMissingEventReturnsError(t *testing.T) {
	db := newTestDB(t)
	user := createTestUser(t, db)

	_, err := services.PurchaseTicket(db, user.ID, services.PurchaseTicketRequest{EventID: 999})

	mustErrorIs(t, err, services.ErrEventNotFound)
}

func TestPurchaseTicketsDoesNotCountCancelledTicketsAsUsedCapacity(t *testing.T) {
	db := newTestDB(t)
	firstUser := createTestUser(t, db, func(user *domain.User) { user.Email = "cancelled-owner@example.com" })
	buyer := createTestUser(t, db, func(user *domain.User) { user.Email = "buyer@example.com" })
	event := createTestEvent(t, db, func(event *domain.Event) { event.Capacity = 1 })
	createTestTicket(t, db, firstUser.ID, event.ID, "CANCELLED")

	response, err := services.PurchaseTicket(db, buyer.ID, services.PurchaseTicketRequest{EventID: event.ID})

	mustNoError(t, err)
	mustEqual(t, buyer.ID, response.UserID)
	mustEqual(t, "ACTIVE", response.Status)
}

func TestGiftTicketSuccessStoresRecipientMetadataAndTrimmedMessage(t *testing.T) {
	db := newTestDB(t)
	giver := createTestUser(t, db, func(user *domain.User) { user.Email = "giver@example.com" })
	target := createTestUser(t, db, func(user *domain.User) { user.Email = "target-gift@example.com" })
	event := createTestEvent(t, db)

	response, err := services.GiftTicket(db, giver.ID, services.GiftTicketRequest{EventID: event.ID, TargetEmail: "  " + target.Email + "  ", GiftMessage: "  Felicidades!  "})

	mustNoError(t, err)
	mustNotNil(t, response)
	mustEqual(t, target.ID, response.UserID)
	mustTrue(t, response.IsGift)
	mustNotNil(t, response.GiftedByID)
	mustEqual(t, giver.ID, *response.GiftedByID)
	mustEqual(t, giver.Email, response.GiftedByEmail)
	mustEqual(t, "Felicidades!", response.GiftMessage)

	var ticket domain.Ticket
	mustNoError(t, db.First(&ticket, response.ID).Error)
	mustEqual(t, target.ID, ticket.UserID)
	mustNotNil(t, ticket.GiftedByID)
	mustEqual(t, giver.ID, *ticket.GiftedByID)
	mustEqual(t, "Felicidades!", ticket.GiftMessage)
}

func TestGiftTicketValidationErrors(t *testing.T) {
	tests := []struct {
		name string
		req  services.GiftTicketRequest
		want error
	}{
		{"empty target email", services.GiftTicketRequest{EventID: 1, TargetEmail: "  "}, services.ErrInvalidRequest},
		{"zero event", services.GiftTicketRequest{TargetEmail: "target@example.com"}, services.ErrInvalidRequest},
		{"long message", services.GiftTicketRequest{EventID: 1, TargetEmail: "target@example.com", GiftMessage: string(make([]byte, 251))}, services.ErrGiftMessageTooLong},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			db := newTestDB(t)
			giver := createTestUser(t, db)
			_, err := services.GiftTicket(db, giver.ID, tc.req)
			mustErrorIs(t, err, tc.want)
		})
	}
}

func TestGiftTicketDomainErrors(t *testing.T) {
	t.Run("missing target email", func(t *testing.T) {
		db := newTestDB(t)
		giver := createTestUser(t, db)
		event := createTestEvent(t, db)
		_, err := services.GiftTicket(db, giver.ID, services.GiftTicketRequest{EventID: event.ID, TargetEmail: "missing@example.com"})
		mustErrorIs(t, err, services.ErrTargetUserNotFound)
	})
	t.Run("gift to self", func(t *testing.T) {
		db := newTestDB(t)
		giver := createTestUser(t, db, func(user *domain.User) { user.Email = "self@example.com" })
		event := createTestEvent(t, db)
		_, err := services.GiftTicket(db, giver.ID, services.GiftTicketRequest{EventID: event.ID, TargetEmail: giver.Email})
		mustErrorIs(t, err, services.ErrGiftToSelf)
	})
	t.Run("no capacity", func(t *testing.T) {
		db := newTestDB(t)
		giver := createTestUser(t, db, func(user *domain.User) { user.Email = "giver-cap@example.com" })
		target := createTestUser(t, db, func(user *domain.User) { user.Email = "target-cap@example.com" })
		occupant := createTestUser(t, db, func(user *domain.User) { user.Email = "occupant@example.com" })
		event := createTestEvent(t, db, func(event *domain.Event) { event.Capacity = 1 })
		createTestTicket(t, db, occupant.ID, event.ID, "ACTIVE")
		_, err := services.GiftTicket(db, giver.ID, services.GiftTicketRequest{EventID: event.ID, TargetEmail: target.Email})
		mustErrorIs(t, err, services.ErrNoTicketCapacity)
	})
	t.Run("target already has active ticket", func(t *testing.T) {
		db := newTestDB(t)
		giver := createTestUser(t, db, func(user *domain.User) { user.Email = "giver-seat@example.com" })
		target := createTestUser(t, db, func(user *domain.User) { user.Email = "target-seat@example.com" })
		event := createTestEvent(t, db, func(event *domain.Event) { event.Capacity = 2 })
		createTestTicket(t, db, target.ID, event.ID, "ACTIVE")
		_, err := services.GiftTicket(db, giver.ID, services.GiftTicketRequest{EventID: event.ID, TargetEmail: target.Email})
		mustErrorIs(t, err, services.ErrTargetUserAlreadyHasSeat)
	})
}

func TestCancelTicketMissingAndInactiveErrors(t *testing.T) {
	db := newTestDB(t)
	user := createTestUser(t, db)
	event := createTestEvent(t, db)
	cancelled := createTestTicket(t, db, user.ID, event.ID, "CANCELLED")

	_, missingErr := services.CancelTicket(db, user.ID, 999)
	mustErrorIs(t, missingErr, services.ErrTicketNotFound)

	_, inactiveErr := services.CancelTicket(db, user.ID, cancelled.ID)
	mustErrorIs(t, inactiveErr, services.ErrTicketNotActive)
}

func TestTransferTicketValidationAndInactiveErrors(t *testing.T) {
	db := newTestDB(t)
	owner := createTestUser(t, db, func(user *domain.User) { user.Email = "transfer-owner@example.com" })
	other := createTestUser(t, db, func(user *domain.User) { user.Email = "transfer-other@example.com" })
	target := createTestUser(t, db, func(user *domain.User) { user.Email = "transfer-target@example.com" })
	event := createTestEvent(t, db)
	active := createTestTicket(t, db, owner.ID, event.ID, "ACTIVE")
	cancelled := createTestTicket(t, db, owner.ID, event.ID, "CANCELLED")

	_, err := services.TransferTicket(db, other.ID, active.ID, services.TransferTicketRequest{TargetEmail: target.Email})
	mustErrorIs(t, err, services.ErrTicketForbidden)
	_, err = services.TransferTicket(db, owner.ID, active.ID, services.TransferTicketRequest{TargetEmail: owner.Email})
	mustErrorIs(t, err, services.ErrTransferToSameUser)
	_, err = services.TransferTicket(db, owner.ID, 999, services.TransferTicketRequest{TargetEmail: target.Email})
	mustErrorIs(t, err, services.ErrTicketNotFound)
	_, err = services.TransferTicket(db, owner.ID, cancelled.ID, services.TransferTicketRequest{TargetEmail: target.Email})
	mustErrorIs(t, err, services.ErrTicketNotActive)
}

func TestTransferTicketKeepsGiftMetadata(t *testing.T) {
	db := newTestDB(t)
	giver := createTestUser(t, db, func(user *domain.User) { user.Email = "gift-giver@example.com" })
	owner := createTestUser(t, db, func(user *domain.User) { user.Email = "gift-owner@example.com" })
	target := createTestUser(t, db, func(user *domain.User) { user.Email = "gift-new-owner@example.com" })
	event := createTestEvent(t, db)
	now := createTestTicket(t, db, owner.ID, event.ID, "ACTIVE")
	now.GiftedByID = &giver.ID
	now.GiftMessage = "Enjoy"
	mustNoError(t, db.Save(&now).Error)

	response, err := services.TransferTicket(db, owner.ID, now.ID, services.TransferTicketRequest{TargetEmail: target.Email})

	mustNoError(t, err)
	mustEqual(t, target.ID, response.UserID)
	mustTrue(t, response.IsGift)
	mustNotNil(t, response.GiftedByID)
	mustEqual(t, giver.ID, *response.GiftedByID)
	mustEqual(t, "Enjoy", response.GiftMessage)
}

func TestGetMyTicketsGiftResponseFieldsAndNilGiftedBy(t *testing.T) {
	db := newTestDB(t)
	user := createTestUser(t, db, func(user *domain.User) { user.Email = "response-owner@example.com" })
	giver := createTestUser(t, db, func(user *domain.User) { user.Email = "response-giver@example.com" })
	event := createTestEvent(t, db)
	normal := createTestTicket(t, db, user.ID, event.ID, "ACTIVE")
	gift := createTestTicket(t, db, user.ID, event.ID, "ACTIVE")
	gift.GiftedByID = &giver.ID
	gift.GiftMessage = "Happy day"
	mustNoError(t, db.Save(&gift).Error)

	responses, err := services.GetMyTickets(db, user.ID)

	mustNoError(t, err)
	byID := map[uint]services.TicketResponse{}
	for _, response := range responses {
		byID[response.ID] = response
	}
	mustEqual(t, false, byID[normal.ID].IsGift)
	mustTrue(t, byID[gift.ID].IsGift)
	mustEqual(t, giver.Email, byID[gift.ID].GiftedByEmail)
	mustEqual(t, "Happy day", byID[gift.ID].GiftMessage)
}
