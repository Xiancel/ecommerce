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

// CancelOrder скасування замовлення
func (s *service) CancelOrder(ctx context.Context, id uuid.UUID) error {
	// валідація
	if id == uuid.Nil {
		return ErrOrderIDRequired
	}

	// отримання ID замовлення
	order, err := s.orderRepo.GetById(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get order: %w", err)
	}
	// валідація
	if order.Status == "cancelled" {
		return ErrOrderAlreadyCanceled
	}
	if order.Status == "delivered" {
		return ErrCannotCancelDelivered
	}

	// оновлення статусу(скасування) замовлення
	if err := s.orderRepo.UpdateStatus(ctx, id, "cancelled"); err != nil {
		return fmt.Errorf("failed to cancel order: %w", err)
	}

	return nil
}

// CreateOrder створення замовлення
func (s *service) CreateOrder(ctx context.Context, userID uuid.UUID, req CreateOrderRequest) (*models.Order, error) {
	// валідація
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

	// створення замовлення
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
		// додавання товарув у замовлення
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

	// створення заказу
	if err := s.orderRepo.Create(ctx, order, items); err != nil {
		fmt.Printf("Service lvl CreateOrder error: %+v\n", err)
		return nil, fmt.Errorf("failed to create order: %w", err)
	}

	return order, nil
}

// GetOrder отримання замовлення
func (s *service) GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	// валідація
	if id == uuid.Nil {
		return nil, ErrOrderIDRequired
	}

	// отримання замовлення
	order, err := s.orderRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get order: %w", err)
	}

	return order, nil
}

// ListOrder повертає список замовлень
func (s *service) ListOrder(ctx context.Context, filter OrderFilter) (*OrderListResponse, error) {
	// пагінація
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

	// якщо в фільтрах присутствує userID
	if filter.UserID != nil {
		// виводимо всі замовлення для користувача
		orders, err = s.orderRepo.ListByUserID(ctx, *filter.UserID, filter.Limit, filter.Offset)
	} else {
		// виводимо всі замовлення
		orders, err = s.orderRepo.ListAll(ctx, filter.Limit, filter.Offset)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to list order: %w", err)
	}

	// фільтрація замовлень за статусом
	filtered := []*models.Order{}
	for _, o := range orders {
		if filter.Status != "" && o.Status != filter.Status {
			continue
		}
		filtered = append(filtered, o)
	}

	// формування відповіді
	resp := &OrderListResponse{
		Order: filtered,
		Total: len(filtered),
	}

	return resp, nil
}

// UpdateOrderStatus оновлення статусу замовлення
func (s *service) UpdateOrderStatus(ctx context.Context, id uuid.UUID, req UpdateOrderRequest) (*models.Order, error) {
	// валідація
	if id == uuid.Nil {
		return nil, ErrOrderIDRequired
	}
	if req.Status == "" {
		return nil, ErrStatusRequired
	}

	// створення мапи з валідними статусами
	validStatus := map[string]bool{
		"pending":   true,
		"paid":      true,
		"shipped":   true,
		"canceled":  true,
		"delivered": true,
	}
	// перевірка на валідність статутсу
	if !validStatus[req.Status] {
		return nil, fmt.Errorf("invalid status: %s", req.Status)
	}

	// отримання ID замовлення
	order, err := s.orderRepo.GetById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update status: %w", err)
	}

	// оновлення статусу замовлення
	order.Status = req.Status
	return order, nil
}
