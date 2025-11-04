package usecase

import (
	"context"
	"errors"

	"github.com/ifs21014-itdel/concurrent-order-processor/internal/domain"
	"github.com/ifs21014-itdel/concurrent-order-processor/internal/repository"
)

type WarehouseUsecase struct {
	repo repository.WarehouseRepository
}

func NewWarehouseUsecase(repo repository.WarehouseRepository) *WarehouseUsecase {
	return &WarehouseUsecase{repo: repo}
}

func (u *WarehouseUsecase) CreateWarehouse(ctx context.Context, warehouse *domain.Warehouse) error {
	if warehouse.Name == "" {
		return errors.New("warehouse name cannot be empty")
	}
	if warehouse.Location == "" {
		return errors.New("warehouse location cannot be empty")
	}
	return u.repo.Create(ctx, warehouse)
}

func (u *WarehouseUsecase) UpdateWarehouse(ctx context.Context, warehouse *domain.Warehouse) error {
	if warehouse.ID == 0 {
		return errors.New("warehouse ID is required")
	}
	return u.repo.Update(ctx, warehouse)
}

func (u *WarehouseUsecase) DeleteWarehouse(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid warehouse ID")
	}
	return u.repo.Delete(ctx, id)
}

func (u *WarehouseUsecase) GetAllWarehouses(ctx context.Context) ([]domain.Warehouse, error) {
	return u.repo.GetAll(ctx)
}

func (u *WarehouseUsecase) GetWareHouseById(ctx context.Context, id uint) (*domain.Warehouse, error) {
	return u.repo.GetById(ctx, id)
}
