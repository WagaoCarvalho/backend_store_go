package dto

import (
	"strconv"
	"time"

	modelFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/filter"
	modelProduct "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
)

type ProductFilterDTO struct {
	ProductName        string  `schema:"product_name"`
	Manufacturer       string  `schema:"manufacturer"`
	Barcode            string  `schema:"barcode"`
	Status             *bool   `schema:"status"`
	SupplierID         *int64  `schema:"supplier_id"`
	Version            *int    `schema:"version"`
	MinCostPrice       *string `schema:"min_cost_price"`
	MaxCostPrice       *string `schema:"max_cost_price"`
	MinSalePrice       *string `schema:"min_sale_price"`
	MaxSalePrice       *string `schema:"max_sale_price"`
	MinStockQuantity   *string `schema:"min_stock_quantity"`
	MaxStockQuantity   *string `schema:"max_stock_quantity"`
	AllowDiscount      *bool   `schema:"allow_discount"`
	MinDiscountPercent *string `schema:"min_discount_percent"`
	MaxDiscountPercent *string `schema:"max_discount_percent"`
	CreatedFrom        *string `schema:"created_from"`
	CreatedTo          *string `schema:"created_to"`
	UpdatedFrom        *string `schema:"updated_from"`
	UpdatedTo          *string `schema:"updated_to"`
	Limit              int     `schema:"limit"`
	Offset             int     `schema:"offset"`
}

func (d *ProductFilterDTO) ToModel() (*modelProduct.ProductFilter, error) {
	parseDate := func(s *string) *time.Time {
		if s == nil || *s == "" {
			return nil
		}
		t, err := time.Parse("2006-01-02", *s)
		if err != nil {
			return nil
		}
		return &t
	}

	parseFloat := func(s *string) *float64 {
		if s == nil || *s == "" {
			return nil
		}
		f, err := strconv.ParseFloat(*s, 64)
		if err != nil {
			return nil
		}
		return &f
	}

	parseInt := func(s *string) *int {
		if s == nil || *s == "" {
			return nil
		}
		i, err := strconv.Atoi(*s)
		if err != nil {
			return nil
		}
		return &i
	}

	filter := &modelProduct.ProductFilter{
		BaseFilter: modelFilter.BaseFilter{
			Limit:  d.Limit,
			Offset: d.Offset,
		},
		ProductName:        d.ProductName,
		Manufacturer:       d.Manufacturer,
		Barcode:            d.Barcode,
		Status:             d.Status,
		SupplierID:         d.SupplierID,
		Version:            d.Version,
		MinCostPrice:       parseFloat(d.MinCostPrice),
		MaxCostPrice:       parseFloat(d.MaxCostPrice),
		MinSalePrice:       parseFloat(d.MinSalePrice),
		MaxSalePrice:       parseFloat(d.MaxSalePrice),
		MinStockQuantity:   parseInt(d.MinStockQuantity),
		MaxStockQuantity:   parseInt(d.MaxStockQuantity),
		AllowDiscount:      d.AllowDiscount,
		MinDiscountPercent: parseFloat(d.MinDiscountPercent),
		MaxDiscountPercent: parseFloat(d.MaxDiscountPercent),
		CreatedFrom:        parseDate(d.CreatedFrom),
		CreatedTo:          parseDate(d.CreatedTo),
		UpdatedFrom:        parseDate(d.UpdatedFrom),
		UpdatedTo:          parseDate(d.UpdatedTo),
	}

	return filter, nil
}
