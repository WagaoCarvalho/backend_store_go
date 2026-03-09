package dto

import (
	"errors"
	"strings"
	"time"

	filterAddress "github.com/WagaoCarvalho/backend_store_go/internal/model/address/filter"
	commonFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
)

type AddressFilterDTO struct {
	UserID       *int64     `schema:"user_id"`
	ClientCpfID  *int64     `schema:"client_cpf_id"`
	SupplierID   *int64     `schema:"supplier_id"`
	Street       string     `schema:"street"`
	StreetNumber string     `schema:"street_number"`
	Complement   string     `schema:"complement"`
	City         string     `schema:"city"`
	State        string     `schema:"state"`
	Country      string     `schema:"country"`
	PostalCode   string     `schema:"postal_code"`
	IsActive     *bool      `schema:"is_active"`
	CreatedFrom  *time.Time `schema:"created_from"`
	CreatedTo    *time.Time `schema:"created_to"`
	UpdatedFrom  *time.Time `schema:"updated_from"`
	UpdatedTo    *time.Time `schema:"updated_to"`
	Limit        int        `schema:"limit"`
	Offset       int        `schema:"offset"`
	SortBy       string     `schema:"sort_by"`
	SortOrder    string     `schema:"sort_order"`
}

func (d *AddressFilterDTO) ToModel() (*filterAddress.AddressFilter, error) {
	filter := &filterAddress.AddressFilter{
		BaseFilter: commonFilter.BaseFilter{
			Limit:     d.Limit,
			Offset:    d.Offset,
			SortBy:    d.SortBy,
			SortOrder: d.SortOrder,
		},
		UserID:       d.UserID,
		ClientCpfID:  d.ClientCpfID,
		SupplierID:   d.SupplierID,
		Street:       d.Street,
		StreetNumber: d.StreetNumber,
		Complement:   d.Complement,
		City:         d.City,
		State:        d.State,
		Country:      d.Country,
		PostalCode:   d.PostalCode,
		CreatedFrom:  d.CreatedFrom,
		CreatedTo:    d.CreatedTo,
		UpdatedFrom:  d.UpdatedFrom,
		UpdatedTo:    d.UpdatedTo,
	}

	// Tratar IsActive - se nil, usar false como padrão
	if d.IsActive != nil {
		filter.IsActive = d.IsActive
	} else {
		filter.IsActive = nil // valor padrão
	}

	return filter, nil
}

