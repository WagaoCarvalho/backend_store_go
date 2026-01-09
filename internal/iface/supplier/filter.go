package iface

import (
	"context"

	modelFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/filter"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
)

type SupplierFilter interface {
	Filter(ctx context.Context, f *modelFilter.SupplierFilter) ([]*models.Supplier, error)
}
