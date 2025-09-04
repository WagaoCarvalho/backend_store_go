package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"

	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) (*models.Product, error)
	GetAll(ctx context.Context, limit, offset int) ([]*models.Product, error)
	GetByID(ctx context.Context, id int64) (*models.Product, error)
	GetByName(ctx context.Context, name string) ([]*models.Product, error)
	GetByManufacturer(ctx context.Context, manufacturer string) ([]*models.Product, error)
	GetVersionByID(ctx context.Context, id int64) (int64, error)
	Update(ctx context.Context, product *models.Product) (*models.Product, error)
	Delete(ctx context.Context, id int64) error

	EnableProduct(ctx context.Context, uid int64) error
	DisableProduct(ctx context.Context, uid int64) error

	UpdateStock(ctx context.Context, id int64, quantity int) error
	IncreaseStock(ctx context.Context, id int64, amount int) error
	DecreaseStock(ctx context.Context, id int64, amount int) error
	GetStock(ctx context.Context, id int64) (int, error)

	EnableDiscount(ctx context.Context, id int64) error
	DisableDiscount(ctx context.Context, id int64) error
	ApplyDiscount(ctx context.Context, id int64, percent float64) (*models.Product, error)
}

type productRepository struct {
	db     *pgxpool.Pool
	logger logger.LogAdapterInterface
}

func NewProductRepository(db *pgxpool.Pool, logger logger.LogAdapterInterface) ProductRepository {
	return &productRepository{db: db, logger: logger}
}

