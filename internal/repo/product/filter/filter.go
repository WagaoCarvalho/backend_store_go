package repo

import (
	"context"
	"fmt"
	"strings"

	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/product/filter"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

var allowedProductSortFields = map[string]string{
	"id":             "id",
	"product_name":   "product_name",
	"manufacturer":   "manufacturer",
	"sale_price":     "sale_price",
	"cost_price":     "cost_price",
	"stock_quantity": "stock_quantity",
	"status":         "status",
	"version":        "version",
	"created_at":     "created_at",
	"updated_at":     "updated_at",
}

func (r *productFilterRepo) Filter(ctx context.Context, filter *filter.ProductFilter) ([]*model.Product, error) {
	// Validação do filtro
	if filter == nil {
		return nil, errMsg.ErrInvalidFilter
	}

	base := filter.BaseFilter.WithDefaults()

	query := `
		SELECT
			id,
			supplier_id,
			product_name,
			manufacturer,
			product_description,
			cost_price,
			sale_price,
			stock_quantity,
			min_stock,
			max_stock,
			barcode,
			status,
			version,
			allow_discount,
			min_discount_percent,
			max_discount_percent,
			created_at,
			updated_at
		FROM products
		WHERE 1=1
	`

	args := []any{}
	argPos := 1

	// Filtros com ILIKE seguro
	if filter.ProductName != "" {
		query += fmt.Sprintf(" AND product_name ILIKE $%d", argPos)
		args = append(args, "%"+filter.ProductName+"%")
		argPos++
	}

	if filter.Manufacturer != "" {
		query += fmt.Sprintf(" AND manufacturer ILIKE $%d", argPos)
		args = append(args, "%"+filter.Manufacturer+"%")
		argPos++
	}

	if filter.Barcode != "" {
		query += fmt.Sprintf(" AND barcode = $%d", argPos)
		args = append(args, filter.Barcode)
		argPos++
	}

	if filter.SupplierID != nil {
		query += fmt.Sprintf(" AND supplier_id = $%d", argPos)
		args = append(args, *filter.SupplierID)
		argPos++
	}

	if filter.Status != nil {
		query += fmt.Sprintf(" AND status = $%d", argPos)
		args = append(args, *filter.Status)
		argPos++
	}

	if filter.AllowDiscount != nil {
		query += fmt.Sprintf(" AND allow_discount = $%d", argPos)
		args = append(args, *filter.AllowDiscount)
		argPos++
	}

	if filter.Version != nil {
		query += fmt.Sprintf(" AND version = $%d", argPos)
		args = append(args, *filter.Version)
		argPos++
	}

	if filter.MinCostPrice != nil {
		query += fmt.Sprintf(" AND cost_price >= $%d", argPos)
		args = append(args, *filter.MinCostPrice)
		argPos++
	}

	if filter.MaxCostPrice != nil {
		query += fmt.Sprintf(" AND cost_price <= $%d", argPos)
		args = append(args, *filter.MaxCostPrice)
		argPos++
	}

	if filter.MinSalePrice != nil {
		query += fmt.Sprintf(" AND sale_price >= $%d", argPos)
		args = append(args, *filter.MinSalePrice)
		argPos++
	}

	if filter.MaxSalePrice != nil {
		query += fmt.Sprintf(" AND sale_price <= $%d", argPos)
		args = append(args, *filter.MaxSalePrice)
		argPos++
	}

	if filter.MinStockQuantity != nil {
		query += fmt.Sprintf(" AND stock_quantity >= $%d", argPos)
		args = append(args, *filter.MinStockQuantity)
		argPos++
	}

	if filter.MaxStockQuantity != nil {
		query += fmt.Sprintf(" AND stock_quantity <= $%d", argPos)
		args = append(args, *filter.MaxStockQuantity)
		argPos++
	}

	if filter.CreatedFrom != nil {
		query += fmt.Sprintf(" AND created_at >= $%d", argPos)
		args = append(args, *filter.CreatedFrom)
		argPos++
	}

	if filter.CreatedTo != nil {
		query += fmt.Sprintf(" AND created_at <= $%d", argPos)
		args = append(args, *filter.CreatedTo)
		argPos++
	}

	if filter.UpdatedFrom != nil {
		query += fmt.Sprintf(" AND updated_at >= $%d", argPos)
		args = append(args, *filter.UpdatedFrom)
		argPos++
	}

	if filter.UpdatedTo != nil {
		query += fmt.Sprintf(" AND updated_at <= $%d", argPos)
		args = append(args, *filter.UpdatedTo)
		argPos++
	}

	// ORDER BY seguro
	sortField := "created_at"
	if filter.SortBy != "" {
		if v, ok := allowedProductSortFields[strings.ToLower(filter.SortBy)]; ok {
			sortField = v
		}
	}

	sortOrder := "ASC"
	if filter.SortOrder != "" {
		if strings.ToUpper(filter.SortOrder) == "DESC" {
			sortOrder = "DESC"
		}
	}

	query += fmt.Sprintf(" ORDER BY %s %s LIMIT %d OFFSET %d",
		sortField,
		sortOrder,
		base.Limit,
		base.Offset,
	)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w [filtro=%+v]: %v", errMsg.ErrGet, filter, err)
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(
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
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		products = append(products, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrIterate, err)
	}

	// Retorna slice vazio, não nil, para consistência
	if products == nil {
		return []*model.Product{}, nil
	}

	return products, nil
}
