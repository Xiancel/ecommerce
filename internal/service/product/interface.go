package product

import (
	"context"

	models "github.com/Xiancel/ecommerce/internal/domain"
	"github.com/google/uuid"
)

type ProductService interface {
	CreateProduct(ctx context.Context, req CreateProductRequest) (*models.Product, error)
	GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error)
	ListProduct(ctx context.Context, filter ProductFilter) (*ProductListResponse, error)
	SearchProduct(ctx context.Context, query string, limit, offset int) ([]*models.Product, error)
	UpdateProduct(ctx context.Context, id uuid.UUID, query string, req UpdateProductRequest) ([]*models.Product, error)
	CheckAvailability(ctx context.Context, id uuid.UUID, quantity int) (bool, error)
	ReserveStock(ctx context.Context, id uuid.UUID, quantity int) error
	ReleaseStock(ctx context.Context, id uuid.UUID, quantity int) error
}
