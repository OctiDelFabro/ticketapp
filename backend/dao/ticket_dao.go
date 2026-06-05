package dao

import (
	"github.com/OctiDelFabro/ticketapp/backend/domain"
	"gorm.io/gorm"
)

const ActiveTicketStatus = "ACTIVE"

func CreateTicket(db *gorm.DB, ticket *domain.Ticket) error {
	return db.Create(ticket).Error
}

func FindTicketsByUserID(db *gorm.DB, userID uint) ([]domain.Ticket, error) {
	var tickets []domain.Ticket
	if err := db.Preload("Event").Preload("User").Where("user_id = ?", userID).Order("purchase_date DESC").Find(&tickets).Error; err != nil {
		return nil, err
	}

	return tickets, nil
}

func FindTicketByID(db *gorm.DB, ticketID uint) (*domain.Ticket, error) {
	var ticket domain.Ticket
	if err := db.Preload("Event").Preload("User").First(&ticket, ticketID).Error; err != nil {
		return nil, err
	}

	return &ticket, nil
}

func CountActiveTicketsByEventID(db *gorm.DB, eventID uint) (int64, error) {
	var count int64
	if err := db.Model(&domain.Ticket{}).Where("event_id = ? AND status = ?", eventID, ActiveTicketStatus).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func FindActiveTicketByUserAndEvent(db *gorm.DB, userID uint, eventID uint) (*domain.Ticket, error) {
	var ticket domain.Ticket
	if err := db.Preload("Event").Preload("User").Where("user_id = ? AND event_id = ? AND status = ?", userID, eventID, ActiveTicketStatus).First(&ticket).Error; err != nil {
		return nil, err
	}

	return &ticket, nil
}

func UpdateTicket(db *gorm.DB, ticket *domain.Ticket) error {
	return db.Model(ticket).Select("UserID", "EventID", "Status", "PurchaseDate").Updates(ticket).Error
}
