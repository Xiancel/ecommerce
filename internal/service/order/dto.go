package order

import (
	models "github.com/Xiancel/ecommerce/internal/domain"
	"github.com/google/uuid"
)

// добавить `json`и валидейт
type CreateOrderRequset struct {
	UserID         uuid.UUID
	Items          []models.OrderItem
	ShippingAdress models.ShippingAddress
	PaymentMethod  string
}

type UpdateOrderRequest struct {
	Status string
}

type OrderListResponse struct {
	Order []*models.Order
	Total int
}

type OrderFilter struct {
	UserID  *uuid.UUID
	Status  string
	Limit   int
	Offset  int
	OrderBy string
}
