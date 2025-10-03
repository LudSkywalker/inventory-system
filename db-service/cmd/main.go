package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/LudSkywalker/inventory-system/db-service/infra/http"
	"github.com/LudSkywalker/inventory-system/db-service/infra/sqlite"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Initialize SQLite in memory
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := sqlite.InitDB(db); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Initialize repository and handler
	repo := sqlite.NewSQLiteRepository(db)
	handler := http.NewGlobalInventoryHandler(repo)

	// Initialize Fiber app
	app := fiber.New()

	// Register routes
	handler.RegisterRoutes(app)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8083"
	}

	log.Printf("DB Service starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}