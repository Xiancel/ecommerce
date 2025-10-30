package user

import (
	"context"

	models "github.com/Xiancel/ecommerce/internal/domain"
	repository "github.com/Xiancel/ecommerce/internal/repository/postgres"
	"github.com/google/uuid"
)


type service struct {
	userRepo repository.UserRepository
}

func NewService(userRepo repository.UserRepository) UserService {
	return &service{userRepo: userRepo}
}

// DeleteUser implements UserService.
func (s *service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	panic("unimplemented")
}

// GetUser implements UserService.
func (s *service) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	panic("unimplemented")
}

// ListUser implements UserService.
func (s *service) ListUser(ctx context.Context, filter UserFilter) (*UserListResponse, error) {
	panic("unimplemented")
}

// Login implements UserService.
func (s *service) Login(ctx context.Context, email string, password string) (*models.User, error) {
	panic("unimplemented")
}

// RegisterUser implements UserService.
func (s *service) RegisterUser(ctx context.Context, req CreateUserRequest) (*models.User, error) {
	panic("unimplemented")
}

// UpdateUser implements UserService.
func (s *service) UpdateUser(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (*models.User, error) {
	panic("unimplemented")
}
