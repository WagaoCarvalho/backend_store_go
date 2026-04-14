package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address/address"
)

type AddressReader interface {
	GetByID(ctx context.Context, id int64) (*models.Address, error)
}

type AddressWriter interface {
	Create(ctx context.Context, address *models.Address) (*models.Address, error)
	Update(ctx context.Context, address *models.Address) error
	Delete(ctx context.Context, id int64) error
}

type AddressStatus interface {
	Disable(ctx context.Context, uid int64) error
	Enable(ctx context.Context, uid int64) error
}
