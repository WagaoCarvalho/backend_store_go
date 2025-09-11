package dto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
)

func TestSaleDTO_ToSaleModel(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339)

	dtoInput := SaleDTO{
		ClientID:      nil,
		UserID:        2,
		SaleDate:      &now,
		TotalAmount:   200.0,
		TotalDiscount: 20.0,
		PaymentType:   "card",
		Status:        "active",
		Notes:         "Pedido teste",
		Version:       1,
		CreatedAt:     &now,
		UpdatedAt:     &now,
	}

	model := ToSaleModel(dtoInput)

	assert.Equal(t, dtoInput.UserID, model.UserID)
	assert.Equal(t, dtoInput.TotalAmount, model.TotalAmount)
	assert.Equal(t, dtoInput.PaymentType, model.PaymentType)
	assert.Equal(t, dtoInput.Status, model.Status)
	assert.Equal(t, dtoInput.Notes, model.Notes)
}

func TestSaleDTO_ToSaleDTO(t *testing.T) {
	now := time.Now().UTC()
	model := &models.Sale{
		ID:            1,
		ClientID:      nil,
		UserID:        3,
		SaleDate:      now,
		TotalAmount:   300.0,
		TotalDiscount: 30.0,
		PaymentType:   "credit",
		Status:        "active",
		Notes:         "Outro teste",
		Version:       2,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	dtoOutput := ToSaleDTO(model)

	assert.NotNil(t, dtoOutput.ID)
	assert.Equal(t, model.UserID, dtoOutput.UserID)
	assert.Equal(t, model.TotalAmount, dtoOutput.TotalAmount)
	assert.Equal(t, model.PaymentType, dtoOutput.PaymentType)
	assert.Equal(t, model.Status, dtoOutput.Status)
	assert.Equal(t, model.Notes, dtoOutput.Notes)
	assert.Equal(t, model.Version, dtoOutput.Version)
	assert.NotNil(t, dtoOutput.SaleDate)
	assert.NotNil(t, dtoOutput.CreatedAt)
	assert.NotNil(t, dtoOutput.UpdatedAt)
}
