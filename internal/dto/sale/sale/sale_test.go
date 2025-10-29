package dto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
)

func TestSaleDTO_ToSaleModel(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339)
	userID := int64(2)

	dtoInput := SaleDTO{
		ClientID:      nil,
		UserID:        &userID,
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

	assert.NotNil(t, model)
	assert.Equal(t, dtoInput.ClientID, model.ClientID)
	assert.Equal(t, dtoInput.UserID, model.UserID)
	assert.Equal(t, dtoInput.TotalAmount, model.TotalAmount)
	assert.Equal(t, dtoInput.TotalDiscount, model.TotalDiscount)
	assert.Equal(t, dtoInput.PaymentType, model.PaymentType)
	assert.Equal(t, dtoInput.Status, model.Status)
	assert.Equal(t, dtoInput.Notes, model.Notes)
	assert.Equal(t, dtoInput.Version, model.Version)
	assert.False(t, model.SaleDate.IsZero())
	assert.False(t, model.CreatedAt.IsZero())
	assert.False(t, model.UpdatedAt.IsZero())
}

func TestSaleDTO_ToSaleModel_DefaultStatus(t *testing.T) {
	userID := int64(1)
	dtoInput := SaleDTO{
		UserID:      &userID,
		TotalAmount: 100,
		PaymentType: "cash",
		// Status vazio → deve virar "active"
	}

	model := ToSaleModel(dtoInput)
	assert.Equal(t, "active", model.Status)
}

func TestSaleDTO_ToSaleDTO(t *testing.T) {
	now := time.Now().UTC()
	userID := int64(3)

	model := &models.Sale{
		ID:            1,
		ClientID:      nil,
		UserID:        &userID,
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

	assert.NotNil(t, dtoOutput)
	assert.NotNil(t, dtoOutput.ID)
	assert.Equal(t, model.UserID, dtoOutput.UserID)
	assert.Equal(t, model.TotalAmount, dtoOutput.TotalAmount)
	assert.Equal(t, model.TotalDiscount, dtoOutput.TotalDiscount)
	assert.Equal(t, model.PaymentType, dtoOutput.PaymentType)
	assert.Equal(t, model.Status, dtoOutput.Status)
	assert.Equal(t, model.Notes, dtoOutput.Notes)
	assert.Equal(t, model.Version, dtoOutput.Version)
	assert.NotNil(t, dtoOutput.SaleDate)
	assert.NotNil(t, dtoOutput.CreatedAt)
	assert.NotNil(t, dtoOutput.UpdatedAt)
}

func TestSaleDTO_ToSaleDTOList(t *testing.T) {
	now := time.Now().UTC()
	userID := int64(5)

	modelsList := []*models.Sale{
		{
			ID:            10,
			ClientID:      nil,
			UserID:        &userID,
			SaleDate:      now,
			TotalAmount:   150.0,
			TotalDiscount: 5.0,
			PaymentType:   "cash",
			Status:        "active",
			Notes:         "Lista 1",
			Version:       1,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
		{
			ID:            11,
			ClientID:      nil,
			UserID:        &userID,
			SaleDate:      now,
			TotalAmount:   250.0,
			TotalDiscount: 10.0,
			PaymentType:   "card",
			Status:        "active",
			Notes:         "Lista 2",
			Version:       1,
			CreatedAt:     now,
			UpdatedAt:     now,
		},
	}

	dtoList := ToSaleDTOList(modelsList)
	assert.Len(t, dtoList, 2)
	assert.Equal(t, modelsList[0].TotalAmount, dtoList[0].TotalAmount)
	assert.Equal(t, modelsList[1].PaymentType, dtoList[1].PaymentType)
	assert.NotNil(t, dtoList[0].SaleDate)
	assert.NotNil(t, dtoList[1].CreatedAt)
}

func TestSaleDTO_SaleDTOListToModelList(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339)
	userID1 := int64(1)
	userID2 := int64(2)

	dtoList := []*SaleDTO{
		{
			ID:            nil,
			ClientID:      nil,
			UserID:        &userID1,
			SaleDate:      &now,
			TotalAmount:   100.0,
			TotalDiscount: 10.0,
			PaymentType:   "cash",
			Status:        "active",
			Notes:         "Venda 1",
			Version:       1,
			CreatedAt:     &now,
			UpdatedAt:     &now,
		},
		{
			ID:            nil,
			ClientID:      nil,
			UserID:        &userID2,
			SaleDate:      &now,
			TotalAmount:   200.0,
			TotalDiscount: 20.0,
			PaymentType:   "card",
			Status:        "active",
			Notes:         "Venda 2",
			Version:       2,
			CreatedAt:     &now,
			UpdatedAt:     &now,
		},
	}

	modelList := SaleDTOListToModelList(dtoList)

	assert.Len(t, modelList, 2)
	assert.Equal(t, dtoList[0].UserID, modelList[0].UserID)
	assert.Equal(t, dtoList[1].PaymentType, modelList[1].PaymentType)
	assert.Equal(t, dtoList[1].TotalAmount, modelList[1].TotalAmount)
	assert.False(t, modelList[0].SaleDate.IsZero())
	assert.False(t, modelList[1].UpdatedAt.IsZero())
}
