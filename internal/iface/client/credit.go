package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client/credit"
)

type ClientCreditReader interface {
	GetByName(ctx context.Context, name string) ([]*models.ClientCredit, error)
	GetVersionByID(ctx context.Context, id int64) (int, error)
	GetAll(ctx context.Context) ([]*models.ClientCredit, error)
}

type ClientCreditWriter interface {
	Create(ctx context.Context, client *models.ClientCredit) (*models.ClientCredit, error)
	Update(ctx context.Context, client *models.ClientCredit) error
	Delete(ctx context.Context, id int64) error
}

type ClientCreditStatus interface {
	Disable(ctx context.Context, id int64) error
	Enable(ctx context.Context, id int64) error
}
