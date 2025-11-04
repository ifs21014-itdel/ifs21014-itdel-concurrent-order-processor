# Concurrent Order Processor

A robust e-commerce backend service built with Golang using Clean Architecture principles. This project provides secure authentication with JWT and Google Authenticator (TOTP), along with powerful APIs for managing products, warehouses, carts, and orders with concurrent processing capabilities.

---

## Features

- User registration and login with secure password handling
- Optional Two-Factor Authentication (2FA) using Google Authenticator (TOTP)
- JWT authentication for protected endpoints
- Product management with user ownership
- Warehouse and stock management
- Shopping cart functionality
- Order processing with concurrent stock updates
- Transaction-based order creation with automatic stock deduction
- Real-time inventory tracking
- Clean Architecture with clear separation of concerns

---

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL
- **Authentication**: JWT + TOTP (Google Authenticator)
- **Architecture**: Clean Architecture (Handler → Usecase → Repository)

---

## Authentication Flow

The authentication system uses JWT tokens and optional Google Authenticator (TOTP) for enhanced security.

### 1. Register User

**Endpoint:**
```http
POST /api/register
```

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "mypassword",
  "name": "John Doe"
}
```

**Response:**
```json
{
  "user": {
    "id": 1,
    "email": "user@example.com",
    "name": "John Doe"
  }
}
```

---

### 2. Setup TOTP (Google Authenticator)

**Endpoint:**
```http
POST /api/totp/setup/:id
```

**Response:**
```json
{
  "secret": "JBSWY3DPEHPK3PXP",
  "otpauth_uri": "otpauth://totp/OrderProcessor:user@example.com?secret=JBSWY3DPEHPK3PXP&issuer=OrderProcessor"
}
```

**Note:** Scan the `otpauth_uri` QR code using the Google Authenticator app to generate your 6-digit verification codes.

---

### 3. Verify TOTP

**Endpoint:**
```http
POST /api/totp/verify/:id
```

**Request Body:**
```json
{
  "code": "123456"
}
```

**Response:**
```json
{
  "enabled": true
}
```

---

### 4. Login

**Endpoint:**
```http
POST /api/login
```

**Request Body (with TOTP):**
```json
{
  "email": "user@example.com",
  "password": "mypassword",
  "totp": "123456"
}
```

**Request Body (without TOTP):**
```json
{
  "email": "user@example.com",
  "password": "mypassword"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "totp_enabled": true
  }
}
```

**Important:** Use this token in the Authorization header for all protected endpoints:
```
Authorization: Bearer <your_token_here>
```

---

## Product Management

### 1. Create Product

**Endpoint:**
```http
POST /api/products/
```

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "name": "Laptop ASUS ROG",
  "description": "Gaming laptop with RTX 4060",
  "price": 15000000,
  "category": "Electronics"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "product created successfully",
  "data": {
    "id": 1,
    "user_id": 1,
    "name": "Laptop ASUS ROG",
    "description": "Gaming laptop with RTX 4060",
    "price": 15000000,
    "category": "Electronics"
  }
}
```

---

### 2. Get All Products

**Endpoint:**
```http
GET /api/products/
```

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "name": "Laptop ASUS ROG",
      "description": "Gaming laptop with RTX 4060",
      "price": 15000000,
      "category": "Electronics"
    }
  ]
}
```

---

### 3. Get Product by Name

**Endpoint:**
```http
GET /api/products/:name
```

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "id": 1,
    "name": "Laptop ASUS ROG",
    "description": "Gaming laptop with RTX 4060",
    "price": 15000000,
    "category": "Electronics"
  }
}
```

---

### 4. Update Product

**Endpoint:**
```http
PUT /api/products/:id
```

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "name": "Laptop ASUS ROG Updated",
  "description": "Gaming laptop with RTX 4070",
  "price": 18000000,
  "category": "Electronics"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "product updated successfully",
  "data": {
    "id": 1,
    "name": "Laptop ASUS ROG Updated",
    "price": 18000000
  }
}
```

---

### 5. Delete Product

**Endpoint:**
```http
DELETE /api/products/:id
```

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "message": "product deleted successfully"
}
```

---

## Warehouse Management

### 1. Create Warehouse

