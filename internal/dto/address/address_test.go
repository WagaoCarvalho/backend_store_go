package dto

import (
	"testing"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	"github.com/stretchr/testify/assert"
)

func TestToAddressModel(t *testing.T) {
	id := int64(1)
	userID := int64(10)
	clientID := int64(20)
	supplierID := int64(30)

	dto := AddressDTO{
		ID:         &id,
		UserID:     &userID,
		ClientID:   &clientID,
		SupplierID: &supplierID,
		Street:     "Rua A",
		City:       "Cidade B",
		State:      "Estado C",
		Country:    "Brasil",
		PostalCode: "12345678",
	}

	model := ToAddressModel(dto)

	assert.Equal(t, int64(1), model.ID)
	assert.Equal(t, &userID, model.UserID)
	assert.Equal(t, &clientID, model.ClientID)
	assert.Equal(t, &supplierID, model.SupplierID)
	assert.Equal(t, "Rua A", model.Street)
	assert.Equal(t, "Cidade B", model.City)
	assert.Equal(t, "Estado C", model.State)
	assert.Equal(t, "Brasil", model.Country)
	assert.Equal(t, "12345678", model.PostalCode)
}

func TestToAddressModel_NilID(t *testing.T) {
	dto := AddressDTO{
		Street:     "Rua A",
		City:       "Cidade B",
		State:      "Estado C",
		Country:    "Brasil",
		PostalCode: "12345678",
	}

	model := ToAddressModel(dto)

	assert.Equal(t, int64(0), model.ID) // getOrDefault retorna 0
	assert.Nil(t, model.UserID)
	assert.Nil(t, model.ClientID)
	assert.Nil(t, model.SupplierID)
}

func TestToAddressDTO(t *testing.T) {
	userID := int64(10)
	clientID := int64(20)
	supplierID := int64(30)

	model := &models.Address{
		ID:         1,
		UserID:     &userID,
		ClientID:   &clientID,
		SupplierID: &supplierID,
		Street:     "Rua A",
		City:       "Cidade B",
		State:      "Estado C",
		Country:    "Brasil",
		PostalCode: "12345678",
	}

	dto := ToAddressDTO(model)

	assert.Equal(t, int64(1), *dto.ID)
	assert.Equal(t, &userID, dto.UserID)
	assert.Equal(t, &clientID, dto.ClientID)
	assert.Equal(t, &supplierID, dto.SupplierID)
	assert.Equal(t, "Rua A", dto.Street)
	assert.Equal(t, "Cidade B", dto.City)
	assert.Equal(t, "Estado C", dto.State)
	assert.Equal(t, "Brasil", dto.Country)
	assert.Equal(t, "12345678", dto.PostalCode)
}
