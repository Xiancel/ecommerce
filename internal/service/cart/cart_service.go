package cart

import (
	"context"
	"fmt"

	models "github.com/Xiancel/ecommerce/internal/domain"
	repository "github.com/Xiancel/ecommerce/internal/repository/postgres"
	"github.com/google/uuid"
)

type service struct {
	CartRepo repository.CartRepository
}

func NewService(cartRepo repository.CartRepository) CartService {
	return &service{CartRepo: cartRepo}
}

// AddItem додавання товару в кошик
func (s *service) AddItem(ctx context.Context, userID uuid.UUID, req AddCartItemRequest) (*models.CartItem, error) {
	// валідація
	if req.ProductID == uuid.Nil {
		return nil, ErrProductNotFound
	}
	if req.Quantity <= 0 {
		return nil, ErrEmptyQuantity
	}

	// отримання товарів
	existItem, err := s.CartRepo.GetItem(ctx, userID, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed get item: %w", err)
	}
	// перевірка на існування товару в кошику
	if existItem != nil {
		// якщо вже такий товар існює оновлюємо кількість
		quant := existItem.Quantity + req.Quantity
		if err := s.CartRepo.UpdateQuantity(ctx, existItem.ID, quant); err != nil {
			return nil, fmt.Errorf("failed to update item quantity: %w", err)
		}
		existItem.Quantity = quant
		return existItem, nil
	}

	// додавання товару в кошик
	item := &models.CartItem{
		ID:        uuid.New(),
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}
	if err := s.CartRepo.AddItem(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to add item: %w", err)
	}

	return item, nil
}

// ClearItem очищення кошику
func (s *service) ClearItem(ctx context.Context, userID uuid.UUID) error {
	// валідація
	if userID == uuid.Nil {
		return ErrUserIDRequired
	}

	// очищення кошика
	if err := s.CartRepo.Clear(ctx, userID); err != nil {
		return fmt.Errorf("failed to clear item: %w", err)
	}
	return nil
}

// DeleteItem видалення товару з кошика
func (s *service) DeleteItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) error {
	// ввалідація
	if userID == uuid.Nil {
		return ErrUserIDRequired
	}
	if itemID == uuid.Nil {
		return ErrItemIDRequired
	}

	// отримання товару за його ID
	existItem, err := s.CartRepo.GetItemByID(ctx, userID, itemID)
	// обробка помилок
	if err != nil {
		return fmt.Errorf("failed get item: %w", err)
	}
	if existItem == nil {
		return ErrProductNotFound
	}

	// видалення товару з кошика
	if err := s.CartRepo.RemoveItem(ctx, userID, itemID); err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}
	return nil
}

// ListItem повененя списку товару у кошику
func (s *service) ListItem(ctx context.Context, userID uuid.UUID) (*CartListResponse, error) {
	// отримання товарів у кошику за ID користувача
	items, err := s.CartRepo.GetByUserId(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart item: %w", err)
	}

	// створення відповіді для корзини
	resp := &CartListResponse{
		Items:      []*models.CartItem{},
		TotalPrice: 0,
	}

	// додаванння товарів у список
	for _, item := range items {
		cartItem := &models.CartItem{
			ID:        item.ID,
			UserID:    item.UserID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		}
		resp.Items = append(resp.Items, cartItem)
		resp.TotalPrice += float64(item.Quantity) * item.ProductPrice
	}

	return resp, nil
}

// UpdateItem оновлення товару в кошику
func (s *service) UpdateItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, req UpdateCartItemRequest) (*models.CartItem, error) {
	// валідація
	if userID == uuid.Nil {
		return nil, ErrUserIDRequired
	}
	if itemID == uuid.Nil {
		return nil, ErrItemIDRequired
	}
	if req.Quantity <= 0 {
		return nil, ErrEmptyQuantity
	}

	// получення товару
	existItem, err := s.CartRepo.GetItem(ctx, userID, itemID)
	// обробка помилок
	if err != nil {
		return nil, fmt.Errorf("failed get item: %w", err)
	}
	if existItem == nil {
		return nil, ErrProductNotFound
	}

	// валідація
	if req.ProductID != nil {
		existItem.ProductID = *req.ProductID
	}
	existItem.Quantity = req.Quantity

	// оновлення кількость товару і корзині
	if err := s.CartRepo.UpdateQuantity(ctx, existItem.ID, existItem.Quantity); err != nil {
		return nil, fmt.Errorf("failed to update quantity: %w", err)
	}

	return existItem, nil
}
