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
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository interface {
	Create(ctx context.Context, product *models.Product) (*models.Product, error)
	GetAll(ctx context.Context, limit, offset int) ([]*models.Product, error)
	GetById(ctx context.Context, id int64) (*models.Product, error)
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
	ref := "[productRepository - Create] - "

	r.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{
		"supplier_id":          utils.Int64OrNil(product.SupplierID),
		"product_name":         product.ProductName,
		"manufacturer":         product.Manufacturer,
		"cost_price":           product.CostPrice,
		"sale_price":           product.SalePrice,
		"stock_quantity":       product.StockQuantity,
		"status":               product.Status,
		"allow_discount":       product.AllowDiscount,
		"max_discount_percent": product.MaxDiscountPercent,
	})

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
			r.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{
				"supplier_id": utils.Int64OrNil(product.SupplierID),
			})
			return nil, errMsg.ErrInvalidForeignKey
		}

		r.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"supplier_id":          utils.Int64OrNil(product.SupplierID),
			"product_name":         product.ProductName,
			"manufacturer":         product.Manufacturer,
			"cost_price":           product.CostPrice,
			"sale_price":           product.SalePrice,
			"stock_quantity":       product.StockQuantity,
			"status":               product.Status,
			"allow_discount":       product.AllowDiscount,
			"max_discount_percent": product.MaxDiscountPercent,
		})
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	r.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"product_id": product.ID,
	})

	return product, nil
}

func (r *productRepository) GetAll(ctx context.Context, limit, offset int) ([]*models.Product, error) {
	ref := "[productRepository - GetAll] - "
	r.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"limit":  limit,
		"offset": offset,
	})

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
		r.logger.Error(ctx, err, ref+logger.LogQueryError, map[string]any{
			"limit":  limit,
			"offset": offset,
		})
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var products []*models.Product

	for rows.Next() {
		var p models.Product
		err := rows.Scan(
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
			r.logger.Error(ctx, err, ref+logger.LogGetErrorScan, nil)
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		products = append(products, &p)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error(ctx, err, ref+logger.LogIterateError, nil)
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"total": len(products),
	})

	return products, nil
}

func (r *productRepository) GetById(ctx context.Context, id int64) (*models.Product, error) {
	ref := "[productRepository - GetById] - "
	r.logger.Info(ctx, ref+logger.LogLoginInit, map[string]any{"id": id})

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
			r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{"id": id})
			return nil, errMsg.ErrNotFound
		}

		r.logger.Error(ctx, err, ref+logger.LogQueryError, map[string]any{"id": id})
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{"id": id})
	return &p, nil
}

func (r *productRepository) GetByName(ctx context.Context, name string) ([]*models.Product, error) {
	ref := "[productRepository - GetByName] - "
	r.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{"name": name})

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
		r.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{"name": name})
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	products := make([]*models.Product, 0)
	for rows.Next() {
		var p models.Product
		err := rows.Scan(
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
			r.logger.Error(ctx, err, ref+logger.LogScanError, nil)
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		products = append(products, &p)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error(ctx, err, ref+logger.LogIterateError, nil)
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{"count": len(products)})
	return products, nil
}

func (r *productRepository) GetByManufacturer(ctx context.Context, manufacturer string) ([]*models.Product, error) {
	ref := "[productRepository - GetByManufacturer] - "
	r.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{"manufacturer": manufacturer})

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
		r.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{"manufacturer": manufacturer})
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	products := make([]*models.Product, 0)
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
			r.logger.Error(ctx, err, ref+logger.LogGetErrorScan, nil)
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		products = append(products, &p)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error(ctx, err, ref+logger.LogIterateError, nil)
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	r.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{"count": len(products)})
	return products, nil
}

func (r *productRepository) GetVersionByID(ctx context.Context, id int64) (int64, error) {
	ref := "[productRepository - GetVersionByID] - "
	r.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"product_id": id,
	})

	const query = `SELECT version FROM products WHERE id = $1`

	var version int64
	err := r.db.QueryRow(ctx, query, id).Scan(&version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"product_id": id,
			})
			return 0, errMsg.ErrNotFound
		}

		r.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"product_id": id,
		})
		return 0, fmt.Errorf("%w: %v", errMsg.ErrGetVersion, err)
	}

	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"product_id": id,
		"version":    version,
	})

	return version, nil
}

