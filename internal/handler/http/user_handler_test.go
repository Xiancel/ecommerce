package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	models "github.com/Xiancel/ecommerce/internal/domain"
	userService "github.com/Xiancel/ecommerce/internal/service/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockUserService) ListUser(ctx context.Context, filter userService.UserFilter) (*userService.UserListResponse, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*userService.UserListResponse), args.Error(1)
}
func (m *MockUserService) UpdateUser(ctx context.Context, id uuid.UUID, req userService.UpdateUserRequest, isAdmin bool) (*models.User, error) {
	args := m.Called(ctx, id, req, isAdmin)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}
func (m *MockUserService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestListUser_Succes(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	userID := uuid.New()

	expected := &models.User{
		ID:           userID,
		Email:        "testuser@test.com",
		PasswordHash: "312l213wewetw6323f324|7.6|:))O",
		FirstName:    "TFName",
		LastName:     "TLName",
		Role:         "customer",
	}

	mockService.On("GetUser", mock.Anything, userID).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req = req.WithContext(context.WithValue(req.Context(), ContextKeyUserID, userID))
	rr := httptest.NewRecorder()

	handler.ListUser(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var user models.User
	err := json.NewDecoder(rr.Body).Decode(&user)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, user.ID)
	assert.Equal(t, expected.FirstName, user.FirstName)
	mockService.AssertExpectations(t)
}

func TestListUser_NotFound(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	userID := uuid.New()

	mockService.On("GetUser", mock.Anything, userID).Return(nil, userService.ErrUserNotFound)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	req = req.WithContext(context.WithValue(req.Context(), ContextKeyUserID, userID))
	rr := httptest.NewRecorder()

	handler.ListUser(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var errResp ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&errResp)
	assert.NoError(t, err)
	assert.Contains(t, errResp.Error, "user not found")
	mockService.AssertExpectations(t)
}

func TestUpdateUser_Succes(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	userID := uuid.New()

	body := `{"email":"new@test.com","first_name":"New","last_name":"User","password":"12345678"}`

	req := httptest.NewRequest(http.MethodPut, "/users", strings.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), ContextKeyUserID, userID))
	rr := httptest.NewRecorder()

	expected := &models.User{
		ID:        userID,
		Email:     "new@test.com",
		FirstName: "New",
		LastName:  "User",
	}

	mockService.On("UpdateUser", mock.Anything, userID, mock.Anything, false).Return(expected, nil)

	handler.UpdateUser(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var user models.User
	err := json.NewDecoder(rr.Body).Decode(&user)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, user.ID)
	assert.Equal(t, expected.FirstName, user.FirstName)
	mockService.AssertExpectations(t)
}

func TestUpdateUser_InvalidBody(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodPut, "/users", strings.NewReader("{invalid body}"))
	req = req.WithContext(context.WithValue(req.Context(), ContextKeyUserID, uuid.New()))
	rr := httptest.NewRecorder()

	handler.UpdateUser(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var errResp ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&errResp)
	assert.NoError(t, err)
	assert.Contains(t, errResp.Error, "Invalid request body")
	mockService.AssertExpectations(t)
}
