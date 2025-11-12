package cart

import (
	"context"
	"testing"

	models "github.com/Xiancel/ecommerce/internal/domain"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockCartRepository struct {
	mock.Mock
}

func (m *MockCartRepository) AddItem(ctx context.Context, item *models.CartItem) error {
	args := m.Called(ctx, item)
	return args.Error(0)
}
func (m *MockCartRepository) GetByUserId(ctx context.Context, userID uuid.UUID) ([]*models.CartItemWithProduct, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.CartItemWithProduct), args.Error(1)
}
func (m *MockCartRepository) UpdateQuantity(ctx context.Context, id uuid.UUID, quantity int) error {
	args := m.Called(ctx, id, quantity)
	return args.Error(0)
}
func (m *MockCartRepository) RemoveItem(ctx context.Context, userID, id uuid.UUID) error {
	args := m.Called(ctx, userID, id)
	return args.Error(0)
}
func (m *MockCartRepository) Clear(ctx context.Context, userId uuid.UUID) error {
	args := m.Called(ctx, userId)
	return args.Error(0)
}
func (m *MockCartRepository) GetItem(ctx context.Context, userId, productId uuid.UUID) (*models.CartItem, error) {
	args := m.Called(ctx, userId, productId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CartItem), args.Error(1)
}
func (m *MockCartRepository) GetItemByID(ctx context.Context, userID, itemID uuid.UUID) (*models.CartItem, error) {
	args := m.Called(ctx, userID, itemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CartItem), args.Error(1)
}

func TestAddItem_Success(t *testing.T) {
	mockRepo := new(MockCartRepository)
	service := NewService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()
	productID := uuid.New()

	req := AddCartItemRequest{
		ProductID: productID,
		Quantity:  2,
	}

	mockRepo.On("GetItem", ctx, userID, productID).Return(nil, nil)
	mockRepo.On("AddItem", ctx, mock.AnythingOfType("*models.CartItem")).Return(nil)

	item, err := service.AddItem(ctx, userID, req)

	assert.NoError(t, err)
	assert.NotNil(t, item)
	assert.Equal(t, req.Quantity, item.Quantity)
	assert.Equal(t, productID, item.ProductID)
	mockRepo.AssertExpectations(t)
}

func TestAddItem_AlreadyExist(t *testing.T) {
	mockRepo := new(MockCartRepository)
	service := NewService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()
	productID := uuid.New()

	existingItem := &models.CartItem{
		ID:        uuid.New(),
		UserID:    userID,
		ProductID: productID,
		Quantity:  1,
	}

	req := AddCartItemRequest{
		ProductID: productID,
		Quantity:  2,
	}

	mockRepo.On("GetItem", ctx, userID, productID).Return(existingItem, nil)
	mockRepo.On("UpdateQuantity", ctx, existingItem.ID, 3).Return(nil)

	item, err := service.AddItem(ctx, userID, req)

	assert.NoError(t, err)
	assert.Equal(t, 3, item.Quantity)
	mockRepo.AssertExpectations(t)
}

func TestUpdateItem_Success(t *testing.T) {
	mockRepo := new(MockCartRepository)
	service := NewService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()
	itemID := uuid.New()

	existingItem := &models.CartItem{
		ID:        itemID,
		UserID:    userID,
		ProductID: uuid.New(),
		Quantity:  1,
	}

	req := UpdateCartItemRequest{
		Quantity: 5,
	}

	mockRepo.On("GetItem", ctx, userID, itemID).Return(existingItem, nil)
	mockRepo.On("UpdateQuantity", ctx, itemID, req.Quantity).Return(nil)

	item, err := service.UpdateItem(ctx, userID, itemID, req)

	assert.NoError(t, err)
	assert.Equal(t, req.Quantity, item.Quantity)
	mockRepo.AssertExpectations(t)
}

func TestUpdateItem_NotFound(t *testing.T) {
	mockRepo := new(MockCartRepository)
	service := NewService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()
	itemID := uuid.New()

	req := UpdateCartItemRequest{
		Quantity: 5,
	}

	mockRepo.On("GetItem", ctx, userID, itemID).Return(nil, ErrItemNotFound)

	item, err := service.UpdateItem(ctx, userID, itemID, req)

	assert.ErrorIs(t, err, ErrItemNotFound)
	assert.Nil(t, item)
	mockRepo.AssertExpectations(t)
}

func TestDeleteItem_Success(t *testing.T) {
	mockRepo := new(MockCartRepository)
	service := NewService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()
	itemID := uuid.New()

	existingItem := &models.CartItem{
		ID:        itemID,
		UserID:    userID,
		ProductID: uuid.New(),
		Quantity:  1,
	}
	mockRepo.On("GetItemByID", ctx, userID, itemID).Return(existingItem, nil)
	mockRepo.On("RemoveItem", ctx, userID, itemID).Return(nil)

	err := service.DeleteItem(ctx, userID, itemID)

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestListItem_Success(t *testing.T) {
	mockRepo := new(MockCartRepository)
	service := NewService(mockRepo)
	ctx := context.Background()
	userID := uuid.New()

	items := []*models.CartItemWithProduct{
		{
			CartItem: models.CartItem{
				ID:        uuid.New(),
				UserID:    userID,
				ProductID: uuid.New(),
				Quantity:  2,
			},
			ProductName:  "Product1",
			ProductPrice: 50,
			ProductStock: 10,
		},
	}

	mockRepo.On("GetByUserId", ctx, userID).Return(items, nil)

	resp, err := service.ListItem(ctx, userID)

	assert.NoError(t, err)
	assert.Len(t, resp.Items, 1)
	mockRepo.AssertExpectations(t)
}
