package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/client"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/filter"
)

type ClientCpfFilter interface {
	Filter(ctx context.Context, f *filter.ClientCpfFilter) ([]*models.ClientCpf, error)
}
