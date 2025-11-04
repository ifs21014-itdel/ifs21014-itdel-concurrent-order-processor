package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ifs21014-itdel/concurrent-order-processor/internal/domain"
)

type CartItemRepository interface {
	AddCartItem(ctx context.Context, item *domain.CartItem) error
	UpdateCartItem(ctx context.Context, item *domain.CartItem) error
	DeleteCartItem(ctx context.Context, id uint) error
	GetCartItemsByCartID(ctx context.Context, cartID uint) ([]domain.CartItem, error)
	ClearCart(ctx context.Context, cartID uint) error
}

type cartItemRepo struct {
	db *sql.DB
}

func NewCartItemRepository(db *sql.DB) CartItemRepository {
	return &cartItemRepo{db: db}
}

func (r *cartItemRepo) AddCartItem(ctx context.Context, item *domain.CartItem) error {
	query := `
		INSERT INTO cart_items (cart_id, product_id, quantity, sub_total)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`
	return r.db.QueryRowContext(ctx, query,
		item.CartID,
		item.ProductID,
		item.Quantity,
		item.SubTotal,
	).Scan(&item.ID)
}

func (r *cartItemRepo) UpdateCartItem(ctx context.Context, item *domain.CartItem) error {
	query := `
		UPDATE cart_items
		SET quantity = $1, sub_total = $2, updated_at = NOW()
		WHERE id = $3
	`
	res, err := r.db.ExecContext(ctx, query, item.Quantity, item.SubTotal, item.ID)
	if err != nil {
		return fmt.Errorf("failed to update cart item: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("cart item not found")
	}
	return nil
}

func (r *cartItemRepo) DeleteCartItem(ctx context.Context, id uint) error {
	query := `DELETE FROM cart_items WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete cart item: %w", err)
	}

	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("cart item not found")
	}
	return nil
}

func (r *cartItemRepo) GetCartItemsByCartID(ctx context.Context, cartID uint) ([]domain.CartItem, error) {
	query := `
		SELECT id, cart_id, product_id, quantity, sub_total, created_at, updated_at
		FROM cart_items
		WHERE cart_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.CartItem
	for rows.Next() {
		var i domain.CartItem
		if err := rows.Scan(
			&i.ID,
			&i.CartID,
			&i.ProductID,
			&i.Quantity,
			&i.SubTotal,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	return items, nil
}

func (r *cartItemRepo) ClearCart(ctx context.Context, cartID uint) error {
	query := `DELETE FROM cart_items WHERE cart_id = $1`
	_, err := r.db.ExecContext(ctx, query, cartID)
	return err
}
