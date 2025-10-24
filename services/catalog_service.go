package services

import (
	"goFirst1/repositories"
)

type catalogService struct {
	productRepo repositories.ProductRepository
}

func NewCatalogService(productRepo repositories.ProductRepository) CatalogService {
	return catalogService{productRepo: productRepo}
}
func (s catalogService) GetProducts() (products []Product, err error) {
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
	return products, nil
}
