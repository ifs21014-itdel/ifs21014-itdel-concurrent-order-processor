package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ifs21014-itdel/concurrent-order-processor/internal/domain"
)

type WarehouseRepository interface {
	Create(ctx context.Context, warehouse *domain.Warehouse) error
	Update(ctx context.Context, warehouse *domain.Warehouse) error
	Delete(ctx context.Context, id int64) error
	GetAll(ctx context.Context) ([]domain.Warehouse, error)
	GetById(ctx context.Context, id uint) (*domain.Warehouse, error)
}

type warehouseRepo struct {
	db *sql.DB
}

func NewWarehouseRepository(db *sql.DB) WarehouseRepository {
	return &warehouseRepo{db: db}
}

// Create implements WarehouseRepository.
func (w *warehouseRepo) Create(ctx context.Context, warehouse *domain.Warehouse) error {
	query := `INSERT INTO warehouses(name, location)
			VALUES($1, $2) RETURNING id`
	return w.db.QueryRowContext(ctx, query, warehouse.Name, warehouse.Location).Scan(&warehouse.ID)

}

// Delete implements WarehouseRepository.
func (w *warehouseRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM warehouses WHERE id =$1`
	res, err := w.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete Warehouse : %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not check affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("product not found")
	}
	return nil
}

// GetAll implements WarehouseRepository.
func (w *warehouseRepo) GetAll(ctx context.Context) ([]domain.Warehouse, error) {
	query := `SELECT id, name, location from warehouses ORDER BY id asc`
	rows, err := w.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query warehouses: %w", err)
	}
	var warehouses []domain.Warehouse
	for rows.Next() {
		var warehouse domain.Warehouse
		if err := rows.Scan(&warehouse.ID, &warehouse.Name, &warehouse.Location); err != nil {
			return nil, fmt.Errorf("failed to scan warehouses: %w", err)
		}
		warehouses = append(warehouses, warehouse)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return warehouses, nil
}

// Update implements WarehouseRepository.
func (w *warehouseRepo) Update(ctx context.Context, warehouse *domain.Warehouse) error {
	query := `UPDATE warehouses SET name = $1, location = $2 WHERE id = $3`
	res, err := w.db.ExecContext(ctx, query, warehouse.Name, warehouse.Location, warehouse.ID)
	if err != nil {
		return fmt.Errorf("failed to update wareHouse : %w", err)
	}
	rowAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not check effected rows : %w", err)
	}
	if rowAffected == 0 {
		return errors.New("warehouse not found")
	}
	return nil
}

func (w *warehouseRepo) GetById(ctx context.Context, id uint) (*domain.Warehouse, error) {
	query := `SELECT id, name, location FROM warehouses WHERE id = $1`
	row := w.db.QueryRowContext(ctx, query, id)

	var warehouse domain.Warehouse
	err := row.Scan(&warehouse.ID, &warehouse.Name, &warehouse.Location)
	if err == sql.ErrNoRows {
		return nil, errors.New("id warehouse not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query warehouse: %w", err)
	}
	return &warehouse, nil
}
