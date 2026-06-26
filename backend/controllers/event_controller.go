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

type EventController struct {
	db *gorm.DB
}

func NewEventController(db *gorm.DB) *EventController {
	return &EventController{db: db}
}

func (controller *EventController) ListEvents(c *gin.Context) {
	filters, err := parseEventFilters(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid available_only value"})
		return
	}

	events, err := services.ListEvents(controller.db, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (controller *EventController) GetEventByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return
	}

	event, err := services.GetEventByID(controller.db, uint(id))
	if err != nil {
		if errors.Is(err, services.ErrEventNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "event not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, event)
}

func parseEventFilters(c *gin.Context) (dao.EventFilters, error) {
	filters := dao.EventFilters{
		Search:   c.Query("search"),
		Category: c.Query("category"),
	}

	availableOnly := c.Query("available_only")
	if availableOnly == "" {
		return filters, nil
	}

	parsedAvailableOnly, err := strconv.ParseBool(availableOnly)
	if err != nil {
		return filters, err
	}

	filters.AvailableOnly = parsedAvailableOnly
	return filters, nil
}
