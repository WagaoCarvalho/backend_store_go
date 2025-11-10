package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
)

type SupplierReader interface {
	GetByID(ctx context.Context, id int64) (*models.Supplier, error)
	GetByName(ctx context.Context, name string) ([]*models.Supplier, error)
	GetVersionByID(ctx context.Context, id int64) (int64, error)
	GetAll(ctx context.Context) ([]*models.Supplier, error)
	SupplierExists(ctx context.Context, supplierID int64) (bool, error)
}

type SupplierWriter interface {
	Create(ctx context.Context, supplier *models.Supplier) (*models.Supplier, error)

	Update(ctx context.Context, supplier *models.Supplier) error
	Delete(ctx context.Context, id int64) error
}

type SupplierStatus interface {
	Disable(ctx context.Context, id int64) error
	Enable(ctx context.Context, id int64) error
}
