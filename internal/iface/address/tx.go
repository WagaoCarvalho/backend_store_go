package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	"github.com/jackc/pgx/v5"
)

type AddressTx interface {
	CreateTx(ctx context.Context, tx pgx.Tx, address *models.Address) (*models.Address, error)
}
