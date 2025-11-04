package domain

import "time"

type WarehouseStock struct {
	ID          uint      `json:"id"`
	WarehouseID uint      `json:"warehouse_id"`
	ProductID   uint      `json:"product_id"`
	Quantity    int32     `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
