package http

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ifs21014-itdel/concurrent-order-processor/internal/domain"
	uc "github.com/ifs21014-itdel/concurrent-order-processor/internal/usecase"
	"github.com/ifs21014-itdel/concurrent-order-processor/pkg/jwt"
)

type ProductHandler struct {
	usecase *uc.ProductUsecase
}

func NewProductHandler(rg *gin.RouterGroup, uc *uc.ProductUsecase) {
	h := &ProductHandler{usecase: uc}

	protected := rg.Group("/products")
	protected.Use(jwt.AuthMiddleware())

	protected.POST("/", h.Create)
	protected.GET("/", h.GetAll)
	protected.GET("/:name", h.GetByName)
	protected.PUT("/:id", h.Update)
	protected.DELETE("/:id", h.Delete)
}

func (h *ProductHandler) Create(c *gin.Context) {
	var input domain.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	userID, _ := c.Get("userID")
	fmt.Printf("All context keys: %+v\n", c.Keys)

	fmt.Println("user Id:", userID)
	input.UserID = userID.(uint)
	if err := h.usecase.CreateProduct(ctx, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "product created successfully",
		"data":    input,
	})
}

func (h *ProductHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	products, err := h.usecase.GetAllProducts(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   products,
	})
}

func (h *ProductHandler) GetByName(c *gin.Context) {
	name := c.Param("name")
	ctx := c.Request.Context()

	product, err := h.usecase.GetProductByName(ctx, name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   product,
	})
}

func (h *ProductHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	var input domain.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input.ID = uint(id)

	ctx := c.Request.Context()
	if err := h.usecase.UpdateProduct(ctx, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "product updated successfully",
		"data":    input,
	})
}

func (h *ProductHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	ctx := c.Request.Context()
	if err := h.usecase.DeleteProduct(ctx, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "product deleted successfully",
	})
}
