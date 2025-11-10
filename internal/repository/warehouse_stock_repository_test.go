package repository

import (
	"context"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ifs21014-itdel/concurrent-order-processor/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestWarehouseStockRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewWarehouseStockRepository(db)
	ctx := context.Background()

	t.Run("Create", func(t *testing.T) {
		stock := &domain.WarehouseStock{
			WarehouseID: 1,
			ProductID:   1,
			Quantity:    10,
		}

		mock.ExpectQuery("INSERT INTO warehouse_stock").
			WithArgs(stock.WarehouseID, stock.ProductID, stock.Quantity).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		err := repo.Create(ctx, stock)
		assert.NoError(t, err)
		assert.Equal(t, uint(1), stock.ID)
	})

	t.Run("GetAll", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{"id", "warehouse_id", "product_id", "quantity", "created_at", "updated_at"}).
			AddRow(1, 1, 1, 10, now, now).
			AddRow(2, 1, 2, 5, now, now)

		mock.ExpectQuery("SELECT id, warehouse_id, product_id, quantity, created_at, updated_at FROM warehouse_stock").
			WillReturnRows(rows)

		stocks, err := repo.GetAll(ctx)
		assert.NoError(t, err)
		assert.Len(t, stocks, 2)
	})

	t.Run("GetByWarehouseID", func(t *testing.T) {
		now := time.Now()
		rows := sqlmock.NewRows([]string{"id", "warehouse_id", "product_id", "quantity", "created_at", "updated_at"}).
			AddRow(1, 1, 1, 10, now, now)

		mock.ExpectQuery("SELECT id, warehouse_id, product_id, quantity, created_at, updated_at FROM warehouse_stock WHERE warehouse_id = \\$1").
			WithArgs(1).
			WillReturnRows(rows)

		stocks, err := repo.GetByWarehouseID(ctx, 1)
		assert.NoError(t, err)
		assert.Len(t, stocks, 1)
		assert.Equal(t, uint(1), stocks[0].ID)
	})

	t.Run("UpdateQuantity", func(t *testing.T) {
		mock.ExpectExec("UPDATE warehouse_stock SET quantity").
			WithArgs(20, 1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

		err := repo.UpdateQuantity(ctx, 1, 1, 20)
		assert.NoError(t, err)
	})

	t.Run("Delete", func(t *testing.T) {
		mock.ExpectExec("DELETE FROM warehouse_stock WHERE id = \\$1").
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.Delete(ctx, 1)
		assert.NoError(t, err)
	})

	t.Run("SafeDecreaseQuantity", func(t *testing.T) {
		tx, err := db.Begin()
		assert.NoError(t, err)

		mock.ExpectQuery("SELECT quantity FROM warehouse_stock WHERE warehouse_id = \\$1 AND product_id = \\$2 FOR UPDATE").
			WithArgs(1, 1).
			WillReturnRows(sqlmock.NewRows([]string{"quantity"}).AddRow(15))

		mock.ExpectExec("UPDATE warehouse_stock SET quantity").
			WithArgs(10, 1, 1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err = repo.SafeDecreaseQuantity(ctx, tx, 1, 1, 5)
		assert.NoError(t, err)

		_ = tx.Rollback()
	})
}
