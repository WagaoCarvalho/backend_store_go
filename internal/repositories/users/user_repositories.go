package repositories

import (
	"context"
	"errors"
	"fmt"

	models_user "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrUserNotFound    = errors.New("usuário não encontrado")
	ErrVersionConflict = errors.New("conflito de versão: os dados foram modificados por outro processo")
	ErrCreateUser      = errors.New("erro ao criar usuário")
	ErrGetUsers        = errors.New("erro ao buscar usuários")
	ErrScanUserRow     = errors.New("erro ao ler os dados do usuário")
	ErrIterateUserRows = errors.New("erro ao iterar sobre os resultados")
	ErrFetchUser       = errors.New("erro ao buscar usuário")
	ErrUpdateUser      = errors.New("erro ao atualizar usuário")
	ErrDeleteUser      = errors.New("erro ao deletar usuário")
)

type UserRepository interface {
	Create(ctx context.Context, user *models_user.User) (*models_user.User, error)
	GetAll(ctx context.Context) ([]*models_user.User, error)
	GetByID(ctx context.Context, id int64) (*models_user.User, error)
	GetVersionByID(ctx context.Context, id int64) (int64, error)
	GetByEmail(ctx context.Context, email string) (*models_user.User, error)
	Update(ctx context.Context, user *models_user.User) (*models_user.User, error)
	Delete(ctx context.Context, id int64) error
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(ctx context.Context, user *models_user.User) (*models_user.User, error) {
	query := `
		INSERT INTO users (username, email, password_hash, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, user.Username, user.Email, user.Password, user.Status).
		Scan(&user.UID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrCreateUser, err)
	}

	return user, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]*models_user.User, error) {
	query := `SELECT id, username, email, password_hash, status, created_at, updated_at FROM users`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGetUsers, err)
	}
	defer rows.Close()

	var users []*models_user.User
	for rows.Next() {
		var user models_user.User
		if err := rows.Scan(&user.UID, &user.Username, &user.Email, &user.Password, &user.Status, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrScanUserRow, err)
		}

		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrIterateUserRows, err)
	}

	return users, nil
}

func (r *userRepository) GetByID(ctx context.Context, uid int64) (*models_user.User, error) {
	user := &models_user.User{}

	query := `SELECT id, username, email, password_hash, status, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.QueryRow(ctx, query, uid).Scan(
		&user.UID, &user.Username, &user.Email, &user.Password,
		&user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrFetchUser, err)
	}

	return user, nil
}

func (r *userRepository) GetVersionByID(ctx context.Context, id int64) (int64, error) {
	const query = `SELECT version FROM users WHERE id = $1`

	var version int64
	err := r.db.QueryRow(ctx, query, id).Scan(&version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrUserNotFound
		}
		return 0, fmt.Errorf("erro ao buscar versão do usuário: %w", err)
	}

	return version, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models_user.User, error) {
	user := &models_user.User{}
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
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrFetchUser, err)
	}

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models_user.User) (*models_user.User, error) {
	const query = `
		UPDATE users 
		SET 
			username   = $1,
			email      = $2,
			status     = $3,
			updated_at = NOW(),
			version    = version + 1
		WHERE 
			id      = $4 AND 
			version = $5
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
				return nil, fmt.Errorf("%w: erro ao verificar existência: %v", ErrUpdateUser, checkErr)
			}
			if !exists {
				return nil, ErrUserNotFound
			}
			return nil, ErrVersionConflict
		}
		return nil, fmt.Errorf("%w: %v", ErrUpdateUser, err)
	}

	return user, nil
}

func (r *userRepository) Delete(ctx context.Context, uid int64) error {
	const query = `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(ctx, query, uid)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteUser, err)
	}

	if result.RowsAffected() == 0 {
		return ErrUserNotFound
	}

	return nil
}
