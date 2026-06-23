package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/OctiDelFabro/ticketapp/backend/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminStatsController struct {
	db *gorm.DB
}

func NewAdminStatsController(db *gorm.DB) *AdminStatsController {
	return &AdminStatsController{db: db}
}

func (controller *AdminStatsController) GetSummary(c *gin.Context) {
	summary, err := services.GetStatsSummary(controller.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, summary)
}

func (controller *AdminStatsController) ListEventStats(c *gin.Context) {
	events, err := services.GetEventStats(controller.db)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, events)
}

func (controller *AdminStatsController) GetEventReport(c *gin.Context) {
	eventID, ok := parseAdminStatsEventID(c)
	if !ok {
		return
	}

	report, err := services.GetEventReport(controller.db, eventID)
	if err != nil {
		handleAdminStatsError(c, err)
		return
	}

	c.JSON(http.StatusOK, report)
}

func parseAdminStatsEventID(c *gin.Context) (uint, bool) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid event id"})
		return 0, false
	}

	return uint(id), true
}

func handleAdminStatsError(c *gin.Context, err error) {
	if errors.Is(err, services.ErrEventNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}