func (r *productRepository) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	const query = `
		INSERT INTO products (
			supplier_id, product_name, manufacturer,
			product_description, cost_price, sale_price,
			stock_quantity, barcode, status,
			allow_discount, max_discount_percent,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
		RETURNING id, created_at, updated_at;
	`

	err := r.db.QueryRow(ctx, query,
		product.SupplierID,
		product.ProductName,
		product.Manufacturer,
		product.Description,
		product.CostPrice,
		product.SalePrice,
		product.StockQuantity,
		product.Barcode,
		product.Status,
		product.AllowDiscount,
		product.MaxDiscountPercent,
	).Scan(&product.ID, &product.CreatedAt, &product.UpdatedAt)

	if err != nil {
		if errMsgPg.IsForeignKeyViolation(err) {
			return nil, errMsg.ErrInvalidForeignKey
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return product, nil
}

func (r *productRepository) GetAll(ctx context.Context, limit, offset int) ([]*models.Product, error) {
	const query = `
	SELECT id, supplier_id, product_name, manufacturer, product_description,
		cost_price, sale_price, stock_quantity, barcode,
		status, allow_discount, max_discount_percent,
		created_at, updated_at
	FROM products
	ORDER BY id
	LIMIT $1 OFFSET $2;
	`

	rows, err := r.db.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var products []*models.Product

	for rows.Next() {
		var p models.Product
		if err := rows.Scan(
			&p.ID,
			&p.SupplierID,
			&p.ProductName,
			&p.Manufacturer,
			&p.Description,
			&p.CostPrice,
			&p.SalePrice,
			&p.StockQuantity,
			&p.Barcode,
			&p.Status,
			&p.AllowDiscount,
			&p.MaxDiscountPercent,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		products = append(products, &p)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return products, nil
}

func (r *productRepository) GetByID(ctx context.Context, id int64) (*models.Product, error) {
	const query = `
	SELECT id, supplier_id, product_name, manufacturer, product_description,
	       cost_price, sale_price, stock_quantity, barcode,
	       status, allow_discount, max_discount_percent,
	       created_at, updated_at
	FROM products
	WHERE id = $1;
	`

	var p models.Product
	err := r.db.QueryRow(ctx, query, id).Scan(
		&p.ID,
		&p.SupplierID,
		&p.ProductName,
		&p.Manufacturer,
		&p.Description,
		&p.CostPrice,
		&p.SalePrice,
		&p.StockQuantity,
		&p.Barcode,
		&p.Status,
		&p.AllowDiscount,
		&p.MaxDiscountPercent,
		&p.CreatedAt,
		&p.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errMsg.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return &p, nil
}

func (r *productRepository) GetByName(ctx context.Context, name string) ([]*models.Product, error) {
	const query = `
	SELECT id, supplier_id, product_name, manufacturer, product_description,
	       cost_price, sale_price, stock_quantity, barcode,
	       status, allow_discount, max_discount_percent,
	       created_at, updated_at
	FROM products
	WHERE product_name ILIKE '%' || $1 || '%'
	ORDER BY product_name;
	`

	rows, err := r.db.Query(ctx, query, name)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(
			&p.ID,
			&p.SupplierID,
			&p.ProductName,
			&p.Manufacturer,
			&p.Description,
			&p.CostPrice,
			&p.SalePrice,
			&p.StockQuantity,
			&p.Barcode,
			&p.Status,
			&p.AllowDiscount,
			&p.MaxDiscountPercent,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		products = append(products, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return products, nil
}

func (r *productRepository) GetByManufacturer(ctx context.Context, manufacturer string) ([]*models.Product, error) {
	const query = `
	SELECT id, supplier_id, product_name, manufacturer, product_description,
	       cost_price, sale_price, stock_quantity, barcode,
	       status, allow_discount, max_discount_percent,
	       created_at, updated_at
	FROM products
	WHERE manufacturer ILIKE '%' || $1 || '%'
	ORDER BY product_name;
	`

	rows, err := r.db.Query(ctx, query, manufacturer)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var products []*models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(
			&p.ID,
			&p.SupplierID,
			&p.ProductName,
			&p.Manufacturer,
			&p.Description,
			&p.CostPrice,
			&p.SalePrice,
			&p.StockQuantity,
			&p.Barcode,
			&p.Status,
			&p.AllowDiscount,
			&p.MaxDiscountPercent,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		products = append(products, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return products, nil
}

func (r *productRepository) GetVersionByID(ctx context.Context, id int64) (int64, error) {
	const query = `SELECT version FROM products WHERE id = $1`

	var version int64
	err := r.db.QueryRow(ctx, query, id).Scan(&version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errMsg.ErrNotFound
		}
		return 0, fmt.Errorf("%w: %v", errMsg.ErrGetVersion, err)
	}

	return version, nil
}

func (r *productRepository) EnableProduct(ctx context.Context, uid int64) error {
	const query = `
		UPDATE products
		SET status = TRUE, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version, updated_at;
	`

	var version int
	var updatedAt time.Time
	err := r.db.QueryRow(ctx, query, uid).Scan(&version, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrEnable, err)
	}

	return nil
}

func (r *productRepository) DisableProduct(ctx context.Context, uid int64) error {
	const query = `
		UPDATE products
		SET status = FALSE, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version, updated_at;
	`

	var version int
	var updatedAt time.Time
	err := r.db.QueryRow(ctx, query, uid).Scan(&version, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrDisable, err)
	}

	return nil
}

func (r *productRepository) Update(ctx context.Context, product *models.Product) (*models.Product, error) {
	const query = `
		UPDATE products
		SET
			supplier_id = $1,
			product_name = $2,
			manufacturer = $3,
			product_description = $4,
			cost_price = $5,
			sale_price = $6,
			stock_quantity = $7,
			barcode = $8,
			status = $9,
			allow_discount = $10,
			max_discount_percent = $11,
			updated_at = NOW(),
			version = version + 1
		WHERE id = $12 AND version = $13
		RETURNING created_at, updated_at, version;
	`

	err := r.db.QueryRow(ctx, query,
		product.SupplierID,
		product.ProductName,
		product.Manufacturer,
		product.Description,
		product.CostPrice,
		product.SalePrice,
		product.StockQuantity,
		product.Barcode,
		product.Status,
		product.AllowDiscount,
		product.MaxDiscountPercent,
		product.ID,
		product.Version,
	).Scan(&product.CreatedAt, &product.UpdatedAt, &product.Version)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errMsg.ErrVersionConflict
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return product, nil
}

func (r *productRepository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM products WHERE id = $1;`

	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	if tag.RowsAffected() == 0 {
		return errMsg.ErrNotFound
	}

	return nil
}

func (r *productRepository) UpdateStock(ctx context.Context, id int64, quantity int) error {
	const query = `
		UPDATE products
		SET stock_quantity = $2, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version;
	`

	var version int
	err := r.db.QueryRow(ctx, query, id, quantity).Scan(&version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (r *productRepository) IncreaseStock(ctx context.Context, id int64, amount int) error {
	const query = `
		UPDATE products
		SET stock_quantity = stock_quantity + $2, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version;
	`

	var version int
	err := r.db.QueryRow(ctx, query, id, amount).Scan(&version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (r *productRepository) DecreaseStock(ctx context.Context, id int64, amount int) error {
	const query = `
		UPDATE products
		SET stock_quantity = GREATEST(COALESCE(stock_quantity, 0) - $2, 0),
		    updated_at = NOW(),
		    version = version + 1
		WHERE id = $1
		RETURNING version;
	`

	var version int
	err := r.db.QueryRow(ctx, query, id, amount).Scan(&version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (r *productRepository) GetStock(ctx context.Context, id int64) (int, error) {
	const query = `
		SELECT COALESCE(stock_quantity, 0)
		FROM products
		WHERE id = $1;
	`

	var stock int
	err := r.db.QueryRow(ctx, query, id).Scan(&stock)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errMsg.ErrNotFound
		}
		return 0, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return stock, nil
}

func (r *productRepository) EnableDiscount(ctx context.Context, id int64) error {
	const query = `
		UPDATE products
		SET allow_discount = TRUE, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version;
	`

	var version int
	err := r.db.QueryRow(ctx, query, id).Scan(&version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrProductEnableDiscount, err)
	}

	return nil
}

func (r *productRepository) DisableDiscount(ctx context.Context, id int64) error {
	const query = `
		UPDATE products
		SET allow_discount = FALSE, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version;
	`

	var version int
	err := r.db.QueryRow(ctx, query, id).Scan(&version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrProductDisableDiscount, err)
	}

	return nil
}

func (r *productRepository) ApplyDiscount(ctx context.Context, id int64, percent float64) (*models.Product, error) {
	const query = `
		UPDATE products
		SET max_discount_percent = $2, updated_at = NOW(), version = version + 1
		WHERE id = $1 AND allow_discount = TRUE
		RETURNING id, product_name, sale_price, max_discount_percent, allow_discount, version, updated_at;
	`

	var p models.Product
	err := r.db.QueryRow(ctx, query, id, percent).Scan(
		&p.ID,
		&p.ProductName,
		&p.SalePrice,
		&p.MaxDiscountPercent,
		&p.AllowDiscount,
		&p.Version,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {

			const checkQuery = `SELECT 1 FROM products WHERE id = $1`
			var exists int
			errCheck := r.db.QueryRow(ctx, checkQuery, id).Scan(&exists)
			if errCheck != nil || exists == 0 {
				return nil, errMsg.ErrNotFound
			}
			return nil, errMsg.ErrProductDiscountNotAllowed
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrProductApplyDiscount, err)
	}

	return &p, nil
}
