package user

import (
	"context"
	"fmt"
	"strings"

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
	if id == uuid.Nil {
		return nil, ErrUserNotFound
	}

	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("faield to get user: %w", err)
	}

	return user, nil
}

// ListUser implements UserService.
func (s *service) ListUser(ctx context.Context, filter UserFilter) (*UserListResponse, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	users, err := s.userRepo.List(ctx, filter.Limit, filter.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	filtered := []*models.User{}
	for _, u := range users {
		if filter.Role != nil && u.Role != *filter.Role {
			continue
		}
		if filter.Search != "" && !strings.Contains(strings.ToLower(u.FirstName),
			strings.ToLower(filter.Search)) && !strings.Contains(strings.ToLower(u.LastName), strings.ToLower(filter.Search)) {
			continue
		}
		filtered = append(filtered, u)
	}
	resp := &UserListResponse{
		Users: filtered,
		Total: len(filtered),
	}

	return resp, nil
}

// UpdateUser implements UserService.
func (s *service) UpdateUser(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (*models.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	if req.Email == nil && req.FirstName == nil && req.LastName == nil &&
		req.Password == nil && req.Role == nil {
		return nil, ErrNoFields
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
