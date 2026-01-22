package dto

import (
	"errors"
	"strings"
	"time"

	filterClient "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/filter"
	commonFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
)

type ClientFilterDTO struct {
	Name        string     `schema:"name"`
	Email       string     `schema:"email"`
	CPF         string     `schema:"cpf"`
	CNPJ        string     `schema:"cnpj"`
	Status      *bool      `schema:"status"`
	Version     *int       `schema:"version"`
	CreatedFrom *time.Time `schema:"created_from"`
	CreatedTo   *time.Time `schema:"created_to"`
	UpdatedFrom *time.Time `schema:"updated_from"`
	UpdatedTo   *time.Time `schema:"updated_to"`
	Limit       int        `schema:"limit"`
	Offset      int        `schema:"offset"`
	SortBy      string     `schema:"sort_by"`
	SortOrder   string     `schema:"sort_order"`
}

func (d *ClientFilterDTO) ToModel() (*filterClient.ClientCpfFilter, error) {
	filter := &filterClient.ClientCpfFilter{
		BaseFilter: commonFilter.BaseFilter{
			Limit:     d.Limit,
			Offset:    d.Offset,
			SortBy:    d.SortBy,
			SortOrder: d.SortOrder,
		},
		Name:        d.Name,
		Email:       d.Email,
		CPF:         d.CPF,
		CNPJ:        d.CNPJ,
		Status:      d.Status,
		Version:     d.Version,
		CreatedFrom: d.CreatedFrom,
		CreatedTo:   d.CreatedTo,
		UpdatedFrom: d.UpdatedFrom,
		UpdatedTo:   d.UpdatedTo,
	}

	return filter, nil
}

func (d *ClientFilterDTO) Validate() error {
	var validationErrors []string

	name := strings.TrimSpace(d.Name)
	email := strings.TrimSpace(d.Email)
	cpf := strings.TrimSpace(d.CPF)
	cnpj := strings.TrimSpace(d.CNPJ)

	// Pelo menos um filtro de conteúdo
	hasContentFilter :=
		name != "" ||
			email != "" ||
			cpf != "" ||
			cnpj != "" ||
			d.Status != nil ||
			d.Version != nil ||
			d.CreatedFrom != nil ||
			d.CreatedTo != nil ||
			d.UpdatedFrom != nil ||
			d.UpdatedTo != nil

	if !hasContentFilter {
		validationErrors = append(validationErrors, "pelo menos um filtro de busca deve ser fornecido")
	}

	// ===== LIMITES MÍNIMOS (anti-abuso) =====
	if name != "" && len(name) < 3 {
		validationErrors = append(validationErrors, "'name' deve conter no mínimo 3 caracteres")
	}

	if email != "" && len(email) < 5 {
		validationErrors = append(validationErrors, "'email' deve conter no mínimo 5 caracteres")
	}

	// ===== CPF / CNPJ =====
	if cpf != "" && len(cpf) != 11 {
		validationErrors = append(validationErrors, "'cpf' inválido")
	}

	if cnpj != "" && len(cnpj) != 14 {
		validationErrors = append(validationErrors, "'cnpj' inválido")
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
		validationErrors = append(validationErrors, "'offset' excede o limite permitido")
	}

	// ===== ORDENAÇÃO =====
	if d.SortBy != "" && !isValidSortField(d.SortBy) {
		validationErrors = append(validationErrors, "'sort_by' inválido")
	}

	if d.SortOrder != "" {
		order := strings.ToLower(d.SortOrder)
		if order != "asc" && order != "desc" {
			validationErrors = append(validationErrors, "'sort_order' inválido")
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

func isValidSortField(field string) bool {
	allowedFields := map[string]bool{
		"id":         true,
		"name":       true,
		"email":      true,
		"status":     true,
		"version":    true,
		"created_at": true,
		"updated_at": true,
	}
	return allowedFields[strings.ToLower(field)]
}
