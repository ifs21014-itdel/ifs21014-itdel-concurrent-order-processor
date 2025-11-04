package main

import (
	"log"
	"os"

	"github.com/ifs21014-itdel/concurrent-order-processor/config"
	"github.com/ifs21014-itdel/concurrent-order-processor/internal/delivery/http"
	repo "github.com/ifs21014-itdel/concurrent-order-processor/internal/repository"
	usecase "github.com/ifs21014-itdel/concurrent-order-processor/internal/usecase"
	"github.com/joho/godotenv"
)

func main() {
	// load env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Validasi JWT_SECRET
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET must be set in environment variables")
	}
	log.Println("JWT_SECRET loaded successfully (length:", len(jwtSecret), ")")

	db, err := config.NewDB()
	if err != nil {
		log.Fatal("db:", err)
	}

	// repo -> usecase -> handler
	userRepo := repo.NewUserRepository(db)
	productRepo := repo.NewProductRepository(db)
	wareHouseRepo := repo.NewWarehouseRepository(db)
	wareHouseStockRepo := repo.NewWarehouseStockRepository(db)
	orderRepo := repo.NewOrderRepository(db)
	orderItemRepo := repo.NewOrderItemRepository(db)
	cartRepo := repo.NewCartRepository(db)
	cartItemRepo := repo.NewCartItemRepository(db)

	authUC := usecase.NewAuthUsecase(userRepo)
	productUC := usecase.NewProductUsecase(productRepo)
	wareHouseUC := usecase.NewWarehouseUsecase(wareHouseRepo)
	wareHouseStockUC := usecase.NewWarehouseStockUsecase(wareHouseStockRepo, wareHouseRepo, productRepo)
	orderUC := usecase.NewOrderUsecase(orderRepo, orderItemRepo, cartRepo, cartItemRepo, wareHouseStockRepo)
	cartUC := usecase.NewCartUsecase(cartRepo, cartItemRepo)
	cartItemUC := usecase.NewCartItemUsecase(cartItemRepo)

	r := http.NewRouter(authUC, productUC, wareHouseUC, wareHouseStockUC, orderUC, cartUC, cartItemUC)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("listen on :", port)
	r.Run(":" + port)
}
