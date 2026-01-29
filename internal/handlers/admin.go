package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/diorshelton/golden-market-api/internal/middleware"
	"github.com/diorshelton/golden-market-api/internal/repository"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminHandler struct {
	db            *pgxpool.Pool
	userRepo      *repository.UserRepository
	inventoryRepo *repository.InventoryRepository
}

func NewAdminHandler(db *pgxpool.Pool, userRepo *repository.UserRepository, inventoryRepo *repository.InventoryRepository) *AdminHandler {
	return &AdminHandler{
		db:            db,
		userRepo:      userRepo,
		inventoryRepo: inventoryRepo,
	}
}

type AdjustCoinsRequest struct {
	Amount int `json:"amount"` // Positive to add, negative to deduct
}

// AdjustCoins handles PATCH /api/v1/admin/users/{id}/coins
func (h *AdminHandler) AdjustCoins(w http.ResponseWriter, r *http.Request) {
	// Verify admin is authenticated
	_, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	targetUserID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	var req AdjustCoinsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Amount == 0 {
		http.Error(w, "amount must be non-zero", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// Start transaction
	tx, err := h.db.Begin(ctx)
	if err != nil {
		http.Error(w, "failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(ctx)

	if req.Amount > 0 {
		err = h.userRepo.AddCoins(ctx, tx, targetUserID, req.Amount)
	} else {
		err = h.userRepo.DeductCoins(ctx, tx, targetUserID, -req.Amount)
	}

	if err != nil {
		if contains(err.Error(), "insufficient") {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if contains(err.Error(), "not found") {
			http.Error(w, "user not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("failed to adjust coins: %v", err), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		http.Error(w, "failed to commit transaction", http.StatusInternalServerError)
		return
	}

	// Get updated user
	user, err := h.userRepo.GetUserByID(targetUserID)
	if err != nil {
		http.Error(w, "failed to get updated user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message":     "coins adjusted successfully",
		"new_balance": user.Balance,
	})
}

// ClearInventory handles DELETE /api/v1/admin/users/{id}/inventory
func (h *AdminHandler) ClearInventory(w http.ResponseWriter, r *http.Request) {
	// Verify admin is authenticated
	_, ok := middleware.GetUserID(r)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	vars := mux.Vars(r)
	targetUserID, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "invalid user ID", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	// Start transaction
	tx, err := h.db.Begin(ctx)
	if err != nil {
		http.Error(w, "failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		tx.Rollback(ctx)
	}(tx, ctx)

	if err := h.inventoryRepo.ClearByUserID(ctx, tx, targetUserID); err != nil {
		http.Error(w, fmt.Sprintf("failed to clear inventory: %v", err), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(ctx); err != nil {
		http.Error(w, "failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "inventory cleared successfully",
	})
}
