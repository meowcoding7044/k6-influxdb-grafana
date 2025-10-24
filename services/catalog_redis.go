package services

import (
	"context"
	"encoding/json"
	"fmt"
	"goFirst1/repositories"
	"time"

	"github.com/redis/go-redis/v9"
)

type catalogServiceRedis struct {
	productRepo repositories.ProductRepository
	redisClient *redis.Client
}

func NewCatalogServiceRedis(productRepo repositories.ProductRepository, redisClient *redis.Client) CatalogService {
	return catalogServiceRedis{productRepo, redisClient}
}
func (s catalogServiceRedis) GetProducts() (products []Product, err error) {

	key := "service::GetProducts"
	//Redis Get
	if productJson, err := s.redisClient.Get(context.Background(), key).Result(); err == nil {
		if json.Unmarshal([]byte(productJson), &products) == nil {
			fmt.Println("product from ram (Redis)")
			return products, nil
		}
	}
	//Data from disk
	productDB, err := s.productRepo.GetProducts()
	if err != nil {
		return nil, err
	}
	for _, p := range productDB {
		products = append(products, Product{
			ID:       p.ID,
			Name:     p.Name,
			Quantity: p.Quantity,
		})
	}
	//Redis Set
	if data, err := json.Marshal(products); err == nil {
		s.redisClient.Set(context.Background(), key, string(data), time.Second*10)
	}
	fmt.Println("product from disk")
	return products, nil
}
