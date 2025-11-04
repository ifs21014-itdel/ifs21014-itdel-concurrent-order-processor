package usecase

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/ifs21014-itdel/concurrent-order-processor/internal/domain"
	"github.com/ifs21014-itdel/concurrent-order-processor/internal/repository"
	"github.com/ifs21014-itdel/concurrent-order-processor/internal/utils"
)

// CartItemUsecase - Jika masih diperlukan untuk operasi spesifik cart item
type CartItemUsecase struct {
	repo repository.CartItemRepository
}

func NewCartItemUsecase(repo repository.CartItemRepository) *CartItemUsecase {
	return &CartItemUsecase{repo: repo}
}

func (u *CartItemUsecase) UpdateCartItem(ctx context.Context, item *domain.CartItem) error {
	if item.ID == 0 {
		return fmt.Errorf("invalid item ID")
	}
	if item.Quantity <= 0 {
		return fmt.Errorf("quantity must be greater than zero")
	}
	if item.SubTotal < 0 {
		return fmt.Errorf("subtotal cannot be negative")
	}
	return u.repo.UpdateCartItem(ctx, item)
}

func (u *CartItemUsecase) GetCartItemsByCartID(ctx context.Context, cartID uint) ([]domain.CartItem, error) {
	if cartID == 0 {
		return nil, fmt.Errorf("invalid cart ID")
	}
	return u.repo.GetCartItemsByCartID(ctx, cartID)
}

func (u *CartItemUsecase) ClearCart(ctx context.Context, cartID uint) error {
	if cartID == 0 {
		return fmt.Errorf("invalid cart ID")
	}
	return u.repo.ClearCart(ctx, cartID)
}

// ConcurrentUpdateItems - Batch update cart items secara concurrent
func (u *CartItemUsecase) ConcurrentUpdateItems(ctx context.Context, updates []domain.CartItem) error {
	if len(updates) == 0 {
		return fmt.Errorf("no items to update")
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(updates))

	for _, item := range updates {
		wg.Add(1)
		localItem := item

		go utils.SafeGoroutine(fmt.Sprintf("Do Update Item: %d", localItem.ID), func() {
			defer wg.Done()

			if err := u.UpdateCartItem(ctx, &localItem); err != nil {
				errChan <- fmt.Errorf("failed to update cart item ID %d: %w", localItem.ID, err)
				log.Printf("[ERROR] Failed to update cart item ID %d: %v", localItem.ID, err)
				return
			}

			log.Printf("[OK] Cart item ID %d has been successfully updated", localItem.ID)
		})
	}

	wg.Wait()
	close(errChan)

	// Collect errors if any
	var errors []error
	for err := range errChan {
		errors = append(errors, err)
	}

	if len(errors) > 0 {
		return fmt.Errorf("some cart items failed to update: %d errors occurred", len(errors))
	}

	log.Printf("[DONE] All %d cart items updated successfully", len(updates))
	return nil
}
