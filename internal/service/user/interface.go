package user

import (
	"context"

	models "github.com/Xiancel/ecommerce/internal/domain"
	"github.com/google/uuid"
)

type UserService interface {
	GetUser(ctx context.Context, id uuid.UUID) (*models.User, error)
	ListUser(ctx context.Context, filter UserFilter) (*UserListResponse, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req UpdateUserRequest, isAdmin bool) (*models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
}
