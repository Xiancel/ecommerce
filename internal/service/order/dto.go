package order

import (
	models "github.com/Xiancel/ecommerce/internal/domain"
	"github.com/google/uuid"
)

type CreateOrderItemRequest struct {
	ProductID uuid.UUID `json:"product_id" validate:"required"`
	Quantity  int       `json:"quantity" validate:"required,min=1"`
}

type CreateOrderRequest struct {
	Items          []CreateOrderItemRequest `json:"items" validate:"required,dive"`
	ShippingAdress models.ShippingAddress   `json:"shipping_address" validate:"required"`
	PaymentMethod  string                   `json:"payment_method" validate:"required,oneof=card cash"`
}

type UpdateOrderRequest struct {
	Status string `json:"status" validate:"required,oneof=pending paid shipped canceled delivered"`
}

type OrderListResponse struct {
	Order []*models.Order `json:"order"`
	Total int             `json:"total"`
}

type OrderFilter struct {
	UserID  *uuid.UUID `json:"user_id" validate:"omitempty,uuid"`
	Status  string     `json:"status" validate:"omitempty,oneof=pending paid shipped canceled delivered"`
	Limit   int        `json:"limit" validate:"required,min=1,max=100"`
	Offset  int        `json:"offset" validate:"gte=0"`
	OrderBy string     `json:"order_by" validate:"omitempty,oneof=created_at_asc created_as_desc status_asc status_desc"`
}
