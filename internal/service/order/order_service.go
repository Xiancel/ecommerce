package order

import (
	"context"

	repository "github.com/Xiancel/ecommerce/internal/repository/postgres"
	"github.com/google/uuid"
)

type service struct {
	orderRepo repository.OrderRepository
}

func NewService(orderRepo repository.OrderRepository) OrderService {
	return &service{orderRepo: orderRepo}
}

// CancelOrder implements OrderService.
func (s *service) CancelOrder(ctx context.Context, id uuid.UUID) error {
	panic("unimplemented")
}

// CreateOrder implements OrderService.
func (s *service) CreateOrder(ctx context.Context, req CreateOrderRequset) (*OrderResponse, error) {
	panic("unimplemented")
}

// GetOrder implements OrderService.
func (s *service) GetOrder(ctx context.Context, id uuid.UUID) (*OrderResponse, error) {
	panic("unimplemented")
}

// ListOrder implements OrderService.
func (s *service) ListOrder(ctx context.Context, filter OrderFilter) ([]*OrderResponse, error) {
	panic("unimplemented")
}

// UpdateOrderStatus implements OrderService.
func (s *service) UpdateOrderStatus(ctx context.Context, id uuid.UUID, status string) (*OrderResponse, error) {
	panic("unimplemented")
}
