package cart

import "github.com/google/uuid"

type AddCartItemRequest struct {
	ProductID uuid.UUID
	Quantity  int
}

type UpdateCartItemRequest struct {
	ProductID *uuid.UUID
	Quantity  int
}

type CartItemResponse struct {
	ProductID    *uuid.UUID
	ProductName  string
	ProductPrice float64
	ProductStock int
	Quantity     int
	TotalPrice   float64
}

type CartListResponse struct {
	Items      []CartItemResponse
	TotalPrice float64
}
