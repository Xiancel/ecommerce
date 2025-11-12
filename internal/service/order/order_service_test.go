package order

import (
	"context"
	"testing"

	models "github.com/Xiancel/ecommerce/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOrderRepository struct {
	mock.Mock
}

func (m *MockOrderRepository) Create(ctx context.Context, order *models.Order, items []*models.OrderItem) error {
	args := m.Called(ctx, order, items)
	return args.Error(0)
}

func (m *MockOrderRepository) GetById(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Order), args.Error(1)
}
func (m *MockOrderRepository) GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]*models.OrderItem, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.OrderItem), args.Error(1)
}
func (m *MockOrderRepository) ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Order, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Order), args.Error(1)
}
func (m *MockOrderRepository) ListAll(ctx context.Context, limit, offset int) ([]*models.Order, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Order), args.Error(1)
}
func (m *MockOrderRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

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

func TestGetOrder_Success(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	mockRepoProduct := new(MockProductRepository)
	service := NewService(mockRepo, mockRepoProduct)
	ctx := context.Background()
	orderID := uuid.New()

	order := &models.Order{
		ID:          orderID,
		Status:      "pending",
		TotalAmount: 100,
	}

	mockRepo.On("GetById", ctx, orderID).Return(order, nil)

	result, err := service.GetOrder(ctx, orderID)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, order.ID, result.ID)
	mockRepo.AssertExpectations(t)
}

func TestGetOrder_NotFound(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	mockRepoProduct := new(MockProductRepository)
	service := NewService(mockRepo, mockRepoProduct)
	ctx := context.Background()
	orderID := uuid.New()

	mockRepo.On("GetById", ctx, orderID).Return(nil, ErrOrderNotFound)

	result, err := service.GetOrder(ctx, orderID)

	assert.Nil(t, result)
	assert.ErrorIs(t, err, ErrOrderNotFound)
	mockRepo.AssertExpectations(t)
}

func TestListOrder_Success(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	mockRepoProduct := new(MockProductRepository)
	service := NewService(mockRepo, mockRepoProduct)
	ctx := context.Background()

	filter := OrderFilter{
		Limit:  10,
		Offset: 0,
	}

	orders := []*models.Order{
		{ID: uuid.New(), Status: "pending"},
		{ID: uuid.New(), Status: "completed"},
	}

	mockRepo.On("ListAll", ctx, filter.Limit, filter.Offset).Return(orders, nil)

	resp, err := service.ListOrder(ctx, filter)

	assert.NoError(t, err)
	assert.Len(t, resp.Order, 2)
	mockRepo.AssertExpectations(t)
}

func TestCancelOrder_Success(t *testing.T) {
	mockRepo := new(MockOrderRepository)
	mockRepoProduct := new(MockProductRepository)
	service := NewService(mockRepo, mockRepoProduct)
	ctx := context.Background()
	orderID := uuid.New()

	order := &models.Order{
		ID:     orderID,
		Status: "pending",
	}

	mockRepo.On("GetById", ctx, orderID).Return(order, nil)
	mockRepo.On("UpdateStatus", ctx, orderID, "cancelled").Return(nil)

	err := service.CancelOrder(ctx, orderID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}
