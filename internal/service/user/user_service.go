package user

import (
	"context"
	"fmt"

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
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delte user: %w", err)
	}

	return nil
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
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to get email: %w", err)
	}
	if user.PasswordHash != password {
		return nil, ErrInvalidPassword
	}

	return user, nil
}

// RegisterUser implements UserService.
func (s *service) RegisterUser(ctx context.Context, req CreateUserRequest) (*models.User, error) {
	if req.FirstName == "" {
		return nil, ErrUserFNameRequired
	}
	if req.LastName == "" {
		return nil, ErrUserLNameRequired
	}
	if req.Password == "" {
		return nil, ErrUserPasswordRequired
	}
	if req.Email == "" {
		return nil, ErrUserEmailRequired
	}

	if req.Role == "" {
		req.Role = "user"
	} else if req.Role != "user" && req.Role != "admin" {
		return nil, ErrInvalidRole
	}

	existUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to get email: %w", err)
	}
	if existUser != nil {
		return nil, ErrUserAlreadyExists
	}

	user := &models.User{
		Email:        req.Email,
		PasswordHash: req.Password,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         req.Role,
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return user, nil

}

// UpdateUser implements UserService.
func (s *service) UpdateUser(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	if req.Email == nil && req.FirstName == nil && req.LastName == nil &&
		req.Password == nil && req.Role == nil {
		return nil, fmt.Errorf("no fields to update")
	}

	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	if req.Password != nil {
		user.PasswordHash = *req.Password
	}
	if req.Role != nil {
		user.Role = *req.Role
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}
