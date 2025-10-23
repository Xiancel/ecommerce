package repository

import (
	"context"
	"fmt"

	database "github.com/Xiancel/ecommerce/internal/db"
	models "github.com/Xiancel/ecommerce/internal/domain"
	"github.com/google/uuid"
)

type OrderRepository interface {
	Create(ctx context.Context, order *models.Order, items []*models.OrderItem) error
	GetById(ctx context.Context, id uuid.UUID) (*models.Order, error)
	GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]*models.OrderItem, error)
	ListByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*models.Order, error)
	ListAll(ctx context.Context, limit, offset int) ([]*models.Order, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}

type orderRepo struct {
	db *database.DB
}

func NewOrderRepository(db *database.DB) OrderRepository {
	return &orderRepo{db: db}
}

// Create implements OrderRepository.
func (o *orderRepo) Create(ctx context.Context, order *models.Order, items []*models.OrderItem) error {
	orderQuery := `
	INSERT INTO orders (id, user_id, status, total_amount, shipping_address, payment_method, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5,$6,NOW(),NOW())
	`

	_, err := o.db.ExecContext(ctx, orderQuery,
		order.ID,
		order.UserID,
		order.Status,
		order.TotalAmount,
		order.ShippingAddress,
		order.PaymentMethod,
	)
	if err != nil {
		return fmt.Errorf("failed to create order: %w", err)
	}

	itemQuery := `
		INSERT INTO order_items (id, order_id, product_id, quantity, price, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
	`

	for _, item := range items {
		_, err := o.db.ExecContext(ctx, itemQuery,
			item.ID,
			item.OrderID,
			item.ProductID,
			item.Quantity,
			item.Price,
		)
		if err != nil {
			return fmt.Errorf("failed to create order items: %w", err)
		}
	}
	return nil
}

// GetById implements OrderRepository.
func (o *orderRepo) GetById(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	var order models.Order
	query := `
	SELECT id, user_id, status, total_amount, shipping_address, payment_method, created_at, updated_at
	FROM orders
	WHERE id = $1
	`

	err := o.db.GetContext(ctx, &order, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order id: %w", err)
	}

	return &order, nil
}

// GetOrderItems implements OrderRepository.
func (o *orderRepo) GetOrderItems(ctx context.Context, orderID uuid.UUID) ([]*models.OrderItem, error) {
	var items []*models.OrderItem

	query := `
	SELECT id, order_id, product_id, quantity, price, created_at
	FROM order_items
	WHERE order_id = $1
	`

	err := o.db.SelectContext(ctx, &items, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order items: %w", err)
	}

	return items, nil
}

// ListAll implements OrderRepository.
func (o *orderRepo) ListAll(ctx context.Context, limit int, offset int) ([]*models.Order, error) {
	query := `
	SELECT id, user_id, status, total_amount, shipping_address, payment_method, created_at, updated_at
	FROM orders
	ORDER BY created_at DESC
	LIMIT $1 OFFSET $2
	`

	var orders []*models.Order

	err := o.db.SelectContext(ctx, &orders, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	return orders, nil
}

// ListByUserID implements OrderRepository.
func (o *orderRepo) ListByUserID(ctx context.Context, userID uuid.UUID, limit int, offset int) ([]*models.Order, error) {
	query := `
	SELECT id, user_id, status, total_amount, shipping_address, payment_method, created_at, updated_at
	FROM orders
	WHERE user_id = $1
	ORDER BY created_at DESC
	LIMIT $2 OFFSET $3
	`

	var orders []*models.Order

	err := o.db.SelectContext(ctx, &orders, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}
	return orders, nil
}

// UpdateStatus implements OrderRepository.
func (o *orderRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	query := `
	UPDATE orders
	SET status = $1,
		updated_at = NOW()
	WHERE id = $2
	`

	res, err := o.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update orders: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("order not found")
	}

	return nil
}
