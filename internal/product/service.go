package product

import (
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

func (s *ProductService) Create() (*models.Product, error) {
	return &models.Product{}, nil
}

func (s *ProductService) Update(id uuid.UUID) {

}

func (s *ProductService) Delete(id uuid.UUID) {

}

func (s *ProductService) GetProduct(id uuid.UUID) error {
return nil
}

func (s *ProductService) GetProducts() ([]*models.Product, error) {
 return []*models.Product{}, nil
}

func (s *ProductService) UpdateStock(id uuid.UUID) error {
	return nil
}

func (s *ProductService) DecrementStock(id uuid.UUID) error {
	return nil
}
func (s *ProductService) SearchProducts(id uuid.UUID) error {
	return nil
}