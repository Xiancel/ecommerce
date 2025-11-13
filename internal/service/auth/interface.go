package auth

import "context"

// AuthService Інтерфейс для роботи з Антефікацією та авторизацією
type AuthService interface {
	Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error)
	Login(ctx context.Context, req LoginRequset) (*AuthResponse, error)
	RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error)
	ValidateToken(tokenString string) (*Claims, error)
}
