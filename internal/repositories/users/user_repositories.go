package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) (*models.User, error)
	GetAll(ctx context.Context) ([]*models.User, error)
	GetByID(ctx context.Context, id int64) (*models.User, error)
	GetVersionByID(ctx context.Context, id int64) (int64, error)
	GetByEmail(ctx context.Context, email string) (*models.User, error)
	Update(ctx context.Context, user *models.User) (*models.User, error)
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
	query := `
		INSERT INTO users (username, email, password_hash, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(ctx, query, user.Username, user.Email, user.Password, user.Status).
		Scan(&user.UID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		r.logger.Error(ctx, err, "Erro ao criar usuário", map[string]interface{}{
			"username": user.Username,
			"email":    user.Email,
			"status":   user.Status,
		})
		return nil, fmt.Errorf("%w: %v", ErrCreateUser, err)
	}

	r.logger.Info(ctx, "Usuário criado com sucesso", map[string]interface{}{
		"user_id":  user.UID,
		"username": user.Username,
		"email":    user.Email,
	})

	return user, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]*models.User, error) {
	query := `SELECT id, username, email, password_hash, status, created_at, updated_at FROM users`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		r.logger.Error(ctx, err, "Erro ao buscar todos os usuários", nil)
		return nil, fmt.Errorf("%w: %v", ErrGetUsers, err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.UID, &user.Username, &user.Email, &user.Password, &user.Status, &user.CreatedAt, &user.UpdatedAt); err != nil {
			r.logger.Error(ctx, err, "Erro ao escanear linha de usuário", nil)
			return nil, fmt.Errorf("%w: %v", ErrScanUserRow, err)
		}
		users = append(users, &user)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(ctx, err, "Erro ao iterar sobre as linhas de usuários", nil)
		return nil, fmt.Errorf("%w: %v", ErrIterateUserRows, err)
	}

	r.logger.Info(ctx, "Usuários obtidos com sucesso", map[string]interface{}{
		"total_users": len(users),
	})

	return users, nil
}

func (r *userRepository) GetByID(ctx context.Context, uid int64) (*models.User, error) {
	user := &models.User{}

	query := `SELECT id, username, email, password_hash, status, created_at, updated_at FROM users WHERE id = $1`
	err := r.db.QueryRow(ctx, query, uid).Scan(
		&user.UID, &user.Username, &user.Email, &user.Password,
		&user.Status, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn(ctx, "Usuário não encontrado", map[string]interface{}{
				"user_id": uid,
			})
			return nil, ErrUserNotFound
		}

		r.logger.Error(ctx, err, "Erro ao buscar usuário por ID", map[string]interface{}{
			"user_id": uid,
		})
		return nil, fmt.Errorf("%w: %v", ErrFetchUser, err)
	}

	r.logger.Info(ctx, "Usuário recuperado com sucesso", map[string]interface{}{
		"user_id": uid,
	})

	return user, nil
}

func (r *userRepository) GetVersionByID(ctx context.Context, id int64) (int64, error) {
	const query = `SELECT version FROM users WHERE id = $1`

	var version int64
	err := r.db.QueryRow(ctx, query, id).Scan(&version)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn(ctx, "Versão não encontrada: usuário não existe", map[string]interface{}{
				"user_id": id,
			})
			return 0, ErrUserNotFound
		}

		r.logger.Error(ctx, err, "Erro ao buscar versão do usuário", map[string]interface{}{
			"user_id": id,
		})
		return 0, fmt.Errorf("erro ao buscar versão do usuário: %w", err)
	}

	r.logger.Info(ctx, "Versão do usuário obtida com sucesso", map[string]interface{}{
		"user_id": id,
		"version": version,
	})

	return version, nil
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user := &models.User{}
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
			r.logger.Warn(ctx, "Usuário não encontrado com o email informado", map[string]interface{}{
				"email": email,
			})
			return nil, ErrUserNotFound
		}

		r.logger.Error(ctx, err, "Erro ao buscar usuário por email", map[string]interface{}{
			"email": email,
		})
		return nil, fmt.Errorf("%w: %v", ErrFetchUser, err)
	}

	r.logger.Info(ctx, "Usuário obtido com sucesso por email", map[string]interface{}{
		"user_id": user.UID,
		"email":   email,
	})

	return user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) (*models.User, error) {
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
				r.logger.Error(ctx, checkErr, "Erro ao verificar existência do usuário durante conflito de versão", map[string]interface{}{
					"user_id": user.UID,
				})
				return nil, fmt.Errorf("%w: erro ao verificar existência: %v", ErrUpdateUser, checkErr)
			}
			if !exists {
				r.logger.Warn(ctx, "Usuário não encontrado ao tentar atualizar", map[string]interface{}{
					"user_id": user.UID,
				})
				return nil, ErrUserNotFound
			}
			r.logger.Warn(ctx, "Conflito de versão ao atualizar usuário", map[string]interface{}{
				"user_id": user.UID,
			})
			return nil, ErrVersionConflict
		}

		r.logger.Error(ctx, err, "Erro ao atualizar usuário", map[string]interface{}{
			"user_id": user.UID,
		})
		return nil, fmt.Errorf("%w: %v", ErrUpdateUser, err)
	}

	r.logger.Info(ctx, "Usuário atualizado com sucesso", map[string]interface{}{
		"user_id": user.UID,
	})

	return user, nil
}

func (r *userRepository) Delete(ctx context.Context, uid int64) error {
	const query = `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(ctx, query, uid)
	if err != nil {
		r.logger.Error(ctx, err, "Erro ao deletar usuário", map[string]interface{}{
			"user_id": uid,
		})
		return fmt.Errorf("%w: %v", ErrDeleteUser, err)
	}

	rows := result.RowsAffected()
	r.logger.Info(ctx, "Resultado do DELETE", map[string]interface{}{
		"user_id": uid,
		"rows":    rows,
	})

	if rows == 0 {
		r.logger.Warn(ctx, "Nenhum usuário encontrado para deletar", map[string]interface{}{
			"user_id": uid,
		})
		return ErrUserNotFound
	}

	r.logger.Info(ctx, "Usuário deletado com sucesso", map[string]interface{}{
		"user_id": uid,
	})

	return nil
}
