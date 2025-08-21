package models

import (
	"strings"
	"time"

	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators"
)

type Supplier struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CNPJ      *string   `json:"cnpj,omitempty"`
	CPF       *string   `json:"cpf,omitempty"`
	Version   int       `json:"version"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Supplier) Validate() error {
	if validators.IsBlank(s.Name) {
		return &validators.ValidationError{Field: "Name", Message: "campo obrigatório"}
	}
	if len(s.Name) > 100 {
		return &validators.ValidationError{Field: "Name", Message: "máximo de 100 caracteres"}
	}

	// Valida CPF e CNPJ mutuamente exclusivos (se quiser aplicar isso):
	if s.CPF != nil && s.CNPJ != nil {
		return &validators.ValidationError{Field: "CPF/CNPJ", Message: "não é permitido preencher ambos"}
	}

	if s.CPF != nil {
		cpf := strings.TrimSpace(*s.CPF)
		if !validators.IsValidCPF(cpf) {
			return &validators.ValidationError{Field: "CPF", Message: "CPF inválido"}
		}
	}

	if s.CNPJ != nil {
		cnpj := strings.TrimSpace(*s.CNPJ)
		if !validators.IsValidCNPJ(cnpj) {
			return &validators.ValidationError{Field: "CNPJ", Message: "CNPJ inválido"}
		}
	}

	return nil
}
