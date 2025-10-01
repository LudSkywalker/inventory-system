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
