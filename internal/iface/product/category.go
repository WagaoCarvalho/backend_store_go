package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category"
)

type ProductCategory interface {
	Create(ctx context.Context, category *models.ProductCategory) (*models.ProductCategory, error)
	GetByID(ctx context.Context, id int64) (*models.ProductCategory, error)
	GetAll(ctx context.Context) ([]*models.ProductCategory, error)
	Update(ctx context.Context, category *models.ProductCategory) error
	Delete(ctx context.Context, id int64) error
}
