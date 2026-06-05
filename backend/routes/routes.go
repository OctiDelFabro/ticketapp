package routes

import (
	"github.com/OctiDelFabro/ticketapp/backend/controllers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	router.GET("/api/health", controllers.HealthCheck)
}
