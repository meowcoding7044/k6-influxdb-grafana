package main

import (
	"goFirst1/configs"
	"goFirst1/handlers"
	"goFirst1/repositories"
	"goFirst1/services"

	"github.com/gofiber/fiber/v3"
)

func init() {
	configs.InitTimeZone()
	configs.InitConfig()
	configs.InitDatabase()
	configs.InitRedis()
}
func main() {

	db := configs.Database.Db
	redisClient := configs.Redis.Client
	_ = redisClient
	productRepo := repositories.NewProductRepositoryDB(db)
	productService := services.NewCatalogServiceRedis(productRepo, redisClient)
	productHandler := handlers.NewCatalogHandler(productService)

	app := fiber.New()
	app.Get("/products", productHandler.GetProducts)

	app.Listen(":8000")
	//productRepo := repositories.NewProductRepositoryDB()
}
