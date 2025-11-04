package domain

import "time"

// Order menyimpan data pesanan pelanggan.
type Order struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	Status       string    `json:"status"` // pending / processed / shipped
	TotalPrice   float64   `json:"total_price"`
	ShippingCost float64   `json:"shipping_cost"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
