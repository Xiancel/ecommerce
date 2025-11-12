package auth

import (
	"context"
	"database/sql"
	"testing"

	models "github.com/Xiancel/ecommerce/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
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

func TestRegisterUser_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewService(mockRepo, "secret")
	ctx := context.Background()

	req := RegisterRequest{
		Email:     "test@test.com",
		FirstName: "testFname",
		LastName:  "testLname",
		Password:  "SuperSecretPassword67",
	}

	mockRepo.On("GetByEmail", ctx, req.Email).Return(nil, sql.ErrNoRows)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.User")).Return(nil)

	resp, err := service.Register(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Email, resp.User.Email)
	mockRepo.AssertExpectations(t)
}

func TestRegisterUser_UserAlreadyExsists(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewService(mockRepo, "secret")
	ctx := context.Background()

	exsistsUser := &models.User{
		ID:    uuid.New(),
		Email: "test@test.com",
	}
	req := RegisterRequest{
		Email:     "test@test.com",
		FirstName: "testFname",
		LastName:  "testLname",
		Password:  "SuperSecretPassword67",
	}

	mockRepo.On("GetByEmail", ctx, req.Email).Return(exsistsUser, nil)

	resp, err := service.Register(ctx, req)

	assert.Nil(t, resp)
	assert.Equal(t, ErrUserAlreadyExists, err)
	mockRepo.AssertExpectations(t)
}

func TestLogin_Succes(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewService(mockRepo, "secret")
	ctx := context.Background()

	password := "SuperSecretPassword123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &models.User{
		ID:           uuid.New(),
		Email:        "test@test.com",
		PasswordHash: string(hashedPassword),
	}

	req := LoginRequset{
		Email:    "test@test.com",
		Password: password,
	}
	mockRepo.On("GetByEmail", ctx, req.Email).Return(user, nil)

	resp, err := service.Login(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, req.Email, resp.User.Email)
	mockRepo.AssertExpectations(t)
}

func TestLogin_InvalidPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewService(mockRepo, "secret")
	ctx := context.Background()

	user := &models.User{
		ID:           uuid.New(),
		Email:        "test@test.com",
		PasswordHash: "123ko12op84gf25g872f38r7w8y0if|6.7z|12",
	}

	req := LoginRequset{
		Email:    "test@test.com",
		Password: "password",
	}
	mockRepo.On("GetByEmail", ctx, req.Email).Return(user, nil)

	resp, err := service.Login(ctx, req)

	assert.Nil(t, resp)
	assert.Equal(t, ErrInvalidCredentials, err)
	mockRepo.AssertExpectations(t)
}
