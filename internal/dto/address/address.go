package dto

import models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"

type AddressDTO struct {
	ID         *int64 `json:"id,omitempty"`
	UserID     *int64 `json:"user_id,omitempty"`
	ClientID   *int64 `json:"client_id,omitempty"`
	SupplierID *int64 `json:"supplier_id,omitempty"`
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	Country    string `json:"country"`
	PostalCode string `json:"postal_code"`
}

func ToAddressModel(dto AddressDTO) *models.Address {
	return &models.Address{
		ID:         getOrDefault(dto.ID),
		UserID:     dto.UserID,
		ClientID:   dto.ClientID,
		SupplierID: dto.SupplierID,
		Street:     dto.Street,
		City:       dto.City,
		State:      dto.State,
		Country:    dto.Country,
		PostalCode: dto.PostalCode,
	}
}

func ToAddressDTO(model *models.Address) AddressDTO {
	return AddressDTO{
		ID:         &model.ID,
		UserID:     model.UserID,
		ClientID:   model.ClientID,
		SupplierID: model.SupplierID,
		Street:     model.Street,
		City:       model.City,
		State:      model.State,
		Country:    model.Country,
		PostalCode: model.PostalCode,
	}
}

func getOrDefault(id *int64) int64 {
	if id == nil {
		return 0
	}
	return *id
}
