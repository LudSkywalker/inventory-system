package http

import (
	"github.com/LudSkywalker/inventory-system/db-service/domain/entity"
	"github.com/LudSkywalker/inventory-system/db-service/infra/sqlite"
	"github.com/gofiber/fiber/v2"
)

type GlobalInventoryHandler struct {
	repo *sqlite.SQLiteRepository
}

func NewGlobalInventoryHandler(repo *sqlite.SQLiteRepository) *GlobalInventoryHandler {
	return &GlobalInventoryHandler{repo: repo}
}

func (h *GlobalInventoryHandler) RegisterRoutes(app *fiber.App) {
	app.Post("/inventories", h.Save)
	app.Get("/inventories/:itemID/:storeID", h.FindByItemAndStore)
	app.Get("/inventories", h.FindAll)
	app.Delete("/inventories/:itemID/:storeID", h.Delete)
}

func (h *GlobalInventoryHandler) Save(c *fiber.Ctx) error {
	var req struct {
		ItemID   string `json:"item_id"`
		ItemName string `json:"item_name"`
		StoreID  string `json:"store_id"`
		Quantity int    `json:"quantity"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid request"})
	}

	inventory := entity.NewGlobalInventory(req.ItemID, req.ItemName, req.StoreID, req.Quantity)

	if err := h.repo.Save(c.Context(), inventory); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(inventory)
}

func (h *GlobalInventoryHandler) FindByItemAndStore(c *fiber.Ctx) error {
	itemID := c.Params("itemID")
	storeID := c.Params("storeID")

	inventory, err := h.repo.FindByItemAndStore(c.Context(), itemID, storeID)
	if err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(inventory)
}

func (h *GlobalInventoryHandler) FindAll(c *fiber.Ctx) error {
	inventories, err := h.repo.FindAll(c.Context())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(inventories)
}

func (h *GlobalInventoryHandler) Delete(c *fiber.Ctx) error {
	itemID := c.Params("itemID")
	storeID := c.Params("storeID")

	if err := h.repo.Delete(c.Context(), itemID, storeID); err != nil {
		return c.Status(404).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(204)
}
