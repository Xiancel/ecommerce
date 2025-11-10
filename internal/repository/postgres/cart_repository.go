package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	database "github.com/Xiancel/ecommerce/internal/db"
	models "github.com/Xiancel/ecommerce/internal/domain"
	"github.com/google/uuid"
)

type CartRepository interface {
	AddItem(ctx context.Context, item *models.CartItem) error
	GetByUserId(ctx context.Context, userID uuid.UUID) ([]*models.CartItemWithProduct, error)
	UpdateQuantity(ctx context.Context, id uuid.UUID, quantity int) error
	RemoveItem(ctx context.Context, userID, id uuid.UUID) error
	Clear(ctx context.Context, userId uuid.UUID) error
	GetItem(ctx context.Context, userId, productId uuid.UUID) (*models.CartItem, error)
	GetItemByID(ctx context.Context, userID, itemID uuid.UUID) (*models.CartItem, error)
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

func (c *CartRepo) GetItemByID(ctx context.Context, userID, itemID uuid.UUID) (*models.CartItem, error) {
	var item models.CartItem
	query := `
	SELECT id, user_id, product_id, quantity, created_at
	FROM cart_items
	WHERE id = $1 AND user_id = $2
	`
	err := c.db.GetContext(ctx, &item, query, itemID, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get item by id: %w", err)
	}
	return &item, nil
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get item: %w", err)
	}
	return &item, nil
}

func (c *CartRepo) RemoveItem(ctx context.Context, userID, id uuid.UUID) error {
	query := `
	DELETE FROM cart_items
	WHERE id = $1 AND user_id = $2
	`
	res, err := c.db.ExecContext(ctx, query, id, userID)
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
