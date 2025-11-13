package cart

import (
	models "github.com/Xiancel/ecommerce/internal/domain"
	"github.com/google/uuid"
)

// DTO структури для кошика

type AddCartItemRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required,uuid"`
	Quantity  int       `json:"quantity" validate:"required,gt=0"`
}

type UpdateCartItemRequest struct {
	ProductID *uuid.UUID `json:"product_id" validate:"required,uuid"`
	Quantity  int        `json:"quantity" validate:"required,gt=0"`
}

type CartListResponse struct {
	Items      []*models.CartItem `json:"items"`
	TotalPrice float64            `json:"total_price"`
}