func (r *productRepository) EnableProduct(ctx context.Context, uid int64) error {
	ref := "[productRepository - Enable] - "
	r.logger.Info(ctx, ref+logger.LogEnableInit, map[string]any{
		"product_id": uid,
	})

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
			r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"product_id": uid,
			})
			return errMsg.ErrNotFound
		}

		r.logger.Error(ctx, err, ref+logger.LogEnableError, map[string]any{
			"product_id": uid,
		})
		return fmt.Errorf("%w: %v", errMsg.ErrEnable, err)
	}

	r.logger.Info(ctx, ref+logger.LogEnableSuccess, map[string]any{
		"product_id": uid,
		"new_status": true,
	})

	return nil
}

func (r *productRepository) DisableProduct(ctx context.Context, uid int64) error {
	ref := "[productRepository - Disable] - "
	r.logger.Info(ctx, ref+logger.LogDisableInit, map[string]any{
		"product_id": uid,
	})

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
			r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"product_id": uid,
			})
			return errMsg.ErrNotFound
		}

		r.logger.Error(ctx, err, ref+logger.LogDisableError, map[string]any{
			"product_id": uid,
		})
		return fmt.Errorf("%w: %v", errMsg.ErrDisable, err)
	}

	r.logger.Info(ctx, ref+logger.LogDisableSuccess, map[string]any{
		"product_id": uid,
		"new_status": false,
	})

	return nil
}

func (r *productRepository) Update(ctx context.Context, product *models.Product) (*models.Product, error) {
	ref := "[productRepository - Update] - "
	r.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"id":                   product.ID,
		"name":                 product.ProductName,
		"manufacturer":         product.Manufacturer,
		"version":              product.Version,
		"allow_discount":       product.AllowDiscount,
		"max_discount_percent": product.MaxDiscountPercent,
	})

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
			r.logger.Warn(ctx, ref+logger.LogUpdateVersionConflict, map[string]any{
				"id":      product.ID,
				"version": product.Version,
			})
			return nil, errMsg.ErrVersionConflict
		}

		r.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"id": product.ID,
		})
		return nil, fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	r.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"id":      product.ID,
		"version": product.Version,
	})

	return product, nil
}

func (r *productRepository) Delete(ctx context.Context, id int64) error {
	ref := "[productRepository - Delete] - "
	r.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{
		"id": id,
	})

	const query = `DELETE FROM products WHERE id = $1;`

	tag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"id": id,
		})
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	if tag.RowsAffected() == 0 {
		r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
			"id": id,
		})
		return errMsg.ErrNotFound
	}

	r.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"id": id,
	})

	return nil
}

func (r *productRepository) UpdateStock(ctx context.Context, id int64, quantity int) error {
	ref := "[productRepository - UpdateStock] - "
	r.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"product_id": id,
		"quantity":   quantity,
	})

	const query = `
		UPDATE products
		SET stock_quantity = $2, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version, updated_at;
	`

	var version int
	var updatedAt time.Time
	err := r.db.QueryRow(ctx, query, id, quantity).Scan(&version, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"product_id": id,
			})
			return errMsg.ErrNotFound
		}

		r.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"product_id": id,
			"quantity":   quantity,
		})
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	r.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"product_id": id,
		"quantity":   quantity,
	})

	return nil
}

func (r *productRepository) IncreaseStock(ctx context.Context, id int64, amount int) error {
	ref := "[productRepository - IncreaseStock] - "
	r.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"product_id":     id,
		"stock_quantity": amount,
	})

	const query = `
		UPDATE products
		SET stock_quantity = stock_quantity + $2, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version, updated_at;
	`

	var version int
	var updatedAt time.Time
	err := r.db.QueryRow(ctx, query, id, amount).Scan(&version, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"product_id": id,
			})
			return errMsg.ErrNotFound
		}

		r.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"product_id":     id,
			"stock_quantity": amount,
		})
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	r.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"product_id":     id,
		"stock_quantity": amount,
	})

	return nil
}

