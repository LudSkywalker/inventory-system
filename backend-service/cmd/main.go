package main

import (
	"log"
	"os"

	"github.com/LudSkywalker/inventory-system/backend-service/app/service"
	"github.com/LudSkywalker/inventory-system/backend-service/infra/http"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Initialize services
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalf("DB_URL environment variable is required")
	}

	inventoryService := service.NewInventoryService(dbURL)
	handler := http.NewInventoryHandler(inventoryService)

	// Initialize Fiber app
	app := fiber.New()

	// Register routes
	handler.RegisterRoutes(app)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
