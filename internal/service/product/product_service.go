package product

import (
	"context"
	"fmt"

	models "github.com/Xiancel/ecommerce/internal/domain"
	repository "github.com/Xiancel/ecommerce/internal/repository/postgres"
	"github.com/google/uuid"
)

type service struct {
	productRepo repository.ProductRepository
}

func NewService(productRepo repository.ProductRepository) ProductService {
	return &service{productRepo: productRepo}
}

// CheckAvailability implements ProductService.
func (s *service) CheckAvailability(ctx context.Context, id uuid.UUID, quantity int) (bool, error) {
	panic("unimplemented")
}

// CreateProduct implements ProductService.
func (s *service) CreateProduct(ctx context.Context, req CreateProductRequest) (*models.Product, error) {
	panic("unimplemented")
}

// GetProduct implements ProductService.
func (s *service) GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	panic("unimplemented")
}

// ListProduct implements ProductService.
func (s *service) ListProduct(ctx context.Context, filter ProductFilter) (*ProductListResponse, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	if filter.Offset < 0 {
		filter.Limit = 0
	}
	if filter.MinPrice != nil && *filter.MinPrice < 0 {
		return nil, ErrInvalidPrice
	}

	if filter.MaxPrice != nil && *filter.MaxPrice < 0 {
		return nil, ErrInvalidPrice
	}

	repoFilter := models.ListFilter{
		CategoryID: filter.CategoryID,
		MinPrice:   filter.MinPrice,
		MaxPrice:   filter.MaxPrice,
		Search:     filter.Search,
		Limit:      filter.Limit,
		Offset:     filter.Offset,
		OrderBy:    filter.OrderBy,
	}

	products, err := s.productRepo.List(ctx, repoFilter)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	if filter.InStock != nil && *filter.InStock {
		filteredProducts := make([]*models.Product, 0)
		for _, p := range products {
			if p.Stock > 0 {
				filteredProducts = append(filteredProducts, p)
			}
		}
		products = filteredProducts
	}
	return &ProductListResponse{
		Products: products,
		Total:    len(products),
		Limit:    filter.Limit,
		Offset:   filter.Offset,
	}, nil
}

// ReleaseStock implements ProductService.
func (s *service) ReleaseStock(ctx context.Context, id uuid.UUID, quantity int) error {
	panic("unimplemented")
}

// ReserveStock implements ProductService.
func (s *service) ReserveStock(ctx context.Context, id uuid.UUID, quantity int) error {
	panic("unimplemented")
}

// SearchProduct implements ProductService.
func (s *service) SearchProduct(ctx context.Context, query string, limit int, offset int) ([]*models.Product, error) {
	panic("unimplemented")
}

// UpdateProduct implements ProductService.
func (s *service) UpdateProduct(ctx context.Context, id uuid.UUID, query string, req UpdateProductRequest) ([]*models.Product, error) {
	panic("unimplemented")
}
