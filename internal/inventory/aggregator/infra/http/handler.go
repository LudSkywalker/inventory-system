package http

import (
	"github.com/LudSkywalker/inventory-system/internal/inventory/aggregator/app/port"
	"github.com/gofiber/fiber/v2"
)

type GlobalInventoryHandler struct {
	useCase port.GlobalInventoryUseCase
}

func NewGlobalInventoryHandler(useCase port.GlobalInventoryUseCase) *GlobalInventoryHandler {
	return &GlobalInventoryHandler{useCase: useCase}
}

func (h *GlobalInventoryHandler) RegisterRoutes(app *fiber.App) {
	group := app.Group("/api/v1/global-inventory")

	group.Get("/", h.GetAllInventories)
	group.Get("/store/:storeId", h.GetStoreInventory)
	group.Get("/item/:itemId", h.GetItemInventory)
	group.Get("/:storeId/:itemId", h.GetInventory)
}

func (h *GlobalInventoryHandler) GetAllInventories(c *fiber.Ctx) error {
	inventories, err := h.useCase.GetAllInventories(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(inventories)
}

func (h *GlobalInventoryHandler) GetStoreInventory(c *fiber.Ctx) error {
	storeID := c.Params("storeId")
	inventories, err := h.useCase.GetStoreInventory(c.Context(), storeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(inventories)
}

func (h *GlobalInventoryHandler) GetItemInventory(c *fiber.Ctx) error {
	itemID := c.Params("itemId")
	inventories, err := h.useCase.GetItemInventory(c.Context(), itemID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(inventories)
}

func (h *GlobalInventoryHandler) GetInventory(c *fiber.Ctx) error {
	storeID := c.Params("storeId")
	itemID := c.Params("itemId")

	inventory, err := h.useCase.GetGlobalInventory(c.Context(), itemID, storeID)
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
