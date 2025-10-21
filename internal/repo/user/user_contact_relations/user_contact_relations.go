package repo

import (
	"context"
	"fmt"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_contact_relations"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserContactRelationRepository interface {
	Create(ctx context.Context, relation *models.UserContactRelations) (*models.UserContactRelations, error)
	CreateTx(ctx context.Context, tx pgx.Tx, relation *models.UserContactRelations) (*models.UserContactRelations, error)
	HasUserContactRelation(ctx context.Context, userID, contactID int64) (bool, error)
	GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserContactRelations, error)
	Delete(ctx context.Context, userID, contactID int64) error
	DeleteAll(ctx context.Context, userID int64) error
}

type userContactRelationRepositories struct {
	db *pgxpool.Pool
}

func NewUserContactRelationRepositories(db *pgxpool.Pool) UserContactRelationRepository {
	return &userContactRelationRepositories{db: db}
}

func (r *userContactRelationRepositories) Create(ctx context.Context, relation *models.UserContactRelations) (*models.UserContactRelations, error) {
	const query = `
		INSERT INTO user_contact_relations (user_id, contact_id, created_at)
		VALUES ($1, $2, NOW());
	`

	_, err := r.db.Exec(ctx, query, relation.UserID, relation.ContactID)
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

func (r *userContactRelationRepositories) CreateTx(ctx context.Context, tx pgx.Tx, relation *models.UserContactRelations) (*models.UserContactRelations, error) {
	const query = `
		INSERT INTO user_contact_relations (user_id, contact_id, created_at)
		VALUES ($1, $2, NOW());
	`

	_, err := tx.Exec(ctx, query, relation.UserID, relation.ContactID)
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

func (r *userContactRelationRepositories) HasUserContactRelation(ctx context.Context, userID, contactID int64) (bool, error) {
	const query = `
		SELECT 1
		FROM user_contact_relations
		WHERE user_id = $1 AND contact_id = $2
		LIMIT 1;
	`

	var dummy int
	err := r.db.QueryRow(ctx, query, userID, contactID).Scan(&dummy)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return false, nil
		}
		return false, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return true, nil
}

func (r *userContactRelationRepositories) GetAllRelationsByUserID(ctx context.Context, userID int64) ([]*models.UserContactRelations, error) {
	const query = `
		SELECT user_id, contact_id, created_at
		FROM user_contact_relations
		WHERE user_id = $1;
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	var relations []*models.UserContactRelations
	for rows.Next() {
		var rel models.UserContactRelations
		if err := rows.Scan(&rel.UserID, &rel.ContactID, &rel.CreatedAt); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		relations = append(relations, &rel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}

	return relations, nil
}

func (r *userContactRelationRepositories) Delete(ctx context.Context, userID, contactID int64) error {
	const query = `
		DELETE FROM user_contact_relations
		WHERE user_id = $1 AND contact_id = $2;
	`

	_, err := r.db.Exec(ctx, query, userID, contactID)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}

func (r *userContactRelationRepositories) DeleteAll(ctx context.Context, userID int64) error {
	const query = `
		DELETE FROM user_contact_relations
		WHERE user_id = $1;
	`

	_, err := r.db.Exec(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("%w: %v", errMsg.ErrDelete, err)
	}

	return nil
}
