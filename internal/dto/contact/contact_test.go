package dto

import (
	"testing"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	"github.com/stretchr/testify/assert"
)

func TestToContactModel(t *testing.T) {
	now := time.Now()
	id := int64(1)
	userID := int64(10)
	clientID := int64(20)
	supplierID := int64(30)

	dto := ContactDTO{
		ID:              &id,
		UserID:          &userID,
		ClientID:        &clientID,
		SupplierID:      &supplierID,
		ContactName:     "John Doe",
		ContactPosition: "Manager",
		Email:           "john@example.com",
		Phone:           "123456789",
		Cell:            "987654321",
		ContactType:     "primary",
		CreatedAt:       &now,
		UpdatedAt:       &now,
	}

	model := ToContactModel(dto)

	assert.Equal(t, int64(1), model.ID)
	assert.Equal(t, &userID, model.UserID)
	assert.Equal(t, &clientID, model.ClientID)
	assert.Equal(t, &supplierID, model.SupplierID)
	assert.Equal(t, "John Doe", model.ContactName)
	assert.Equal(t, "Manager", model.ContactPosition)
	assert.Equal(t, "john@example.com", model.Email)
	assert.Equal(t, "123456789", model.Phone)
	assert.Equal(t, "987654321", model.Cell)
	assert.Equal(t, "primary", model.ContactType)
	assert.Equal(t, now, model.CreatedAt)
	assert.Equal(t, now, model.UpdatedAt)
}

func TestToContactModel_NilPointers(t *testing.T) {
	dto := ContactDTO{
		ContactName: "Jane Doe",
	}

	model := ToContactModel(dto)

	assert.Equal(t, "Jane Doe", model.ContactName)
	assert.Zero(t, model.ID)
	assert.Nil(t, model.UserID)
	assert.Nil(t, model.ClientID)
	assert.Nil(t, model.SupplierID)
	assert.Equal(t, "", model.ContactPosition)
	assert.Equal(t, "", model.Email)
	assert.Equal(t, "", model.Phone)
	assert.Equal(t, "", model.Cell)
	assert.Equal(t, "", model.ContactType)
}

func TestToContactDTO(t *testing.T) {
	now := time.Now()
	userID := int64(10)
	clientID := int64(20)
	supplierID := int64(30)

	model := &models.Contact{
		ID:              1,
		UserID:          &userID,
		ClientID:        &clientID,
		SupplierID:      &supplierID,
		ContactName:     "John Doe",
		ContactPosition: "Manager",
		Email:           "john@example.com",
		Phone:           "123456789",
		Cell:            "987654321",
		ContactType:     "primary",
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	dto := ToContactDTO(model)

	assert.Equal(t, int64(1), *dto.ID)
	assert.Equal(t, &userID, dto.UserID)
	assert.Equal(t, &clientID, dto.ClientID)
	assert.Equal(t, &supplierID, dto.SupplierID)
	assert.Equal(t, "John Doe", dto.ContactName)
	assert.Equal(t, "Manager", dto.ContactPosition)
	assert.Equal(t, "john@example.com", dto.Email)
	assert.Equal(t, "123456789", dto.Phone)
	assert.Equal(t, "987654321", dto.Cell)
	assert.Equal(t, "primary", dto.ContactType)
	assert.Equal(t, &now, dto.CreatedAt)
	assert.Equal(t, &now, dto.UpdatedAt)
}
