package repository

import (
	"context"
	"fmt"

	database "github.com/Xiancel/ecommerce/internal/db"
	models "github.com/Xiancel/ecommerce/internal/domain"
	"github.com/google/uuid"
)

type CartRepository interface {
	AddItem(ctx context.Context, item *models.CartItem) error
	GetByUserId(ctx context.Context, userID uuid.UUID) ([]*models.CartItemWithProduct, error)
	UpdateQuantity(ctx context.Context, id uuid.UUID, quantity int) error
	RemoveItem(ctx context.Context, id uuid.UUID) error
	Clear(ctx context.Context, userId uuid.UUID) error
	GetItem(ctx context.Context, userId, productId uuid.UUID) (*models.CartItem, error)
}

type CartRepo struct {
	db *database.DB
}

func NewCartRepository(db *database.DB) CartRepository {
	return &CartRepo{db: db}
}

// AddItem implements CartRepository.
func (c *CartRepo) AddItem(ctx context.Context, item *models.CartItem) error {
	query := `
	INSERT INTO cart_items (id, user_id, product_id, quantity, created_at)
	VALUES ($1, $2, $3, $4, NOW())
	`

	_, err := c.db.ExecContext(ctx, query,
		item.ID,
		item.UserID,
		item.ProductID,
		item.Quantity,
		item.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to add item: %w", err)
	}

	return nil
}

// Clear implements CartRepository.
func (c *CartRepo) Clear(ctx context.Context, userId uuid.UUID) error {
	query := `
	DELETE FROM cart_items WHERE user_id = $1
	`

	_, err := c.db.ExecContext(ctx, query, userId)
	if err != nil {
		return fmt.Errorf("failed to clear cart: %w", err)
	}

	return nil
}

// GetByUserId implements CartRepository.
func (c *CartRepo) GetByUserId(ctx context.Context, userID uuid.UUID) ([]*models.CartItemWithProduct, error) {
	var items []*models.CartItemWithProduct

	query := `
	SELECT
		ci.id,
		ci.user_id,
		ci.product_id,
		ci.quantity,
		ci.created_at,
		p.name AS product_name,
		p.price AS product_price,
		p.stock AS product_stock
	FROM cart_items ci
	JOIN products p ON ci.product_id = p.id
	WHERE ci.user_id = $1
	`

	err := c.db.SelectContext(ctx, &items, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cart items: %w", err)
	}
	return items, nil
}

// GetItem implements CartRepository.
func (c *CartRepo) GetItem(ctx context.Context, userId uuid.UUID, productId uuid.UUID) (*models.CartItem, error) {
	var item models.CartItem
	query := `
	SELECT id, user_id, product_id, quantity, created_at
	FROM cart_items
	WHERE user_id = $1 AND product_id = $2
	`

	err := c.db.GetContext(ctx, &item, query, userId, productId)
	if err != nil {
		return nil, fmt.Errorf("failed to get item: %w", err)
	}
	return &item, nil
}

// RemoveItem implements CartRepository.
func (c *CartRepo) RemoveItem(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE FROM cart_items
	WHERE id = $1
	`
	res, err := c.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to remove item: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("item not found")
	}
	return nil
}

// UpdateQuantity implements CartRepository.
func (c *CartRepo) UpdateQuantity(ctx context.Context, id uuid.UUID, quantity int) error {
	query := `
	UPDATE cart_items
	SET quantity = $1
	WHERE id = $2
	`

	res, err := c.db.ExecContext(ctx, query, quantity, id)
	if err != nil {
		return fmt.Errorf("failed to update quantity: %w", err)
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("cart item not found")
	}

	return nil
}
