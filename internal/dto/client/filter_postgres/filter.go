package dto

import (
	"time"

	filterClient "github.com/WagaoCarvalho/backend_store_go/internal/model/client/filter"
	commonFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	parser "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/parser"
)

type ClientFilterDTO struct {
	Name        string  `schema:"name"`
	Email       string  `schema:"email"`
	CPF         string  `schema:"cpf"`
	CNPJ        string  `schema:"cnpj"`
	Status      *bool   `schema:"status"`
	Version     *int    `schema:"version"`
	CreatedFrom *string `schema:"created_from"`
	CreatedTo   *string `schema:"created_to"`
	UpdatedFrom *string `schema:"updated_from"`
	UpdatedTo   *string `schema:"updated_to"`
	Limit       int     `schema:"limit"`
	Offset      int     `schema:"offset"`
	SortBy      string  `schema:"sort_by"`
	SortOrder   string  `schema:"sort_order"`
	SearchTerm  string  `schema:"search_term"`
}

func (d *ClientFilterDTO) Validate() error {
	// Adicione validações específicas do DTO aqui
	return nil
}

func (d *ClientFilterDTO) ToModel() (*filterClient.ClientFilter, error) {
	dateParser := parser.NewDateParser()

	parseDate := func(s *string) (*time.Time, error) {
		if s == nil || *s == "" {
			return nil, nil
		}
		return dateParser.Parse(*s)
	}

	createdFrom, err := parseDate(d.CreatedFrom)
	if err != nil {
		return nil, err
	}

	createdTo, err := parseDate(d.CreatedTo)
	if err != nil {
		return nil, err
	}

	updatedFrom, err := parseDate(d.UpdatedFrom)
	if err != nil {
		return nil, err
	}

	updatedTo, err := parseDate(d.UpdatedTo)
	if err != nil {
		return nil, err
	}

	// CORREÇÃO AQUI: Use NewBaseFilter em vez de conversão direta
	baseFilter := commonFilter.NewBaseFilter(
		d.Limit,
		d.Offset,
		d.SortBy,
		d.SortOrder,
		d.SearchTerm,
	)

	filter := &filterClient.ClientFilter{
		BaseFilter:  *baseFilter, // Desreferencie o ponteiro
		Name:        d.Name,
		Email:       d.Email,
		CPF:         d.CPF,
		CNPJ:        d.CNPJ,
		Status:      d.Status,
		Version:     d.Version,
		CreatedFrom: createdFrom,
		CreatedTo:   createdTo,
		UpdatedFrom: updatedFrom,
		UpdatedTo:   updatedTo,
	}

	return filter, nil
}
