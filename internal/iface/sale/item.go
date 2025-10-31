package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/item"
)

type SaleItemReader interface {
	GetByID(ctx context.Context, id int64) (*models.SaleItem, error)
	GetBySaleID(ctx context.Context, saleID int64, limit, offset int) ([]*models.SaleItem, error)
	GetByProductID(ctx context.Context, productID int64, limit, offset int) ([]*models.SaleItem, error)
}

type SaleItemWriter interface {
	Create(ctx context.Context, item *models.SaleItem) (*models.SaleItem, error)
	Update(ctx context.Context, item *models.SaleItem) error
	Delete(ctx context.Context, id int64) error
	DeleteBySaleID(ctx context.Context, saleID int64) error
}

type SaleItemChecker interface {
	ItemExists(ctx context.Context, id int64) (bool, error)
}
