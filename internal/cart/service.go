package cart

import (
	"context"
	"fmt"

	"github.com/diorshelton/golden-market-api/internal/models"
	"github.com/diorshelton/golden-market-api/internal/repository"
	"github.com/google/uuid"
)

type CartService struct {
	CartRepository    *repository.CartRepository
	ProductRepository *repository.ProductRepository
}

func NewCartService(cartRepo *repository.CartRepository, productRepo *repository.ProductRepository) *CartService {
	return &CartService{
		CartRepository:    cartRepo,
		ProductRepository: productRepo,
	}
}

func (s *CartService) AddToCart(ctx context.Context, userID, productID uuid.UUID, quantity int) error {
	//verify product is available
	product, err := s.ProductRepository.GetByID(ctx, productID)
	if err != nil {
		return fmt.Errorf("product not found or unavailable %w", err)
	}

	//Check stock availability
	if product.Stock < quantity {
		return fmt.Errorf("insufficient stock: only %d available", product.Stock)
	}
	return s.CartRepository.AddToCart(ctx, userID, productID, quantity)
}

func (s *CartService) GetCart(ctx context.Context, userID uuid.UUID) (*models.CartSummary, error) {
	return s.CartRepository.GetCart(ctx, userID)
}

func (s *CartService) UpdateCartItemQuantity(ctx context.Context, userID, cartItemID uuid.UUID, quantity int) error {
	// Get the user's cart
	cart, err := s.CartRepository.GetCart(ctx, userID)
	if err != nil {
		return fmt.Errorf("cart item not found %w", err)
	}

	// Find the specific cart item to verify ownership
	var cartItem *models.CartItemDetail
	for i := range cart.Items {
		if cart.Items[i].CartItemID == cartItemID {
			cartItem = &cart.Items[i]
			break
		}
	}

	if cartItem == nil {
		return fmt.Errorf("cart item not found or does not belong to user %w", err)
	}

	// Verify stock availability for the new quantity
	if cartItem.Product.Stock < quantity {
		return fmt.Errorf("insufficient stock: only %d available, err:%w", cartItem.Product.Stock, err)
	}

	return s.CartRepository.UpdateCartItemQuantity(ctx, cartItemID, quantity)
}

func (s *CartService) RemoveFromCart(ctx context.Context, userID, cartItemID uuid.UUID) error {
	return s.CartRepository.RemoveFromCart(ctx, userID, cartItemID)
}
