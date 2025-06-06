package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/product"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository interface {
	GetAll(ctx context.Context) ([]models.Product, error)
	GetById(ctx context.Context, id int64) (models.Product, error)
	GetByName(ctx context.Context, name string) ([]models.Product, error)
	GetByManufacturer(ctx context.Context, manufacturer string) ([]models.Product, error)
	Create(ctx context.Context, product models.Product) (models.Product, error)
	Update(ctx context.Context, product models.Product) (models.Product, error)
	DeleteById(ctx context.Context, id int64) error
	GetByCostPriceRange(ctx context.Context, min, max float64) ([]models.Product, error)
	GetBySalePriceRange(ctx context.Context, min, max float64) ([]models.Product, error)
	GetLowInStock(ctx context.Context, threshold int) ([]models.Product, error)
}

type productRepository struct {
	db *pgxpool.Pool
}

func NewProductRepository(db *pgxpool.Pool) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) GetAll(ctx context.Context) ([]models.Product, error) {
	query := `SELECT 
		id, product_name, manufacturer, product_description, 
		cost_price, sale_price, stock_quantity, barcode, 
		created_at, updated_at 
	FROM products`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar produtos: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(
			&product.ID,
			&product.ProductName,
			&product.Manufacturer,
			&product.Description,
			&product.CostPrice,
			&product.SalePrice,
			&product.StockQuantity,
			&product.Barcode,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler os dados do produto: %w", err)
		}
		products = append(products, product)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("erro ao iterar sobre os resultados: %w", rows.Err())
	}

	return products, nil
}

func (r *productRepository) GetById(ctx context.Context, id int64) (models.Product, error) {
	var product models.Product
	query := `SELECT 
		id, product_name, manufacturer, product_description, 
		cost_price, sale_price, stock_quantity, barcode, 
		created_at, updated_at 
	FROM products WHERE id = $1`

	err := r.db.QueryRow(ctx, query, id).Scan(
		&product.ID,
		&product.ProductName,
		&product.Manufacturer,
		&product.Description,
		&product.CostPrice,
		&product.SalePrice,
		&product.StockQuantity,
		&product.Barcode,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return product, fmt.Errorf("produto não encontrado")
		}
		return product, fmt.Errorf("erro ao buscar produto: %w", err)
	}

	return product, nil
}

func (r *productRepository) GetByName(ctx context.Context, name string) ([]models.Product, error) {
	query := `SELECT 
		id, product_name, manufacturer, product_description, 
		cost_price, sale_price, stock_quantity, barcode, 
		created_at, updated_at 
	FROM products WHERE product_name LIKE '%' || $1 || '%'`

	rows, err := r.db.Query(ctx, query, name)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar produtos por nome: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(
			&product.ID,
			&product.ProductName,
			&product.Manufacturer,
			&product.Description,
			&product.CostPrice,
			&product.SalePrice,
			&product.StockQuantity,
			&product.Barcode,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler os dados do produto: %w", err)
		}
		products = append(products, product)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("erro ao iterar sobre os resultados: %w", rows.Err())
	}

	return products, nil
}

func (r *productRepository) GetByManufacturer(ctx context.Context, manufacturer string) ([]models.Product, error) {
	query := `SELECT 
		id, product_name, manufacturer, product_description, 
		cost_price, sale_price, stock_quantity, barcode, 
		created_at, updated_at 
	FROM products WHERE manufacturer = $1`

	rows, err := r.db.Query(ctx, query, manufacturer)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar produtos por fabricante: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(
			&product.ID,
			&product.ProductName,
			&product.Manufacturer,
			&product.Description,
			&product.CostPrice,
			&product.SalePrice,
			&product.StockQuantity,
			&product.Barcode,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler os dados do produto: %w", err)
		}
		products = append(products, product)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("erro ao iterar sobre os resultados: %w", rows.Err())
	}

	return products, nil
}

