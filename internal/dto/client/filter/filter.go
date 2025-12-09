package dto

import (
	"time"

	filterClient "github.com/WagaoCarvalho/backend_store_go/internal/model/client/filter"
	commonFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
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
}

func (d *ClientFilterDTO) ToModel() (*filterClient.ClientFilter, error) {
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

	filter := &filterClient.ClientFilter{
		BaseFilter: commonFilter.BaseFilter{
			Limit:  d.Limit,
			Offset: d.Offset,
		},
		Name:        d.Name,
		Email:       d.Email,
		CPF:         d.CPF,
		CNPJ:        d.CNPJ,
		Status:      d.Status,
		Version:     d.Version,
		CreatedFrom: parseDate(d.CreatedFrom),
		CreatedTo:   parseDate(d.CreatedTo),
		UpdatedFrom: parseDate(d.UpdatedFrom),
		UpdatedTo:   parseDate(d.UpdatedTo),
	}

	return filter, nil
}
