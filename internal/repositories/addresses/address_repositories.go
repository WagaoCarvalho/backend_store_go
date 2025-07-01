package repositories

import (
	"context"
	"errors"
	"fmt"

	logger "github.com/WagaoCarvalho/backend_store_go/internal/logger"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AddressRepository interface {
	Create(ctx context.Context, address *models.Address) (*models.Address, error)
	GetByID(ctx context.Context, id int64) (*models.Address, error)
	GetByUserID(ctx context.Context, userID int64) ([]*models.Address, error)
	GetByClientID(ctx context.Context, clientID int64) ([]*models.Address, error)
	GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.Address, error)
	Update(ctx context.Context, address *models.Address) error
	Delete(ctx context.Context, id int64) error
}

type addressRepository struct {
	db     *pgxpool.Pool
	logger *logger.LoggerAdapter
}

func NewAddressRepository(db *pgxpool.Pool, logger *logger.LoggerAdapter) AddressRepository {
	return &addressRepository{db: db, logger: logger}
}

func (r *addressRepository) Create(ctx context.Context, address *models.Address) (*models.Address, error) {
	r.logger.Info(ctx, "[addressRepository] - Iniciando criação de endereço", map[string]interface{}{
		"user_id":     utils.Int64OrNil(address.UserID),
		"client_id":   utils.Int64OrNil(address.ClientID),
		"supplier_id": utils.Int64OrNil(address.SupplierID),
		"street":      address.Street,
	})
	const query = `
		INSERT INTO addresses (
			user_id, client_id, supplier_id,
			street, city, state, country, postal_code,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		RETURNING id, created_at, updated_at;
	`

	err := r.db.QueryRow(ctx, query,
		address.UserID,
		address.ClientID,
		address.SupplierID,
		address.Street,
		address.City,
		address.State,
		address.Country,
		address.PostalCode,
	).Scan(&address.ID, &address.CreatedAt, &address.UpdatedAt)

	if err != nil {
		if IsForeignKeyViolation(err) {
			r.logger.Warn(ctx, "[addressRepository] - Falha por chave estrangeira ao criar endereço", map[string]interface{}{
				"user_id":     utils.Int64OrNil(address.UserID),
				"client_id":   utils.Int64OrNil(address.ClientID),
				"supplier_id": utils.Int64OrNil(address.SupplierID),
			})
			return nil, ErrInvalidForeignKey
		}

		r.logger.Error(ctx, err, "[addressRepository] - Erro ao criar endereço", map[string]interface{}{
			"user_id":     utils.Int64OrNil(address.UserID),
			"client_id":   utils.Int64OrNil(address.ClientID),
			"supplier_id": utils.Int64OrNil(address.SupplierID),
			"street":      address.Street,
		})
		return nil, fmt.Errorf("%w: %v", ErrCreateAddress, err)
	}

	r.logger.Info(ctx, "[addressRepository] - Endereço criado com sucesso", map[string]interface{}{
		"address_id": address.ID,
		"user_id":    utils.Int64OrNil(address.UserID),
	})

	return address, nil
}

func (r *addressRepository) GetByID(ctx context.Context, id int64) (*models.Address, error) {
	r.logger.Info(ctx, "[addressRepository] - Iniciando busca de endereço", map[string]interface{}{
		"address_id": id,
	})
	const query = `
		SELECT 
			id, user_id, client_id, supplier_id,
			street, city, state, country, postal_code,
			created_at, updated_at
		FROM addresses
		WHERE id = $1;
	`

	var address models.Address
	err := r.db.QueryRow(ctx, query, id).Scan(
		&address.ID,
		&address.UserID,
		&address.ClientID,
		&address.SupplierID,
		&address.Street,
		&address.City,
		&address.State,
		&address.Country,
		&address.PostalCode,
		&address.CreatedAt,
		&address.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Info(ctx, "[addressRepository] - Endereço não encontrado", map[string]interface{}{
				"address_id": id,
			})
			return nil, ErrAddressNotFound
		}

		r.logger.Error(ctx, err, "[addressRepository] - Erro ao buscar endereço", map[string]interface{}{
			"address_id": id,
		})
		return nil, fmt.Errorf("%w: %v", ErrFetchAddress, err)
	}

	r.logger.Info(ctx, "[addressRepository] - Endereço recuperado com sucesso", map[string]interface{}{
		"address_id": address.ID,
		"user_id":    utils.Int64OrNil(address.UserID),
	})

	return &address, nil
}

