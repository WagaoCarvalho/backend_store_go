package services

import (
	"context"
	"errors"
	"testing"

	mock_supplier "github.com/WagaoCarvalho/backend_store_go/infra/mock/supplier"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSupplierService_Create(t *testing.T) {

	setup := func() (*mock_supplier.MockSupplier, Supplier) {
		mockRepo := new(mock_supplier.MockSupplier)
		service := NewSupplierService(mockRepo)
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
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
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
		mockRepo := new(mock_supplier.MockSupplier)
		service := NewSupplierService(mockRepo)

		input := validSupplier()
		input.ID = 0

		err := service.Update(ctx, input)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("falha: validação inválida", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplier)
		service := NewSupplierService(mockRepo)

		input := validSupplier()
		input.Name = ""

		err := service.Update(ctx, input)

		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("falha: versão inválida", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplier)
		service := NewSupplierService(mockRepo)

		input := validSupplier()
		input.Version = 0

		err := service.Update(ctx, input)

		assert.ErrorIs(t, err, errMsg.ErrVersionConflict)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("falha: fornecedor não encontrado (SupplierExists = false)", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplier)
		service := NewSupplierService(mockRepo)

		input := validSupplier()

		mockRepo.On("Update", ctx, input).Return(errMsg.ErrNotFound).Once()
		mockRepo.On("SupplierExists", ctx, input.ID).Return(false, nil).Once()

		err := service.Update(ctx, input)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha: conflito de versão (SupplierExists = true)", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplier)
		service := NewSupplierService(mockRepo)

		input := validSupplier()

		mockRepo.On("Update", ctx, input).Return(errMsg.ErrNotFound).Once()
		mockRepo.On("SupplierExists", ctx, input.ID).Return(true, nil).Once()

		err := service.Update(ctx, input)

		assert.ErrorIs(t, err, errMsg.ErrVersionConflict)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha: erro ao verificar existência", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplier)
		service := NewSupplierService(mockRepo)

		input := validSupplier()

		mockRepo.On("Update", ctx, input).Return(errMsg.ErrNotFound).Once()
		mockRepo.On("SupplierExists", ctx, input.ID).Return(false, errors.New("erro banco")).Once()

		err := service.Update(ctx, input)

		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, "erro banco")
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha: erro genérico no repositório", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplier)
		service := NewSupplierService(mockRepo)

		input := validSupplier()

		mockRepo.On("Update", ctx, input).Return(errors.New("erro genérico")).Once()

		err := service.Update(ctx, input)

		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.ErrorContains(t, err, "erro genérico")
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplier)
		service := NewSupplierService(mockRepo)

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
		mockRepo    func(m *mock_supplier.MockSupplier)
		expectedErr error
	}{
		{
			name: "sucesso ao deletar",
			args: args{id: 1},
			mockRepo: func(m *mock_supplier.MockSupplier) {
				m.On("Delete", mock.Anything, int64(1)).Return(nil).Once()
			},
			expectedErr: nil,
		},
		{
			name: "id inválido para deleção",
			args: args{id: 0},
			mockRepo: func(_ *mock_supplier.MockSupplier) {
				// não deve chamar Delete
			},
			expectedErr: errMsg.ErrZeroID,
		},
		{
			name: "erro ao deletar",
			args: args{id: 2},
			mockRepo: func(m *mock_supplier.MockSupplier) {
				m.On("Delete", mock.Anything, int64(2)).Return(errors.New("erro banco")).Once()
			},
			expectedErr: errors.New("erro banco"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mock_supplier.MockSupplier{}
			if tt.mockRepo != nil {
				tt.mockRepo(mockRepo)
			}

			service := NewSupplierService(mockRepo)

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
