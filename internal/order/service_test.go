package order

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/diorshelton/golden-market-api/internal/database"
	"github.com/diorshelton/golden-market-api/internal/models"
	"github.com/diorshelton/golden-market-api/internal/repository"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type testDeps struct {
	orderService  *OrderService
	userRepo      *repository.UserRepository
	productRepo   *repository.ProductRepository
	cartRepo      *repository.CartRepository
	inventoryRepo *repository.InventoryRepository
}

func setupOrderTest(t *testing.T) *testDeps {
	t.Helper()

	_ = godotenv.Load("../../.env")
	if os.Getenv("TEMP_DB_URL") == "" {
		t.Skip("TEMP_DB_URL not set, skipping database tests")
	}

	db, err := database.SetupTestDB()
	if err != nil {
		t.Fatalf("failed to set up test db: %v", err)
	}
	t.Cleanup(func() { db.Close() })

	userRepo := repository.NewUserRepository(db)
	productRepo := repository.NewProductRepository(db)
	cartRepo := repository.NewCartRepository(db)
	inventoryRepo := repository.NewInventoryRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	orderItemRepo := repository.NewOrderItemRepository(db)

	orderService := NewOrderService(db, orderRepo, orderItemRepo, inventoryRepo, userRepo, productRepo, cartRepo)

	return &testDeps{
		orderService:  orderService,
		userRepo:      userRepo,
		productRepo:   productRepo,
		cartRepo:      cartRepo,
		inventoryRepo: inventoryRepo,
	}
}

var testCounter int

func uniqueSuffix() string {
	testCounter++
	return fmt.Sprintf("%d_%d", time.Now().UnixNano(), testCounter)
}

func createTestUser(t *testing.T, deps *testDeps, balance int) *models.User {
	t.Helper()
	suffix := uniqueSuffix()

	user, err := deps.userRepo.CreateUser(
		"user_"+suffix,
		"Test",
		"User",
		"user_"+suffix+"@example.com",
		"hashed_password",
	)
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	if err := deps.userRepo.UpdateBalance(user.ID, models.Coins(balance)); err != nil {
		t.Fatalf("failed to set test user balance: %v", err)
	}
	user.Balance = models.Coins(balance)

	return user
}

func createTestProduct(t *testing.T, deps *testDeps, price, stock int) *models.Product {
	t.Helper()
	suffix := uniqueSuffix()

	product := &models.Product{
		ID:          uuid.New(),
		Name:        "Test Product " + suffix,
		Description: "A product used for order service tests",
		Price:       models.Coins(price),
		Stock:       stock,
		Category:    "test",
		IsAvailable: true,
	}

	if err := deps.productRepo.Create(context.Background(), product); err != nil {
		t.Fatalf("failed to create test product: %v", err)
	}

	return product
}

// TestCreateOrder_HappyPath verifies a full checkout: balance deducted,
// stock decremented, inventory credited, order+items persisted, cart cleared.
func TestCreateOrder_HappyPath(t *testing.T) {
	deps := setupOrderTest(t)
	ctx := context.Background()

	user := createTestUser(t, deps, 1000)
	product := createTestProduct(t, deps, 300, 5)

	if err := deps.cartRepo.AddToCart(ctx, user.ID, product.ID, 2); err != nil {
		t.Fatalf("failed to add to cart: %v", err)
	}

	order, err := deps.orderService.CreateOrder(ctx, user.ID)
	if err != nil {
		t.Fatalf("CreateOrder returned unexpected error: %v", err)
	}

	if order.TotalAmount != 600 {
		t.Errorf("expected order total 600, got %d", order.TotalAmount)
	}
	if len(order.Items) != 1 {
		t.Fatalf("expected 1 order item, got %d", len(order.Items))
	}
	if order.Items[0].Quantity != 2 {
		t.Errorf("expected order item quantity 2, got %d", order.Items[0].Quantity)
	}

	updatedUser, err := deps.userRepo.GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("failed to reload user: %v", err)
	}
	if updatedUser.Balance != 400 {
		t.Errorf("expected balance 400 after checkout, got %d", updatedUser.Balance)
	}

	updatedProduct, err := deps.productRepo.GetByID(ctx, product.ID)
	if err != nil {
		t.Fatalf("failed to reload product: %v", err)
	}
	if updatedProduct.Stock != 3 {
		t.Errorf("expected stock 3 after checkout, got %d", updatedProduct.Stock)
	}

	inventory, err := deps.inventoryRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to load inventory: %v", err)
	}
	if len(inventory) != 1 || inventory[0].Quantity != 2 {
		t.Errorf("expected inventory with quantity 2, got %+v", inventory)
	}

	cart, err := deps.cartRepo.GetCart(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to reload cart: %v", err)
	}
	if len(cart.Items) != 0 {
		t.Errorf("expected cart to be cleared, got %d items", len(cart.Items))
	}
}

