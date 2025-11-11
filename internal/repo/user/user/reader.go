package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

func (r *userRepo) GetAll(ctx context.Context) ([]*models.User, error) {
	const query = `
		SELECT 
			id,
			username,
			email,
			description,
			status,
			created_at,
			updated_at
		FROM users
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.UID,
			&user.Username,
			&user.Email,
			&user.Description,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrIterate, err)
	}

	return users, nil
}

func (r *userRepo) GetByID(ctx context.Context, id int64) (*models.User, error) {
	const query = `
		SELECT 
			id,
			username,
			email,
			description,
			status,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.UID,
		&user.Username,
		&user.Email,
		&user.Description,
		&user.Status,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errMsg.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return user, nil
}

func (r *userRepo) GetVersionByID(ctx context.Context, id int64) (int64, error) {
	const query = `SELECT version FROM users WHERE id = $1`

	var version int64
	err := r.db.QueryRow(ctx, query, id).Scan(&version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errMsg.ErrNotFound
		}
		return 0, fmt.Errorf("%w: %v", errMsg.ErrGetVersion, err)
	}

	return version, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	const query = `
		SELECT id, username, email, password_hash, status, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &models.User{}
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errMsg.ErrNotFound
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return user, nil
}

func (r *userRepo) GetByName(ctx context.Context, name string) ([]*models.User, error) {
	const query = `
		SELECT 
			id,
			username,
			email,
			description,
			status,
			created_at,
			updated_at
		FROM users
		WHERE username ILIKE $1
		ORDER BY username ASC
	`

	rows, err := r.db.Query(ctx, query, "%"+name+"%")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := new(models.User)
		if err := rows.Scan(
			&user.UID,
			&user.Username,
			&user.Email,
			&user.Description,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrIterate, err)
	}

	if len(users) == 0 {
		return nil, errMsg.ErrNotFound
	}

	return users, nil
}

func (r *userRepo) UserExists(ctx context.Context, userID int64) (bool, error) {
	const query = `
		SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE id = $1
		)
	`

	var exists bool
	err := r.db.QueryRow(ctx, query, userID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return exists, nil
}
