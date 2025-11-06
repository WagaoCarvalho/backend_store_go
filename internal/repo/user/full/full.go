package repo

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/repo"
	"github.com/jackc/pgx/v5"
)

type UserFull interface {
	CreateTx(ctx context.Context, tx pgx.Tx, user *models.User) (*models.User, error)
	BeginTx(ctx context.Context) (pgx.Tx, error)
}

type userFull struct {
	db repo.DBTransactor
}

func NewUser(db repo.DBTransactor) UserFull {
	return &userFull{
		db: db,
	}
}

func (r *userFull) BeginTx(ctx context.Context) (pgx.Tx, error) {
	return r.db.BeginTx(ctx, pgx.TxOptions{})
}

func (r *userFull) CreateTx(ctx context.Context, tx pgx.Tx, user *models.User) (*models.User, error) {
	const query = `
		INSERT INTO users (username, email, password_hash, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := tx.QueryRow(ctx, query, user.Username, user.Email, user.Password, user.Status).
		Scan(&user.UID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return user, nil
}
