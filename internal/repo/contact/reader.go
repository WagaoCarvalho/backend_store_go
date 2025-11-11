package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

func (r *contactRepo) GetByID(ctx context.Context, id int64) (*models.Contact, error) {
	const query = `
		SELECT 
			id, contact_name, contact_description,
			email, phone, cell, contact_type, created_at, updated_at
		FROM contacts 
		WHERE id = $1
	`

	var contact models.Contact
	err := r.db.QueryRow(ctx, query, id).Scan(
		&contact.ID,
		&contact.ContactName,
		&contact.ContactDescription,
		&contact.Email,
		&contact.Phone,
		&contact.Cell,
		&contact.ContactType,
		&contact.CreatedAt,
		&contact.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errMsg.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return &contact, nil
}
