package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/diorshelton/golden-market-api/internal/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepository(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

// Create inserts a new order into the database (within a transaction)
func (r *OrderRepository) Create(ctx context.Context, tx DBTX, order *models.Order) error {
	query := `
		INSERT INTO orders (id, user_id, order_number, total_amount, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := tx.Exec(ctx, query,
		order.ID,
		order.UserID,
		order.OrderNumber,
		order.TotalAmount,
		order.Status,
		order.CreatedAt,
		order.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	return nil
}

// GetByID retrieves an order by its ID
func (r *OrderRepository) GetByID(ctx context.Context, orderID uuid.UUID) (*models.Order, error) {
	query := `
		SELECT id, user_id, order_number, total_amount, status, created_at, updated_at
		FROM orders
		WHERE id = $1
	`

	var order models.Order
	err := r.db.QueryRow(ctx, query, orderID).Scan(
		&order.ID,
		&order.UserID,
		&order.OrderNumber,
		&order.TotalAmount,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("order not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return &order, nil
}

// GetByUserID retrieves all orders for a user
func (r *OrderRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*models.Order, error) {
	query := `
		SELECT id, user_id, order_number, total_amount, status, created_at, updated_at
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query orders: %w", err)
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.OrderNumber,
			&order.TotalAmount,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating orders: %w", err)
	}

	return orders, nil
}

// GetRecentByUserID retrieves orders from the last N seconds for duplicate prevention
func (r *OrderRepository) GetRecentByUserID(ctx context.Context, userID uuid.UUID, seconds int) ([]*models.Order, error) {
	query := `
		SELECT id, user_id, order_number, total_amount, status, created_at, updated_at
		FROM orders
		WHERE user_id = $1 AND created_at > NOW() - INTERVAL '1 second' * $2
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, query, userID, seconds)
	if err != nil {
		return nil, fmt.Errorf("failed to query recent orders: %w", err)
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var order models.Order
		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.OrderNumber,
			&order.TotalAmount,
			&order.Status,
			&order.CreatedAt,
			&order.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan order: %w", err)
		}
		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating orders: %w", err)
	}

	return orders, nil
}

// GetPool returns the database pool for transaction management
func (r *OrderRepository) GetPool() *pgxpool.Pool {
	return r.db
}

// BeginTx starts a new database transaction
func (r *OrderRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.db.Begin(ctx)
}

// GenerateOrderNumber creates a unique order number (max 20 chars)
func GenerateOrderNumber() string {
	now := time.Now()
	return fmt.Sprintf("ORD-%s-%s",
		now.Format("060102"),
		uuid.New().String()[:6],
	)
}
