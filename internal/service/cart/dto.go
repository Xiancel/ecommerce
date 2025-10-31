package cart

import (
	models "github.com/Xiancel/ecommerce/internal/domain"
	"github.com/google/uuid"
)

type AddCartItemRequest struct {
	ProductID uuid.UUID
	Quantity  int
}

type UpdateCartItemRequest struct {
	ProductID *uuid.UUID
	Quantity  int
}

type CartListResponse struct {
	Items      []*models.CartItem
	TotalPrice float64
}
