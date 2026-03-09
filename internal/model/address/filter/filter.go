package model

import (
	"time"

	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type AddressFilter struct {
	filter.BaseFilter
	UserID       *int64
	ClientCpfID  *int64
	SupplierID   *int64
	Street       string
	StreetNumber string
	Complement   string
	City         string
	State        string
	Country      string
	PostalCode   string
	IsActive     *bool
	CreatedFrom  *time.Time
	CreatedTo    *time.Time
	UpdatedFrom  *time.Time
	UpdatedTo    *time.Time
}

func (f *AddressFilter) Validate() error {
	if err := f.BaseFilter.Validate(); err != nil {
		return err
	}

	// IsActive agora é opcional - só valida se foi fornecido
	if f.IsActive != nil {
		// pode adicionar validações se necessário
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