func (d *AddressFilterDTO) Validate() error {
	var validationErrors []string

	// Trim em todos os campos de string
	d.Street = strings.TrimSpace(d.Street)
	d.StreetNumber = strings.TrimSpace(d.StreetNumber)
	d.Complement = strings.TrimSpace(d.Complement)
	d.City = strings.TrimSpace(d.City)
	d.State = strings.TrimSpace(d.State)
	d.Country = strings.TrimSpace(d.Country)
	d.PostalCode = strings.TrimSpace(d.PostalCode)

	// Pelo menos um filtro de conteúdo
	hasContentFilter :=
		d.UserID != nil ||
			d.ClientCpfID != nil ||
			d.SupplierID != nil ||
			d.Street != "" ||
			d.StreetNumber != "" ||
			d.Complement != "" ||
			d.City != "" ||
			d.State != "" ||
			d.Country != "" ||
			d.PostalCode != "" ||
			d.IsActive != nil ||
			d.CreatedFrom != nil ||
			d.CreatedTo != nil ||
			d.UpdatedFrom != nil ||
			d.UpdatedTo != nil

	if !hasContentFilter {
		validationErrors = append(validationErrors, "pelo menos um filtro de busca deve ser fornecido")
	}

	// ===== VALIDAÇÕES ESPECÍFICAS =====
	if d.Street != "" && len(d.Street) < 2 {
		validationErrors = append(validationErrors, "'street' deve conter no mínimo 2 caracteres")
	}

	if d.StreetNumber != "" && len(d.StreetNumber) > 20 {
		validationErrors = append(validationErrors, "'street_number' não pode ter mais que 20 caracteres")
	}

	if d.Complement != "" && len(d.Complement) > 100 {
		validationErrors = append(validationErrors, "'complement' não pode ter mais que 100 caracteres")
	}

	if d.City != "" && len(d.City) < 2 {
		validationErrors = append(validationErrors, "'city' deve conter no mínimo 2 caracteres")
	}

	if d.State != "" {
		if len(d.State) != 2 {
			validationErrors = append(validationErrors, "'state' deve conter exatamente 2 caracteres (UF)")
		}
		// Converter para maiúsculas para padronização
		d.State = strings.ToUpper(d.State)
	}

	if d.Country != "" && len(d.Country) < 2 {
		validationErrors = append(validationErrors, "'country' deve conter no mínimo 2 caracteres")
	}

	if d.PostalCode != "" {
		cleanPostalCode := strings.ReplaceAll(d.PostalCode, "-", "")
		cleanPostalCode = strings.ReplaceAll(cleanPostalCode, ".", "")
		cleanPostalCode = strings.ReplaceAll(cleanPostalCode, " ", "")
		if len(cleanPostalCode) != 8 {
			validationErrors = append(validationErrors, "'postal_code' inválido - deve conter 8 dígitos")
		}
		// Se for válido, armazenar sem formatação
		if len(validationErrors) == 0 {
			d.PostalCode = cleanPostalCode
		}
	}

	// ===== VALIDAÇÃO DE IDs =====
	if d.UserID != nil && *d.UserID <= 0 {
		validationErrors = append(validationErrors, "'user_id' deve ser maior que zero")
	}

	if d.ClientCpfID != nil && *d.ClientCpfID <= 0 {
		validationErrors = append(validationErrors, "'client_cpf_id' deve ser maior que zero")
	}

	if d.SupplierID != nil && *d.SupplierID <= 0 {
		validationErrors = append(validationErrors, "'supplier_id' deve ser maior que zero")
	}

	// ===== PAGINAÇÃO =====
	if d.Limit < 1 {
		validationErrors = append(validationErrors, "'limit' deve ser maior que zero")
	}
	if d.Limit > 100 {
		validationErrors = append(validationErrors, "'limit' não pode ser maior que 100")
	}
	if d.Offset < 0 {
		validationErrors = append(validationErrors, "'offset' não pode ser negativo")
	}
	if d.Offset > 10_000 {
		validationErrors = append(validationErrors, "'offset' excede o limite permitido (máximo 10.000)")
	}

	// ===== ORDENAÇÃO =====
	if d.SortBy != "" && !isValidAddressSortField(d.SortBy) {
		validationErrors = append(validationErrors, "'sort_by' inválido. Campos permitidos: id, user_id, client_cpf_id, supplier_id, street, street_number, city, state, country, postal_code, is_active, created_at, updated_at")
	}

	if d.SortOrder != "" {
		order := strings.ToLower(d.SortOrder)
		if order != "asc" && order != "desc" {
			validationErrors = append(validationErrors, "'sort_order' inválido. Use 'asc' ou 'desc'")
		} else {
			d.SortOrder = order // normalizar para minúsculas
		}
	}

	// ===== DATAS =====
	if d.CreatedFrom != nil && d.CreatedTo != nil && d.CreatedFrom.After(*d.CreatedTo) {
		validationErrors = append(validationErrors, "'created_from' não pode ser maior que 'created_to'")
	}

	if d.UpdatedFrom != nil && d.UpdatedTo != nil && d.UpdatedFrom.After(*d.UpdatedTo) {
		validationErrors = append(validationErrors, "'updated_from' não pode ser maior que 'updated_to'")
	}

	if len(validationErrors) > 0 {
		return errors.New(strings.Join(validationErrors, "; "))
	}
	return nil
}

func isValidAddressSortField(field string) bool {
	allowedFields := map[string]bool{
		"id":            true,
		"user_id":       true,
		"client_cpf_id": true,
		"supplier_id":   true,
		"street":        true,
		"street_number": true,
		"city":          true,
		"state":         true,
		"country":       true,
		"postal_code":   true,
		"is_active":     true,
		"created_at":    true,
		"updated_at":    true,
	}
	return allowedFields[strings.ToLower(field)]
}

// Função auxiliar para criar DTO a partir de query params (opcional)
func NewAddressFilterDTOFromQuery(params map[string][]string) (*AddressFilterDTO, error) {
	dto := &AddressFilterDTO{}

	// Implementar se necessário
	// Pode usar uma biblioteca como "github.com/gorilla/schema" para decodificar

	return dto, nil
}
