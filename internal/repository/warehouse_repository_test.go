package repository

import (
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ifs21014-itdel/concurrent-order-processor/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestWarehouseRepository_Create(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewWarehouseRepository(db)
	ctx := context.Background()

	warehouse := &domain.Warehouse{
		Name:     "Warehouse A",
		Location: "Jakarta",
	}

	// Mock behavior
	mock.ExpectQuery("INSERT INTO warehouses").
		WithArgs(warehouse.Name, warehouse.Location).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	err := repo.Create(ctx, warehouse)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), warehouse.ID)
}

func TestWarehouseRepository_GetAll(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewWarehouseRepository(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "name", "location"}).
		AddRow(1, "Warehouse A", "Jakarta").
		AddRow(2, "Warehouse B", "Bandung")

	mock.ExpectQuery("SELECT id, name, location from warehouses").
		WillReturnRows(rows)

	result, err := repo.GetAll(ctx)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, "Warehouse B", result[1].Name)
}

func TestWarehouseRepository_GetById(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewWarehouseRepository(db)
	ctx := context.Background()

	rows := sqlmock.NewRows([]string{"id", "name", "location"}).
		AddRow(1, "Warehouse A", "Jakarta")

	mock.ExpectQuery("SELECT id, name, location FROM warehouses WHERE id =").
		WithArgs(1).
		WillReturnRows(rows)

	w, err := repo.GetById(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, "Warehouse A", w.Name)
}

func TestWarehouseRepository_Update(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewWarehouseRepository(db)
	ctx := context.Background()

	warehouse := &domain.Warehouse{
		ID:       1,
		Name:     "Warehouse Updated",
		Location: "Jakarta",
	}

	mock.ExpectExec("UPDATE warehouses SET name =").
		WithArgs(warehouse.Name, warehouse.Location, warehouse.ID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	err := repo.Update(ctx, warehouse)
	assert.NoError(t, err)
}

func TestWarehouseRepository_Delete(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewWarehouseRepository(db)
	ctx := context.Background()

	mock.ExpectExec("DELETE FROM warehouses WHERE id =").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	err := repo.Delete(ctx, 1)
	assert.NoError(t, err)
}

func TestWarehouseRepository_Delete_NotFound(t *testing.T) {
	db, mock, _ := sqlmock.New()
	defer db.Close()

	repo := NewWarehouseRepository(db)
	ctx := context.Background()

	mock.ExpectExec("DELETE FROM warehouses WHERE id =").
		WithArgs(999).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

	err := repo.Delete(ctx, 999)
	assert.Error(t, err)
	assert.Equal(t, "product not found", err.Error())
}
