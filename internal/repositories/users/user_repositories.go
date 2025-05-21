package repositories

import (
	"context"
	"errors"
	"fmt"

	models_user "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
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
	ErrVersionConflict   = errors.New("conflito de versão: os dados foram modificados por outro processo")
	ErrRecordNotFound    = errors.New("registro não encontrado")
)

type UserRepository interface {
	Create(ctx context.Context, user *models_user.User) (models_user.User, error)
	GetAll(ctx context.Context) ([]models_user.User, error)
	GetById(ctx context.Context, uid int64) (models_user.User, error)
	GetByEmail(ctx context.Context, email string) (models_user.User, error)
	Delete(ctx context.Context, uid int64) error
	Update(ctx context.Context, user models_user.User) (models_user.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models_user.User) (models_user.User, error) {
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

	// Criptografar senha
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return models_user.User{}, ErrPasswordHash
	}
	user.Password = string(hashedPassword)

	query := `
		INSERT INTO users (username, email, password_hash, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`
	err = r.db.QueryRow(ctx, query, user.Username, user.Email, user.Password, user.Status).
		Scan(&user.UID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return models_user.User{}, fmt.Errorf("erro ao criar usuário: %w", err)
	}

	return *user, nil
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

	query := `
		SELECT id, username, email, password_hash, status, created_at, updated_at 
		FROM users WHERE id = $1`
	err := r.db.QueryRow(ctx, query, uid).Scan(
		&user.UID, &user.Username, &user.Email, &user.Password,
		&user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return user, ErrUserNotFound
		}
		return user, fmt.Errorf("erro ao buscar usuário: %w", err)
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

func (r *userRepository) Update(ctx context.Context, user models_user.User) (models_user.User, error) {
	if !utils.IsValidEmail(user.Email) {
		return models_user.User{}, ErrInvalidEmail
	}

	query := `
		UPDATE users 
		SET username = $1, email = $2, status = $3, updated_at = NOW(), version = version + 1
		WHERE id = $4 AND version = $5
		RETURNING updated_at, version
	`

	err := r.db.QueryRow(ctx, query,
		user.Username,
		user.Email,
		user.Status,
		user.UID,
		user.Version,
	).Scan(&user.UpdatedAt, &user.Version)

	if err != nil {
		// Importante: use o erro correto (sql.ErrNoRows) se estiver usando database/sql
		if errors.Is(err, pgx.ErrNoRows) {
			return models_user.User{}, ErrVersionConflict
		}
		return models_user.User{}, fmt.Errorf("erro ao atualizar usuário: %w", err)
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