**Endpoint:**
```http
POST /api/warehouses/
```

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "name": "Gudang Jakarta Pusat",
  "location": "Jl. Sudirman No. 123, Jakarta"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "warehouse created successfully",
  "data": {
    "id": 1,
    "name": "Gudang Jakarta Pusat",
    "location": "Jl. Sudirman No. 123, Jakarta"
  }
}
```

---

### 2. Get All Warehouses

**Endpoint:**
```http
GET /api/warehouses/
```

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "name": "Gudang Jakarta Pusat",
      "location": "Jl. Sudirman No. 123, Jakarta"
    }
  ]
}
```

---

### 3. Update Warehouse

**Endpoint:**
```http
PUT /api/warehouses/:id
```

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "name": "Gudang Jakarta Selatan",
  "location": "Jl. TB Simatupang No. 456, Jakarta"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "warehouse updated successfully"
}
```

---

### 4. Delete Warehouse

**Endpoint:**
```http
DELETE /api/warehouses/:id
```

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "message": "warehouse deleted successfully"
}
```

---

## Warehouse Stock Management

### 1. Create Warehouse Stock

**Endpoint:**
```http
POST /api/warehouseStocks/
```

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "warehouse_id": 1,
  "product_id": 1,
  "quantity": 100
}
```

**Response:**
```json
{
  "status": "success",
  "message": "warehouse stock created successfully",
  "data": {
    "id": 1,
    "warehouse_id": 1,
    "product_id": 1,
    "quantity": 100
  }
}
```

---

### 2. Get All Warehouse Stocks

**Endpoint:**
```http
GET /api/warehouseStocks/
```

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "warehouse_id": 1,
      "product_id": 1,
      "quantity": 100
    }
  ]
}
```

---

### 3. Get Warehouse Stock by Warehouse ID

**Endpoint:**
```http
GET /api/warehouseStocks/:warehouseId
```

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "data": {
    "warehouse": {
      "id": 1,
      "name": "Gudang Jakarta Pusat"
    },
    "stocks": [
      {
        "product_id": 1,
        "product_name": "Laptop ASUS ROG",
        "quantity": 100
      }
    ]
  }
}
```

---

### 4. Update Stock Quantity

**Endpoint:**
```http
PUT /api/warehouseStocks/:id
```

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "warehouse_id": 1,
  "product_id": 1,
  "quantity": 150
}
```

**Response:**
```json
{
  "status": "success",
  "message": "stock quantity updated successfully"
}
```

---

### 5. Concurrent Update Multiple Stocks

**Endpoint:**
```http
PUT /api/warehouseStocks/concurrent
```

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
[
  {
    "warehouse_id": 1,
    "product_id": 1,
    "quantity": 150
  },
  {
    "warehouse_id": 1,
    "product_id": 2,
    "quantity": 200
  }
]
```

**Response:**
```json
{
  "status": "success",
  "message": "all stocks updated successfully",
  "results": [
    {
      "warehouse_id": 1,
      "product_id": 1,
      "success": true
    },
    {
      "warehouse_id": 1,
      "product_id": 2,
      "success": true
    }
  ]
}
```

---

### 6. Delete Warehouse Stock

**Endpoint:**
```http
DELETE /api/warehouseStocks/:id
```

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "message": "warehouse stock deleted successfully"
}
```

---

## Cart Management

### 1. Add Item to Cart

**Endpoint:**
```http
POST /api/cart/item
```

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "product_id": 1,
  "quantity": 2
}
```

**Response:**
```json
{
  "status": "success",
  "message": "item berhasil ditambahkan ke cart"
}
```

---

### 2. Get User's Carts with Items

**Endpoint:**
```http
GET /api/cart/
```

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "cart": {
        "id": 1,
        "user_id": 1,
        "created_at": "2025-01-15T10:00:00Z"
      },
      "items": [
        {
          "id": 1,
          "cart_id": 1,
          "product_id": 1,
          "product_name": "Laptop ASUS ROG",
          "quantity": 2,
          "price": 15000000,
          "sub_total": 30000000
        }
      ],
      "total": 30000000
    }
  ]
}
```

---

### 3. Delete Cart

**Endpoint:**
```http
DELETE /api/cart/:id
```

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "message": "cart berhasil dihapus"
}
```

---

### 4. Delete Cart Item

**Endpoint:**
```http
DELETE /api/cart/item/:id
```

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "message": "cart item berhasil dihapus"
}
```

---

## Order Management

### 1. Create Order from Cart

