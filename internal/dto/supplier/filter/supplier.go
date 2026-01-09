package dto

import (
	"fmt"
	"time"

	modelFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	modelSupplier "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/filter"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
)

type SupplierFilterDTO struct {
	Name        string  `schema:"name"`
	CPF         string  `schema:"cpf"`
	CNPJ        string  `schema:"cnpj"`
	Status      *bool   `schema:"status"`
	CreatedFrom *string `schema:"created_from"`
	CreatedTo   *string `schema:"created_to"`
	UpdatedFrom *string `schema:"updated_from"`
	UpdatedTo   *string `schema:"updated_to"`
	Limit       int     `schema:"limit"`
	Offset      int     `schema:"offset"`
}

func (d *SupplierFilterDTO) ToModel() (*modelSupplier.SupplierFilter, error) {
	// Parse de datas
	parseDate := func(s *string, fieldName string) (*time.Time, error) {
		if s == nil || *s == "" {
			return nil, nil
		}
		t, err := time.Parse("2006-01-02", *s)
		if err != nil {
			return nil, fmt.Errorf(
				"%w: campo '%s' com valor inválido '%s' - formato esperado: YYYY-MM-DD",
				errMsg.ErrInvalidFilter, fieldName, *s,
			)
		}
		return &t, nil
	}

	// Validação de paginação
	if d.Limit < 1 {
		return nil, fmt.Errorf("%w: 'limit' deve ser maior que 0", errMsg.ErrInvalidFilter)
	}
	if d.Limit > 100 {
		return nil, fmt.Errorf("%w: 'limit' máximo é 100", errMsg.ErrInvalidFilter)
	}
	if d.Offset < 0 {
		return nil, fmt.Errorf("%w: 'offset' não pode ser negativo", errMsg.ErrInvalidFilter)
	}

	// CPF e CNPJ são mutuamente exclusivos
	if d.CPF != "" && d.CNPJ != "" {
		return nil, fmt.Errorf(
			"%w: informe apenas 'cpf' ou 'cnpj', nunca ambos",
			errMsg.ErrInvalidFilter,
		)
	}

	baseFilter := modelFilter.BaseFilter{
		Limit:  d.Limit,
		Offset: d.Offset,
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
		return nil, fmt.Errorf(
			"%w: 'created_from' não pode ser depois de 'created_to'",
			errMsg.ErrInvalidFilter,
		)
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
		return nil, fmt.Errorf(
			"%w: 'updated_from' não pode ser depois de 'updated_to'",
			errMsg.ErrInvalidFilter,
		)
	}

	filter := &modelSupplier.SupplierFilter{
		BaseFilter: baseFilter,

		Name:   d.Name,
		CPF:    d.CPF,
		CNPJ:   d.CNPJ,
		Status: d.Status,

		CreatedFrom: createdFrom,
		CreatedTo:   createdTo,
		UpdatedFrom: updatedFrom,
		UpdatedTo:   updatedTo,
	}

	return filter, nil
}
