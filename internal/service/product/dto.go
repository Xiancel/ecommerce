package product

import (
	models "github.com/Xiancel/ecommerce/internal/domain"
	"github.com/google/uuid"
)

// CreateProductRequest is the DTO for creating a product
type CreateProductRequest struct {
	Name        string     `json:"name" validate:"required,min=3,max=255"`
	Description string     `json:"description" validate:"max=1000"`
	Price       float64    `json:"price" validate:"required,gt=0"`
	Stock       int        `json:"stock" validate:"required,gte=0"`
	CategoryID  *uuid.UUID `json:"category_id" validate:"omitempty,uuid"`
	ImageURL    string     `json:"image_url" validate:"omitempty,url"`
}

// UpdateProductRequest is the DTO for updating a product
type UpdateProductRequest struct {
	Name        *string    `json:"name" validate:"omitempty,min=3,max=255"`
	Description *string    `json:"description" validate:"omitempty,max=1000"`
	Price       *float64   `json:"price" validate:"omitempty,gt=0"`
	Stock       *int       `json:"stock" validate:"omitempty,gte=0"`
	CategoryID  *uuid.UUID `json:"category_id" validate:"omitempty,uuid"`
	ImageURL    *string    `json:"image_url" validate:"omitempty,url"`
}

// ProductFilter is the DTO for filtering products
type ProductFilter struct {
	CategoryID *uuid.UUID `json:"category_id"`
	MinPrice   *float64   `json:"min_price" validate:"omitempty,gte=0"`
	MaxPrice   *float64   `json:"max_price" validate:"omitempty,gte=0"`
	Search     string     `json:"search"`
	InStock    *bool      `json:"in_stock"`
	OrderBy    string     `json:"order_by" validate:"omitempty,oneof=price_asc price_desc name_asc name_desc created_at_asc created_at_desc"`
	Limit      int        `json:"limit" validate:"required,min=1,max=100"`
	Offset     int        `json:"offset" validate:"gte=0"`
}

// ProductListResponse contains paginated products and metadata
type ProductListResponse struct {
	Products []*models.Product `json:"products"`
	Total    int               `json:"total"`
	Limit    int               `json:"limit"`
	Offset   int               `json:"offset"`
}
