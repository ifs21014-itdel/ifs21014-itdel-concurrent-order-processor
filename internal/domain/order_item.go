package domain

import "time"

// OrderItem menyimpan daftar item produk dalam satu order.
type OrderItem struct {
	ID        uint      `json:"id"`
	OrderID   uint      `json:"order_id"`
	ProductID uint      `json:"product_id"`
	Quantity  int32     `json:"quantity"`
	SubTotal  float64   `json:"subtotal"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
