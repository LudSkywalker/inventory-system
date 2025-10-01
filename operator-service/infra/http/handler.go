package http

import (
	"github.com/LudSkywalker/inventory-system/operator-service/app/port/input"
	"github.com/gofiber/fiber/v2"
)

type GlobalInventoryHandler struct {
	useCase input.GlobalInventoryUseCase
}

func NewGlobalInventoryHandler(useCase input.GlobalInventoryUseCase) *GlobalInventoryHandler {
	return &GlobalInventoryHandler{useCase: useCase}
}

func (h *GlobalInventoryHandler) RegisterRoutes(app *fiber.App) {
	group := app.Group("/api/v1/global-inventory")

	group.Get("/", h.GetAllInventories)
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

func (h *GlobalInventoryHandler) GetInventory(c *fiber.Ctx) error {
	storeID := c.Params("storeId")
	itemID := c.Params("itemId")

	inventory, err := h.useCase.GetGlobalInventory(c.Context(), itemID, storeID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(inventory)
}
