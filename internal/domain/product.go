package domain

import "time"

type Product struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	UserID    uint      `json:"user_id"`
	Price     float64   `json:"price"`
	Stock     int32     `json:"stock"`            // optional, bisa dihapus jika stok per gudang
	Weight    *float64  `json:"weight,omitempty"` // bisa NULL
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
