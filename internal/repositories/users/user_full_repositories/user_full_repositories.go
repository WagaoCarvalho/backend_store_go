package repositories

import (
	"context"
	"fmt"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserFullRepository interface {
	CreateTx(ctx context.Context, tx pgx.Tx, user *models.User) (*models.User, error)
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type userFullRepository struct {
	db     *pgxpool.Pool
	logger *logger.LoggerAdapter
}

func NewUserRepository(db *pgxpool.Pool, logger *logger.LoggerAdapter) UserFullRepository {
	return &userFullRepository{
		db:     db,
		logger: logger,
	}
}

func (r *userFullRepository) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.db.BeginTx(ctx, pgx.TxOptions{})
}

func (r *userFullRepository) CreateTx(ctx context.Context, tx pgx.Tx, user *models.User) (*models.User, error) {
	ref := "[userRepository - CreateTx] - "
	r.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{
		"username": user.Username,
		"email":    user.Email,
		"status":   user.Status,
	})

	const query = `
		INSERT INTO users (username, email, password_hash, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := tx.QueryRow(ctx, query, user.Username, user.Email, user.Password, user.Status).
		Scan(&user.UID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
			"username": user.Username,
			"email":    user.Email,
			"status":   user.Status,
		})
		return nil, fmt.Errorf("%w: %v", ErrCreateUser, err)
	}

	r.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"user_id":  user.UID,
		"username": user.Username,
		"email":    user.Email,
	})

	return user, nil
}
