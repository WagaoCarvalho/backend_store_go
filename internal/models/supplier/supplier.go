package models

import (
	"strings"
	"time"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	models_supplier_categories "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_categories"
	utils_errors "github.com/WagaoCarvalho/backend_store_go/internal/utils"
	utils_validators "github.com/WagaoCarvalho/backend_store_go/internal/utils/validators"
)

type Supplier struct {
	ID          int64                                         `json:"id"`
	Name        string                                        `json:"name"`
	CNPJ        *string                                       `json:"cnpj,omitempty"`
	CPF         *string                                       `json:"cpf,omitempty"`
	ContactInfo string                                        `json:"contact_info"`
	Version     int                                           `json:"version"`
	CreatedAt   time.Time                                     `json:"created_at"`
	UpdatedAt   time.Time                                     `json:"updated_at"`
	Categories  []models_supplier_categories.SupplierCategory `json:"categories,omitempty"`
	Address     *models_address.Address                       `json:"address,omitempty"`
	Contact     *models_contact.Contact                       `json:"contact,omitempty"`
}

func (s *Supplier) Validate() error {
	if utils_validators.IsBlank(s.Name) {
		return &utils_errors.ValidationError{Field: "Name", Message: "campo obrigatório"}
	}
	if len(s.Name) > 100 {
		return &utils_errors.ValidationError{Field: "Name", Message: "máximo de 100 caracteres"}
	}

	if utils_validators.IsBlank(s.ContactInfo) {
		return &utils_errors.ValidationError{Field: "ContactInfo", Message: "campo obrigatório"}
	}
	if len(s.ContactInfo) > 100 {
		return &utils_errors.ValidationError{Field: "ContactInfo", Message: "máximo de 100 caracteres"}
	}

	// Valida CPF e CNPJ mutuamente exclusivos (se quiser aplicar isso):
	if s.CPF != nil && s.CNPJ != nil {
		return &utils_errors.ValidationError{Field: "CPF/CNPJ", Message: "não é permitido preencher ambos"}
	}

	if s.CPF != nil {
		cpf := strings.TrimSpace(*s.CPF)
		if !utils_validators.IsValidCPF(cpf) {
			return &utils_errors.ValidationError{Field: "CPF", Message: "CPF inválido"}
		}
	}

	if s.CNPJ != nil {
		cnpj := strings.TrimSpace(*s.CNPJ)
		if !utils_validators.IsValidCNPJ(cnpj) {
			return &utils_errors.ValidationError{Field: "CNPJ", Message: "CNPJ inválido"}
		}
	}

	// Validação opcional para Address e Contact
	if s.Address != nil {
		if err := s.Address.Validate(); err != nil {
			return err
		}
	}

	if s.Contact != nil {
		if err := s.Contact.Validate(); err != nil {
			return err
		}
	}

	return nil
}
