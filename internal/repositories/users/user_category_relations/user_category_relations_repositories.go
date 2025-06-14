package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_category_relations"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrRelationNotFound       = errors.New("relação usuário-categoria não encontrada")
	ErrRelationExists         = errors.New("relação já existe")
	ErrCreateRelation         = errors.New("erro ao criar relação")
	ErrGetRelationsByUser     = errors.New("erro ao buscar relações por usuário")
	ErrGetRelationsByCategory = errors.New("erro ao buscar relações por categoria")
	ErrScanRelation           = errors.New("erro ao ler relação")
	ErrIterateRelations       = errors.New("erro após ler relações")
	ErrDeleteRelation         = errors.New("erro ao deletar relação")
	ErrDeleteAllUserRelations = errors.New("erro ao deletar todas as relações do usuário")
	ErrUpdateRelation         = errors.New("erro ao atualizar relação")
	ErrVersionConflict        = errors.New("conflito de versão: os dados foram modificados por outro processo")
	ErrGetVersionByUserID     = errors.New("erro ao obter versão das relações de categoria do usuário")
)

type UserCategoryRelationRepository interface {
	Create(ctx context.Context, relation *models.UserCategoryRelations) (*models.UserCategoryRelations, error)
	GetAll(ctx context.Context, userID int64) ([]*models.UserCategoryRelations, error)
	GetByUserID(ctx context.Context, userID int64) ([]*models.UserCategoryRelations, error)
	GetByCategoryID(ctx context.Context, categoryID int64) ([]*models.UserCategoryRelations, error)
	Delete(ctx context.Context, userID, categoryID int64) error
	DeleteAll(ctx context.Context, userID int64) error
}

type userCategoryRelationRepositories struct {
	db *pgxpool.Pool
}

func NewUserCategoryRelationRepositories(db *pgxpool.Pool) UserCategoryRelationRepository {
	return &userCategoryRelationRepositories{db: db}
}

func (r *userCategoryRelationRepositories) Create(ctx context.Context, relation *models.UserCategoryRelations) (*models.UserCategoryRelations, error) {
	const query = `
		INSERT INTO user_category_relations (user_id, category_id, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW());
	`

	_, err := r.db.Exec(ctx, query, relation.UserID, relation.CategoryID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return nil, ErrRelationExists
		}
		return nil, fmt.Errorf("%w: %v", ErrCreateRelation, err)
	}

	return relation, nil
}

func (r *userCategoryRelationRepositories) GetAll(ctx context.Context, userID int64) ([]*models.UserCategoryRelations, error) {
	const query = `
		SELECT user_id, category_id, created_at, updated_at
		FROM user_category_relations
		WHERE user_id = $1;
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGetRelationsByUser, err)
	}
	defer rows.Close()

	var relations []*models.UserCategoryRelations
	for rows.Next() {
		var rel models.UserCategoryRelations
		if err := rows.Scan(&rel.UserID, &rel.CategoryID, &rel.CreatedAt, &rel.UpdatedAt); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrScanRelation, err)
		}
		relations = append(relations, &rel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrIterateRelations, err)
	}

	return relations, nil
}

func (r *userCategoryRelationRepositories) GetByUserID(ctx context.Context, userID int64) ([]*models.UserCategoryRelations, error) {
	const query = `
		SELECT user_id, category_id, created_at, updated_at
		FROM user_category_relations
		WHERE user_id = $1;
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGetRelationsByUser, err)
	}
	defer rows.Close()

	var relations []*models.UserCategoryRelations
	for rows.Next() {
		var rel models.UserCategoryRelations
		if err := rows.Scan(&rel.UserID, &rel.CategoryID, &rel.CreatedAt, &rel.UpdatedAt); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrScanRelation, err)
		}
		relations = append(relations, &rel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrIterateRelations, err)
	}

	return relations, nil
}

func (r *userCategoryRelationRepositories) GetVersionByUserID(ctx context.Context, userID int64) (int, error) {
	const query = `
		SELECT COALESCE(SUM(version) + COUNT(*), 0)
		FROM user_category_relations
		WHERE user_id = $1;
	`

	var combinedVersion int
	err := r.db.QueryRow(ctx, query, userID).Scan(&combinedVersion)
	if err != nil {
		return 0, fmt.Errorf("%w: %v", ErrGetVersionByUserID, err)
	}

	// Se o usuário não tem relações, retorna erro (caso deseje tratar isso explicitamente)
	if combinedVersion == 0 {
		return 0, ErrRelationNotFound
	}

	return combinedVersion, nil
}

func (r *userCategoryRelationRepositories) GetByCategoryID(ctx context.Context, categoryID int64) ([]*models.UserCategoryRelations, error) {
	const query = `
		SELECT user_id, category_id, created_at, updated_at
		FROM user_category_relations
		WHERE category_id = $1;
	`

	rows, err := r.db.Query(ctx, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrGetRelationsByCategory, err)
	}
	defer rows.Close()

	var relations []*models.UserCategoryRelations
	for rows.Next() {
		var rel models.UserCategoryRelations
		if err := rows.Scan(&rel.UserID, &rel.CategoryID, &rel.CreatedAt, &rel.UpdatedAt); err != nil {
			return nil, fmt.Errorf("%w: %v", ErrScanRelation, err)
		}
		relations = append(relations, &rel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrIterateRelations, err)
	}

	return relations, nil
}

func (r *userCategoryRelationRepositories) Delete(ctx context.Context, userID, categoryID int64) error {
	const query = `
		DELETE FROM user_category_relations 
		WHERE user_id = $1 AND category_id = $2;
	`

	result, err := r.db.Exec(ctx, query, userID, categoryID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteRelation, err)
	}

	if result.RowsAffected() == 0 {
		return ErrRelationNotFound
	}

	return nil
}

func (r *userCategoryRelationRepositories) DeleteAll(ctx context.Context, userID int64) error {
	const query = `
		DELETE FROM user_category_relations
		WHERE user_id = $1;
	`

	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrDeleteAllUserRelations, err)
	}

	return nil
}
