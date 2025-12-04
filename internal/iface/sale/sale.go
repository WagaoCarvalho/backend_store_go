package iface

import (
	"context"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
)

type SaleReader interface {
	GetByID(ctx context.Context, id int64) (*models.Sale, error)
	GetByClientID(ctx context.Context, clientID int64, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error)
	GetByUserID(ctx context.Context, userID int64, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error)
	GetByStatus(ctx context.Context, status string, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error)
	GetByDateRange(ctx context.Context, start, end time.Time, limit, offset int, orderBy, orderDir string) ([]*models.Sale, error)
}

type SaleWriter interface {
	Create(ctx context.Context, sale *models.Sale) (*models.Sale, error)
	Update(ctx context.Context, sale *models.Sale) error
	Delete(ctx context.Context, id int64) error
}

type SaleStatus interface {
	Activate(ctx context.Context, id int64) error
	Cancel(ctx context.Context, id int64) error
	Complete(ctx context.Context, id int64) error
	Returned(ctx context.Context, id int64) error
}
