package repositories

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_category_relations"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserCategoryRelationRepository interface {
	Create(ctx context.Context, relation *models.UserCategoryRelations) (*models.UserCategoryRelations, error)
	CreateTx(ctx context.Context, tx pgx.Tx, relation *models.UserCategoryRelations) (*models.UserCategoryRelations, error)
	HasUserCategoryRelation(ctx context.Context, userID, categoryID int64) (bool, error)
	GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserCategoryRelations, error)
	Delete(ctx context.Context, userID, categoryID int64) error
	DeleteAll(ctx context.Context, userID int64) error
}

type userCategoryRelationRepositories struct {
	db     *pgxpool.Pool
	logger logger.LoggerAdapterInterface
}

func NewUserCategoryRelationRepositories(db *pgxpool.Pool, logger logger.LoggerAdapterInterface) UserCategoryRelationRepository {
	return &userCategoryRelationRepositories{db: db, logger: logger}
}

func (r *userCategoryRelationRepositories) Create(ctx context.Context, relation *models.UserCategoryRelations) (*models.UserCategoryRelations, error) {
	ref := "[userCategoryRelationRepositories - Create] - "
	r.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{
		"user_id":     relation.UserID,
		"category_id": relation.CategoryID,
	})

	const query = `
		INSERT INTO user_category_relations (user_id, category_id, created_at)
		VALUES ($1, $2, NOW());
	`

	_, err := r.db.Exec(ctx, query, relation.UserID, relation.CategoryID)
	if err != nil {
		switch {
		case errMsgPg.IsDuplicateKey(err):
			r.logger.Warn(ctx, ref+logger.LogForeignKeyHasExists, map[string]any{
				"user_id":     relation.UserID,
				"category_id": relation.CategoryID,
			})
			return nil, errMsg.ErrRelationExists

		case errMsgPg.IsForeignKeyViolation(err):
			r.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{
				"user_id":     relation.UserID,
				"category_id": relation.CategoryID,
			})
			return nil, errMsg.ErrInvalidForeignKey

		default:
			r.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
				"user_id":     relation.UserID,
				"category_id": relation.CategoryID,
			})
			return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
		}
	}

	r.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"user_id":     relation.UserID,
		"category_id": relation.CategoryID,
	})

	return relation, nil
}

func (r *userCategoryRelationRepositories) CreateTx(ctx context.Context, tx pgx.Tx, relation *models.UserCategoryRelations) (*models.UserCategoryRelations, error) {
	ref := "[userCategoryRelationRepositories - CreateTx] - "
	r.logger.Info(ctx, ref+logger.LogCreateInit, map[string]any{
		"user_id":     relation.UserID,
		"category_id": relation.CategoryID,
	})

	const query = `
		INSERT INTO user_category_relations (user_id, category_id, created_at)
		VALUES ($1, $2, NOW());
	`

	_, err := tx.Exec(ctx, query, relation.UserID, relation.CategoryID)
	if err != nil {
		switch {
		case errMsgPg.IsDuplicateKey(err):
			r.logger.Warn(ctx, ref+logger.LogForeignKeyHasExists, map[string]any{
				"user_id":     relation.UserID,
				"category_id": relation.CategoryID,
			})
			return nil, errMsg.ErrRelationExists

		case errMsgPg.IsForeignKeyViolation(err):
			r.logger.Warn(ctx, ref+logger.LogForeignKeyViolation, map[string]any{
				"user_id":     relation.UserID,
				"category_id": relation.CategoryID,
			})
			return nil, errMsg.ErrInvalidForeignKey

		default:
			r.logger.Error(ctx, err, ref+logger.LogCreateError, map[string]any{
				"user_id":     relation.UserID,
				"category_id": relation.CategoryID,
			})
			return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
		}
	}

	r.logger.Info(ctx, ref+logger.LogCreateSuccess, map[string]any{
		"user_id":     relation.UserID,
		"category_id": relation.CategoryID,
	})

	return relation, nil
}

