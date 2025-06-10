package repositories

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrContactNotFound         = errors.New("contato não encontrado")
	ErrCreateContact           = errors.New("erro ao criar contato")
	ErrFetchContact            = errors.New("erro ao buscar contato")
	ErrFetchContactsByUser     = errors.New("erro ao buscar contatos por user_id")
	ErrFetchContactsByClient   = errors.New("erro ao buscar contatos por client_id")
	ErrFetchContactsBySupplier = errors.New("erro ao buscar contatos por supplier_id")
	ErrScanContact             = errors.New("erro ao escanear contato")
	ErrUpdateContact           = errors.New("erro ao atualizar contato")
	ErrDeleteContact           = errors.New("erro ao deletar contato")
	ErrVersionConflict         = errors.New("conflito de versão: o endereço foi modificado por outra operação")
)

type ContactRepository interface {
	Create(ctx context.Context, contact *models.Contact) (*models.Contact, error)
	GetByID(ctx context.Context, id int64) (*models.Contact, error)
	GetByUserID(ctx context.Context, userID int64) ([]*models.Contact, error)
	GetByClientID(ctx context.Context, clientID int64) ([]*models.Contact, error)
	GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.Contact, error)
	GetVersionByID(ctx context.Context, id int64) (int, error)
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
		return nil, fmt.Errorf("%w: %v", ErrCreateContact, err)
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
			return nil, ErrContactNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrFetchContact, err)
	}

	return &contact, nil
}

func (r *contactRepository) GetVersionByID(ctx context.Context, id int64) (int, error) {
	const query = `
		SELECT version
		FROM contacts
		WHERE id = $1;
	`

	var version int
	err := r.db.QueryRow(ctx, query, id).Scan(&version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrContactNotFound
		}
		return 0, fmt.Errorf("%w: %v", ErrFetchContact, err)
	}

	return version, nil
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
		return nil, fmt.Errorf("%w: %v", ErrFetchContactsByUser, err)
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
			return nil, fmt.Errorf("%w: %v", ErrScanContact, err)
		}
		contacts = append(contacts, &contact)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFetchContactsByUser, err)
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
		return nil, fmt.Errorf("%w: %v", ErrFetchContactsByClient, err)
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
			return nil, fmt.Errorf("%w: %v", ErrScanContact, err)
		}
		contacts = append(contacts, &contact)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFetchContactsByClient, err)
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
		return nil, fmt.Errorf("%w: %v", ErrFetchContactsBySupplier, err)
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
			return nil, fmt.Errorf("%w: %v", ErrScanContact, err)
		}
		contacts = append(contacts, &contact)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrFetchContactsBySupplier, err)
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
			version          = version + 1,
			updated_at       = NOW()
		WHERE
			id      = $10
			AND version = $11
		RETURNING
			version,
			updated_at;
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
		contact.Version,
	).Scan(&contact.Version, &contact.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {

			return ErrVersionConflict
		}
		return fmt.Errorf("%w: %v", ErrUpdateContact, err)
	}

	return nil
}

func (r *contactRepository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM contacts WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteContact, err)
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return ErrContactNotFound
	}

	return nil
}
