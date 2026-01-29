package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/diorshelton/golden-market-api/internal/middleware"
	"github.com/diorshelton/golden-market-api/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type OrderServiceInterface interface {
	CreateOrder(ctx context.Context, userID uuid.UUID) (*models.Order, error)
	GetOrderByID(ctx context.Context, orderID uuid.UUID) (*models.Order, error)
	GetUserOrders(ctx context.Context, userID uuid.UUID) ([]*models.Order, error)
}

type OrderHandler struct {
	orderService OrderServiceInterface
}

func NewOrderHandler(service OrderServiceInterface) *OrderHandler {
	return &OrderHandler{
		orderService: service,
	}
}

// CreateOrder handles POST /api/v1/orders
func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	order, err := h.orderService.CreateOrder(r.Context(), userID)
	if err != nil {
		// Log the full error for debugging
		log.Printf("CreateOrder error for user %s: %v", userID, err)

		// Check for specific error types
		errMsg := err.Error()
		switch {
		case contains(errMsg, "cart is empty"):
			http.Error(w, errMsg, http.StatusBadRequest)
		case contains(errMsg, "insufficient coins"):
			http.Error(w, errMsg, http.StatusPaymentRequired)
		case contains(errMsg, "insufficient stock"):
			http.Error(w, errMsg, http.StatusConflict)
		default:
			http.Error(w, fmt.Sprintf("failed to create order: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// GetOrders handles GET /api/v1/orders
func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	orders, err := h.orderService.GetUserOrders(r.Context(), userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get orders: %v", err), http.StatusInternalServerError)
		return
	}

	// Return empty array instead of null
	if orders == nil {
		orders = []*models.Order{}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(orders)
}

// GetOrder handles GET /api/v1/orders/{id}
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	orderID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "invalid order ID", http.StatusBadRequest)
		return
	}

	order, err := h.orderService.GetOrderByID(r.Context(), orderID)
	if err != nil {
		if contains(err.Error(), "not found") {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("failed to get order: %v", err), http.StatusInternalServerError)
		return
	}

	// Verify the order belongs to the requesting user
	if order.UserID != userID {
		http.Error(w, "order not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsImpl(s, substr))
}

func containsImpl(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