func (r *userCategoryRelationRepositories) GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserCategoryRelations, error) {
	ref := "[userCategoryRelationRepositories - GetAllRelationsByUserID] - "
	r.logger.Info(ctx, ref+logger.LogGetInit, map[string]any{
		"user_id": userID,
	})

	const query = `
		SELECT user_id, category_id, created_at
		FROM user_category_relations
		WHERE user_id = $1;
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogGetError, map[string]any{
			"user_id": userID,
		})
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var relations []*models.UserCategoryRelations
	for rows.Next() {
		var rel models.UserCategoryRelations
		if err := rows.Scan(&rel.UserID, &rel.CategoryID, &rel.CreatedAt); err != nil {
			r.logger.Error(ctx, err, ref+logger.LogGetErrorScan, map[string]any{
				"user_id": userID,
			})
			return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		relations = append(relations, &rel)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error(ctx, err, ref+logger.LogIterateError, map[string]any{
			"user_id": userID,
		})
		return nil, fmt.Errorf("%w: %v", errMsg.ErrIterate, err)
	}

	r.logger.Info(ctx, ref+logger.LogGetSuccess, map[string]any{
		"user_id":         userID,
		"relations_count": len(relations),
	})

	return relations, nil
}

func (r *userCategoryRelationRepositories) HasUserCategoryRelation(ctx context.Context, userID, categoryID int64) (bool, error) {
	ref := "[userCategoryRelationRepositories - HasUserCategoryRelation] - "
	r.logger.Info(ctx, ref+logger.LogCheckInit, map[string]any{
		"user_id":     userID,
		"category_id": categoryID,
	})

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
			r.logger.Info(ctx, ref+logger.LogCheckNotFound, map[string]any{
				"user_id":     userID,
				"category_id": categoryID,
			})
			return false, nil
		}

		r.logger.Error(ctx, err, ref+logger.LogCheckError, map[string]any{
			"user_id":     userID,
			"category_id": categoryID,
		})
		return false, fmt.Errorf("%w: %v", errMsg.ErrRelationCheck, err)
	}

	r.logger.Info(ctx, ref+logger.LogCheckSuccess, map[string]any{
		"user_id":     userID,
		"category_id": categoryID,
	})

	return true, nil
}

func (r *userCategoryRelationRepositories) Delete(ctx context.Context, userID, categoryID int64) error {
	ref := "[userCategoryRelationRepositories - Delete] - "
	r.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{
		"user_id":     userID,
		"category_id": categoryID,
	})

	const query = `
		DELETE FROM user_category_relations 
		WHERE user_id = $1 AND category_id = $2;
	`

	result, err := r.db.Exec(ctx, query, userID, categoryID)
	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"user_id":     userID,
			"category_id": categoryID,
		})
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	if result.RowsAffected() == 0 {
		r.logger.Warn(ctx, ref+logger.LogNotFound, map[string]any{
			"user_id":     userID,
			"category_id": categoryID,
		})
		return errMsg.ErrNotFound
	}

	r.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"user_id":     userID,
		"category_id": categoryID,
	})

	return nil
}

func (r *userCategoryRelationRepositories) DeleteAll(ctx context.Context, userID int64) error {
	ref := "[userCategoryRelationRepositories - DeleteAll] - "
	r.logger.Info(ctx, ref+logger.LogDeleteInit, map[string]any{
		"user_id": userID,
	})

	const query = `
		DELETE FROM user_category_relations
		WHERE user_id = $1;
	`

	result, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		r.logger.Error(ctx, err, ref+logger.LogDeleteError, map[string]any{
			"user_id": userID,
		})
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	r.logger.Info(ctx, ref+logger.LogDeleteSuccess, map[string]any{
		"user_id":       userID,
		"rows_affected": result.RowsAffected(),
	})

	return nil
}
