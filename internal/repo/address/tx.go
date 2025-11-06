package repo

import (
	"context"
	"fmt"

	ifaceTx "github.com/WagaoCarvalho/backend_store_go/internal/iface/address"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type addressTx struct {
	db *pgxpool.Pool
}

func NewAddressTx(db *pgxpool.Pool) ifaceTx.AddressTx {
	return &addressTx{db: db}
}

func (r *addressTx) CreateTx(ctx context.Context, tx pgx.Tx, address *models.Address) (*models.Address, error) {
	const query = `
		INSERT INTO addresses (
			user_id, client_id, supplier_id,
			street, street_number, complement, city, state, country, postal_code,
			is_active, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())
		RETURNING id, created_at, updated_at;
	`

	err := tx.QueryRow(ctx, query,
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
