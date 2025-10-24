package handlers

import (
	"goFirst1/services"

	"github.com/gofiber/fiber/v3"
)

type catalogHandler struct {
	catalogSrv services.CatalogService
}

func NewCatalogHandler(catalogSrv services.CatalogService) CatalogHandler {
	return catalogHandler{catalogSrv}
}

func (h catalogHandler) GetProducts(c fiber.Ctx) error {
	products, err := h.catalogSrv.GetProducts()
	if err != nil {
		return err
	}
	response := fiber.Map{
		"status":   "ok",
		"products": products,
	}
	return c.JSON(response)
}
