package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	mock_supplier "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/supplier"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	convert "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSupplierService_Create(t *testing.T) {

	setup := func() (*mock_supplier.MockSupplierRepository, Supplier) {
		mockRepo := new(mock_supplier.MockSupplierRepository)
		service := NewSupplier(mockRepo)
		return mockRepo, service
	}

	t.Run("sucesso", func(t *testing.T) {
		mockRepo, service := setup()

		input := &models.Supplier{Name: "Fornecedor X"}
		expected := &models.Supplier{ID: 1, Name: "Fornecedor X"}

		mockRepo.On("Create", mock.Anything, input).Return(expected, nil).Once()

		result, err := service.Create(context.Background(), input)

		assert.NoError(t, err)
		assert.Equal(t, expected, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro nome vazio", func(t *testing.T) {
		mockRepo, service := setup()

		input := &models.Supplier{Name: ""}

		result, err := service.Create(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData) // em vez de Contains na string
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("erro repo", func(t *testing.T) {
		mockRepo, service := setup()

		input := &models.Supplier{Name: "Fornecedor X"}

		mockRepo.On("Create", mock.Anything, input).Return(nil, errors.New("erro DB")).Once()

		result, err := service.Create(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "erro ao criar")
		mockRepo.AssertExpectations(t)
	})
}

func TestSupplierService_GetAll(t *testing.T) {

	setup := func() (*mock_supplier.MockSupplierRepository, Supplier) {
		mockRepo := new(mock_supplier.MockSupplierRepository)
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

	setup := func() (*mock_supplier.MockSupplierRepository, Supplier) {
		mockRepo := new(mock_supplier.MockSupplierRepository)
		service := NewSupplier(mockRepo)
		return mockRepo, service
	}

	t.Run("sucesso ao buscar por ID", func(t *testing.T) {
		mockRepo, service := setup()

		expected := &models.Supplier{
			ID:   1,
			Name: "Fornecedor A",
			CNPJ: convert.StrToPtr("111"),
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

	setup := func() (*mock_supplier.MockSupplierRepository, Supplier) {
		mockRepo := new(mock_supplier.MockSupplierRepository)
		service := NewSupplier(mockRepo)
		return mockRepo, service
	}

	t.Run("sucesso ao buscar por nome", func(t *testing.T) {
		mockRepo, service := setup()

		expected := []*models.Supplier{
			{
				ID:   1,
				Name: "Fornecedor A",
				CNPJ: convert.StrToPtr("111"),
			},
			{
				ID:   2,
				Name: "Fornecedor AB",
				CNPJ: convert.StrToPtr("222"),
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
		mockRepo       func(m *mock_supplier.MockSupplierRepository)
		expectedResult int64
		expectedErr    error
	}{
		{
			name: "sucesso ao obter versão",
			args: args{id: 1},
			mockRepo: func(m *mock_supplier.MockSupplierRepository) {
				m.On("GetVersionByID", mock.Anything, int64(1)).Return(int64(3), nil).Once()
			},
			expectedResult: 3,
			expectedErr:    nil,
		},
		{
			name: "id inválido",
			args: args{id: 0},
			mockRepo: func(_ *mock_supplier.MockSupplierRepository) {
				// não deve chamar o repo
			},
			expectedResult: 0,
			expectedErr:    errMsg.ErrZeroID,
		},
		{
			name: "erro ao buscar versão",
			args: args{id: 2},
			mockRepo: func(m *mock_supplier.MockSupplierRepository) {
				m.On("GetVersionByID", mock.Anything, int64(2)).Return(int64(0), errors.New("erro ao buscar versão")).Once()
			},
			expectedResult: 0,
			expectedErr:    errMsg.ErrGetVersion,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock_supplier.MockSupplierRepository{}
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

func TestSupplierService_Update(t *testing.T) {
	ctx := context.Background()

	validSupplier := func() *models.Supplier {
		return &models.Supplier{
			ID:      1,
			Name:    "Fornecedor Teste",
			Version: 1,
		}
	}

	t.Run("falha: ID inválido", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplierRepository)
		service := NewSupplier(mockRepo)

		input := validSupplier()
		input.ID = 0

		err := service.Update(ctx, input)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("falha: validação inválida", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplierRepository)
		service := NewSupplier(mockRepo)

		input := validSupplier()
		input.Name = ""

		err := service.Update(ctx, input)

		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("falha: versão inválida", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplierRepository)
		service := NewSupplier(mockRepo)

		input := validSupplier()
		input.Version = 0

		err := service.Update(ctx, input)

		assert.ErrorIs(t, err, errMsg.ErrVersionConflict)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("falha: fornecedor não encontrado (SupplierExists = false)", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplierRepository)
		service := NewSupplier(mockRepo)

		input := validSupplier()

		mockRepo.On("Update", ctx, input).Return(errMsg.ErrNotFound).Once()
		mockRepo.On("SupplierExists", ctx, input.ID).Return(false, nil).Once()

		err := service.Update(ctx, input)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha: conflito de versão (SupplierExists = true)", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplierRepository)
		service := NewSupplier(mockRepo)

		input := validSupplier()

		mockRepo.On("Update", ctx, input).Return(errMsg.ErrNotFound).Once()
		mockRepo.On("SupplierExists", ctx, input.ID).Return(true, nil).Once()

		err := service.Update(ctx, input)

		assert.ErrorIs(t, err, errMsg.ErrVersionConflict)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha: erro ao verificar existência", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplierRepository)
		service := NewSupplier(mockRepo)

		input := validSupplier()

		mockRepo.On("Update", ctx, input).Return(errMsg.ErrNotFound).Once()
		mockRepo.On("SupplierExists", ctx, input.ID).Return(false, errors.New("erro banco")).Once()

		err := service.Update(ctx, input)

		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, "erro banco")
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha: erro genérico no repositório", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplierRepository)
		service := NewSupplier(mockRepo)

		input := validSupplier()

		mockRepo.On("Update", ctx, input).Return(errors.New("erro genérico")).Once()

		err := service.Update(ctx, input)

		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.ErrorContains(t, err, "erro genérico")
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplierRepository)
		service := NewSupplier(mockRepo)

		input := validSupplier()

		mockRepo.On("Update", ctx, input).Return(nil).Once()

		err := service.Update(ctx, input)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSupplierService_Delete(t *testing.T) {
	type args struct {
		id int64
	}

	tests := []struct {
		name        string
		args        args
		mockRepo    func(m *mock_supplier.MockSupplierRepository)
		expectedErr error
	}{
		{
			name: "sucesso ao deletar",
			args: args{id: 1},
			mockRepo: func(m *mock_supplier.MockSupplierRepository) {
				m.On("Delete", mock.Anything, int64(1)).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "id inválido para deleção",
			args: args{id: 0},
			mockRepo: func(_ *mock_supplier.MockSupplierRepository) {
				// não deve chamar Delete
			},
			expectedErr: errMsg.ErrZeroID,
		},
		{
			name: "erro ao deletar",
			args: args{id: 2},
			mockRepo: func(m *mock_supplier.MockSupplierRepository) {
				m.On("Delete", mock.Anything, int64(2)).Return(errors.New("erro banco")).Once()
			},
			expectedErr: errors.New("erro banco"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock_supplier.MockSupplierRepository{}
			if tt.mockRepo != nil {
				tt.mockRepo(mockRepo)
			}

			service := NewSupplier(mockRepo)

			err := service.Delete(context.Background(), tt.args.id)

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

func TestSupplierService_Disable(t *testing.T) {
	type args struct {
		id int64
	}

	tests := []struct {
		name        string
		args        args
		mockRepo    func(m *mock_supplier.MockSupplierRepository)
		expectedErr error
	}{
		{
			name: "sucesso ao desabilitar fornecedor",
			args: args{id: 1},
			mockRepo: func(m *mock_supplier.MockSupplierRepository) {
				m.On("GetByID", mock.Anything, int64(1)).Return(&models.Supplier{
					ID:     1,
					Status: true,
				}, nil).Once()
				m.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
					return s.ID == 1 && s.Status == false
				})).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "id inválido para desabilitar",
			args: args{id: 0},
			mockRepo: func(_ *mock_supplier.MockSupplierRepository) {
				// não deve chamar o repo
			},
			expectedErr: errMsg.ErrZeroID,
		},
		{
			name: "erro ao obter fornecedor",
			args: args{id: 2},
			mockRepo: func(m *mock_supplier.MockSupplierRepository) {
				m.On("GetByID", mock.Anything, int64(2)).Return(nil, errors.New("erro banco")).Once()
			},
			expectedErr: fmt.Errorf("%w: %v", errMsg.ErrGet, errors.New("erro banco")),
		},
		{
			name: "erro ao atualizar fornecedor",
			args: args{id: 3},
			mockRepo: func(m *mock_supplier.MockSupplierRepository) {
				m.On("GetByID", mock.Anything, int64(3)).Return(&models.Supplier{
					ID:     3,
					Status: true,
				}, nil).Once()
				m.On("Update", mock.Anything, mock.Anything).Return(errors.New("erro update")).Once()
			},
			expectedErr: fmt.Errorf("%w: %v", errMsg.ErrDisable, errors.New("erro update")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock_supplier.MockSupplierRepository{}
			if tt.mockRepo != nil {
				tt.mockRepo(mockRepo)
			}

			service := NewSupplier(mockRepo)

			err := service.Disable(context.Background(), tt.args.id)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestSupplierService_Enable(t *testing.T) {
	type args struct {
		id int64
	}

	tests := []struct {
		name        string
		args        args
		mockRepo    func(m *mock_supplier.MockSupplierRepository)
		expectedErr error
	}{
		{
			name: "sucesso ao habilitar fornecedor",
			args: args{id: 1},
			mockRepo: func(m *mock_supplier.MockSupplierRepository) {
				m.On("GetByID", mock.Anything, int64(1)).Return(&models.Supplier{
					ID:     1,
					Status: false,
				}, nil).Once()
				m.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
					return s.ID == 1 && s.Status == true
				})).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "id inválido para habilitar",
			args: args{id: 0},
			mockRepo: func(_ *mock_supplier.MockSupplierRepository) {
				// não deve chamar o repo
			},
			expectedErr: errMsg.ErrZeroID,
		},
		{
			name: "erro ao obter fornecedor",
			args: args{id: 2},
			mockRepo: func(m *mock_supplier.MockSupplierRepository) {
				m.On("GetByID", mock.Anything, int64(2)).Return(nil, errors.New("erro banco")).Once()
			},
			expectedErr: fmt.Errorf("%w: %v", errMsg.ErrGet, errors.New("erro banco")),
		},
		{
			name: "erro ao atualizar fornecedor",
			args: args{id: 3},
			mockRepo: func(m *mock_supplier.MockSupplierRepository) {
				m.On("GetByID", mock.Anything, int64(3)).Return(&models.Supplier{
					ID:     3,
					Status: false,
				}, nil).Once()
				m.On("Update", mock.Anything, mock.Anything).Return(errors.New("erro update")).Once()
			},
			expectedErr: fmt.Errorf("%w: %v", errMsg.ErrEnable, errors.New("erro update")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock_supplier.MockSupplierRepository{}
			if tt.mockRepo != nil {
				tt.mockRepo(mockRepo)
			}

			service := NewSupplier(mockRepo)

			err := service.Enable(context.Background(), tt.args.id)

			if tt.expectedErr != nil {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr.Error())
			} else {
				assert.NoError(t, err)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
