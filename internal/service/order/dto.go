package order

import (
	"github.com/google/uuid"
)

// добавить `json`и валидейт
type CreateOrderRequset struct {
	UserID         uuid.UUID
	Items          []OrderItemInput
	ShippingAdress ShippingAdress
	PaymentMethod  string
}

type OrderItemInput struct {
	ProductID uuid.UUID
	Quantity  int
}

type ShippingAdress struct {
	Street     string
	City       string
	PostalCode string
	Country    string
}

type OrderResponse struct {
	ID             uuid.UUID
	UserID         *uuid.UUID
	Status         string
	TotalAmount    float64
	ShippingAdress ShippingAdress
	PaymentMethod  string
}

type OrderItemResponse struct {
	ProductID uuid.UUID
	Quantity  int
	Price     float64
}

type OrderFilter struct {
	UserID  *uuid.UUID
	Status  string
	Limit   int
	Offset  int
	OrderBy string
}
