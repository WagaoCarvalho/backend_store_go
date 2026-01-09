package dto

import (
	"testing"

	modelFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSupplierFilterDTO_ToModel(t *testing.T) {

	t.Run("Converte todos os campos preenchidos corretamente com CPF", func(t *testing.T) {
		createdFrom := "2024-01-01"
		createdTo := "2024-12-31"
		updatedFrom := "2024-06-01"
		updatedTo := "2024-06-30"
		status := true

		dto := SupplierFilterDTO{
			Name:        "Fornecedor Teste",
			CPF:         "123.456.789-00",
			Status:      &status,
			CreatedFrom: &createdFrom,
			CreatedTo:   &createdTo,
			UpdatedFrom: &updatedFrom,
			UpdatedTo:   &updatedTo,
			Limit:       20,
			Offset:      10,
		}

		model, err := dto.ToModel()
		require.NoError(t, err)
		require.NotNil(t, model)

		assert.Equal(t, "Fornecedor Teste", model.Name)
		assert.Equal(t, "123.456.789-00", model.CPF)
		assert.Empty(t, model.CNPJ)
		assert.Equal(t, &status, model.Status)
		assert.Equal(t, modelFilter.BaseFilter{Limit: 20, Offset: 10}, model.BaseFilter)

		require.NotNil(t, model.CreatedFrom)
		require.NotNil(t, model.CreatedTo)
		require.NotNil(t, model.UpdatedFrom)
		require.NotNil(t, model.UpdatedTo)

		assert.Equal(t, "2024-01-01", model.CreatedFrom.Format("2006-01-02"))
		assert.Equal(t, "2024-12-31", model.CreatedTo.Format("2006-01-02"))
		assert.Equal(t, "2024-06-01", model.UpdatedFrom.Format("2006-01-02"))
		assert.Equal(t, "2024-06-30", model.UpdatedTo.Format("2006-01-02"))
	})

	t.Run("Converte todos os campos preenchidos corretamente com CNPJ", func(t *testing.T) {
		status := false

		dto := SupplierFilterDTO{
			Name:   "Fornecedor PJ",
			CNPJ:   "12.345.678/0001-00",
			Status: &status,
			Limit:  10,
			Offset: 0,
		}

		model, err := dto.ToModel()
		require.NoError(t, err)
		require.NotNil(t, model)

		assert.Equal(t, "Fornecedor PJ", model.Name)
		assert.Equal(t, "12.345.678/0001-00", model.CNPJ)
		assert.Empty(t, model.CPF)
		assert.Equal(t, &status, model.Status)
		assert.Equal(t, modelFilter.BaseFilter{Limit: 10, Offset: 0}, model.BaseFilter)
	})

	t.Run("Erro ao informar CPF e CNPJ juntos", func(t *testing.T) {
		dto := SupplierFilterDTO{
			CPF:   "123.456.789-00",
			CNPJ:  "12.345.678/0001-00",
			Limit: 10,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
	})

	t.Run("Erro com created_from após created_to", func(t *testing.T) {
		createdFrom := "2024-12-31"
		createdTo := "2024-01-01"

		dto := SupplierFilterDTO{
			CreatedFrom: &createdFrom,
			CreatedTo:   &createdTo,
			Limit:       10,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
	})

	t.Run("Erro com updated_from após updated_to", func(t *testing.T) {
		updatedFrom := "2024-12-31"
		updatedTo := "2024-01-01"

		dto := SupplierFilterDTO{
			UpdatedFrom: &updatedFrom,
			UpdatedTo:   &updatedTo,
			Limit:       10,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
	})

	t.Run("Erro com limit inválido", func(t *testing.T) {
		dto := SupplierFilterDTO{
			Limit: 0,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
	})

	t.Run("Erro com offset negativo", func(t *testing.T) {
		dto := SupplierFilterDTO{
			Limit:  10,
			Offset: -1,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
	})

	t.Run("Erro ao converter data com formato inválido", func(t *testing.T) {
		invalidDate := "31-12-2024"

		dto := SupplierFilterDTO{
			CreatedFrom: &invalidDate,
			Limit:       10,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "campo 'created_from'")
		assert.Contains(t, err.Error(), "formato esperado: YYYY-MM-DD")
	})

	t.Run("Erro quando limit é maior que 100", func(t *testing.T) {
		dto := SupplierFilterDTO{
			Limit: 101,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "limit")
		assert.Contains(t, err.Error(), "máximo é 100")
	})

	t.Run("Erro ao converter created_to com formato inválido", func(t *testing.T) {
		invalidDate := "2024/12/31"

		dto := SupplierFilterDTO{
			CreatedTo: &invalidDate,
			Limit:     10,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "campo 'created_to'")
		assert.Contains(t, err.Error(), "formato esperado: YYYY-MM-DD")
	})

	t.Run("Erro ao converter updated_from com formato inválido", func(t *testing.T) {
		invalidDate := "01-06-2024"

		dto := SupplierFilterDTO{
			UpdatedFrom: &invalidDate,
			Limit:       10,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "campo 'updated_from'")
		assert.Contains(t, err.Error(), "formato esperado: YYYY-MM-DD")
	})

	t.Run("Erro ao converter updated_to com formato inválido", func(t *testing.T) {
		invalidDate := "30/06/2024"

		dto := SupplierFilterDTO{
			UpdatedTo: &invalidDate,
			Limit:     10,
		}

		model, err := dto.ToModel()
		require.Error(t, err)
		assert.Nil(t, model)
		assert.Contains(t, err.Error(), "campo 'updated_to'")
		assert.Contains(t, err.Error(), "formato esperado: YYYY-MM-DD")
	})

}
