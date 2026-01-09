package iface

import (
	"context"

	modelFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/filter"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
)

type SaleFilter interface {
	Filter(ctx context.Context, f *modelFilter.SaleFilter) ([]*models.Sale, error)
}
