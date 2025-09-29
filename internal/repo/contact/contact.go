package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContactRepository interface {
	Create(ctx context.Context, contact *models.Contact) (*models.Contact, error)
	CreateTx(ctx context.Context, tx pgx.Tx, contact *models.Contact) (*models.Contact, error)
	GetByID(ctx context.Context, id int64) (*models.Contact, error)
	GetByUserID(ctx context.Context, userID int64) ([]*models.Contact, error)
	GetByClientID(ctx context.Context, clientID int64) ([]*models.Contact, error)
	GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.Contact, error)
	Update(ctx context.Context, contact *models.Contact) error
	Delete(ctx context.Context, id int64) error
}

type contactRepository struct {
	db *pgxpool.Pool
}

func NewContactRepository(db *pgxpool.Pool) ContactRepository {
	return &contactRepository{db: db}
}

func (r *contactRepository) Create(ctx context.Context, contact *models.Contact) (*models.Contact, error) {
	const query = `
		INSERT INTO contacts (
			user_id, client_id, supplier_id, contact_name, contact_position,
			email, phone, cell, contact_type, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query,
		contact.UserID,
		contact.ClientID,
		contact.SupplierID,
		contact.ContactName,
		contact.ContactPosition,
		contact.Email,
		contact.Phone,
		contact.Cell,
		contact.ContactType,
	).Scan(&contact.ID, &contact.CreatedAt, &contact.UpdatedAt)

	if err != nil {
		if errMsgPg.IsForeignKeyViolation(err) {
			return nil, errMsg.ErrInvalidForeignKey
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return contact, nil
}

func (r *contactRepository) CreateTx(ctx context.Context, tx pgx.Tx, contact *models.Contact) (*models.Contact, error) {
	const query = `
		INSERT INTO contacts (
			user_id, client_id, supplier_id, contact_name, contact_position,
			email, phone, cell, contact_type, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := tx.QueryRow(ctx, query,
		contact.UserID,
		contact.ClientID,
		contact.SupplierID,
		contact.ContactName,
		contact.ContactPosition,
		contact.Email,
		contact.Phone,
		contact.Cell,
		contact.ContactType,
	).Scan(&contact.ID, &contact.CreatedAt, &contact.UpdatedAt)

	if err != nil {
		if errMsgPg.IsForeignKeyViolation(err) {
			return nil, errMsg.ErrInvalidForeignKey
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return contact, nil
}

func (r *contactRepository) GetByID(ctx context.Context, id int64) (*models.Contact, error) {
	const query = `
		SELECT 
			id, user_id, client_id, supplier_id, contact_name, contact_position,
			email, phone, cell, contact_type, created_at, updated_at
		FROM contacts 
		WHERE id = $1
	`

	var contact models.Contact
	err := r.db.QueryRow(ctx, query, id).Scan(
		&contact.ID,
		&contact.UserID,
		&contact.ClientID,
		&contact.SupplierID,
		&contact.ContactName,
		&contact.ContactPosition,
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

func (r *contactRepository) GetByUserID(ctx context.Context, userID int64) ([]*models.Contact, error) {
	const query = `
		SELECT 
			id, user_id, client_id, supplier_id, contact_name, contact_position,
			email, phone, cell, contact_type, created_at, updated_at
		FROM contacts 
		WHERE user_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var contacts []*models.Contact
	for rows.Next() {
		var contact models.Contact
		if err := rows.Scan(
			&contact.ID,
			&contact.UserID,
			&contact.ClientID,
			&contact.SupplierID,
			&contact.ContactName,
			&contact.ContactPosition,
			&contact.Email,
			&contact.Phone,
			&contact.Cell,
			&contact.ContactType,
			&contact.CreatedAt,
			&contact.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		contacts = append(contacts, &contact)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return contacts, nil
}

func (r *contactRepository) GetByClientID(ctx context.Context, clientID int64) ([]*models.Contact, error) {
	const query = `
		SELECT 
			id, user_id, client_id, supplier_id, contact_name, contact_position,
			email, phone, cell, contact_type, created_at, updated_at
		FROM contacts 
		WHERE client_id = $1
	`

	rows, err := r.db.Query(ctx, query, clientID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var contacts []*models.Contact
	for rows.Next() {
		var contact models.Contact
		if err := rows.Scan(
			&contact.ID,
			&contact.UserID,
			&contact.ClientID,
			&contact.SupplierID,
			&contact.ContactName,
			&contact.ContactPosition,
			&contact.Email,
			&contact.Phone,
			&contact.Cell,
			&contact.ContactType,
			&contact.CreatedAt,
			&contact.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		contacts = append(contacts, &contact)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return contacts, nil
}

func (r *contactRepository) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.Contact, error) {
	const query = `
		SELECT 
			id, user_id, client_id, supplier_id, contact_name, contact_position,
			email, phone, cell, contact_type, created_at, updated_at
		FROM contacts 
		WHERE supplier_id = $1
	`

	rows, err := r.db.Query(ctx, query, supplierID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var contacts []*models.Contact
	for rows.Next() {
		var contact models.Contact
		if err := rows.Scan(
			&contact.ID,
			&contact.UserID,
			&contact.ClientID,
			&contact.SupplierID,
			&contact.ContactName,
			&contact.ContactPosition,
			&contact.Email,
			&contact.Phone,
			&contact.Cell,
			&contact.ContactType,
			&contact.CreatedAt,
			&contact.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		contacts = append(contacts, &contact)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return contacts, nil
}

func (r *contactRepository) Update(ctx context.Context, contact *models.Contact) error {
	const query = `
		UPDATE contacts
		SET
			user_id          = $1,
			client_id        = $2,
			supplier_id      = $3,
			contact_name     = $4,
			contact_position = $5,
			email            = $6,
			phone            = $7,
			cell             = $8,
			contact_type     = $9,
			updated_at       = NOW()
		WHERE id = $10
		RETURNING updated_at
	`

	err := r.db.QueryRow(ctx, query,
		contact.UserID,
		contact.ClientID,
		contact.SupplierID,
		contact.ContactName,
		contact.ContactPosition,
		contact.Email,
		contact.Phone,
		contact.Cell,
		contact.ContactType,
		contact.ID,
	).Scan(&contact.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errMsg.ErrNotFound
		}
		if errMsgPg.IsForeignKeyViolation(err) {
			return errMsg.ErrInvalidForeignKey
		}
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (r *contactRepository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM contacts WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	if result.RowsAffected() == 0 {
		return errMsg.ErrNotFound
	}

	return nil
}
