// @title Store Service API
// @version 1.0
// @description This is the store service for the inventory management system.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
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
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/LudSkywalker/inventory-system/store-service/app/service"
	"github.com/LudSkywalker/inventory-system/store-service/infra/http"
	"github.com/LudSkywalker/inventory-system/store-service/infra/kafka"
	"github.com/LudSkywalker/inventory-system/store-service/infra/sqlite"
	"github.com/gofiber/fiber/v2"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func main() {
	// Initialize SQLite
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := sqlite.InitDB(db); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize Kafka producer
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	producer, err := kafka.NewProducer(brokers)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// Initialize repositories and services
	inventoryRepo := sqlite.NewSQLiteRepository(db)
	inventoryService := service.NewInventoryService(inventoryRepo, producer)
	inventoryHandler := http.NewInventoryHandler(inventoryService)

	// Initialize Fiber app
	app := fiber.New()

	// Register routes
	inventoryHandler.RegisterRoutes(app)

	// Swagger routes
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.SendFile("./store-service/docs/docs.json")
	})
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
