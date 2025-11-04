package usecase

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/ifs21014-itdel/concurrent-order-processor/internal/domain"
	"github.com/ifs21014-itdel/concurrent-order-processor/internal/repository"
	"github.com/ifs21014-itdel/concurrent-order-processor/internal/utils"
)

type WarehouseStockUsecase struct {
	repo          repository.WarehouseStockRepository
	warehouseRepo repository.WarehouseRepository
	productRepo   repository.ProductRepository
}

func NewWarehouseStockUsecase(
	repo repository.WarehouseStockRepository,
	warehouseRepo repository.WarehouseRepository,
	productRepo repository.ProductRepository,
) *WarehouseStockUsecase {
	return &WarehouseStockUsecase{repo, warehouseRepo, productRepo}
}

// Semua response dikembalikan dalam bentuk map agar handler tidak perlu mikir
func (u *WarehouseStockUsecase) Create(ctx context.Context, stock *domain.WarehouseStock) (map[string]interface{}, error) {
	if err := u.repo.Create(ctx, stock); err != nil {
		return nil, fmt.Errorf("failed to create warehouse stock: %w", err)
	}
	return map[string]interface{}{
		"status":  "success",
		"message": "warehouseStock created successfully",
		"data":    stock,
	}, nil
}

func (u *WarehouseStockUsecase) Delete(ctx context.Context, stockID uint) (map[string]interface{}, error) {
	if err := u.repo.Delete(ctx, stockID); err != nil {
		return nil, fmt.Errorf("failed to delete warehouse stock: %w", err)
	}
	return map[string]interface{}{
		"status":  "success",
		"message": "warehouseStock deleted successfully",
	}, nil
}

func (u *WarehouseStockUsecase) GetAll(ctx context.Context) (map[string]interface{}, error) {
	data, err := u.repo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get all warehouse stocks: %w", err)
	}
	return map[string]interface{}{
		"status": "success",
		"data":   data,
	}, nil
}

func (u *WarehouseStockUsecase) GetWarehouseDetail(ctx context.Context, warehouseID uint) (map[string]interface{}, error) {
	stocks, err := u.repo.GetByWarehouseID(ctx, warehouseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get stocks for warehouse ID %d: %w", warehouseID, err)
	}

	warehouse, err := u.warehouseRepo.GetById(ctx, warehouseID)
	if err != nil {
		return nil, fmt.Errorf("failed to get warehouse detail: %w", err)
	}

	var result []map[string]interface{}
	for _, s := range stocks {
		product, err := u.productRepo.FindById(ctx, s.ProductID)
		if err != nil {
			log.Printf("[WARN] Failed to get product ID %d: %v", s.ProductID, err)
			continue
		}
		result = append(result, map[string]interface{}{
			"id":     s.ID,
			"name":   product.Name,
			"price":  product.Price,
			"stock":  s.Quantity,
			"weight": product.Weight,
		})
	}

	return map[string]interface{}{
		"status":    "success",
		"warehouse": warehouse,
		"stocks":    result,
	}, nil
}

func (u *WarehouseStockUsecase) UpdateQuantity(ctx context.Context, warehouseID uint, productID uint, quantity int32) (map[string]interface{}, error) {
	if err := u.repo.UpdateQuantity(ctx, warehouseID, productID, quantity); err != nil {
		return nil, fmt.Errorf("failed to update stock quantity: %w", err)
	}
	return map[string]interface{}{
		"status":  "success",
		"message": fmt.Sprintf("Quantity for stock Product ID %d and Warehouse ID %d updated successfully", warehouseID, productID),
		"data": gin.H{
			"Product_id":   productID,
			"warehouse_id": warehouseID,
			"quantity":     quantity,
		},
	}, nil
}

func (u *WarehouseStockUsecase) ConcurrentUpdateQuantities(ctx context.Context, updates []domain.WarehouseStock) map[string]interface{} {
	var wg sync.WaitGroup
	for _, stock := range updates {
		wg.Add(1)
		go utils.SafeGoroutine(fmt.Sprintf("Update stock where warehoseID %d and Product Id %d", stock.WarehouseID, stock.ProductID), func() {
			defer wg.Done()
			if err := u.repo.UpdateQuantity(ctx, stock.WarehouseID, stock.ProductID, stock.Quantity); err != nil {
				log.Printf("[ERROR] Failed to update stock ID %d: %v", stock.ID, err)
			} else {
				log.Printf("[OK] Stock ID %d updated successfully", stock.ProductID)
			}
		})
	}
	wg.Wait()
	return map[string]interface{}{
		"status":  "done",
		"message": "All stock quantities updated successfully",
	}
}
