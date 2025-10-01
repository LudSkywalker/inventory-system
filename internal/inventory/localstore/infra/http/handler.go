package http

import (
	"github.com/LudSkywalker/inventory-system/internal/inventory/localstore/app/dto"
	"github.com/LudSkywalker/inventory-system/internal/inventory/localstore/app/port"
	"github.com/gofiber/fiber/v2"
)

type InventoryHandler struct {
	useCase port.InventoryUseCase
}

func NewInventoryHandler(useCase port.InventoryUseCase) *InventoryHandler {
	return &InventoryHandler{useCase: useCase}
}

func (h *InventoryHandler) RegisterRoutes(app *fiber.App) {
	group := app.Group("/api/v1/inventory")

	group.Post("/", h.UpdateStock)
	group.Get("/:storeId/:itemId", h.GetStock)
	group.Delete("/:storeId/:itemId", h.DeleteStock)
}

func (h *InventoryHandler) UpdateStock(c *fiber.Ctx) error {
	var cmd dto.UpdateStockCommand
	if err := c.BodyParser(&cmd); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.useCase.UpdateStock(c.Context(), cmd); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

func (h *InventoryHandler) GetStock(c *fiber.Ctx) error {
	query := dto.GetStockQuery{
		StoreID: c.Params("storeId"),
		ItemID:  c.Params("itemId"),
	}

	inventory, err := h.useCase.GetStock(c.Context(), query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(inventory)
}

func (h *InventoryHandler) DeleteStock(c *fiber.Ctx) error {
	cmd := dto.DeleteStockCommand{
		StoreID: c.Params("storeId"),
		ItemID:  c.Params("itemId"),
	}

	if err := h.useCase.DeleteStock(c.Context(), cmd); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}
