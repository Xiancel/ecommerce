package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/Xiancel/ecommerce/internal/domain"
	productService "github.com/Xiancel/ecommerce/internal/service/product"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductService struct {
	mock.Mock
}

func (m *MockProductService) CreateProduct(ctx context.Context, req productService.CreateProductRequest) (*models.Product, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}
func (m *MockProductService) GetProduct(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}
func (m *MockProductService) ListProduct(ctx context.Context, filter productService.ProductFilter) (*productService.ProductListResponse, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*productService.ProductListResponse), args.Error(1)
}
func (m *MockProductService) SearchProduct(ctx context.Context, query string, limit, offset int) ([]*models.Product, error) {
	args := m.Called(ctx, query, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Product), args.Error(1)
}
func (m *MockProductService) UpdateProduct(ctx context.Context, id uuid.UUID, req productService.UpdateProductRequest) (*models.Product, error) {
	args := m.Called(ctx, id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}
func (m *MockProductService) CheckAvailability(ctx context.Context, id uuid.UUID, quantity int) (bool, error) {
	args := m.Called(ctx, id, quantity)
	return args.Bool(0), args.Error(1)
}
func (m *MockProductService) ReserveStock(ctx context.Context, id uuid.UUID, quantity int) error {
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}
func (m *MockProductService) ReleaseStock(ctx context.Context, id uuid.UUID, quantity int) error {
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

// GetProduct
func TestGetProduct_Success(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	productID := uuid.New()

	expectedProduct := &models.Product{
		ID:    productID,
		Name:  "Test Product",
		Price: 6.7,
		Stock: 10,
	}

	mockService.On("GetProduct", mock.Anything, productID).Return(expectedProduct, nil)

	req := httptest.NewRequest(http.MethodGet, "/products/"+productID.String(), nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", productID.String())
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetProduct(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var product models.Product
	err := json.NewDecoder(rr.Body).Decode(&product)
	assert.NoError(t, err)
	assert.Equal(t, expectedProduct.ID, product.ID)
	assert.Equal(t, expectedProduct.Name, product.Name)
	mockService.AssertExpectations(t)
}

func TestGetProduct_InvalidID(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/products/invalid-id", nil)
	rr := httptest.NewRecorder()

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "invalid-id")
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))

	handler.GetProduct(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var errResp ErrorResponse
	err := json.NewDecoder(rr.Body).Decode(&errResp)
	assert.NoError(t, err)
	assert.Contains(t, errResp.Error, "InvalidID")
	mockService.AssertNotCalled(t, "GetProduct")
}

// ListProducts

func TestListProducts_Success(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	filter := productService.ProductFilter{
		Limit:  20,
		Offset: 0,
	}

	expected := &productService.ProductListResponse{
		Products: []*models.Product{
			{ID: uuid.New(), Name: "PT1", Price: 10, Stock: 5},
			{ID: uuid.New(), Name: "PT2", Price: 6, Stock: 7},
		},
		Total: 2,
	}

	mockService.On("ListProduct", mock.Anything, filter).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/products", nil)
	rr := httptest.NewRecorder()

	handler.ListProducts(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp productService.ProductListResponse
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, 2, resp.Total)
	mockService.AssertExpectations(t)
}

func TestListProducts_InvalidMinPrice(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/products?min_price=abc", nil)
	rr := httptest.NewRecorder()

	handler.ListProducts(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var errResp ErrorResponse
	_ = json.NewDecoder(rr.Body).Decode(&errResp)
	assert.Contains(t, errResp.Error, "Invalid min_price")
	mockService.AssertNotCalled(t, "ListProduct")
}

func TestSearchProduct_Success(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	products := []*models.Product{
		{ID: uuid.New(), Name: "PT1", Price: 5, Stock: 2},
	}

	mockService.On("SearchProduct", mock.Anything, "test", 20, 0).Return(products, nil)

	req := httptest.NewRequest(http.MethodGet, "/products/search?q=test", nil)
	rr := httptest.NewRecorder()

	handler.SearchProduct(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp []*models.Product
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 1)

	mockService.AssertExpectations(t)
}

func TestSearchProduct_MissingQuery(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/products/search", nil)
	rr := httptest.NewRecorder()

	handler.SearchProduct(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)

	var errResp ErrorResponse
	_ = json.NewDecoder(rr.Body).Decode(&errResp)
	assert.Contains(t, errResp.Error, "Search query is required")

	mockService.AssertNotCalled(t, "SearchProduct")
}

func TestListCategories_Success(t *testing.T) {
	mockService := new(MockProductService)
	handler := NewProductHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	rr := httptest.NewRecorder()

	handler.ListCategories(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var resp []interface{}
	err := json.NewDecoder(rr.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Len(t, resp, 0) 
}
