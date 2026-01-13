package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/diorshelton/golden-market-api/internal/middleware"
	"github.com/diorshelton/golden-market-api/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type CartServiceInterface interface {
	AddToCart(ctx context.Context, userID, productID uuid.UUID, quantity int) error
	GetCart(ctx context.Context, userID uuid.UUID) (*models.CartSummary, error)
	UpdateCartItemQuantity(ctx context.Context, userID, cartItemID uuid.UUID, quantity int) error
	RemoveFromCart(ctx context.Context, userID, cartItemID uuid.UUID) error
}

type CartHandler struct {
	cartService CartServiceInterface
}

func NewCartHandler(service CartServiceInterface) *CartHandler {
	return &CartHandler{
		cartService: service,
	}
}

type AddToCartRequest struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type UpdateCartItemRequest struct {
	Quantity int `json:"quantity"`
}

func (h *CartHandler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	cart, err := h.cartService.GetCart(r.Context(), userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to get cart: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cart)
}

func (h *CartHandler) AddToCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	var req AddToCartRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Quantity <= 0 {
		http.Error(w, "quantity must be greater than 0", http.StatusBadRequest)
		return
	}

	productID, err := uuid.Parse(req.ProductID)
	if err != nil {
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	if err := h.cartService.AddToCart(r.Context(), userID, productID, req.Quantity); err != nil {
		http.Error(w, fmt.Sprintf("failed to add to cart: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "item added to cart",
	})
}

func (h *CartHandler) UpdateCartItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	cartItemID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "invalid cart item ID", http.StatusBadRequest)
		return
	}

	var req UpdateCartItemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Quantity <= 0 {
		http.Error(w, "quantity must be greater than 0", http.StatusBadRequest)
		return
	}

	if err := h.cartService.UpdateCartItemQuantity(r.Context(), userID, cartItemID, req.Quantity); err != nil {
		http.Error(w, fmt.Sprintf("failed to update cart item: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "cart item updated",
	})
}

func (h *CartHandler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	cartItemID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "invalid cart item ID", http.StatusBadRequest)
		return
	}

	if err := h.cartService.RemoveFromCart(r.Context(), userID, cartItemID); err != nil {
		http.Error(w, fmt.Sprintf("failed to remove from cart: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "item removed from cart",
	})
}
