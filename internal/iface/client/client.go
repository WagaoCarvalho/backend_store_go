package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client"
)

type ClientReader interface {
	GetByID(ctx context.Context, id int64) (*models.Client, error)
	GetByName(ctx context.Context, name string) ([]*models.Client, error)
	GetVersionByID(ctx context.Context, id int64) (int, error)
	ClientExists(ctx context.Context, clientID int64) (bool, error)
}

type ClientWriter interface {
	Create(ctx context.Context, client *models.Client) (*models.Client, error)
	Update(ctx context.Context, client *models.Client) error
	Delete(ctx context.Context, id int64) error
}

type ClientStatus interface {
	Disable(ctx context.Context, id int64) error
	Enable(ctx context.Context, id int64) error
}

type ClientFilter interface {
	GetAll(ctx context.Context, f *models.ClientFilter) ([]*models.Client, error)
}
