package services

import (
	"errors"
	"strings"
	"time"

	"github.com/OctiDelFabro/ticketapp/backend/dao"
	"github.com/OctiDelFabro/ticketapp/backend/domain"
	"gorm.io/gorm"
)

const (
	ticketStatusActive    = "ACTIVE"
	ticketStatusCancelled = "CANCELLED"
)

var (
	ErrTicketNotFound           = errors.New("ticket not found")
	ErrTicketForbidden          = errors.New("ticket does not belong to authenticated user")
	ErrTicketAlreadyExists      = errors.New("user already has an active ticket for this event")
	ErrNoTicketCapacity         = errors.New("no capacity available for this event")
	ErrInvalidTicketQuantity    = errors.New("ticket quantity must be greater than zero")
	ErrTicketNotActive          = errors.New("ticket is not active")
	ErrTargetUserNotFound       = errors.New("target user not found")
	ErrTransferToSameUser       = errors.New("target user must be different from authenticated user")
	ErrGiftToSelf               = errors.New("target user must be different from authenticated user")
	ErrGiftMessageTooLong       = errors.New("gift message must be 250 characters or fewer")
	ErrTargetUserAlreadyHasSeat = errors.New("target user already has an active ticket for this event")
)

type PurchaseTicketRequest struct {
	EventID  uint `json:"event_id"`
	Quantity *int `json:"quantity,omitempty"`
}

type PurchaseTicketsResponse struct {
	Tickets  []TicketResponse `json:"tickets"`
	Quantity int              `json:"quantity"`
}

type TransferTicketRequest struct {
	TargetEmail string `json:"target_email"`
}

type GiftTicketRequest struct {
	EventID     uint   `json:"event_id"`
	TargetEmail string `json:"target_email"`
	GiftMessage string `json:"message"`
}

type TicketResponse struct {
	ID             uint       `json:"id"`
	EventID        uint       `json:"event_id"`
	EventTitle     string     `json:"event_title"`
	EventImageURL  string     `json:"image_url"`
	EventStartDate time.Time  `json:"event_start_date"`
	EventLocation  string     `json:"event_location"`
	EventPrice     float64    `json:"event_price"`
	Status         string     `json:"status"`
	PurchaseDate   time.Time  `json:"purchase_date"`
	UserID         uint       `json:"user_id"`
	UserEmail      string     `json:"user_email"`
	IsGift         bool       `json:"is_gift"`
	GiftedByID     *uint      `json:"gifted_by_id,omitempty"`
	GiftedByEmail  string     `json:"gifted_by_email"`
	GiftMessage    string     `json:"gift_message"`
	GiftedAt       *time.Time `json:"gifted_at,omitempty"`
}

func PurchaseTicket(db *gorm.DB, userID uint, req PurchaseTicketRequest) (*TicketResponse, error) {
	response, err := PurchaseTickets(db, userID, req)
	if err != nil {
		return nil, err
	}

	return &response.Tickets[0], nil
}

