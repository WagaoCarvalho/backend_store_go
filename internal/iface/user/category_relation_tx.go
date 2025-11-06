package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/category_relation"
	"github.com/jackc/pgx/v5"
)

type UserCategoryRelationTx interface {
	CreateTx(ctx context.Context, tx pgx.Tx, relation *models.UserCategoryRelation) (*models.UserCategoryRelation, error)
}
