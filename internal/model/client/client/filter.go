package model

import (
	"time"

	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/filter"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type ClientFilter struct {
	filter.BaseFilter

	Name        string
	Email       string
	CPF         string
	CNPJ        string
	Status      *bool
	Version     *int
	CreatedFrom *time.Time
	CreatedTo   *time.Time
	UpdatedFrom *time.Time
	UpdatedTo   *time.Time
}

func (f *ClientFilter) Validate() error {
	if err := f.BaseFilter.Validate(); err != nil {
		return err
	}

	if f.CreatedFrom != nil && f.CreatedTo != nil && f.CreatedFrom.After(*f.CreatedTo) {
		return &validators.ValidationError{
			Field:   "CreatedFrom/CreatedTo",
			Message: "intervalo de criação inválido",
		}
	}
	if f.UpdatedFrom != nil && f.UpdatedTo != nil && f.UpdatedFrom.After(*f.UpdatedTo) {
		return &validators.ValidationError{
			Field:   "UpdatedFrom/UpdatedTo",
			Message: "intervalo de atualização inválido",
		}
	}

	return nil
}
