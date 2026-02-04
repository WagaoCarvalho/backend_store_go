package repo

import (
	"context"
	"fmt"

	ifaceTx "github.com/WagaoCarvalho/backend_store_go/internal/iface/contact"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

type contactTx struct{}

func NewContactTx() ifaceTx.ContactTx {
	return &contactTx{}
}

func (r *contactTx) CreateTx(
	ctx context.Context,
	tx pgx.Tx,
	contact *models.Contact,
) (*models.Contact, error) {
	if contact == nil {
		return nil, errMsg.ErrNilModel
	}

	const query = `
		INSERT INTO contacts (
			contact_name,
			contact_description,
			email,
			phone,
			cell,
			contact_type,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := tx.QueryRow(ctx, query,
		contact.ContactName,
		contact.ContactDescription,
		contact.Email,
		contact.Phone,
		contact.Cell,
		contact.ContactType,
	).Scan(
		&contact.ID,
		&contact.CreatedAt,
		&contact.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("%w: %w", errMsg.ErrCreate, err)
	}

	return contact, nil
}
