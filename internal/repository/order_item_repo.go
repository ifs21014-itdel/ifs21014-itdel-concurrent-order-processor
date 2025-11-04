package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ifs21014-itdel/concurrent-order-processor/internal/domain"
)

type OrderItemRepository interface {
	CreateOrderItem(ctx context.Context, orderItem *domain.OrderItem) error
	UpdateOrderItem(ctx context.Context, orderItem *domain.OrderItem) error
	DeleteOrderItem(ctx context.Context, id uint) error
	GetOrderItemByIdOrder(ctx context.Context, id uint) ([]domain.OrderItem, error)
	CreateOrderItemTx(ctx context.Context, tx *sql.Tx, item *domain.OrderItem) error
}

type orderItemRepo struct {
	db *sql.DB
}

// CreateOrderItem implements OrderItemRepository.
func (o *orderItemRepo) CreateOrderItem(ctx context.Context, orderItem *domain.OrderItem) error {
	query := `INSERT INTO order_items (order_id, product_id, quantity, sub_total)
	VALUES($1, $2, $3, $4) RETURNING id`
	return o.db.QueryRowContext(ctx, query, orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.SubTotal).Scan(&orderItem.ID)
}

// DeleteOrderItem implements OrderItemRepository.
func (o *orderItemRepo) DeleteOrderItem(ctx context.Context, id uint) error {
	query := `DELETE FROM order_items WHERE id = $1`
	res, err := o.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete item :%w", err)
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not check affected rows :%w", err)
	}

	if rowAffected == 0 {
		return errors.New("item not found")
	}
	return nil
}

// GetOrderItemByIdOrder implements OrderItemRepository.
func (o *orderItemRepo) GetOrderItemByIdOrder(ctx context.Context, id uint) ([]domain.OrderItem, error) {
	query := `SELECT id, order_id, product_id,quantity, sub_total FROM order_items WHERE order_id = $1`
	rows, err := o.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orderItems []domain.OrderItem
	for rows.Next() {
		var orderItem domain.OrderItem
		if err := rows.Scan(&orderItem.ID, orderItem.OrderID, orderItem.ProductID, orderItem.Quantity, orderItem.SubTotal); err != nil {
			return nil, err
		}
		orderItems = append(orderItems, orderItem)
	}
	return orderItems, nil
}

// UpdateOrderItem implements OrderItemRepository.
func (o *orderItemRepo) UpdateOrderItem(ctx context.Context, orderItem *domain.OrderItem) error {
	query := `UPDATE order_items 
	SET quantity = $1, sub_total = $2 WHERE id = $3`
	res, err := o.db.ExecContext(ctx, query, orderItem.Quantity, orderItem.SubTotal, orderItem.ID)
	if err != nil {
		return fmt.Errorf("failed to update ordeer item : %w", err)
	}

	rowAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not check affected rows :%w", err)
	}
	if rowAffected == 0 {
		return errors.New("order item not found")
	}
	return nil
}
func (r *orderItemRepo) CreateOrderItemTx(ctx context.Context, tx *sql.Tx, item *domain.OrderItem) error {
	query := `
		INSERT INTO order_items (order_id, product_id, quantity, sub_total, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id
	`
	err := tx.QueryRowContext(ctx, query, item.OrderID, item.ProductID, item.Quantity, item.SubTotal).Scan(&item.ID)
	if err != nil {
		return fmt.Errorf("failed to insert order item (tx): %w", err)
	}
	return nil
}

func NewOrderItemRepository(db *sql.DB) OrderItemRepository {
	return &orderItemRepo{db: db}
}
