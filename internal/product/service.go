package product

import (
	"context"
	"log"

	"github.com/diorshelton/golden-market-api/internal/models"
	"github.com/diorshelton/golden-market-api/internal/repository"
	"github.com/google/uuid"
)

type ProductService struct {
	ProductRepository *repository.ProductRepository
}

func NewProductService(productRepo *repository.ProductRepository) *ProductService {
	return &ProductService{
	ProductRepository: productRepo,
	}
}

func (s *ProductService) Create(product *models.Product) error {
	product.ID = uuid.New()

	err := s.ProductRepository.Create(context.Background(), product)
	if err != nil {
		log.Printf("An error occurred while adding products %v", err)
		return err
	}
	return nil
}

func (s *ProductService) Update(id uuid.UUID) {

}

func (s *ProductService) Delete(id uuid.UUID) {

}

func (s *ProductService) GetProduct(id uuid.UUID) error {
	return nil
}

func (s *ProductService) GetProducts() ([]*models.Product, error) {
	// s.ProductRepository.GetAll(context.Background())
	return []*models.Product{}, nil
}