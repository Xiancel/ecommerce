package order

import "errors"

var (
	//Order validate errors
	ErrOrderIDRequired         = errors.New("order id is required")
	ErrUserIDRequired          = errors.New("user id is required")
	ErrShippingAddressRequired = errors.New("shipping address is required")
	ErrPaymentMethodInvalid    = errors.New("invalid payment method")
	ErrStatusRequired          = errors.New("status is required")

	//Order item errors
	ErrOrderMustContainItem   = errors.New("order must contain at least one item")
	ErrProductIDRequired      = errors.New("product id is required")
	ErrInvalidProductQuantity = errors.New("product quantity must be greater than 0")

	//Order status errors
	ErrOrderAlreadyCanceled  = errors.New("order already canceled")
	ErrCannotCancelDelivered = errors.New("cannot cancel a delivered order")
)
