package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetVersionByID(ctx context.Context, id int64) (int64, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByName(ctx context.Context, name string) ([]*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
	Disable(ctx context.Context, uid int64) error
	Enable(ctx context.Context, uid int64) error
	Delete(ctx context.Context, id int64) error
	UserExists(ctx context.Context, userID int64) (bool, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	const query = `
		INSERT INTO users (username, email, password_hash, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, user.Username, user.Email, user.Password, user.Status).
		Scan(&user.UID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
	}

	return user, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	const query = `
		SELECT id, username, email, status, created_at, updated_at
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

func (r *userRepository) GetByID(ctx context.Context, uid int64) (*models.User, error) {
	const query = `
		SELECT id, username, email, status, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(ctx, query, uid).Scan(
		&user.UID,
		&user.Username,
		&user.Email,
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

func (r *userRepository) GetVersionByID(ctx context.Context, id int64) (int64, error) {
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

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	const query = `
		SELECT id, username, email, password_hash, status, version, created_at, updated_at
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
		&user.Version,
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

func (r *userRepository) GetByName(ctx context.Context, name string) ([]*models.User, error) {
	const query = `
		SELECT id, username, email, status, created_at, updated_at
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
			&user.Status,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	if len(users) == 0 {
		return nil, errMsg.ErrNotFound
	}

	return users, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
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
			if errCheck := r.db.QueryRow(ctx, checkQuery, user.UID).Scan(&exists); errCheck != nil {
				return nil, fmt.Errorf("%w: erro ao verificar existência: %v", errMsg.ErrUpdate, errCheck)
			}
			if !exists {
				return nil, errMsg.ErrNotFound
			}
			return nil, errMsg.ErrVersionConflict
		}
		return nil, fmt.Errorf("%w: %v", errMsg.ErrUpdate, err)
	}

	return user, nil
}

func (r *userRepository) Disable(ctx context.Context, uid int64) error {
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
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrDisable, err)
	}

	return nil
}

func (r *userRepository) Enable(ctx context.Context, uid int64) error {
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
			return errMsg.ErrNotFound
		}
		return fmt.Errorf("%w: %v", errMsg.ErrEnable, err)
	}

	return nil
}

func (r *userRepository) Delete(ctx context.Context, uid int64) error {
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

func (r *userRepository) UserExists(ctx context.Context, userID int64) (bool, error) {
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
