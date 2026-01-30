package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"

	"github.com/jackc/pgx/v5"
)

const baseSelectAddress = `
	SELECT 
		id, user_id, client_cpf_id, supplier_id,
		street, street_number, complement, city, state, country, postal_code,
		is_active, created_at, updated_at
	FROM addresses
`

// =======================
// Public methods
// =======================

func (r *addressRepo) GetByID(ctx context.Context, id int64) (*models.Address, error) {
	const query = baseSelectAddress + `
		WHERE id = $1;
	`

	row := r.db.QueryRow(ctx, query, id)

	addr, err := scanAddress(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errMsg.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return addr, nil
}

func (r *addressRepo) GetByUserID(ctx context.Context, userID int64) ([]*models.Address, error) {
	return r.getByField(ctx, "user_id", userID)
}

func (r *addressRepo) GetByClientCpfID(ctx context.Context, clientCpfID int64) ([]*models.Address, error) {
	return r.getByField(ctx, "client_cpf_id", clientCpfID)
}

func (r *addressRepo) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.Address, error) {
	return r.getByField(ctx, "supplier_id", supplierID)
}

// =======================
// Internal helpers
// =======================

func (r *addressRepo) getByField(ctx context.Context, field string, value any) ([]*models.Address, error) {
	if !isAllowedAddressField(field) {
		return nil, errMsg.ErrInvalidField
	}

	query := fmt.Sprintf(`
		%s
		WHERE %s = $1;
	`, baseSelectAddress, field)

	rows, err := r.db.Query(ctx, query, value)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var results []*models.Address

	for rows.Next() {
		addr, err := scanAddress(rows)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		results = append(results, addr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return results, nil
}

func isAllowedAddressField(field string) bool {
	switch field {
	case "user_id", "client_cpf_id", "supplier_id":
		return true
	default:
		return false
	}
}

// =======================
// Scanner
// =======================

type scanner interface {
	Scan(dest ...any) error
}

func scanAddress(s scanner) (*models.Address, error) {
	var addr models.Address

	err := s.Scan(
		&addr.ID,
		&addr.UserID,
		&addr.ClientCpfID,
		&addr.SupplierID,
		&addr.Street,
		&addr.StreetNumber,
		&addr.Complement,
		&addr.City,
		&addr.State,
		&addr.Country,
		&addr.PostalCode,
		&addr.IsActive,
		&addr.CreatedAt,
		&addr.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &addr, nil
}
