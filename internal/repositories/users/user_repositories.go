package repositories

import (
	"context"
	"errors"
	"fmt"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	models_user "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	models_user_categories "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_categories"
	"github.com/WagaoCarvalho/backend_store_go/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound      = errors.New("usuário não encontrado")
	ErrInvalidEmail      = errors.New("email inválido")
	ErrPasswordHash      = errors.New("erro ao criptografar senha")
	ErrUserAlreadyExists = errors.New("usuário já existe")
)

type UserRepository interface {
	GetAll(ctx context.Context) ([]models_user.User, error)
	GetById(ctx context.Context, uid int64) (models_user.User, error)
	GetByEmail(ctx context.Context, email string) (models_user.User, error)
	Delete(ctx context.Context, uid int64) error
	Update(ctx context.Context, user models_user.User, contact *models_contact.Contact) (models_user.User, error)
	Create(
		ctx context.Context,
		user models_user.User,
		categoryID int64,
		address models_address.Address,
		contact models_contact.Contact,
	) (models_user.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetAll(ctx context.Context) ([]models_user.User, error) {
	query := `SELECT id, username, email, password_hash, status, created_at, updated_at FROM users`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuários: %w", err)
	}
	defer rows.Close()

	var users []models_user.User
	for rows.Next() {
		var user models_user.User
		if err := rows.Scan(&user.UID, &user.Username, &user.Email, &user.Password, &user.Status, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, fmt.Errorf("erro ao ler os dados do usuário: %w", err)
		}
		users = append(users, user)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("erro ao iterar sobre os resultados: %w", rows.Err())
	}

	return users, nil
}

func (r *userRepository) GetById(ctx context.Context, uid int64) (models_user.User, error) {
	var user models_user.User
	var address models_address.Address
	var contact models_contact.Contact
	var categories []models_user_categories.UserCategory

	// 🔹 Consulta dados do usuário
	userQuery := `
		SELECT id, username, email, password_hash, status, created_at, updated_at 
		FROM users WHERE id = $1`
	err := r.db.QueryRow(ctx, userQuery, uid).Scan(
		&user.UID, &user.Username, &user.Email, &user.Password,
		&user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return user, ErrUserNotFound
		}
		return user, fmt.Errorf("erro ao buscar usuário: %w", err)
	}

	// 🔹 Consulta categorias associadas
	categoryQuery := `
		SELECT c.id, c.name, c.description, c.created_at, c.updated_at 
		FROM user_category_relations ucr
		JOIN user_categories c ON ucr.category_id = c.id
		WHERE ucr.user_id = $1`
	rows, err := r.db.Query(ctx, categoryQuery, uid)
	if err != nil {
		return user, fmt.Errorf("erro ao buscar categorias do usuário: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var category models_user_categories.UserCategory
		if err := rows.Scan(&category.ID, &category.Name, &category.Description, &category.CreatedAt, &category.UpdatedAt); err != nil {
			return user, fmt.Errorf("erro ao escanear categorias do usuário: %w", err)
		}
		categories = append(categories, category)
	}
	user.Categories = categories

	// 🔹 Consulta endereço
	addressQuery := `
		SELECT id, street, city, state, country, postal_code, created_at, updated_at 
		FROM addresses WHERE user_id = $1`
	err = r.db.QueryRow(ctx, addressQuery, uid).Scan(
		&address.ID, &address.Street, &address.City, &address.State,
		&address.Country, &address.PostalCode, &address.CreatedAt, &address.UpdatedAt,
	)
	if err != nil && err != pgx.ErrNoRows {
		return user, fmt.Errorf("erro ao buscar endereço do usuário: %w", err)
	} else if err == nil {
		user.Address = &address
	}

	// 🔹 Consulta contato
	contactQuery := `
		SELECT id, user_id, client_id, supplier_id, contact_name, contact_position, email, phone, cell, contact_type, created_at, updated_at 
		FROM contacts WHERE user_id = $1`
	err = r.db.QueryRow(ctx, contactQuery, uid).Scan(
		&contact.ID, &contact.UserID, &contact.ClientID, &contact.SupplierID,
		&contact.ContactName, &contact.ContactPosition, &contact.Email,
		&contact.Phone, &contact.Cell, &contact.ContactType,
		&contact.CreatedAt, &contact.UpdatedAt,
	)
	if err != nil && err != pgx.ErrNoRows {
		return user, fmt.Errorf("erro ao buscar contato do usuário: %w", err)
	} else if err == nil {
		user.Contact = &contact
	}

	return user, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (models_user.User, error) {
	var user models_user.User
	query := `SELECT id, username, email, password_hash, status, created_at, updated_at FROM users WHERE email = $1`

	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.UID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return user, ErrUserNotFound
		}
		return user, fmt.Errorf("erro ao buscar usuário: %w", err)
	}

	return user, nil
}

