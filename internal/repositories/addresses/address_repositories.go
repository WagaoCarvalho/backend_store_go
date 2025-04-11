package repositories

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// AddressRepository define a interface para o repositório de endereços.
type AddressRepository interface {
	Create(ctx context.Context, address models.Address) (models.Address, error)
	GetByID(ctx context.Context, id int) (models.Address, error)
	Update(ctx context.Context, address models.Address) error
	Delete(ctx context.Context, id int) error
}

// addressRepository é a implementação da interface AddressRepository.
type addressRepository struct {
	db *pgxpool.Pool
}

// NewAddressRepository cria uma nova instância de AddressRepository.
func NewAddressRepository(db *pgxpool.Pool) AddressRepository {
	return &addressRepository{db: db}
}

// Create insere um novo endereço no banco de dados.
func (r *addressRepository) Create(ctx context.Context, address models.Address) (models.Address, error) {
	query := `
		INSERT INTO addresses (user_id, client_id, supplier_id, street, city, state, country, postal_code, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRow(ctx, query,
		address.UserID, address.ClientID, address.SupplierID,
		address.Street, address.City, address.State,
		address.Country, address.PostalCode,
	).Scan(&address.ID, &address.CreatedAt, &address.UpdatedAt)

	if err != nil {
		return models.Address{}, fmt.Errorf("erro ao criar endereço: %w", err)
	}

	return address, nil
}

// GetByID retorna um endereço pelo ID.
func (r *addressRepository) GetByID(ctx context.Context, id int) (models.Address, error) {
	query := `
		SELECT id, user_id, client_id, supplier_id, street, city, state, country, postal_code, created_at, updated_at
		FROM addresses WHERE id = $1
	`
	var address models.Address
	err := r.db.QueryRow(ctx, query, id).Scan(
		&address.ID, &address.UserID, &address.ClientID, &address.SupplierID,
		&address.Street, &address.City, &address.State,
		&address.Country, &address.PostalCode,
		&address.CreatedAt, &address.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Address{}, fmt.Errorf("endereço não encontrado")
		}
		return models.Address{}, fmt.Errorf("erro ao buscar endereço: %w", err)
	}

	return address, nil
}

// Update atualiza um endereço existente.
func (r *addressRepository) Update(ctx context.Context, address models.Address) error {
	query := `
		UPDATE addresses
		SET user_id = $1, client_id = $2, supplier_id = $3, street = $4, city = $5, state = $6, 
			country = $7, postal_code = $8, updated_at = NOW()
		WHERE id = $9
	`
	ct, err := r.db.Exec(ctx, query,
		address.UserID, address.ClientID, address.SupplierID,
		address.Street, address.City, address.State,
		address.Country, address.PostalCode,
		address.ID,
	)

	if err != nil {
		return fmt.Errorf("erro ao atualizar endereço: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("endereço não encontrado")
	}

	return nil
}

// Delete remove um endereço pelo ID.
func (r *addressRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM addresses WHERE id = $1`
	ct, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("erro ao excluir endereço: %w", err)
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("endereço não encontrado")
	}

	return nil
}
