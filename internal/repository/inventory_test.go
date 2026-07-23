package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/diorshelton/golden-market-api/internal/database"
	"github.com/diorshelton/golden-market-api/internal/models"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func setupInventoryTest(t *testing.T) (*InventoryRepository, *UserRepository, *ProductRepository) {
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

	return NewInventoryRepository(db), NewUserRepository(db), NewProductRepository(db)
}

var testCounter int

func uniqueSuffix() string {
	testCounter++
	return fmt.Sprintf("%d_%d", time.Now().UnixNano(), testCounter)
}

func seedInventoryUser(t *testing.T, userRepo *UserRepository) *models.User {
	t.Helper()
	suffix := uniqueSuffix()

	user, err := userRepo.CreateUser("inv_user_"+suffix, "Test", "User", "inv_"+suffix+"@example.com", "hash")
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}
	return user
}

func seedInventoryProduct(t *testing.T, productRepo *ProductRepository) *models.Product {
	t.Helper()
	suffix := uniqueSuffix()

	product := &models.Product{
		ID:          uuid.New(),
		Name:        "Inventory Test Product " + suffix,
		Description: "used for inventory repository tests",
		Price:       models.Coins(100),
		Stock:       10,
		Category:    "test",
		IsAvailable: true,
	}
	if err := productRepo.Create(context.Background(), product); err != nil {
		t.Fatalf("failed to create test product: %v", err)
	}
	return product
}

// TestInventoryRepository_AddOrUpdate_Insert verifies a first-time
// AddOrUpdate call creates a new inventory row with the given quantity.
func TestInventoryRepository_AddOrUpdate_Insert(t *testing.T) {
	inventoryRepo, userRepo, productRepo := setupInventoryTest(t)
	ctx := context.Background()

	user := seedInventoryUser(t, userRepo)
	product := seedInventoryProduct(t, productRepo)

	if err := inventoryRepo.AddOrUpdate(ctx, inventoryRepo.db, user.ID, product.ID, 3); err != nil {
		t.Fatalf("AddOrUpdate returned unexpected error: %v", err)
	}

	items, err := inventoryRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to load inventory: %v", err)
	}
	if len(items) != 1 || items[0].Quantity != 3 {
		t.Errorf("expected 1 item with quantity 3, got %+v", items)
	}
}

// TestInventoryRepository_AddOrUpdate_Accumulates verifies a second
// AddOrUpdate call for the same user/product ADDS to the existing quantity
// rather than replacing it - this is the exact behavior CreateOrder relies
// on when a user buys the same product across two separate orders.
func TestInventoryRepository_AddOrUpdate_Accumulates(t *testing.T) {
	inventoryRepo, userRepo, productRepo := setupInventoryTest(t)
	ctx := context.Background()

	user := seedInventoryUser(t, userRepo)
	product := seedInventoryProduct(t, productRepo)

	if err := inventoryRepo.AddOrUpdate(ctx, inventoryRepo.db, user.ID, product.ID, 3); err != nil {
		t.Fatalf("first AddOrUpdate returned unexpected error: %v", err)
	}
	if err := inventoryRepo.AddOrUpdate(ctx, inventoryRepo.db, user.ID, product.ID, 2); err != nil {
		t.Fatalf("second AddOrUpdate returned unexpected error: %v", err)
	}

	items, err := inventoryRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to load inventory: %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected exactly 1 inventory row (upsert, not duplicate), got %d", len(items))
	}
	if items[0].Quantity != 5 {
		t.Errorf("expected accumulated quantity 5, got %d", items[0].Quantity)
	}
}

// TestInventoryRepository_ClearByUserID verifies all inventory rows for a
// user are removed - used by the guest-login reset flow.
func TestInventoryRepository_ClearByUserID(t *testing.T) {
	inventoryRepo, userRepo, productRepo := setupInventoryTest(t)
	ctx := context.Background()

	user := seedInventoryUser(t, userRepo)
	product := seedInventoryProduct(t, productRepo)

	if err := inventoryRepo.AddOrUpdate(ctx, inventoryRepo.db, user.ID, product.ID, 3); err != nil {
		t.Fatalf("AddOrUpdate returned unexpected error: %v", err)
	}

	if err := inventoryRepo.ClearByUserID(ctx, inventoryRepo.db, user.ID); err != nil {
		t.Fatalf("ClearByUserID returned unexpected error: %v", err)
	}

	items, err := inventoryRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed to load inventory: %v", err)
	}
	if len(items) != 0 {
		t.Errorf("expected inventory to be empty after clear, got %+v", items)
	}
}
