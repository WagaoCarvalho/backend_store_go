package repositories

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_categories"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrCategoryNotFound  = errors.New("categoria não encontrada")
	ErrGetCategories     = errors.New("erro ao buscar categorias")
	ErrScanCategory      = errors.New("erro ao ler os dados da categoria")
	ErrIterateCategories = errors.New("erro ao iterar sobre os resultados")
	ErrGetCategoryByID   = errors.New("erro ao buscar categoria por ID")
	ErrCreateCategory    = errors.New("erro ao criar categoria")
	ErrUpdateCategory    = errors.New("erro ao atualizar categoria")
	ErrDeleteCategory    = errors.New("erro ao deletar categoria")
	ErrVersionConflict   = errors.New("conflito de versão: os dados foram modificados por outro processo")
)

type UserCategoryRepository interface {
	GetAll(ctx context.Context) ([]*models.UserCategory, error)
	GetByID(ctx context.Context, id int64) (*models.UserCategory, error)
	Create(ctx context.Context, category *models.UserCategory) (*models.UserCategory, error)
	Update(ctx context.Context, category *models.UserCategory) error
	Delete(ctx context.Context, id int64) error
}

type userCategoryRepository struct {
	db *pgxpool.Pool
}

func NewUserCategoryRepository(db *pgxpool.Pool) UserCategoryRepository {
	return &userCategoryRepository{db: db}
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
		return nil, fmt.Errorf("%w: %v", ErrCreateCategory, err)
	}

	return category, nil
}

func (r *userCategoryRepository) GetAll(ctx context.Context) ([]*models.UserCategory, error) {
	query := `SELECT id, name, description, created_at, updated_at FROM user_categories`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGetCategories, err)
	}
	defer rows.Close()

	var categories []*models.UserCategory
	for rows.Next() {
		category := new(models.UserCategory)
		if err := rows.Scan(&category.ID, &category.Name, &category.Description,
			&category.CreatedAt, &category.UpdatedAt); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrScanCategory, err)
		}
		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrIterateCategories, err)
	}

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
			return nil, ErrCategoryNotFound
		}
		return nil, fmt.Errorf("%w: %v", ErrGetCategoryByID, err)
	}

	return &category, nil
}

func (r *userCategoryRepository) Update(ctx context.Context, category *models.UserCategory) error {
	const query = `
		UPDATE user_categories
		SET 
			name        = $1,
			description = $2,
			updated_at  = NOW(),
			version     = version + 1
		WHERE 
			id      = $3
		AND 
			version = $4
		RETURNING 
			updated_at,
			version;
	`

	row := r.db.QueryRow(ctx, query,
		category.Name,
		category.Description,
		category.ID,
		category.Version,
	)

	if err := row.Scan(&category.UpdatedAt, &category.Version); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {

			var exists bool
			checkQuery := `SELECT EXISTS(SELECT 1 FROM user_categories WHERE id = $1)`
			checkErr := r.db.QueryRow(ctx, checkQuery, category.ID).Scan(&exists)
			if checkErr != nil {
				return fmt.Errorf("%w: erro ao verificar existência: %v", ErrUpdateCategory, checkErr)
			}
			if !exists {
				return ErrCategoryNotFound
			}
			return ErrVersionConflict
		}
		return fmt.Errorf("%w: %v", ErrUpdateCategory, err)
	}

	return nil
}

func (r *userCategoryRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM user_categories WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteCategory, err)
	}

	if result.RowsAffected() == 0 {
		return ErrCategoryNotFound
	}

	return nil
}
