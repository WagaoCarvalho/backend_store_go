package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category_relation"
)

type ProductCategoryRelationReader interface {
	GetAllRelationsByProductID(ctx context.Context, productID int64) ([]*models.ProductCategoryRelation, error)
}

type ProductCategoryRelationWriter interface {
	Create(ctx context.Context, relation *models.ProductCategoryRelation) (*models.ProductCategoryRelation, error)
	Delete(ctx context.Context, productID, categoryID int64) error
	DeleteAll(ctx context.Context, productID int64) error
}

type ProductCategoryRelationChecker interface {
	HasProductCategoryRelation(ctx context.Context, productID, categoryID int64) (bool, error)
}
