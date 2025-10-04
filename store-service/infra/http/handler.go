package http

import (
	"github.com/LudSkywalker/inventory-system/store-service/app/dto"
	port "github.com/LudSkywalker/inventory-system/store-service/app/port/input"
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
	group.Get("/", h.ListInventory)
	group.Get("/:storeId/:itemId", h.GetStock)
	group.Delete("/:storeId/:itemId", h.DeleteStock)
}

// UpdateStock godoc
// @Summary Update stock quantity
// @Description Update the stock quantity for an item in a store
// @Tags inventory
// @Accept json
// @Produce json
// @Param request body dto.UpdateStockCommand true "Stock update request"
// @Success 200
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/inventory [post]
func (h *InventoryHandler) UpdateStock(c *fiber.Ctx) error {
	var req dto.UpdateStockCommand
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := h.useCase.UpdateStock(c.Context(), req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusOK)
}

// GetStock godoc
// @Summary Get stock information
// @Description Get stock quantity for a specific item in a store
// @Tags inventory
// @Accept json
// @Produce json
// @Param storeId path string true "Store ID"
// @Param itemId path string true "Item ID"
// @Success 200 {object} dto.InventoryDTO
// @Failure 500 {object} map[string]string
// @Router /api/v1/inventory/{storeId}/{itemId} [get]
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

// ListInventory godoc
// @Summary List all inventory
// @Description Get all inventory items across all stores
// @Tags inventory
// @Accept json
// @Produce json
// @Success 200 {array} dto.InventoryDTO
// @Failure 500 {object} map[string]string
// @Router /api/v1/inventory [get]
func (h *InventoryHandler) ListInventory(c *fiber.Ctx) error {
	inventories, err := h.useCase.ListInventory(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(inventories)
}

// DeleteStock godoc
// @Summary Delete stock
// @Description Delete stock information for a specific item in a store
// @Tags inventory
// @Accept json
// @Produce json
// @Param storeId path string true "Store ID"
// @Param itemId path string true "Item ID"
// @Success 200
// @Failure 500 {object} map[string]string
// @Router /api/v1/inventory/{storeId}/{itemId} [delete]
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
