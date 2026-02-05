package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

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

		// Check for specific error types - these are safe to expose
		errMsg := err.Error()
		switch {
		case strings.Contains(errMsg, "cart is empty"):
			http.Error(w, errMsg, http.StatusBadRequest)
		case strings.Contains(errMsg, "insufficient coins"):
			http.Error(w, errMsg, http.StatusPaymentRequired)
		case strings.Contains(errMsg, "insufficient stock"):
			http.Error(w, errMsg, http.StatusConflict)
		case strings.Contains(errMsg, "no longer available"):
			http.Error(w, errMsg, http.StatusConflict)
		default:
			http.Error(w, "failed to create order", http.StatusInternalServerError)
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
		log.Printf("GetOrders error for user %s: %v", userID, err)
		http.Error(w, "failed to get orders", http.StatusInternalServerError)
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
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "order not found", http.StatusNotFound)
			return
		}
		log.Printf("GetOrder error for order %s: %v", orderID, err)
		http.Error(w, "failed to get order", http.StatusInternalServerError)
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
