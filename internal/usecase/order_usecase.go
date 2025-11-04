package usecase

import (
	"context"
	"fmt"

	"github.com/ifs21014-itdel/concurrent-order-processor/internal/domain"
	"github.com/ifs21014-itdel/concurrent-order-processor/internal/repository"
)

type CreateOrderRequest struct {
	UserID       uint
	CartID       uint
	Status       string
	ShippingCost float64
}

type OrderWithItemsResponse struct {
	Order domain.Order       `json:"order"`
	Items []domain.OrderItem `json:"items"`
}

type OrderUsecase struct {
	orderRepo          repository.OrderRepository
	orderItemRepo      repository.OrderItemRepository
	cartRepo           repository.CartRepository
	cartItemRepo       repository.CartItemRepository
	warehouseStockRepo repository.WarehouseStockRepository
}

func NewOrderUsecase(
	orderRepo repository.OrderRepository,
	orderItemRepo repository.OrderItemRepository,
	cartRepo repository.CartRepository,
	cartItemRepo repository.CartItemRepository,
	warehouseStockRepo repository.WarehouseStockRepository,
) *OrderUsecase {
	return &OrderUsecase{
		orderRepo:          orderRepo,
		orderItemRepo:      orderItemRepo,
		cartRepo:           cartRepo,
		cartItemRepo:       cartItemRepo,
		warehouseStockRepo: warehouseStockRepo,
	}
}

func (o *OrderUsecase) CreateOrder(ctx context.Context, req CreateOrderRequest) error {
	if err := o.validateCreateOrderRequest(req); err != nil {
		return err
	}

	cartItems, err := o.cartItemRepo.GetCartItemsByCartID(ctx, req.CartID)
	if err != nil {
		return fmt.Errorf("failed to get cart items: %w", err)
	}
	if len(cartItems) == 0 {
		return fmt.Errorf("cart is empty")
	}

	tx, err := o.orderRepo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	var totalPrice float64
	for _, item := range cartItems {
		totalPrice += item.SubTotal
	}

	order := &domain.Order{
		UserID:       req.UserID,
		Status:       req.Status,
		TotalPrice:   totalPrice,
		ShippingCost: req.ShippingCost,
	}

	if err := o.orderRepo.CreateOrderTx(ctx, tx, order); err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	for _, cartItem := range cartItems {
		warehouseID := uint(1)
		if err := o.orderItemRepo.CreateOrderItemTx(ctx, tx, &domain.OrderItem{
			OrderID:   order.ID,
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
			SubTotal:  cartItem.SubTotal,
		}); err != nil {
			return fmt.Errorf("failed to create order item: %w", err)
		}

		if err := o.warehouseStockRepo.SafeDecreaseQuantity(ctx, tx, warehouseID, cartItem.ProductID, cartItem.Quantity); err != nil {
			return fmt.Errorf("failed to decrease stock: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (o *OrderUsecase) DeleteOrder(ctx context.Context, id uint) error {
	if err := o.orderRepo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func (o *OrderUsecase) GetOrderByUserId(ctx context.Context, userID uint) ([]domain.Order, error) {
	orders, err := o.orderRepo.GetOrderByUserId(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders: %w", err)
	}
	return orders, nil
}

func (o *OrderUsecase) GetOrderByUserIdAndStatus(ctx context.Context, userID uint, status string) ([]domain.Order, error) {
	if err := o.validateStatus(status); err != nil {
		return nil, err
	}

	orders, err := o.orderRepo.GetOrderByUserIdAndStatus(ctx, userID, status)
	if err != nil {
		return nil, fmt.Errorf("failed to get orders by status: %w", err)
	}
	return orders, nil
}

func (o *OrderUsecase) UpdateOrderStatus(ctx context.Context, id uint, status string) error {
	if err := o.validateStatus(status); err != nil {
		return err
	}

	if err := o.orderRepo.UpdateOrderStatus(ctx, id, status); err != nil {
		return err
	}
	return nil
}

func (o *OrderUsecase) validateCreateOrderRequest(req CreateOrderRequest) error {
	if req.UserID == 0 {
		return fmt.Errorf("user ID is required")
	}
	if req.CartID == 0 {
		return fmt.Errorf("cart ID is required")
	}
	if req.Status == "" {
		return fmt.Errorf("order status is required")
	}
	if err := o.validateStatus(req.Status); err != nil {
		return err
	}
	if req.ShippingCost < 0 {
		return fmt.Errorf("shipping cost cannot be negative")
	}
	return nil
}

func (o *OrderUsecase) validateStatus(status string) error {
	validStatuses := map[string]bool{
		"pending":   true,
		"processed": true,
		"shipped":   true,
		"delivered": true,
		"cancelled": true,
	}
	if !validStatuses[status] {
		return fmt.Errorf("invalid status: must be 'pending', 'processed', 'shipped', 'delivered', or 'cancelled'")
	}
	return nil
}
