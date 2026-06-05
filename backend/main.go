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

	router := gin.Default()
	routes.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
