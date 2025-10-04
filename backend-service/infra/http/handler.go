package http

import (
	"log"

	"github.com/LudSkywalker/inventory-system/backend-service/app/dto"
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
	group.Get("/inventory/grouped", h.GetGroupedInventory)

	// Health check endpoint
	app.Get("/health", h.HealthCheck)
}

// GetGlobalInventory godoc
// @Summary Get global inventory
// @Description Get all inventory items across all stores
// @Tags inventory
// @Accept json
// @Produce json
// @Success 200 {array} dto.InventoryDTO
// @Failure 500 {object} map[string]string
// @Router /api/v1/inventory [get]
func (h *InventoryHandler) GetGlobalInventory(c *fiber.Ctx) error {
	log.Println("HTTP Handler: GetGlobalInventory called")
	var inventories []dto.InventoryDTO
	var err error
	inventories, err = h.service.GetGlobalInventory(c.Context())
	if err != nil {
		log.Printf("HTTP Handler: GetGlobalInventory error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	log.Printf("HTTP Handler: GetGlobalInventory returning %d items", len(inventories))
	return c.JSON(inventories)
}

// HealthCheck godoc
// @Summary Health check
// @Description Check if the service is running
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /health [get]
func (h *InventoryHandler) HealthCheck(c *fiber.Ctx) error {
	log.Println("HTTP Handler: Health check called")
	return c.JSON(fiber.Map{
		"status":  "ok",
		"service": "backend-service",
	})
}

// GetStoreInventory godoc
// @Summary Get store inventory
// @Description Get all inventory items for a specific store
// @Tags inventory
// @Accept json
// @Produce json
// @Param storeId path string true "Store ID"
// @Success 200 {array} dto.InventoryDTO
// @Failure 500 {object} map[string]string
// @Router /api/v1/inventory/store/{storeId} [get]
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

// GetItemInventory godoc
// @Summary Get item inventory
// @Description Get all inventory information for a specific item across all stores
// @Tags inventory
// @Accept json
// @Produce json
// @Param itemId path string true "Item ID"
// @Success 200 {array} dto.InventoryDTO
// @Failure 500 {object} map[string]string
// @Router /api/v1/inventory/item/{itemId} [get]
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

// GetGroupedInventory godoc
// @Summary Get grouped inventory
// @Description Get all inventory items grouped by item ID with total quantity and store details
// @Tags inventory
// @Accept json
// @Produce json
// @Success 200 {array} dto.GroupedInventoryDTO
// @Failure 500 {object} map[string]string
// @Router /api/v1/inventory/grouped [get]
func (h *InventoryHandler) GetGroupedInventory(c *fiber.Ctx) error {
	groupedInventories, err := h.service.GetGroupedInventory(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(groupedInventories)
}
