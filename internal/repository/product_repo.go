package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ifs21014-itdel/concurrent-order-processor/internal/domain"
)

type ProductRepository interface {
	Create(ctx context.Context, product *domain.Product) error
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id int64) error
	FindByName(ctx context.Context, name string) (*domain.Product, error)
	FindById(ctx context.Context, id uint) (*domain.Product, error)
	GetAll(ctx context.Context) ([]domain.Product, error)
}

type productRepo struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) ProductRepository {
	return &productRepo{db: db}
}

// Create a new product
func (p *productRepo) Create(ctx context.Context, product *domain.Product) error {
	query := `INSERT INTO products (name, price, stock, weight,user_id)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return p.db.QueryRowContext(ctx, query,
		product.Name, product.Price, product.Stock, product.Weight, product.UserID).Scan(&product.ID)
}

// Update product data
func (p *productRepo) Update(ctx context.Context, product *domain.Product) error {
	query := `UPDATE products
	          SET name = $1, price = $2, stock = $3, weight = $4, user_id = $5
	          WHERE id = $6`
	res, err := p.db.ExecContext(ctx, query,
		product.Name, product.Price, product.Stock, product.Weight, product.UserID, product.ID)
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
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

// Delete product by ID
func (p *productRepo) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM products WHERE id = $1`
	res, err := p.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
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

func (p *productRepo) GetAll(ctx context.Context) ([]domain.Product, error) {
	query := `SELECT id, name, price, stock, weight FROM products ORDER BY id ASC`
	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.Weight); err != nil {
			return nil, fmt.Errorf("failed to scan product: %w", err)
		}
		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return products, nil
}

// Find product by name
func (p *productRepo) FindByName(ctx context.Context, name string) (*domain.Product, error) {
	query := `SELECT id, name, price, stock, weight FROM products WHERE name = $1`
	row := p.db.QueryRowContext(ctx, query, name)

	var product domain.Product
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.Weight)
	if err == sql.ErrNoRows {
		return nil, errors.New("product not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query product: %w", err)
	}
	return &product, nil
}

func (p *productRepo) FindById(ctx context.Context, id uint) (*domain.Product, error) {
	query := `SELECT id, name, price, stock, weight FROM products WHERE id = $1`
	row := p.db.QueryRowContext(ctx, query, id)

	var product domain.Product
	err := row.Scan(&product.ID, &product.Name, &product.Price, &product.Stock, &product.Weight)
	if err == sql.ErrNoRows {
		return nil, errors.New("id product not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query product: %w", err)
	}
	return &product, nil
}
