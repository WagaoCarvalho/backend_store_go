package repositories

import (
	"context"
	"errors"
	"fmt"

	logger "github.com/WagaoCarvalho/backend_store_go/internal/logger"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ContactRepository interface {
	Create(ctx context.Context, contact *models.Contact) (*models.Contact, error)
	GetByID(ctx context.Context, id int64) (*models.Contact, error)
	GetByUserID(ctx context.Context, userID int64) ([]*models.Contact, error)
	GetByClientID(ctx context.Context, clientID int64) ([]*models.Contact, error)
	GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.Contact, error)
	Update(ctx context.Context, contact *models.Contact) error
	Delete(ctx context.Context, id int64) error
}

type contactRepository struct {
	db     *pgxpool.Pool
	logger *logger.LoggerAdapter
}

func NewContactRepository(db *pgxpool.Pool, logger *logger.LoggerAdapter) ContactRepository {
	return &contactRepository{db: db, logger: logger}
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

	r.logger.Info(ctx, "Contato criado com sucesso", map[string]interface{}{
		"contact_id":  contact.ID,
		"user_id":     utils.Int64OrNil(contact.UserID),
		"client_id":   utils.Int64OrNil(contact.ClientID),
		"supplier_id": utils.Int64OrNil(contact.SupplierID),
		"email":       contact.Email,
	})

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
			r.logger.Error(ctx, ErrContactNotFound, "Contato não encontrado", map[string]interface{}{
				"contact_id": id,
			})
			return nil, ErrContactNotFound
		}

		r.logger.Error(ctx, err, "Erro ao buscar contato por ID", map[string]interface{}{
			"contact_id": id,
		})
		return nil, fmt.Errorf("%w: %v", ErrFetchContact, err)
	}

	r.logger.Info(ctx, "Contato encontrado com sucesso", map[string]interface{}{
		"contact_id":  contact.ID,
		"user_id":     utils.Int64OrNil(contact.UserID),
		"client_id":   utils.Int64OrNil(contact.ClientID),
		"supplier_id": utils.Int64OrNil(contact.SupplierID),
		"email":       contact.Email,
	})

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
		r.logger.Error(ctx, err, "Erro ao buscar contatos por user_id", map[string]interface{}{
			"user_id": userID,
		})
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
			r.logger.Error(ctx, err, "Erro ao fazer scan do contato", map[string]interface{}{
				"user_id": userID,
			})
			return nil, fmt.Errorf("%w: %v", ErrScanContact, err)
		}
		contacts = append(contacts, &contact)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(ctx, err, "Erro ao iterar sobre os contatos", map[string]interface{}{
			"user_id": userID,
		})
		return nil, fmt.Errorf("%w: %v", ErrFetchContactsByUser, err)
	}

	r.logger.Info(ctx, "Contatos buscados com sucesso", map[string]interface{}{
		"user_id":      userID,
		"total_result": len(contacts),
	})

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
		r.logger.Error(ctx, err, "Erro ao buscar contatos por client_id", map[string]interface{}{
			"client_id": clientID,
		})
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
			r.logger.Error(ctx, err, "Erro ao fazer scan do contato por client_id", map[string]interface{}{
				"client_id": clientID,
			})
			return nil, fmt.Errorf("%w: %v", ErrScanContact, err)
		}
		contacts = append(contacts, &contact)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(ctx, err, "Erro ao iterar sobre os contatos por client_id", map[string]interface{}{
			"client_id": clientID,
		})
		return nil, fmt.Errorf("%w: %v", ErrFetchContactsByClient, err)
	}

	r.logger.Info(ctx, "Contatos buscados com sucesso por client_id", map[string]interface{}{
		"client_id":    clientID,
		"total_result": len(contacts),
	})

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
		r.logger.Error(ctx, err, "Erro ao buscar contatos por supplier_id", map[string]interface{}{
			"supplier_id": supplierID,
		})
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
			r.logger.Error(ctx, err, "Erro ao fazer scan do contato por supplier_id", map[string]interface{}{
				"supplier_id": supplierID,
			})
			return nil, fmt.Errorf("%w: %v", ErrScanContact, err)
		}
		contacts = append(contacts, &contact)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(ctx, err, "Erro ao iterar sobre os contatos por supplier_id", map[string]interface{}{
			"supplier_id": supplierID,
		})
		return nil, fmt.Errorf("%w: %v", ErrFetchContactsBySupplier, err)
	}

	r.logger.Info(ctx, "Contatos buscados com sucesso por supplier_id", map[string]interface{}{
		"supplier_id":  supplierID,
		"total_result": len(contacts),
	})

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
		WHERE
			id = $10
	`

	cmdTag, err := r.db.Exec(ctx, query,
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
		r.logger.Error(ctx, err, "Erro ao atualizar contato", map[string]interface{}{
			"contact_id":  contact.ID,
			"user_id":     utils.Int64OrNil(contact.UserID),
			"client_id":   utils.Int64OrNil(contact.ClientID),
			"supplier_id": utils.Int64OrNil(contact.SupplierID),
			"email":       contact.Email,
		})
		return fmt.Errorf("%w: %v", ErrUpdateContact, err)
	}

	if cmdTag.RowsAffected() == 0 {
		r.logger.Info(ctx, "Contato não encontrado para atualização", map[string]interface{}{
			"contact_id": contact.ID,
		})
		return ErrContactNotFound
	}

	r.logger.Info(ctx, "Contato atualizado com sucesso", map[string]interface{}{
		"contact_id": contact.ID,
	})

	return nil
}

func (r *contactRepository) Delete(ctx context.Context, id int64) error {
	const query = `DELETE FROM contacts WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error(ctx, err, "Erro ao deletar contato", map[string]interface{}{
			"contact_id": id,
		})
		return fmt.Errorf("%w: %v", ErrDeleteContact, err)
	}

	if rowsAffected := result.RowsAffected(); rowsAffected == 0 {
		r.logger.Info(ctx, "Contato não encontrado para exclusão", map[string]interface{}{
			"contact_id": id,
		})
		return ErrContactNotFound
	}

	r.logger.Info(ctx, "Contato deletado com sucesso", map[string]interface{}{
		"contact_id": id,
	})

	return nil
}
