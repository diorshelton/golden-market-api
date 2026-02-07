package order

import (
	"context"
	"fmt"
	"time"

	"github.com/diorshelton/golden-market-api/internal/models"
	"github.com/diorshelton/golden-market-api/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderService struct {
	db                  *pgxpool.Pool
	orderRepo           *repository.OrderRepository
	orderItemRepo       *repository.OrderItemRepository
	inventoryRepo       *repository.InventoryRepository
	userRepo            *repository.UserRepository
	productRepo         *repository.ProductRepository
	cartRepo            *repository.CartRepository
	duplicateWindowSecs int
}

func NewOrderService(
	db *pgxpool.Pool,
	orderRepo *repository.OrderRepository,
	orderItemRepo *repository.OrderItemRepository,
	inventoryRepo *repository.InventoryRepository,
	userRepo *repository.UserRepository,
	productRepo *repository.ProductRepository,
	cartRepo *repository.CartRepository,
) *OrderService {
	return &OrderService{
		db:                  db,
		orderRepo:           orderRepo,
		orderItemRepo:       orderItemRepo,
		inventoryRepo:       inventoryRepo,
		userRepo:            userRepo,
		productRepo:         productRepo,
		cartRepo:            cartRepo,
		duplicateWindowSecs: 15, // Prevent duplicate orders within 15 seconds
	}
}

// CreateOrder processes the checkout atomically:
// 1. Validate cart is not empty
// 2. Check for duplicate recent order
// 3. Begin transaction
// 4. Verify user has sufficient coins (with row lock)
// 5. Verify all products have sufficient stock (with row locks)
// 6. Deduct coins from user
// 7. Decrement stock for each product
// 8. Create order and order items
// 9. Add items to user inventory
// 10. Clear cart
// 11. Commit transaction
func (s *OrderService) CreateOrder(ctx context.Context, userID uuid.UUID) (*models.Order, error) {
	// Get cart items (outside transaction - just for empty check and duplicate detection)
	cart, err := s.cartRepo.GetCart(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart: %w", err)
	}

	if len(cart.Items) == 0 {
		return nil, fmt.Errorf("cart is empty")
	}

	// Check for recent duplicate orders
	recentOrders, err := s.orderRepo.GetRecentByUserID(ctx, userID, s.duplicateWindowSecs)
	if err != nil {
		return nil, fmt.Errorf("failed to check recent orders: %w", err)
	}

	// If there's a recent order with the same total, return it (idempotency)
	totalAmount := int(cart.TotalPrice)
	for _, recentOrder := range recentOrders {
		if recentOrder.TotalAmount == totalAmount {
			// Load order items
			items, err := s.orderItemRepo.GetByOrderID(ctx, recentOrder.ID)
			if err == nil {
				recentOrder.Items = items
			}
			return recentOrder, nil
		}
	}

	// Begin transaction BEFORE validation to ensure consistency
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Get user with row lock to verify balance
	user, err := s.userRepo.GetUserByIDTx(ctx, tx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if int(user.Balance) < totalAmount {
		return nil, fmt.Errorf("insufficient coins: have %d, need %d", user.Balance, totalAmount)
	}

	// Verify stock availability with row locks and get fresh product data
	for i, item := range cart.Items {
		product, err := s.productRepo.GetByIDForUpdate(ctx, tx, item.Product.ID)
		if err != nil {
			return nil, fmt.Errorf("product %s is no longer available", item.Product.Name)
		}
		if product.Stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for %s: available %d, requested %d",
				product.Name, product.Stock, item.Quantity)
		}
		// Update cart item with fresh product data
		cart.Items[i].Product = *product
	}

	// Create order
	now := time.Now().UTC()
	order := &models.Order{
		ID:          uuid.New(),
		UserID:      userID,
		OrderNumber: repository.GenerateOrderNumber(),
		TotalAmount: totalAmount,
		Status:      models.OrderStatusCompleted,
		CreatedAt:   now,
		UpdatedAt:   now,
		Items:       make([]models.OrderItem, 0, len(cart.Items)),
	}

	// Deduct coins from user
	if err := s.userRepo.DeductCoins(ctx, tx, userID, totalAmount); err != nil {
		return nil, err
	}

	// Create order record FIRST (before order items due to foreign key constraint)
	if err := s.orderRepo.Create(ctx, tx, order); err != nil {
		return nil, err
	}

	// Process each cart item
	for _, cartItem := range cart.Items {
		// Decrement product stock
		if err := s.productRepo.DecrementStockTx(ctx, tx, cartItem.Product.ID, cartItem.Quantity); err != nil {
			return nil, fmt.Errorf("failed to decrement stock for %s: %w", cartItem.Product.Name, err)
		}

		// Add to user inventory
		if err := s.inventoryRepo.AddOrUpdate(ctx, tx, userID, cartItem.Product.ID, cartItem.Quantity); err != nil {
			return nil, fmt.Errorf("failed to add to inventory: %w", err)
		}

		// Create order item
		orderItem := models.OrderItem{
			ID:           uuid.New(),
			OrderID:      order.ID,
			ProductID:    cartItem.Product.ID,
			ProductName:  cartItem.Product.Name,
			Quantity:     cartItem.Quantity,
			PricePerUnit: int(cartItem.Product.Price),
			Subtotal:     int(cartItem.Subtotal),
			CreatedAt:    now,
		}

		if err := s.orderItemRepo.Create(ctx, tx, &orderItem); err != nil {
			return nil, fmt.Errorf("failed to create order item: %w", err)
		}

		order.Items = append(order.Items, orderItem)
	}

	// Clear cart
	if err := s.cartRepo.ClearCart(ctx, tx, userID); err != nil {
		return nil, fmt.Errorf("failed to clear cart: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return order, nil
}

// GetOrderByID retrieves an order by ID with its items
func (s *OrderService) GetOrderByID(ctx context.Context, orderID uuid.UUID) (*models.Order, error) {
	order, err := s.orderRepo.GetByID(ctx, orderID)
	if err != nil {
		return nil, err
	}

	items, err := s.orderItemRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}

	order.Items = items
	return order, nil
}

// GetUserOrders retrieves all orders for a user with their items
func (s *OrderService) GetUserOrders(ctx context.Context, userID uuid.UUID) ([]*models.Order, error) {
	orders, err := s.orderRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Load items for each order
	for _, order := range orders {
		items, err := s.orderItemRepo.GetByOrderID(ctx, order.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get items for order %s: %w", order.ID, err)
		}
		order.Items = items
	}

	return orders, nil
}

// BeginTx starts a new database transaction (for admin operations)
func (s *OrderService) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return s.db.Begin(ctx)
}
