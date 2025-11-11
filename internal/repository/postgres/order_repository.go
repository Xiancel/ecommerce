package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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

	shippingJSON, errjson := json.Marshal(order.ShippingAddress)
	if errjson != nil {
		return fmt.Errorf("failed to marshal shipping address: %w", errjson)
	}

	orderQuery := `
	INSERT INTO orders (id, user_id, status, total_amount, shipping_address, payment_method, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5,$6,NOW(),NOW())
	`

	_, err := o.db.ExecContext(ctx, orderQuery,
		order.ID,
		order.UserID,
		order.Status,
		order.TotalAmount,
		shippingJSON,
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
	var shippingBytes []byte

	query := `
	SELECT id, user_id, status, total_amount, shipping_address, payment_method, created_at, updated_at
	FROM orders
	WHERE id = $1
	`
	temp := struct {
		ID            uuid.UUID  `db:"id"`
		UserID        *uuid.UUID `db:"user_id"`
		Status        string     `db:"status"`
		TotalAmount   float64    `db:"total_amount"`
		Shipping      []byte     `db:"shipping_address"`
		PaymentMethod string     `db:"payment_method"`
		CreatedAt     time.Time  `db:"created_at"`
		UpdatedAt     time.Time  `db:"updated_at"`
	}{}
	err := o.db.GetContext(ctx, &temp, query, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order id: %w", err)
	}

	order.ID = temp.ID
	order.UserID = temp.UserID
	order.Status = temp.Status
	order.TotalAmount = temp.TotalAmount
	order.PaymentMethod = temp.PaymentMethod
	order.CreatedAt = temp.CreatedAt
	order.UpdatedAt = temp.UpdatedAt
	shippingBytes = temp.Shipping

	if err := json.Unmarshal(shippingBytes, &order.ShippingAddress); err != nil {
		return nil, fmt.Errorf("failed to unmarshal shipping address: %w", err)
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
	// временна структура для роботи з shipping adress
	var tempOrders []struct {
		ID            uuid.UUID  `db:"id"`
		UserID        *uuid.UUID `db:"user_id"`
		Status        string     `db:"status"`
		TotalAmount   float64    `db:"total_amount"`
		Shipping      []byte     `db:"shipping_address"`
		PaymentMethod string     `db:"payment_method"`
		CreatedAt     time.Time  `db:"created_at"`
		UpdatedAt     time.Time  `db:"updated_at"`
	}

	err := o.db.SelectContext(ctx, &tempOrders, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	orders := make([]*models.Order, len(tempOrders))
	for i, temp := range tempOrders {
		order := &models.Order{
			ID:            temp.ID,
			UserID:        temp.UserID,
			Status:        temp.Status,
			TotalAmount:   temp.TotalAmount,
			PaymentMethod: temp.PaymentMethod,
			CreatedAt:     temp.CreatedAt,
			UpdatedAt:     temp.UpdatedAt,
		}
		if err := json.Unmarshal(temp.Shipping, &order.ShippingAddress); err != nil {
			return nil, fmt.Errorf("failed to unmarshal shipping address: %w", err)
		}
		orders[i] = order
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

	var tempOrders []struct {
		ID            uuid.UUID  `db:"id"`
		UserID        *uuid.UUID `db:"user_id"`
		Status        string     `db:"status"`
		TotalAmount   float64    `db:"total_amount"`
		Shipping      []byte     `db:"shipping_address"`
		PaymentMethod string     `db:"payment_method"`
		CreatedAt     time.Time  `db:"created_at"`
		UpdatedAt     time.Time  `db:"updated_at"`
	}

	err := o.db.SelectContext(ctx, &tempOrders, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list orders: %w", err)
	}

	orders := make([]*models.Order, len(tempOrders))
	for i, temp := range tempOrders {
		order := &models.Order{
			ID:            temp.ID,
			UserID:        temp.UserID,
			Status:        temp.Status,
			TotalAmount:   temp.TotalAmount,
			PaymentMethod: temp.PaymentMethod,
			CreatedAt:     temp.CreatedAt,
			UpdatedAt:     temp.UpdatedAt,
		}
		if err := json.Unmarshal(temp.Shipping, &order.ShippingAddress); err != nil {
			return nil, fmt.Errorf("failed to unmarshal shipping address: %w", err)
		}
		orders[i] = order
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
