package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ifs21014-itdel/concurrent-order-processor/internal/domain"
)

type OrderRepository interface {
	UpdatePriceOrder(ctx context.Context, id uint, price float64) error
	UpdateShippingCostOrder(ctx context.Context, id uint, shipping_cost float64) error
	GetOrderByUserId(ctx context.Context, id uint) ([]domain.Order, error)
	GetOrderByUserIdAndStatus(ctx context.Context, userID uint, status string) ([]domain.Order, error)
	UpdateOrderStatus(ctx context.Context, id uint, status string) error
	Delete(ctx context.Context, id uint) error
	BeginTx(ctx context.Context) (*sql.Tx, error)
	CreateOrderTx(ctx context.Context, tx *sql.Tx, order *domain.Order) error
}

type orderRepo struct {
	db *sql.DB
}

func (o *orderRepo) UpdatePriceOrder(ctx context.Context, id uint, price float64) error {
	query := `UPDATE orders SET total_price = $1, updated_at = NOW() WHERE id = $2`
	result, err := o.db.ExecContext(ctx, query, price, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no rows updated")
	}
	return nil
}

func (o *orderRepo) UpdateShippingCostOrder(ctx context.Context, id uint, shipping_cost float64) error {
	query := `UPDATE orders SET shipping_cost = $1, updated_at = NOW() WHERE id = $2`
	result, err := o.db.ExecContext(ctx, query, shipping_cost, id)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no rows updated")
	}
	return nil
}

func (o *orderRepo) Delete(ctx context.Context, id uint) error {
	query := `DELETE FROM orders WHERE id = $1`
	res, err := o.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete order: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not check affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("order not found")
	}
	return nil
}

func (o *orderRepo) GetOrderByUserId(ctx context.Context, id uint) ([]domain.Order, error) {
	query := `SELECT id, user_id, status, total_price, shipping_cost, created_at, updated_at 
	          FROM orders WHERE user_id = $1 ORDER BY created_at DESC`
	rows, err := o.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.Status, &order.TotalPrice,
			&order.ShippingCost, &order.CreatedAt, &order.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (o *orderRepo) GetOrderByUserIdAndStatus(ctx context.Context, userID uint, status string) ([]domain.Order, error) {
	query := `SELECT id, user_id, status, total_price, shipping_cost, created_at, updated_at 
	          FROM orders WHERE user_id = $1 AND status = $2 ORDER BY created_at DESC`
	rows, err := o.db.QueryContext(ctx, query, userID, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.Status, &order.TotalPrice,
			&order.ShippingCost, &order.CreatedAt, &order.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (o *orderRepo) UpdateOrderStatus(ctx context.Context, id uint, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = NOW() WHERE id = $2`
	res, err := o.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update order status: %w", err)
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not check affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return errors.New("order not found")
	}
	return nil
}

func (o *orderRepo) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return o.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
}

func (o *orderRepo) CreateOrderTx(ctx context.Context, tx *sql.Tx, order *domain.Order) error {
	query := `INSERT INTO orders (user_id, status, total_price, shipping_cost, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id`
	return tx.QueryRowContext(ctx, query, order.UserID, order.Status, order.TotalPrice, order.ShippingCost).Scan(&order.ID)
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepo{db: db}
}
