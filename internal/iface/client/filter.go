package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/client/filter"
)

type ClientFilter interface {
	GetAll(ctx context.Context, f *filter.ClientFilter) ([]*models.Client, error)
}
