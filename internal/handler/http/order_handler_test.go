package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	models "github.com/Xiancel/ecommerce/internal/domain"
	orderService "github.com/Xiancel/ecommerce/internal/service/order"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) CreateOrder(ctx context.Context, userID uuid.UUID, req orderService.CreateOrderRequest) (*models.Order, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}
func (m *MockOrderService) GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}
func (m *MockOrderService) ListOrder(ctx context.Context, filter orderService.OrderFilter) (*orderService.OrderListResponse, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*orderService.OrderListResponse), args.Error(1)
}
func (m *MockOrderService) UpdateOrderStatus(ctx context.Context, id uuid.UUID, req orderService.UpdateOrderRequest) (*models.Order, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}
func (m *MockOrderService) CancelOrder(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestGetOrder_Succes(t *testing.T) {
	mockService := new(MockOrderService)
	handler := NewOrderHandler(mockService)

	orderID := uuid.New()
	expected := &models.Order{ID: orderID}

	mockService.On("GetOrder", mock.Anything, orderID).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/orders/"+orderID.String(), nil)
	rr := httptest.NewRecorder()

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", orderID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	handler.GetOrder(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var order models.Order
	err := json.NewDecoder(rr.Body).Decode(&order)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, order.ID)
	mockService.AssertExpectations(t)
}
func TestGetOrder_NotFound(t *testing.T) {
	mockSrv := new(MockOrderService)
	handler := NewOrderHandler(mockSrv)

	orderID := uuid.New()
	mockSrv.On("GetOrder", mock.Anything, orderID).Return(nil, orderService.ErrOrderNotFound)

	req := httptest.NewRequest(http.MethodGet, "/orders/"+orderID.String(), nil)
	rr := httptest.NewRecorder()

	handler.GetOrder(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCreateOrder_Success(t *testing.T) {
	mockSrv := new(MockOrderService)
	handler := NewOrderHandler(mockSrv)

	userID := uuid.New()
	orderReq := orderService.CreateOrderRequest{}
	expected := &models.Order{ID: uuid.New()}

	mockSrv.On("CreateOrder", mock.Anything, userID, orderReq).Return(expected, nil)

	body, _ := json.Marshal(orderReq)
	req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
	req = req.WithContext(context.WithValue(req.Context(), ContextKeyUserID, userID))
	rr := httptest.NewRecorder()

	handler.CreateOrder(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	var resp models.Order
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, expected.ID, resp.ID)
}
func TestCreateOrder_InvalidBody(t *testing.T) {
	mockSrv := new(MockOrderService)
	handler := NewOrderHandler(mockSrv)

	userID := uuid.New()
	req := httptest.NewRequest(http.MethodPost, "/orders", strings.NewReader("invalid json"))
	req = req.WithContext(context.WithValue(req.Context(), ContextKeyUserID, userID))
	rr := httptest.NewRecorder()

	handler.CreateOrder(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCancelOrder_Success(t *testing.T) {
	mockSrv := new(MockOrderService)
	handler := NewOrderHandler(mockSrv)

	orderID := uuid.New()
	mockSrv.On("CancelOrder", mock.Anything, orderID).Return(nil)

	req := httptest.NewRequest(http.MethodPut, "/orders/"+orderID.String()+"/cancel", nil)
	rr := httptest.NewRecorder()

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", orderID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	handler.CancelOrder(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestListOrder_Success(t *testing.T) {
	mockSrv := new(MockOrderService)
	handler := NewOrderHandler(mockSrv)

	userID := uuid.New()
	expected := &orderService.OrderListResponse{}
	mockSrv.On("ListOrder", mock.Anything, mock.AnythingOfType("order.OrderFilter")).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	req = req.WithContext(context.WithValue(req.Context(), ContextKeyUserID, userID))
	rr := httptest.NewRecorder()

	handler.ListOrder(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestListOrder_InvalidLimit(t *testing.T) {
	mockSrv := new(MockOrderService)
	handler := NewOrderHandler(mockSrv)

	userID := uuid.New()
	req := httptest.NewRequest(http.MethodGet, "/orders?limit=invalid", nil)
	req = req.WithContext(context.WithValue(req.Context(), ContextKeyUserID, userID))
	rr := httptest.NewRecorder()

	handler.ListOrder(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
