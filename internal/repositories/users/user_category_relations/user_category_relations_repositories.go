package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_category_relations"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserCategoryRelationRepository interface {
	Create(ctx context.Context, relation *models.UserCategoryRelations) (*models.UserCategoryRelations, error)
	HasUserCategoryRelation(ctx context.Context, userID, categoryID int64) (bool, error)
	GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserCategoryRelations, error)
	Delete(ctx context.Context, userID, categoryID int64) error
	DeleteAll(ctx context.Context, userID int64) error
}

type userCategoryRelationRepositories struct {
	db     *pgxpool.Pool
	logger *logger.LoggerAdapter
}

func NewUserCategoryRelationRepositories(db *pgxpool.Pool, logger *logger.LoggerAdapter) UserCategoryRelationRepository {
	return &userCategoryRelationRepositories{db: db, logger: logger}
}

func (r *userCategoryRelationRepositories) Create(ctx context.Context, relation *models.UserCategoryRelations) (*models.UserCategoryRelations, error) {
	const query = `
		INSERT INTO user_category_relations (user_id, category_id, created_at)
		VALUES ($1, $2, NOW());
	`

	_, err := r.db.Exec(ctx, query, relation.UserID, relation.CategoryID)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			r.logger.Warn(ctx, "Relação já existente entre usuário e categoria", map[string]interface{}{
				"user_id":     relation.UserID,
				"category_id": relation.CategoryID,
			})
			return nil, ErrRelationExists
		}
		if strings.Contains(err.Error(), "violates foreign key constraint") {
			r.logger.Warn(ctx, "Chave estrangeira inválida na criação da relação", map[string]interface{}{
				"user_id":     relation.UserID,
				"category_id": relation.CategoryID,
			})
			return nil, ErrInvalidForeignKey
		}

		r.logger.Error(ctx, err, "Erro ao criar relação entre usuário e categoria", map[string]interface{}{
			"user_id":     relation.UserID,
			"category_id": relation.CategoryID,
		})
		return nil, fmt.Errorf("%w: %v", ErrCreateRelation, err)
	}

	r.logger.Info(ctx, "Relação entre usuário e categoria criada com sucesso", map[string]interface{}{
		"user_id":     relation.UserID,
		"category_id": relation.CategoryID,
	})

	return relation, nil
}

func (r *userCategoryRelationRepositories) GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserCategoryRelations, error) {
	const query = `
		SELECT user_id, category_id, created_at
		FROM user_category_relations
		WHERE user_id = $1;
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		r.logger.Error(ctx, err, "Erro ao buscar relações por user_id", map[string]interface{}{
			"user_id": userID,
		})
		return nil, fmt.Errorf("%w: %v", ErrGetRelationsByUser, err)
	}
	defer rows.Close()

	var relations []*models.UserCategoryRelations
	for rows.Next() {
		var rel models.UserCategoryRelations
		if err := rows.Scan(&rel.UserID, &rel.CategoryID, &rel.CreatedAt); err != nil {
			r.logger.Error(ctx, err, "Erro ao escanear relação de categoria de usuário", map[string]interface{}{
				"user_id": userID,
			})
			return nil, fmt.Errorf("%w: %v", ErrScanRelation, err)
		}
		relations = append(relations, &rel)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(ctx, err, "Erro ao iterar sobre as relações de categoria de usuário", map[string]interface{}{
			"user_id": userID,
		})
		return nil, fmt.Errorf("%w: %v", ErrIterateRelations, err)
	}

	r.logger.Info(ctx, "Relações de categorias do usuário obtidas com sucesso", map[string]interface{}{
		"user_id":         userID,
		"relations_count": len(relations),
	})

	return relations, nil
}

func (r *userCategoryRelationRepositories) HasUserCategoryRelation(ctx context.Context, userID, categoryID int64) (bool, error) {
	const query = `
		SELECT 1
		FROM user_category_relations
		WHERE user_id = $1 AND category_id = $2
		LIMIT 1;
	`

	var exists int
	err := r.db.QueryRow(ctx, query, userID, categoryID).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Info(ctx, "Relação entre usuário e categoria não existe", map[string]interface{}{
				"user_id":     userID,
				"category_id": categoryID,
			})
			return false, nil
		}

		r.logger.Error(ctx, err, "Erro ao verificar existência de relação entre usuário e categoria", map[string]interface{}{
			"user_id":     userID,
			"category_id": categoryID,
		})
		return false, fmt.Errorf("%w: %v", ErrCheckRelationExists, err)
	}

	r.logger.Info(ctx, "Relação entre usuário e categoria existe", map[string]interface{}{
		"user_id":     userID,
		"category_id": categoryID,
	})

	return true, nil
}

func (r *userCategoryRelationRepositories) Delete(ctx context.Context, userID, categoryID int64) error {
	const query = `
		DELETE FROM user_category_relations 
		WHERE user_id = $1 AND category_id = $2;
	`

	result, err := r.db.Exec(ctx, query, userID, categoryID)
	if err != nil {
		r.logger.Error(ctx, err, "Erro ao excluir relação entre usuário e categoria", map[string]interface{}{
			"user_id":     userID,
			"category_id": categoryID,
		})
		return fmt.Errorf("%w: %v", ErrDeleteRelation, err)
	}

	if result.RowsAffected() == 0 {
		r.logger.Warn(ctx, "Relação entre usuário e categoria não encontrada para exclusão", map[string]interface{}{
			"user_id":     userID,
			"category_id": categoryID,
		})
		return ErrRelationNotFound
	}

	r.logger.Info(ctx, "Relação entre usuário e categoria excluída com sucesso", map[string]interface{}{
		"user_id":     userID,
		"category_id": categoryID,
	})

	return nil
}

func (r *userCategoryRelationRepositories) DeleteAll(ctx context.Context, userID int64) error {
	const query = `
		DELETE FROM user_category_relations
		WHERE user_id = $1;
	`

	result, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		r.logger.Error(ctx, err, "Erro ao excluir todas as relações do usuário", map[string]interface{}{
			"user_id": userID,
		})
		return fmt.Errorf("%w: %v", ErrDeleteAllUserRelations, err)
	}

	r.logger.Info(ctx, "Todas as relações do usuário foram excluídas com sucesso", map[string]interface{}{
		"user_id":       userID,
		"rows_affected": result.RowsAffected(),
	})

	return nil
}
