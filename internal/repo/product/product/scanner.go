package repo

import (
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	"github.com/jackc/pgx/v5"
)

func scanProductRow(row pgx.Row, p *models.Product) error {
	return row.Scan(
		&p.ID,
		&p.SupplierID,
		&p.ProductName,
		&p.Manufacturer,
		&p.Description,
		&p.CostPrice,
		&p.SalePrice,
		&p.StockQuantity,
		&p.MinStock,
		&p.MaxStock,
		&p.Barcode,
		&p.Status,
		&p.Version,
		&p.AllowDiscount,
		&p.MinDiscountPercent,
		&p.MaxDiscountPercent,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
}

func scanProductRowLimited(row pgx.Row, p *models.Product) error {
	return row.Scan(
		&p.ID,
		&p.SupplierID,
		&p.ProductName,
		&p.Manufacturer,
		&p.Description,
		&p.CostPrice,
		&p.SalePrice,
		&p.StockQuantity,
		&p.Barcode,
		&p.Status,
		&p.AllowDiscount,
		&p.MaxDiscountPercent,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
}

func scanProductDiscountRow(row pgx.Row, p *models.Product) error {
	return row.Scan(
		&p.ID,
		&p.ProductName,
		&p.SalePrice,
		&p.MaxDiscountPercent,
		&p.AllowDiscount,
		&p.Version,
		&p.UpdatedAt,
	)
}
