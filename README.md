# Golden Market API

Backend for Golden Market, a fantasy-themed marketplace app. Handles registration and auth, a product catalog, a virtual coin economy, cart and checkout, and inventory — all backed by Postgres.

## Stack

- Go 1.25
- PostgreSQL via pgx v5
- JWT auth (access + refresh tokens), bcrypt for passwords
- Gorilla Mux
- `cmd`/`internal` layout, service/handler/repository layers per domain

## Prerequisites

- Go 1.25+
- A running PostgreSQL instance (local, Docker, or hosted — production runs on Neon)
- Make (optional, but assumed below)
- Air, for hot-reload (optional): `go install github.com/air-verse/air@latest`

Windows: the commands below assume a Unix shell. Use WSL and everything works as written.

## Getting started

Clone the repo and pull dependencies:

```bash
git clone https://github.com/diorshelton/golden-market-api.git
cd golden-market-api
go mod download
```

Copy the example env file:

```bash
cp .env.example .env
```

`.env.example` looks like this:

```env
PORT=8080
ENVIRONMENT=development
DATABASE_URL=postgres://youruser@localhost:5432/golden_market?sslmode=disable
JWT_SECRET=replace-with-a-random-string
REFRESH_SECRET=replace-with-a-different-random-string
ACCESS_TOKEN_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=168h
ALLOWED_ORIGINS=http://localhost:5173
```

Generate a value for `JWT_SECRET` and another for `REFRESH_SECRET`:

```bash
openssl rand -hex 32
```

Run that twice and drop the two results into `.env`.

Set `DATABASE_URL` to match your local Postgres. If you're on a default Homebrew install (trust auth, role = your OS username), the value in `.env.example` works as-is once the database exists:

```bash
createdb golden_market
```

If you're using Air for hot-reload, create `.air.toml` in the repo root:

```toml
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/main ./cmd/api"
bin = "./tmp/main"
include_ext = ["go"]
exclude_dir = ["assets", "tmp", "vendor", "testdata", "bin"]
delay = 1000
stop_on_error = true
log = "build-errors.log"

[color]
main = "magenta"
watcher = "cyan"
build = "yellow"
runner = "green"

[log]
time = false

[misc]
clean_on_exit = false
```

Then run the server:

```bash
make dev    # hot-reload, requires Air
make run    # plain go run
```

## Project layout

```
golden-market-api/
├── cmd/api/          entry point
├── internal/
│   ├── auth/          JWT generation/validation, auth service
│   ├── cart/          cart service
│   ├── database/      schema setup and migrations
│   ├── handlers/      HTTP handlers
│   ├── inventory/     inventory service
│   ├── middleware/    auth, CORS, rate limiting
│   ├── models/        data models
│   ├── order/         checkout, atomic transaction processing
│   ├── product/       product service
│   └── repository/    database access layer
├── Makefile
└── go.mod
```

## Auth

Access tokens expire in 15 minutes, refresh tokens in 7 days (both configurable via `.env`), rotated on each use. Passwords are hashed with bcrypt.

## API endpoints

### Public
- `GET /` — welcome message
- `GET /health` — status and environment info
- `GET /api/v1/products` — list products
- `GET /api/v1/products/{id}` — get one product

### Auth
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/login`
- `POST /api/v1/auth/guest-login` — logs into a single shared guest account; its cart, inventory, and orders reset on every login
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`

### Protected (bearer token required)
- `GET /api/v1/profile`
- `GET /api/v1/cart`
- `POST /api/v1/cart/items`
- `PUT /api/v1/cart/items/{id}`
- `DELETE /api/v1/cart/items/{id}`
- `POST /api/v1/orders` — atomic checkout: deducts coins, updates stock, populates inventory
- `GET /api/v1/orders`
- `GET /api/v1/orders/{id}`
- `GET /api/v1/inventory`

## Status

Registration, login, guest login, product catalog, cart, atomic checkout, order history, and inventory are all live, deployed on Render against Neon Postgres.

Not built yet: admin panel, API docs, server-side product search/filtering.

## Contributing

This is a personal portfolio project, but feedback and suggestions are welcome. Feel free to open an issue if you notice any bugs or have ideas for improvements.

## License

MIT

## Author

Dior Shelton — [@diorshelton](https://github.com/diorshelton)

---

Part of the [Golden Market](https://github.com/diorshelton/golden-market) project.
