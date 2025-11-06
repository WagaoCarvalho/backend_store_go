package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

func (r *user) Create(ctx context.Context, user *models.User) (*models.User, error) {
	const query = `
		INSERT INTO users (
			username,
			email,
			password_hash,
			description,
			status,
			version,
			created_at,
			updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, version, description, created_at, updated_at
	`

	err := r.db.QueryRow(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Password,
		user.Description,
		user.Status,
		user.Version,
	).Scan(&user.UID, &user.Version, &user.Description, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return user, nil
}

func (r *user) Update(ctx context.Context, user *models.User) error {
	const query = `
		UPDATE users
		SET 
			username = $1,
			email = $2,
			description = $3,
			status = $4,
			updated_at = NOW(),
			version = version + 1
		WHERE id = $5 AND version = $6
		RETURNING updated_at, version
	`

	err := r.db.QueryRow(ctx, query,
		user.Username,
		user.Email,
		user.Description,
		user.Status,
		user.UID,
		user.Version,
	).Scan(&user.UpdatedAt, &user.Version)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errMsg.ErrVersionConflict
		}
		return fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (r *user) Delete(ctx context.Context, uid int64) error {
	const query = `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(ctx, query, uid)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	if result.RowsAffected() == 0 {
		return errMsg.ErrNotFound
	}

	return nil
}
