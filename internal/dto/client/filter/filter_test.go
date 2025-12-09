package dto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestClientFilterDTO_ToModel(t *testing.T) {
	createdFromStr := "2025-01-01"
	createdToStr := "2025-01-31"
	updatedFromStr := "2025-02-01"
	updatedToStr := "2025-02-15"

	status := true
	version := 3

	dto := ClientFilterDTO{
		Name:        "Cliente XPTO",
		Email:       "cliente@exemplo.com",
		CPF:         "12345678901",
		CNPJ:        "12345678000199",
		Status:      &status,
		Version:     &version,
		CreatedFrom: &createdFromStr,
		CreatedTo:   &createdToStr,
		UpdatedFrom: &updatedFromStr,
		UpdatedTo:   &updatedToStr,
		Limit:       50,
		Offset:      10,
	}

	model, err := dto.ToModel()
	assert.NoError(t, err)
	assert.NotNil(t, model)

	assert.Equal(t, "Cliente XPTO", model.Name)
	assert.Equal(t, "cliente@exemplo.com", model.Email)
	assert.Equal(t, "12345678901", model.CPF)
	assert.Equal(t, "12345678000199", model.CNPJ)
	assert.Equal(t, &status, model.Status)
	assert.Equal(t, &version, model.Version)
	assert.Equal(t, 50, model.Limit)
	assert.Equal(t, 10, model.Offset)

	expectedCreatedFrom, _ := time.Parse("2006-01-02", createdFromStr)
	expectedCreatedTo, _ := time.Parse("2006-01-02", createdToStr)
	expectedUpdatedFrom, _ := time.Parse("2006-01-02", updatedFromStr)
	expectedUpdatedTo, _ := time.Parse("2006-01-02", updatedToStr)

	assert.Equal(t, expectedCreatedFrom, *model.CreatedFrom)
	assert.Equal(t, expectedCreatedTo, *model.CreatedTo)
	assert.Equal(t, expectedUpdatedFrom, *model.UpdatedFrom)
	assert.Equal(t, expectedUpdatedTo, *model.UpdatedTo)
}

func TestClientFilterDTO_ToModel_InvalidAndEmptyDates(t *testing.T) {
	invalidDate := "invalid-date"
	empty := ""

	dto := ClientFilterDTO{
		CreatedFrom: &invalidDate,
		CreatedTo:   &empty,
		UpdatedFrom: nil,
		UpdatedTo:   nil,
	}

	model, err := dto.ToModel()
	assert.NoError(t, err)
	assert.NotNil(t, model)

	assert.Nil(t, model.CreatedFrom)
	assert.Nil(t, model.CreatedTo)
	assert.Nil(t, model.UpdatedFrom)
	assert.Nil(t, model.UpdatedTo)
}

func TestClientFilterDTO_ToModel_NilAndZeroValues(t *testing.T) {
	dto := ClientFilterDTO{}

	model, err := dto.ToModel()
	assert.NoError(t, err)
	assert.NotNil(t, model)

	assert.Equal(t, "", model.Name)
	assert.Equal(t, "", model.Email)
	assert.Nil(t, model.Status)
	assert.Nil(t, model.Version)
	assert.Equal(t, 0, model.Limit)
	assert.Equal(t, 0, model.Offset)
	assert.Nil(t, model.CreatedFrom)
	assert.Nil(t, model.CreatedTo)
	assert.Nil(t, model.UpdatedFrom)
	assert.Nil(t, model.UpdatedTo)
}

func TestClientFilterDTO_ToModel_OnlyDatesFilled(t *testing.T) {
	dateStr := "2025-05-10"

	dto := ClientFilterDTO{
		CreatedFrom: &dateStr,
		UpdatedTo:   &dateStr,
	}

	model, err := dto.ToModel()
	assert.NoError(t, err)
	assert.NotNil(t, model)

	expectedDate, _ := time.Parse("2006-01-02", dateStr)

	assert.Equal(t, expectedDate, *model.CreatedFrom)
	assert.Equal(t, expectedDate, *model.UpdatedTo)
	assert.Nil(t, model.CreatedTo)
	assert.Nil(t, model.UpdatedFrom)
}
