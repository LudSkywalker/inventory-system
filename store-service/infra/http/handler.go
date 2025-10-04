package http

import (
	"log"

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
		log.Printf("UpdateStock: Invalid request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	log.Printf("UpdateStock: Updating item %s in store %s to quantity %d", req.ItemID, req.StoreID, req.Quantity)
	if err := h.useCase.UpdateStock(c.Context(), req); err != nil {
		log.Printf("UpdateStock: Error updating stock: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	log.Printf("UpdateStock: Successfully updated item %s in store %s", req.ItemID, req.StoreID)
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

	log.Printf("GetStock: Retrieving item %s from store %s", query.ItemID, query.StoreID)
	inventory, err := h.useCase.GetStock(c.Context(), query)
	if err != nil {
		log.Printf("GetStock: Error retrieving stock: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	log.Printf("GetStock: Successfully retrieved item %s from store %s with quantity %d", query.ItemID, query.StoreID, inventory.Quantity)
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
	log.Println("ListInventory: Retrieving all inventory items")
	inventories, err := h.useCase.ListInventory(c.Context())
	if err != nil {
		log.Printf("ListInventory: Error retrieving inventory list: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	log.Printf("ListInventory: Successfully retrieved %d inventory items", len(inventories))
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

	log.Printf("DeleteStock: Deleting item %s from store %s", cmd.ItemID, cmd.StoreID)
	if err := h.useCase.DeleteStock(c.Context(), cmd); err != nil {
		log.Printf("DeleteStock: Error deleting stock: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	log.Printf("DeleteStock: Successfully deleted item %s from store %s", cmd.ItemID, cmd.StoreID)
	return c.SendStatus(fiber.StatusOK)
}
