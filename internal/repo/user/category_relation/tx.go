package repo

import (
	"context"
	"fmt"

	ifaceTx "github.com/WagaoCarvalho/backend_store_go/internal/iface/user"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/db"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/category_relation"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/jackc/pgx/v5"
)

type userCategoryRelationTx struct {
	db repo.DBExecutor
}

func NewUserCategoryRelationTx(db repo.DBExecutor) ifaceTx.UserCategoryRelationTx {
	return &userCategoryRelationTx{db: db}
}

func (r *userCategoryRelationTx) CreateTx(ctx context.Context, tx pgx.Tx, relation *models.UserCategoryRelation) (*models.UserCategoryRelation, error) {
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
			return nil, fmt.Errorf("userCategoryRelationTx: %w", errMsg.ErrDBInvalidForeignKey)
		default:
			return nil, fmt.Errorf("%w: %v", errMsg.ErrCreate, err)
		}
	}

	return relation, nil
}
