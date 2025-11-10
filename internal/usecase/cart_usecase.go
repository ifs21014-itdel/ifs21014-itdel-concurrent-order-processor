package usecase

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ifs21014-itdel/concurrent-order-processor/internal/domain"
	"github.com/ifs21014-itdel/concurrent-order-processor/internal/repository"
	"github.com/redis/go-redis/v9"
)

type AddItemToCartRequest struct {
	UserID      uint    `json:"-"`
	WarehouseID uint    `json:"warehouse_id"`
	ProductID   uint    `json:"product_id"`
	Quantity    int     `json:"quantity"`
	SubTotal    float64 `json:"sub_total"`
}

type CartWithItemsResponse struct {
	Cart  domain.Cart       `json:"cart"`
	Items []domain.CartItem `json:"items"`
}

type CartUsecase struct {
	cartRepo     repository.CartRepository
	cartItemRepo repository.CartItemRepository
	cache        *redis.Client
}

func NewCartUsecase(cartRepo repository.CartRepository, cartItemRepo repository.CartItemRepository, cache *redis.Client) *CartUsecase {
	return &CartUsecase{
		cartRepo:     cartRepo,
		cartItemRepo: cartItemRepo,
		cache:        cache,
	}
}

// Cache keys and TTL
const (
	cartsByUserKeyPrefix = "carts:user:" // carts:user:<userID>
	cartTTL              = 5 * time.Minute
)

// AddItemToCart - Add an item to a user's cart
func (u *CartUsecase) AddItemToCart(ctx context.Context, req AddItemToCartRequest) error {
	// Validate input
	if err := u.validateAddItemRequest(req); err != nil {
		return err
	}

	// Get or create cart
	cart, err := u.getOrCreateCart(ctx, req.UserID, req.WarehouseID)
	if err != nil {
		return fmt.Errorf("failed to get or create cart: %w", err)
	}

	// Create cart item
	item := &domain.CartItem{
		CartID:    cart.ID,
		ProductID: req.ProductID,
		Quantity:  int32(req.Quantity),
		SubTotal:  req.SubTotal,
	}

	if err := u.cartItemRepo.AddCartItem(ctx, item); err != nil {
		return fmt.Errorf("failed to add cart item: %w", err)
	}

	// Invalidate cache after adding item
	u.invalidateUserCartCache(ctx, req.UserID)

	log.Printf("[OK] Added product %d to cart ID %d", req.ProductID, cart.ID)
	return nil
}

// DeleteCart - Delete cart by ID
func (u *CartUsecase) DeleteCart(ctx context.Context, id uint, userID uint) error {
	if err := u.validateID(id, "cart"); err != nil {
		return err
	}

	if err := u.cartRepo.DeleteCart(ctx, id); err != nil {
		return err
	}

	// Invalidate cache
	u.invalidateUserCartCache(ctx, userID)
	return nil
}

// DeleteCartItem - Delete cart item by ID
func (u *CartUsecase) DeleteCartItem(ctx context.Context, id uint, userID uint) error {
	if err := u.validateID(id, "cart item"); err != nil {
		return err
	}

	if err := u.cartItemRepo.DeleteCartItem(ctx, id); err != nil {
		return err
	}

	// Invalidate cache
	u.invalidateUserCartCache(ctx, userID)
	return nil
}

// GetCartsWithItemsByUserID - Get all carts with items for a user (cached)
func (u *CartUsecase) GetCartsWithItemsByUserID(ctx context.Context, userID uint) ([]CartWithItemsResponse, error) {
	if err := u.validateID(userID, "user"); err != nil {
		return nil, err
	}

	cacheKey := fmt.Sprintf("%s%d", cartsByUserKeyPrefix, userID)
	if u.cache != nil {
		cached, err := u.cache.Get(ctx, cacheKey).Result()
		if err == nil {
			var carts []CartWithItemsResponse
			if json.Unmarshal([]byte(cached), &carts) == nil {
				log.Println("✅ [CACHE HIT] Carts retrieved from Redis")
				return carts, nil
			}
		} else {
			log.Println("ℹ️ [CACHE MISS] Cart cache not found for user", userID)
		}
	}

	// Fetch from DB
	carts, err := u.cartRepo.GetCartByUserId(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get carts: %w", err)
	}

	if len(carts) == 0 {
		return []CartWithItemsResponse{}, nil
	}

	result := make([]CartWithItemsResponse, 0, len(carts))
	for _, cart := range carts {
		items, err := u.cartItemRepo.GetCartItemsByCartID(ctx, cart.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get cart items for cart ID %d: %w", cart.ID, err)
		}
		result = append(result, CartWithItemsResponse{
			Cart:  cart,
			Items: items,
		})
	}

	// Cache the result
	if u.cache != nil {
		if data, err := json.Marshal(result); err == nil {
			if err := u.cache.Set(ctx, cacheKey, data, cartTTL).Err(); err != nil {
				log.Printf("⚠️ [CACHE SET ERROR] Failed to save cart data to Redis: %v\n", err)
			} else {
				log.Println("✅ [CACHE SET] Cart data cached in Redis")
			}
		}
	}

	return result, nil
}

// getOrCreateCart - Get existing cart or create a new one
func (u *CartUsecase) getOrCreateCart(ctx context.Context, userID, warehouseID uint) (*domain.Cart, error) {
	cart, err := u.cartRepo.FindByUserAndWarehouse(ctx, userID, warehouseID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if cart != nil {
		return cart, nil
	}

	newCart := &domain.Cart{
		UserID:      userID,
		WarehouseID: warehouseID,
	}

	if err := u.cartRepo.CreateCart(ctx, newCart); err != nil {
		return nil, err
	}

	// Invalidate cache in case user had empty cart before
	u.invalidateUserCartCache(ctx, userID)

	log.Printf("[OK] Created new cart for user %d at warehouse %d", userID, warehouseID)
	return newCart, nil
}

// Helper to invalidate user cart cache
func (u *CartUsecase) invalidateUserCartCache(ctx context.Context, userID uint) {
	if u.cache == nil {
		return
	}
	cacheKey := fmt.Sprintf("%s%d", cartsByUserKeyPrefix, userID)
	u.cache.Del(ctx, cacheKey)
}

// Validation methods
func (u *CartUsecase) validateAddItemRequest(req AddItemToCartRequest) error {
	if req.UserID == 0 {
		return fmt.Errorf("user ID is required")
	}
	if req.WarehouseID == 0 {
		return fmt.Errorf("warehouse ID is required")
	}
	if req.ProductID == 0 {
		return fmt.Errorf("product ID is required")
	}
	if req.Quantity <= 0 {
		return fmt.Errorf("quantity must be greater than zero")
	}
	if req.SubTotal < 0 {
		return fmt.Errorf("subtotal cannot be negative")
	}
	return nil
}

func (u *CartUsecase) validateID(id uint, entityName string) error {
	if id == 0 {
		return fmt.Errorf("invalid %s ID", entityName)
	}
	return nil
}
