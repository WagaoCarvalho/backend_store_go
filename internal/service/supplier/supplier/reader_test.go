package services

import (
	"context"
	"errors"
	"testing"

	mockSupplier "github.com/WagaoCarvalho/backend_store_go/infra/mock/supplier"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSupplierService_GetAll(t *testing.T) {

	setup := func() (*mockSupplier.MockSupplier, Supplier) {
		mockRepo := new(mockSupplier.MockSupplier)
		service := NewSupplier(mockRepo)
		return mockRepo, service
	}

	t.Run("sucesso ao buscar fornecedores", func(t *testing.T) {
		mockRepo, service := setup()

		expected := []*models.Supplier{
			{ID: 1, Name: "Fornecedor 1"},
			{ID: 2, Name: "Fornecedor 2"},
		}

		mockRepo.On("GetAll", mock.Anything).Return(expected, nil).Once()

		result, err := service.GetAll(context.Background())

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, expected, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("erro ao buscar fornecedores", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("GetAll", mock.Anything).Return(([]*models.Supplier)(nil), errors.New("erro na validação dos dados")).Once()

		result, err := service.GetAll(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, errMsg.ErrGet))
		assert.Contains(t, err.Error(), "erro na validação dos dados")

		mockRepo.AssertExpectations(t)
	})

}

func TestSupplierService_GetByID(t *testing.T) {

	setup := func() (*mockSupplier.MockSupplier, Supplier) {
		mockRepo := new(mockSupplier.MockSupplier)
		service := NewSupplier(mockRepo)
		return mockRepo, service
	}

	t.Run("sucesso ao buscar por ID", func(t *testing.T) {
		mockRepo, service := setup()

		expected := &models.Supplier{
			ID:   1,
			Name: "Fornecedor A",
			CNPJ: utils.StrToPtr("111"),
		}

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expected, nil).Once()

		result, err := service.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro ao buscar por ID inválido", func(t *testing.T) {
		_, service := setup()

		result, err := service.GetByID(context.Background(), 0)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, errMsg.ErrZeroID, err)
	})

	t.Run("erro do repositório ao buscar por ID", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("GetByID", mock.Anything, int64(2)).Return(nil, errors.New("erro na validação dos dados")).Once()

		result, err := service.GetByID(context.Background(), 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, errMsg.ErrGet))
		assert.Contains(t, err.Error(), "erro na validação dos dados")
		mockRepo.AssertExpectations(t)
	})

	t.Run("fornecedor não encontrado", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("GetByID", mock.Anything, int64(10)).Return(nil, nil).Once()

		result, err := service.GetByID(context.Background(), 10)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, errMsg.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSupplierService_GetByName(t *testing.T) {

	setup := func() (*mockSupplier.MockSupplier, Supplier) {
		mockRepo := new(mockSupplier.MockSupplier)
		service := NewSupplier(mockRepo)
		return mockRepo, service
	}

	t.Run("sucesso ao buscar por nome", func(t *testing.T) {
		mockRepo, service := setup()

		expected := []*models.Supplier{
			{
				ID:   1,
				Name: "Fornecedor A",
				CNPJ: utils.StrToPtr("111"),
			},
			{
				ID:   2,
				Name: "Fornecedor AB",
				CNPJ: utils.StrToPtr("222"),
			},
		}

		mockRepo.On("GetByName", mock.Anything, "Fornecedor").Return(expected, nil).Once()

		result, err := service.GetByName(context.Background(), "Fornecedor")

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro ao buscar por nome inválido (vazio)", func(t *testing.T) {
		_, service := setup()

		result, err := service.GetByName(context.Background(), "")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, errMsg.ErrInvalidData, err)
	})

	t.Run("erro do repositório ao buscar por nome", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("GetByName", mock.Anything, "Fornecedor").Return(nil, errors.New("db down")).Once()

		result, err := service.GetByName(context.Background(), "Fornecedor")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, errMsg.ErrGet))
		mockRepo.AssertExpectations(t)
	})

	t.Run("fornecedor não encontrado", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("GetByName", mock.Anything, "Inexistente").Return(nil, nil).Once()

		result, err := service.GetByName(context.Background(), "Inexistente")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, errMsg.ErrNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSupplierService_GetVersionByID(t *testing.T) {
	type args struct {
		id int64
	}

	tests := []struct {
		name           string
		args           args
		mockRepo       func(m *mockSupplier.MockSupplier)
		expectedResult int64
		expectedErr    error
	}{
		{
			name: "sucesso ao obter versão",
			args: args{id: 1},
			mockRepo: func(m *mockSupplier.MockSupplier) {
				m.On("GetVersionByID", mock.Anything, int64(1)).Return(int64(3), nil).Once()
			},
			expectedResult: 3,
			expectedErr:    nil,
		},
		{
			name: "id inválido",
			args: args{id: 0},
			mockRepo: func(_ *mockSupplier.MockSupplier) {
				// não deve chamar o repo
			},
			expectedResult: 0,
			expectedErr:    errMsg.ErrZeroID,
		},
		{
			name: "erro ao buscar versão",
			args: args{id: 2},
			mockRepo: func(m *mockSupplier.MockSupplier) {
				m.On("GetVersionByID", mock.Anything, int64(2)).Return(int64(0), errors.New("erro ao buscar versão")).Once()
			},
			expectedResult: 0,
			expectedErr:    errMsg.ErrGetVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockSupplier.MockSupplier{}
			if tt.mockRepo != nil {
				tt.mockRepo(mockRepo)
			}

			service := NewSupplier(mockRepo)

			result, err := service.GetVersionByID(context.Background(), tt.args.id)

			assert.Equal(t, tt.expectedResult, result)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSupplierService_SupplierExists(t *testing.T) {

	setup := func() (*mockSupplier.MockSupplier, Supplier) {
		mockRepo := new(mockSupplier.MockSupplier)
		service := NewSupplier(mockRepo)
		return mockRepo, service
	}

	t.Run("sucesso ao verificar existência do fornecedor", func(t *testing.T) {
		mockRepo, service := setup()
		supplierID := int64(1)

		mockRepo.On("SupplierExists", mock.Anything, supplierID).Return(true, nil).Once()

		exists, err := service.SupplierExists(context.Background(), supplierID)

		assert.NoError(t, err)
		assert.True(t, exists)

		mockRepo.AssertExpectations(t)
	})

	t.Run("retorna falso quando fornecedor não existe", func(t *testing.T) {
		mockRepo, service := setup()
		supplierID := int64(2)

		mockRepo.On("SupplierExists", mock.Anything, supplierID).Return(false, nil).Once()

		exists, err := service.SupplierExists(context.Background(), supplierID)

		assert.NoError(t, err)
		assert.False(t, exists)

		mockRepo.AssertExpectations(t)
	})

	t.Run("erro ao verificar existência do fornecedor", func(t *testing.T) {
		mockRepo, service := setup()
		supplierID := int64(3)
		dbErr := errors.New("falha no banco de dados")

		mockRepo.On("SupplierExists", mock.Anything, supplierID).Return(false, dbErr).Once()

		exists, err := service.SupplierExists(context.Background(), supplierID)

		assert.Error(t, err)
		assert.False(t, exists)
		assert.True(t, errors.Is(err, errMsg.ErrGet))
		assert.Contains(t, err.Error(), dbErr.Error())

		mockRepo.AssertExpectations(t)
	})

	t.Run("erro quando supplierID é zero ou negativo", func(t *testing.T) {
		_, service := setup()

		exists, err := service.SupplierExists(context.Background(), 0)

		assert.Error(t, err)
		assert.False(t, exists)
		assert.True(t, errors.Is(err, errMsg.ErrZeroID))
	})
}
