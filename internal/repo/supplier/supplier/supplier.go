package repositories

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SupplierRepository interface {
	Create(ctx context.Context, supplier *models.Supplier) (*models.Supplier, error)
	GetByID(ctx context.Context, id int64) (*models.Supplier, error)
	GetByName(ctx context.Context, name string) ([]*models.Supplier, error)
	GetVersionByID(ctx context.Context, id int64) (int64, error)
	GetAll(ctx context.Context) ([]*models.Supplier, error)
	Update(ctx context.Context, supplier *models.Supplier) error
	Delete(ctx context.Context, id int64) error
	Disable(ctx context.Context, id int64) error
	Enable(ctx context.Context, id int64) error
}

type supplierRepository struct {
	db     *pgxpool.Pool
	logger logger.LoggerAdapterInterface
}

func NewSupplierRepository(db *pgxpool.Pool, logger logger.LoggerAdapterInterface) SupplierRepository {
	return &supplierRepository{db: db, logger: logger}
}

func (r *supplierRepository) Create(ctx context.Context, supplier *models.Supplier) (*models.Supplier, error) {
	ref := "[supplierRepository - Create] - "

	r.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{
		"name":   supplier.Name,
		"cnpj":   supplier.CNPJ,
		"cpf":    supplier.CPF,
		"status": supplier.Status,
	})

	const query = `
		INSERT INTO suppliers (name, cnpj, cpf, status)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		supplier.Name,
		supplier.CNPJ,
		supplier.CPF,
		supplier.Status,
	).Scan(&supplier.ID, &supplier.CreatedAt, &supplier.UpdatedAt)

	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"name":   supplier.Name,
			"cnpj":   supplier.CNPJ,
			"cpf":    supplier.CPF,
			"status": supplier.Status,
		})
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	r.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"supplier_id": supplier.ID,
	})

	return supplier, nil
}

func (r *supplierRepository) GetByID(ctx context.Context, id int64) (*models.Supplier, error) {
	ref := "[supplierRepository - GetByID] - "

	r.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"supplier_id": id,
	})

	const query = `
		SELECT id, name, cnpj, cpf, status, created_at, updated_at
		FROM suppliers
		WHERE id = $1
	`

	var supplier models.Supplier
	err := r.db.QueryRow(ctx, query, id).Scan(
		&supplier.ID,
		&supplier.Name,
		&supplier.CNPJ,
		&supplier.CPF,
		&supplier.Status,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{"supplier_id": id})
		return nil, errMsg.ErrNotFound
	}
	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogGetError, nil)
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"supplier_id": supplier.ID,
		"name":        supplier.Name,
	})

	return &supplier, nil
}

func (r *supplierRepository) GetByName(ctx context.Context, name string) ([]*models.Supplier, error) {
	ref := "[supplierRepository - GetByName] - "
	r.logger.Info(ctx, ref+"início da busca", map[string]any{
		"name_partial": name,
	})

	const query = `
		SELECT id, name, cnpj, cpf, status, created_at, updated_at
		FROM suppliers
		WHERE name ILIKE $1
		ORDER BY name ASC
	`

	rows, err := r.db.Query(ctx, query, "%"+name+"%")
	if err != nil {
		r.logger.Error(ctx, err, ref+"erro na busca", map[string]any{
			"name_partial": name,
		})
		return nil, fmt.Errorf("falha ao buscar fornecedores: %w", err)
	}
	defer rows.Close()

	var suppliers []*models.Supplier
	for rows.Next() {
		supplier := new(models.Supplier)
		if err := rows.Scan(
			&supplier.ID,
			&supplier.Name,
			&supplier.CNPJ,
			&supplier.CPF,
			&supplier.Status,
			&supplier.CreatedAt,
			&supplier.UpdatedAt,
		); err != nil {
			r.logger.Error(ctx, err, ref+"erro no scan", map[string]any{
				"name_partial": name,
			})
			return nil, fmt.Errorf("falha ao ler dados do fornecedor: %w", err)
		}
		suppliers = append(suppliers, supplier)
	}

	if len(suppliers) == 0 {
		r.logger.Warn(ctx, ref+"não encontrado", map[string]any{
			"name_partial": name,
		})
		return nil, errMsg.ErrNotFound
	}

	r.logger.Info(ctx, ref+"sucesso na busca", map[string]any{
		"count": len(suppliers),
	})

	return suppliers, nil
}

func (r *supplierRepository) GetAll(ctx context.Context) ([]*models.Supplier, error) {
	ref := "[supplierRepository - GetAll] - "

	r.logger.Info(ctx, ref+logger.LogGetInit, nil)

	const query = `
		SELECT id, name, cnpj, cpf, status, created_at, updated_at
		FROM suppliers
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogGetError, nil)
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var suppliers []*models.Supplier
	for rows.Next() {
		var s models.Supplier
		if err := rows.Scan(
			&s.ID,
			&s.Name,
			&s.CNPJ,
			&s.CPF,
			&s.Status,
			&s.CreatedAt,
			&s.UpdatedAt,
		); err != nil {
			r.logger.Error(ctx, err, ref+logger.LogGetError, nil)
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		suppliers = append(suppliers, &s)
	}

	if err = rows.Err(); err != nil {
		r.logger.Error(ctx, err, ref+logger.LogGetError, nil)
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"total": len(suppliers),
	})

	return suppliers, nil
}

