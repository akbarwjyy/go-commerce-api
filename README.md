# Go-Commerce API

Backend E-Commerce REST API built with Golang (Gin) using **Modular Monolith** architecture.

## Architecture

This project implements a Modular Monolith pattern - mimicking microservices boundaries while running as a single application. Each business domain is isolated in its own module with dedicated models, repositories, and services.

### Modules

| Module | Description |
|--------|-------------|
| **Auth** | User registration, login, logout, JWT authentication, role management |
| **Product** | Product CRUD, categories, stock management |
| **Order** | Checkout, price calculation, order history |
| **Payment** | Payment simulation with async processing (Goroutines) |

## Tech Stack

- **Language:** Go 1.21+
- **Framework:** Gin Web Framework
- **Database:** PostgreSQL
- **ORM:** GORM v2
- **Caching:** Redis (Token Blacklist)
- **Authentication:** JWT (HMAC/RSA)
- **Infrastructure:** Docker & Docker Compose

##  Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL
- Redis

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/akbar/go-commerce-api.git
   cd go-commerce-api
   ```

2. **Setup environment variables**
   ```bash
   cp .env.example .env
   # Edit .env with your database credentials
   ```

3. **Install dependencies**
   ```bash
   go mod tidy
   ```

4. **Run the application**
   ```bash
   go run cmd/api/main.go
   ```

5. **Test the API**
   ```bash
   curl http://localhost:8080/health
   ```
## API Endpoints

### Health Check
```
GET /health
```

### Auth
```
POST /api/v1/auth/register   - Register new user
POST /api/v1/auth/login      - Login user
POST /api/v1/auth/logout     - Logout (blacklist token)
GET  /api/v1/auth/me         - Get current user profile (protected)
```

### Products (Coming Soon)
```
GET    /api/v1/products
POST   /api/v1/products
GET    /api/v1/products/:id
PUT    /api/v1/products/:id
DELETE /api/v1/products/:id
```

### Orders (Coming Soon)
```
GET  /api/v1/orders
POST /api/v1/orders
GET  /api/v1/orders/:id
```

### Payments (Coming Soon)
```
POST /api/v1/payments
GET  /api/v1/payments/:id
```

## User Roles

| Role | Description |
|------|-------------|
| `user` | Default role, can browse products and make orders |
| `seller` | Can manage own products |
| `admin` | Full access to all resources |
