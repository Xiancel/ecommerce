package cart

import (
	"context"

	models "github.com/Xiancel/ecommerce/internal/domain"
	"github.com/google/uuid"
)

type CartService interface {
	AddItem(ctx context.Context, userID uuid.UUID, req AddCartItemRequest) (*models.CartItem, error)
	UpdateItem(ctx context.Context, userID, itemID uuid.UUID, req UpdateCartItemRequest) (*models.CartItem, error)
	DeleteItem(ctx context.Context, userID, itemID uuid.UUID) error
	ListItem(ctx context.Context, userID uuid.UUID) (*CartListResponse, error)
	ClearItem(ctx context.Context, userID uuid.UUID) error
}