// TestCreateOrder_InsufficientBalance verifies the transaction rolls back
// completely — no partial state changes — when the user can't afford the cart.
func TestCreateOrder_InsufficientBalance(t *testing.T) {
	deps := setupOrderTest(t)
	ctx := context.Background()

	user := createTestUser(t, deps, 100)
	product := createTestProduct(t, deps, 300, 5)

	if err := deps.cartRepo.AddToCart(ctx, user.ID, product.ID, 1); err != nil {
		t.Fatalf("failed to add to cart: %v", err)
	}

	_, err := deps.orderService.CreateOrder(ctx, user.ID)
	if err == nil {
		t.Fatal("expected error for insufficient balance, got nil")
	}

	updatedUser, err := deps.userRepo.GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("failed to reload user: %v", err)
	}
	if updatedUser.Balance != 100 {
		t.Errorf("expected balance unchanged at 100, got %d", updatedUser.Balance)
	}

	updatedProduct, err := deps.productRepo.GetByID(ctx, product.ID)
	if err != nil {
		t.Fatalf("failed to reload product: %v", err)
	}
	if updatedProduct.Stock != 5 {
		t.Errorf("expected stock unchanged at 5, got %d", updatedProduct.Stock)
	}

	cart, err := deps.cartRepo.GetCart(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to reload cart: %v", err)
	}
	if len(cart.Items) != 1 {
		t.Errorf("expected cart to still have 1 item, got %d", len(cart.Items))
	}
}

// TestCreateOrder_InsufficientStock verifies the transaction rolls back
// when stock is insufficient, even though the user can afford it.
func TestCreateOrder_InsufficientStock(t *testing.T) {
	deps := setupOrderTest(t)
	ctx := context.Background()

	user := createTestUser(t, deps, 10000)
	product := createTestProduct(t, deps, 100, 1)

	if err := deps.cartRepo.AddToCart(ctx, user.ID, product.ID, 1); err != nil {
		t.Fatalf("failed to add to cart: %v", err)
	}

	// Drain stock to 0 after the cart item was added, so the service's
	// fresh-stock check inside the transaction is what catches it.
	if err := deps.productRepo.DecrementStock(ctx, product.ID, 1); err != nil {
		t.Fatalf("failed to drain stock: %v", err)
	}

	_, err := deps.orderService.CreateOrder(ctx, user.ID)
	if err == nil {
		t.Fatal("expected error for insufficient stock, got nil")
	}

	updatedUser, err := deps.userRepo.GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("failed to reload user: %v", err)
	}
	if updatedUser.Balance != 10000 {
		t.Errorf("expected balance unchanged at 10000, got %d", updatedUser.Balance)
	}

	inventory, err := deps.inventoryRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to load inventory: %v", err)
	}
	if len(inventory) != 0 {
		t.Errorf("expected no inventory added, got %+v", inventory)
	}
}

// TestCreateOrder_EmptyCart verifies checkout is rejected before any
// transaction is opened when the cart has nothing in it.
func TestCreateOrder_EmptyCart(t *testing.T) {
	deps := setupOrderTest(t)
	ctx := context.Background()

	user := createTestUser(t, deps, 1000)

	_, err := deps.orderService.CreateOrder(ctx, user.ID)
	if err == nil {
		t.Fatal("expected error for empty cart, got nil")
	}
}

// TestCreateOrder_DuplicateIdempotency verifies that calling CreateOrder
// again within the duplicate-detection window, for the same cart total,
// returns the original order instead of creating a second one.
func TestCreateOrder_DuplicateIdempotency(t *testing.T) {
	deps := setupOrderTest(t)
	ctx := context.Background()

	user := createTestUser(t, deps, 1000)
	product := createTestProduct(t, deps, 300, 5)

	if err := deps.cartRepo.AddToCart(ctx, user.ID, product.ID, 1); err != nil {
		t.Fatalf("failed to add to cart: %v", err)
	}

	firstOrder, err := deps.orderService.CreateOrder(ctx, user.ID)
	if err != nil {
		t.Fatalf("first CreateOrder failed: %v", err)
	}

	// Re-add the same item/quantity so the cart total matches again -
	// simulates a duplicate submit (e.g. double-click) before the cart
	// was visibly cleared client-side.
	if err := deps.cartRepo.AddToCart(ctx, user.ID, product.ID, 1); err != nil {
		t.Fatalf("failed to re-add to cart: %v", err)
	}

	secondOrder, err := deps.orderService.CreateOrder(ctx, user.ID)
	if err != nil {
		t.Fatalf("second CreateOrder failed: %v", err)
	}

	if secondOrder.ID != firstOrder.ID {
		t.Errorf("expected duplicate call to return original order %s, got new order %s",
			firstOrder.ID, secondOrder.ID)
	}

	updatedUser, err := deps.userRepo.GetUserByID(user.ID)
	if err != nil {
		t.Fatalf("failed to reload user: %v", err)
	}
	if updatedUser.Balance != 700 {
		t.Errorf("expected balance charged only once (700), got %d", updatedUser.Balance)
	}
}
