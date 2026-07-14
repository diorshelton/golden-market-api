package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/diorshelton/golden-market-api/internal/auth"
	"github.com/diorshelton/golden-market-api/internal/cart"
	"github.com/diorshelton/golden-market-api/internal/config"
	"github.com/diorshelton/golden-market-api/internal/database"
	"github.com/diorshelton/golden-market-api/internal/handlers"
	"github.com/diorshelton/golden-market-api/internal/inventory"
	"github.com/diorshelton/golden-market-api/internal/middleware"
	"github.com/diorshelton/golden-market-api/internal/order"
	"github.com/diorshelton/golden-market-api/internal/product"
	"github.com/diorshelton/golden-market-api/internal/repository"
	"github.com/gorilla/mux"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Set up databases
	database, err := database.SetupDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	// Create repositories
	tokenRepo := repository.NewRefreshTokenRepository(database)
	userRepo := repository.NewUserRepository(database)
	productRepo := repository.NewProductRepository(database)
	cartRepo := repository.NewCartRepository(database)
	orderRepo := repository.NewOrderRepository(database)
	orderItemRepo := repository.NewOrderItemRepository(database)
	inventoryRepo := repository.NewInventoryRepository(database)

	// Create  auth service
	authService := auth.NewAuthService(
		userRepo,
		tokenRepo,
		cfg.JWTSecret,
		cfg.RefreshSecret,
		cfg.AccessTokenExpiry,
		cfg.RefreshTokenExpiry,
	)

	// Create product service
	productService := product.NewProductService(productRepo)

	// Create cart service
	cartService := cart.NewCartService(cartRepo, productRepo)

	// Create order service
	orderService := order.NewOrderService(
		database,
		orderRepo,
		orderItemRepo,
		inventoryRepo,
		userRepo,
		productRepo,
		cartRepo,
	)

	// Create inventory service
	inventoryService := inventory.NewInventoryService(inventoryRepo)

	// Create handlers
	authHandler := handlers.NewAuthHandler(authService, cfg)
	userHandler := handlers.NewUserHandler(userRepo)
	productHandler := handlers.NewProductHandler(productService)
	cartHandler := handlers.NewCartHandler(cartService)
	orderHandler := handlers.NewOrderHandler(orderService)
	inventoryHandler := handlers.NewInventoryHandler(inventoryService)
	//adminHandler := handlers.NewAdminHandler(database, userRepo, inventoryRepo)

	// Create router
	r := mux.NewRouter()

	//Apply CORS middleware
	corsMiddleware := middleware.CORS(cfg)
	r.Use(corsMiddleware)

	// --- Public API Endpoints --
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Welcome to Golden Market!\n")
	})

	// Health check
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status":      "ok",
			"port":        cfg.Port,
			"environment": cfg.Environment,
		})
	}).Methods("GET")

	// --- Product API Endpoints (Public - Read Only) --
	r.HandleFunc("/api/v1/products", productHandler.GetProducts).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/v1/products/{id}", productHandler.GetProduct).Methods("GET", "OPTIONS")

	// --- Auth API Endpoints (rate limited) ---
	authRouter := r.PathPrefix("/api/v1/auth").Subrouter()
	authRouter.Use(corsMiddleware)       // Apply CORS to Subrouter
	authRouter.Use(middleware.RateLimit) // Apply ratelimiting

	authRouter.HandleFunc("/register", authHandler.Register).Methods("POST", "OPTIONS")
	authRouter.HandleFunc("/login", authHandler.Login).Methods("POST", "OPTIONS")
	authRouter.HandleFunc("/guest-login", authHandler.GuestLogin).Methods("POST", "OPTIONS")
	authRouter.HandleFunc("/refresh", authHandler.Refresh).Methods("POST", "OPTIONS")
	authRouter.HandleFunc("/logout", authHandler.Logout).Methods("POST", "OPTIONS")

	// --- Protected routes ---
	protected := r.PathPrefix("/api/v1").Subrouter()
	protected.Use(corsMiddleware) // Apply CORS to Subrouter
	protected.Use(middleware.Auth(authService))
	protected.HandleFunc("/profile", userHandler.Profile).Methods("GET", "OPTIONS")

	// Product write operations (protected)
	protected.HandleFunc("/products", productHandler.Create).Methods("POST", "OPTIONS")
	protected.HandleFunc("/products/{id}", productHandler.Update).Methods("PUT", "PATCH", "OPTIONS")
	protected.HandleFunc("/products/{id}", productHandler.Delete).Methods("DELETE", "OPTIONS")

	// Cart operations (protected)
	protected.HandleFunc("/cart", cartHandler.GetCart).Methods("GET", "OPTIONS")
	protected.HandleFunc("/cart/items", cartHandler.AddToCart).Methods("POST", "OPTIONS")
	protected.HandleFunc("/cart/items/{id}", cartHandler.UpdateCartItem).Methods("PUT", "PATCH", "OPTIONS")
	protected.HandleFunc("/cart/items/{id}", cartHandler.RemoveFromCart).Methods("DELETE", "OPTIONS")

	// Order operations (protected)
	protected.HandleFunc("/orders", orderHandler.CreateOrder).Methods("POST", "OPTIONS")
	protected.HandleFunc("/orders", orderHandler.GetOrders).Methods("GET", "OPTIONS")
	protected.HandleFunc("/orders/{id}", orderHandler.GetOrder).Methods("GET", "OPTIONS")

	// Inventory operations (protected)
	protected.HandleFunc("/inventory", inventoryHandler.GetInventory).Methods("GET", "OPTIONS")
	/*
		TODO IMPLEMENT ADMIN FEATURES

		// Admin operations (protected)
		protected.HandleFunc("/admin/users/{id}/coins", adminHandler.AdjustCoins).Methods("PATCH", "OPTIONS")
		protected.HandleFunc("/admin/users/{id}/inventory", adminHandler.ClearInventory).Methods("DELETE", "OPTIONS")
	*/
	// Start server
	addr := ":" + cfg.Port
	log.Printf("Server starting on port %s", cfg.Port)
	log.Printf("Environment: %s", cfg.Environment)

	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatal(err)
	}
}
