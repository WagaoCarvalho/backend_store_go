package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/category_relation"
)

type SupplierCategoryRelationReader interface {
	GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.SupplierCategoryRelation, error)
	GetByCategoryID(ctx context.Context, categoryID int64) ([]*models.SupplierCategoryRelation, error)
	HasRelation(ctx context.Context, supplierID, categoryID int64) (bool, error)
}

type SupplierCategoryRelationWriter interface {
	Create(ctx context.Context, relation *models.SupplierCategoryRelation) (*models.SupplierCategoryRelation, error)
	Delete(ctx context.Context, supplierID, categoryID int64) error
	DeleteAllBySupplierID(ctx context.Context, supplierID int64) error
}
