package dto

import (
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

type SaleDTO struct {
	ID            *int64  `json:"id,omitempty"`
	ClientID      *int64  `json:"client_id,omitempty"`
	UserID        *int64  `json:"user_id,omitempty"`
	SaleDate      *string `json:"sale_date,omitempty"`
	TotalAmount   float64 `json:"total_amount"`
	TotalDiscount float64 `json:"total_discount,omitempty"`
	PaymentType   string  `json:"payment_type"`
	Status        string  `json:"status,omitempty"`
	Notes         string  `json:"notes,omitempty"`
	Version       int     `json:"version,omitempty"`
	CreatedAt     *string `json:"created_at,omitempty"`
	UpdatedAt     *string `json:"updated_at,omitempty"`
}

func ToSaleModel(dto SaleDTO) *models.Sale {
	model := &models.Sale{
		ID:            utils.NilToZero(dto.ID),
		ClientID:      dto.ClientID,
		UserID:        dto.UserID,
		TotalAmount:   dto.TotalAmount,
		TotalDiscount: dto.TotalDiscount,
		PaymentType:   dto.PaymentType,
		Status:        utils.DefaultString(dto.Status, "active"),
		Notes:         dto.Notes,
		Version:       dto.Version,
	}

	if dto.SaleDate != nil {
		if t, err := time.Parse(time.RFC3339, *dto.SaleDate); err == nil {
			model.SaleDate = t
		}
	}
	if dto.CreatedAt != nil {
		if t, err := time.Parse(time.RFC3339, *dto.CreatedAt); err == nil {
			model.CreatedAt = t
		}
	}
	if dto.UpdatedAt != nil {
		if t, err := time.Parse(time.RFC3339, *dto.UpdatedAt); err == nil {
			model.UpdatedAt = t
		}
	}

	return model
}

func ToSaleDTO(model *models.Sale) SaleDTO {
	saleDate := model.SaleDate.Format(time.RFC3339)
	createdAt := model.CreatedAt.Format(time.RFC3339)
	updatedAt := model.UpdatedAt.Format(time.RFC3339)

	return SaleDTO{
		ID:            &model.ID,
		ClientID:      model.ClientID,
		UserID:        model.UserID,
		SaleDate:      &saleDate,
		TotalAmount:   model.TotalAmount,
		TotalDiscount: model.TotalDiscount,
		PaymentType:   model.PaymentType,
		Status:        model.Status,
		Notes:         model.Notes,
		Version:       model.Version,
		CreatedAt:     &createdAt,
		UpdatedAt:     &updatedAt,
	}
}

func ToSaleDTOList(modelsList []*models.Sale) []SaleDTO {
	result := make([]SaleDTO, len(modelsList))
	for i, m := range modelsList {
		result[i] = ToSaleDTO(m)
	}
	return result
}

func SaleDTOListToModelList(dtos []*SaleDTO) []*models.Sale {
	result := make([]*models.Sale, len(dtos))
	for i, dto := range dtos {
		result[i] = ToSaleModel(*dto)
	}
	return result
}
