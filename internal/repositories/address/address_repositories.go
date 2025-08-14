package repositories

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
	logger "github.com/WagaoCarvalho/backend_store_go/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AddressRepository interface {
	Create(ctx context.Context, address *models.Address) (*models.Address, error)
	CreateTx(ctx context.Context, tx pgx.Tx, address *models.Address) (*models.Address, error)
	GetByID(ctx context.Context, id int64) (*models.Address, error)
	GetByUserID(ctx context.Context, userID int64) ([]*models.Address, error)
	GetByClientID(ctx context.Context, clientID int64) ([]*models.Address, error)
	GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.Address, error)
	Update(ctx context.Context, address *models.Address) error
	Delete(ctx context.Context, id int64) error
}

type addressRepository struct {
	db     *pgxpool.Pool
	logger logger.LoggerAdapterInterface
}

func NewAddressRepository(db *pgxpool.Pool, logger logger.LoggerAdapterInterface) AddressRepository {
	return &addressRepository{db: db, logger: logger}
}

func (r *addressRepository) Create(ctx context.Context, address *models.Address) (*models.Address, error) {
	ref := "[addressRepository - Create] - "
	r.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{
		"user_id":     utils.Int64OrNil(address.UserID),
		"client_id":   utils.Int64OrNil(address.ClientID),
		"supplier_id": utils.Int64OrNil(address.SupplierID),
		"street":      address.Street,
	})
	const query = `
		INSERT INTO addresses (
			user_id, client_id, supplier_id,
			street, city, state, country, postal_code,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING id, created_at, updated_at;
	`

	err := r.db.QueryRow(ctx, query,
		address.UserID,
		address.ClientID,
		address.SupplierID,
		address.Street,
		address.City,
		address.State,
		address.Country,
		address.PostalCode,
	).Scan(&address.ID, &address.CreatedAt, &address.UpdatedAt)

	if err != nil {
		if IsForeignKeyViolation(err) {
			r.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{
				"user_id":     utils.Int64OrNil(address.UserID),
				"client_id":   utils.Int64OrNil(address.ClientID),
				"supplier_id": utils.Int64OrNil(address.SupplierID),
			})
			return nil, ErrInvalidForeignKey
		}

		r.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"user_id":     utils.Int64OrNil(address.UserID),
			"client_id":   utils.Int64OrNil(address.ClientID),
			"supplier_id": utils.Int64OrNil(address.SupplierID),
			"street":      address.Street,
		})
		return nil, fmt.Errorf("%w: %v", ErrCreateAddress, err)
	}

	r.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"address_id": address.ID,
	})

	return address, nil
}

func (r *addressRepository) CreateTx(ctx context.Context, tx pgx.Tx, address *models.Address) (*models.Address, error) {
	ref := "[addressRepository - CreateTx] - "
	r.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{
		"user_id":     utils.Int64OrNil(address.UserID),
		"client_id":   utils.Int64OrNil(address.ClientID),
		"supplier_id": utils.Int64OrNil(address.SupplierID),
		"street":      address.Street,
	})

	const query = `
		INSERT INTO addresses (
			user_id, client_id, supplier_id,
			street, city, state, country, postal_code,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING id, created_at, updated_at;
	`

	err := tx.QueryRow(ctx, query,
		address.UserID,
		address.ClientID,
		address.SupplierID,
		address.Street,
		address.City,
		address.State,
		address.Country,
		address.PostalCode,
	).Scan(&address.ID, &address.CreatedAt, &address.UpdatedAt)

	if err != nil {
		if IsForeignKeyViolation(err) {
			r.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{
				"user_id":     utils.Int64OrNil(address.UserID),
				"client_id":   utils.Int64OrNil(address.ClientID),
				"supplier_id": utils.Int64OrNil(address.SupplierID),
			})
			return nil, ErrInvalidForeignKey
		}

		r.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"user_id":     utils.Int64OrNil(address.UserID),
			"client_id":   utils.Int64OrNil(address.ClientID),
			"supplier_id": utils.Int64OrNil(address.SupplierID),
			"street":      address.Street,
		})
		return nil, fmt.Errorf("%w: %v", ErrCreateAddress, err)
	}

	r.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"address_id": address.ID,
	})

	return address, nil
}

func (r *addressRepository) GetByID(ctx context.Context, id int64) (*models.Address, error) {
	ref := "[addressRepository - GetByID] - "
	r.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"address_id": id,
	})
	const query = `
		SELECT 
			id, user_id, client_id, supplier_id,
			street, city, state, country, postal_code,
			created_at, updated_at
		FROM addresses
		WHERE id = $1;
	`

	var address models.Address
	err := r.db.QueryRow(ctx, query, id).Scan(
		&address.ID,
		&address.UserID,
		&address.ClientID,
		&address.SupplierID,
		&address.Street,
		&address.City,
		&address.State,
		&address.Country,
		&address.PostalCode,
		&address.CreatedAt,
		&address.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Info(ctx, ref+logger.LogNotFound, map[string]any{
				"address_id": id,
			})
			return nil, ErrAddressNotFound
		}

		r.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"address_id": id,
		})
		return nil, fmt.Errorf("%w: %v", ErrFetchAddress, err)
	}

	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"address_id": address.ID,
	})

	return &address, nil
}

