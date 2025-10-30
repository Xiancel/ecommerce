package cart

import (
	"context"

	"github.com/google/uuid"
)

type CartService interface {
	AddItem(ctx context.Context, userID uuid.UUID, req AddCartItemRequest) (*CartItemResponse, error)
	UpdateItem(ctx context.Context, userID, itemID uuid.UUID, req UpdateCartItemRequest) (*CartItemResponse, error)
	DeleteItem(ctx context.Context, userID, itemID uuid.UUID) error
	ListItem(ctx context.Context, userID uuid.UUID) (*CartListResponse, error)
	ClearItem(ctx context.Context, userID uuid.UUID) error
}