func (r *productRepository) Create(ctx context.Context, product models.Product) (models.Product, error) {
	query := `INSERT INTO products (
		product_name, manufacturer, product_description, 
		cost_price, sale_price, stock_quantity, barcode,
		created_at, updated_at
	) VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
	RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		product.ProductName,
		product.Manufacturer,
		product.Description,
		product.CostPrice,
		product.SalePrice,
		product.StockQuantity,
		product.Barcode,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		return models.Product{}, fmt.Errorf("erro ao criar produto: %w", err)
	}

	return product, nil
}

func (r *productRepository) Update(ctx context.Context, product models.Product) (models.Product, error) {
	const query = `
		UPDATE products
		SET
			product_name        = $1,
			manufacturer        = $2,
			product_description = $3,
			cost_price          = $4,
			sale_price          = $5,
			stock_quantity      = $6,
			barcode             = $7,
			updated_at          = NOW()
		WHERE 
			id = $8
		RETURNING
			id,
			product_name,
			manufacturer,
			product_description,
			cost_price,
			sale_price,
			stock_quantity,
			barcode,
			created_at,
			updated_at;
	`

	var (
		id            int64
		productName   string
		manufacturer  string
		description   string
		costPrice     float64
		salePrice     float64
		stockQuantity int
		barcode       string
		createdAt     time.Time
		updatedAt     time.Time
	)

	err := r.db.QueryRow(ctx, query,
		product.ProductName,
		product.Manufacturer,
		product.Description,
		product.CostPrice,
		product.SalePrice,
		product.StockQuantity,
		product.Barcode,
		product.ID,
	).Scan(
		&id,
		&productName,
		&manufacturer,
		&description,
		&costPrice,
		&salePrice,
		&stockQuantity,
		&barcode,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Product{}, fmt.Errorf("produto não encontrado")
		}
		return models.Product{}, fmt.Errorf("erro ao atualizar produto: %w", err)
	}

	product.ID = id
	product.ProductName = productName
	product.Manufacturer = manufacturer
	product.Description = description
	product.CostPrice = costPrice
	product.SalePrice = salePrice
	product.StockQuantity = stockQuantity
	product.Barcode = barcode
	product.CreatedAt = createdAt
	product.UpdatedAt = updatedAt

	return product, nil
}

func (r *productRepository) DeleteById(ctx context.Context, id int64) error {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar produto: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("produto não encontrado")
	}

	return nil
}

func (r *productRepository) GetByCostPriceRange(ctx context.Context, min, max float64) ([]models.Product, error) {
	query := `SELECT 
		id, product_name, manufacturer, product_description, 
		cost_price, sale_price, stock_quantity, barcode, 
		created_at, updated_at 
	FROM products WHERE cost_price BETWEEN $1 AND $2`

	return r.getProductsByPriceRange(ctx, query, min, max)
}

func (r *productRepository) GetBySalePriceRange(ctx context.Context, min, max float64) ([]models.Product, error) {
	query := `SELECT 
		id, product_name, manufacturer, product_description, 
		cost_price, sale_price, stock_quantity, barcode, 
		created_at, updated_at 
	FROM products WHERE sale_price BETWEEN $1 AND $2`

	return r.getProductsByPriceRange(ctx, query, min, max)
}

func (r *productRepository) getProductsByPriceRange(ctx context.Context, query string, min, max float64) ([]models.Product, error) {
	rows, err := r.db.Query(ctx, query, min, max)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar produtos por faixa de preço: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(
			&product.ID,
			&product.ProductName,
			&product.Manufacturer,
			&product.Description,
			&product.CostPrice,
			&product.SalePrice,
			&product.StockQuantity,
			&product.Barcode,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler os dados do produto: %w", err)
		}
		products = append(products, product)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("erro ao iterar sobre os resultados: %w", rows.Err())
	}

	return products, nil
}

func (r *productRepository) GetLowInStock(ctx context.Context, threshold int) ([]models.Product, error) {
	query := `SELECT 
		id, product_name, manufacturer, product_description, 
		cost_price, sale_price, stock_quantity, barcode, 
		created_at, updated_at 
	FROM products WHERE stock_quantity <= $1`

	rows, err := r.db.Query(ctx, query, threshold)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar produtos com estoque baixo: %w", err)
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(
			&product.ID,
			&product.ProductName,
			&product.Manufacturer,
			&product.Description,
			&product.CostPrice,
			&product.SalePrice,
			&product.StockQuantity,
			&product.Barcode,
			&product.CreatedAt,
			&product.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("erro ao ler os dados do produto: %w", err)
		}
		products = append(products, product)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("erro ao iterar sobre os resultados: %w", rows.Err())
	}

	return products, nil
}
