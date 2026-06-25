package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/OctiDelFabro/ticketapp/backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TicketController struct {
	db *gorm.DB
}

func NewTicketController(db *gorm.DB) *TicketController {
	return &TicketController{db: db}
}

func (controller *TicketController) Purchase(c *gin.Context) {
	userID, ok := getAuthenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authenticated user not found"})
		return
	}

	var req services.PurchaseTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	response, err := services.PurchaseTickets(controller.db, userID, req)
	if err != nil {
		handleTicketError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (controller *TicketController) Gift(c *gin.Context) {
	userID, ok := getAuthenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authenticated user not found"})
		return
	}

	var req services.GiftTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	response, err := services.GiftTicket(controller.db, userID, req)
	if err != nil {
		handleTicketError(c, err)
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (controller *TicketController) GetMyTickets(c *gin.Context) {
	userID, ok := getAuthenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authenticated user not found"})
		return
	}

	response, err := services.GetMyTickets(controller.db, userID)
	if err != nil {
		handleTicketError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (controller *TicketController) Cancel(c *gin.Context) {
	userID, ok := getAuthenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authenticated user not found"})
		return
	}

	ticketID, err := parseTicketID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket id"})
		return
	}

	response, err := services.CancelTicket(controller.db, userID, ticketID)
	if err != nil {
		handleTicketError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func (controller *TicketController) Transfer(c *gin.Context) {
	userID, ok := getAuthenticatedUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "authenticated user not found"})
		return
	}

	ticketID, err := parseTicketID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ticket id"})
		return
	}

	var req services.TransferTicketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	response, err := services.TransferTicket(controller.db, userID, ticketID, req)
	if err != nil {
		handleTicketError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}

func getAuthenticatedUserID(c *gin.Context) (uint, bool) {
	value, exists := c.Get("userID")
	if !exists {
		return 0, false
	}

	userID, ok := value.(uint)
	return userID, ok
}

func parseTicketID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return 0, err
	}
	if id == 0 {
		return 0, strconv.ErrSyntax
	}

	return uint(id), nil
}

func handleTicketError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrInvalidRequest), errors.Is(err, services.ErrInvalidTicketQuantity), errors.Is(err, services.ErrGiftMessageTooLong), errors.Is(err, services.ErrGiftToSelf):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, services.ErrTicketForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, services.ErrTicketNotFound), errors.Is(err, services.ErrEventNotFound), errors.Is(err, services.ErrTargetUserNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, services.ErrNoTicketCapacity), errors.Is(err, services.ErrTicketAlreadyExists), errors.Is(err, services.ErrTicketNotActive), errors.Is(err, services.ErrTransferToSameUser), errors.Is(err, services.ErrTargetUserAlreadyHasSeat):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
