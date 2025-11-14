package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	models "github.com/Xiancel/ecommerce/internal/domain"
	cartService "github.com/Xiancel/ecommerce/internal/service/cart"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCartService struct {
	mock.Mock
}

func (m *MockCartService) AddItem(ctx context.Context, userID uuid.UUID, req cartService.AddCartItemRequest) (*models.CartItem, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CartItem), args.Error(1)
}
func (m *MockCartService) UpdateItem(ctx context.Context, userID, itemID uuid.UUID, req cartService.UpdateCartItemRequest) (*models.CartItem, error) {
	args := m.Called(ctx, userID, itemID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CartItem), args.Error(1)
}
func (m *MockCartService) DeleteItem(ctx context.Context, userID, itemID uuid.UUID) error {
	args := m.Called(ctx, userID, itemID)
	return args.Error(0)
}
func (m *MockCartService) ListItem(ctx context.Context, userID uuid.UUID) (*cartService.CartListResponse, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cartService.CartListResponse), args.Error(1)
}
func (m *MockCartService) ClearItem(ctx context.Context, userID uuid.UUID) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func TestListItems_Succes(t *testing.T) {
	mockService := new(MockCartService)
	handler := NewCartHandler(mockService)

	userID := uuid.New()

	expected := &cartService.CartListResponse{
		Items:      []*models.CartItem{},
		TotalPrice: 0,
	}

	mockService.On("ListItem", mock.Anything, userID).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/cart", nil)
	req = req.WithContext(context.WithValue(req.Context(), ContextKeyUserID, userID))
	rr := httptest.NewRecorder()

	handler.ListItems(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var cart models.CartItem
	err := json.NewDecoder(rr.Body).Decode(&cart)
	assert.NoError(t, err)
	mockService.AssertExpectations(t)
}

func TestListItems_Error(t *testing.T) {
	mockSrv := new(MockCartService)
	handler := NewCartHandler(mockSrv)

	userID := uuid.New()

	mockSrv.On("ListItem", mock.Anything, userID).
		Return(nil, errors.New("internal server error"))

	req := httptest.NewRequest(http.MethodGet, "/cart", nil)
	req = req.WithContext(context.WithValue(req.Context(), ContextKeyUserID, userID))
	rr := httptest.NewRecorder()

	handler.ListItems(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockSrv.AssertExpectations(t)
}

func TestAddItem_Success(t *testing.T) {
	mockSrv := new(MockCartService)
	handler := NewCartHandler(mockSrv)

	userID := uuid.New()

	reqBody := `{"product_id":"` + uuid.New().String() + `","quantity":2}`

	item := &models.CartItem{Quantity: 2}

	mockSrv.On("AddItem", mock.Anything, userID, mock.Anything).
		Return(item, nil)

	req := httptest.NewRequest(http.MethodPost, "/cart/items", strings.NewReader(reqBody))
	req = req.WithContext(context.WithValue(req.Context(), ContextKeyUserID, userID))

	rr := httptest.NewRecorder()
	handler.AddItem(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	mockSrv.AssertExpectations(t)
}

func TestAddItem_BadRequest(t *testing.T) {
	mockSrv := new(MockCartService)
	handler := NewCartHandler(mockSrv)

	userID := uuid.New()

	req := httptest.NewRequest(http.MethodPost, "/cart/items", strings.NewReader("invalid json"))
	req = req.WithContext(context.WithValue(req.Context(), ContextKeyUserID, userID))
	rr := httptest.NewRecorder()

	handler.AddItem(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUpdateItem_Success(t *testing.T) {
	mockSrv := new(MockCartService)
	handler := NewCartHandler(mockSrv)

	userID := uuid.New()
	itemID := uuid.New()

	reqBody := `{"quantity":3}`

	expected := &models.CartItem{Quantity: 3}

	mockSrv.On("UpdateItem", mock.Anything, userID, itemID, mock.Anything).
		Return(expected, nil)

	req := httptest.NewRequest(http.MethodPut, "/cart/items/"+itemID.String(), strings.NewReader(reqBody))
	req = req.WithContext(context.WithValue(req.Context(), ContextKeyUserID, userID))
	rr := httptest.NewRecorder()

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", itemID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	handler.UpdateItem(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockSrv.AssertExpectations(t)
}

func TestUpdateItem_InvalidID(t *testing.T) {
	mockSrv := new(MockCartService)
	handler := NewCartHandler(mockSrv)

	userID := uuid.New()

	req := httptest.NewRequest(http.MethodPut, "/cart/items/invalid-id", strings.NewReader(`{}`))
	req = req.WithContext(context.WithValue(req.Context(), ContextKeyUserID, userID))

	rr := httptest.NewRecorder()

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", "invalid-id")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	handler.UpdateItem(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestDeleteItem_Success(t *testing.T) {
	mockSrv := new(MockCartService)
	handler := NewCartHandler(mockSrv)

	userID := uuid.New()
	itemID := uuid.New()

	mockSrv.On("DeleteItem", mock.Anything, userID, itemID).
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/cart/items/"+itemID.String(), nil)
	req = req.WithContext(context.WithValue(req.Context(), ContextKeyUserID, userID))

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", itemID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	rr := httptest.NewRecorder()
	handler.DeleteItem(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockSrv.AssertExpectations(t)
}

func TestDeleteItem_NotFound(t *testing.T) {
	mockSrv := new(MockCartService)
	handler := NewCartHandler(mockSrv)

	userID := uuid.New()
	itemID := uuid.New()

	mockSrv.On("DeleteItem", mock.Anything, userID, itemID).
		Return(cartService.ErrItemNotFound)

	req := httptest.NewRequest(http.MethodDelete, "/cart/items/"+itemID.String(), nil)
	req = req.WithContext(context.WithValue(req.Context(), ContextKeyUserID, userID))

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", itemID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	rr := httptest.NewRecorder()
	handler.DeleteItem(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	mockSrv.AssertExpectations(t)
}

func TestClearCart_Success(t *testing.T) {
	mockSrv := new(MockCartService)
	handler := NewCartHandler(mockSrv)

	userID := uuid.New()

	mockSrv.On("ClearItem", mock.Anything, userID).
		Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/cart", nil)
	req = req.WithContext(context.WithValue(req.Context(), ContextKeyUserID, userID))

	rr := httptest.NewRecorder()
	handler.ClearCart(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockSrv.AssertExpectations(t)
}
