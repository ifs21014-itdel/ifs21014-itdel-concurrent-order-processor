package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	uc "github.com/ifs21014-itdel/concurrent-order-processor/internal/usecase"
	"github.com/ifs21014-itdel/concurrent-order-processor/pkg/jwt"
)

type CartHandler struct {
	cartUsecase *uc.CartUsecase
}

func NewCartHandler(rg *gin.RouterGroup, cartUC *uc.CartUsecase) {
	h := &CartHandler{
		cartUsecase: cartUC,
	}

	protected := rg.Group("/cart")
	protected.Use(jwt.AuthMiddleware())
	protected.POST("/item", h.AddItemToCart)
	protected.GET("/", h.GetCartsByUser)
	protected.DELETE("/:id", h.DeleteCart)
	protected.DELETE("/item/:id", h.DeleteCartItem)
}

func (h *CartHandler) AddItemToCart(c *gin.Context) {
	var req uc.AddItemToCartRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON input"})
		return
	}

	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	req.UserID = userID
	ctx := c.Request.Context()

	if err := h.cartUsecase.AddItemToCart(ctx, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "item berhasil ditambahkan ke cart",
	})
}

func (h *CartHandler) DeleteCart(c *gin.Context) {
	id, err := h.getIDFromParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cart ID"})
		return
	}
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	if err := h.cartUsecase.DeleteCart(ctx, id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "cart berhasil dihapus",
	})
}

func (h *CartHandler) DeleteCartItem(c *gin.Context) {

	id, err := h.getIDFromParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cart item ID"})
		return
	}
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	if err := h.cartUsecase.DeleteCartItem(ctx, id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "cart item berhasil dihapus",
	})
}

func (h *CartHandler) GetCartsByUser(c *gin.Context) {
	userID, err := h.getUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	result, err := h.cartUsecase.GetCartsWithItemsByUserID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(result) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no carts found for this user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}

// Helper methods
func (h *CartHandler) getUserIDFromContext(c *gin.Context) (uint, error) {
	userIDValue, exists := c.Get("userID")
	if !exists {
		return 0, fmt.Errorf("userID not found in context")
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		return 0, fmt.Errorf("invalid userID type")
	}

	return userID, nil
}

func (h *CartHandler) getIDFromParam(c *gin.Context, param string) (uint, error) {
	idParam := c.Param(param)
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}
