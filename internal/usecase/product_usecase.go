package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/ifs21014-itdel/concurrent-order-processor/internal/domain"
	"github.com/ifs21014-itdel/concurrent-order-processor/internal/repository"
	"github.com/redis/go-redis/v9"
)

type ProductUsecase struct {
	repo  repository.ProductRepository
	cache *redis.Client
}

func NewProductUsecase(repo repository.ProductRepository, cache *redis.Client) *ProductUsecase {
	return &ProductUsecase{
		repo:  repo,
		cache: cache,
	}
}

// Cache keys
const (
	productByIDPrefix   = "product:id:"
	productByNamePrefix = "product:name:"
	allProductsKey      = "products:all"
	cacheTTL            = 5 * time.Minute
)

func (u *ProductUsecase) CreateProduct(ctx context.Context, product *domain.Product) error {
	existing, err := u.repo.FindByName(ctx, product.Name)
	if err == nil && existing != nil {
		return errors.New("product already exists")
	}

	err = u.repo.Create(ctx, product)
	if err != nil {
		return err
	}

	// Invalidate cache after create
	u.invalidateCache(ctx, product)
	return nil
}

func (u *ProductUsecase) UpdateProduct(ctx context.Context, product *domain.Product) error {
	if product.ID == 0 {
		return errors.New("product ID is required")
	}

	err := u.repo.Update(ctx, product)
	if err != nil {
		return err
	}

	// Invalidate cache after update
	u.invalidateCache(ctx, product)
	return nil
}

func (u *ProductUsecase) DeleteProduct(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid product ID")
	}

	// Get product first to invalidate cache properly
	product, _ := u.repo.FindById(ctx, uint(id))

	err := u.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	// Invalidate cache after delete
	if product != nil {
		u.invalidateCache(ctx, product)
	}
	return nil
}

func (u *ProductUsecase) GetProductByName(ctx context.Context, name string) (*domain.Product, error) {
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}

	// Check cache first
	if u.cache != nil {
		cacheKey := productByNamePrefix + name
		cachedData, err := u.cache.Get(ctx, cacheKey).Result()
		if err == nil {
			var product domain.Product
			if json.Unmarshal([]byte(cachedData), &product) == nil {
				return &product, nil
			}
		}
	}

	// Cache miss - get from DB
	product, err := u.repo.FindByName(ctx, name)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if u.cache != nil && product != nil {
		log.Println("✅ Data diambil dari Redis cache")
		u.cacheProduct(ctx, product)
	}

	return product, nil
}

func (u *ProductUsecase) GetProductById(ctx context.Context, id uint) (*domain.Product, error) {
	// Check cache first
	if u.cache != nil {
		cacheKey := fmt.Sprintf("%s%d", productByIDPrefix, id)
		cachedData, err := u.cache.Get(ctx, cacheKey).Result()
		if err == nil {
			var product domain.Product
			if json.Unmarshal([]byte(cachedData), &product) == nil {
				return &product, nil
			}
		}
	}

	// Cache miss - get from DB
	product, err := u.repo.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	// Store in cache
	if u.cache != nil && product != nil {
		u.cacheProduct(ctx, product)
	}

	return product, nil
}

func (u *ProductUsecase) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
	// Try to get products from Redis cache first
	if u.cache != nil {
		cachedData, err := u.cache.Get(ctx, allProductsKey).Result()
		if err == nil {
			var products []domain.Product
			if json.Unmarshal([]byte(cachedData), &products) == nil {
				log.Println("✅ [CACHE HIT] Products retrieved from Redis")
				return products, nil
			}
		} else {
			log.Println("ℹ️ [CACHE MISS] Key not found in Redis")
		}
	}

	// Cache miss — fetch products from the database
	log.Println("💾 [DB QUERY] Fetching products from the database")
	products, err := u.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	if u.cache != nil {
		if data, err := json.Marshal(products); err == nil {
			err := u.cache.Set(ctx, allProductsKey, data, cacheTTL).Err()
			if err != nil {
				log.Printf("⚠️ [CACHE SET ERROR] Failed to save data to Redis: %v\n", err)
			} else {
				log.Println("✅ [CACHE SET] Products cached successfully in Redis")
			}
		}
	}

	return products, nil
}

func (u *ProductUsecase) cacheProduct(ctx context.Context, product *domain.Product) {
	if u.cache == nil {
		return
	}

	data, err := json.Marshal(product)
	if err != nil {
		return
	}

	// Cache by ID
	idKey := fmt.Sprintf("%s%d", productByIDPrefix, product.ID)
	u.cache.Set(ctx, idKey, data, cacheTTL)

	// Cache by Name
	nameKey := productByNamePrefix + product.Name
	u.cache.Set(ctx, nameKey, data, cacheTTL)
}

// Helper function to invalidate cache
func (u *ProductUsecase) invalidateCache(ctx context.Context, product *domain.Product) {
	if u.cache == nil {
		return
	}

	// Delete specific product caches
	idKey := fmt.Sprintf("%s%d", productByIDPrefix, product.ID)
	nameKey := productByNamePrefix + product.Name

	u.cache.Del(ctx, idKey, nameKey, allProductsKey)
}
