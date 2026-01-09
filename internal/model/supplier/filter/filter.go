package model

import (
	"time"

	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type SupplierFilter struct {
	filter.BaseFilter

	Name        string
	CNPJ        string
	CPF         string
	Status      *bool
	CreatedFrom *time.Time
	CreatedTo   *time.Time
	UpdatedFrom *time.Time
	UpdatedTo   *time.Time
}

func (f *SupplierFilter) Validate() error {
	if err := f.BaseFilter.Validate(); err != nil {
		return err
	}

	// CPF e CNPJ não podem ser informados juntos
	if f.CPF != "" && f.CNPJ != "" {
		return &validators.ValidationError{
			Field:   "CPF/CNPJ",
			Message: "informe apenas CPF ou CNPJ, nunca ambos",
		}
	}

	// Validação básica de tamanho CPF
	if f.CPF != "" && len(f.CPF) != 14 {
		return &validators.ValidationError{
			Field:   "CPF",
			Message: "formato de CPF inválido",
		}
	}

	// Validação básica de tamanho CNPJ
	if f.CNPJ != "" && len(f.CNPJ) != 18 {
		return &validators.ValidationError{
			Field:   "CNPJ",
			Message: "formato de CNPJ inválido",
		}
	}

	// Valida intervalo de criação
	if f.CreatedFrom != nil && f.CreatedTo != nil && f.CreatedFrom.After(*f.CreatedTo) {
		return &validators.ValidationError{
			Field:   "CreatedFrom/CreatedTo",
			Message: "intervalo de criação inválido",
		}
	}

	// Valida intervalo de atualização
	if f.UpdatedFrom != nil && f.UpdatedTo != nil && f.UpdatedFrom.After(*f.UpdatedTo) {
		return &validators.ValidationError{
			Field:   "UpdatedFrom/UpdatedTo",
			Message: "intervalo de atualização inválido",
		}
	}

	// Datas não podem estar no futuro
	now := time.Now()

	if f.CreatedFrom != nil && f.CreatedFrom.After(now) {
		return &validators.ValidationError{
			Field:   "CreatedFrom",
			Message: "data de criação não pode estar no futuro",
		}
	}

	if f.CreatedTo != nil && f.CreatedTo.After(now) {
		return &validators.ValidationError{
			Field:   "CreatedTo",
			Message: "data de criação não pode estar no futuro",
		}
	}

	if f.UpdatedFrom != nil && f.UpdatedFrom.After(now) {
		return &validators.ValidationError{
			Field:   "UpdatedFrom",
			Message: "data de atualização não pode estar no futuro",
		}
	}

	if f.UpdatedTo != nil && f.UpdatedTo.After(now) {
		return &validators.ValidationError{
			Field:   "UpdatedTo",
			Message: "data de atualização não pode estar no futuro",
		}
	}

	return nil
}
