package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address/address"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/address/filter"
)

type AddressFilter interface {
	Filter(ctx context.Context, filter *filter.AddressFilter) ([]*models.Address, error)
}
