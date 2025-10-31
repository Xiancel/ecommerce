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
	product, err := s.productRepo.GetById(ctx, id)
	if err != nil {
		return false, nil
	}
	return product.Stock >= quantity, nil
}

// CreateProduct implements ProductService.
func (s *service) CreateProduct(ctx context.Context, req CreateProductRequest) (*models.Product, error) {
	if req.Name == "" {
		return nil, ErrProductNameRequired
	}

	if req.Price <= 0 {
		return nil, ErrInvalidPrice
	}

	if req.Stock < 0 {
		return nil, ErrInvalidStock
	}

	var description *string
	if req.Description != "" {
		description = &req.Description
	}

	var imageURL *string
	if req.ImageURL != "" {
		imageURL = &req.ImageURL
	}

	product := &models.Product{
		Name:        req.Name,
		Description: description,
		Price:       req.Price,
		Stock:       req.Stock,
		CategoryID:  req.CategoryID,
		ImageURL:    imageURL,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	return product, nil
}

// GetProduct implements ProductService.
func (s *service) GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	product, err := s.productRepo.GetById(ctx, id)
	if err != nil {
		return nil, ErrProductNotFound
	}

	return product, nil
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
		filter.Offset = 0
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
	product, err := s.productRepo.GetById(ctx, id)
	if err != nil {
		return err
	}

	product.Stock += quantity
	return s.productRepo.Update(ctx, product)
}

// ReserveStock implements ProductService.
func (s *service) ReserveStock(ctx context.Context, id uuid.UUID, quantity int) error {
	product, err := s.productRepo.GetById(ctx, id)
	if err != nil {
		return err
	}

	if product.Stock < quantity {
		return ErrInvalidStock
	}

	product.Stock -= quantity

	return s.productRepo.Update(ctx, product)
}

// SearchProduct implements ProductService.
func (s *service) SearchProduct(ctx context.Context, query string, limit int, offset int) ([]*models.Product, error) {
	if query == "" {
		return []*models.Product{}, nil
	}

	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		limit = 0
	}

	filter := models.ListFilter{
		Search: query,
		Limit:  limit,
		Offset: offset,
	}

	products, err := s.productRepo.List(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to search product: %w", err)
	}

	return products, nil
}

// UpdateProduct implements ProductService.
func (s *service) UpdateProduct(ctx context.Context, id uuid.UUID, query string, req UpdateProductRequest) (*models.Product, error) {
	product, err := s.productRepo.GetById(ctx, id)
	if err != nil {
		return nil, ErrProductNotFound
	}

	if req.Name != nil {
		if *req.Name == "" {
			return nil, ErrProductNameRequired
		}
		product.Name = *req.Name
	}
	if req.Description != nil {
		product.Description = req.Description
	}
	if req.Price != nil {
		if *req.Price <= 0 {
			return nil, ErrInvalidPrice
		}
		product.Price = *req.Price
	}
	if req.Stock != nil {
		if *req.Stock <= 0 {
			return nil, ErrInvalidStock
		}
		product.Stock = *req.Stock
	}

	if req.CategoryID != nil {
		product.CategoryID = req.CategoryID
	}

	if req.ImageURL != nil {
		product.ImageURL = req.ImageURL
	}

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	return product, nil
}
