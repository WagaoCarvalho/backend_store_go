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
	ErrCreateAddress   = errors.New("erro ao criar endereço")
	ErrFetchAddress    = errors.New("erro ao buscar endereço")
	ErrAddressNotFound = errors.New("endereço não encontrado")
	ErrUpdateAddress   = errors.New("erro ao atualizar endereço")
	ErrDeleteAddress   = errors.New("erro ao excluir endereço")
)

type AddressRepository interface {
	Create(ctx context.Context, address *models.Address) (*models.Address, error)
	GetByID(ctx context.Context, id int64) (*models.Address, error)
	GetByUserID(ctx context.Context, userID int64) ([]*models.Address, error)
	GetByClientID(ctx context.Context, clientID int64) ([]*models.Address, error)
	GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.Address, error)
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
			return nil, fmt.Errorf("%w: %v", ErrFetchAddress, err)
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
			return nil, fmt.Errorf("%w: %v", ErrFetchAddress, err)
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
			return nil, fmt.Errorf("%w: %v", ErrFetchAddress, err)
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
			return ErrAddressNotFound
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