func (r *supplierRepository) Update(ctx context.Context, supplier *models.Supplier) error {
	ref := "[supplierRepository - Update] - "

	r.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"supplier_id": supplier.ID,
	})

	const query = `
		UPDATE suppliers
		SET
			name       = $1,
			cnpj       = $2,
			cpf        = $3,
			status     = $4,
			updated_at = NOW(),
			version    = version + 1
		WHERE
			id      = $5 AND
			version = $6
		RETURNING id, name, cnpj, cpf, status, created_at, updated_at, version
	`

	err := r.db.QueryRow(ctx, query,
		supplier.Name,
		supplier.CNPJ,
		supplier.CPF,
		supplier.Status,
		supplier.ID,
		supplier.Version,
	).Scan(
		&supplier.ID,
		&supplier.Name,
		&supplier.CNPJ,
		&supplier.CPF,
		&supplier.Status,
		&supplier.CreatedAt,
		&supplier.UpdatedAt,
		&supplier.Version,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn(ctx, ref+logger.LogUpdateVersionConflict, map[string]any{
				"supplier_id": supplier.ID,
				"version":     supplier.Version,
			})
			return errMsg.ErrVersionConflict
		}
		r.logger.Error(ctx, err, ref+logger.LogUpdateError, nil)
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	r.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"supplier_id": supplier.ID,
	})

	return nil
}

func (r *supplierRepository) Delete(ctx context.Context, id int64) error {
	ref := "[supplierRepository - Delete] - "

	r.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{
		"supplier_id": id,
	})

	const query = `DELETE FROM suppliers WHERE id = $1`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogDeleteError, nil)
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	if cmdTag.RowsAffected() == 0 {
		r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
			"supplier_id": id,
		})
		return errMsg.ErrNotFound
	}

	r.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"supplier_id": id,
	})

	return nil
}

func (r *supplierRepository) Disable(ctx context.Context, id int64) error {
	ref := "[supplierRepository - Disable] - "

	r.logger.Info(ctx, ref+logger.LogDisableInit, map[string]any{
		"supplier_id": id,
	})

	const query = `
		UPDATE suppliers
		SET status = false,
		    updated_at = NOW()
		WHERE id = $1
	`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogDisableError, nil)
		return fmt.Errorf("%w: %v", errMsg.ErrDisable, err)
	}

	if cmdTag.RowsAffected() == 0 {
		r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{"supplier_id": id})
		return errMsg.ErrNotFound
	}

	r.logger.Info(ctx, ref+logger.LogDisableSuccess, map[string]any{"supplier_id": id})
	return nil
}

func (r *supplierRepository) Enable(ctx context.Context, id int64) error {
	ref := "[supplierRepository - Enable] - "

	r.logger.Info(ctx, ref+logger.LogEnableInit, map[string]any{
		"supplier_id": id,
	})

	const query = `
		UPDATE suppliers
		SET status = true,
		    updated_at = NOW()
		WHERE id = $1
	`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogEnableError, nil)
		return fmt.Errorf("%w: %v", errMsg.ErrEnable, err)
	}

	if cmdTag.RowsAffected() == 0 {
		r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{"supplier_id": id})
		return errMsg.ErrNotFound
	}

	r.logger.Info(ctx, ref+logger.LogEnableSuccess, map[string]any{"supplier_id": id})
	return nil
}

func (r *supplierRepository) GetVersionByID(ctx context.Context, id int64) (int64, error) {
	ref := "[supplierRepository - GetVersionByID] - "

	r.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"supplier_id": id,
	})

	const query = `
		SELECT version
		FROM suppliers
		WHERE id = $1
	`

	var version int64
	err := r.db.QueryRow(ctx, query, id).Scan(&version)

	if errors.Is(err, pgx.ErrNoRows) {
		r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{"supplier_id": id})
		return 0, errMsg.ErrNotFound
	}
	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogGetError, nil)
		return 0, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"supplier_id": id,
		"version":     version,
	})

	return version, nil
}
