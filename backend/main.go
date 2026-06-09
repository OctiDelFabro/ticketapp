package main

import (
	"log"
	"os"

	"github.com/OctiDelFabro/ticketapp/backend/config"
	"github.com/OctiDelFabro/ticketapp/backend/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("warning: could not load .env file: %v", err)
	}

	db, err := config.ConnectDatabase()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := config.RunMigrations(db); err != nil {
		log.Fatalf("failed to run database migrations: %v", err)
	}

	if err := config.SeedEvents(db); err != nil {
		log.Fatalf("failed to seed initial events: %v", err)
	}

	if err := config.SeedDemoData(db); err != nil {
		log.Fatalf("failed to seed demo data: %v", err)
	}

	router := gin.Default()
	router.Use(corsMiddleware())
	routes.SetupRoutes(router, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PATCH, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
