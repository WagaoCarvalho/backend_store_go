package model

import (
	"strings"
	"time"

	valCpfCnpj "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/cpf_cnpj"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type Client struct {
	ID        int64
	Name      string
	Email     *string
	CPF       *string
	CNPJ      *string
	Version   int
	Status    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (c *Client) Validate() error {
	if validators.IsBlank(c.Name) {
		return &validators.ValidationError{Field: "Name", Message: "campo obrigatório"}
	}
	if len(c.Name) > 255 {
		return &validators.ValidationError{Field: "Name", Message: "máximo de 255 caracteres"}
	}

	// Não pode preencher ambos
	if c.CPF != nil && c.CNPJ != nil {
		return &validators.ValidationError{Field: "CPF/CNPJ", Message: "não é permitido preencher ambos"}
	}

	// Deve preencher um dos dois
	if c.CPF == nil && c.CNPJ == nil {
		return &validators.ValidationError{Field: "CPF/CNPJ", Message: "deve informar CPF ou CNPJ"}
	}

	// Valida CPF
	if c.CPF != nil {
		cpf := strings.TrimSpace(*c.CPF)
		if !valCpfCnpj.IsValidCPF(cpf) {
			return &validators.ValidationError{Field: "CPF", Message: "CPF inválido"}
		}
	}

	// Valida CNPJ
	if c.CNPJ != nil {
		cnpj := strings.TrimSpace(*c.CNPJ)
		if !valCpfCnpj.IsValidCNPJ(cnpj) {
			return &validators.ValidationError{Field: "CNPJ", Message: "CNPJ inválido"}
		}
	}

	return nil
}
