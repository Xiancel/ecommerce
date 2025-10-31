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

// type CartItemResponse struct {
// 	ProductID    *uuid.UUID
// 	ProductName  string
// 	ProductPrice float64
// 	ProductStock int
// 	Quantity     int
// 	TotalPrice   float64
// }

type CartListResponse struct {
	Items      []*models.CartItem
	TotalPrice float64
}
