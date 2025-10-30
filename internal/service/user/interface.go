package user

import (
	"context"

	models "github.com/Xiancel/ecommerce/internal/domain"
	"github.com/google/uuid"
)

type UserService interface {
	RegisterUser(ctx context.Context, req CreateUserRequest) (*models.User, error) // поменять *models.User на JWT
	Login(ctx context.Context, email, password string) (*models.User, error)       // поменять *models.User на JWT
	GetUser(ctx context.Context, id uuid.UUID) (*models.User, error)
	ListUser(ctx context.Context, filter UserFilter) (*UserListResponse, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (*models.User, error)
	DeleteUser(ctx context.Context, id uuid.UUID) error
	// добавить Refresh token
}
