package repositories

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	logger "github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
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
	logger logger.LogAdapterInterface
}

func NewAddressRepository(db *pgxpool.Pool, logger logger.LogAdapterInterface) AddressRepository {
	return &addressRepository{db: db, logger: logger}
}

func (r *addressRepository) Create(ctx context.Context, address *models.Address) (*models.Address, error) {
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
		if errMsgPg.IsForeignKeyViolation(err) {
			return nil, errMsg.ErrInvalidForeignKey
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return address, nil
}

func (r *addressRepository) CreateTx(ctx context.Context, tx pgx.Tx, address *models.Address) (*models.Address, error) {
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
		if errMsgPg.IsForeignKeyViolation(err) {
			return nil, errMsg.ErrInvalidForeignKey
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return address, nil
}

func (r *addressRepository) GetByID(ctx context.Context, id int64) (*models.Address, error) {
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
			return nil, errMsg.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return &address, nil
}

func (r *addressRepository) GetByUserID(ctx context.Context, userID int64) ([]*models.Address, error) {
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
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
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
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		addresses = append(addresses, &address)
	}

	return addresses, nil
}

func (r *addressRepository) GetByClientID(ctx context.Context, clientID int64) ([]*models.Address, error) {
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
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
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
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		addresses = append(addresses, &address)
	}

	return addresses, nil
}

func (r *addressRepository) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.Address, error) {
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
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
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
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		addresses = append(addresses, &address)
	}

	return addresses, nil
}

func (r *addressRepository) Update(ctx context.Context, address *models.Address) error {
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
			return errMsg.ErrNotFound
		}

		if errMsgPg.IsForeignKeyViolation(err) {
			return errMsg.ErrInvalidForeignKey
		}

		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (r *addressRepository) Delete(ctx context.Context, id int64) error {
	const query = `
		DELETE FROM addresses 
		WHERE id = $1
	`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return errMsg.ErrNotFound
	}

	return nil
}
