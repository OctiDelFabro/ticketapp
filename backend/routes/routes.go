package routes

import (
	"github.com/OctiDelFabro/ticketapp/backend/controllers"
	"github.com/OctiDelFabro/ticketapp/backend/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	authController := controllers.NewAuthController(db)
	eventController := controllers.NewEventController(db)
	ticketController := controllers.NewTicketController(db)

	router.GET("/api/health", controllers.HealthCheck)
	router.POST("/api/auth/register", authController.Register)
	router.POST("/api/auth/login", authController.Login)
	router.GET("/api/events", eventController.ListEvents)
	router.GET("/api/events/:id", eventController.GetEventByID)

	tickets := router.Group("/api/tickets")
	tickets.Use(middlewares.AuthMiddleware())
	tickets.POST("/purchase", ticketController.Purchase)
	tickets.GET("/me", ticketController.GetMyTickets)
	tickets.PATCH("/:id/cancel", ticketController.Cancel)
	tickets.PATCH("/:id/transfer", ticketController.Transfer)
}
