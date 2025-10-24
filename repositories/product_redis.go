package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type productRepositoryRedis struct {
	db          *gorm.DB
	redisClient *redis.Client
}

func NewProductRepositoryRedis(db *gorm.DB, redisClient *redis.Client) ProductRepository {
	db.AutoMigrate(&product{})
	mockData(db)
	return productRepositoryRedis{db, redisClient}
}

func (r productRepositoryRedis) GetProducts() (products []product, err error) {
	key := "repository::GetProduct"
	//Redis Get
	projectJson, err := r.redisClient.Get(context.Background(), key).Result()
	if err == nil {
		err = json.Unmarshal([]byte(projectJson), &products)
		if err == nil {
			fmt.Println("product from ram (Redis)")
			return products, err
		}
	}
	//Database
	err = r.db.Order("quantity desc").Limit(30).Find(&products).Error
	if err != nil {
		return nil, err
	}
	data, err := json.Marshal(products)
	if err != nil {
		return nil, err
	}
	//Redis Set
	err = r.redisClient.Set(context.Background(), key, string(data), time.Second*10).Err()
	if err != nil {
		return nil, err
	}
	fmt.Println("product from disk")
	return products, err
}
