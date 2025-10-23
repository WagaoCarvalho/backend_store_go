package repo

import (
	"context"
	"errors"
	"fmt"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Contact interface {
	Create(ctx context.Context, contact *models.Contact) (*models.Contact, error)
	CreateTx(ctx context.Context, tx pgx.Tx, contact *models.Contact) (*models.Contact, error)
	GetByID(ctx context.Context, id int64) (*models.Contact, error)
	Update(ctx context.Context, contact *models.Contact) error
	Delete(ctx context.Context, id int64) error
}

type contact struct {
	db *pgxpool.Pool
}

func NewContact(db *pgxpool.Pool) Contact {
	return &contact{db: db}
}

func (r *contact) Create(ctx context.Context, contact *models.Contact) (*models.Contact, error) {
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505":
				return nil, errMsg.ErrDuplicate
			case "23514":
				return nil, errMsg.ErrInvalidData
			default:
				return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
			}
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return contact, nil
}

func (r *contact) CreateTx(ctx context.Context, tx pgx.Tx, contact *models.Contact) (*models.Contact, error) {
	const query = `
		INSERT INTO contacts (
			contact_name, contact_description,
			email, phone, cell, contact_type, created_at, updated_at
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
	).Scan(&contact.ID, &contact.CreatedAt, &contact.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return contact, nil
}

func (r *contact) GetByID(ctx context.Context, id int64) (*models.Contact, error) {
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

func (r *contact) Update(ctx context.Context, contact *models.Contact) error {
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

func (r *contact) Delete(ctx context.Context, id int64) error {
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
