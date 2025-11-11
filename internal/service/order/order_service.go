package order

import (
	"context"
	"fmt"
	"time"

	models "github.com/Xiancel/ecommerce/internal/domain"
	repository "github.com/Xiancel/ecommerce/internal/repository/postgres"
	"github.com/google/uuid"
)

type service struct {
	orderRepo   repository.OrderRepository
	productRepo repository.ProductRepository
}

func NewService(orderRepo repository.OrderRepository, productRepo repository.ProductRepository) OrderService {
	return &service{orderRepo: orderRepo,
		productRepo: productRepo}
}

// CancelOrder implements OrderService.
func (s *service) CancelOrder(ctx context.Context, id uuid.UUID) error {
	if id == uuid.Nil {
		return ErrOrderIDRequired
	}

	order, err := s.orderRepo.GetById(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	if order.Status == "cancelled" {
		return ErrOrderAlreadyCanceled
	}
	if order.Status == "delivered" {
		return ErrCannotCancelDelivered
	}

	if err := s.orderRepo.UpdateStatus(ctx, id, "cancelled"); err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	return nil
}

// CreateOrder implements OrderService.
func (s *service) CreateOrder(ctx context.Context, userID uuid.UUID, req CreateOrderRequest) (*models.Order, error) {
	if userID == uuid.Nil {
		return nil, ErrUserIDRequired
	}
	if req.ShippingAdress.City == "" || req.ShippingAdress.Country == "" ||
		req.ShippingAdress.PostalCode == "" || req.ShippingAdress.Street == "" {
		return nil, ErrShippingAddressRequired
	}
	if req.PaymentMethod != "cash" && req.PaymentMethod != "card" {
		return nil, ErrPaymentMethodInvalid
	}

	if len(req.Items) == 0 {
		return nil, ErrOrderMustContainItem
	}

	order := &models.Order{
		ID:              uuid.New(),
		UserID:          &userID,
		Status:          "pending",
		ShippingAddress: req.ShippingAdress,
		PaymentMethod:   req.PaymentMethod,
	}

	var total float64
	items := make([]*models.OrderItem, len(req.Items))

	for i, item := range req.Items {
		product, err := s.productRepo.GetById(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product not found: %w", err)
		}

		items[i] = &models.OrderItem{
			ID:        uuid.New(),
			OrderID:   order.ID,
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			Price:     product.Price,
			CreatedAt: time.Now(),
		}
		total += product.Price * float64(item.Quantity)
	}
	order.TotalAmount = total

	if err := s.orderRepo.Create(ctx, order, items); err != nil {
		fmt.Printf("Service lvl CreateOrder error: %+v\n", err)
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

// GetOrder implements OrderService.
func (s *service) GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	if id == uuid.Nil {
		return nil, ErrOrderIDRequired
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
		return nil, ErrOrderIDRequired
	}
	if req.Status == "" {
		return nil, ErrStatusRequired
	}

	validStatus := map[string]bool{
		"pending":   true,
		"paid":      true,
		"shipped":   true,
		"canceled":  true,
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
