package services

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	mockSale "github.com/WagaoCarvalho/backend_store_go/infra/mock/sale"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
)

func TestSaleService_Create(t *testing.T) {
	mockRepo := new(mockSale.MockSale)
	svc := NewSaleService(mockRepo)
	ctx := context.Background()

	t.Run("should return ErrInvalidData when sale is nil", func(t *testing.T) {
		result, err := svc.Create(ctx, nil)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("should return ErrInvalidData when structural validation fails", func(t *testing.T) {
		// Criando uma venda com dados inválidos para falhar na validação estrutural
		sale := &models.Sale{
			ClientID:           utils.Int64Ptr(1),
			UserID:             utils.Int64Ptr(1),
			SaleDate:           time.Now(),
			TotalItemsAmount:   100.00,
			TotalItemsDiscount: -10.00, // Valor negativo - falha na validação estrutural
			TotalSaleDiscount:  5.00,
			TotalAmount:        95.00,
			PaymentType:        "invalid_payment", // Tipo de pagamento inválido
			Status:             "active",
			Notes:              "Test sale",
			Version:            0, // Version < 1 - falha na validação
		}

		result, err := svc.Create(ctx, sale)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("should return ErrInvalidData when business validation fails", func(t *testing.T) {
		// Criando uma venda que falha nas regras de negócio
		sale := &models.Sale{
			ClientID:           utils.Int64Ptr(1),
			UserID:             utils.Int64Ptr(1),
			SaleDate:           time.Time{}, // Data zero - falha na validação de negócio
			TotalItemsAmount:   100.00,
			TotalItemsDiscount: 50.00,
			TotalSaleDiscount:  60.00, // Total de descontos (110) > TotalAmount (95) - falha
			TotalAmount:        95.00,
			PaymentType:        "cash",
			Status:             "active",
			Notes:              "Test sale",
			Version:            1,
		}

		result, err := svc.Create(ctx, sale)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("should propagate repository error", func(t *testing.T) {
		// Criando uma venda válida
		sale := &models.Sale{
			ClientID:           utils.Int64Ptr(1),
			UserID:             utils.Int64Ptr(1),
			SaleDate:           time.Now(),
			TotalItemsAmount:   100.00,
			TotalItemsDiscount: 10.00,
			TotalSaleDiscount:  5.00,
			TotalAmount:        85.00,
			PaymentType:        "cash",
			Status:             "active",
			Notes:              "Test sale",
			Version:            1,
		}

		// Configurando o mock para retornar um erro
		expectedErr := fmt.Errorf("%w: database error", errMsg.ErrCreate)
		mockRepo.On("Create", ctx, sale).Return(nil, expectedErr).Once()

		result, err := svc.Create(ctx, sale)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, expectedErr, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should return ErrDBInvalidForeignKey when repository returns foreign key error", func(t *testing.T) {
		// Criando uma venda válida
		sale := &models.Sale{
			ClientID:           utils.Int64Ptr(1),
			UserID:             utils.Int64Ptr(1),
			SaleDate:           time.Now(),
			TotalItemsAmount:   100.00,
			TotalItemsDiscount: 10.00,
			TotalSaleDiscount:  5.00,
			TotalAmount:        85.00,
			PaymentType:        "cash",
			Status:             "active",
			Notes:              "Test sale",
			Version:            1,
		}

		// Configurando o mock para retornar erro de chave estrangeira
		expectedErr := errMsg.ErrDBInvalidForeignKey
		mockRepo.On("Create", ctx, sale).Return(nil, expectedErr).Once()

		result, err := svc.Create(ctx, sale)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, expectedErr)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should successfully create sale with valid data", func(t *testing.T) {
		now := time.Now()
		sale := &models.Sale{
			ClientID:           utils.Int64Ptr(1),
			UserID:             utils.Int64Ptr(2),
			SaleDate:           now,
			TotalItemsAmount:   200.00,
			TotalItemsDiscount: 20.00,
			TotalSaleDiscount:  10.00,
			TotalAmount:        170.00, // 200 - 20 - 10 = 170 (válido)
			PaymentType:        "cash",
			Status:             "active",
			Notes:              "Test sale with valid data",
			Version:            1,
		}

		// Criando o objeto que será retornado pelo repositório
		createdSale := &models.Sale{
			ID:                 1,
			ClientID:           sale.ClientID,
			UserID:             sale.UserID,
			SaleDate:           sale.SaleDate,
			TotalItemsAmount:   sale.TotalItemsAmount,
			TotalItemsDiscount: sale.TotalItemsDiscount,
			TotalSaleDiscount:  sale.TotalSaleDiscount,
			TotalAmount:        sale.TotalAmount,
			PaymentType:        sale.PaymentType,
			Status:             sale.Status,
			Notes:              sale.Notes,
			Version:            sale.Version,
			CreatedAt:          now,
			UpdatedAt:          now,
		}

		// Configurando o mock para retornar sucesso
		mockRepo.On("Create", ctx, sale).Return(createdSale, nil).Once()

		result, err := svc.Create(ctx, sale)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, createdSale.ID, result.ID)
		assert.Equal(t, createdSale.TotalAmount, result.TotalAmount)
		assert.Equal(t, createdSale.Status, result.Status)
		assert.Equal(t, createdSale.PaymentType, result.PaymentType)
		assert.Equal(t, createdSale.Version, result.Version)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should successfully create sale with different payment types", func(t *testing.T) {
		testCases := []struct {
			name        string
			paymentType string
		}{
			{"cash payment", "cash"},
			{"card payment", "card"},
			{"credit payment", "credit"},
			{"pix payment", "pix"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				sale := &models.Sale{
					ClientID:           utils.Int64Ptr(1),
					UserID:             utils.Int64Ptr(1),
					SaleDate:           time.Now(),
					TotalItemsAmount:   100.00,
					TotalItemsDiscount: 10.00,
					TotalSaleDiscount:  5.00,
					TotalAmount:        85.00,
					PaymentType:        tc.paymentType,
					Status:             "active",
					Notes:              "Test sale",
					Version:            1,
				}

				createdSale := &models.Sale{
					ID:       2,
					ClientID: sale.ClientID,
					UserID:   sale.UserID,
					// ... outros campos copiados
				}

				mockRepo.On("Create", ctx, sale).Return(createdSale, nil).Once()

				result, err := svc.Create(ctx, sale)

				assert.NoError(t, err)
				assert.NotNil(t, result)

				mockRepo.AssertExpectations(t)
			})
		}
	})

	t.Run("should successfully create sale with different statuses", func(t *testing.T) {
		testCases := []struct {
			name   string
			status string
		}{
			{"active status", "active"},
			{"canceled status", "canceled"},
			{"returned status", "returned"},
			{"completed status", "completed"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				sale := &models.Sale{
					ClientID:           utils.Int64Ptr(1),
					UserID:             utils.Int64Ptr(1),
					SaleDate:           time.Now(),
					TotalItemsAmount:   100.00,
					TotalItemsDiscount: 10.00,
					TotalSaleDiscount:  5.00,
					TotalAmount:        85.00,
					PaymentType:        "cash",
					Status:             tc.status,
					Notes:              "Test sale",
					Version:            1,
				}

				createdSale := &models.Sale{
					ID:       3,
					ClientID: sale.ClientID,
					UserID:   sale.UserID,
					// ... outros campos copiados
				}

				mockRepo.On("Create", ctx, sale).Return(createdSale, nil).Once()

				result, err := svc.Create(ctx, sale)

				assert.NoError(t, err)
				assert.NotNil(t, result)

				mockRepo.AssertExpectations(t)
			})
		}
	})

	t.Run("should handle sale with empty notes", func(t *testing.T) {
		sale := &models.Sale{
			ClientID:           utils.Int64Ptr(1),
			UserID:             utils.Int64Ptr(1),
			SaleDate:           time.Now(),
			TotalItemsAmount:   100.00,
			TotalItemsDiscount: 10.00,
			TotalSaleDiscount:  5.00,
			TotalAmount:        85.00,
			PaymentType:        "cash",
			Status:             "active",
			Notes:              "", // Notes vazio é permitido
			Version:            1,
		}

		createdSale := &models.Sale{
			ID:       4,
			ClientID: sale.ClientID,
			UserID:   sale.UserID,
			// ... outros campos copiados
		}

		mockRepo.On("Create", ctx, sale).Return(createdSale, nil).Once()

		result, err := svc.Create(ctx, sale)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should handle sale with long notes (up to 500 characters)", func(t *testing.T) {
		longNotes := strings.Repeat("a", 500) // Máximo permitido

		sale := &models.Sale{
			ClientID:           utils.Int64Ptr(1),
			UserID:             utils.Int64Ptr(1),
			SaleDate:           time.Now(),
			TotalItemsAmount:   100.00,
			TotalItemsDiscount: 10.00,
			TotalSaleDiscount:  5.00,
			TotalAmount:        85.00,
			PaymentType:        "cash",
			Status:             "active",
			Notes:              longNotes,
			Version:            1,
		}

		createdSale := &models.Sale{
			ID:       5,
			ClientID: sale.ClientID,
			UserID:   sale.UserID,
			// ... outros campos copiados
		}

		mockRepo.On("Create", ctx, sale).Return(createdSale, nil).Once()

		result, err := svc.Create(ctx, sale)

		assert.NoError(t, err)
		assert.NotNil(t, result)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should fail when notes exceed 500 characters", func(t *testing.T) {
		longNotes := strings.Repeat("a", 501) // Excede o máximo permitido

		sale := &models.Sale{
			ClientID:           utils.Int64Ptr(1),
			UserID:             utils.Int64Ptr(1),
			SaleDate:           time.Now(),
			TotalItemsAmount:   100.00,
			TotalItemsDiscount: 10.00,
			TotalSaleDiscount:  5.00,
			TotalAmount:        85.00,
			PaymentType:        "cash",
			Status:             "active",
			Notes:              longNotes,
			Version:            1,
		}

		result, err := svc.Create(ctx, sale)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})
}

func TestSaleService_Update(t *testing.T) {
	mockRepo := new(mockSale.MockSale)
	svc := NewSaleService(mockRepo)
	ctx := context.Background()

	t.Run("sale nil", func(t *testing.T) {
		err := svc.Update(ctx, nil)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("id zero", func(t *testing.T) {
		s := &models.Sale{
			ID:                0,
			Version:           1,
			TotalAmount:       100,
			TotalSaleDiscount: 10,
			PaymentType:       "cash",
			Status:            "active",
			SaleDate:          time.Now(),
		}
		err := svc.Update(ctx, s)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("version zero", func(t *testing.T) {
		s := &models.Sale{
			ID:                1,
			Version:           0,
			TotalAmount:       100,
			TotalSaleDiscount: 10,
			PaymentType:       "cash",
			Status:            "active",
			SaleDate:          time.Now(),
		}
		err := svc.Update(ctx, s)
		assert.ErrorIs(t, err, errMsg.ErrVersionConflict)
	})

	t.Run("structural validation fails", func(t *testing.T) {
		s := &models.Sale{
			ID:                1,
			Version:           1,
			TotalAmount:       -1,
			TotalSaleDiscount: -1,
			PaymentType:       "",
			Status:            "",
		}
		err := svc.Update(ctx, s)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("business validation fails", func(t *testing.T) {
		s := &models.Sale{
			ID:                1,
			Version:           1,
			TotalAmount:       10,
			TotalSaleDiscount: 20, // maior que o total
			PaymentType:       "cash",
			Status:            "active",
			SaleDate:          time.Now(),
		}
		err := svc.Update(ctx, s)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("repo returns error", func(t *testing.T) {
		s := &models.Sale{
			ID:                1,
			Version:           1,
			TotalAmount:       100,
			TotalSaleDiscount: 10,
			PaymentType:       "cash",
			Status:            "active",
			SaleDate:          time.Now(),
		}
		mockRepo.On("Update", ctx, s).Return(errors.New("repo error")).Once()

		err := svc.Update(ctx, s)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errMsg.ErrUpdate.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		s := &models.Sale{
			ID:                1,
			Version:           1,
			TotalAmount:       100,
			TotalSaleDiscount: 10,
			PaymentType:       "cash",
			Status:            "active",
			SaleDate:          time.Now(),
			Notes:             "Teste",
		}
		mockRepo.On("Update", ctx, s).Return(nil).Once()

		err := svc.Update(ctx, s)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestSaleService_Delete(t *testing.T) {
	mockRepo := new(mockSale.MockSale)
	svc := NewSaleService(mockRepo)
	ctx := context.Background()

	t.Run("id zero", func(t *testing.T) {
		err := svc.Delete(ctx, 0)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("repo returns error", func(t *testing.T) {
		mockRepo.On("Delete", ctx, int64(1)).Return(errors.New("repo error")).Once()
		err := svc.Delete(ctx, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errMsg.ErrDelete.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		mockRepo.On("Delete", ctx, int64(2)).Return(nil).Once()
		err := svc.Delete(ctx, 2)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
