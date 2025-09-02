package dto

import (
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

type SupplierDTO struct {
	ID      *int64  `json:"id,omitempty"`
	Name    string  `json:"name"`
	CNPJ    *string `json:"cnpj,omitempty"`
	CPF     *string `json:"cpf,omitempty"`
	Version int     `json:"version"`
	Status  bool    `json:"status"`

	CreatedAt string `json:"created_at,omitempty"`
	UpdatedAt string `json:"updated_at,omitempty"`
}

func ToSupplierModel(dto SupplierDTO) *models.Supplier {
	return &models.Supplier{
		ID:      utils.NilToZero(dto.ID),
		Name:    dto.Name,
		CNPJ:    dto.CNPJ,
		CPF:     dto.CPF,
		Version: dto.Version,
		Status:  dto.Status,
	}
}

func ToSupplierDTO(m *models.Supplier) SupplierDTO {
	if m == nil {
		return SupplierDTO{}
	}

	return SupplierDTO{
		ID:        &m.ID,
		Name:      m.Name,
		CNPJ:      m.CNPJ,
		CPF:       m.CPF,
		Version:   m.Version,
		Status:    m.Status,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
		UpdatedAt: m.UpdatedAt.Format(time.RFC3339),
	}
}
