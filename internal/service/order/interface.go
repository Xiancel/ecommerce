package order

import (
	"context"

	models "github.com/Xiancel/ecommerce/internal/domain"
	"github.com/google/uuid"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req CreateOrderRequset) (*models.Order, error)
	GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error)
	ListOrder(ctx context.Context, filter OrderFilter) (*OrderListResponse, error)
	UpdateOrderStatus(ctx context.Context, id uuid.UUID, req UpdateOrderRequest) (*models.Order, error)
	CancelOrder(ctx context.Context, id uuid.UUID) error
}