**Endpoint:**
```http
POST /api/orders/
```

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "cart_id": 1,
  "status": "pending",
  "shipping_cost": 50000
}
```

**Response:**
```json
{
  "status": "success",
  "message": "order created successfully from cart"
}
```

**Note:** 
- Valid status values: `pending`, `processed`, `shipped`, `delivered`, `cancelled`
- Order creation will automatically:
  - Calculate total price from cart items
  - Create order items
  - Deduct stock from warehouse
  - All operations are transaction-based (rollback on failure)

---

### 2. Get All User Orders

**Endpoint:**
```http
GET /api/orders/
```

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "status": "pending",
      "total_price": 30000000,
      "shipping_cost": 50000,
      "created_at": "2025-01-15T10:30:00Z",
      "updated_at": "2025-01-15T10:30:00Z"
    }
  ]
}
```

---

### 3. Get Orders by User ID and Status

**Endpoint:**
```http
GET /api/orders/status/:status
```

**Headers:**
```
Authorization: Bearer <token>
```

**Example:**
```http
GET /api/orders/status/pending
```

**Response:**
```json
{
  "status": "success",
  "data": [
    {
      "id": 1,
      "user_id": 1,
      "status": "pending",
      "total_price": 30000000,
      "shipping_cost": 50000,
      "created_at": "2025-01-15T10:30:00Z"
    }
  ]
}
```

**Valid Status Values:**
- `pending` - Order menunggu pembayaran
- `processed` - Order sedang diproses
- `shipped` - Order dalam pengiriman
- `delivered` - Order sudah sampai
- `cancelled` - Order dibatalkan

---

### 4. Update Order Status

**Endpoint:**
```http
PATCH /api/orders/:id/status
```

**Headers:**
```
Authorization: Bearer <token>
```

**Request Body:**
```json
{
  "status": "processed"
}
```

**Response:**
```json
{
  "status": "success",
  "message": "order status updated successfully"
}
```

---

### 5. Delete Order

**Endpoint:**
```http
DELETE /api/orders/:id
```

**Headers:**
```
Authorization: Bearer <token>
```

**Response:**
```json
{
  "status": "success",
  "message": "order deleted successfully"
}
```

---

## Error Responses

All endpoints may return the following error responses:

### 400 Bad Request
```json
{
  "status": "error",
  "message": "invalid JSON input",
  "error": "detailed error message"
}
```

### 401 Unauthorized
```json
{
  "error": "unauthorized"
}
```

### 404 Not Found
```json
{
  "status": "error",
  "message": "order not found"
}
```

### 500 Internal Server Error
```json
{
  "status": "error",
  "message": "failed to create order",
  "error": "detailed error message"
}
```

---

## Installation & Setup

### Prerequisites
- Go 1.21 or higher
- PostgreSQL 14+
- Git

### 1. Clone Repository
```bash
git clone https://github.com/yourusername/concurrent-order-processor.git
cd concurrent-order-processor
```

### 2. Install Dependencies
```bash
go mod download
```

### 3. Setup Database
```sql
CREATE DATABASE order_processor;
```

### 4. Configure Environment
Create `.env` file:
```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=yourpassword
DB_NAME=order_processor
JWT_SECRET=your-secret-key
PORT=8080
```

### 5. Run Migrations
```bash
go run cmd/migrate/main.go
```

### 6. Run Application
```bash
go run cmd/main.go
```

Server will start on `http://localhost:8080`

---

## Database Schema

### Users
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    name VARCHAR(255),
    totp_secret VARCHAR(255),
    totp_enabled BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Products
