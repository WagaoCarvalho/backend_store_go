package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address/address"
)

type AddressFilter interface {
	Filter(ctx context.Context, f *models.Address) ([]*models.Address, error)
}
