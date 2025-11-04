package usecase

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/ifs21014-itdel/concurrent-order-processor/internal/domain"
	"github.com/ifs21014-itdel/concurrent-order-processor/internal/repository"
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
}

func NewCartUsecase(cartRepo repository.CartRepository, cartItemRepo repository.CartItemRepository) *CartUsecase {
	return &CartUsecase{
		cartRepo:     cartRepo,
		cartItemRepo: cartItemRepo,
	}
}

// AddItemToCart - Logic utama untuk menambahkan item ke cart
func (u *CartUsecase) AddItemToCart(ctx context.Context, req AddItemToCartRequest) error {
	// Validasi input
	if err := u.validateAddItemRequest(req); err != nil {
		return err
	}

	// Get atau create cart
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

	log.Printf("[OK] Added product %d to cart ID %d", req.ProductID, cart.ID)
	return nil
}

// DeleteCart - Menghapus cart berdasarkan ID
func (u *CartUsecase) DeleteCart(ctx context.Context, id uint) error {
	if err := u.validateID(id, "cart"); err != nil {
		return err
	}
	return u.cartRepo.DeleteCart(ctx, id)
}

// DeleteCartItem - Menghapus cart item berdasarkan ID
func (u *CartUsecase) DeleteCartItem(ctx context.Context, id uint) error {
	if err := u.validateID(id, "cart item"); err != nil {
		return err
	}
	return u.cartItemRepo.DeleteCartItem(ctx, id)
}

// GetCartsWithItemsByUserID - Mengambil semua cart beserta items untuk user tertentu
func (u *CartUsecase) GetCartsWithItemsByUserID(ctx context.Context, userID uint) ([]CartWithItemsResponse, error) {
	if err := u.validateID(userID, "user"); err != nil {
		return nil, err
	}

	// Get all carts by user ID
	carts, err := u.cartRepo.GetCartByUserId(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get carts: %w", err)
	}

	if len(carts) == 0 {
		return []CartWithItemsResponse{}, nil
	}

	// Build response with cart items
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

	return result, nil
}

// getOrCreateCart - Internal method untuk get atau create cart
func (u *CartUsecase) getOrCreateCart(ctx context.Context, userID, warehouseID uint) (*domain.Cart, error) {
	cart, err := u.cartRepo.FindByUserAndWarehouse(ctx, userID, warehouseID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if cart != nil {
		return cart, nil
	}

	// Create new cart
	newCart := &domain.Cart{
		UserID:      userID,
		WarehouseID: warehouseID,
	}

	if err := u.cartRepo.CreateCart(ctx, newCart); err != nil {
		return nil, err
	}

	log.Printf("[OK] Created new cart for user %d at warehouse %d", userID, warehouseID)
	return newCart, nil
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
