package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	mock_supplier "github.com/WagaoCarvalho/backend_store_go/infra/mock/supplier"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSupplierService_Disable(t *testing.T) {
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
			name: "sucesso ao desabilitar fornecedor",
			args: args{id: 1},
			mockRepo: func(m *mock_supplier.MockSupplier) {
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
			mockRepo: func(_ *mock_supplier.MockSupplier) {
				// não deve chamar o repo
			},
			expectedErr: errMsg.ErrZeroID,
		},
		{
			name: "erro ao obter fornecedor",
			args: args{id: 2},
			mockRepo: func(m *mock_supplier.MockSupplier) {
				m.On("GetByID", mock.Anything, int64(2)).Return(nil, errors.New("erro banco")).Once()
			},
			expectedErr: fmt.Errorf("%w: %v", errMsg.ErrGet, errors.New("erro banco")),
		},
		{
			name: "erro ao atualizar fornecedor",
			args: args{id: 3},
			mockRepo: func(m *mock_supplier.MockSupplier) {
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
			mockRepo := &mock_supplier.MockSupplier{}
			if tt.mockRepo != nil {
				tt.mockRepo(mockRepo)
			}

			service := NewSupplierService(mockRepo)

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
		mockRepo    func(m *mock_supplier.MockSupplier)
		expectedErr error
	}{
		{
			name: "sucesso ao habilitar fornecedor",
			args: args{id: 1},
			mockRepo: func(m *mock_supplier.MockSupplier) {
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
			mockRepo: func(_ *mock_supplier.MockSupplier) {
				// não deve chamar o repo
			},
			expectedErr: errMsg.ErrZeroID,
		},
		{
			name: "erro ao obter fornecedor",
			args: args{id: 2},
			mockRepo: func(m *mock_supplier.MockSupplier) {
				m.On("GetByID", mock.Anything, int64(2)).Return(nil, errors.New("erro banco")).Once()
			},
			expectedErr: fmt.Errorf("%w: %v", errMsg.ErrGet, errors.New("erro banco")),
		},
		{
			name: "erro ao atualizar fornecedor",
			args: args{id: 3},
			mockRepo: func(m *mock_supplier.MockSupplier) {
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
			mockRepo := &mock_supplier.MockSupplier{}
			if tt.mockRepo != nil {
				tt.mockRepo(mockRepo)
			}

			service := NewSupplierService(mockRepo)

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
