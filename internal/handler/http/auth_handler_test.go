package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	authService "github.com/Xiancel/ecommerce/internal/service/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockAuthService struct {
	mock.Mock
}

func (m *MockAuthService) Register(ctx context.Context, req authService.RegisterRequest) (*authService.AuthResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authService.AuthResponse), args.Error(1)
}
func (m *MockAuthService) Login(ctx context.Context, req authService.LoginRequset) (*authService.AuthResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authService.AuthResponse), args.Error(1)
}
func (m *MockAuthService) RefreshToken(ctx context.Context, refreshToken string) (*authService.AuthResponse, error) {
	args := m.Called(ctx, refreshToken)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authService.AuthResponse), args.Error(1)
}
func (m *MockAuthService) ValidateToken(tokenString string) (*authService.Claims, error) {
	args := m.Called(tokenString)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*authService.Claims), args.Error(1)
}

func TestRegister_Success(t *testing.T) {
	mockSrv := new(MockAuthService)
	handler := NewAuthHandler(mockSrv)

	reqBody := authService.RegisterRequest{
		Email:    "test@test.com",
		Password: "supersecretpassword30",
	}
	respBody := &authService.AuthResponse{
		AccessToken:  "access",
		RefreshToken: "refresh",
	}
	mockSrv.On("Register", mock.Anything, reqBody).Return(respBody, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var resp authService.AuthResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, respBody.AccessToken, resp.AccessToken)
	mockSrv.AssertExpectations(t)
}

func TestRegister_Error(t *testing.T) {
	mockSrv := new(MockAuthService)
	handler := NewAuthHandler(mockSrv)

	reqBody := authService.RegisterRequest{
		Email:    "exists@test.com",
		Password: "supersecretpassword30",
	}
	mockSrv.On("Register", mock.Anything, reqBody).Return(nil, authService.ErrUserAlreadyExists)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.Register(rr, req)

	assert.Equal(t, http.StatusConflict, rr.Code)
	mockSrv.AssertExpectations(t)
}
func TestLogin_Success(t *testing.T) {
	mockSrv := new(MockAuthService)
	handler := NewAuthHandler(mockSrv)

	reqBody := authService.LoginRequset{
		Email:    "test@test.com",
		Password: "supersecretpassword30",
	}
	respBody := &authService.AuthResponse{
		AccessToken:  "access",
		RefreshToken: "refresh",
	}

	mockSrv.On("Login", mock.Anything, reqBody).Return(respBody, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp authService.AuthResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, respBody.AccessToken, resp.AccessToken)
	mockSrv.AssertExpectations(t)
}

func TestLogin_Error(t *testing.T) {
	mockSrv := new(MockAuthService)
	handler := NewAuthHandler(mockSrv)

	reqBody := authService.LoginRequset{
		Email:    "wrong@test.com",
		Password: "wrong",
	}
	mockSrv.On("Login", mock.Anything, reqBody).Return(nil, authService.ErrInvalidCredentials)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.Login(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	mockSrv.AssertExpectations(t)
}

func TestRefreshToken_Success(t *testing.T) {
	mockSrv := new(MockAuthService)
	handler := NewAuthHandler(mockSrv)

	reqBody := authService.RefreshRequest{
		RefreshToken: "refresh_token",
	}
	respBody := &authService.AuthResponse{
		AccessToken:  "new_access",
		RefreshToken: "new_refresh",
	}

	mockSrv.On("RefreshToken", mock.Anything, "refresh_token").Return(respBody, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.RefreshToken(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp authService.AuthResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, respBody.AccessToken, resp.AccessToken)
	mockSrv.AssertExpectations(t)
}

func TestRefreshToken_Error(t *testing.T) {
	mockSrv := new(MockAuthService)
	handler := NewAuthHandler(mockSrv)

	reqBody := authService.RefreshRequest{
		RefreshToken: "expired_token",
	}
	mockSrv.On("RefreshToken", mock.Anything, "expired_token").Return(nil, authService.ErrExpiredToken)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.RefreshToken(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	mockSrv.AssertExpectations(t)
}
