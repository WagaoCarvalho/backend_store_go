package dto

import (
	"fmt"
	"strconv"
	"time"

	modelFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	modelSale "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/filter"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
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
	// Função para parsear datas com validação
	parseDate := func(s *string, fieldName string) (*time.Time, error) {
		if s == nil || *s == "" {
			return nil, nil
		}
		t, err := time.Parse("2006-01-02", *s)
		if err != nil {
			return nil, fmt.Errorf("%w: campo '%s' com valor inválido '%s' - formato esperado: YYYY-MM-DD",
				errMsg.ErrInvalidFilter, fieldName, *s)
		}
		return &t, nil
	}

	// Função para parsear valores float com validação
	parseFloat := func(s *string, fieldName string, allowNegative bool) (*float64, error) {
		if s == nil || *s == "" {
			return nil, nil
		}
		f, err := strconv.ParseFloat(*s, 64)
		if err != nil {
			return nil, fmt.Errorf("%w: campo '%s' com valor inválido '%s' - valor numérico esperado",
				errMsg.ErrInvalidFilter, fieldName, *s)
		}
		if !allowNegative && f < 0 {
			return nil, fmt.Errorf("%w: campo '%s' não pode ser negativo",
				errMsg.ErrInvalidFilter, fieldName)
		}
		return &f, nil
	}

	// Validar parâmetros de paginação
	if d.Limit < 1 {
		return nil, fmt.Errorf("%w: 'limit' deve ser maior que 0", errMsg.ErrInvalidFilter)
	}
	if d.Limit > 100 {
		return nil, fmt.Errorf("%w: 'limit' máximo é 100", errMsg.ErrInvalidFilter)
	}
	if d.Offset < 0 {
		return nil, fmt.Errorf("%w: 'offset' não pode ser negativo", errMsg.ErrInvalidFilter)
	}

	baseFilter := modelFilter.BaseFilter{
		Limit:  d.Limit,
		Offset: d.Offset,
	}

	// Parsear e validar datas
	saleDateFrom, err := parseDate(d.SaleDateFrom, "sale_date_from")
	if err != nil {
		return nil, err
	}

	saleDateTo, err := parseDate(d.SaleDateTo, "sale_date_to")
	if err != nil {
		return nil, err
	}

	// Validar intervalo de datas (from <= to)
	if saleDateFrom != nil && saleDateTo != nil && saleDateFrom.After(*saleDateTo) {
		return nil, fmt.Errorf("%w: 'sale_date_from' não pode ser depois de 'sale_date_to'",
			errMsg.ErrInvalidFilter)
	}

	createdFrom, err := parseDate(d.CreatedFrom, "created_from")
	if err != nil {
		return nil, err
	}

	createdTo, err := parseDate(d.CreatedTo, "created_to")
	if err != nil {
		return nil, err
	}

	if createdFrom != nil && createdTo != nil && createdFrom.After(*createdTo) {
		return nil, fmt.Errorf("%w: 'created_from' não pode ser depois de 'created_to'",
			errMsg.ErrInvalidFilter)
	}

	updatedFrom, err := parseDate(d.UpdatedFrom, "updated_from")
	if err != nil {
		return nil, err
	}

	updatedTo, err := parseDate(d.UpdatedTo, "updated_to")
	if err != nil {
		return nil, err
	}

	if updatedFrom != nil && updatedTo != nil && updatedFrom.After(*updatedTo) {
		return nil, fmt.Errorf("%w: 'updated_from' não pode ser depois de 'updated_to'",
			errMsg.ErrInvalidFilter)
	}

	// Parsear e validar valores monetários (não permitir negativos)
	minTotalAmount, err := parseFloat(d.MinTotalAmount, "min_total_amount", false)
	if err != nil {
		return nil, err
	}

	maxTotalAmount, err := parseFloat(d.MaxTotalAmount, "max_total_amount", false)
	if err != nil {
		return nil, err
	}

	// Validar intervalos numéricos (min <= max)
	if minTotalAmount != nil && maxTotalAmount != nil && *minTotalAmount > *maxTotalAmount {
		return nil, fmt.Errorf("%w: 'min_total_amount' não pode ser maior que 'max_total_amount'",
			errMsg.ErrInvalidFilter)
	}

	// Parsear outros valores monetários
	minItemsAmount, err := parseFloat(d.MinItemsAmount, "min_items_amount", false)
	if err != nil {
		return nil, err
	}

	maxItemsAmount, err := parseFloat(d.MaxItemsAmount, "max_items_amount", false)
	if err != nil {
		return nil, err
	}

	if minItemsAmount != nil && maxItemsAmount != nil && *minItemsAmount > *maxItemsAmount {
		return nil, fmt.Errorf("%w: 'min_items_amount' não pode ser maior que 'max_items_amount'",
			errMsg.ErrInvalidFilter)
	}

	// Descontos podem ser negativos? Depende da regra de negócio
	// Assumindo que não podem ser negativos por padrão
	minItemsDiscount, err := parseFloat(d.MinItemsDiscount, "min_items_discount", false)
	if err != nil {
		return nil, err
	}

	maxItemsDiscount, err := parseFloat(d.MaxItemsDiscount, "max_items_discount", false)
	if err != nil {
		return nil, err
	}

	if minItemsDiscount != nil && maxItemsDiscount != nil && *minItemsDiscount > *maxItemsDiscount {
		return nil, fmt.Errorf("%w: 'min_items_discount' não pode ser maior que 'max_items_discount'",
			errMsg.ErrInvalidFilter)
	}

	minSaleDiscount, err := parseFloat(d.MinSaleDiscount, "min_sale_discount", false)
	if err != nil {
		return nil, err
	}

	maxSaleDiscount, err := parseFloat(d.MaxSaleDiscount, "max_sale_discount", false)
	if err != nil {
		return nil, err
	}

	if minSaleDiscount != nil && maxSaleDiscount != nil && *minSaleDiscount > *maxSaleDiscount {
		return nil, fmt.Errorf("%w: 'min_sale_discount' não pode ser maior que 'max_sale_discount'",
			errMsg.ErrInvalidFilter)
	}

	// Validar payment_type se houver valores permitidos
	if d.PaymentType != "" {
		validPaymentTypes := map[string]bool{
			"credit":    true,
			"debit":     true,
			"cash":      true,
			"pix":       true,
			"bank_slip": true,
			"":          true, // string vazia é válida (sem filtro)
		}
		if !validPaymentTypes[d.PaymentType] {
			return nil, fmt.Errorf("%w: 'payment_type' com valor inválido '%s'",
				errMsg.ErrInvalidFilter, d.PaymentType)
		}
	}

	// Validar status se houver valores permitidos
	if d.Status != "" {
		validStatuses := map[string]bool{
			"pending":   true,
			"completed": true,
			"cancelled": true,
			"refunded":  true,
			"":          true, // string vazia é válida (sem filtro)
		}
		if !validStatuses[d.Status] {
			return nil, fmt.Errorf("%w: 'status' com valor inválido '%s'",
				errMsg.ErrInvalidFilter, d.Status)
		}
	}

	filter := &modelSale.SaleFilter{
		BaseFilter: baseFilter,

		ClientID:    d.ClientID,
		UserID:      d.UserID,
		PaymentType: d.PaymentType,
		Status:      d.Status,
		Notes:       d.Notes,

		MinTotalAmount: minTotalAmount,
		MaxTotalAmount: maxTotalAmount,

		MinTotalItemsAmount: minItemsAmount,
		MaxTotalItemsAmount: maxItemsAmount,

		MinTotalItemsDiscount: minItemsDiscount,
		MaxTotalItemsDiscount: maxItemsDiscount,

		MinTotalSaleDiscount: minSaleDiscount,
		MaxTotalSaleDiscount: maxSaleDiscount,

		SaleDateFrom: saleDateFrom,
		SaleDateTo:   saleDateTo,

		CreatedFrom: createdFrom,
		CreatedTo:   createdTo,
		UpdatedFrom: updatedFrom,
		UpdatedTo:   updatedTo,
	}

	return filter, nil
}