func (r *addressRepository) GetByUserID(ctx context.Context, userID int64) ([]*models.Address, error) {
	r.logger.Info(ctx, "[addressRepository] - Iniciando busca de endereços por usuário", map[string]interface{}{
		"user_id": userID,
	})
	const query = `
		SELECT 
			id, user_id, client_id, supplier_id,
			street, city, state, country, postal_code,
			created_at, updated_at
		FROM addresses
		WHERE user_id = $1;
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		r.logger.Error(ctx, err, "[addressRepository] - Erro ao buscar endereços por user_id", map[string]interface{}{
			"user_id": userID,
		})
		return nil, fmt.Errorf("%w: %v", ErrFetchAddress, err)
	}
	defer rows.Close()

	var addresses []*models.Address
	for rows.Next() {
		var address models.Address
		if err := rows.Scan(
			&address.ID,
			&address.UserID,
			&address.ClientID,
			&address.SupplierID,
			&address.Street,
			&address.City,
			&address.State,
			&address.Country,
			&address.PostalCode,
			&address.CreatedAt,
			&address.UpdatedAt,
		); err != nil {
			r.logger.Error(ctx, err, "[addressRepository] - Erro ao fazer scan dos endereços", map[string]interface{}{
				"user_id": userID,
			})
			return nil, fmt.Errorf("%w: %v", ErrFetchAddress, err)
		}
		addresses = append(addresses, &address)
	}

	r.logger.Info(ctx, "[addressRepository] - Endereços recuperados com sucesso", map[string]interface{}{
		"user_id":       userID,
		"total_results": len(addresses),
	})

	return addresses, nil
}

func (r *addressRepository) GetByClientID(ctx context.Context, clientID int64) ([]*models.Address, error) {
	r.logger.Info(ctx, "[addressRepository] - Iniciando busca de endereços por cliente", map[string]interface{}{
		"client_id": clientID,
	})

	const query = `
		SELECT 
			id, user_id, client_id, supplier_id,
			street, city, state, country, postal_code,
			created_at, updated_at
		FROM addresses
		WHERE client_id = $1;
	`

	rows, err := r.db.Query(ctx, query, clientID)
	if err != nil {
		r.logger.Error(ctx, err, "[addressRepository] - Erro ao buscar endereços por client_id", map[string]interface{}{
			"client_id": clientID,
		})
		return nil, fmt.Errorf("%w: %v", ErrFetchAddress, err)
	}
	defer rows.Close()

	var addresses []*models.Address
	for rows.Next() {
		var address models.Address
		if err := rows.Scan(
			&address.ID,
			&address.UserID,
			&address.ClientID,
			&address.SupplierID,
			&address.Street,
			&address.City,
			&address.State,
			&address.Country,
			&address.PostalCode,
			&address.CreatedAt,
			&address.UpdatedAt,
		); err != nil {
			r.logger.Error(ctx, err, "[addressRepository] - Erro ao fazer scan dos endereços por client_id", map[string]interface{}{
				"client_id": clientID,
			})
			return nil, fmt.Errorf("%w: %v", ErrFetchAddress, err)
		}
		addresses = append(addresses, &address)
	}

	r.logger.Info(ctx, "[addressRepository] - Endereços recuperados com sucesso por client_id", map[string]interface{}{
		"client_id":     clientID,
		"total_results": len(addresses),
	})

	return addresses, nil
}

func (r *addressRepository) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.Address, error) {
	r.logger.Info(ctx, "[addressRepository] - Iniciando busca de endereços por fornecedor", map[string]interface{}{
		"supplier_id": supplierID,
	})

	const query = `
		SELECT 
			id, user_id, client_id, supplier_id,
			street, city, state, country, postal_code,
			created_at, updated_at
		FROM addresses
		WHERE supplier_id = $1;
	`

	rows, err := r.db.Query(ctx, query, supplierID)
	if err != nil {
		r.logger.Error(ctx, err, "[addressRepository] - Erro ao buscar endereços por supplier_id", map[string]interface{}{
			"supplier_id": supplierID,
		})
		return nil, fmt.Errorf("%w: %v", ErrFetchAddress, err)
	}
	defer rows.Close()

	var addresses []*models.Address
	for rows.Next() {
		var address models.Address
		if err := rows.Scan(
			&address.ID,
			&address.UserID,
			&address.ClientID,
			&address.SupplierID,
			&address.Street,
			&address.City,
			&address.State,
			&address.Country,
			&address.PostalCode,
			&address.CreatedAt,
			&address.UpdatedAt,
		); err != nil {
			r.logger.Error(ctx, err, "[addressRepository] - Erro ao fazer scan dos endereços por supplier_id", map[string]interface{}{
				"supplier_id": supplierID,
			})
			return nil, fmt.Errorf("%w: %v", ErrFetchAddress, err)
		}
		addresses = append(addresses, &address)
	}

	r.logger.Info(ctx, "[addressRepository] - Endereços recuperados com sucesso por supplier_id", map[string]interface{}{
		"supplier_id":   supplierID,
		"total_results": len(addresses),
	})

	return addresses, nil
}

func (r *addressRepository) Update(ctx context.Context, address *models.Address) error {
	r.logger.Info(ctx, "[addressRepository] - Iniciando atualização de endereço", map[string]interface{}{
		"address_id":  address.ID,
		"user_id":     utils.Int64OrNil(address.UserID),
		"client_id":   utils.Int64OrNil(address.ClientID),
		"supplier_id": utils.Int64OrNil(address.SupplierID),
		"street":      address.Street,
	})

	const query = `
		UPDATE addresses
		SET 
			user_id     = $1,
			client_id   = $2,
			supplier_id = $3,
			street      = $4,
			city        = $5,
			state       = $6,
			country     = $7,
			postal_code = $8,
			updated_at  = NOW()
		WHERE id = $9
		RETURNING updated_at;
	`

	err := r.db.QueryRow(ctx, query,
		address.UserID,
		address.ClientID,
		address.SupplierID,
		address.Street,
		address.City,
		address.State,
		address.Country,
		address.PostalCode,
		address.ID,
	).Scan(&address.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Info(ctx, "[addressRepository] - Endereço não encontrado para atualização", map[string]interface{}{
				"address_id": address.ID,
				"user_id":    utils.Int64OrNil(address.UserID),
			})
			return ErrAddressNotFound
		}

		r.logger.Error(ctx, err, "[addressRepository] - Erro ao atualizar endereço", map[string]interface{}{
			"address_id": address.ID,
			"user_id":    utils.Int64OrNil(address.UserID),
		})
		return fmt.Errorf("%w: %v", ErrUpdateAddress, err)
	}

	r.logger.Info(ctx, "[addressRepository] - Endereço atualizado com sucesso", map[string]interface{}{
		"address_id": address.ID,
		"user_id":    utils.Int64OrNil(address.UserID),
	})

	return nil
}

func (r *addressRepository) Delete(ctx context.Context, id int64) error {
	r.logger.Info(ctx, "[addressRepository] - Iniciando exclusão de endereço", map[string]interface{}{
		"address_id": id,
	})
	const query = `
		DELETE FROM addresses 
		WHERE id = $1
	`

	cmdTag, err := r.db.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error(ctx, err, "[addressRepository] - Erro ao excluir endereço", map[string]interface{}{
			"address_id": id,
		})
		return fmt.Errorf("%w: %v", ErrDeleteAddress, err)
	}

	if cmdTag.RowsAffected() == 0 {
		r.logger.Info(ctx, "[addressRepository] - Endereço não encontrado para exclusão", map[string]interface{}{
			"address_id": id,
		})
		return ErrAddressNotFound
	}

	r.logger.Info(ctx, "[addressRepository] - Endereço excluído com sucesso", map[string]interface{}{
		"address_id": id,
	})

	return nil
}