func (r *productRepository) DecreaseStock(ctx context.Context, id int64, amount int) error {
	ref := "[productRepository - DecreaseStock] - "
	r.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"product_id":     id,
		"stock_quantity": amount,
	})

	const query = `
		UPDATE products
		SET stock_quantity = GREATEST(COALESCE(stock_quantity, 0) - $2, 0),
		    updated_at = NOW(),
		    version = version + 1
		WHERE id = $1
		RETURNING version, updated_at;
	`

	var version int
	var updatedAt time.Time
	err := r.db.QueryRow(ctx, query, id, amount).Scan(&version, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"product_id": id,
			})
			return errMsg.ErrNotFound
		}

		r.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"product_id":     id,
			"stock_quantity": amount,
		})
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	r.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"product_id":     id,
		"stock_quantity": amount,
	})

	return nil
}

func (r *productRepository) GetStock(ctx context.Context, id int64) (int, error) {
	ref := "[productRepository - GetStock] - "
	r.logger.Info(ctx, ref+logger.LogGetError, map[string]any{
		"product_id": id,
	})

	const query = `
		SELECT COALESCE(stock_quantity, 0)
		FROM products
		WHERE id = $1;
	`

	var stock int
	err := r.db.QueryRow(ctx, query, id).Scan(&stock)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"product_id": id,
			})
			return 0, errMsg.ErrNotFound
		}

		r.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"product_id": id,
		})
		return 0, fmt.Errorf("erro ao obter estoque: %w", err)
	}

	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"product_id": id,
		"stock":      stock,
	})

	return stock, nil
}

func (r *productRepository) EnableDiscount(ctx context.Context, id int64) error {
	ref := "[productRepository - EnableDiscount] - "
	r.logger.Info(ctx, ref+logger.LogEnableInit, map[string]any{
		"product_id": id,
	})

	const query = `
		UPDATE products
		SET allow_discount = TRUE, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version, updated_at;
	`

	var version int
	var updatedAt time.Time
	err := r.db.QueryRow(ctx, query, id).Scan(&version, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"product_id": id,
			})
			return errMsg.ErrNotFound
		}

		r.logger.Error(ctx, err, ref+logger.LogEnableError, map[string]any{
			"product_id": id,
		})
		return fmt.Errorf("%w: %v", errMsg.ErrProductEnableDiscount, err)
	}

	r.logger.Info(ctx, ref+logger.LogEnableSuccess, map[string]any{
		"product_id":     id,
		"allow_discount": true,
	})

	return nil
}

func (r *productRepository) DisableDiscount(ctx context.Context, id int64) error {
	ref := "[productRepository - DisableDiscount] - "
	r.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"product_id": id,
	})

	const query = `
		UPDATE products
		SET allow_discount = FALSE, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version, updated_at;
	`

	var version int
	var updatedAt time.Time
	err := r.db.QueryRow(ctx, query, id).Scan(&version, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"product_id": id,
			})
			return errMsg.ErrNotFound
		}

		r.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"product_id": id,
		})
		return fmt.Errorf("%w: %v", errMsg.ErrProductDisableDiscount, err)
	}

	r.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"product_id":     id,
		"allow_discount": false,
	})

	return nil
}

func (r *productRepository) ApplyDiscount(ctx context.Context, id int64, percent float64) (*models.Product, error) {
	ref := "[productRepository - ApplyDiscount] - "
	r.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"product_id": id,
		"percent":    percent,
	})

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
			// Verifica se o produto existe independentemente de allow_discount
			const checkQuery = `SELECT 1 FROM products WHERE id = $1`
			var exists int
			errCheck := r.db.QueryRow(ctx, checkQuery, id).Scan(&exists)
			if errCheck != nil || exists == 0 {
				r.logger.Warn(ctx, ref+"produto não encontrado", map[string]any{"product_id": id})
				return nil, errMsg.ErrNotFound
			}
			r.logger.Warn(ctx, ref+"desconto não permitido", map[string]any{"product_id": id})
			return nil, errMsg.ErrProductDiscountNotAllowed
		}

		r.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"product_id": id,
			"percent":    percent,
		})
		return nil, fmt.Errorf("%w: %v", errMsg.ErrProductApplyDiscount, err)
	}

	r.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"product_id": id,
		"percent":    percent,
	})

	return &p, nil
}
