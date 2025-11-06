package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/contact_relation"
	"github.com/jackc/pgx/v5"
)

type UserContactRelationTx interface {
	CreateTx(ctx context.Context, tx pgx.Tx, relation *models.UserContactRelation) (*models.UserContactRelation, error)
}
