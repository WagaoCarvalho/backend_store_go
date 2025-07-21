package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	CreateTx(ctx context.Context, tx pgx.Tx, user *models.User) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetVersionByID(ctx context.Context, id int64) (int64, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Disable(ctx context.Context, uid int64) error
	Enable(ctx context.Context, uid int64) error
	Delete(ctx context.Context, id int64) error
}

type userRepository struct {
	db     *pgxpool.Pool
	logger *logger.LoggerAdapter
}

func NewUserRepository(db *pgxpool.Pool, logger *logger.LoggerAdapter) UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	ref := "[userRepository - Create] - "
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

	err := r.db.QueryRow(ctx, query, user.Username, user.Email, user.Password, user.Status).
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

func (r *userRepository) CreateTx(ctx context.Context, tx pgx.Tx, user *models.User) (*models.User, error) {
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

func (r *userRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	ref := "[userRepository - GetAll] - "
	r.logger.Info(ctx, ref+logger.LogGetInit, nil)

	const query = `
		SELECT id, username, email, password_hash, status, created_at, updated_at
		FROM users
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogGetError, nil)
		return nil, fmt.Errorf("%w: %v", ErrGetUsers, err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.UID,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			r.logger.Error(ctx, err, ref+logger.LogGetErrorScan, nil)
			return nil, fmt.Errorf("%w: %v", ErrScanUserRow, err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(ctx, err, ref+logger.LogIterateError, nil)
		return nil, fmt.Errorf("%w: %v", ErrIterateUserRows, err)
	}

	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"total_users": len(users),
	})

	return users, nil
}

func (r *userRepository) GetByID(ctx context.Context, uid int64) (*models.User, error) {
	ref := "[userRepository - GetByID] - "
	r.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"user_id": uid,
	})

	user := &models.User{}

	const query = `
		SELECT id, username, email, password_hash, status, created_at, updated_at
		FROM users
		WHERE id = $1
	`

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
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"user_id": uid,
			})
			return nil, ErrUserNotFound
		}

		r.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"user_id": uid,
		})
		return nil, fmt.Errorf("%w: %v", ErrFetchUser, err)
	}

	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"user_id": uid,
	})

	return user, nil
}

func (r *userRepository) GetVersionByID(ctx context.Context, id int64) (int64, error) {
	ref := "[userRepository - GetVersionByID] - "
	r.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"user_id": id,
	})

	const query = `SELECT version FROM users WHERE id = $1`

	var version int64
	err := r.db.QueryRow(ctx, query, id).Scan(&version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"user_id": id,
			})
			return 0, ErrUserNotFound
		}

		r.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"user_id": id,
		})
		return 0, fmt.Errorf("%w: %v", ErrFetchUserVersion, err)
	}

	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"user_id": id,
		"version": version,
	})

	return version, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	ref := "[userRepository - GetByEmail] - "
	r.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"email": email,
	})

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
			r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"email": email,
			})
			return nil, ErrUserNotFound
		}

		r.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"email": email,
		})
		return nil, fmt.Errorf("%w: %v", ErrFetchUser, err)
	}

	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"user_id": user.UID,
		"email":   email,
	})

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
	ref := "[userRepository - Update] - "
	r.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"user_id":  user.UID,
		"username": user.Username,
		"email":    user.Email,
		"status":   user.Status,
	})

	const query = `
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
		if errors.Is(err, pgx.ErrNoRows) {
			var exists bool
			checkQuery := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1)`
			checkErr := r.db.QueryRow(ctx, checkQuery, user.UID).Scan(&exists)
			if checkErr != nil {
				r.logger.Error(ctx, checkErr, ref+logger.LogUpdateError, map[string]any{
					"user_id": user.UID,
				})
				return nil, fmt.Errorf("%w: erro ao verificar existência: %v", ErrUpdateUser, checkErr)
			}
			if !exists {
				r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
					"user_id": user.UID,
				})
				return nil, ErrUserNotFound
			}
			r.logger.Warn(ctx, ref+logger.LogUpdateVersionConflict, map[string]any{
				"user_id": user.UID,
			})
			return nil, ErrVersionConflict
		}

		r.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"user_id": user.UID,
		})
		return nil, fmt.Errorf("%w: %v", ErrUpdateUser, err)
	}

	r.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"user_id": user.UID,
	})

	return user, nil
}

func (r *userRepository) Disable(ctx context.Context, uid int64) error {
	ref := "[userRepository - Disable] - "
	r.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"user_id": uid,
	})

	const query = `
		UPDATE users
		SET status = FALSE, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version, updated_at;
	`

	var version int
	var updatedAt time.Time
	err := r.db.QueryRow(ctx, query, uid).Scan(&version, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"user_id": uid,
			})
			return ErrUserNotFound
		}

		r.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"user_id": uid,
		})
		return fmt.Errorf("%w: %v", ErrUpdateUser, err)
	}

	r.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"user_id":    uid,
		"new_status": false,
	})

	return nil
}

func (r *userRepository) Enable(ctx context.Context, uid int64) error {
	ref := "[userRepository - Enable] - "
	r.logger.Info(ctx, ref+logger.LogUpdateInit, map[string]any{
		"user_id": uid,
	})

	const query = `
		UPDATE users
		SET status = TRUE, updated_at = NOW(), version = version + 1
		WHERE id = $1
		RETURNING version, updated_at;
	`

	var version int
	var updatedAt time.Time
	err := r.db.QueryRow(ctx, query, uid).Scan(&version, &updatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
				"user_id": uid,
			})
			return ErrUserNotFound
		}

		r.logger.Error(ctx, err, ref+logger.LogUpdateError, map[string]any{
			"user_id": uid,
		})
		return fmt.Errorf("%w: %v", ErrUpdateUser, err)
	}

	r.logger.Info(ctx, ref+logger.LogUpdateSuccess, map[string]any{
		"user_id":    uid,
		"new_status": true,
	})

	return nil
}

func (r *userRepository) Delete(ctx context.Context, uid int64) error {
	ref := "[userRepository - Delete] - "
	r.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{
		"user_id": uid,
	})

	const query = `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(ctx, query, uid)
	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"user_id": uid,
		})
		return fmt.Errorf("%w: %v", ErrDeleteUser, err)
	}

	if result.RowsAffected() == 0 {
		r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
			"user_id": uid,
		})
		return ErrUserNotFound
	}

	r.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"user_id": uid,
	})

	return nil
}
