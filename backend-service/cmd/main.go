// Package main Backend Service API
//
// This is the backend service for the inventory management system.
// It provides REST API endpoints to query inventory data.
//
// Terms Of Service:
//
// Schemes: http, https
// Host: localhost:8081
// BasePath: /api/v1
// Version: 1.0.0
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
// @title Backend Service API
// @version 1.0
// @description This is the backend service for the inventory management system.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8081
// @BasePath /api/v1
// @schemes http https

// @externalDocs.description OpenAPI
// @externalDocs.url https://swagger.io/resources/open-api/

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

// @openapi 3.0.0
package main

import (
	"log"
	"os"

	"github.com/LudSkywalker/inventory-system/backend-service/app/service"
	"github.com/LudSkywalker/inventory-system/backend-service/infra/http"
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func main() {
	// Initialize services
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalf("DB_URL environment variable is required")
	}
	log.Printf("Using DB_URL: %s", dbURL)

	inventoryService := service.NewInventoryService(dbURL)
	handler := http.NewInventoryHandler(inventoryService)

	// Initialize Fiber app
	app := fiber.New()

	// Register routes
	handler.RegisterRoutes(app)

	// Swagger routes
	log.Println("Setting up Swagger routes")
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.SendFile("./backend-service/docs/docs.json")
	})
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

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
