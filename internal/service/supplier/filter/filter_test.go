package services

import (
	"context"
	"errors"
	"testing"
	"time"

	mockSupplier "github.com/WagaoCarvalho/backend_store_go/infra/mock/supplier"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	supplierFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/filter"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSupplierService_Filter(t *testing.T) {

	t.Run("falha quando filtro é nulo", func(t *testing.T) {
		mockRepo := new(mockSupplier.MockSupplier)
		service := NewSupplierFilterService(mockRepo)

		result, err := service.Filter(context.Background(), nil)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidFilter)
		mockRepo.AssertNotCalled(t, "Filter", mock.Anything, mock.Anything)
	})

	t.Run("falha na validação do filtro", func(t *testing.T) {
		mockRepo := new(mockSupplier.MockSupplier)
		service := NewSupplierFilterService(mockRepo)

		invalidFilter := &supplierFilter.SupplierFilter{
			BaseFilter: filter.BaseFilter{
				Limit: -1, // inválido
			},
		}

		result, err := service.Filter(context.Background(), invalidFilter)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidFilter)
		mockRepo.AssertNotCalled(t, "Filter", mock.Anything, mock.Anything)
	})

	t.Run("falha ao buscar no repositório", func(t *testing.T) {
		mockRepo := new(mockSupplier.MockSupplier)
		service := NewSupplierFilterService(mockRepo)

		validFilter := &supplierFilter.SupplierFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		dbErr := errors.New("falha no banco de dados")

		mockRepo.
			On("Filter", mock.Anything, validFilter).
			Return(nil, dbErr).
			Once()

		result, err := service.Filter(context.Background(), validFilter)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, dbErr.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso ao retornar lista de fornecedores", func(t *testing.T) {
		mockRepo := new(mockSupplier.MockSupplier)
		service := NewSupplierFilterService(mockRepo)

		validFilter := &supplierFilter.SupplierFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		now := time.Now()

		mockSuppliers := []*model.Supplier{
			{
				ID:        1,
				Name:      "Fornecedor A",
				CPF:       utils.StrToPtr("123.456.789-00"),
				Status:    true,
				Version:   1,
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        2,
				Name:      "Fornecedor B",
				CNPJ:      utils.StrToPtr("12.345.678/0001-00"),
				Status:    false,
				Version:   1,
				CreatedAt: now,
				UpdatedAt: now,
			},
		}

		mockRepo.
			On("Filter", mock.Anything, validFilter).
			Return(mockSuppliers, nil).
			Once()

		result, err := service.Filter(context.Background(), validFilter)

		assert.NoError(t, err)
		assert.Len(t, result, 2)

		assert.Equal(t, int64(1), result[0].ID)
		assert.Equal(t, "Fornecedor A", result[0].Name)
		assert.True(t, result[0].Status)

		assert.Equal(t, int64(2), result[1].ID)
		assert.Equal(t, "Fornecedor B", result[1].Name)
		assert.False(t, result[1].Status)

		mockRepo.AssertExpectations(t)
	})
}
