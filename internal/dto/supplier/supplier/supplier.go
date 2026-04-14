package dto

import (
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

type SupplierDTO struct {
	ID          *int64  `json:"id,omitempty"`
	Name        string  `json:"name"`
	CNPJ        *string `json:"cnpj,omitempty"`
	CPF         *string `json:"cpf,omitempty"`
	Description string  `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
	Version     int     `json:"version"` // ← Tag JSON correta
}

func ToSupplierModel(dto SupplierDTO) *models.Supplier {
	model := &models.Supplier{
		ID:          utils.NilToZero(dto.ID),
		Name:        dto.Name,
		CNPJ:        dto.CNPJ,
		CPF:         dto.CPF,
		Description: dto.Description,
		Status:      true,
		Version:     dto.Version, // ← CORRIGIDO: Copiar o version do DTO
	}

	if dto.IsActive != nil {
		model.Status = *dto.IsActive
	}

	return model
}

func ToSupplierDTO(model *models.Supplier) SupplierDTO {
	if model == nil {
		return SupplierDTO{}
	}

	return SupplierDTO{
		ID:          &model.ID,
		Name:        model.Name,
		CNPJ:        model.CNPJ,
		CPF:         model.CPF,
		Description: model.Description,
		IsActive:    &model.Status,
		Version:     model.Version, // ← CORRIGIDO: Incluir version no DTO de saída
	}
}

func ToSupplierDTOs(models []*models.Supplier) []SupplierDTO {
	if len(models) == 0 {
		return []SupplierDTO{}
	}

	dtos := make([]SupplierDTO, 0, len(models))
	for _, m := range models {
		if m != nil {
			dtos = append(dtos, ToSupplierDTO(m))
		}
	}
	return dtos
}
