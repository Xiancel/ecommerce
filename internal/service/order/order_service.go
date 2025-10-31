package order

import (
	"context"
	"fmt"

	models "github.com/Xiancel/ecommerce/internal/domain"
	repository "github.com/Xiancel/ecommerce/internal/repository/postgres"
	"github.com/google/uuid"
)

type service struct {
	orderRepo repository.OrderRepository
}

func NewService(orderRepo repository.OrderRepository) OrderService {
	return &service{orderRepo: orderRepo}
}

// CancelOrder implements OrderService.
func (s *service) CancelOrder(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return fmt.Errorf("order id is required")
	}

	order, err := s.orderRepo.GetById(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	if order.Status == "cancelled" {
		return fmt.Errorf("order already cancelled")
	}
	if order.Status == "delivered" {
		return fmt.Errorf("cannot cancel a delivered order")
	}

	if err := s.orderRepo.UpdateStatus(ctx, id, "cancelled"); err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	return nil
}

// CreateOrder implements OrderService.
func (s *service) CreateOrder(ctx context.Context, req CreateOrderRequset) (*models.Order, error) {
	if req.UserID == uuid.Nil {
		return nil, fmt.Errorf("user id required")
	}
	if req.ShippingAdress.City == "" || req.ShippingAdress.Country == "" ||
		req.ShippingAdress.PostalCode == "" || req.ShippingAdress.Street == "" {
		return nil, fmt.Errorf("shipingadress required")
	}
	if req.PaymentMethod != "Cash" && req.PaymentMethod != "Card" {
		return nil, fmt.Errorf("invalid paymentMethod")
	}

	if len(req.Items) == 0 {
		return nil, fmt.Errorf("order must contain one item")
	}

	for _, item := range req.Items {
		if item.ProductID == uuid.Nil {
			return nil, fmt.Errorf("product id is required")
		}
		if item.Quantity <= 0 {
			return nil, fmt.Errorf("invalid quantity product")
		}
	}

	order := &models.Order{
		UserID:          &req.UserID,
		Status:          "pending",
		PaymentMethod:   req.PaymentMethod,
		ShippingAddress: req.ShippingAdress,
	}

	items := make([]*models.OrderItem, len(req.Items))
	for i := range req.Items {
		items[i] = &req.Items[i]
	}

	if err := s.orderRepo.Create(ctx, order, items); err != nil {
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

// GetOrder implements OrderService.
func (s *service) GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("order id is required")
	}

	order, err := s.orderRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return order, nil
}

// ListOrder implements OrderService.
func (s *service) ListOrder(ctx context.Context, filter OrderFilter) (*OrderListResponse, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 100 {
		filter.Limit = 100
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	var orders []*models.Order
	var err error

	if filter.UserID != nil {
		orders, err = s.orderRepo.ListByUserID(ctx, *filter.UserID, filter.Limit, filter.Offset)
	} else {
		orders, err = s.orderRepo.ListAll(ctx, filter.Limit, filter.Offset)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to list order: %w", err)
	}

	filtered := []*models.Order{}
	for _, o := range orders {
		if filter.Status != "" && o.Status != filter.Status {
			continue
		}
		filtered = append(filtered, o)
	}

	resp := &OrderListResponse{
		Order: filtered,
		Total: len(filtered),
	}

	return resp, nil
}

// UpdateOrderStatus implements OrderService.
func (s *service) UpdateOrderStatus(ctx context.Context, id uuid.UUID, req UpdateOrderRequest) (*models.Order, error) {
	if id == uuid.Nil {
		return nil, fmt.Errorf("order id is required")
	}
	if req.Status == "" {
		return nil, fmt.Errorf("status is required")
	}

	validStatus := map[string]bool{
		"pending":   true,
		"paid":      true,
		"shipped":   true,
		"cancelled": true,
		"delivered": true,
	}
	if !validStatus[req.Status] {
		return nil, fmt.Errorf("invalid status: %s", req.Status)
	}

	order, err := s.orderRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update status: %w", err)
	}

	order.Status = req.Status
	return order, nil
}
