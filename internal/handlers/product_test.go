package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/diorshelton/golden-market-api/internal/models"
	"github.com/google/uuid"
)

type MockProductService struct {
	CreateFunc      func(*models.Product) error
	GetProductsFunc func() ([]*models.Product, error)
	GetProductFunc  func(id uuid.UUID) ([]*models.Product, error)
	UpdateFunc      func(id uuid.UUID)
	DeleteFunc      func(id uuid.UUID) error
}

/*
Create(*models.Product) error
GetProducts() ([]*models.Product, error)
GetProduct(id uuid.UUID) (*models.Product, error)
Update(id uuid.UUID)
Delete(id uuid.UUID)
*/
func (m *MockProductService) Create(product *models.Product) error {
	return m.CreateFunc(product)
}

func (m *MockProductService) GetProducts() ([]*models.Product, error) {
	return m.GetProductsFunc()
}

func (m *MockProductService) GetProduct(id uuid.UUID) (*models.Product, error) {
	if m.GetProductFunc != nil {
		return nil, nil
	}
	return nil, nil
}

func (m *MockProductService) Update(id uuid.UUID) {
	if m.UpdateFunc != nil {
		m.UpdateFunc(id)
	}
}

func (m *MockProductService) Delete(id uuid.UUID) {
	if m.DeleteFunc != nil {
		m.DeleteFunc(id)
	}
}

func TestCreateHandler(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]string
		mockCreate     func(*models.Product) error
		expectedStatus int
		checkResponse  func(t *testing.T, resp *httptest.ResponseRecorder)
	}{
		{
			name: "successful product creation",
			requestBody: map[string]string{
				"product_name":        "coffee mug",
				"product_description": "artesenal ceramic mug made by hand",
				"price":               "45",
				"stock":               "12",
				"image_url":           "https://example.com/coffee-mug.jpg",
				"category":            "Misc. Items",
			},
			mockCreate: func(product *models.Product) error {
				return nil
			},
			expectedStatus: http.StatusCreated,
			checkResponse: func(t *testing.T, resp *httptest.ResponseRecorder) {
				var response models.Product
				err := json.NewDecoder(resp.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Failed to decode response: %v", err)
				}
				if response.Name != "coffee mug" {
					t.Fatalf("Failed to create product")
				}
			},
		},
		{
			name: "invalid price format",
			requestBody: map[string]string{
				"product_name":        "coffee mug",
				"product_description": "artesenal ceramic mug made by hand",
				"price":               "fourty-five",
				"stock":               "12",
				"image_url":           "https://example.com/coffee-mug.jpg",
				"category":            "Misc. Items",
			},
			mockCreate:     nil,
			expectedStatus: http.StatusBadRequest,
			checkResponse:  nil,
		},
		{
			name: "service returns error",
			requestBody: map[string]string{
				"product_name":        "coffee mug",
				"product_description": "artesenal ceramic mug made by hand",
				"price":               "45",
				"stock":               "12",
				"image_url":           "https://example.com/coffee-mug.jpg",
				"category":            "Misc. Items",
			},
			mockCreate: func(p *models.Product) error {
				return errors.New("database error")
			},
			expectedStatus: http.StatusConflict,
			checkResponse:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &MockProductService{
				CreateFunc: tt.mockCreate,
			}
			handler := NewProductHandler(mockService)

			jsonBody, _ := json.Marshal(tt.requestBody)
			req := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			handler.Create(rr, req)

			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			if tt.checkResponse != nil {
				tt.checkResponse(t, rr)
			}
		})
	}
}
