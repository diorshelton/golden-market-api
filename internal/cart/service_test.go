package cart

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
	cartService *CartService
	userRepo    *repository.UserRepository
	productRepo *repository.ProductRepository
	cartRepo    *repository.CartRepository
}

func setupCartTest(t *testing.T) *testDeps {
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

	cartService := NewCartService(cartRepo, productRepo)

	return &testDeps{
		cartService: cartService,
		userRepo:    userRepo,
		productRepo: productRepo,
		cartRepo:    cartRepo,
	}
}

var testCounter int

func uniqueSuffix() string {
	testCounter++
	return fmt.Sprintf("%d_%d", time.Now().UnixNano(), testCounter)
}

func createTestUser(t *testing.T, deps *testDeps) *models.User {
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

	return user
}

func createTestProduct(t *testing.T, deps *testDeps, price, stock int) *models.Product {
	t.Helper()
	suffix := uniqueSuffix()

	product := &models.Product{
		ID:          uuid.New(),
		Name:        "Test Product " + suffix,
		Description: "A product used for cart service tests",
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

// TestAddToCart_HappyPath verifies a valid add-to-cart call actually
// persists a cart item.
func TestAddToCart_HappyPath(t *testing.T) {
	deps := setupCartTest(t)
	ctx := context.Background()

	user := createTestUser(t, deps)
	product := createTestProduct(t, deps, 100, 5)

	if err := deps.cartService.AddToCart(ctx, user.ID, product.ID, 2); err != nil {
		t.Fatalf("AddToCart returned unexpected error: %v", err)
	}

	cart, err := deps.cartRepo.GetCart(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to load cart: %v", err)
	}
	if len(cart.Items) != 1 || cart.Items[0].Quantity != 2 {
		t.Errorf("expected 1 cart item with quantity 2, got %+v", cart.Items)
	}
}

// TestAddToCart_InsufficientStock verifies the service rejects the add
// before it ever reaches the repository when stock is too low.
func TestAddToCart_InsufficientStock(t *testing.T) {
	deps := setupCartTest(t)
	ctx := context.Background()

	user := createTestUser(t, deps)
	product := createTestProduct(t, deps, 100, 1)

	err := deps.cartService.AddToCart(ctx, user.ID, product.ID, 2)
	if err == nil {
		t.Fatal("expected error for insufficient stock, got nil")
	}

	cart, err := deps.cartRepo.GetCart(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to load cart: %v", err)
	}
	if len(cart.Items) != 0 {
		t.Errorf("expected no cart item to be added, got %+v", cart.Items)
	}
}

// TestAddToCart_ProductNotFound verifies a nonexistent product is rejected
// with a clear error instead of a database-level failure.
func TestAddToCart_ProductNotFound(t *testing.T) {
	deps := setupCartTest(t)
	ctx := context.Background()

	user := createTestUser(t, deps)

	err := deps.cartService.AddToCart(ctx, user.ID, uuid.New(), 1)
	if err == nil {
		t.Fatal("expected error for nonexistent product, got nil")
	}
}

// TestUpdateCartItemQuantity_HappyPath verifies a valid quantity update
// against the caller's own cart item succeeds.
func TestUpdateCartItemQuantity_HappyPath(t *testing.T) {
	deps := setupCartTest(t)
	ctx := context.Background()

	user := createTestUser(t, deps)
	product := createTestProduct(t, deps, 100, 5)

	if err := deps.cartRepo.AddToCart(ctx, user.ID, product.ID, 1); err != nil {
		t.Fatalf("failed to seed cart item: %v", err)
	}

	cart, err := deps.cartRepo.GetCart(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to load cart: %v", err)
	}
	cartItemID := cart.Items[0].CartItemID

	if err := deps.cartService.UpdateCartItemQuantity(ctx, user.ID, cartItemID, 3); err != nil {
		t.Fatalf("UpdateCartItemQuantity returned unexpected error: %v", err)
	}

	updatedCart, err := deps.cartRepo.GetCart(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to reload cart: %v", err)
	}
	if updatedCart.Items[0].Quantity != 3 {
		t.Errorf("expected quantity 3, got %d", updatedCart.Items[0].Quantity)
	}
}

// TestUpdateCartItemQuantity_InsufficientStock verifies the service rejects
// a quantity update that exceeds available stock, and leaves it unchanged.
func TestUpdateCartItemQuantity_InsufficientStock(t *testing.T) {
	deps := setupCartTest(t)
	ctx := context.Background()

	user := createTestUser(t, deps)
	product := createTestProduct(t, deps, 100, 2)

	if err := deps.cartRepo.AddToCart(ctx, user.ID, product.ID, 1); err != nil {
		t.Fatalf("failed to seed cart item: %v", err)
	}

	cart, err := deps.cartRepo.GetCart(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to load cart: %v", err)
	}
	cartItemID := cart.Items[0].CartItemID

	err = deps.cartService.UpdateCartItemQuantity(ctx, user.ID, cartItemID, 5)
	if err == nil {
		t.Fatal("expected error for insufficient stock, got nil")
	}

	updatedCart, err := deps.cartRepo.GetCart(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to reload cart: %v", err)
	}
	if updatedCart.Items[0].Quantity != 1 {
		t.Errorf("expected quantity to remain 1, got %d", updatedCart.Items[0].Quantity)
	}
}

// TestUpdateCartItemQuantity_WrongOwner verifies a user cannot update a
// cart item that belongs to a different user.
func TestUpdateCartItemQuantity_WrongOwner(t *testing.T) {
	deps := setupCartTest(t)
	ctx := context.Background()

	owner := createTestUser(t, deps)
	intruder := createTestUser(t, deps)
	product := createTestProduct(t, deps, 100, 5)

	if err := deps.cartRepo.AddToCart(ctx, owner.ID, product.ID, 1); err != nil {
		t.Fatalf("failed to seed cart item: %v", err)
	}

	cart, err := deps.cartRepo.GetCart(ctx, owner.ID)
	if err != nil {
		t.Fatalf("failed to load cart: %v", err)
	}
	cartItemID := cart.Items[0].CartItemID

	err = deps.cartService.UpdateCartItemQuantity(ctx, intruder.ID, cartItemID, 2)
	if err == nil {
		t.Fatal("expected error when updating another user's cart item, got nil")
	}

	ownerCart, err := deps.cartRepo.GetCart(ctx, owner.ID)
	if err != nil {
		t.Fatalf("failed to reload owner cart: %v", err)
	}
	if ownerCart.Items[0].Quantity != 1 {
		t.Errorf("expected owner's cart item quantity unchanged at 1, got %d", ownerCart.Items[0].Quantity)
	}
}
