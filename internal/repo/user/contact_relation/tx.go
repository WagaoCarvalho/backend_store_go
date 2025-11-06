package repo

import (
	"context"
	"fmt"

	ifaceTx "github.com/WagaoCarvalho/backend_store_go/internal/iface/user"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/contact_relation"
	errMsgPg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/repo"
	"github.com/jackc/pgx/v5"
)

type userContactRelationTx struct {
	db repo.DBExecutor
}

func NewUserContactRelationTx(db repo.DBExecutor) ifaceTx.UserContactRelationTx {
	return &userContactRelationTx{db: db}
}

func (r *userContactRelationTx) CreateTx(ctx context.Context, tx pgx.Tx, relation *models.UserContactRelation) (*models.UserContactRelation, error) {
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
