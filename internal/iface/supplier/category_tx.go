package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/category_relation"
	"github.com/jackc/pgx/v5"
)

type SupplierCategoryRelationTx interface {
	CreateTx(ctx context.Context, tx pgx.Tx, relation *models.SupplierCategoryRelation) (*models.SupplierCategoryRelation, error)
}
