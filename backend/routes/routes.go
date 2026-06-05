package routes

import (
	"github.com/OctiDelFabro/ticketapp/backend/controllers"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	authController := controllers.NewAuthController(db)

	router.GET("/api/health", controllers.HealthCheck)
	router.POST("/api/auth/register", authController.Register)
	router.POST("/api/auth/login", authController.Login)
}
