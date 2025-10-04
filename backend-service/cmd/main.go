// Package main Backend Service API
//
// This is the backend service for the inventory management system.
// It provides REST API endpoints to query inventory data.
//
// Terms Of Service:
//
// Schemes: http, https
// Host: localhost:8081
// BasePath: /
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
// @BasePath /
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
	"github.com/gofiber/fiber/v2/middleware/cors"
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

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000,http://127.0.0.1:3000,http://localhost:8080,http://127.0.0.1:8080,http://localhost:8081,http://127.0.0.1:8081,http://localhost:8082,http://127.0.0.1:8082,http://localhost:8083,http://127.0.0.1:8083",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

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
