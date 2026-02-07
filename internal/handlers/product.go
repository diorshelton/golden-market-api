package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/diorshelton/golden-market-api/internal/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type ProductServiceInterface interface {
	Create(*models.Product) error
	GetProducts() ([]*models.Product, error)
	GetProduct(id uuid.UUID) (*models.Product, error)
	Update(id uuid.UUID)
	Delete(id uuid.UUID)
}

type ProductHandler struct {
	productService ProductServiceInterface
}

func NewProductHandler(service ProductServiceInterface) *ProductHandler {
	return &ProductHandler{
		productService: service,
	}
}

type ProductRequest struct {
	Name        string `json:"product_name"`
	Description string `json:"product_description"`
	Price       string `json:"price"`
	Stock       string `json:"stock"`
	ImageURL    string `json:"image_url"`
	Category    string `json:"category"`
}

// Validate Product input
func (r *ProductRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("product name is required")
	}
	if r.Price == "" {
		return errors.New("price is required")
	}
	price, err := strconv.ParseInt(r.Price, 10, 64)
	if err != nil {
		return errors.New("invalid price format")
	}
	if price <= 0 {
		return errors.New("price must be greater than 0")
	}
	if r.Stock != "" {
		stock, err := strconv.Atoi(r.Stock)
		if err != nil {
			return errors.New("invalid stock format")
		}
		if stock < 0 {
			return errors.New("stock cannot be negative")
		}
	}
	return nil
}

type ProductResponse struct {
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req ProductRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON format: %v", err), http.StatusBadRequest)
		return
	}

	if err := req.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	priceInt, err := strconv.ParseInt(req.Price, 10, 64)
	if err != nil {
		http.Error(w, "Invalid price format", http.StatusBadRequest)
		return
	}

	stockInt := 0
	if req.Stock != "" {
		stockInt, err = strconv.Atoi(req.Stock)
		if err != nil {
			http.Error(w, "Error parsing stock string", http.StatusBadRequest)
			return
		}
	}

	product := &models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       models.Coins(priceInt),
		Stock:       stockInt,
		ImageURL:    req.ImageURL,
		Category:    req.Category,
		IsAvailable: true,
	}

	// Call Product Service
	err = h.productService.Create(product)
	if err != nil {
		log.Printf("Error creating product: %v", err)
		http.Error(w, "failed to create product", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.productService.GetProducts()
	if err != nil {
		log.Printf("Error retrieving products: %v", err)
		http.Error(w, "failed to retrieve products", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := uuid.Parse(vars["id"])
	if err != nil {
		http.Error(w, "invalid product ID", http.StatusBadRequest)
		return
	}

	product, err := h.productService.GetProduct(id)
	if err != nil {
		http.Error(w, "product not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Update endpoint not yet implemented",
	})
}

func (h *ProductHandler) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Delete endpoint not yet implemented",
	})
}
