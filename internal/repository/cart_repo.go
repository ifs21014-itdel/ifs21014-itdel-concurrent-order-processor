package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ifs21014-itdel/concurrent-order-processor/internal/domain"
)

type CartRepository interface {
	CreateCart(ctx context.Context, cart *domain.Cart) error
	DeleteCart(ctx context.Context, id uint) error
	GetCartByUserId(ctx context.Context, userId uint) ([]domain.Cart, error)
	FindByUserAndWarehouse(ctx context.Context, userId uint, warehouseID uint) (*domain.Cart, error)
}

type cartRepo struct {
	db *sql.DB
}

func (c *cartRepo) FindByUserAndWarehouse(ctx context.Context, userId uint, warehouseID uint) (*domain.Cart, error) {
	query := `SELECT id, user_id, warehouse_id, created_at FROM carts WHERE user_id = $1 AND warehouse_id = $2 LIMIT 1`
	row := c.db.QueryRowContext(ctx, query, userId, warehouseID)

	var cart domain.Cart
	err := row.Scan(&cart.ID, &cart.UserID, &cart.WarehouseID, &cart.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("chart not found")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to query cart :%w", err)
	}

	return &cart, nil
}

// AddCart implements CartRepository.
func (c *cartRepo) CreateCart(ctx context.Context, cart *domain.Cart) error {
	query := `INSERT INTO carts (user_id , warehouse_id) VALUES ($1, $2) RETURNING id`
	return c.db.QueryRowContext(ctx, query, cart.UserID, cart.WarehouseID).Scan(&cart.ID)
}

// DeleteCart implements CartRepository.
func (c *cartRepo) DeleteCart(ctx context.Context, id uint) error {
	query := `DELETE FROM carts WHERE id = $1`
	res, err := c.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete cart :%w", err)
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not check affecred rows : %w", err)
	}
	if rowAffected == 0 {
		return fmt.Errorf("cart not found")
	}
	return nil
}

// GetCartByUserId implements CartRepository.
func (c *cartRepo) GetCartByUserId(ctx context.Context, userId uint) ([]domain.Cart, error) {
	query := `SELECT * FROM carts WHERE user_id = $1`
	rows, err := c.db.Query(query, userId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var carts []domain.Cart
	for rows.Next() {
		var cart domain.Cart
		err := rows.Scan(&cart.ID, &cart.UserID, &cart.WarehouseID, &cart.CreatedAt, &cart.UpdatedAt)
		if err != nil {
			return nil, err
		}
		carts = append(carts, cart)
	}
	return carts, nil
}

func NewCartRepository(db *sql.DB) CartRepository {
	return &cartRepo{db: db}
}
