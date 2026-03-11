package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address/address"
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

func (r *addressRepo) GetByID(ctx context.Context, id int64) (*models.Address, error) {
	const query = baseSelectAddress + ` WHERE id = $1`

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
