package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	models "github.com/Xiancel/ecommerce/internal/domain"
	repository "github.com/Xiancel/ecommerce/internal/repository/postgres"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func NewService(authRepo repository.UserRepository, jwtSecret string) AuthService {
	return &service{
		userRepo:      authRepo,
		jwtSecret:     jwtSecret,
		tokenDuration: 24 * time.Hour,
	}
}

// Login авторизація користувача
func (s *service) Login(ctx context.Context, req LoginRequset) (*AuthResponse, error) {
	// валідація
	if req.Email == "" {
		return nil, ErrEmailRequired
	}
	if req.Password == "" {
		return nil, ErrPasswordRequired
	}

	// перевірка користувача за Email
	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	// обробка помилок
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// хешування пароля
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// генерація і оновлення токену
	accessToken, err := s.generateToken(user, s.tokenDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate acces token: %w", err)
	}
	refreshToken, err := s.generateToken(user, 7*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// повертає данні користувача
	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// RefreshToken оновлення JWT токену
func (s *service) RefreshToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	// переірка токена на валідність
	claims, err := s.ValidateToken(refreshToken)
	if err != nil {
		return nil, err
	}

	// отримання користувача за ID
	user, err := s.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// генерація нового токену
	newAccessToken, err := s.generateToken(user, 7*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate acces token: %w", err)
	}
	newRefreshToken, err := s.generateToken(user, 7*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// повертає данні користувача
	return &AuthResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
		User:         user,
	}, nil
}

// Register реєстрація користувача
func (s *service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, error) {
	// валідація данних
	if req.Email == "" {
		return nil, ErrEmailRequired
	}
	if req.Password == "" {
		return nil, ErrPasswordRequired
	}
	if len(req.Password) < 8 {
		return nil, ErrWeakPassword
	}

	// перевірка на існування користувача за єлектроною адрессою
	existUser, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	if existUser != nil {
		return nil, ErrUserAlreadyExists
	}

	//hash password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// створення користувача
	user := &models.User{
		Email:        req.Email,
		PasswordHash: string(passwordHash),
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Role:         "customer",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	//Generate tokens
	accessToken, err := s.generateToken(user, s.tokenDuration)
	if err != nil {
		return nil, fmt.Errorf("failed to generate acces token: %w", err)
	}
	refreshToken, err := s.generateToken(user, 7*24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	// повертає данні користувача
	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

// ValidateToken валідація токену
func (s *service) ValidateToken(tokenString string) (*Claims, error) {
	// парсинг токена з claims
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		// перевірка алгоритму підпису
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil {
		return nil, ErrInvalidToken
	}

	// перевірка claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	// перевірка на просроченість токену
	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, ErrTokenExpired
	}
	// повертання claims
	return claims, nil
}

// generateToken генерація нового JWT токену
func (s *service) generateToken(user *models.User, duration time.Duration) (string, error) {
	// створення claims з інформацією про користувача та терміном дії
	claims := &Claims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)), // дата закінчення дії
			IssuedAt:  jwt.NewNumericDate(time.Now()),               // дата створення
			Issuer:    "eccomerce-api",
		},
	}

	// створення новго JWT токену
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// підписання токену
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	// повертає токен
	return tokenString, nil
}