```sql
CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(15,2) NOT NULL,
    category VARCHAR(100),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Warehouses
```sql
CREATE TABLE warehouses (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    location TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Warehouse Stocks
```sql
CREATE TABLE warehouse_stocks (
    id SERIAL PRIMARY KEY,
    warehouse_id INTEGER REFERENCES warehouses(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(warehouse_id, product_id)
);
```

### Carts
```sql
CREATE TABLE carts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Cart Items
```sql
CREATE TABLE cart_items (
    id SERIAL PRIMARY KEY,
    cart_id INTEGER REFERENCES carts(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL,
    sub_total DECIMAL(15,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Orders
```sql
CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL,
    total_price DECIMAL(15,2) NOT NULL,
    shipping_cost DECIMAL(15,2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Order Items
```sql
CREATE TABLE order_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id) ON DELETE CASCADE,
    product_id INTEGER REFERENCES products(id),
    quantity INTEGER NOT NULL,
    sub_total DECIMAL(15,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## Architecture

This project follows Clean Architecture principles:

```

.
├── .env                               # Environment variables
├── .git/                              # Git repository metadata
├── cmd/
│   └── server/
│       └── main.go                    # Application entry point
├── config/
│   └── db.go                          # Database connection setup
├── docker-compose.yml                 # Docker services configuration
├── go.mod                             # Go module definition
├── go.sum                             # Go dependencies lock file
├── internal/
│   ├── delivery/
│   │   └── http/
│   │       ├── auth_handler.go        # Auth HTTP handler
│   │       ├── cart_handler.go        # Cart HTTP handler
│   │       ├── order_handler.go       # Order HTTP handler
│   │       ├── product_handler.go     # Product HTTP handler
│   │       ├── router.go              # HTTP router configuration
│   │       ├── wareHouse_handler.go   # Warehouse HTTP handler
│   │       └── warehouse_stock_handler.go # Warehouse stock HTTP handler
│   ├── domain/
│   │   ├── cart.go                    # Cart entity
│   │   ├── cart_item.go               # Cart item entity
│   │   ├── order.go                   # Order entity
│   │   ├── order_item.go              # Order item entity
│   │   ├── product.go                 # Product entity
│   │   ├── user.go                    # User entity
│   │   ├── wareHouse.go               # Warehouse entity
│   │   └── warehouse_stock.go         # Warehouse stock entity
│   ├── repository/
│   │   ├── cart_item_repo.go          # Cart item repository
│   │   ├── cart_repo.go               # Cart repository
│   │   ├── order_item_repo.go         # Order item repository
│   │   ├── order_repo.go              # Order repository
│   │   ├── product_repo.go            # Product repository
│   │   ├── user_repo.go               # User repository
│   │   ├── warehouse_repo.go          # Warehouse repository
│   │   └── warehouse_stock_repo.go    # Warehouse stock repository
│   ├── usecase/
│   │   ├── auth_usecase.go            # Auth business logic
│   │   ├── cart_item_usecase.go       # Cart item business logic
│   │   ├── cart_usecase.go            # Cart business logic
│   │   ├── order_usecase.go           # Order business logic
│   │   ├── product_usecase.go         # Product business logic
│   │   ├── wareHouse_usecase.go       # Warehouse business logic
│   │   └── warehouse_stock_usecase.go # Warehouse stock business logic
│   └── utils/
│       └── concurrency.go             # Concurrency utilities (goroutines, channels, etc.)
├── migrations/
│   ├── 001_create_users.sql
│   ├── 002_create_products.sql
│   ├── 003_ceate_warehouse.sql
│   ├── 004_create_warehouse_stock.sql
│   ├── 005_create_orders.sql
│   ├── 006_create_carts.sql
│   ├── 006_create_order_items.sql
│   └── 007_create_cart_items.sql
└── pkg/
├── jwt/
│   └── jwt.go                     # JWT utilities
└── totp/
└── totp.go                    # TOTP (2FA) utilities

````



## Key Features Explained

### 1. Concurrent Stock Updates
The system uses goroutines to update multiple warehouse stocks concurrently, improving performance when processing bulk operations.

### 2. Transaction-Based Order Creation
When creating an order:
- All operations are wrapped in a database transaction
- If any step fails, all changes are rolled back
- Stock is automatically deducted from warehouse
- Order items are created from cart items

### 3. JWT Authentication
- Token-based authentication for stateless API
- Token contains user ID for authorization
- Protected routes automatically validate token

### 4. Two-Factor Authentication (Optional)
- Uses TOTP (Time-based One-Time Password)
- Compatible with Google Authenticator, Authy, etc.
- Can be enabled/disabled per user

---

## Testing

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

---

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## License

This project is licensed under the MIT License - see the LICENSE file for details.

---

## Contact

For questions or support, please contact:
- Email: dediandree22@gmail.com
- GitHub: [ifs21014-itdel](https://github.com/ifs21014-itdel)

---

## Changelog

### Version 1.0.0 (2025-01-15)
- Initial release
- User authentication with JWT and TOTP
- Product management
- Warehouse and stock management
- Shopping cart functionality
- Order processing with concurrent stock updates
