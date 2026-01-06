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
- **Documentation:** Swagger (swaggo)
- **Infrastructure:** Docker & Docker Compose

##  Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL
- Redis

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/akbarwjyy/go-commerce-api.git
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

5. **Access Swagger Documentation**
   ```
   http://localhost:8080/swagger/index.html
   ```

## API Documentation

### Swagger UI
Interactive API documentation is available at:
```
http://localhost:8080/swagger/index.html
```

### Swagger UI Documentation

![Swagger UI - Endpoints](docs/images/swagger2.png)

![Swagger UI - API Overview](docs/images/swagger1.png)

![Swagger UI - Request Testing](docs/images/swagger3.png)

### Endpoints Overview

#### Auth
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/v1/auth/register` | Register new user | ❌ |
| POST | `/api/v1/auth/login` | Login user | ❌ |
| POST | `/api/v1/auth/logout` | Logout (blacklist token) | ✅ |
| GET | `/api/v1/auth/me` | Get current user profile | ✅ |

#### Categories
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/v1/categories` | Get all categories | ❌ |
| GET | `/api/v1/categories/:id` | Get category by ID | ❌ |
| POST | `/api/v1/categories` | Create category | 🔐 Admin |
| PUT | `/api/v1/categories/:id` | Update category | 🔐 Admin |
| DELETE | `/api/v1/categories/:id` | Delete category | 🔐 Admin |

#### Products
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/v1/products` | Get all products | ❌ |
| GET | `/api/v1/products/:id` | Get product by ID | ❌ |
| POST | `/api/v1/products` | Create product | 🔐 Seller |
| PUT | `/api/v1/products/:id` | Update product | 🔐 Owner |
| DELETE | `/api/v1/products/:id` | Delete product | 🔐 Owner |
| PATCH | `/api/v1/products/:id/stock` | Update stock | 🔐 Owner |

#### Orders
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/v1/orders/checkout` | Create order | ✅ |
| GET | `/api/v1/orders` | Get my orders | ✅ |
| GET | `/api/v1/orders/:id` | Get order by ID | ✅ |
| PATCH | `/api/v1/orders/:id/status` | Update status | ✅ |
| POST | `/api/v1/orders/:id/cancel` | Cancel order | ✅ |

#### Payments
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/v1/payments` | Create payment | ✅ |
| GET | `/api/v1/payments` | Get my payments | ✅ |
| GET | `/api/v1/payments/:id` | Get payment by ID | ✅ |

#### Admin
| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| GET | `/api/v1/admin/orders` | Get all orders | 🔐 Admin |
| GET | `/api/v1/admin/payments` | Get all payments | 🔐 Admin |

**Legend:** ❌ Public | ✅ Authenticated | 🔐 Role-based

## User Roles

| Role | Description |
|------|-------------|
| `user` | Default role, can browse products and make orders |
| `seller` | Can manage own products |
| `admin` | Full access to all resources |
