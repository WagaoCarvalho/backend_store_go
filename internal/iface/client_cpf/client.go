package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/client"
)

type ClientCpfReader interface {
	GetByID(ctx context.Context, id int64) (*models.ClientCpf, error)
}

type ClientCpfWriter interface {
	Create(ctx context.Context, clientCpf *models.ClientCpf) (*models.ClientCpf, error)
	Update(ctx context.Context, clientCpf *models.ClientCpf) error
	Delete(ctx context.Context, id int64) error
}

type ClientCpfStatus interface {
	Disable(ctx context.Context, id int64) error
	Enable(ctx context.Context, id int64) error
}

type ClientCpfChecker interface {
	ClientCpfExists(ctx context.Context, clientCpfID int64) (bool, error)
}

type ClientCpfVersion interface {
	GetVersionByID(ctx context.Context, id int64) (int, error)
}
