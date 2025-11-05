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

// AddItem implements CartService.
func (s *service) AddItem(ctx context.Context, userID uuid.UUID, req AddCartItemRequest) (*models.CartItem, error) {
	if req.ProductID == uuid.Nil {
		return nil, ErrProductNotFound
	}
	if req.Quantity <= 0 {
		return nil, ErrEmptyQuantity
	}

	existItem, err := s.CartRepo.GetItem(ctx, userID, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("failed get item: %w", err)
	}
	if existItem != nil {
		quant := existItem.Quantity + req.Quantity
		if err := s.CartRepo.UpdateQuantity(ctx, existItem.ID, quant); err != nil {
			return nil, fmt.Errorf("failed to update item quantity: %w", err)
		}
		existItem.Quantity = quant
		return existItem, nil
	}

	item := &models.CartItem{
		UserID:    userID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}

	if err := s.CartRepo.AddItem(ctx, item); err != nil {
		return nil, fmt.Errorf("failed to add item: %w", err)
	}

	return item, nil
}

// ClearItem implements CartService.
func (s *service) ClearItem(ctx context.Context, userID uuid.UUID) error {
	if userID == uuid.Nil {
		return ErrUserIDRequired
	}

	if err := s.CartRepo.Clear(ctx, userID); err != nil {
		return fmt.Errorf("failed to clear item: %w", err)
	}
	return nil
}

// DeleteItem implements CartService.
func (s *service) DeleteItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) error {
	if userID == uuid.Nil {
		return ErrUserIDRequired
	}
	if itemID == uuid.Nil {
		return ErrItemIDRequired
	}

	existItem, err := s.CartRepo.GetItem(ctx, userID, itemID)
	if err != nil {
		return fmt.Errorf("failed get item: %w", err)
	}
	if existItem == nil {
		return ErrProductNotFound
	}

	if err := s.CartRepo.RemoveItem(ctx, itemID); err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}
	return nil
}

// ListItem implements CartService.
func (s *service) ListItem(ctx context.Context, userID uuid.UUID) (*CartListResponse, error) {
	items, err := s.CartRepo.GetByUserId(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart item: %w", err)
	}

	resp := &CartListResponse{
		Items:      []*models.CartItem{},
		TotalPrice: 0,
	}

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

// UpdateItem implements CartService.
func (s *service) UpdateItem(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, req UpdateCartItemRequest) (*models.CartItem, error) {
	if userID == uuid.Nil {
		return nil, ErrUserIDRequired
	}
	if itemID == uuid.Nil {
		return nil, ErrItemIDRequired
	}
	if req.Quantity <= 0 {
		return nil, ErrEmptyQuantity
	}

	existItem, err := s.CartRepo.GetItem(ctx, userID, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed get item: %w", err)
	}
	if existItem == nil {
		return nil, ErrProductNotFound
	}

	if req.ProductID != nil {
		existItem.ProductID = *req.ProductID
	}
	existItem.Quantity = req.Quantity

	if err := s.CartRepo.UpdateQuantity(ctx, existItem.ID, existItem.Quantity); err != nil {
		return nil, fmt.Errorf("failed to update quantity: %w", err)
	}

	return existItem, nil
}
