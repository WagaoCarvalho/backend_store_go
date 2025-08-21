package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func strPtr(s string) *string {
	return &s
}

func TestSupplierService_Create(t *testing.T) {
	baseLogger := logrus.New()
	logger := logger.NewLoggerAdapter(baseLogger)

	setup := func() (*repo.MockSupplierRepository, SupplierService) {
		mockRepo := new(repo.MockSupplierRepository)
		service := NewSupplierService(mockRepo, logger)
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
		assert.Equal(t, ErrSupplierNameRequired, err)
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("erro repo", func(t *testing.T) {
		mockRepo, service := setup()

		input := &models.Supplier{Name: "Fornecedor X"}

		mockRepo.On("Create", mock.Anything, input).Return(nil, errors.New("erro DB")).Once()

		result, err := service.Create(context.Background(), input)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "erro ao criar fornecedor")
		mockRepo.AssertExpectations(t)
	})
}

func TestSupplierService_GetAll(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	setup := func() (*repo.MockSupplierRepository, SupplierService) {
		mockRepo := new(repo.MockSupplierRepository)
		service := NewSupplierService(mockRepo, logger)
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

		mockRepo.On("GetAll", mock.Anything).Return(([]*models.Supplier)(nil), errors.New("erro interno")).Once()

		result, err := service.GetAll(context.Background())

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, ErrGetSupplier))
		assert.Contains(t, err.Error(), "erro na validação dos dados")

		mockRepo.AssertExpectations(t)
	})

}

