package domain

import "time"

type Cart struct {
	ID          uint      `json:"id"`
	UserID      uint      `json:"user_id"`
	WarehouseID uint      `json:"warehouse_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
