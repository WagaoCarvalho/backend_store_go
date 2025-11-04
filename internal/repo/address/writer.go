package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

func (r *address) Create(ctx context.Context, address *models.Address) (*models.Address, error) {
	const query = `
		INSERT INTO addresses (
			user_id, client_id, supplier_id,
			street, street_number, complement, city, state, country, postal_code,
			is_active, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
		RETURNING id, created_at, updated_at;
	`

	err := r.db.QueryRow(ctx, query,
		address.UserID,
		address.ClientID,
		address.SupplierID,
		address.Street,
		address.StreetNumber,
		address.Complement,
		address.City,
		address.State,
		address.Country,
		address.PostalCode,
		address.IsActive,
	).Scan(&address.ID, &address.CreatedAt, &address.UpdatedAt)

	if err != nil {
		if errMsgPg.IsForeignKeyViolation(err) {
			return nil, errMsg.ErrDBInvalidForeignKey
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return address, nil
}

func (r *address) Update(ctx context.Context, address *models.Address) error {
	const query = `
		UPDATE addresses
		SET 
			user_id       = $1,
			client_id     = $2,
			supplier_id   = $3,
			street        = $4,
			street_number = $5,
			complement    = $6,
			city          = $7,
			state         = $8,
			country       = $9,
			postal_code   = $10,
			is_active     = $11,
			updated_at    = NOW()
		WHERE id = $12
		RETURNING updated_at;
	`

	err := r.db.QueryRow(ctx, query,
		address.UserID,
		address.ClientID,
		address.SupplierID,
		address.Street,
		address.StreetNumber,
		address.Complement,
		address.City,
		address.State,
		address.Country,
		address.PostalCode,
		address.IsActive,
		address.ID,
	).Scan(&address.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errMsg.ErrNotFound
		}
		if errMsgPg.IsForeignKeyViolation(err) {
			return errMsg.ErrDBInvalidForeignKey
		}
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (r *address) Delete(ctx context.Context, id int64) error {
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
