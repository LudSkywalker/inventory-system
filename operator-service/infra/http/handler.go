package http

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (h *HealthHandler) RegisterRoutes(app *fiber.App) {
	app.Get("/health", h.HealthCheck)
}

func (h *HealthHandler) HealthCheck(c *fiber.Ctx) error {
	log.Println("Operator Service: Health check called")
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": "operator-service",
	})
}
