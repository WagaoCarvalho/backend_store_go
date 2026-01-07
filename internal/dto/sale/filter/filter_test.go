package dto

import (
	"testing"

	modelFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	"github.com/stretchr/testify/assert"
)

func TestSaleFilterDTO_ToModel(t *testing.T) {
	t.Run("Converte todos os campos preenchidos corretamente", func(t *testing.T) {
		clientID := int64(100)
		userID := int64(50)
		minTotalAmount := "150.75"
		maxTotalAmount := "500.00"
		minItemsAmount := "100.00"
		maxItemsAmount := "300.00"
		minItemsDiscount := "5.0"
		maxItemsDiscount := "15.5"
		minSaleDiscount := "10.0"
		maxSaleDiscount := "20.0"
		saleDateFrom := "2024-01-01"
		saleDateTo := "2024-01-31"
		createdFrom := "2024-01-01"
		createdTo := "2024-12-31"
		updatedFrom := "2024-06-01"
		updatedTo := "2024-06-30"

		dto := SaleFilterDTO{
			ClientID:         &clientID,
			UserID:           &userID,
			PaymentType:      "credit_card",
			Status:           "completed",
			MinTotalAmount:   &minTotalAmount,
			MaxTotalAmount:   &maxTotalAmount,
			MinItemsAmount:   &minItemsAmount,
			MaxItemsAmount:   &maxItemsAmount,
			MinItemsDiscount: &minItemsDiscount,
			MaxItemsDiscount: &maxItemsDiscount,
			MinSaleDiscount:  &minSaleDiscount,
			MaxSaleDiscount:  &maxSaleDiscount,
			SaleDateFrom:     &saleDateFrom,
			SaleDateTo:       &saleDateTo,
			CreatedFrom:      &createdFrom,
			CreatedTo:        &createdTo,
			UpdatedFrom:      &updatedFrom,
			UpdatedTo:        &updatedTo,
			Limit:            20,
			Offset:           10,
		}

		model, err := dto.ToModel()
		assert.NoError(t, err)

		// Verificar campos básicos
		assert.Equal(t, &clientID, model.ClientID)
		assert.Equal(t, &userID, model.UserID)
		assert.Equal(t, "credit_card", model.PaymentType)
		assert.Equal(t, "completed", model.Status)
		assert.Equal(t, modelFilter.BaseFilter{Limit: 20, Offset: 10}, model.BaseFilter)

		// Verificar valores float - CORREÇÃO AQUI
		assert.Equal(t, float64(150.75), *model.MinTotalAmount)
		assert.Equal(t, float64(500.00), *model.MaxTotalAmount)

		// Correção: verificar MinTotalItemsAmount e MaxTotalItemsAmount separadamente
		assert.Equal(t, float64(100.00), *model.MinTotalItemsAmount) // Antes estava MaxTotalItemsAmount
		assert.Equal(t, float64(300.00), *model.MaxTotalItemsAmount) // OK

		assert.Equal(t, float64(5.0), *model.MinTotalItemsDiscount)
		assert.Equal(t, float64(15.5), *model.MaxTotalItemsDiscount)
		assert.Equal(t, float64(10.0), *model.MinTotalSaleDiscount)
		assert.Equal(t, float64(20.0), *model.MaxTotalSaleDiscount)

		// Verificar datas
		assert.NotNil(t, model.SaleDateFrom)
		assert.NotNil(t, model.SaleDateTo)
		assert.NotNil(t, model.CreatedFrom)
		assert.NotNil(t, model.CreatedTo)
		assert.NotNil(t, model.UpdatedFrom)
		assert.NotNil(t, model.UpdatedTo)

		assert.Equal(t, "2024-01-01", model.SaleDateFrom.Format("2006-01-02"))
		assert.Equal(t, "2024-01-31", model.SaleDateTo.Format("2006-01-02"))
		assert.Equal(t, "2024-01-01", model.CreatedFrom.Format("2006-01-02"))
		assert.Equal(t, "2024-12-31", model.CreatedTo.Format("2006-01-02"))
		assert.Equal(t, "2024-06-01", model.UpdatedFrom.Format("2006-01-02"))
		assert.Equal(t, "2024-06-30", model.UpdatedTo.Format("2006-01-02"))
	})

	t.Run("Retorna nil para campos vazios ou não preenchidos", func(t *testing.T) {
		dto := SaleFilterDTO{}
		model, err := dto.ToModel()
		assert.NoError(t, err)

		// Verificar que ponteiros são nil
		assert.Nil(t, model.ClientID)
		assert.Nil(t, model.UserID)
		assert.Nil(t, model.MinTotalAmount)
		assert.Nil(t, model.MaxTotalAmount)
		assert.Nil(t, model.MinTotalItemsAmount)
		assert.Nil(t, model.MaxTotalItemsAmount)
		assert.Nil(t, model.MinTotalItemsDiscount)
		assert.Nil(t, model.MaxTotalItemsDiscount)
		assert.Nil(t, model.MinTotalSaleDiscount)
		assert.Nil(t, model.MaxTotalSaleDiscount)
		assert.Nil(t, model.SaleDateFrom)
		assert.Nil(t, model.SaleDateTo)
		assert.Nil(t, model.CreatedFrom)
		assert.Nil(t, model.CreatedTo)
		assert.Nil(t, model.UpdatedFrom)
		assert.Nil(t, model.UpdatedTo)

		// Verificar campos de string vazios
		assert.Equal(t, "", model.PaymentType)
		assert.Equal(t, "", model.Status)

		// Verificar BaseFilter com valores padrão
		assert.Equal(t, 0, model.Limit)
		assert.Equal(t, 0, model.Offset)
	})

	t.Run("Ignora valores inválidos em campos de parse", func(t *testing.T) {
		invalidFloat := "abc"
		invalidDate := "31/12/2024"
		validClientID := int64(1)
		validUserID := int64(2)

		dto := SaleFilterDTO{
			ClientID:       &validClientID,
			UserID:         &validUserID,
			PaymentType:    "valid_payment",
			Status:         "valid_status",
			MinTotalAmount: &invalidFloat,
			MaxTotalAmount: &invalidFloat,
			CreatedFrom:    &invalidDate,
			Limit:          10,
			Offset:         0,
		}

		model, err := dto.ToModel()
		assert.NoError(t, err)

		// Campos válidos devem ser mantidos
		assert.Equal(t, &validClientID, model.ClientID)
		assert.Equal(t, &validUserID, model.UserID)
		assert.Equal(t, "valid_payment", model.PaymentType)
		assert.Equal(t, "valid_status", model.Status)
		assert.Equal(t, modelFilter.BaseFilter{Limit: 10, Offset: 0}, model.BaseFilter)

		// Campos inválidos devem ser nil
		assert.Nil(t, model.MinTotalAmount)
		assert.Nil(t, model.MaxTotalAmount)
		assert.Nil(t, model.CreatedFrom)
	})

	t.Run("Converte strings vazias para nil", func(t *testing.T) {
		emptyString := ""
		clientID := int64(1)

		dto := SaleFilterDTO{
			ClientID:       &clientID,
			MinTotalAmount: &emptyString,
			CreatedFrom:    &emptyString,
		}

		model, err := dto.ToModel()
		assert.NoError(t, err)

		// ClientID deve ser mantido
		assert.Equal(t, &clientID, model.ClientID)

		// Strings vazias devem resultar em nil
		assert.Nil(t, model.MinTotalAmount)
		assert.Nil(t, model.CreatedFrom)
	})

	t.Run("Mantém valores zero válidos", func(t *testing.T) {
		zeroFloat := "0.0"
		validDate := "2024-01-01"
		clientID := int64(0) // Zero é um ID válido

		dto := SaleFilterDTO{
			ClientID:       &clientID,
			MinTotalAmount: &zeroFloat,
			CreatedFrom:    &validDate,
		}

		model, err := dto.ToModel()
		assert.NoError(t, err)

		// Zero values devem ser mantidos
		assert.Equal(t, &clientID, model.ClientID)
		assert.Equal(t, float64(0.0), *model.MinTotalAmount)
		assert.NotNil(t, model.CreatedFrom)
		assert.Equal(t, "2024-01-01", model.CreatedFrom.Format("2006-01-02"))
	})

	t.Run("Converte apenas alguns campos preenchidos", func(t *testing.T) {
		clientID := int64(10)
		minTotalAmount := "100.00"
		saleDateFrom := "2024-03-01"

		dto := SaleFilterDTO{
			ClientID:       &clientID,
			PaymentType:    "cash",
			MinTotalAmount: &minTotalAmount,
			SaleDateFrom:   &saleDateFrom,
			Limit:          15,
			Offset:         5,
		}

		model, err := dto.ToModel()
		assert.NoError(t, err)

		// Campos preenchidos
		assert.Equal(t, &clientID, model.ClientID)
		assert.Equal(t, "cash", model.PaymentType)
		assert.Equal(t, float64(100.00), *model.MinTotalAmount)
		assert.Equal(t, "2024-03-01", model.SaleDateFrom.Format("2006-01-02"))
		assert.Equal(t, modelFilter.BaseFilter{Limit: 15, Offset: 5}, model.BaseFilter)

		// Campos não preenchidos devem ser nil
		assert.Nil(t, model.UserID)
		assert.Equal(t, "", model.Status)
		assert.Nil(t, model.MaxTotalAmount)
		assert.Nil(t, model.MinTotalItemsAmount)
		assert.Nil(t, model.MaxTotalItemsAmount)
		assert.Nil(t, model.MinTotalItemsDiscount)
		assert.Nil(t, model.MaxTotalItemsDiscount)
		assert.Nil(t, model.MinTotalSaleDiscount)
		assert.Nil(t, model.MaxTotalSaleDiscount)
		assert.Nil(t, model.SaleDateTo)
		assert.Nil(t, model.CreatedFrom)
		assert.Nil(t, model.CreatedTo)
		assert.Nil(t, model.UpdatedFrom)
		assert.Nil(t, model.UpdatedTo)
	})
}
