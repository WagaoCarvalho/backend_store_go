package iface

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
)

type ProductReader interface {
	GetByID(ctx context.Context, id int64) (*models.Product, error)
	GetVersionByID(ctx context.Context, id int64) (int64, error)
}

type ProductGetAll interface {
	GetAll(ctx context.Context) ([]*models.Product, error)
}

type ProductWriter interface {
	Create(ctx context.Context, product *models.Product) (*models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id int64) error
}

type ProductFilter interface {
	Filter(ctx context.Context, f *models.ProductFilter) ([]*models.Product, error)
}

type ProductStock interface {
	UpdateStock(ctx context.Context, id int64, quantity int) error
	IncreaseStock(ctx context.Context, id int64, amount int) error
	DecreaseStock(ctx context.Context, id int64, amount int) error
	GetStock(ctx context.Context, id int64) (int, error)
}

type ProductDiscount interface {
	EnableDiscount(ctx context.Context, id int64) error
	DisableDiscount(ctx context.Context, id int64) error
	ApplyDiscount(ctx context.Context, id int64, percent float64) (*models.Product, error)
}

type ProductStatus interface {
	EnableProduct(ctx context.Context, uid int64) error
	DisableProduct(ctx context.Context, uid int64) error
}

type ProductChecker interface {
	ProductExists(ctx context.Context, productID int64) (bool, error)
}
