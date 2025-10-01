package http

import (
	"github.com/LudSkywalker/inventory-system/internal/inventory/query/app/dto"
	"github.com/LudSkywalker/inventory-system/internal/inventory/query/app/port"
	"github.com/gofiber/fiber/v2"
)

type InventoryQueryHandler struct {
	useCase port.InventoryQueryUseCase
}

func NewInventoryQueryHandler(useCase port.InventoryQueryUseCase) *InventoryQueryHandler {
	return &InventoryQueryHandler{useCase: useCase}
}

func (h *InventoryQueryHandler) RegisterRoutes(app *fiber.App) {
	group := app.Group("/api/v1")

	group.Get("/inventory", h.GetAllInventories)
	group.Get("/inventory/store/:storeId", h.GetStoreInventory)
	group.Get("/inventory/item/:itemId", h.GetItemInventory)
	group.Get("/inventory/:storeId/:itemId", h.GetInventory)
}

func (h *InventoryQueryHandler) GetAllInventories(c *fiber.Ctx) error {
	inventories, err := h.useCase.GetAllInventories(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(inventories)
}

func (h *InventoryQueryHandler) GetStoreInventory(c *fiber.Ctx) error {
	storeID := c.Params("storeId")
	inventories, err := h.useCase.GetStoreInventory(c.Context(), storeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(inventories)
}

func (h *InventoryQueryHandler) GetItemInventory(c *fiber.Ctx) error {
	itemID := c.Params("itemId")
	inventories, err := h.useCase.GetItemInventory(c.Context(), itemID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(inventories)
}

func (h *InventoryQueryHandler) GetInventory(c *fiber.Ctx) error {
	query := dto.InventoryQuery{
		StoreID: c.Params("storeId"),
		ItemID:  c.Params("itemId"),
	}

	inventory, err := h.useCase.GetInventory(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if inventory == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Inventory not found",
		})
	}

	return c.JSON(inventory)
}
