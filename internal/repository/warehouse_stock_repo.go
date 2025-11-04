package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ifs21014-itdel/concurrent-order-processor/internal/domain"
)

type WarehouseStockRepository interface {
	Create(ctx context.Context, stock *domain.WarehouseStock) error
	GetAll(ctx context.Context) ([]domain.WarehouseStock, error)
	GetByWarehouseID(ctx context.Context, warehouseID uint) ([]domain.WarehouseStock, error)
	UpdateQuantity(ctx context.Context, warehouseID uint, productID uint, quantity int32) error
	Delete(ctx context.Context, stockID uint) error
	SafeDecreaseQuantity(ctx context.Context, tx *sql.Tx, warehouseID uint, productID uint, qtyToDecrease int32) error
}

type warehouseStockRepo struct {
	db *sql.DB
}

func NewWarehouseStockRepository(db *sql.DB) WarehouseStockRepository {
	return &warehouseStockRepo{db: db}
}

func (r *warehouseStockRepo) Create(ctx context.Context, stock *domain.WarehouseStock) error {
	query := `
		INSERT INTO warehouse_stock (warehouse_id, product_id, quantity, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
		RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query, stock.WarehouseID, stock.ProductID, stock.Quantity).Scan(&stock.ID)
	if err != nil {
		return fmt.Errorf("failed to create warehouse stock: %w", err)
	}
	return nil
}

func (r *warehouseStockRepo) GetAll(ctx context.Context) ([]domain.WarehouseStock, error) {
	query := `SELECT id, warehouse_id, product_id, quantity, created_at, updated_at FROM warehouse_stock`
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query warehouse stocks: %w", err)
	}
	defer rows.Close()

	var stocks []domain.WarehouseStock
	for rows.Next() {
		var s domain.WarehouseStock
		if err := rows.Scan(&s.ID, &s.WarehouseID, &s.ProductID, &s.Quantity, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan warehouse stock: %w", err)
		}
		stocks = append(stocks, s)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating warehouse stocks: %w", err)
	}
	return stocks, nil
}

func (r *warehouseStockRepo) GetByWarehouseID(ctx context.Context, warehouseID uint) ([]domain.WarehouseStock, error) {
	query := `SELECT id, warehouse_id, product_id, quantity, created_at, updated_at FROM warehouse_stock WHERE warehouse_id = $1`
	rows, err := r.db.QueryContext(ctx, query, warehouseID)
	if err != nil {
		return nil, fmt.Errorf("failed to query warehouse stock by warehouse_id: %w", err)
	}
	defer rows.Close()

	var stocks []domain.WarehouseStock
	for rows.Next() {
		var s domain.WarehouseStock
		if err := rows.Scan(&s.ID, &s.WarehouseID, &s.ProductID, &s.Quantity, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan warehouse stock: %w", err)
		}
		stocks = append(stocks, s)
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating warehouse stock rows: %w", err)
	}
	return stocks, nil
}

func (r *warehouseStockRepo) UpdateQuantity(ctx context.Context, warehouseID uint, productID uint, quantity int32) error {
	query := `UPDATE warehouse_stock SET quantity = $1, updated_at = NOW() WHERE warehouse_id = $2 and product_id =$3`
	result, err := r.db.ExecContext(ctx, query, quantity, warehouseID, productID)

	if err != nil {
		return fmt.Errorf("failed to update warehouse stock quantity: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("no warehouse stock updated for warehouse_id=%d product_id=%d", warehouseID, productID)
	}

	return nil
}

func (r *warehouseStockRepo) Delete(ctx context.Context, stockID uint) error {
	query := `DELETE FROM warehouse_stock WHERE id = $1`
	result, err := r.db.ExecContext(ctx, query, stockID)
	if err != nil {
		return fmt.Errorf("failed to delete warehouse stock: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no warehouse stock deleted")
	}
	return nil
}

func (r *warehouseStockRepo) SafeDecreaseQuantity(ctx context.Context, tx *sql.Tx, warehouseID uint, productID uint, qtyToDecrease int32) error {

	querySelect := `SELECT quantity FROM warehouse_stock WHERE warehouse_id = $1 AND product_id = $2 FOR UPDATE`
	var currentQty int32
	err := tx.QueryRowContext(ctx, querySelect, warehouseID, productID).Scan(&currentQty)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("no stock found for warehouse_id=%d and product_id=%d", warehouseID, productID)
		}
		return fmt.Errorf("failed to get current stock: %w", err)
	}

	if currentQty < qtyToDecrease {
		return fmt.Errorf("not enough stock for product_id=%d in warehouse_id=%d", productID, warehouseID)
	}

	newQty := currentQty - qtyToDecrease
	queryUpdate := `UPDATE warehouse_stock SET quantity = $1, updated_at = NOW() WHERE warehouse_id = $2 AND product_id = $3`
	_, err = tx.ExecContext(ctx, queryUpdate, newQty, warehouseID, productID)
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}

	return nil
}
