package http

import (
	"github.com/LudSkywalker/inventory-system/backend-service/app/service"
	"github.com/gofiber/fiber/v2"
)

type InventoryHandler struct {
	service *service.InventoryService
}

func NewInventoryHandler(service *service.InventoryService) *InventoryHandler {
	return &InventoryHandler{service: service}
}

func (h *InventoryHandler) RegisterRoutes(app *fiber.App) {
	group := app.Group("/api/v1")

	group.Get("/inventory", h.GetGlobalInventory)
	group.Get("/inventory/store/:storeId", h.GetStoreInventory)
	group.Get("/inventory/item/:itemId", h.GetItemInventory)
}

func (h *InventoryHandler) GetGlobalInventory(c *fiber.Ctx) error {
	inventories, err := h.service.GetGlobalInventory(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(inventories)
}

func (h *InventoryHandler) GetStoreInventory(c *fiber.Ctx) error {
	storeID := c.Params("storeId")
	inventories, err := h.service.GetStoreInventory(c.Context(), storeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(inventories)
}

func (h *InventoryHandler) GetItemInventory(c *fiber.Ctx) error {
	itemID := c.Params("itemId")
	inventories, err := h.service.GetItemInventory(c.Context(), itemID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(inventories)
}