func (r *userRepository) Create(
	ctx context.Context,
	user models_user.User,
	categoryID int64,
	address models_address.Address,
	contact models_contact.Contact,
) (models_user.User, error) {

	if !utils.IsValidEmail(user.Email) {
		return models_user.User{}, ErrInvalidEmail
	}

	// Verifica se o usuário já existe
	_, err := r.GetByEmail(ctx, user.Email)
	if err == nil {
		return models_user.User{}, ErrUserAlreadyExists
	}
	if !errors.Is(err, ErrUserNotFound) {
		return models_user.User{}, fmt.Errorf("erro ao verificar usuário existente: %w", err)
	}

	// 🔐 Criptografar senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models_user.User{}, ErrPasswordHash
	}
	user.Password = string(hashedPassword)

	// 🔹 Iniciar transação
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return models_user.User{}, fmt.Errorf("erro ao iniciar transação: %w", err)
	}
	defer tx.Rollback(ctx)

	// 🔹 Criar usuário
	userQuery := `
		INSERT INTO users (username, email, password_hash, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err = tx.QueryRow(ctx, userQuery, user.Username, user.Email, user.Password, user.Status).
		Scan(&user.UID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return models_user.User{}, fmt.Errorf("erro ao criar usuário: %w", err)
	}

	// 🔹 Criar endereço
	addressQuery := `
		INSERT INTO addresses (user_id, street, city, state, country, postal_code, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW()) RETURNING id
	`
	err = tx.QueryRow(ctx, addressQuery, user.UID, address.Street, address.City, address.State, address.Country, address.PostalCode).
		Scan(&address.ID)
	if err != nil {
		return models_user.User{}, fmt.Errorf("erro ao criar endereço: %w", err)
	}

	// 🔹 Criar contato associado ao usuário
	contactQuery := `
		INSERT INTO contacts (user_id, contact_name, contact_position, email, phone, cell, contact_type, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id
	`
	err = tx.QueryRow(ctx, contactQuery,
		user.UID,
		contact.ContactName,
		contact.ContactPosition,
		contact.Email,
		contact.Phone,
		contact.Cell,
		contact.ContactType,
	).Scan(&contact.ID)
	if err != nil {
		return models_user.User{}, fmt.Errorf("erro ao criar contato: %w", err)
	}

	// 🔹 Criar relação usuário-categoria
	relationQuery := `
		INSERT INTO user_category_relations (user_id, category_id, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
	`
	_, err = tx.Exec(ctx, relationQuery, user.UID, categoryID)
	if err != nil {
		return models_user.User{}, fmt.Errorf("erro ao criar relação usuário-categoria: %w", err)
	}

	// 🔹 Commitar transação
	if err := tx.Commit(ctx); err != nil {
		return models_user.User{}, fmt.Errorf("erro ao confirmar transação: %w", err)
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user models_user.User, contact *models_contact.Contact) (models_user.User, error) {
	if !utils.IsValidEmail(user.Email) {
		return models_user.User{}, ErrInvalidEmail
	}

	// Atualiza os dados do usuário
	query := `UPDATE users 
			  SET username = $1, email = $2, status = $3, updated_at = NOW() 
			  WHERE id = $4 
			  RETURNING updated_at`

	err := r.db.QueryRow(ctx, query,
		user.Username,
		user.Email,
		user.Status,
		user.UID,
	).Scan(&user.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return models_user.User{}, ErrUserNotFound
		}
		return models_user.User{}, fmt.Errorf("erro ao atualizar usuário: %w", err)
	}

	// Atualiza categorias
	deleteCategoriesQuery := `DELETE FROM user_category_relations WHERE user_id = $1`
	_, err = r.db.Exec(ctx, deleteCategoriesQuery, user.UID)
	if err != nil {
		return models_user.User{}, fmt.Errorf("erro ao remover categorias antigas do usuário: %w", err)
	}

	insertCategoryQuery := `INSERT INTO user_category_relations (user_id, category_id, created_at, updated_at) VALUES ($1, $2, NOW(), NOW())`
	for _, category := range user.Categories {
		_, err = r.db.Exec(ctx, insertCategoryQuery, user.UID, category.ID)
		if err != nil {
			return models_user.User{}, fmt.Errorf("erro ao adicionar categorias ao usuário: %w", err)
		}
	}

	// Atualiza endereço
	if user.Address != nil {
		addressQuery := `UPDATE addresses 
						SET street = $1, city = $2, state = $3, country = $4, postal_code = $5, updated_at = NOW() 
						WHERE user_id = $6`
		_, err = r.db.Exec(ctx, addressQuery,
			user.Address.Street,
			user.Address.City,
			user.Address.State,
			user.Address.Country,
			user.Address.PostalCode,
			user.UID,
		)
		if err != nil {
			return models_user.User{}, fmt.Errorf("erro ao atualizar endereço do usuário: %w", err)
		}
	}

	// Atualiza contato
	if contact != nil {
		contactQuery := `
			UPDATE contacts
			SET contact_name = $1, contact_position = $2, email = $3, phone = $4, cell = $5, contact_type = $6, updated_at = NOW()
			WHERE user_id = $7`

		_, err = r.db.Exec(ctx, contactQuery,
			contact.ContactName,
			contact.ContactPosition,
			contact.Email,
			contact.Phone,
			contact.Cell,
			contact.ContactType,
			user.UID,
		)

		if err != nil {
			return models_user.User{}, fmt.Errorf("erro ao atualizar contato do usuário: %w", err)
		}
	}

	return user, nil
}

func (r *userRepository) Delete(ctx context.Context, uid int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(ctx, query, uid)
	if err != nil {
		return fmt.Errorf("erro ao deletar usuário: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}
