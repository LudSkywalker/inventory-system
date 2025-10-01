package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/LudSkywalker/inventory-system/operator-service/app/service"
	"github.com/LudSkywalker/inventory-system/operator-service/infra/http"
	"github.com/LudSkywalker/inventory-system/operator-service/infra/kafka"
	"github.com/LudSkywalker/inventory-system/operator-service/infra/sqlite"
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

	// Initialize repositories and services
	repo := sqlite.NewSQLiteRepository(db)
	inventoryService := service.NewGlobalInventoryService(repo)
	handler := http.NewGlobalInventoryHandler(inventoryService)

	// Initialize Kafka consumer
	brokers := strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
	consumer, err := kafka.NewConsumer(brokers, inventoryService)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer consumer.Close()

	// Start Kafka consumer in a goroutine
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := consumer.Start(ctx); err != nil {
			log.Printf("Error starting consumer: %v", err)
		}
	}()

	// Initialize Fiber app
	app := fiber.New()

	// Register routes
	handler.RegisterRoutes(app)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
