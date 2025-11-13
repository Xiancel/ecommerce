package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Xiancel/ecommerce/internal/db"
	"github.com/joho/godotenv"

	httpHandler "github.com/Xiancel/ecommerce/internal/handler/http"
	postgres "github.com/Xiancel/ecommerce/internal/repository/postgres"
	authService "github.com/Xiancel/ecommerce/internal/service/auth"
	cartService "github.com/Xiancel/ecommerce/internal/service/cart"
	orderService "github.com/Xiancel/ecommerce/internal/service/order"
	productService "github.com/Xiancel/ecommerce/internal/service/product"
	userService "github.com/Xiancel/ecommerce/internal/service/user"
)

// @title E-Commerce API
// @version 1.0
// @description API для e-commerce платформи з управлінням продуктами, кошиком та замовленнями
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	// завантаження .env файлів
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment varibles")
	}

	// отримання конфігурації з env
	dbHost := getEnv("DB_HOST", "db")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "user")
	dbPassword := getEnv("DB_PASSWORD", "1234!")
	dbName := getEnv("DB_NAME", "ecommerce_db")
	dbSSLMode := getEnv("DB_SSLMODE", "disable")
	serverPort := getEnv("APP_PORT", "8080")
	jwtSecret := getEnv("JWT_SECRET", "kfJ+JpWThVtZ5p0hIM9s7jFGucNvHdn59aTfzT7fQ2iqlt3rH2bnSKTwsm4B3Q3P")

	// конфігурація бази данних
	dbConfig := db.Config{
		Host:     dbHost,
		Port:     dbPort,
		User:     dbUser,
		Password: dbPassword,
		DBName:   dbName,
		SSLMode:  dbSSLMode,
	}

	// ініціалізація бази данних
	database, err := db.NewDB(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	log.Println("✅ Database connetion established")

	// ініціалізація репозеторіїв
	productRepo := postgres.NewProductRepository(database)
	userRepo := postgres.NewUserRepository(database)
	cartRepo := postgres.NewCartRepository(database)
	orderRepo := postgres.NewOrderRepository(database)

	log.Println("✅ Repository initialized")

	// ініціалізація сервісів
	productSrv := productService.NewService(productRepo)
	userSrv := userService.NewService(userRepo)
	authSrv := authService.NewService(userRepo, jwtSecret)
	cartSrv := cartService.NewService(cartRepo)
	orderService := orderService.NewService(orderRepo, productRepo)

	log.Println("✅ Services initialized")

	// ініціалізація http router
	router := httpHandler.NewRouter(httpHandler.RouterConfig{
		AuthService:    authSrv,
		ProductService: productSrv,
		CartService:    cartSrv,
		OrderService:   orderService,
		UserService:    userSrv,
	})

	log.Println("✅ HTTP router initialized")

	// створення HTTP серверу
	server := &http.Server{
		Addr:         ":" + serverPort,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// запуск сервера в горутинах
	go func() {
		log.Printf("🚀 Server starting on http://localhost:%s", serverPort)
		log.Printf("📚 API documentation: http://localhost:%s/api/v1", serverPort)
		log.Printf("🏥 Health check: http://localhost:%s/health", serverPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// повернення повідомлень
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("⚠️ Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server stopped gracefully")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
