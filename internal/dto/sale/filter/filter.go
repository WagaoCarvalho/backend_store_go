package dto

import (
	"strconv"
	"time"

	modelFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	modelSale "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/filter"
)

type SaleFilterDTO struct {
	ClientID         *int64  `schema:"client_id"`
	UserID           *int64  `schema:"user_id"`
	PaymentType      string  `schema:"payment_type"`
	Status           string  `schema:"status"`
	Notes            string  `schema:"notes"`
	MinTotalAmount   *string `schema:"min_total_amount"`
	MaxTotalAmount   *string `schema:"max_total_amount"`
	MinItemsAmount   *string `schema:"min_items_amount"`
	MaxItemsAmount   *string `schema:"max_items_amount"`
	MinItemsDiscount *string `schema:"min_items_discount"`
	MaxItemsDiscount *string `schema:"max_items_discount"`
	MinSaleDiscount  *string `schema:"min_sale_discount"`
	MaxSaleDiscount  *string `schema:"max_sale_discount"`
	SaleDateFrom     *string `schema:"sale_date_from"`
	SaleDateTo       *string `schema:"sale_date_to"`
	CreatedFrom      *string `schema:"created_from"`
	CreatedTo        *string `schema:"created_to"`
	UpdatedFrom      *string `schema:"updated_from"`
	UpdatedTo        *string `schema:"updated_to"`
	Limit            int     `schema:"limit"`
	Offset           int     `schema:"offset"`
}

func (d *SaleFilterDTO) ToModel() (*modelSale.SaleFilter, error) {
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

	// Criar BaseFilter com Version se necessário
	baseFilter := modelFilter.BaseFilter{
		Limit:  d.Limit,
		Offset: d.Offset,
	}

	filter := &modelSale.SaleFilter{
		BaseFilter: baseFilter,

		ClientID:    d.ClientID,
		UserID:      d.UserID,
		PaymentType: d.PaymentType,
		Status:      d.Status,
		Notes:       d.Notes,

		MinTotalAmount: parseFloat(d.MinTotalAmount),
		MaxTotalAmount: parseFloat(d.MaxTotalAmount),

		MinTotalItemsAmount: parseFloat(d.MinItemsAmount),
		MaxTotalItemsAmount: parseFloat(d.MaxItemsAmount),

		MinTotalItemsDiscount: parseFloat(d.MinItemsDiscount),
		MaxTotalItemsDiscount: parseFloat(d.MaxItemsDiscount),

		MinTotalSaleDiscount: parseFloat(d.MinSaleDiscount),
		MaxTotalSaleDiscount: parseFloat(d.MaxSaleDiscount),

		SaleDateFrom: parseDate(d.SaleDateFrom),
		SaleDateTo:   parseDate(d.SaleDateTo),

		CreatedFrom: parseDate(d.CreatedFrom),
		CreatedTo:   parseDate(d.CreatedTo),
		UpdatedFrom: parseDate(d.UpdatedFrom),
		UpdatedTo:   parseDate(d.UpdatedTo),
	}

	return filter, nil
}
