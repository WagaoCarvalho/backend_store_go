package repositories

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrCreateAddress       = errors.New("erro ao criar endereço")
	ErrFetchAddress        = errors.New("erro ao buscar endereço")
	ErrAddressNotFound     = errors.New("endereço não encontrado")
	ErrUpdateAddress       = errors.New("erro ao atualizar endereço")
	ErrDeleteAddress       = errors.New("erro ao excluir endereço")
	ErrVersionConflict     = errors.New("conflito de versão: o endereço foi modificado por outra operação")
	ErrFetchAddressVersion = errors.New("erro ao buscar a versão do endereço")
)

type AddressRepository interface {
	Create(ctx context.Context, address *models.Address) (*models.Address, error)
	GetByID(ctx context.Context, id int64) (*models.Address, error)
	GetByUserID(ctx context.Context, userID int64) (*models.Address, error)
	GetByClientID(ctx context.Context, clientID int64) (*models.Address, error)
	GetBySupplierID(ctx context.Context, supplierID int64) (*models.Address, error)
	GetVersionByID(ctx context.Context, id int64) (int, error)
	Update(ctx context.Context, address *models.Address) error
	Delete(ctx context.Context, id int64) error
}

type addressRepository struct {
	db *pgxpool.Pool
}

func NewAddressRepository(db *pgxpool.Pool) AddressRepository {
	return &addressRepository{db: db}
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
		return nil, fmt.Errorf("%w: %v", ErrCreateAddress, err)
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
			return nil, ErrAddressNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrFetchAddress, err)
	}

	return &address, nil
}

func (r *addressRepository) GetByUserID(ctx context.Context, userID int64) (*models.Address, error) {
	const query = `
		SELECT 
			id, user_id, client_id, supplier_id,
			street, city, state, country, postal_code,
			created_at, updated_at
		FROM addresses
		WHERE user_id = $1;
	`

	var address models.Address
	err := r.db.QueryRow(ctx, query, userID).Scan(
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
			return nil, ErrAddressNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrFetchAddress, err)
	}

	return &address, nil
}

func (r *addressRepository) GetByClientID(ctx context.Context, clientID int64) (*models.Address, error) {
	const query = `
		SELECT 
			id, user_id, client_id, supplier_id,
			street, city, state, country, postal_code,
			created_at, updated_at
		FROM addresses
		WHERE client_id = $1;
	`

	var address models.Address
	err := r.db.QueryRow(ctx, query, clientID).Scan(
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
			return nil, ErrAddressNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrFetchAddress, err)
	}

	return &address, nil
}

func (r *addressRepository) GetBySupplierID(ctx context.Context, supplierID int64) (*models.Address, error) {
	const query = `
		SELECT 
			id, user_id, client_id, supplier_id,
			street, city, state, country, postal_code,
			created_at, updated_at
		FROM addresses
		WHERE supplier_id = $1;
	`

	var address models.Address
	err := r.db.QueryRow(ctx, query, supplierID).Scan(
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
			return nil, ErrAddressNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrFetchAddress, err)
	}

	return &address, nil
}

func (r *addressRepository) GetVersionByID(ctx context.Context, id int64) (int, error) {
	const query = `
		SELECT version
		FROM addresses
		WHERE id = $1;
	`

	var version int
	err := r.db.QueryRow(ctx, query, id).Scan(&version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrAddressNotFound
		}
		return 0, fmt.Errorf("%w: %v", ErrFetchAddressVersion, err)
	}

	return version, nil
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
			version     = version + 1,
			updated_at  = NOW()
		WHERE 
			id      = $9
			AND version = $10
		RETURNING 
			version,
			updated_at;
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
		address.Version,
	).Scan(&address.Version, &address.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {

			var exists bool
			checkQuery := `SELECT EXISTS(SELECT 1 FROM addresses WHERE id = $1)`
			checkErr := r.db.QueryRow(ctx, checkQuery, address.ID).Scan(&exists)
			if checkErr != nil {
				return fmt.Errorf("%w: erro ao verificar existência: %v", ErrUpdateAddress, checkErr)
			}
			if !exists {
				return ErrAddressNotFound
			}
			return ErrVersionConflict
		}
		return fmt.Errorf("%w: %v", ErrUpdateAddress, err)
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
		return fmt.Errorf("%w: %v", ErrDeleteAddress, err)
	}

	if cmdTag.RowsAffected() == 0 {
		return ErrAddressNotFound
	}

	return nil
}
