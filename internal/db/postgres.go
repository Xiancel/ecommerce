package db

import (
	"context"
	"fmt"
	"time"

	models "github.com/Xiancel/ecommerce/internal/domain"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	*sqlx.DB
}

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// Create implements repository.ProductRepository.
func (c Config) Create(ctx context.Context, product *models.Product) error {
	panic("unimplemented")
}

// Delete implements repository.ProductRepository.
func (c Config) Delete(ctx context.Context, id uuid.UUID) error {
	panic("unimplemented")
}

// GetById implements repository.ProductRepository.
func (c Config) GetById(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	panic("unimplemented")
}

// List implements repository.ProductRepository.
func (c Config) List(ctx context.Context, filter models.ListFilter) ([]*models.Product, error) {
	panic("unimplemented")
}

// Update implements repository.ProductRepository.
func (c Config) Update(ctx context.Context, product *models.Product) error {
	panic("unimplemented")
}

// UpdateStock implements repository.ProductRepository.
func (c Config) UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error {
	panic("unimplemented")
}

func NewDB(ctg Config) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		ctg.Host, ctg.Port, ctg.User, ctg.Password, ctg.DBName, ctg.SSLMode)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}