func PurchaseTickets(db *gorm.DB, userID uint, req PurchaseTicketRequest) (*PurchaseTicketsResponse, error) {
	if req.EventID == 0 {
		return nil, ErrInvalidRequest
	}

	quantity := 1
	if req.Quantity != nil {
		quantity = *req.Quantity
	}
	if quantity <= 0 {
		return nil, ErrInvalidTicketQuantity
	}

	var createdTickets []TicketResponse
	err := db.Transaction(func(tx *gorm.DB) error {
		event, err := dao.FindEventByID(tx, req.EventID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrEventNotFound
			}
			return err
		}

		activeTickets, err := dao.CountActiveTicketsByEventID(tx, req.EventID)
		if err != nil {
			return err
		}

		if event.Capacity-int(activeTickets) < quantity {
			return ErrNoTicketCapacity
		}

		createdTickets = make([]TicketResponse, 0, quantity)
		purchaseDate := time.Now()
		for range quantity {
			ticket := domain.Ticket{
				UserID:       userID,
				EventID:      req.EventID,
				Status:       ticketStatusActive,
				PurchaseDate: purchaseDate,
			}

			if err := dao.CreateTicket(tx, &ticket); err != nil {
				return err
			}

			createdTicket, err := dao.FindTicketByID(tx, ticket.ID)
			if err != nil {
				return err
			}

			createdTickets = append(createdTickets, buildTicketResponse(*createdTicket))
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &PurchaseTicketsResponse{Tickets: createdTickets, Quantity: quantity}, nil
}

func GiftTicket(db *gorm.DB, giverUserID uint, req GiftTicketRequest) (*TicketResponse, error) {
	req.TargetEmail = strings.TrimSpace(req.TargetEmail)
	req.GiftMessage = strings.TrimSpace(req.GiftMessage)
	if req.EventID == 0 || req.TargetEmail == "" {
		return nil, ErrInvalidRequest
	}
	if len(req.GiftMessage) > 250 {
		return nil, ErrGiftMessageTooLong
	}

	var response TicketResponse
	err := db.Transaction(func(tx *gorm.DB) error {
		targetUser, err := dao.FindUserByEmail(tx, req.TargetEmail)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrTargetUserNotFound
			}
			return err
		}
		if targetUser.ID == giverUserID {
			return ErrGiftToSelf
		}

		event, err := dao.FindEventByID(tx, req.EventID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrEventNotFound
			}
			return err
		}

		activeTickets, err := dao.CountActiveTicketsByEventID(tx, req.EventID)
		if err != nil {
			return err
		}
		if event.Capacity-int(activeTickets) < 1 {
			return ErrNoTicketCapacity
		}

		if _, err := dao.FindActiveTicketByUserAndEvent(tx, targetUser.ID, req.EventID); err == nil {
			return ErrTargetUserAlreadyHasSeat
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		now := time.Now()
		ticket := domain.Ticket{
			UserID:       targetUser.ID,
			EventID:      req.EventID,
			Status:       ticketStatusActive,
			PurchaseDate: now,
			GiftedByID:   &giverUserID,
			GiftMessage:  req.GiftMessage,
			GiftedAt:     &now,
		}
		if err := dao.CreateTicket(tx, &ticket); err != nil {
			return err
		}
		createdTicket, err := dao.FindTicketByID(tx, ticket.ID)
		if err != nil {
			return err
		}
		response = buildTicketResponse(*createdTicket)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func GetMyTickets(db *gorm.DB, userID uint) ([]TicketResponse, error) {
	tickets, err := dao.FindTicketsByUserID(db, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]TicketResponse, 0, len(tickets))
	for _, ticket := range tickets {
		responses = append(responses, buildTicketResponse(ticket))
	}

	return responses, nil
}

func CancelTicket(db *gorm.DB, userID uint, ticketID uint) (*TicketResponse, error) {
	ticket, err := findOwnedActiveTicket(db, userID, ticketID)
	if err != nil {
		return nil, err
	}

	ticket.Status = ticketStatusCancelled
	if err := dao.UpdateTicket(db, ticket); err != nil {
		return nil, err
	}

	updatedTicket, err := dao.FindTicketByID(db, ticket.ID)
	if err != nil {
		return nil, err
	}

	response := buildTicketResponse(*updatedTicket)
	return &response, nil
}

func TransferTicket(db *gorm.DB, userID uint, ticketID uint, req TransferTicketRequest) (*TicketResponse, error) {
	req.TargetEmail = strings.TrimSpace(req.TargetEmail)
	if req.TargetEmail == "" {
		return nil, ErrInvalidRequest
	}

	ticket, err := findOwnedActiveTicket(db, userID, ticketID)
	if err != nil {
		return nil, err
	}

	targetUser, err := dao.FindUserByEmail(db, req.TargetEmail)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTargetUserNotFound
		}
		return nil, err
	}

	if targetUser.ID == userID {
		return nil, ErrTransferToSameUser
	}

	ticket.UserID = targetUser.ID
	ticket.User = *targetUser
	if err := dao.UpdateTicket(db, ticket); err != nil {
		return nil, err
	}

	updatedTicket, err := dao.FindTicketByID(db, ticket.ID)
	if err != nil {
		return nil, err
	}

	response := buildTicketResponse(*updatedTicket)
	return &response, nil
}

func findOwnedActiveTicket(db *gorm.DB, userID uint, ticketID uint) (*domain.Ticket, error) {
	ticket, err := dao.FindTicketByID(db, ticketID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrTicketNotFound
		}
		return nil, err
	}

	if ticket.UserID != userID {
		return nil, ErrTicketForbidden
	}

	if ticket.Status != ticketStatusActive {
		return nil, ErrTicketNotActive
	}

	return ticket, nil
}

func buildTicketResponse(ticket domain.Ticket) TicketResponse {
	return TicketResponse{
		ID:             ticket.ID,
		EventID:        ticket.EventID,
		EventTitle:     ticket.Event.Title,
		EventImageURL:  ticket.Event.ImageURL,
		EventStartDate: ticket.Event.StartDate,
		EventLocation:  ticket.Event.Location,
		EventPrice:     ticket.Event.Price,
		Status:         ticket.Status,
		PurchaseDate:   ticket.PurchaseDate,
		UserID:         ticket.UserID,
		UserEmail:      ticket.User.Email,
		IsGift:         ticket.GiftedByID != nil,
		GiftedByID:     ticket.GiftedByID,
		GiftedByEmail:  giftedByEmail(ticket),
		GiftMessage:    ticket.GiftMessage,
		GiftedAt:       ticket.GiftedAt,
	}
}

func giftedByEmail(ticket domain.Ticket) string {
	if ticket.GiftedBy == nil {
		return ""
	}
	return ticket.GiftedBy.Email
}
