package repositories

import (
	"context"
	"fmt"

	"github.com/WagaoCarvalho/backend_store_go/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserById(ctx context.Context, uid int64) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
	CreateUser(ctx context.Context, user models.User) (models.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) GetUsers(ctx context.Context) ([]models.User, error) {
	query := `SELECT id, username, email, password_hash, status, created_at, updated_at FROM users`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar usuários: %w", err)
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		scanErr := rows.Scan(
			&user.UID, &user.Username, &user.Email,
			&user.Password, &user.Status, &user.CreatedAt, &user.UpdatedAt,
		)
		if scanErr != nil {
			return nil, fmt.Errorf("erro ao ler os dados do usuário: %w", scanErr)
		}
		users = append(users, user)
	}

	if rowsErr := rows.Err(); rowsErr != nil {
		return nil, fmt.Errorf("erro ao iterar sobre os resultados: %w", rowsErr)
	}

	return users, nil
}

func (r *userRepository) GetUserById(ctx context.Context, uid int64) (models.User, error) {
	var user models.User
	query := `SELECT id, username, email, password_hash, status, created_at, updated_at FROM users WHERE id = $1`

	err := r.db.QueryRow(ctx, query, uid).Scan(
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
			return user, fmt.Errorf("usuário não encontrado")
		}
		return user, fmt.Errorf("erro ao buscar usuário: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	var user models.User
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
			return user, fmt.Errorf("usuário não encontrado")
		}
		return user, fmt.Errorf("erro ao buscar usuário: %w", err)
	}

	return user, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	query := `INSERT INTO users (username, email, password_hash, status, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, NOW(), NOW()) RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(ctx, query,
		user.Username,
		user.Email,
		user.Password,
		user.Status,
	).Scan(&user.UID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return models.User{}, fmt.Errorf("erro ao criar usuário: %w", err)
	}

	return user, nil
}
