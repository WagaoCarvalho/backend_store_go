package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	"github.com/jackc/pgx/v5"
)

type ContactTx interface {
	CreateTx(ctx context.Context, tx pgx.Tx, contact *models.Contact) (*models.Contact, error)
}
