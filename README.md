# Golden Market API

A RESTful API for an online farmers market simulation where users can purchase items from vendors using virtual coins and manage their inventory.

## 🎯 Project Overview

Golden Market is a portfolio project demonstrating full-stack API development with Go. The platform simulates a farmers market economy where users can:

- Create accounts and authenticate securely
- Browse products from various vendors
- Purchase items using virtual coins
- Manage personal inventory
- (More features in development)

## 🛠 Tech Stack

- **Language:** Go
- **Database:** SQLite (development), PostgreSQL (production)
- **Authentication:** JWT with refresh tokens
- **Architecture:** Clean architecture with cmd/internal structure

## 📋 Prerequisites

- Go 1.21 or higher
- Make (optional, for using Makefile commands)
- Air (optional, for hot-reloading in development)

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

3. Set up environment variables

Create a `.env` file in the root directory:

```env
JWT_SECRET=your-secret-key-here
REFRESH_SECRET=your-refresh-secret-here
ACCESS_TOKEN_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=7d
```

### Running the Application

**Using Make (recommended):**

```bash
# Run the server
make run

# Run with hot-reload
make dev

# Run tests
make test

# Build binary
make build
```

**Without Make:**

```bash
# Run the server
go run ./cmd/api

# Run tests
go test ./...

# Build binary
go build -o bin/api ./cmd/api
```

## 📁 Project Structure

```
golden-market-api/
├── cmd/
│   └── api/           # Application entry point
├── internal/
│   ├── auth/          # Authentication logic
│   ├── database/      # Database connection
│   ├── handlers/      # HTTP handlers
│   ├── middleware/    # HTTP middleware
│   └── models/        # Data models
├── .air.toml          # Hot-reload configuration
├── Makefile           # Build commands
└── go.mod
```

## 🔐 Authentication

The API uses JWT-based authentication with access and refresh tokens:

- **Access tokens** expire in 15 minutes (configurable)
- **Refresh tokens** expire in 7 days (configurable)
- Passwords are hashed using bcrypt

## 🗺 Roadmap

### ✅ Completed

- [x] Project structure setup
- [x] User authentication (register/login)
- [x] JWT token management
- [x] Basic user models

### 🚧 In Progress

- [ ] Vendor system
- [ ] Product catalog
- [ ] Virtual coin economy

### 📅 Planned Features

- [ ] Purchase transactions
- [ ] User inventory management
- [ ] Vendor dashboard
- [ ] Product reviews/ratings
- [ ] Search and filtering
- [ ] PostgreSQL migration
- [ ] API documentation (Swagger)
- [ ] Rate limiting
- [ ] Comprehensive test coverage

## 🧪 Testing

Run the test suite:

```bash
make test

# With verbose output
make test-verbose

# With coverage
make test-coverage
```

## 📝 API Endpoints

### Authentication

- `POST /auth/register` - Create new user account
- `POST /auth/login` - Authenticate and receive tokens
- `POST /auth/refresh` - Refresh access token

_(More endpoints coming soon)_

## 🤝 Contributing

This is a personal portfolio project, but feedback and suggestions are welcome! Feel free to open an issue if you notice any bugs or have ideas for improvements.

## 📄 License

This project is open source and available under the [MIT License](LICENSE).

## 👤 Author

**Dior Shelton**

- GitHub: [@diorshelton](https://github.com/diorshelton)

---

_This project is actively being developed as part of my portfolio. Check back for updates!_
