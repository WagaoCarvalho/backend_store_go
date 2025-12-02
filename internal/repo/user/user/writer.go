package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

func (r *userRepo) Create(ctx context.Context, user *models.User) (*models.User, error) {
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

func (r *userRepo) Update(ctx context.Context, user *models.User) error {

	const querySelect = `
		SELECT version
		FROM users
		WHERE id = $1
	`

	var currentVersion int
	err := r.db.QueryRow(ctx, querySelect, user.UID).Scan(&currentVersion)

	if errors.Is(err, pgx.ErrNoRows) {
		return errMsg.ErrNotFound
	}
	if err != nil {
		return fmt.Errorf("%w: erro ao consultar usuário: %v", errMsg.ErrUpdate, err)
	}

	if currentVersion != user.Version {
		return errMsg.ErrZeroVersion
	}

	const queryUpdate = `
		UPDATE users
		SET 
			username = $1,
			email = $2,
			description = $3,
			status = $4,
			updated_at = NOW(),
			version = version + 1
		WHERE id = $5
		RETURNING updated_at, version
	`

	err = r.db.QueryRow(ctx, queryUpdate,
		user.Username,
		user.Email,
		user.Description,
		user.Status,
		user.UID,
	).Scan(&user.UpdatedAt, &user.Version)

	if err != nil {
		return fmt.Errorf("%w: erro ao atualizar usuário: %v", errMsg.ErrUpdate, err)
	}

	return nil
}

func (r *userRepo) Delete(ctx context.Context, uid int64) error {
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
