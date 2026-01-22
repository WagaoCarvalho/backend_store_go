package dto

import (
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

type AddressDTO struct {
	ID           *int64 `json:"id,omitempty"`
	UserID       *int64 `json:"user_id,omitempty"`
	ClientCpfID  *int64 `json:"client_cpf_id,omitempty"`
	SupplierID   *int64 `json:"supplier_id,omitempty"`
	Street       string `json:"street"`
	StreetNumber string `json:"street_number,omitempty"`
	Complement   string `json:"complement,omitempty"`
	City         string `json:"city"`
	State        string `json:"state"`
	Country      string `json:"country"`
	PostalCode   string `json:"postal_code"`
	IsActive     *bool  `json:"is_active,omitempty"`
}

func ToAddressModel(dto AddressDTO) *models.Address {
	model := &models.Address{
		ID:           utils.NilToZero(dto.ID),
		UserID:       dto.UserID,
		ClientCpfID:  dto.ClientCpfID,
		SupplierID:   dto.SupplierID,
		Street:       dto.Street,
		StreetNumber: dto.StreetNumber,
		Complement:   dto.Complement,
		City:         dto.City,
		State:        dto.State,
		Country:      dto.Country,
		PostalCode:   dto.PostalCode,
		IsActive:     true,
	}
	if dto.IsActive != nil {
		model.IsActive = *dto.IsActive
	}
	return model
}

func ToAddressDTO(model *models.Address) AddressDTO {
	return AddressDTO{
		ID:           &model.ID,
		UserID:       model.UserID,
		ClientCpfID:  model.ClientCpfID,
		SupplierID:   model.SupplierID,
		Street:       model.Street,
		StreetNumber: model.StreetNumber,
		Complement:   model.Complement,
		City:         model.City,
		State:        model.State,
		Country:      model.Country,
		PostalCode:   model.PostalCode,
		IsActive:     &model.IsActive,
	}
}

func ToAddressDTOs(models []*models.Address) []AddressDTO {
	if len(models) == 0 {
		return []AddressDTO{}
	}

	dtos := make([]AddressDTO, 0, len(models))
	for _, m := range models {
		if m != nil {
			dtos = append(dtos, ToAddressDTO(m))
		}
	}
	return dtos
}
