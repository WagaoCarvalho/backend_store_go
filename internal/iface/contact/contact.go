package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
)

type ContactReader interface {
	GetByID(ctx context.Context, id int64) (*models.Contact, error)
}

type ContactWriter interface {
	Create(ctx context.Context, contact *models.Contact) (*models.Contact, error)
	Update(ctx context.Context, contact *models.Contact) error
	Delete(ctx context.Context, id int64) error
}
