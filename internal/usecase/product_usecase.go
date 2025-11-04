package usecase

import (
	"context"
	"errors"

	"github.com/ifs21014-itdel/concurrent-order-processor/internal/domain"
	"github.com/ifs21014-itdel/concurrent-order-processor/internal/repository"
)

type ProductUsecase struct {
	repo repository.ProductRepository
}

func NewProductUsecase(repo repository.ProductRepository) *ProductUsecase {
	return &ProductUsecase{repo: repo}
}

func (u *ProductUsecase) CreateProduct(ctx context.Context, product *domain.Product) error {
	existing, err := u.repo.FindByName(ctx, product.Name)
	if err == nil && existing != nil {
		return errors.New("product already exists")
	}

	return u.repo.Create(ctx, product)
}

func (u *ProductUsecase) UpdateProduct(ctx context.Context, product *domain.Product) error {
	if product.ID == 0 {
		return errors.New("product ID is required")
	}
	return u.repo.Update(ctx, product)
}

func (u *ProductUsecase) DeleteProduct(ctx context.Context, id int64) error {
	if id <= 0 {
		return errors.New("invalid product ID")
	}
	return u.repo.Delete(ctx, id)
}

func (u *ProductUsecase) GetProductByName(ctx context.Context, name string) (*domain.Product, error) {
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}
	return u.repo.FindByName(ctx, name)
}
func (u *ProductUsecase) GetProductById(ctx context.Context, id uint) (*domain.Product, error) {

	return u.repo.FindById(ctx, id)
}

func (u *ProductUsecase) GetAllProducts(ctx context.Context) ([]domain.Product, error) {
	return u.repo.GetAll(ctx)
}
