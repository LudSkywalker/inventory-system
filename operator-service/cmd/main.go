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
	// Initialize HTTP repository for Kafka consumer
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatalf("DB_URL environment variable is required")
	}

	repo := http.NewHTTPRepository(dbURL)
	inventoryService := service.NewGlobalInventoryService(repo)

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
		log.Println("Starting Kafka consumer...")
		if err := consumer.Start(ctx); err != nil {
			log.Printf("Error starting consumer: %v", err)
		}
	}()

	// Initialize Fiber app with only health check
	app := fiber.New()
	handler := http.NewHealthHandler()
	handler.RegisterRoutes(app)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Operator service starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