func TestSupplierService_GetByID(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	setup := func() (*repo.MockSupplierRepository, SupplierService) {
		mockRepo := new(repo.MockSupplierRepository)
		service := NewSupplierService(mockRepo, logger)
		return mockRepo, service
	}

	t.Run("sucesso ao buscar por ID", func(t *testing.T) {
		mockRepo, service := setup()

		expected := &models.Supplier{
			ID:   1,
			Name: "Fornecedor A",
			CNPJ: strPtr("111"),
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
		assert.Equal(t, ErrInvalidSupplierID, err)
	})

	t.Run("erro do repositório ao buscar por ID", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("GetByID", mock.Anything, int64(2)).Return(nil, errors.New("db down")).Once()

		result, err := service.GetByID(context.Background(), 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, ErrGetSupplier))
		assert.Contains(t, err.Error(), "erro na validação dos dados")
		mockRepo.AssertExpectations(t)
	})

	t.Run("fornecedor não encontrado", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("GetByID", mock.Anything, int64(10)).Return(nil, nil).Once()

		result, err := service.GetByID(context.Background(), 10)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrSupplierNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSupplierService_GetByName(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	setup := func() (*repo.MockSupplierRepository, SupplierService) {
		mockRepo := new(repo.MockSupplierRepository)
		service := NewSupplierService(mockRepo, logger)
		return mockRepo, service
	}

	t.Run("sucesso ao buscar por nome", func(t *testing.T) {
		mockRepo, service := setup()

		expected := []*models.Supplier{
			{
				ID:   1,
				Name: "Fornecedor A",
				CNPJ: strPtr("111"),
			},
			{
				ID:   2,
				Name: "Fornecedor AB",
				CNPJ: strPtr("222"),
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
		assert.Equal(t, ErrInvalidSupplierName, err)
	})

	t.Run("erro do repositório ao buscar por nome", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("GetByName", mock.Anything, "Fornecedor").Return(nil, errors.New("db down")).Once()

		result, err := service.GetByName(context.Background(), "Fornecedor")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.True(t, errors.Is(err, ErrGetSupplier))
		mockRepo.AssertExpectations(t)
	})

	t.Run("fornecedor não encontrado", func(t *testing.T) {
		mockRepo, service := setup()

		mockRepo.On("GetByName", mock.Anything, "Inexistente").Return(nil, nil).Once()

		result, err := service.GetByName(context.Background(), "Inexistente")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, ErrSupplierNotFound, err)
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
		mockRepo       func(m *repo.MockSupplierRepository)
		expectedResult int64
		expectedErr    error
	}{
		{
			name: "sucesso ao obter versão",
			args: args{id: 1},
			mockRepo: func(m *repo.MockSupplierRepository) {
				m.On("GetVersionByID", mock.Anything, int64(1)).Return(int64(3), nil).Once()
			},
			expectedResult: 3,
			expectedErr:    nil,
		},
		{
			name: "id inválido",
			args: args{id: 0},
			mockRepo: func(m *repo.MockSupplierRepository) {
				// não deve chamar o repo
			},
			expectedResult: 0,
			expectedErr:    ErrInvalidSupplierID,
		},
		{
			name: "erro ao buscar versão",
			args: args{id: 2},
			mockRepo: func(m *repo.MockSupplierRepository) {
				m.On("GetVersionByID", mock.Anything, int64(2)).Return(int64(0), errors.New("erro banco")).Once()
			},
			expectedResult: 0,
			expectedErr:    ErrGetSupplierVersion, // testará se está contido na composição do erro
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repo.MockSupplierRepository{}
			if tt.mockRepo != nil {
				tt.mockRepo(mockRepo)
			}

			logger := logger.NewLoggerAdapter(logrus.New())

			service := NewSupplierService(mockRepo, logger)

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
	type args struct {
		supplier *models.Supplier
	}

	tests := []struct {
		name        string
		args        args
		mockRepo    func(m *repo.MockSupplierRepository)
		expected    *models.Supplier
		expectedErr error
	}{
		{
			name: "sucesso na atualização",
			args: args{
				supplier: &models.Supplier{ID: 1, Name: "Fornecedor A", Version: 1},
			},
			mockRepo: func(m *repo.MockSupplierRepository) {
				m.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
					return s.ID == 1 && s.Name == "Fornecedor A" && s.Version == 1
				})).Return(nil).Once()
			},
			expected:    &models.Supplier{ID: 1, Name: "Fornecedor A", Version: 1},
			expectedErr: nil,
		},
		{
			name: "id inválido",
			args: args{
				supplier: &models.Supplier{ID: 0, Name: "Fornecedor B", Version: 1},
			},
			mockRepo:    nil,
			expected:    nil,
			expectedErr: ErrInvalidSupplierID,
		},
		{
			name: "nome obrigatório",
			args: args{
				supplier: &models.Supplier{ID: 1, Name: "", Version: 1},
			},
			mockRepo:    nil,
			expected:    nil,
			expectedErr: ErrSupplierNameRequired,
		},
		{
			name: "versão obrigatória",
			args: args{
				supplier: &models.Supplier{ID: 1, Name: "Fornecedor C", Version: 0},
			},
			mockRepo:    nil,
			expected:    nil,
			expectedErr: ErrSupplierVersionRequired,
		},
		{
			name: "conflito de versão",
			args: args{
				supplier: &models.Supplier{ID: 1, Name: "Fornecedor D", Version: 2},
			},
			mockRepo: func(m *repo.MockSupplierRepository) {
				m.On("Update", mock.Anything, mock.Anything).Return(repo.ErrVersionConflict).Once()
			},
			expected:    nil,
			expectedErr: ErrSupplierVersionConflict,
		},
		{
			name: "fornecedor não encontrado",
			args: args{
				supplier: &models.Supplier{ID: 10, Name: "Fornecedor X", Version: 1},
			},
			mockRepo: func(m *repo.MockSupplierRepository) {
				m.On("Update", mock.Anything, mock.Anything).Return(repo.ErrSupplierNotFound).Once()
			},
			expected:    nil,
			expectedErr: ErrSupplierNotFound,
		},
		{
			name: "erro genérico na atualização",
			args: args{
				supplier: &models.Supplier{ID: 1, Name: "Fornecedor Z", Version: 1},
			},
			mockRepo: func(m *repo.MockSupplierRepository) {
				m.On("Update", mock.Anything, mock.Anything).Return(errors.New("erro banco")).Once()
			},
			expected:    nil,
			expectedErr: ErrSupplierUpdate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repo.MockSupplierRepository{}
			if tt.mockRepo != nil {
				tt.mockRepo(mockRepo)
			}

			logger := logger.NewLoggerAdapter(logrus.New())
			service := NewSupplierService(mockRepo, logger)

			result, err := service.Update(context.Background(), tt.args.supplier)

			assert.Equal(t, tt.expected, result)

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

func TestSupplierService_Delete(t *testing.T) {
	type args struct {
		id int64
	}

	tests := []struct {
		name        string
		args        args
		mockRepo    func(m *repo.MockSupplierRepository)
		expectedErr error
	}{
		{
			name: "sucesso ao deletar",
			args: args{id: 1},
			mockRepo: func(m *repo.MockSupplierRepository) {
				m.On("Delete", mock.Anything, int64(1)).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "id inválido para deleção",
			args: args{id: 0},
			mockRepo: func(m *repo.MockSupplierRepository) {
				// não deve chamar Delete
			},
			expectedErr: ErrInvalidSupplierIDForDeletion,
		},
		{
			name: "erro ao deletar",
			args: args{id: 2},
			mockRepo: func(m *repo.MockSupplierRepository) {
				m.On("Delete", mock.Anything, int64(2)).Return(errors.New("erro banco")).Once()
			},
			expectedErr: errors.New("erro banco"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repo.MockSupplierRepository{}
			if tt.mockRepo != nil {
				tt.mockRepo(mockRepo)
			}

			logger := logger.NewLoggerAdapter(logrus.New())
			service := NewSupplierService(mockRepo, logger)

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
		mockRepo    func(m *repo.MockSupplierRepository)
		expectedErr error
	}{
		{
			name: "sucesso ao desabilitar fornecedor",
			args: args{id: 1},
			mockRepo: func(m *repo.MockSupplierRepository) {
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
			mockRepo: func(m *repo.MockSupplierRepository) {
				// não deve chamar o repo
			},
			expectedErr: ErrInvalidSupplierID,
		},
		{
			name: "erro ao obter fornecedor",
			args: args{id: 2},
			mockRepo: func(m *repo.MockSupplierRepository) {
				m.On("GetByID", mock.Anything, int64(2)).Return(nil, errors.New("erro banco")).Once()
			},
			expectedErr: fmt.Errorf("%w: %v", ErrGetSupplier, errors.New("erro banco")),
		},
		{
			name: "erro ao atualizar fornecedor",
			args: args{id: 3},
			mockRepo: func(m *repo.MockSupplierRepository) {
				m.On("GetByID", mock.Anything, int64(3)).Return(&models.Supplier{
					ID:     3,
					Status: true,
				}, nil).Once()
				m.On("Update", mock.Anything, mock.Anything).Return(errors.New("erro update")).Once()
			},
			expectedErr: fmt.Errorf("%w: %v", ErrDisableSupplier, errors.New("erro update")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repo.MockSupplierRepository{}
			if tt.mockRepo != nil {
				tt.mockRepo(mockRepo)
			}

			logger := logger.NewLoggerAdapter(logrus.New())
			service := NewSupplierService(mockRepo, logger)

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
		mockRepo    func(m *repo.MockSupplierRepository)
		expectedErr error
	}{
		{
			name: "sucesso ao habilitar fornecedor",
			args: args{id: 1},
			mockRepo: func(m *repo.MockSupplierRepository) {
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
			mockRepo: func(m *repo.MockSupplierRepository) {
				// não deve chamar o repo
			},
			expectedErr: ErrInvalidSupplierID,
		},
		{
			name: "erro ao obter fornecedor",
			args: args{id: 2},
			mockRepo: func(m *repo.MockSupplierRepository) {
				m.On("GetByID", mock.Anything, int64(2)).Return(nil, errors.New("erro banco")).Once()
			},
			expectedErr: fmt.Errorf("%w: %v", ErrGetSupplier, errors.New("erro banco")),
		},
		{
			name: "erro ao atualizar fornecedor",
			args: args{id: 3},
			mockRepo: func(m *repo.MockSupplierRepository) {
				m.On("GetByID", mock.Anything, int64(3)).Return(&models.Supplier{
					ID:     3,
					Status: false,
				}, nil).Once()
				m.On("Update", mock.Anything, mock.Anything).Return(errors.New("erro update")).Once()
			},
			expectedErr: fmt.Errorf("%w: %v", ErrEnableSupplier, errors.New("erro update")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &repo.MockSupplierRepository{}
			if tt.mockRepo != nil {
				tt.mockRepo(mockRepo)
			}

			logger := logger.NewLoggerAdapter(logrus.New())
			service := NewSupplierService(mockRepo, logger)

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
