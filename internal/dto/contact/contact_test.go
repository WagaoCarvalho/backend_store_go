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

	dto := ContactDTO{
		ID:                 &id,
		ContactName:        "John Doe",
		ContactDescription: "Manager of sales",
		Email:              "john@example.com",
		Phone:              "123456789",
		Cell:               "987654321",
		ContactType:        "primary",
		CreatedAt:          &now,
		UpdatedAt:          &now,
	}

	model := ToContactModel(dto)

	assert.Equal(t, int64(1), model.ID)
	assert.Equal(t, "John Doe", model.ContactName)
	assert.Equal(t, "Manager of sales", model.ContactDescription)
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
	assert.Equal(t, "", model.ContactDescription)
	assert.Equal(t, "", model.Email)
	assert.Equal(t, "", model.Phone)
	assert.Equal(t, "", model.Cell)
	assert.Equal(t, "", model.ContactType)
}

func TestToContactDTO(t *testing.T) {
	now := time.Now()

	model := &models.Contact{
		ID:                 1,
		ContactName:        "John Doe",
		ContactDescription: "Manager of sales",
		Email:              "john@example.com",
		Phone:              "123456789",
		Cell:               "987654321",
		ContactType:        "primary",
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	dto := ToContactDTO(model)

	assert.Equal(t, int64(1), *dto.ID)
	assert.Equal(t, "John Doe", dto.ContactName)
	assert.Equal(t, "Manager of sales", dto.ContactDescription)
	assert.Equal(t, "john@example.com", dto.Email)
	assert.Equal(t, "123456789", dto.Phone)
	assert.Equal(t, "987654321", dto.Cell)
	assert.Equal(t, "primary", dto.ContactType)
	assert.Equal(t, &now, dto.CreatedAt)
	assert.Equal(t, &now, dto.UpdatedAt)
}

func TestToContactDTOs_EmptyList(t *testing.T) {
	result := ToContactDTOs([]*models.Contact{})
	assert.Empty(t, result)
}

func TestToContactDTOs_SingleContact(t *testing.T) {
	now := time.Now()
	contact := &models.Contact{
		ID:                 1,
		ContactName:        "Alice",
		ContactDescription: "Finance manager",
		Email:              "alice@example.com",
		Phone:              "123456789",
		Cell:               "987654321",
		ContactType:        "financeiro",
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	result := ToContactDTOs([]*models.Contact{contact})

	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), *result[0].ID)
	assert.Equal(t, "Alice", result[0].ContactName)
	assert.Equal(t, "Finance manager", result[0].ContactDescription)
	assert.Equal(t, "alice@example.com", result[0].Email)
	assert.Equal(t, "123456789", result[0].Phone)
	assert.Equal(t, "987654321", result[0].Cell)
	assert.Equal(t, "financeiro", result[0].ContactType)
	assert.Equal(t, &now, result[0].CreatedAt)
	assert.Equal(t, &now, result[0].UpdatedAt)
}

func TestToContactDTOs_MultipleContacts(t *testing.T) {
	now := time.Now()
	contact1 := &models.Contact{
		ID:          1,
		ContactName: "Bob",
		Email:       "bob@example.com",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	contact2 := &models.Contact{
		ID:          2,
		ContactName: "Charlie",
		Email:       "charlie@example.com",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	result := ToContactDTOs([]*models.Contact{contact1, contact2})

	assert.Len(t, result, 2)
	assert.Equal(t, int64(1), *result[0].ID)
	assert.Equal(t, "Bob", result[0].ContactName)
	assert.Equal(t, int64(2), *result[1].ID)
	assert.Equal(t, "Charlie", result[1].ContactName)
}

func TestToContactDTOs_WithNilContact(t *testing.T) {
	now := time.Now()
	contact := &models.Contact{
		ID:          1,
		ContactName: "Valid",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	result := ToContactDTOs([]*models.Contact{contact, nil})

	assert.Len(t, result, 1)
	assert.Equal(t, "Valid", result[0].ContactName)
}
