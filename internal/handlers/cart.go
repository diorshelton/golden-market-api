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
		log.Printf("Failed to get cart for user %s: %v", userID, err)
		http.Error(w, "failed to get cart", http.StatusInternalServerError)
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
		log.Printf("Failed to add to cart for user %s: %v", userID, err)
		// Check for user-friendly errors
		errMsg := err.Error()
		if strings.Contains(errMsg, "insufficient stock") || strings.Contains(errMsg, "not found") {
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}
		http.Error(w, "failed to add to cart", http.StatusInternalServerError)
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
		log.Printf("Failed to update cart item %s for user %s: %v", cartItemID, userID, err)
		errMsg := err.Error()
		if strings.Contains(errMsg, "insufficient stock") || strings.Contains(errMsg, "not found") {
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}
		http.Error(w, "failed to update cart item", http.StatusInternalServerError)
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
		log.Printf("Failed to remove cart item %s for user %s: %v", cartItemID, userID, err)
		http.Error(w, "failed to remove from cart", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "item removed from cart",
	})
}
