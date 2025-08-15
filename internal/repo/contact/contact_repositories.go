package repositories

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
	logger "github.com/WagaoCarvalho/backend_store_go/pkg/logger"

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
	db     *pgxpool.Pool
	logger logger.LoggerAdapterInterface
}

func NewContactRepository(db *pgxpool.Pool, logger logger.LoggerAdapterInterface) ContactRepository {
	return &contactRepository{db: db, logger: logger}
}

func (r *contactRepository) Create(ctx context.Context, contact *models.Contact) (*models.Contact, error) {
	ref := "[contactRepository - Create] - "

	r.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{
		"contact_name": contact.ContactName,
		"user_id":      utils.Int64OrNil(contact.UserID),
		"client_id":    utils.Int64OrNil(contact.ClientID),
		"supplier_id":  utils.Int64OrNil(contact.SupplierID),
	})

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
		if IsForeignKeyViolation(err) {
			r.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{
				"user_id":     utils.Int64OrNil(contact.UserID),
				"client_id":   utils.Int64OrNil(contact.ClientID),
				"supplier_id": utils.Int64OrNil(contact.SupplierID),
			})
			return nil, ErrInvalidForeignKey
		}

		r.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"user_id":      utils.Int64OrNil(contact.UserID),
			"client_id":    utils.Int64OrNil(contact.ClientID),
			"supplier_id":  utils.Int64OrNil(contact.SupplierID),
			"contact_name": contact.ContactName,
			"email":        contact.Email,
		})
		return nil, fmt.Errorf("%w: %v", ErrCreateContact, err)
	}

	r.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"contact_id":  contact.ID,
		"user_id":     utils.Int64OrNil(contact.UserID),
		"client_id":   utils.Int64OrNil(contact.ClientID),
		"supplier_id": utils.Int64OrNil(contact.SupplierID),
		"email":       contact.Email,
	})

	return contact, nil
}

func (r *contactRepository) CreateTx(ctx context.Context, tx pgx.Tx, contact *models.Contact) (*models.Contact, error) {
	ref := "[contactRepository - CreateTx] - "

	r.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{
		"contact_name": contact.ContactName,
		"user_id":      utils.Int64OrNil(contact.UserID),
		"client_id":    utils.Int64OrNil(contact.ClientID),
		"supplier_id":  utils.Int64OrNil(contact.SupplierID),
	})

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
		if IsForeignKeyViolation(err) {
			r.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{
				"user_id":     utils.Int64OrNil(contact.UserID),
				"client_id":   utils.Int64OrNil(contact.ClientID),
				"supplier_id": utils.Int64OrNil(contact.SupplierID),
			})
			return nil, ErrInvalidForeignKey
		}

		r.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"user_id":      utils.Int64OrNil(contact.UserID),
			"client_id":    utils.Int64OrNil(contact.ClientID),
			"supplier_id":  utils.Int64OrNil(contact.SupplierID),
			"contact_name": contact.ContactName,
			"email":        contact.Email,
		})
		return nil, fmt.Errorf("%w: %v", ErrCreateContact, err)
	}

	r.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"contact_id":  contact.ID,
		"user_id":     utils.Int64OrNil(contact.UserID),
		"client_id":   utils.Int64OrNil(contact.ClientID),
		"supplier_id": utils.Int64OrNil(contact.SupplierID),
		"email":       contact.Email,
	})

	return contact, nil
}

func (r *contactRepository) GetByID(ctx context.Context, id int64) (*models.Contact, error) {
	ref := "[contactRepository - GetByID] - "

	r.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"contact_id": id,
	})

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
			r.logger.Info(ctx, ref+logger.LogNotFound, map[string]any{
				"contact_id": id,
			})
			return nil, ErrContactNotFound
		}

		r.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"contact_id": id,
		})
		return nil, fmt.Errorf("%w: %v", ErrFetchContact, err)
	}

	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"contact_id":  contact.ID,
		"user_id":     utils.Int64OrNil(contact.UserID),
		"client_id":   utils.Int64OrNil(contact.ClientID),
		"supplier_id": utils.Int64OrNil(contact.SupplierID),
		"email":       contact.Email,
	})

	return &contact, nil
}

