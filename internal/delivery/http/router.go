package http

import (
	"github.com/gin-gonic/gin"
	"github.com/ifs21014-itdel/concurrent-order-processor/internal/usecase"
)

func NewRouter(
	authUC *usecase.AuthUsecase,
	productUC *usecase.ProductUsecase,
	wareHouseUC *usecase.WarehouseUsecase,
	warehouseStockUC *usecase.WarehouseStockUsecase,
	orderUC *usecase.OrderUsecase,
	cartUC *usecase.CartUsecase,
	cartItemUC *usecase.CartItemUsecase,
) *gin.Engine {
	r := gin.Default()

	// Group API routes
	api := r.Group("/api")

	// Initialize handlers
	NewAuthHandler(api, authUC)
	NewProductHandler(api, productUC)
	NewWarehouseHandler(api, wareHouseUC)
	NewWarehouseStockHandler(api, warehouseStockUC)
	NewOrderHandler(api, orderUC)
	NewCartHandler(api, cartUC)

	return r
}
