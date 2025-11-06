package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category"
)

type ProductCategoryWriter interface {
	Create(ctx context.Context, category *models.ProductCategory) (*models.ProductCategory, error)

	Update(ctx context.Context, category *models.ProductCategory) error
	Delete(ctx context.Context, id int64) error
}

type ProductCategoryReader interface {
	GetByID(ctx context.Context, id int64) (*models.ProductCategory, error)
	GetAll(ctx context.Context) ([]*models.ProductCategory, error)
}
