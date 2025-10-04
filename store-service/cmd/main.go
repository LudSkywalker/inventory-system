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
	"context"
	"database/sql"
	"log"
	"os"
	"strings"
	"time"

	"github.com/LudSkywalker/inventory-system/store-service/app/service"
	"github.com/LudSkywalker/inventory-system/store-service/infra/http"
	"github.com/LudSkywalker/inventory-system/store-service/infra/kafka"
	"github.com/LudSkywalker/inventory-system/store-service/infra/sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func main() {
	// Initialize SQLite
	dbMode := os.Getenv("DB_MODE")
	if dbMode == "" {
		dbMode = "memory"
	}
	log.Printf("Using DB_MODE: %s", dbMode)

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

	// CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000,http://127.0.0.1:3000,http://localhost:8080,http://127.0.0.1:8080,http://localhost:8081,http://127.0.0.1:8081,http://localhost:8082,http://127.0.0.1:8082,http://localhost:8083,http://127.0.0.1:8083",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	// Register routes
	inventoryHandler.RegisterRoutes(app)

	// Swagger routes
	log.Println("Setting up Swagger routes")
	app.Get("/swagger/doc.json", func(c *fiber.Ctx) error {
		return c.SendFile("./docs/swagger.json")
	})
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)

	// Start periodic sync goroutine
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				log.Println("Starting periodic inventory sync")
				ctx := context.Background()
				if err := inventoryService.SyncAllInventories(ctx); err != nil {
					log.Printf("Error during periodic sync: %v", err)
				} else {
					log.Println("Periodic inventory sync completed")
				}
			}
		}
	}()

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
