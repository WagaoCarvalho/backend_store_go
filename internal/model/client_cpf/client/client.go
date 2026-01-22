package model

import (
	"strings"
	"time"

	valCpfCnpj "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/cpf_cnpj"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type ClientCpf struct {
	ID          int64
	Name        string
	Email       string
	CPF         string
	Description string
	Version     int
	Status      bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (c *ClientCpf) Validate() error {
	// Name
	if validators.IsBlank(c.Name) {
		return &validators.ValidationError{Field: "Name", Message: "campo obrigatório"}
	}
	if len(c.Name) > 255 {
		return &validators.ValidationError{Field: "Name", Message: "máximo de 255 caracteres"}
	}

	// Email
	c.Email = strings.TrimSpace(strings.ToLower(c.Email))
	if validators.IsBlank(c.Email) {
		return &validators.ValidationError{Field: "Email", Message: "campo obrigatório"}
	}
	if !validators.IsEmail(c.Email) {
		return &validators.ValidationError{Field: "Email", Message: "email inválido"}
	}

	// CPF
	c.CPF = strings.TrimSpace(c.CPF)
	if validators.IsBlank(c.CPF) {
		return &validators.ValidationError{Field: "CPF", Message: "campo obrigatório"}
	}
	if !valCpfCnpj.IsValidCPF(c.CPF) {
		return &validators.ValidationError{Field: "CPF", Message: "CPF inválido"}
	}

	// Description
	if len(c.Description) > 1000 {
		return &validators.ValidationError{Field: "Description", Message: "máximo de 1000 caracteres"}
	}

	// Version
	if c.Version <= 0 {
		return &validators.ValidationError{Field: "Version", Message: "versão inválida"}
	}

	return nil
}
