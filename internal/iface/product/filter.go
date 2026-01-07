package iface

import (
	"context"

	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/product/filter"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
)

type ProductFilter interface {
	Filter(ctx context.Context, f *filter.ProductFilter) ([]*models.Product, error)
}
