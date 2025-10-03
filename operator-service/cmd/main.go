package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/LudSkywalker/inventory-system/operator-service/app/service"
	"github.com/LudSkywalker/inventory-system/operator-service/infra/http"
	"github.com/LudSkywalker/inventory-system/operator-service/infra/kafka"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Initialize HTTP repository
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalf("DB_URL environment variable is required")
	}

	repo := http.NewHTTPRepository(dbURL)
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
