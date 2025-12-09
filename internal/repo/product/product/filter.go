package repo

import (
	"context"
	"fmt"

	modelFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/product/filter"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	commonFilter "github.com/WagaoCarvalho/backend_store_go/internal/repo/common/filter"
)

// ----------------------------
// Abstrações de banco e scanner
// ----------------------------

type scanner interface {
	Next() bool
	Scan(dest ...any) error
	Err() error
	Close()
}

// Filter retorna produtos filtrados e paginados usando o builder genérico.
func (r *productRepo) Filter(ctx context.Context, filterData *modelFilter.ProductFilter) ([]*model.Product, error) {
	qb := commonFilter.NewSQLQueryBuilder(
		"products",
		[]string{
			"id", "supplier_id", "product_name", "manufacturer", "product_description",
			"cost_price", "sale_price", "stock_quantity", "min_stock", "max_stock",
			"barcode", "status", "version", "allow_discount",
			"min_discount_percent", "max_discount_percent",
			"created_at", "updated_at",
		},
		"created_at DESC",
	)

	// Criação dos filtros dinâmicos genéricos
	filters := []commonFilter.FilterCondition{
		&commonFilter.TextFilter{Field: "product_name", Value: filterData.ProductName},
		&commonFilter.TextFilter{Field: "manufacturer", Value: filterData.Manufacturer},
	}

	// Filtro para barcode (se necessário)
	if filterData.Barcode != "" {
		filters = append(filters, &commonFilter.TextFilter{Field: "barcode", Value: filterData.Barcode})
	}

	// Filtros com ponteiros usando os tipos genéricos
	if filterData.SupplierID != nil {
		filters = append(filters, &commonFilter.EqualFilter[int64]{Field: "supplier_id", Value: filterData.SupplierID})
	}

	if filterData.Status != nil {
		filters = append(filters, &commonFilter.EqualFilter[bool]{Field: "status", Value: filterData.Status})
	}

	if filterData.Version != nil {
		filters = append(filters, &commonFilter.EqualFilter[int]{Field: "version", Value: filterData.Version})
	}

	if filterData.AllowDiscount != nil {
		filters = append(filters, &commonFilter.EqualFilter[bool]{Field: "allow_discount", Value: filterData.AllowDiscount})
	}

	// Filtros de range para preços
	if filterData.MinCostPrice != nil || filterData.MaxCostPrice != nil {
		filters = append(filters, &commonFilter.RangeFilter[float64]{
			FieldMin: "cost_price",
			FieldMax: "cost_price",
			Min:      filterData.MinCostPrice,
			Max:      filterData.MaxCostPrice,
		})
	}

	if filterData.MinSalePrice != nil || filterData.MaxSalePrice != nil {
		filters = append(filters, &commonFilter.RangeFilter[float64]{
			FieldMin: "sale_price",
			FieldMax: "sale_price",
			Min:      filterData.MinSalePrice,
			Max:      filterData.MaxSalePrice,
		})
	}

	// Filtros de range para estoque - usando int que é o tipo correto
	if filterData.MinStockQuantity != nil || filterData.MaxStockQuantity != nil {
		filters = append(filters, &commonFilter.RangeFilter[int]{
			FieldMin: "stock_quantity",
			FieldMax: "stock_quantity",
			Min:      filterData.MinStockQuantity,
			Max:      filterData.MaxStockQuantity,
		})
	}

	// Filtros de range para desconto
	if filterData.MinDiscountPercent != nil || filterData.MaxDiscountPercent != nil {
		filters = append(filters, &commonFilter.RangeFilter[float64]{
			FieldMin: "min_discount_percent",
			FieldMax: "max_discount_percent",
			Min:      filterData.MinDiscountPercent,
			Max:      filterData.MaxDiscountPercent,
		})
	}

	// Filtros de data - convertendo time.Time para string no formato ISO
	if filterData.CreatedFrom != nil || filterData.CreatedTo != nil {
		var createdFromStr, createdToStr *string

		if filterData.CreatedFrom != nil {
			fromStr := filterData.CreatedFrom.Format("2006-01-02 15:04:05")
			createdFromStr = &fromStr
		}
		if filterData.CreatedTo != nil {
			toStr := filterData.CreatedTo.Format("2006-01-02 15:04:05")
			createdToStr = &toStr
		}

		filters = append(filters, &commonFilter.RangeFilter[string]{
			FieldMin: "created_at",
			FieldMax: "created_at",
			Min:      createdFromStr,
			Max:      createdToStr,
		})
	}

	if filterData.UpdatedFrom != nil || filterData.UpdatedTo != nil {
		var updatedFromStr, updatedToStr *string

		if filterData.UpdatedFrom != nil {
			fromStr := filterData.UpdatedFrom.Format("2006-01-02 15:04:05")
			updatedFromStr = &fromStr
		}
		if filterData.UpdatedTo != nil {
			toStr := filterData.UpdatedTo.Format("2006-01-02 15:04:05")
			updatedToStr = &toStr
		}

		filters = append(filters, &commonFilter.RangeFilter[string]{
			FieldMin: "updated_at",
			FieldMax: "updated_at",
			Min:      updatedFromStr,
			Max:      updatedToStr,
		})
	}

	// Aplicar todos os filtros
	for _, f := range filters {
		f.Apply(qb)
	}

	query, args := qb.Build(filterData.Limit, filterData.Offset)

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrGet, err)
	}
	defer rows.Close()

	return scanProducts(rows)
}

// ----------------------------
// Scaneamento separado para testabilidade
// ----------------------------

func scanProducts(rows scanner) ([]*model.Product, error) {
	var products []*model.Product

	for rows.Next() {
		var p model.Product
		if err := rows.Scan(
			&p.ID, &p.SupplierID, &p.ProductName, &p.Manufacturer,
			&p.Description, &p.CostPrice, &p.SalePrice, &p.StockQuantity,
			&p.MinStock, &p.MaxStock, &p.Barcode, &p.Status,
			&p.Version, &p.AllowDiscount, &p.MinDiscountPercent,
			&p.MaxDiscountPercent, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("%w: %v", errMsg.ErrScan, err)
		}
		products = append(products, &p)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%w: %v", errMsg.ErrIterate, err)
	}

	return products, nil
}
