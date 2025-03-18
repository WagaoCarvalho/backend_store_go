package repositories

import (
	"context"
	"fmt"

	"github.com/WagaoCarvalho/backend_store_go/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetUsers(ctx context.Context) ([]models.User, error)
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
