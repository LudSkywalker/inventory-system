package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/LudSkywalker/inventory-system/internal/inventory/localstore/app/service"
	"github.com/LudSkywalker/inventory-system/internal/inventory/localstore/infra/http"
	"github.com/LudSkywalker/inventory-system/internal/inventory/localstore/infra/sqlite"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Initialize SQLite in-memory database
	db, err := sql.Open("sqlite3", "file::memory:?cache=shared")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := sqlite.InitDB(db); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize repository and service
	repo := sqlite.NewSQLiteRepository(db)
	inventoryService := service.NewInventoryService(repo)

	// Initialize HTTP handler
	handler := http.NewInventoryHandler(inventoryService)

	// Initialize Fiber app with longer timeout for shared memory access
	app := fiber.New(fiber.Config{
		ReadTimeout:  60,
		WriteTimeout: 60,
	})

	// Register routes
	handler.RegisterRoutes(app)

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
