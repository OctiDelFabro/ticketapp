package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/OctiDelFabro/ticketapp/backend/dao"
	"github.com/OctiDelFabro/ticketapp/backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminEventController struct {
	db *gorm.DB
}

func NewAdminEventController(db *gorm.DB) *AdminEventController {
	return &AdminEventController{db: db}
}

func (controller *AdminEventController) ListEvents(c *gin.Context) {
	events, err := services.ListEvents(controller.db, dao.EventFilters{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (controller *AdminEventController) CreateEvent(c *gin.Context) {
	var req services.CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	event, err := services.CreateEvent(controller.db, req)
	if err != nil {
		handleAdminEventError(c, err)
		return
	}

	c.JSON(http.StatusCreated, event)
}

func (controller *AdminEventController) UpdateEvent(c *gin.Context) {
	id, ok := parseAdminEventID(c)
	if !ok {
		return
	}

	var req services.UpdateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	event, err := services.UpdateEvent(controller.db, id, req)
	if err != nil {
		handleAdminEventError(c, err)
		return
	}

	c.JSON(http.StatusOK, event)
}

func (controller *AdminEventController) DeleteEvent(c *gin.Context) {
	id, ok := parseAdminEventID(c)
	if !ok {
		return
	}

	if err := services.DisableEvent(controller.db, id); err != nil {
		handleAdminEventError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "event disabled successfully"})
}

func parseAdminEventID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return 0, false
	}

	return uint(id), true
}

func handleAdminEventError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, services.ErrEventNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, services.ErrInvalidEventRequest), errors.Is(err, services.ErrCapacityBelowTickets):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
