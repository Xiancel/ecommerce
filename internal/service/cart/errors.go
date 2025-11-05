package cart

import "errors"

var (
	//Validate errors
	ErrUserIDRequired = errors.New("user id is required")
	ErrItemIDRequired = errors.New("item id is required")
	ErrEmptyQuantity  = errors.New("quantity must be greater than 0")

	// Product error
	ErrProductNotFound = errors.New("product not found")

	// Cart item error
	ErrItemNotFound = errors.New("cart item not found")
)
