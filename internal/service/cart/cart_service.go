package cart

import (
	"context"

	repository "github.com/Xiancel/ecommerce/internal/repository/postgres"
	"github.com/google/uuid"
)

type service struct {
	CartRepo repository.CartRepository
}

func NewService(cartRepo repository.CartRepository) CartService {
	return &service{CartRepo: cartRepo}
}

// AddItem implements CartService.
func (s *service) AddItem(ctx context.Context, userID uuid.UUID, req AddCartItemRequest) (*CartItemResponse, error) {
	panic("unimplemented")
}

// ClearItem implements CartService.
func (s *service) ClearItem(ctx context.Context, userID uuid.UUID) error {
	panic("unimplemented")
}

// DeleteItem implements CartService.
func (s *service) DeleteItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) error {
	panic("unimplemented")
}

// ListItem implements CartService.
func (s *service) ListItem(ctx context.Context, userID uuid.UUID) (*CartListResponse, error) {
	panic("unimplemented")
}

// UpdateItem implements CartService.
func (s *service) UpdateItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, req UpdateCartItemRequest) (*CartItemResponse, error) {
	panic("unimplemented")
}
