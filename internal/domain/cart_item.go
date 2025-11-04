package domain

import "time"

type CartItem struct {
	ID        uint      `json:"id"`
	CartID    uint      `json :"cart_id"`
	ProductID uint      `json:"product_id"`
	Quantity  int32     `json :"quantity"`
	SubTotal  float64   `json:"sub_total"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
