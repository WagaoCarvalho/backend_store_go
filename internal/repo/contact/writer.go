package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (r *contactRepo) Create(ctx context.Context, contact *models.Contact) (*models.Contact, error) {
	const query = `
        INSERT INTO contacts (
            contact_name, contact_description,
            email, phone, cell, contact_type, created_at, updated_at
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
        RETURNING id, created_at, updated_at
    `

	now := time.Now()
	contact.CreatedAt = now
	contact.UpdatedAt = now

	err := r.db.QueryRow(ctx, query,
		contact.ContactName,
		contact.ContactDescription,
		contact.Email,
		contact.Phone,
		contact.Cell,
		contact.ContactType,
		contact.CreatedAt,
		contact.UpdatedAt,
	).Scan(&contact.ID, &contact.CreatedAt, &contact.UpdatedAt)

	if err != nil {
		if errMsgPg.IsForeignKeyViolation(err) {
			return nil, errMsg.ErrDBInvalidForeignKey
		}

		if ok, constraint := errMsgPg.IsUniqueViolation(err); ok {
			return nil, fmt.Errorf("%w: %s", errMsg.ErrDuplicate, constraint)
		}

		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return contact, nil
}

func (r *contactRepo) Update(ctx context.Context, contact *models.Contact) error {
	const query = `
		UPDATE contacts
		SET contact_name = $1,
		    contact_description = $2,
		    email = $3,
		    phone = $4,
		    cell = $5,
		    contact_type = $6,
		    updated_at = $7
		WHERE id = $8
		RETURNING updated_at
	`

	now := time.Now()
	contact.UpdatedAt = now

	err := r.db.QueryRow(ctx, query,
		contact.ContactName,
		contact.ContactDescription,
		contact.Email,
		contact.Phone,
		contact.Cell,
		contact.ContactType,
		contact.UpdatedAt,
		contact.ID,
	).Scan(&contact.UpdatedAt)

	if err != nil {
		var pgErr *pgconn.PgError
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return errMsg.ErrNotFound
		case errors.As(err, &pgErr):
			switch pgErr.Code {
			case "23505":
				return errMsg.ErrDuplicate
			case "23514":
				return errMsg.ErrInvalidData
			default:
				return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
			}
		default:
			return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
		}
	}

	return nil
}

func (r *contactRepo) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM contacts WHERE id = $1 RETURNING id`

	var deletedID int64
	err := r.db.QueryRow(ctx, query, id).Scan(&deletedID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}
