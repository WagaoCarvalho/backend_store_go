package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/category"
)

type SupplierCategoryReader interface {
	GetByID(ctx context.Context, id int64) (*models.SupplierCategory, error)
	GetAll(ctx context.Context) ([]*models.SupplierCategory, error)
}

type SupplierCategoryWriter interface {
	Create(ctx context.Context, category *models.SupplierCategory) (*models.SupplierCategory, error)
	Update(ctx context.Context, category *models.SupplierCategory) error
	Delete(ctx context.Context, id int64) error
}
