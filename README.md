# Golden Market API

A RESTful API for a fantasy-themed marketplace where users can register, authenticate, browse products, purchase items with virtual coins, and manage their inventory.

---

## 🛠 Tech Stack

- **Language:** Go 1.25
- **Database:** PostgreSQL (via pgx v5)
- **Authentication:** JWT with refresh tokens, bcrypt password hashing
- **Router:** Gorilla Mux
- **Architecture:** Clean architecture with cmd/internal structure

---

## 📋 Prerequisites

- Go 1.25 or higher
- PostgreSQL database (local or hosted — project uses [Neon](https://neon.tech) in production)
- Make (optional)
- Air (optional, for hot-reloading)

---

## 🚀 Getting Started

### Installation

1. Clone the repository
```bash
git clone https://github.com/diorshelton/golden-market-api.git
cd golden-market-api
```

2. Install dependencies
```bash
go mod download
```

3. Set up environment variables — create a `.env` file in the root directory:
```env
PORT=8080
ENVIRONMENT=development
DATABASE_URL=postgres://user:password@localhost:5432/golden_market?sslmode=disable
JWT_SECRET=your-secret-key-here
REFRESH_SECRET=your-refresh-secret-here
ACCESS_TOKEN_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=7d
ALLOWED_ORIGINS=http://localhost:5173
```

### Running the Application

```bash
make run          # Run the server
make dev          # Run with hot-reload (requires Air)
make build        # Build binary
make test         # Run tests
```

Without Make:
```bash
go run ./cmd/api
```

---

## 📁 Project Structure

```
golden-market-api/
├── cmd/
│   └── api/              # Application entry point
├── internal/
│   ├── auth/             # JWT generation, validation, and auth service
│   ├── cart/             # Cart service
│   ├── database/         # Database setup and schema migrations
│   ├── handlers/         # HTTP handlers (auth, products, cart, orders, inventory, user)
│   ├── inventory/        # Inventory service
│   ├── middleware/        # Auth, CORS, and rate limiting middleware
│   ├── models/           # Data models (user, product, cart, order, inventory)
│   ├── order/            # Order service with atomic transaction processing
│   ├── product/          # Product service
│   └── repository/       # Database access layer
├── Makefile
└── go.mod
```

---

## 🔐 Authentication

- **Access tokens** expire in 15 minutes (configurable)
- **Refresh tokens** expire in 7 days (configurable), rotated on each use
- Passwords hashed with bcrypt

---

## 📝 API Endpoints

### Public
- `GET /` — Welcome message
- `GET /health` — Status and environment info
- `GET /api/v1/products` — List all available products
- `GET /api/v1/products/{id}` — Get a single product

### Auth
- `POST /api/v1/auth/register` — Create a new account
- `POST /api/v1/auth/login` — Authenticate and receive tokens
- `POST /api/v1/auth/refresh` — Rotate access and refresh tokens
- `POST /api/v1/auth/logout` — Invalidate refresh token

### Protected
- `GET /api/v1/profile` — Get authenticated user profile and coin balance
- `GET /api/v1/cart` — Get current cart
- `POST /api/v1/cart/items` — Add item to cart
- `PUT /api/v1/cart/items/{id}` — Update cart item quantity
- `DELETE /api/v1/cart/items/{id}` — Remove item from cart
- `POST /api/v1/orders` — Place an order (atomic: deducts coins, updates stock, populates inventory)
- `GET /api/v1/orders` — Get order history
- `GET /api/v1/orders/{id}` — Get a single order
- `GET /api/v1/inventory` — Get user inventory

---

## 🗺 Roadmap

### Done
- User registration and login
- JWT access and refresh token management
- Rate limiting on auth endpoints
- Product catalog
- Virtual coin economy (5,000 starting balance)
- Shopping cart
- Atomic checkout (coins, stock, inventory, and order record — all in one transaction)
- Order history
- User inventory management
- Deployed on Render with Neon PostgreSQL

### Future Ideas
- Admin panel for product and coin management
- API documentation (Swagger)
- Product search and filtering at the API level

---

## 🤝 Contributing

This is a personal portfolio project, but feedback and suggestions are welcome. Feel free to open an issue if you notice any bugs or have ideas for improvements.

---

## 📄 License

MIT

---

## 👤 Author

**Dior Shelton** — [@diorshelton](https://github.com/diorshelton)

---

_Part of the [Golden Market](https://github.com/diorshelton/golden-market) full-stack project._
