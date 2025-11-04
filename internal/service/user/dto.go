package user

import models "github.com/Xiancel/ecommerce/internal/domain"

type UpdateUserRequest struct {
	Email     *string `json:"email" validate:"omitempty,email"`
	Password  *string `json:"password" validate:"omitempty,min=6,max=100"`
	FirstName *string `json:"first_name" validate:"omitempty,min=2,max=100"`
	LastName  *string `json:"last_name" validate:"omitempty,min=2,max=100"`
	Role      *string `json:"role" validate:"omitempty,oneof=user admin"`
}

type UserFilter struct {
	Search  string  `json:"search"`
	Role    *string `json:"role" validate:"omitempty,oneof=user admin"`
	Limit   int     `json:"limit" validate:"required,min=1,max=100"`
	Offset  int     `json:"offset" validate:"gte=0"`
	OrderBy string  `json:"order_by" validate:"omitempty,oneof=created_at_asc created_at_desc name_asc name_desc"`
}
type UserListResponse struct {
	Users  []*models.User `json:"users"`
	Total  int            `json:"total"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}
