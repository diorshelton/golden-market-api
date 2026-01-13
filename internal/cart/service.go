package cart

import (
	"context"

	"github.com/diorshelton/golden-market-api/internal/models"
	"github.com/diorshelton/golden-market-api/internal/repository"
	"github.com/google/uuid"
)

type CartService struct {
	CartRepository *repository.CartRepository
}

func NewCartService(cartRepo *repository.CartRepository) *CartService {
	return &CartService{
		CartRepository: cartRepo,
	}
}

func (s *CartService) AddToCart(ctx context.Context, userID, productID uuid.UUID, quantity int) error {
	return s.CartRepository.AddToCart(ctx, userID, productID, quantity)
}

func (s *CartService) GetCart(ctx context.Context, userID uuid.UUID) (*models.CartSummary, error) {
	return s.CartRepository.GetCart(ctx, userID)
}

func (s *CartService) UpdateCartItemQuantity(ctx context.Context, userID, cartItemID uuid.UUID, quantity int) error {
	return s.CartRepository.UpdateCartItemQuantity(ctx, userID, cartItemID, quantity)
}

func (s *CartService) RemoveFromCart(ctx context.Context, userID, cartItemID uuid.UUID) error {
	return s.CartRepository.RemoveFromCart(ctx, userID, cartItemID)
}
