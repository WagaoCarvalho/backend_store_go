package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_categories"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserCategoryRepository interface {
	GetAll(ctx context.Context) ([]*models.UserCategory, error)
	GetByID(ctx context.Context, id int64) (*models.UserCategory, error)
	Create(ctx context.Context, category *models.UserCategory) (*models.UserCategory, error)
	Update(ctx context.Context, category *models.UserCategory) error
	Delete(ctx context.Context, id int64) error
}

type userCategoryRepository struct {
	db     *pgxpool.Pool
	logger *logger.LoggerAdapter
}

func NewUserCategoryRepository(db *pgxpool.Pool, logger *logger.LoggerAdapter) UserCategoryRepository {
	return &userCategoryRepository{db: db, logger: logger}
}

func (r *userCategoryRepository) Create(ctx context.Context, category *models.UserCategory) (*models.UserCategory, error) {
	const query = `
		INSERT INTO user_categories (name, description, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW())
		RETURNING id, created_at, updated_at;
	`

	err := r.db.QueryRow(ctx, query, category.Name, category.Description).
		Scan(&category.ID, &category.CreatedAt, &category.UpdatedAt)

	if err != nil {
		r.logger.Error(ctx, err, "Erro ao criar categoria de usuário", map[string]interface{}{
			"name":        category.Name,
			"description": category.Description,
		})
		return nil, fmt.Errorf("%w: %v", ErrCreateCategory, err)
	}

	r.logger.Info(ctx, "Categoria de usuário criada com sucesso", map[string]interface{}{
		"category_id": category.ID,
		"name":        category.Name,
	})

	return category, nil
}

func (r *userCategoryRepository) GetAll(ctx context.Context) ([]*models.UserCategory, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM user_categories`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		r.logger.Error(ctx, err, "Erro ao buscar todas as categorias de usuário", nil)
		return nil, fmt.Errorf("%w: %v", ErrGetCategories, err)
	}
	defer rows.Close()

	var categories []*models.UserCategory
	for rows.Next() {
		category := new(models.UserCategory)
		if err := rows.Scan(&category.ID, &category.Name, &category.Description,
			&category.CreatedAt, &category.UpdatedAt); err != nil {
			r.logger.Error(ctx, err, "Erro ao fazer scan de categoria de usuário", nil)
			return nil, fmt.Errorf("%w: %v", ErrScanCategory, err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(ctx, err, "Erro ao iterar categorias de usuário", nil)
		return nil, fmt.Errorf("%w: %v", ErrIterateCategories, err)
	}

	r.logger.Info(ctx, "Categorias de usuário buscadas com sucesso", map[string]interface{}{
		"total": len(categories),
	})

	return categories, nil
}

func (r *userCategoryRepository) GetByID(ctx context.Context, id int64) (*models.UserCategory, error) {
	const query = `
		SELECT id, name, description, created_at, updated_at
		FROM user_categories
		WHERE id = $1;
	`

	var category models.UserCategory

	err := r.db.QueryRow(ctx, query, id).Scan(
		&category.ID,
		&category.Name,
		&category.Description,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn(ctx, "Categoria de usuário não encontrada", map[string]interface{}{
				"category_id": id,
			})
			return nil, ErrCategoryNotFound
		}

		r.logger.Error(ctx, err, "Erro ao buscar categoria de usuário por ID", map[string]interface{}{
			"category_id": id,
		})
		return nil, fmt.Errorf("%w: %v", ErrGetCategoryByID, err)
	}

	r.logger.Info(ctx, "Categoria de usuário buscada com sucesso", map[string]interface{}{
		"category_id": category.ID,
	})

	return &category, nil
}

func (r *userCategoryRepository) Update(ctx context.Context, category *models.UserCategory) error {
	const query = `
		UPDATE user_categories
		SET 
			name        = $1,
			description = $2,
			updated_at  = NOW()
		WHERE 
			id = $3
		RETURNING 
			updated_at;
	`

	err := r.db.QueryRow(ctx, query,
		category.Name,
		category.Description,
		category.ID,
	).Scan(&category.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn(ctx, "Categoria de usuário não encontrada para atualização", map[string]interface{}{
				"category_id": category.ID,
			})
			return ErrCategoryNotFound
		}

		r.logger.Error(ctx, err, "Erro ao atualizar categoria de usuário", map[string]interface{}{
			"category_id": category.ID,
			"name":        category.Name,
		})
		return fmt.Errorf("%w: %v", ErrUpdateCategory, err)
	}

	r.logger.Info(ctx, "Categoria de usuário atualizada com sucesso", map[string]interface{}{
		"category_id": category.ID,
	})

	return nil
}

func (r *userCategoryRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM user_categories WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		r.logger.Error(ctx, err, "Erro ao excluir categoria de usuário", map[string]interface{}{
			"category_id": id,
		})
		return fmt.Errorf("%w: %v", ErrDeleteCategory, err)
	}

	if result.RowsAffected() == 0 {
		r.logger.Warn(ctx, "Categoria de usuário não encontrada para exclusão", map[string]interface{}{
			"category_id": id,
		})
		return ErrCategoryNotFound
	}

	r.logger.Info(ctx, "Categoria de usuário excluída com sucesso", map[string]interface{}{
		"category_id": id,
	})

	return nil
}
