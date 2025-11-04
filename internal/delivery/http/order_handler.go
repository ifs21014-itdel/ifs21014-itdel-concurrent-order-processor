package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	uc "github.com/ifs21014-itdel/concurrent-order-processor/internal/usecase"
	"github.com/ifs21014-itdel/concurrent-order-processor/pkg/jwt"
)

type OrderHandler struct {
	usecase *uc.OrderUsecase
}

func NewOrderHandler(rg *gin.RouterGroup, orderUc *uc.OrderUsecase) {
	h := &OrderHandler{usecase: orderUc}
	protected := rg.Group("orders")
	protected.Use(jwt.AuthMiddleware())
	protected.POST("/", h.CreateOrder)
	protected.DELETE("/:id", h.DeleteOrder)
	protected.GET("/", h.GetOrderByUserId)
	protected.GET("/status/:status", h.GetOrderByUserIdAndStatus)
	protected.PATCH("/:id/status", h.UpdateOrderStatus)
}

type CreateOrderInput struct {
	CartID       uint    `json:"cart_id" binding:"required"`
	Status       string  `json:"status" binding:"required"`
	ShippingCost float64 `json:"shipping_cost"`
}

type UpdateStatusInput struct {
	Status string `json:"status" binding:"required"`
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var input CreateOrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid JSON input",
			"error":   err.Error(),
		})
		return
	}

	userID, _ := c.Get("userID")

	orderReq := uc.CreateOrderRequest{
		UserID:       userID.(uint),
		CartID:       input.CartID,
		Status:       input.Status,
		ShippingCost: input.ShippingCost,
	}

	ctx := c.Request.Context()
	if err := h.usecase.CreateOrder(ctx, orderReq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "failed to create order",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "order created successfully from cart",
	})
}

func (h *OrderHandler) DeleteOrder(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid order ID",
		})
		return
	}

	ctx := c.Request.Context()
	if err := h.usecase.DeleteOrder(ctx, uint(id)); err != nil {
		if err.Error() == "order not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "order deleted successfully",
	})
}

func (h *OrderHandler) GetOrderByUserId(c *gin.Context) {
	userID, _ := c.Get("userID")

	ctx := c.Request.Context()
	orders, err := h.usecase.GetOrderByUserId(ctx, userID.(uint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   orders,
	})
}

func (h *OrderHandler) GetOrderByUserIdAndStatus(c *gin.Context) {
	userID, _ := c.Get("userID")
	status := c.Param("status")

	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "status parameter is required",
		})
		return
	}

	ctx := c.Request.Context()
	orders, err := h.usecase.GetOrderByUserIdAndStatus(ctx, userID.(uint), status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   orders,
	})
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid order ID",
		})
		return
	}

	var input UpdateStatusInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid JSON input",
			"error":   err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	if err := h.usecase.UpdateOrderStatus(ctx, uint(id), input.Status); err != nil {
		if err.Error() == "order not found" {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "order status updated successfully",
	})
}
