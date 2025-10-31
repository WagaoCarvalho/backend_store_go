package repo

import (
	"context"
	"errors"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/category_relation"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserCategoryRelation interface {
	Create(ctx context.Context, relation *models.UserCategoryRelation) (*models.UserCategoryRelation, error)
	CreateTx(ctx context.Context, tx pgx.Tx, relation *models.UserCategoryRelation) (*models.UserCategoryRelation, error)
	HasUserCategoryRelation(ctx context.Context, userID, categoryID int64) (bool, error)
	GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserCategoryRelation, error)
	Delete(ctx context.Context, userID, categoryID int64) error
	DeleteAll(ctx context.Context, userID int64) error
}

type userCategoryRelation struct {
	db *pgxpool.Pool
}

func NewUserCategoryRelation(db *pgxpool.Pool) UserCategoryRelation {
	return &userCategoryRelation{db: db}
}

func (r *userCategoryRelation) Create(ctx context.Context, relation *models.UserCategoryRelation) (*models.UserCategoryRelation, error) {
	const query = `
		INSERT INTO user_category_relations (user_id, category_id, created_at)
		VALUES ($1, $2, NOW());
	`

	_, err := r.db.Exec(ctx, query, relation.UserID, relation.CategoryID)
	if err != nil {
		switch {
		case errMsgPg.IsDuplicateKey(err):
			return nil, errMsg.ErrRelationExists
		case errMsgPg.IsForeignKeyViolation(err):
			return nil, errMsg.ErrDBInvalidForeignKey
		default:
			return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
		}
	}

	return relation, nil
}

func (r *userCategoryRelation) CreateTx(ctx context.Context, tx pgx.Tx, relation *models.UserCategoryRelation) (*models.UserCategoryRelation, error) {
	const query = `
		INSERT INTO user_category_relations (user_id, category_id, created_at)
		VALUES ($1, $2, NOW());
	`

	_, err := tx.Exec(ctx, query, relation.UserID, relation.CategoryID)
	if err != nil {
		switch {
		case errMsgPg.IsDuplicateKey(err):
			return nil, errMsg.ErrRelationExists
		case errMsgPg.IsForeignKeyViolation(err):
			return nil, errMsg.ErrDBInvalidForeignKey
		default:
			return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
		}
	}

	return relation, nil
}

func (r *userCategoryRelation) GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserCategoryRelation, error) {
	if userID <= 0 {
		return []*models.UserCategoryRelation{}, errMsg.ErrZeroID
	}

	const query = `
        SELECT user_id, category_id, created_at
        FROM user_category_relations
        WHERE user_id = $1;
    `

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return []*models.UserCategoryRelation{}, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var relations []*models.UserCategoryRelation
	for rows.Next() {
		var rel models.UserCategoryRelation
		if err := rows.Scan(&rel.UserID, &rel.CategoryID, &rel.CreatedAt); err != nil {
			return []*models.UserCategoryRelation{}, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		relations = append(relations, &rel)
	}

	if err := rows.Err(); err != nil {
		return []*models.UserCategoryRelation{}, fmt.Errorf("%w: %v", errMsg.ErrIterate, err)
	}

	// Garantir que nunca retorne nil
	if relations == nil {
		relations = []*models.UserCategoryRelation{}
	}

	return relations, nil
}

func (r *userCategoryRelation) HasUserCategoryRelation(ctx context.Context, userID, categoryID int64) (bool, error) {
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
			return false, nil
		}
		return false, fmt.Errorf("%w: %v", errMsg.ErrRelationCheck, err)
	}

	return true, nil
}

func (r *userCategoryRelation) Delete(ctx context.Context, userID, categoryID int64) error {
	const query = `
		DELETE FROM user_category_relations 
		WHERE user_id = $1 AND category_id = $2;
	`

	result, err := r.db.Exec(ctx, query, userID, categoryID)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	if result.RowsAffected() == 0 {
		return errMsg.ErrNotFound
	}

	return nil
}

func (r *userCategoryRelation) DeleteAll(ctx context.Context, userID int64) error {
	const query = `
		DELETE FROM user_category_relations
		WHERE user_id = $1;
	`

	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}
