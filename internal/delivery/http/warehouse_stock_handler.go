package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ifs21014-itdel/concurrent-order-processor/internal/domain"
	uc "github.com/ifs21014-itdel/concurrent-order-processor/internal/usecase"
	"github.com/ifs21014-itdel/concurrent-order-processor/pkg/jwt"
)

type WarehouseStockHandler struct {
	usecase *uc.WarehouseStockUsecase
}

func NewWarehouseStockHandler(
	rg *gin.RouterGroup,
	warehouseStockUc *uc.WarehouseStockUsecase,
) {
	h := &WarehouseStockHandler{usecase: warehouseStockUc}

	protected := rg.Group("/warehouseStocks")
	protected.Use(jwt.AuthMiddleware())
	{
		protected.POST("/", h.Create)
		protected.DELETE("/:id", h.Delete)
		protected.GET("/", h.GetAll)
		protected.PUT("/concurrent", h.ConcurrentUpdateQuantities)
		protected.GET("/:warehouseId", h.GetByWareHouseId)
		protected.PUT("/:id", h.UpdateQuantity)
	}
}

func (h *WarehouseStockHandler) Create(c *gin.Context) {
	var input domain.WarehouseStock
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	res, err := h.usecase.Create(ctx, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, res)
}

func (h *WarehouseStockHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid warehouseStock ID"})
		return
	}

	ctx := c.Request.Context()
	res, err := h.usecase.Delete(ctx, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *WarehouseStockHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	res, err := h.usecase.GetAll(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *WarehouseStockHandler) GetByWareHouseId(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("warehouseId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid warehouse ID"})
		return
	}

	ctx := c.Request.Context()
	res, err := h.usecase.GetWarehouseDetail(ctx, uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *WarehouseStockHandler) UpdateQuantity(c *gin.Context) {

	var input domain.WarehouseStock
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	res, err := h.usecase.UpdateQuantity(ctx, input.WarehouseID, input.ProductID, input.Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *WarehouseStockHandler) ConcurrentUpdateQuantities(c *gin.Context) {
	var updates []domain.WarehouseStock
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	res := h.usecase.ConcurrentUpdateQuantities(ctx, updates)
	c.JSON(http.StatusOK, res)
}
