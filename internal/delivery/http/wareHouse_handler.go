package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ifs21014-itdel/concurrent-order-processor/internal/domain"
	uc "github.com/ifs21014-itdel/concurrent-order-processor/internal/usecase"
	"github.com/ifs21014-itdel/concurrent-order-processor/pkg/jwt"
)

type WarehouseHandler struct {
	usecase *uc.WarehouseUsecase
}

func NewWarehouseHandler(rg *gin.RouterGroup, uc *uc.WarehouseUsecase) {
	h := &WarehouseHandler{usecase: uc}

	protected := rg.Group("/warehouses")
	protected.Use(jwt.AuthMiddleware())
	protected.POST("/", h.Create)
	protected.PUT("/:id", h.Update)
	protected.GET("/", h.GetAll)
	protected.DELETE("/:id", h.Delete)
}

func (h *WarehouseHandler) Create(c *gin.Context) {
	var input domain.Warehouse
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx := c.Request.Context()
	if err := h.usecase.CreateWarehouse(ctx, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "warehouse created successfully",
		"data":    input,
	})
}

func (h *WarehouseHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid warehouse ID"})
		return
	}

	var input domain.Warehouse
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	input.ID = uint(id)
	ctx := c.Request.Context()
	if err := h.usecase.UpdateWarehouse(ctx, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "warehouse updated successfully",
		"data":    "input",
	})
}

func (h *WarehouseHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	warehouses, err := h.usecase.GetAllWarehouses(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   warehouses,
	})
}

func (h *WarehouseHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid warehouse ID"})
		return
	}

	ctx := c.Request.Context()
	if err := h.usecase.DeleteWarehouse(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "warehouse deleted successfully",
	})

}