func (r *addressRepository) GetByUserID(ctx context.Context, userID int64) ([]*models.Address, error) {
	ref := "[addressRepository - GetByUserID] - "
	r.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"user_id": userID,
	})
	const query = `
		SELECT 
			id, user_id, client_id, supplier_id,
			street, city, state, country, postal_code,
			created_at, updated_at
		FROM addresses
		WHERE user_id = $1;
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"user_id": userID,
		})
		return nil, fmt.Errorf("%w: %v", ErrFetchAddress, err)
	}
	defer rows.Close()

	var addresses []*models.Address
	for rows.Next() {
		var address models.Address
		if err := rows.Scan(
			&address.ID,
			&address.UserID,
			&address.ClientID,
			&address.SupplierID,
			&address.Street,
			&address.City,
			&address.State,
			&address.Country,
			&address.PostalCode,
			&address.CreatedAt,
			&address.UpdatedAt,
		); err != nil {
			r.logger.Error(ctx, err, ref+logger.LogGetErrorScan, map[string]any{
				"user_id": userID,
			})
			return nil, fmt.Errorf("%w: %v", ErrFetchAddress, err)
		}
		addresses = append(addresses, &address)
	}

	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"user_id":       userID,
		"total_results": len(addresses),
	})

	return addresses, nil
}

func (r *addressRepository) GetByClientID(ctx context.Context, clientID int64) ([]*models.Address, error) {
	ref := "[addressRepository - GetByClientID] - "
	r.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"client_id": clientID,
	})

	const query = `
		SELECT 
			id, user_id, client_id, supplier_id,
			street, city, state, country, postal_code,
			created_at, updated_at
		FROM addresses
		WHERE client_id = $1;
	`

	rows, err := r.db.Query(ctx, query, clientID)
	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"client_id": clientID,
		})
		return nil, fmt.Errorf("%w: %v", ErrFetchAddress, err)
	}
	defer rows.Close()

	var addresses []*models.Address
	for rows.Next() {
		var address models.Address
		if err := rows.Scan(
			&address.ID,
			&address.UserID,
			&address.ClientID,
			&address.SupplierID,
			&address.Street,
			&address.City,
			&address.State,
			&address.Country,
			&address.PostalCode,
			&address.CreatedAt,
			&address.UpdatedAt,
		); err != nil {
			r.logger.Error(ctx, err, ref+logger.LogGetErrorScan, map[string]any{
				"client_id": clientID,
			})
			return nil, fmt.Errorf("%w: %v", ErrFetchAddress, err)
		}
		addresses = append(addresses, &address)
	}

	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"client_id":     clientID,
		"total_results": len(addresses),
	})

	return addresses, nil
}

func (r *addressRepository) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.Address, error) {
	ref := "[addressRepository - GetBySupplierID] - "
	r.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"supplier_id": supplierID,
	})

	const query = `
		SELECT 
			id, user_id, client_id, supplier_id,
			street, city, state, country, postal_code,
			created_at, updated_at
		FROM addresses
		WHERE supplier_id = $1;
	`

	rows, err := r.db.Query(ctx, query, supplierID)
	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"supplier_id": supplierID,
		})
		return nil, fmt.Errorf("%w: %v", ErrFetchAddress, err)
	}
	defer rows.Close()

	var addresses []*models.Address
	for rows.Next() {
		var address models.Address
		if err := rows.Scan(
			&address.ID,
			&address.UserID,
			&address.ClientID,
			&address.SupplierID,
			&address.Street,
			&address.City,
			&address.State,
			&address.Country,
			&address.PostalCode,
			&address.CreatedAt,
			&address.UpdatedAt,
		); err != nil {
			r.logger.Error(ctx, err, ref+logger.LogGetErrorScan, map[string]any{
				"supplier_id": supplierID,
			})
			return nil, fmt.Errorf("%w: %v", ErrFetchAddress, err)
		}
		addresses = append(addresses, &address)
	}

	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"supplier_id":   supplierID,
		"total_results": len(addresses),
	})

	return addresses, nil
}

func (r *addressRepository) Update(ctx context.Context, address *models.Address) error {
	ref := "[addressRepository - Update] - "
	r.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"address_id":  address.ID,
		"user_id":     utils.Int64OrNil(address.UserID),
		"client_id":   utils.Int64OrNil(address.ClientID),
		"supplier_id": utils.Int64OrNil(address.SupplierID),
		"street":      address.Street,
	})

	const query = `
		UPDATE addresses
		SET 
			user_id     = $1,
			client_id   = $2,
			supplier_id = $3,
			street      = $4,
			city        = $5,
			state       = $6,
			country     = $7,
			postal_code = $8,
			updated_at  = NOW()
		WHERE id = $9
		RETURNING updated_at;
	`

	err := r.db.QueryRow(ctx, query,
		address.UserID,
		address.ClientID,
		address.SupplierID,
		address.Street,
		address.City,
		address.State,
		address.Country,
		address.PostalCode,
		address.ID,
	).Scan(&address.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Info(ctx, ref+logger.LogNotFound, map[string]any{
				"address_id": address.ID,
			})
			return ErrAddressNotFound
		}

		r.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"address_id": address.ID,
		})
		return fmt.Errorf("%w: %v", ErrUpdateAddress, err)
	}

	r.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"address_id": address.ID,
	})

	return nil
}

func (r *addressRepository) Delete(ctx context.Context, id int64) error {
	ref := "[addressRepository - Delete] - "
	r.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{
		"address_id": id,
	})
	const query = `
		DELETE FROM addresses 
		WHERE id = $1
	`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"address_id": id,
		})
		return fmt.Errorf("%w: %v", ErrDeleteAddress, err)
	}

	if cmdTag.RowsAffected() == 0 {
		r.logger.Info(ctx, ref+logger.LogNotFound, map[string]any{
			"address_id": id,
		})
		return ErrAddressNotFound
	}

	r.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"address_id": id,
	})

	return nil
}
