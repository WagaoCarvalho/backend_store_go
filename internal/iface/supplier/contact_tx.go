package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/contact_relation"
	"github.com/jackc/pgx/v5"
)

type SupplierContactRelationTx interface {
	CreateTx(ctx context.Context, tx pgx.Tx, relation *models.SupplierContactRelation) (*models.SupplierContactRelation, error)
}
