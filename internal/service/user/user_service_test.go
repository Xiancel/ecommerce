package user

import (
	"context"
	"testing"

	models "github.com/Xiancel/ecommerce/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}
func (m *MockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
func (m *MockUserRepository) List(ctx context.Context, limit, offset int) ([]*models.User, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.User), args.Error(1)
}

func TestGetUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()
	user := &models.User{
		ID:           userID,
		Email:        "test@test.com",
		FirstName:    "TestName",
		LastName:     "TestLastName",
		PasswordHash: "supersecretpassword",
	}

	mockRepo.On("GetByID", ctx, userID).Return(user, nil)

	result, err := service.GetUser(ctx, userID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, user.ID, result.ID)
	mockRepo.AssertExpectations(t)
}

func TestGetUser_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("GetByID", ctx, userID).Return(nil, ErrUserNotFound)

	result, err := service.GetUser(ctx, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrUserNotFound)
	mockRepo.AssertExpectations(t)
}

func TestListUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewService(mockRepo)
	ctx := context.Background()

	users := []*models.User{
		{ID: uuid.New(), FirstName: "FUser 1", LastName: "LUser 1", Email: "u1@example.com"},
		{ID: uuid.New(), FirstName: "FUser 2", LastName: "LUser 2", Email: "u2@example.com"},
	}
	filter := UserFilter{
		Limit:  20,
		Offset: 0,
	}
	mockRepo.On("List", ctx, filter.Limit, filter.Offset).Return(users, nil)

	resp, err := service.ListUser(ctx, filter)

	assert.NoError(t, err)
	assert.Len(t, resp.Users, 2)
	mockRepo.AssertExpectations(t)
}

func TestListUser_Empty(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewService(mockRepo)
	ctx := context.Background()

	filter := UserFilter{
		Limit:  20,
		Offset: 0,
	}
	mockRepo.On("List", ctx, filter.Limit, filter.Offset).Return([]*models.User{}, nil)

	resp, err := service.ListUser(ctx, filter)

	assert.NoError(t, err)
	assert.Len(t, resp.Users, 0)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	firstName := "new UserFN"
	lastName := "new UserLN"
	email := "updated@test.com"

	updateReq := UpdateUserRequest{
		FirstName: &firstName,
		LastName:  &lastName,
		Email:     &email,
	}

	existingUser := &models.User{
		ID:        userID,
		FirstName: "old firstName",
		LastName:  "old lastName",
		Email:     "old@test.com",
	}

	mockRepo.On("GetByID", ctx, userID).Return(existingUser, nil)
	mockRepo.On("Update", ctx, mock.AnythingOfType("*models.User")).Return(nil)

	updatedUser, err := service.UpdateUser(ctx, userID, updateReq, false)

	assert.NoError(t, err)
	assert.Equal(t, updateReq.FirstName, &updatedUser.FirstName)
	assert.Equal(t, updateReq.Email, &updatedUser.Email)
	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	firstName := "new UserFN"
	lastName := "new UserLN"
	email := "updated@test.com"

	updateReq := UpdateUserRequest{
		FirstName: &firstName,
		LastName:  &lastName,
		Email:     &email,
	}

	mockRepo.On("GetByID", ctx, userID).Return(nil, ErrUserNotFound)

	updatedUser, err := service.UpdateUser(ctx, userID, updateReq, false)

	assert.Error(t, err)
	assert.Nil(t, updatedUser)
	assert.Equal(t, ErrUserNotFound, err)
	mockRepo.AssertNotCalled(t, "Update")
}

func TestDeletUser_Succes(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("GetByID", ctx, userID).Return(&models.User{
		ID:        userID,
		FirstName: "User",
		Email:     "test@test.com",
	}, nil)
	mockRepo.On("Delete", ctx, userID).Return(nil)

	err := service.DeleteUser(ctx, userID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteUser_NotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	mockRepo.On("GetByID", ctx, userID).Return(&models.User{
		ID:        userID,
		FirstName: "User",
		Email:     "test@test.com",
	}, nil)
	mockRepo.On("Delete", ctx, userID).Return(ErrUserNotFound)

	err := service.DeleteUser(ctx, userID)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrUserNotFound)
	mockRepo.AssertExpectations(t)
}
