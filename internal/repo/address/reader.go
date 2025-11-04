package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

func (r *address) GetByID(ctx context.Context, id int64) (*models.Address, error) {
	const query = `
		SELECT 
			id, user_id, client_id, supplier_id,
			street, street_number, complement, city, state, country, postal_code,
			is_active, created_at, updated_at
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
		&address.StreetNumber,
		&address.Complement,
		&address.City,
		&address.State,
		&address.Country,
		&address.PostalCode,
		&address.IsActive,
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

func (r *address) GetByUserID(ctx context.Context, userID int64) ([]*models.Address, error) {
	const query = `
		SELECT 
			id, user_id, client_id, supplier_id,
			street, street_number, complement, city, state, country, postal_code,
			is_active, created_at, updated_at
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
			&address.StreetNumber,
			&address.Complement,
			&address.City,
			&address.State,
			&address.Country,
			&address.PostalCode,
			&address.IsActive,
			&address.CreatedAt,
			&address.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		addresses = append(addresses, &address)
	}

	return addresses, nil
}

func (r *address) GetByClientID(ctx context.Context, clientID int64) ([]*models.Address, error) {
	const query = `
		SELECT 
			id, user_id, client_id, supplier_id,
			street, street_number, complement, city, state, country, postal_code,
			is_active, created_at, updated_at
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
			&address.StreetNumber,
			&address.Complement,
			&address.City,
			&address.State,
			&address.Country,
			&address.PostalCode,
			&address.IsActive,
			&address.CreatedAt,
			&address.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		addresses = append(addresses, &address)
	}

	return addresses, nil
}

func (r *address) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.Address, error) {
	const query = `
		SELECT 
			id, user_id, client_id, supplier_id,
			street, street_number, complement, city, state, country, postal_code,
			is_active, created_at, updated_at
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
			&address.StreetNumber,
			&address.Complement,
			&address.City,
			&address.State,
			&address.Country,
			&address.PostalCode,
			&address.IsActive,
			&address.CreatedAt,
			&address.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		addresses = append(addresses, &address)
	}

	return addresses, nil
}