func (r *contactRepository) GetByUserID(ctx context.Context, userID int64) ([]*models.Contact, error) {
	ref := "[contactRepository - GetByUserID] - "

	r.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"user_id": userID,
	})

	const query = `
		SELECT 
			id, user_id, client_id, supplier_id, contact_name, contact_position,
			email, phone, cell, contact_type, created_at, updated_at
		FROM contacts 
		WHERE user_id = $1
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
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
			r.logger.Error(ctx, err, ref+logger.LogGetErrorScan, map[string]any{
				"user_id": userID,
			})
			return nil, fmt.Errorf("%w: %v", ErrScanContact, err)
		}
		contacts = append(contacts, &contact)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"user_id": userID,
		})
		return nil, fmt.Errorf("%w: %v", ErrFetchContactsByUser, err)
	}

	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"user_id":       userID,
		"total_results": len(contacts),
	})

	return contacts, nil
}

func (r *contactRepository) GetByClientID(ctx context.Context, clientID int64) ([]*models.Contact, error) {
	ref := "[contactRepository - GetByClientID] - "

	r.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"client_id": clientID,
	})

	const query = `
		SELECT 
			id, user_id, client_id, supplier_id, contact_name, contact_position,
			email, phone, cell, contact_type, created_at, updated_at
		FROM contacts 
		WHERE client_id = $1
	`

	rows, err := r.db.Query(ctx, query, clientID)
	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
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
			r.logger.Error(ctx, err, ref+logger.LogGetErrorScan, map[string]any{
				"client_id": clientID,
			})
			return nil, fmt.Errorf("%w: %v", ErrScanContact, err)
		}
		contacts = append(contacts, &contact)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"client_id": clientID,
		})
		return nil, fmt.Errorf("%w: %v", ErrFetchContactsByClient, err)
	}

	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"client_id":     clientID,
		"total_results": len(contacts),
	})

	return contacts, nil
}

func (r *contactRepository) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.Contact, error) {
	ref := "[contactRepository - GetBySupplierID] - "

	r.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"supplier_id": supplierID,
	})

	const query = `
		SELECT 
			id, user_id, client_id, supplier_id, contact_name, contact_position,
			email, phone, cell, contact_type, created_at, updated_at
		FROM contacts 
		WHERE supplier_id = $1
	`

	rows, err := r.db.Query(ctx, query, supplierID)
	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
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
			r.logger.Error(ctx, err, ref+logger.LogGetErrorScan, map[string]any{
				"supplier_id": supplierID,
			})
			return nil, fmt.Errorf("%w: %v", ErrScanContact, err)
		}
		contacts = append(contacts, &contact)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"supplier_id": supplierID,
		})
		return nil, fmt.Errorf("%w: %v", ErrFetchContactsBySupplier, err)
	}

	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"supplier_id":   supplierID,
		"total_results": len(contacts),
	})

	return contacts, nil
}

func (r *contactRepository) Update(ctx context.Context, contact *models.Contact) error {
	ref := "[contactRepository - Update] - "

	r.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"contact_id":   contact.ID,
		"contact_name": contact.ContactName,
		"user_id":      utils.Int64OrNil(contact.UserID),
		"client_id":    utils.Int64OrNil(contact.ClientID),
		"supplier_id":  utils.Int64OrNil(contact.SupplierID),
	})

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
		r.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"contact_id":  contact.ID,
			"user_id":     utils.Int64OrNil(contact.UserID),
			"client_id":   utils.Int64OrNil(contact.ClientID),
			"supplier_id": utils.Int64OrNil(contact.SupplierID),
			"email":       contact.Email,
		})
		return fmt.Errorf("%w: %v", ErrUpdateContact, err)
	}

	if cmdTag.RowsAffected() == 0 {
		r.logger.Info(ctx, ref+logger.LogNotFound, map[string]any{
			"contact_id": contact.ID,
		})
		return ErrContactNotFound
	}

	r.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"contact_id": contact.ID,
	})

	return nil
}

func (r *contactRepository) Delete(ctx context.Context, id int64) error {
	ref := "[contactRepository - Delete] - "

	r.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{
		"contact_id": id,
	})

	const query = `DELETE FROM contacts WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"contact_id": id,
		})
		return fmt.Errorf("%w: %v", ErrDeleteContact, err)
	}

	if result.RowsAffected() == 0 {
		r.logger.Info(ctx, ref+logger.LogNotFound, map[string]any{
			"contact_id": id,
		})
		return ErrContactNotFound
	}

	r.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"contact_id": id,
	})

	return nil
}
