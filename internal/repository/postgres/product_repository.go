package repository

import (
	"context"
	"fmt"
	"strings"

	database "github.com/Xiancel/ecommerce/internal/db"
	models "github.com/Xiancel/ecommerce/internal/domain"
	"github.com/google/uuid"
)

// ProductRepository інтерфейс для роботи з продуктами
type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) error
	GetById(ctx context.Context, id uuid.UUID) (*models.Product, error)
	List(ctx context.Context, filter models.ListFilter) ([]*models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type productRepo struct {
	db *database.DB
}

func NewProductRepository(db *database.DB) ProductRepository {
	return &productRepo{db: db}
}

// Create створює новий продукт
func (p *productRepo) Create(ctx context.Context, product *models.Product) error {
	query := `
	INSERT INTO products (id, name, description, price, stock, category_id, image_url, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	`

	// присвоєння айді продукту
	product.ID = uuid.New()
	// створення продукту
	_, err := p.db.ExecContext(ctx, query,
		product.ID,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.CategoryID,
		product.ImageURL,
	)
	// обробка помилки
	if err != nil {
		return fmt.Errorf("failed to created product: %w", err)
	}

	return nil
}

// Delete видаляє продукт
func (p *productRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
	DELETE FROM products WHERE id = $1
	`
	// видалення продукту за його ID
	res, err := p.db.ExecContext(ctx, query, id)
	// обробка помилок
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

// GetById повертає продукт за його ID
func (p *productRepo) GetById(ctx context.Context, id uuid.UUID) (*models.Product, error) {
	var product models.Product
	query := `
	SELECT id, name, description, price, stock, category_id, image_url, created_at, updated_at
	FROM products
	WHERE id = $1
	`

	// передання інформації про продукт за його ID
	err := p.db.GetContext(ctx, &product, query, id)
	// обробка помилок
	if err != nil {
		return nil, fmt.Errorf("failed to get product: %w", err)
	}

	return &product, nil
}

// List повертає список продуктів
func (p *productRepo) List(ctx context.Context, filter models.ListFilter) ([]*models.Product, error) {
	query := `
	SELECT id, name, description, price, stock, category_id, image_url, created_at, updated_at
	FROM products
	WHERE 1=1
	`

	args := []interface{}{}
	argsCount := 1

	// фільтрація
	if filter.CategoryID != nil && *filter.CategoryID != uuid.Nil {
		query += fmt.Sprintf(" AND category_id = $%d", argsCount)
		args = append(args, *filter.CategoryID)
		argsCount++
	}

	if filter.MinPrice != nil {
		query += fmt.Sprintf(" AND price >= $%d", argsCount)
		args = append(args, *filter.MinPrice)
		argsCount++
	}

	if filter.MaxPrice != nil {
		query += fmt.Sprintf(" AND price <= $%d", argsCount)
		args = append(args, *filter.MaxPrice)
		argsCount++
	}

	if filter.Search != "" {
		query += fmt.Sprintf(" AND (name ILIKE $%d OR description ILIKE $%d) ", argsCount, argsCount)
		args = append(args, "%"+filter.Search+"%")
		argsCount++
	}

	orderBy := "created_at DESC"
	if filter.OrderBy != "" {
		allowedOrders := map[string]bool{
			"price_asc":  true,
			"price_desc": true,
			"name_asc":   true,
			"name_desc":  true,
		}
		if allowedOrders[filter.OrderBy] {
			parts := strings.Split(filter.OrderBy, "_")
			orderBy = parts[0] + " " + strings.ToUpper(parts[1])
		}
	}
	query += " ORDER BY " + orderBy

	// Pagination
	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argsCount)
		args = append(args, filter.Limit)
		argsCount++
	}

	if filter.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argsCount)
		args = append(args, filter.Offset)
		argsCount++
	}

	// отримання продуктів
	var products []*models.Product
	err := p.db.SelectContext(ctx, &products, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list products: %w", err)
	}
	return products, nil
}

// Update оновлення продукту
func (p *productRepo) Update(ctx context.Context, product *models.Product) error {
	query := `
	UPDATE products 
	SET name = $1,
		description = $2,
		price = $3,
		stock = $4,
		category_id = $5, 
		image_url = $6,  
		updated_at = NOW()
	WHERE id = $7
	`

	// оновлення даних продукта
	res, err := p.db.ExecContext(ctx, query,
		product.Name,
		product.Description,
		product.Price,
		product.Stock,
		product.CategoryID,
		product.ImageURL,
		product.ID,
	)
	// обробка помилок
	if err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

// UpdateStock оновлення количества продукту
func (p *productRepo) UpdateStock(ctx context.Context, id uuid.UUID, quantity int) error {
	query := `
	UPDATE products
	SET stock = $1
		updated_at = NOW()
	WHERE id = $2
	`
	// оновлення количество продукта по ID
	res, err := p.db.ExecContext(ctx, query, quantity, id)
	// обробка помилок
	if err != nil {
		return fmt.Errorf("failed to update stock: %w", err)
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}
