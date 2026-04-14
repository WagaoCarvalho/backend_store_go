package model

import (
	"strings"
	"time"

	valCpfCnpj "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/cpf_cnpj"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type Supplier struct {
	ID          int64
	Name        string
	CNPJ        *string
	CPF         *string
	Description string
	Version     int
	Status      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (s *Supplier) Validate() error {
	var errs validators.ValidationErrors

	// Validação do Nome
	if validators.IsBlank(s.Name) {
		errs = append(errs, validators.ValidationError{
			Field:   "name",
			Message: validators.MsgRequiredField,
		})
	} else if len(s.Name) > 100 {
		errs = append(errs, validators.ValidationError{
			Field:   "name",
			Message: validators.MsgMax100,
		})
	}

	// Validação mutuamente exclusiva de CPF/CNPJ
	if s.CPF != nil && s.CNPJ != nil {
		errs = append(errs, validators.ValidationError{
			Field:   "cpf_cnpj",
			Message: validators.MsgInvalidAssociation, // ou criar MsgMutuallyExclusive
		})
	}

	// Validação do CPF
	if s.CPF != nil {
		cpf := strings.TrimSpace(*s.CPF)
		if !valCpfCnpj.IsValidCPF(cpf) {
			errs = append(errs, validators.ValidationError{
				Field:   "cpf",
				Message: "CPF inválido", // Sugiro criar constante MsgInvalidCPF
			})
		}
	}

	// Validação do CNPJ
	if s.CNPJ != nil {
		cnpj := strings.TrimSpace(*s.CNPJ)
		if !valCpfCnpj.IsValidCNPJ(cnpj) {
			errs = append(errs, validators.ValidationError{
				Field:   "cnpj",
				Message: "CNPJ inválido", // Sugiro criar constante MsgInvalidCNPJ
			})
		}
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}
