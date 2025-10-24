package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"goFirst1/services"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/redis/go-redis/v9"
)

type catalogHandlerRedis struct {
	catalogSrv  services.CatalogService
	redisClient *redis.Client
}

func NewCatalogHandlerRedis(catalogSrv services.CatalogService, redisClient *redis.Client) CatalogHandler {
	return catalogHandlerRedis{catalogSrv, redisClient}
}

func (h catalogHandlerRedis) GetProducts(c fiber.Ctx) error {
	key := "handler::GetProducts"
	//Redis Get
	if products, err := h.redisClient.Get(context.Background(), key).Result(); err == nil {
		fmt.Println("product from ram (redis)")
		c.Set("Content-Type", "application/json")
		return c.SendString(products)
	}
	//Service
	products, err := h.catalogSrv.GetProducts()
	if err != nil {
		return err
	}
	response := fiber.Map{
		"status":   "ok",
		"products": products,
	}
	//Redis Set
	if data, err := json.Marshal(response); err == nil {
		h.redisClient.Set(context.Background(), key, string(data), time.Second*10)
	}
	fmt.Println("product from disk (DB)")
	return c.JSON(response)
}
