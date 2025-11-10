package dto

import (
	"time"

	modelClient "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client"
	modelFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/filter"
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

// ToModel converte o DTO para o modelo de filtro.
// Garante limite mínimo e parsing seguro das datas.
func (d *ClientFilterDTO) ToModel() (*modelClient.ClientFilter, error) {
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

	// Define limites padrão se não informados ou inválidos
	limit := d.Limit
	if limit <= 0 {
		limit = 100
	}
	offset := d.Offset
	if offset < 0 {
		offset = 0
	}

	filter := &modelClient.ClientFilter{
		BaseFilter: modelFilter.BaseFilter{
			Limit:  limit,
			Offset: offset,
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
