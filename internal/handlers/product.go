package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/diorshelton/golden-market-api/internal/models"
)

type ProductServiceInterface interface {
	Create(*models.Product) error
	GetProducts() ([]*models.Product, error)
	// Delete(id uuid.UUID)
	// GetProduct(id uuid.UUID) error
	// Update(id uuid.UUID)
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
	Stock       string    `json:"stock"`
	ImageURL    string `json:"image_url"`
	Category    string `json:"category"`
}

// Validate Product input
func (r *ProductRequest) Validate() error {
	return nil
}

type ProductResponse struct {
}

func (h *ProductHandler) Create(w http.ResponseWriter, r *http.Request) {

	var req ProductRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON format%v", err), http.StatusBadRequest)
		return
	}

	priceInt, err := strconv.ParseInt(req.Price, 10, 64)
	if err != nil {
		http.Error(w, "Invalid price format", http.StatusBadRequest)
		return
	}

	stockInt, err := strconv.Atoi(req.Stock)
	if err != nil {
		http.Error(w, "Error parsing stock string", http.StatusBadRequest)
	}

	product :=&models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       models.Coins(priceInt),
		Stock:       stockInt,
		ImageURL:    req.ImageURL,
		Category:    req.Category,
	}

	// Call Product Service
	err = h.productService.Create(product)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating product: %v", err), http.StatusConflict)
		log.Printf("Error occurred: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.productService.GetProducts()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving products: %v", err), http.StatusConflict)
		log.Printf("Error occurred: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "feature not yet implemented",
	})
}

func (h *ProductHandler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "feature not yet implemented",
	})
}
