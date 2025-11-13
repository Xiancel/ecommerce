package user

import (
	"context"
	"fmt"
	"strings"

	models "github.com/Xiancel/ecommerce/internal/domain"
	repository "github.com/Xiancel/ecommerce/internal/repository/postgres"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	userRepo repository.UserRepository
}

func NewService(userRepo repository.UserRepository) UserService {
	return &service{userRepo: userRepo}
}

// DeleteUser видалення користувача
func (s *service) DeleteUser(ctx context.Context, id uuid.UUID) error {
	// перевірка користувача на наявність по ID
	_, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return ErrUserNotFound
	}

	// видалення користувача
	if err := s.userRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delte user: %w", err)
	}

	return nil
}

// GetUser отримання користувача
func (s *service) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	// валідація
	if id == uuid.Nil {
		return nil, ErrUserNotFound
	}

	// отримання інформації про користувача
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("faield to get user: %w", err)
	}

	return user, nil
}

// ListUser повертаж список користувачів
func (s *service) ListUser(ctx context.Context, filter UserFilter) (*UserListResponse, error) {
	// пагінація
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}
	// отримання списку користувачів
	users, err := s.userRepo.List(ctx, filter.Limit, filter.Offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	filtered := []*models.User{}
	// отримання користувачів по фільтрам
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
	// формування відповіді
	resp := &UserListResponse{
		Users: filtered,
		Total: len(filtered),
	}

	return resp, nil
}

// UpdateUser новлення данних користувача
func (s *service) UpdateUser(ctx context.Context, id uuid.UUID, req UpdateUserRequest, isAdmin bool) (*models.User, error) {
	// перевірка на наявність користувача
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, ErrUserNotFound
	}
	// валідація
	if req.Email == nil && req.FirstName == nil && req.LastName == nil &&
		req.Password == nil && req.Role == nil {
		return nil, ErrNoFields
	}

	// пагінація/формування оновлених данних
	if req.Email != nil {
		user.Email = *req.Email
	}
	if req.FirstName != nil {
		user.FirstName = *req.FirstName
	}
	if req.LastName != nil {
		user.LastName = *req.LastName
	}
	// формування нового паролля
	if req.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		user.PasswordHash = string(hash)
	}

	if req.Role != nil && isAdmin {
		user.Role = *req.Role
	}

	// оновлення користувача
	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}
