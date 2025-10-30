package order

import (
	"context"

	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req CreateOrderRequset) (*OrderResponse, error)
	GetOrder(ctx context.Context, id uuid.UUID) (*OrderResponse, error)
	ListOrder(ctx context.Context, filter OrderFilter) ([]*OrderResponse, error)
	UpdateOrderStatus(ctx context.Context, id uuid.UUID, status string) (*OrderResponse, error)
	CancelOrder(ctx context.Context, id uuid.UUID) error
}
