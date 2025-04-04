package repositories

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContactRepository interface {
	CreateContact(ctx context.Context, contact *models.Contact) error
	GetContactByID(ctx context.Context, id int64) (*models.Contact, error)
	GetContactByUserID(ctx context.Context, userID int64) ([]*models.Contact, error)
	GetContactByClientID(ctx context.Context, clientID int64) ([]*models.Contact, error)
	GetContactBySupplierID(ctx context.Context, supplierID int64) ([]*models.Contact, error)
	Updatecontac(ctx context.Context, contact *models.Contact) error
	Deletecontact(ctx context.Context, id int64) error
}

type contactRepository struct {
	db *pgxpool.Pool
}

func NewContactRepository(db *pgxpool.Pool) ContactRepository {
	return &contactRepository{db: db}
}

func (r *contactRepository) CreateContact(ctx context.Context, contact *models.Contact) error {
	query := `
		INSERT INTO contacts (
			user_id, client_id, supplier_id, contact_name, contact_position,
			email, phone, cell, contact_type
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
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
		return fmt.Errorf("erro ao criar contato: %w", err)
	}

	return nil
}

func (r *contactRepository) GetContactByID(ctx context.Context, id int64) (*models.Contact, error) {
	query := `
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
		return nil, fmt.Errorf("erro ao buscar contato: %w", err)
	}

	return &contact, nil
}

func (r *contactRepository) GetContactByUserID(ctx context.Context, userID int64) ([]*models.Contact, error) {
	query := `
		SELECT 
			id, user_id, client_id, supplier_id, contact_name, contact_position,
			email, phone, cell, contact_type, created_at, updated_at
		FROM contacts 
		WHERE user_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar contatos por user_id: %w", err)
	}
	defer rows.Close()

	var contacts []*models.Contact
	for rows.Next() {
		var contact models.Contact
		err := rows.Scan(
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
			return nil, fmt.Errorf("erro ao escanear contato: %w", err)
		}
		contacts = append(contacts, &contact)
	}

	return contacts, nil
}

func (r *contactRepository) GetContactByClientID(ctx context.Context, clientID int64) ([]*models.Contact, error) {
	query := `
		SELECT 
			id, user_id, client_id, supplier_id, contact_name, contact_position,
			email, phone, cell, contact_type, created_at, updated_at
		FROM contacts 
		WHERE client_id = $1
	`

	rows, err := r.db.Query(ctx, query, clientID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar contatos por client_id: %w", err)
	}
	defer rows.Close()

	var contacts []*models.Contact
	for rows.Next() {
		var contact models.Contact
		err := rows.Scan(
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
			return nil, fmt.Errorf("erro ao escanear contato: %w", err)
		}
		contacts = append(contacts, &contact)
	}

	return contacts, nil
}

func (r *contactRepository) GetContactBySupplierID(ctx context.Context, supplierID int64) ([]*models.Contact, error) {
	query := `
		SELECT 
			id, user_id, client_id, supplier_id, contact_name, contact_position,
			email, phone, cell, contact_type, created_at, updated_at
		FROM contacts 
		WHERE supplier_id = $1
	`

	rows, err := r.db.Query(ctx, query, supplierID)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar contatos por supplier_id: %w", err)
	}
	defer rows.Close()

	var contacts []*models.Contact
	for rows.Next() {
		var contact models.Contact
		err := rows.Scan(
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
			return nil, fmt.Errorf("erro ao escanear contato: %w", err)
		}
		contacts = append(contacts, &contact)
	}

	return contacts, nil
}

func (r *contactRepository) Updatecontac(ctx context.Context, contact *models.Contact) error {
	query := `
		UPDATE contacts SET
			user_id = $1,
			client_id = $2,
			supplier_id = $3,
			contact_name = $4,
			contact_position = $5,
			email = $6,
			phone = $7,
			cell = $8,
			contact_type = $9,
			updated_at = NOW()
		WHERE id = $10
	`

	result, err := r.db.Exec(ctx, query,
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
	)

	if err != nil {
		return fmt.Errorf("erro ao atualizar contato: %w", err)
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return ErrContactNotFound
	}

	return nil
}

func (r *contactRepository) Deletecontact(ctx context.Context, id int64) error {
	query := `DELETE FROM contacts WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("erro ao deletar contato: %w", err)
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		return ErrContactNotFound
	}

	return nil
}

var ErrContactNotFound = errors.New("contato não encontrado")
