package dto

import (
	"strings"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

type ProductCategoryDTO struct {
	ID          *int64  `json:"id,omitempty"`
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	CreatedAt   *string `json:"created_at,omitempty"`
	UpdatedAt   *string `json:"updated_at,omitempty"`
}

// Validate valida os campos obrigatórios do DTO
func (dto *ProductCategoryDTO) Validate() error {
	name := strings.TrimSpace(dto.Name)
	if validators.IsBlank(name) {
		return &validators.ValidationError{Field: "name", Message: "nome é obrigatório"}
	}
	if len(name) < 2 {
		return &validators.ValidationError{Field: "name", Message: validators.MsgMin2}
	}
	if len(name) > 255 {
		return &validators.ValidationError{Field: "name", Message: "nome máximo 255 caracteres"}
	}

	if dto.Description != nil {
		desc := strings.TrimSpace(*dto.Description)
		if len(desc) > 255 {
			return &validators.ValidationError{Field: "description", Message: "descrição máxima 255 caracteres"}
		}
		// Atualiza o ponteiro com o valor trimmed
		*dto.Description = desc
	}

	return nil
}

// ToProductCategoryModel converte DTO para Model
func ToProductCategoryModel(dto ProductCategoryDTO) *models.ProductCategory {
	var id int64
	if dto.ID != nil {
		id = *dto.ID
	}

	var description string
	if dto.Description != nil {
		description = *dto.Description
	}

	return &models.ProductCategory{
		ID:          id,
		Name:        dto.Name,
		Description: description,
	}
}

// ToProductCategoryDTO converte Model para DTO
func ToProductCategoryDTO(m *models.ProductCategory) ProductCategoryDTO {
	if m == nil {
		return ProductCategoryDTO{}
	}

	dto := ProductCategoryDTO{
		ID:          &m.ID,
		Name:        m.Name,
		Description: &m.Description,
	}

	// Formatar timestamps apenas se não forem zero values
	if !m.CreatedAt.IsZero() {
		createdAt := m.CreatedAt.Format(time.RFC3339)
		dto.CreatedAt = &createdAt
	}

	if !m.UpdatedAt.IsZero() {
		updatedAt := m.UpdatedAt.Format(time.RFC3339)
		dto.UpdatedAt = &updatedAt
	}

	return dto
}

// ToProductCategoryDTOs converte slice de Models para slice de DTOs
func ToProductCategoryDTOs(models []*models.ProductCategory) []ProductCategoryDTO {
	if len(models) == 0 {
		return []ProductCategoryDTO{}
	}

	dtos := make([]ProductCategoryDTO, 0, len(models))
	for _, m := range models {
		if m != nil {
			dtos = append(dtos, ToProductCategoryDTO(m))
		}
	}
	return dtos
}

// ToProductCategoryModelPtr converte DTO pointer para Model (útil para handlers)
func ToProductCategoryModelPtr(dto *ProductCategoryDTO) *models.ProductCategory {
	if dto == nil {
		return nil
	}
	return ToProductCategoryModel(*dto)
}
