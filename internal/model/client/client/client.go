package model

import (
	"strings"
	"time"

	valCpfCnpj "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/cpf_cnpj"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type Client struct {
	ID         int64
	Name       string
	Email      *string
	CPF        *string
	CNPJ       *string
	ClientType string
	Version    int
	Status     bool
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (c *Client) Validate() error {
	if validators.IsBlank(c.Name) {
		return &validators.ValidationError{Field: "Name", Message: "campo obrigatório"}
	}
	if len(c.Name) > 255 {
		return &validators.ValidationError{Field: "Name", Message: "máximo de 255 caracteres"}
	}

	if validators.IsBlank(c.ClientType) {
		return &validators.ValidationError{Field: "ClientType", Message: "campo obrigatório"}
	}
	if c.ClientType != "PF" && c.ClientType != "PJ" {
		return &validators.ValidationError{Field: "ClientType", Message: "deve ser PF ou PJ"}
	}

	// Regras de CPF e CNPJ
	if c.CPF != nil && c.CNPJ != nil {
		return &validators.ValidationError{Field: "CPF/CNPJ", Message: "não é permitido preencher ambos"}
	}

	if c.ClientType == "PF" {
		if c.CPF == nil {
			return &validators.ValidationError{Field: "CPF", Message: "obrigatório para pessoa física"}
		}
		cpf := strings.TrimSpace(*c.CPF)
		if !valCpfCnpj.IsValidCPF(cpf) {
			return &validators.ValidationError{Field: "CPF", Message: "CPF inválido"}
		}
	}

	if c.ClientType == "PJ" {
		if c.CNPJ == nil {
			return &validators.ValidationError{Field: "CNPJ", Message: "obrigatório para pessoa jurídica"}
		}
		cnpj := strings.TrimSpace(*c.CNPJ)
		if !valCpfCnpj.IsValidCNPJ(cnpj) {
			return &validators.ValidationError{Field: "CNPJ", Message: "CNPJ inválido"}
		}
	}

	return nil
}
