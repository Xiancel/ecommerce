package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/Xiancel/ecommerce/internal/domain"
	orderSrv "github.com/Xiancel/ecommerce/internal/service/order"
	productSrv "github.com/Xiancel/ecommerce/internal/service/product"
	userSrv "github.com/Xiancel/ecommerce/internal/service/user"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAdmin_CreateProduct_Success(t *testing.T) {
	mockSrv := new(MockProductService)
	handler := NewAdminHandler(mockSrv, nil, nil)

	reqBody := productSrv.CreateProductRequest{
		Name:  "Test Product",
		Price: 100,
	}
	respBody := &models.Product{
		ID:    uuid.New(),
		Name:  "Test Product",
		Price: 100,
	}

	mockSrv.On("CreateProduct", mock.Anything, reqBody).Return(respBody, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/admin/products", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.CreateProduct(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp models.Product
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, respBody.ID, resp.ID)
	mockSrv.AssertExpectations(t)
}

func TestAdmin_CreateProduct_Error(t *testing.T) {
	mockSrv := new(MockProductService)
	handler := NewAdminHandler(mockSrv, nil, nil)

	reqBody := productSrv.CreateProductRequest{
		Name:  "Test Product",
		Price: 100,
	}

	mockSrv.On("CreateProduct", mock.Anything, reqBody).Return(nil, errors.New("Internal server error"))

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/admin/products", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.CreateProduct(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var resp models.Product
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockSrv.AssertExpectations(t)
}

func TestAdmin_UpdateProduct_Success(t *testing.T) {
	mockSrv := new(MockProductService)
	handler := NewAdminHandler(mockSrv, nil, nil)

	productName := "TestProd"
	productID := uuid.New()
	reqBody := productSrv.UpdateProductRequest{
		Name: &productName,
	}
	respBody := &models.Product{
		ID:   productID,
		Name: "NewProdName",
	}

	mockSrv.On("UpdateProduct", mock.Anything, productID, reqBody).Return(respBody, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/admin/products/"+productID.String(), bytes.NewReader(body))
	rr := httptest.NewRecorder()

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", productID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	handler.UpdateProduct(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp models.Product
	_ = json.NewDecoder(rr.Body).Decode(&resp)
	assert.Equal(t, respBody.Name, resp.Name)
	mockSrv.AssertExpectations(t)
}

func TestAdmin_UpdateProduct_Error(t *testing.T) {
	mockSrv := new(MockProductService)
	handler := NewAdminHandler(mockSrv, nil, nil)

	productID := uuid.New()
	productName := "NewProd"
	reqBody := productSrv.UpdateProductRequest{
		Name: &productName,
	}

	mockSrv.On("UpdateProduct", mock.Anything, productID, reqBody).Return(nil, errors.New("db error"))

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/admin/products/"+productID.String(), bytes.NewReader(body))
	rr := httptest.NewRecorder()

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", productID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	handler.UpdateProduct(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)
	mockSrv.AssertExpectations(t)
}

func TestAdmin_ListAllOrder_Success(t *testing.T) {
	mockSrv := new(MockOrderService)
	handler := NewAdminHandler(nil, mockSrv, nil)

	orders := &orderSrv.OrderListResponse{
		Order: []*models.Order{
			{ID: uuid.New()},
		},
	}

	mockSrv.On("ListOrder", mock.Anything, mock.Anything).Return(orders, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin/orders", nil)
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/admin/orders", handler.ListAllOrder)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockSrv.AssertExpectations(t)
}

func TestAdmin_ListAllOrder_Error_InvalidLimit(t *testing.T) {
	handler := NewAdminHandler(nil, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin/orders?limit=invalid", nil)
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Get("/admin/orders", handler.ListAllOrder)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAdmin_UpdateOrderStatus_Success(t *testing.T) {
	mockSrv := new(MockOrderService)
	handler := NewAdminHandler(nil, mockSrv, nil)

	orderID := uuid.New()
	reqBody := orderSrv.UpdateOrderRequest{Status: "shipped"}
	respBody := &models.Order{ID: orderID, Status: "shipped"}

	mockSrv.On("UpdateOrderStatus", mock.Anything, orderID, reqBody).Return(respBody, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/admin/orders/"+orderID.String()+"/status", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	r := chi.NewRouter()
	r.Put("/admin/orders/{id}/status", handler.UpdateOrderStatus)
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAdmin_GetUser_Success(t *testing.T) {
	mockSrv := new(MockUserService)
	handler := NewAdminHandler(nil, nil, mockSrv)

	userID := uuid.New()
	respBody := &models.User{ID: userID, FirstName: "TestUser"}

	mockSrv.On("GetUser", mock.Anything, userID).Return(respBody, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin/users/"+userID.String(), nil)
	rr := httptest.NewRecorder()

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", userID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	handler.GetUser(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockSrv.AssertExpectations(t)
}

func TestAdmin_GetUser_Error_InvalidID(t *testing.T) {
	handler := NewAdminHandler(nil, nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/admin/users/invalid-id", nil)
	rr := httptest.NewRecorder()

	handler.GetUser(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAdmin_DeleteUser_Success(t *testing.T) {
	mockSrv := new(MockUserService)
	handler := NewAdminHandler(nil, nil, mockSrv)

	userID := uuid.New()

	mockSrv.On("DeleteUser", mock.Anything, userID).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/admin/users/"+userID.String(), nil)
	rr := httptest.NewRecorder()

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", userID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	handler.DeleteUser(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockSrv.AssertExpectations(t)
}

func TestAdmin_UpdateUser_Success(t *testing.T) {
	mockSrv := new(MockUserService)
	handler := NewAdminHandler(nil, nil, mockSrv)

	userName := "TestUser"
	userID := uuid.New()
	reqBody := userSrv.UpdateUserRequest{
		FirstName: &userName,
	}
	respBody := &models.User{
		ID:        userID,
		FirstName: "UpdatedName",
	}

	mockSrv.On("UpdateUser", mock.Anything, userID, reqBody, true).Return(respBody, nil)

	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPut, "/admin/users/"+userID.String(), bytes.NewReader(body))
	rr := httptest.NewRecorder()

	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", userID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
	
	handler.UpdateUser(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	var resp models.User
	_ = json.NewDecoder(rr.Body).Decode(&resp)
	assert.Equal(t, respBody.FirstName, resp.FirstName)
	mockSrv.AssertExpectations(t)
}
