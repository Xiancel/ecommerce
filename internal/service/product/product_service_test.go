package product

import (
	"context"
	"testing"

	models "github.com/Xiancel/ecommerce/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockProductRepository struct {
	mock.Mock
}

func (m *MockProductRepository) Create(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) GetById(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Product), args.Error(1)
}

func (m *MockProductRepository) List(ctx context.Context, filter models.ListFilter) ([]*models.Product, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Product), args.Error(1)
}

func (m *MockProductRepository) Update(ctx context.Context, product *models.Product) error {
	args := m.Called(ctx, product)
	return args.Error(0)
}

func (m *MockProductRepository) UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error {
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}

func (m *MockProductRepository) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateProduct_Success(t *testing.T) {
	//Arrange
	mockRepo := new(MockProductRepository)
	service := NewService(mockRepo)
	ctx := context.Background()
	req := CreateProductRequest{
		Name:        "Test Product",
		Description: "Test description",
		Price:       99.99,
		Stock:       10,
	}

	mockRepo.On("Create", ctx, mock.AnythingOfType("*models.Product")).Return(nil)
	//Act
	product, err := service.CreateProduct(ctx, req)

	//Assert
	assert.NoError(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, req.Name, product.Name)
	assert.Equal(t, req.Stock, product.Stock)
	assert.Equal(t, req.Price, product.Price)
	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_EmptyName(t *testing.T) {
	//Arrange
	mockRepo := new(MockProductRepository)
	service := NewService(mockRepo)
	ctx := context.Background()
	req := CreateProductRequest{
		Name:  "",
		Price: 99.99,
		Stock: 10,
	}

	//Act
	product, err := service.CreateProduct(ctx, req)

	//Assert
	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Equal(t, ErrProductNameRequired, err)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateProduct_InvalidPrice(t *testing.T) {
	//Arrange
	mockRepo := new(MockProductRepository)
	service := NewService(mockRepo)
	ctx := context.Background()
	req := CreateProductRequest{
		Name:  "Test Product",
		Price: -99.99,
		Stock: 10,
	}

	//Act
	product, err := service.CreateProduct(ctx, req)

	//Assert
	assert.Error(t, err)
	assert.Nil(t, product)
	assert.Equal(t, ErrInvalidPrice, err)
	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateProduct_NotAvailable(t *testing.T) {
	//Arrange
	mockRepo := new(MockProductRepository)
	service := NewService(mockRepo)
	ctx := context.Background()

	productID := uuid.New()
	product := &models.Product{
		ID:    productID,
		Name:  "Test Product",
		Stock: 3,
	}

	mockRepo.On("GetById", ctx, productID).Return(product, nil)
	//Act
	available, err := service.CheckAvailability(ctx, productID, 5)

	//Assert
	assert.NoError(t, err)
	assert.False(t, available)

	mockRepo.AssertExpectations(t)
}

func TestCreateProduct_NotQuantity(t *testing.T) {
	//Arrange
	mockRepo := new(MockProductRepository)
	service := NewService(mockRepo)
	ctx := context.Background()

	productID := uuid.New()

	//Act
	available, err := service.CheckAvailability(ctx, productID, 0)

	//Assert
	assert.Error(t, err)
	assert.False(t, available)
	assert.Equal(t, ErrInvalidQuantity, err)

	mockRepo.AssertNotCalled(t, "GetById")
}
