package services

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	modelsAddress "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	modelsContact "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	modelsSupplier "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	modelsSupplierCategories "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_category"
	modelsCatRel "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_category_relation"
	modelsSupplierFull "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_full"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"

	mockTX "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo"
	mockAddress "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/address"
	mockContact "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/contact"
	mockCatRel "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/supplier"
	mockSupplier "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/supplier"
)

func TestSupplierFullService_CreateFull(t *testing.T) {
	ctx := context.Background()

	// ------------------------
	// Grupo: Supplier
	// ------------------------
	t.Run("Supplier: falha quando supplierFull é nil", func(t *testing.T) {
		service := NewSupplierFull(
			new(mockSupplier.MockSupplierFullRepository),
			nil, nil, nil, nil,
		)

		result, err := service.CreateFull(ctx, nil)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errMsg.ErrInvalidData.Error())
	})

	t.Run("Supplier: falha ao validar supplier inválido", func(t *testing.T) {
		service := NewSupplierFull(
			new(mockSupplier.MockSupplierFullRepository),
			nil, nil, nil, nil,
		)

		supplierFull := &modelsSupplierFull.SupplierFull{
			Supplier: &modelsSupplier.Supplier{Name: ""},
		}

		result, err := service.CreateFull(ctx, supplierFull)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), errMsg.ErrInvalidData.Error())
	})

	t.Run("Supplier: sucesso ao criar supplier completo", func(t *testing.T) {
		mockRepoSupplier := new(mockSupplier.MockSupplierFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockRepoContact := new(mockContact.MockContactRepository)
		mockRepoCatRel := new(mockCatRel.MockSupplierCategoryRelationRepo)
		tx := new(mockTX.MockTx)

		service := NewSupplierFull(
			mockRepoSupplier,
			mockRepoAddress,
			mockRepoContact,
			mockRepoCatRel,
			nil,
		)

		supplierFull := &modelsSupplierFull.SupplierFull{
			Supplier: &modelsSupplier.Supplier{
				ID:   1,
				Name: "Fornecedor Teste",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "123",
				City:         "Cidade Teste",
				State:        "SP",
				PostalCode:   "01234567",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@teste.com",
				Phone:       "1112345678",
			},
			Categories: []modelsSupplierCategories.SupplierCategory{
				{ID: 100},
			},
		}

		mockRepoSupplier.On("BeginTx", ctx).Return(tx, nil)
		mockRepoSupplier.On("CreateTx", ctx, tx, supplierFull.Supplier).Return(supplierFull.Supplier, nil)
		mockRepoAddress.On("CreateTx", ctx, tx, supplierFull.Address).Return(supplierFull.Address, nil)
		mockRepoContact.On("CreateTx", ctx, tx, supplierFull.Contact).Return(supplierFull.Contact, nil)
		mockRepoCatRel.On("CreateTx", ctx, tx, mock.AnythingOfType("*model.SupplierCategoryRelation")).Return(&modelsCatRel.SupplierCategoryRelation{
			SupplierID: supplierFull.Supplier.ID,
			CategoryID: 100,
			CreatedAt:  time.Now(),
		}, nil)

		// Aqui adiciona as expectativas do TX
		tx.On("Commit", ctx).Return(nil)
		tx.On("Rollback", ctx).Return(nil)

		result, err := service.CreateFull(ctx, supplierFull)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, supplierFull.Supplier.Name, result.Supplier.Name)
		assert.Equal(t, supplierFull.Address.Street, result.Address.Street)
		assert.Equal(t, supplierFull.Contact.Email, result.Contact.Email)
		assert.Len(t, result.Categories, 1)

		// Garante que rollback e commit foram chamados corretamente
		tx.AssertCalled(t, "Commit", ctx)
		tx.AssertNotCalled(t, "Rollback", ctx)
	})

	// ------------------------
	// Grupo: Transaction
	// ------------------------

	t.Run("Transaction: falha ao iniciar transação", func(t *testing.T) {
		mockRepoSupplier := new(mockSupplier.MockSupplierFullRepository)
		service := NewSupplierFull(
			mockRepoSupplier,
			new(mockAddress.MockAddressRepository),
			new(mockContact.MockContactRepository),
			new(mockCatRel.MockSupplierCategoryRelationRepo),
			nil,
		)

		supplierFull := &modelsSupplierFull.SupplierFull{
			Supplier: &modelsSupplier.Supplier{
				ID:   1,
				Name: "Fornecedor Teste",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				SupplierID:   utils.Int64Ptr(1),
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelsSupplierCategories.SupplierCategory{
				{ID: 100},
			},
		}

		mockRepoSupplier.On("BeginTx", ctx).Return(nil, errors.New("begin error"))

		result, err := service.CreateFull(ctx, supplierFull)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao iniciar transação")
	})

	t.Run("Transaction: falha quando transação é nil", func(t *testing.T) {
		mockRepoSupplier := new(mockSupplier.MockSupplierFullRepository)
		service := NewSupplierFull(
			mockRepoSupplier,
			new(mockAddress.MockAddressRepository),
			new(mockContact.MockContactRepository),
			new(mockCatRel.MockSupplierCategoryRelationRepo),
			nil,
		)

		supplierFull := &modelsSupplierFull.SupplierFull{
			Supplier: &modelsSupplier.Supplier{
				ID: 1, Name: "Fornecedor Teste",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				SupplierID:   utils.Int64Ptr(1),
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelsSupplierCategories.SupplierCategory{
				{ID: 100},
			},
		}

		mockRepoSupplier.On("BeginTx", ctx).Return(nil, nil)

		result, err := service.CreateFull(ctx, supplierFull)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "transação inválida")
	})

	t.Run("Transaction: falha no rollback quando há erro", func(t *testing.T) {
		mockRepoSupplier := new(mockSupplier.MockSupplierFullRepository)
		tx := new(mockTX.MockTx)

		service := NewSupplierFull(
			mockRepoSupplier,
			new(mockAddress.MockAddressRepository),
			new(mockContact.MockContactRepository),
			new(mockCatRel.MockSupplierCategoryRelationRepo),
			nil,
		)

		supplierFull := &modelsSupplierFull.SupplierFull{
			Supplier: &modelsSupplier.Supplier{
				ID: 1, Name: "Fornecedor Teste",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				SupplierID:   utils.Int64Ptr(1),
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelsSupplierCategories.SupplierCategory{
				{ID: 100},
			},
		}

		mockRepoSupplier.On("BeginTx", ctx).Return(tx, nil)
		mockRepoSupplier.On("CreateTx", ctx, tx, supplierFull.Supplier).Return(nil, errors.New("erro ao criar fornecedor"))
		tx.On("Rollback", ctx).Return(errors.New("erro no rollback"))

		result, err := service.CreateFull(ctx, supplierFull)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar fornecedor")
		assert.Contains(t, err.Error(), "rollback error")
	})

	t.Run("Transaction: falha no commit e rollback também falha", func(t *testing.T) {
		mockRepoSupplier := new(mockSupplier.MockSupplierFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockRepoContact := new(mockContact.MockContactRepository)
		mockRepoCatRel := new(mockCatRel.MockSupplierCategoryRelationRepo)
		tx := new(mockTX.MockTx)

		service := NewSupplierFull(
			mockRepoSupplier,
			mockRepoAddress,
			mockRepoContact,
			mockRepoCatRel,
			nil,
		)

		supplierFull := &modelsSupplierFull.SupplierFull{
			Supplier: &modelsSupplier.Supplier{
				ID: 1, Name: "Fornecedor Teste",
			},
			Address: &modelsAddress.Address{
				SupplierID:   utils.Int64Ptr(1),
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelsSupplierCategories.SupplierCategory{
				{ID: 100},
			},
		}

		mockRepoSupplier.On("BeginTx", ctx).Return(tx, nil)
		mockRepoSupplier.On("CreateTx", ctx, tx, supplierFull.Supplier).Return(supplierFull.Supplier, nil)
		mockRepoAddress.On("CreateTx", ctx, tx, supplierFull.Address).Return(supplierFull.Address, nil)
		mockRepoContact.On("CreateTx", ctx, tx, supplierFull.Contact).Return(supplierFull.Contact, nil)
		mockRepoCatRel.On("CreateTx", ctx, tx, mock.Anything).Return(&modelsCatRel.SupplierCategoryRelation{
			SupplierID: supplierFull.Supplier.ID,
			CategoryID: 100,
			CreatedAt:  time.Now(),
		}, nil)

		commitErr := errors.New("erro no commit")
		rollbackErr := errors.New("erro no rollback")
		tx.On("Commit", ctx).Return(commitErr)
		tx.On("Rollback", ctx).Return(rollbackErr)

		result, err := service.CreateFull(ctx, supplierFull)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao commitar transação")
		assert.Contains(t, err.Error(), commitErr.Error())
		assert.Contains(t, err.Error(), rollbackErr.Error())
	})

	t.Run("Transaction: rollback é chamado e panic é propagado (SupplierFullService)", func(t *testing.T) {
		mockRepoSupplier := new(mockSupplier.MockSupplierFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockRepoContact := new(mockContact.MockContactRepository)
		mockRepoCatRel := new(mockCatRel.MockSupplierCategoryRelationRepo)
		tx := new(mockTX.MockTx)

		service := NewSupplierFull(
			mockRepoSupplier,
			mockRepoAddress,
			mockRepoContact,
			mockRepoCatRel,
			nil,
		)

		supplierFull := &modelsSupplierFull.SupplierFull{
			Supplier: &modelsSupplier.Supplier{
				ID:   1,
				Name: "Fornecedor Teste",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "123",
				City:         "Cidade Teste",
				State:        "SP",
				PostalCode:   "01234567",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@teste.com",
				Phone:       "11999999999",
			},
			Categories: []modelsSupplierCategories.SupplierCategory{
				{ID: 100},
			},
		}

		mockRepoSupplier.On("BeginTx", ctx).Return(tx, nil)

		// Simula panic durante a criação do supplier
		mockRepoSupplier.On("CreateTx", ctx, tx, supplierFull.Supplier).Run(func(args mock.Arguments) {
			panic("panic simulado durante criação do supplier")
		}).Return(nil, nil)

		// Rollback deve ser chamado no defer
		tx.On("Rollback", ctx).Return(nil)

		defer func() {
			if r := recover(); r != nil {
				tx.AssertCalled(t, "Rollback", ctx)
				assert.Equal(t, "panic simulado durante criação do supplier", r)
			} else {
				t.Errorf("Esperado panic, mas não ocorreu")
			}
		}()

		// Chamada que dispara o panic
		_, _ = service.CreateFull(ctx, supplierFull)
	})

	t.Run("Transaction: commit falha mas rollback funciona", func(t *testing.T) {
		mockRepoSupplier := new(mockSupplier.MockSupplierFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockRepoContact := new(mockContact.MockContactRepository)
		mockRepoCatRel := new(mockCatRel.MockSupplierCategoryRelationRepo)
		tx := new(mockTX.MockTx)

		service := NewSupplierFull(
			mockRepoSupplier,
			mockRepoAddress,
			mockRepoContact,
			mockRepoCatRel,
			nil,
		)

		supplierFull := &modelsSupplierFull.SupplierFull{
			Supplier: &modelsSupplier.Supplier{
				ID:   1,
				Name: "Fornecedor Teste",
			},
			Address: &modelsAddress.Address{
				SupplierID:   utils.Int64Ptr(1),
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelsSupplierCategories.SupplierCategory{
				{ID: 100},
			},
		}

		// Configura os mocks de criação
		mockRepoSupplier.On("BeginTx", ctx).Return(tx, nil)
		mockRepoSupplier.On("CreateTx", ctx, tx, supplierFull.Supplier).Return(supplierFull.Supplier, nil)
		mockRepoAddress.On("CreateTx", ctx, tx, supplierFull.Address).Return(supplierFull.Address, nil)
		mockRepoContact.On("CreateTx", ctx, tx, supplierFull.Contact).Return(supplierFull.Contact, nil)

		// Mock para a relação de categoria
		mockRepoCatRel.On(
			"CreateTx",
			ctx,
			tx,
			mock.AnythingOfType("*model.SupplierCategoryRelation"),
		).Return(&modelsCatRel.SupplierCategoryRelation{
			SupplierID: supplierFull.Supplier.ID,
			CategoryID: 100,
			CreatedAt:  time.Now(),
		}, nil)

		// Cenário: commit falha MAS rollback funciona
		commitErr := errors.New("erro no commit")
		tx.On("Commit", ctx).Return(commitErr)
		tx.On("Rollback", ctx).Return(nil) // Rollback funciona

		// Executa o serviço
		result, err := service.CreateFull(ctx, supplierFull)

		// Valida resultados
		assert.Nil(t, result)
		assert.Error(t, err)

		// Deve conter a mensagem de erro do commit
		assert.Contains(t, err.Error(), "erro ao commitar transação")
		assert.Contains(t, err.Error(), commitErr.Error())

		// NÃO deve conter "rollback error" pois o rollback foi bem-sucedido
		assert.NotContains(t, err.Error(), "rollback error")

		// Verifica que o erro é exatamente o esperado: "erro ao commitar transação: {commitErr}"
		expectedError := fmt.Errorf("erro ao commitar transação: %w", commitErr)
		assert.Equal(t, expectedError.Error(), err.Error())
	})

	t.Run("Transaction: erro na operação mas rollback funciona", func(t *testing.T) {
		mockRepoSupplier := new(mockSupplier.MockSupplierFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockRepoContact := new(mockContact.MockContactRepository)
		mockRepoCatRel := new(mockCatRel.MockSupplierCategoryRelationRepo)
		tx := new(mockTX.MockTx)

		service := NewSupplierFull(
			mockRepoSupplier,
			mockRepoAddress,
			mockRepoContact,
			mockRepoCatRel,
			nil,
		)

		supplierFull := &modelsSupplierFull.SupplierFull{
			Supplier: &modelsSupplier.Supplier{
				ID:   1,
				Name: "Fornecedor Teste",
			},
			Address: &modelsAddress.Address{
				SupplierID:   utils.Int64Ptr(1),
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelsSupplierCategories.SupplierCategory{
				{ID: 100},
			},
		}

		// Configura os mocks
		mockRepoSupplier.On("BeginTx", ctx).Return(tx, nil)

		// Simula erro na criação do supplier
		operationErr := errors.New("erro ao criar fornecedor")
		mockRepoSupplier.On("CreateTx", ctx, tx, supplierFull.Supplier).Return(nil, operationErr)

		// Rollback funciona
		tx.On("Rollback", ctx).Return(nil)

		// Executa o serviço
		result, err := service.CreateFull(ctx, supplierFull)

		// Valida resultados
		assert.Nil(t, result)
		assert.Error(t, err)

		// Deve retornar exatamente o mesmo erro da operação (sem modificação)
		assert.Equal(t, operationErr, err)
		assert.Contains(t, err.Error(), "erro ao criar fornecedor")

		// Não deve conter "rollback error" pois o rollback foi bem-sucedido
		assert.NotContains(t, err.Error(), "rollback error")
	})

	// Address
	t.Run("Address: falha na validação após setar SupplierID", func(t *testing.T) {
		mockRepoSupplier := new(mockSupplier.MockSupplierFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		tx := new(mockTX.MockTx)

		service := NewSupplierFull(
			mockRepoSupplier,
			mockRepoAddress,
			nil, // não precisa do contact para este teste
			nil, // não precisa do category relation para este teste
			nil,
		)

		supplierFull := &modelsSupplierFull.SupplierFull{
			Supplier: &modelsSupplier.Supplier{
				ID:   1,
				Name: "Fornecedor Teste",
			},
			Address: &modelsAddress.Address{
				// Address sem campos obrigatórios para forçar erro na validação
				ID:           10,
				Street:       "", // Campo obrigatório vazio
				StreetNumber: "45",
				City:         "", // Campo obrigatório vazio
				State:        "SP",
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelsSupplierCategories.SupplierCategory{
				{ID: 100},
			},
		}

		// Configura os mocks
		mockRepoSupplier.On("BeginTx", ctx).Return(tx, nil)

		// Criação do supplier funciona
		createdSupplier := &modelsSupplier.Supplier{
			ID:   1,
			Name: "Fornecedor Teste",
		}
		mockRepoSupplier.On("CreateTx", ctx, tx, supplierFull.Supplier).Return(createdSupplier, nil)

		// Rollback funciona
		tx.On("Rollback", ctx).Return(nil)

		// Executa o serviço
		result, err := service.CreateFull(ctx, supplierFull)

		// Valida resultados
		assert.Nil(t, result)
		assert.Error(t, err)

		// Deve conter a mensagem de erro específica do endereço
		assert.Contains(t, err.Error(), "endereço inválido")

		// CORREÇÃO: Verificar a mensagem real baseada no erro que está aparecendo
		assert.Contains(t, err.Error(), "street")
		assert.Contains(t, err.Error(), "city")
		assert.Contains(t, err.Error(), "campo obrigatório")

		// Verifica que o SupplierID foi setado antes da validação
		// (isso é verificado indiretamente pelo fato de que a validação foi executada)
	})

	t.Run("Address: falha ao criar endereço e rollback também falha", func(t *testing.T) {
		mockRepoSupplier := new(mockSupplier.MockSupplierFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		tx := new(mockTX.MockTx)

		service := NewSupplierFull(
			mockRepoSupplier,
			mockRepoAddress,
			nil,
			nil,
			nil,
		)

		supplierFull := &modelsSupplierFull.SupplierFull{
			Supplier: &modelsSupplier.Supplier{
				ID:   1,
				Name: "Fornecedor Teste",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelsSupplierCategories.SupplierCategory{
				{ID: 100},
			},
		}

		mockRepoSupplier.On("BeginTx", ctx).Return(tx, nil)
		createdSupplier := &modelsSupplier.Supplier{ID: 1, Name: "Fornecedor Teste"}
		mockRepoSupplier.On("CreateTx", ctx, tx, supplierFull.Supplier).Return(createdSupplier, nil)

		// Criação do address falha
		addressErr := errors.New("erro ao criar endereço no banco")
		mockRepoAddress.On("CreateTx", ctx, tx, supplierFull.Address).Return(nil, addressErr)

		// Rollback também falha
		rollbackErr := errors.New("erro no rollback")
		tx.On("Rollback", ctx).Return(rollbackErr)

		result, err := service.CreateFull(ctx, supplierFull)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar endereço no banco")
		assert.Contains(t, err.Error(), "rollback error")
		assert.Contains(t, err.Error(), rollbackErr.Error())
	})

	// Contato

	t.Run("Contact: falha na validação do contato", func(t *testing.T) {
		mockRepoSupplier := new(mockSupplier.MockSupplierFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockRepoContact := new(mockContact.MockContactRepository)
		tx := new(mockTX.MockTx)

		service := NewSupplierFull(
			mockRepoSupplier,
			mockRepoAddress,
			mockRepoContact,
			nil,
			nil,
		)

		supplierFull := &modelsSupplierFull.SupplierFull{
			Supplier: &modelsSupplier.Supplier{
				ID:   1,
				Name: "Fornecedor Teste",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "",               // Campo obrigatório vazio
				Email:       "email-invalido", // Email inválido
				Phone:       "123",            // Phone inválido
			},
			Categories: []modelsSupplierCategories.SupplierCategory{
				{ID: 100},
			},
		}

		mockRepoSupplier.On("BeginTx", ctx).Return(tx, nil)
		createdSupplier := &modelsSupplier.Supplier{ID: 1, Name: "Fornecedor Teste"}
		mockRepoSupplier.On("CreateTx", ctx, tx, supplierFull.Supplier).Return(createdSupplier, nil)

		createdAddress := &modelsAddress.Address{
			ID:           10,
			Street:       "Rua Teste",
			StreetNumber: "45",
			City:         "Cidade Teste",
			State:        "SP",
			PostalCode:   "03459808",
			IsActive:     true,
			Country:      "Brasil",
			SupplierID:   utils.Int64Ptr(1),
		}
		mockRepoAddress.On("CreateTx", ctx, tx, supplierFull.Address).Return(createdAddress, nil)

		tx.On("Rollback", ctx).Return(nil)

		result, err := service.CreateFull(ctx, supplierFull)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "contato inválido")
	})

	t.Run("Contact: falha ao criar contato no banco", func(t *testing.T) {
		mockRepoSupplier := new(mockSupplier.MockSupplierFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockRepoContact := new(mockContact.MockContactRepository)
		tx := new(mockTX.MockTx)

		service := NewSupplierFull(
			mockRepoSupplier,
			mockRepoAddress,
			mockRepoContact,
			nil,
			nil,
		)

		supplierFull := &modelsSupplierFull.SupplierFull{
			Supplier: &modelsSupplier.Supplier{
				ID:   1,
				Name: "Fornecedor Teste",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelsSupplierCategories.SupplierCategory{
				{ID: 100},
			},
		}

		mockRepoSupplier.On("BeginTx", ctx).Return(tx, nil)
		createdSupplier := &modelsSupplier.Supplier{ID: 1, Name: "Fornecedor Teste"}
		mockRepoSupplier.On("CreateTx", ctx, tx, supplierFull.Supplier).Return(createdSupplier, nil)

		createdAddress := &modelsAddress.Address{
			ID:           10,
			Street:       "Rua Teste",
			StreetNumber: "45",
			City:         "Cidade Teste",
			State:        "SP",
			PostalCode:   "03459808",
			IsActive:     true,
			Country:      "Brasil",
			SupplierID:   utils.Int64Ptr(1),
		}
		mockRepoAddress.On("CreateTx", ctx, tx, supplierFull.Address).Return(createdAddress, nil)

		contactErr := errors.New("erro ao criar contato no banco")
		mockRepoContact.On("CreateTx", ctx, tx, supplierFull.Contact).Return(nil, contactErr)

		tx.On("Rollback", ctx).Return(nil)

		result, err := service.CreateFull(ctx, supplierFull)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, contactErr, err)
		assert.Contains(t, err.Error(), "erro ao criar contato no banco")
	})

	t.Run("Contact: falha ao criar contato e rollback também falha", func(t *testing.T) {
		mockRepoSupplier := new(mockSupplier.MockSupplierFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockRepoContact := new(mockContact.MockContactRepository)
		tx := new(mockTX.MockTx)

		service := NewSupplierFull(
			mockRepoSupplier,
			mockRepoAddress,
			mockRepoContact,
			nil,
			nil,
		)

		supplierFull := &modelsSupplierFull.SupplierFull{
			Supplier: &modelsSupplier.Supplier{
				ID:   1,
				Name: "Fornecedor Teste",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelsSupplierCategories.SupplierCategory{
				{ID: 100},
			},
		}

		mockRepoSupplier.On("BeginTx", ctx).Return(tx, nil)
		createdSupplier := &modelsSupplier.Supplier{ID: 1, Name: "Fornecedor Teste"}
		mockRepoSupplier.On("CreateTx", ctx, tx, supplierFull.Supplier).Return(createdSupplier, nil)

		createdAddress := &modelsAddress.Address{
			ID:           10,
			Street:       "Rua Teste",
			StreetNumber: "45",
			City:         "Cidade Teste",
			State:        "SP",
			PostalCode:   "03459808",
			IsActive:     true,
			Country:      "Brasil",
			SupplierID:   utils.Int64Ptr(1),
		}
		mockRepoAddress.On("CreateTx", ctx, tx, supplierFull.Address).Return(createdAddress, nil)

		contactErr := errors.New("erro ao criar contato no banco")
		mockRepoContact.On("CreateTx", ctx, tx, supplierFull.Contact).Return(nil, contactErr)

		rollbackErr := errors.New("erro no rollback")
		tx.On("Rollback", ctx).Return(rollbackErr)

		result, err := service.CreateFull(ctx, supplierFull)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar contato no banco")
		assert.Contains(t, err.Error(), "rollback error")
		assert.Contains(t, err.Error(), rollbackErr.Error())
	})

	t.Run("Contact: sucesso ao criar contato", func(t *testing.T) {
		mockRepoSupplier := new(mockSupplier.MockSupplierFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockRepoContact := new(mockContact.MockContactRepository)
		mockRepoCatRel := new(mockCatRel.MockSupplierCategoryRelationRepo)
		tx := new(mockTX.MockTx)

		service := NewSupplierFull(
			mockRepoSupplier,
			mockRepoAddress,
			mockRepoContact,
			mockRepoCatRel,
			nil,
		)

		supplierFull := &modelsSupplierFull.SupplierFull{
			Supplier: &modelsSupplier.Supplier{
				ID:   1,
				Name: "Fornecedor Teste",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelsSupplierCategories.SupplierCategory{
				{ID: 100},
			},
		}

		mockRepoSupplier.On("BeginTx", ctx).Return(tx, nil)
		createdSupplier := &modelsSupplier.Supplier{ID: 1, Name: "Fornecedor Teste"}
		mockRepoSupplier.On("CreateTx", ctx, tx, supplierFull.Supplier).Return(createdSupplier, nil)

		createdAddress := &modelsAddress.Address{
			ID:           10,
			Street:       "Rua Teste",
			StreetNumber: "45",
			City:         "Cidade Teste",
			State:        "SP",
			PostalCode:   "03459808",
			IsActive:     true,
			Country:      "Brasil",
			SupplierID:   utils.Int64Ptr(1),
		}
		mockRepoAddress.On("CreateTx", ctx, tx, supplierFull.Address).Return(createdAddress, nil)

		createdContact := &modelsContact.Contact{
			ID:          20,
			ContactName: "Contato Teste",
			Email:       "contato@example.com",
			Phone:       "1234567895",
		}
		mockRepoContact.On("CreateTx", ctx, tx, supplierFull.Contact).Return(createdContact, nil)

		mockRepoCatRel.On("CreateTx", ctx, tx, mock.AnythingOfType("*model.SupplierCategoryRelation")).
			Return(&modelsCatRel.SupplierCategoryRelation{
				SupplierID: 1,
				CategoryID: 100,
				CreatedAt:  time.Now(),
			}, nil)

		tx.On("Commit", ctx).Return(nil)

		result, err := service.CreateFull(ctx, supplierFull)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, createdContact, result.Contact)
	})

	// Relação Categoria

	t.Run("CategoryRelation: falha na validação da relação", func(t *testing.T) {
		mockRepoSupplier := new(mockSupplier.MockSupplierFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockRepoContact := new(mockContact.MockContactRepository)
		mockRepoCatRel := new(mockCatRel.MockSupplierCategoryRelationRepo)
		tx := new(mockTX.MockTx)

		service := NewSupplierFull(
			mockRepoSupplier,
			mockRepoAddress,
			mockRepoContact,
			mockRepoCatRel,
			nil,
		)

		supplierFull := &modelsSupplierFull.SupplierFull{
			Supplier: &modelsSupplier.Supplier{
				ID:   1,
				Name: "Fornecedor Teste",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelsSupplierCategories.SupplierCategory{
				{ID: 0}, // ID inválido para forçar erro na validação
			},
		}

		mockRepoSupplier.On("BeginTx", ctx).Return(tx, nil)
		createdSupplier := &modelsSupplier.Supplier{ID: 1, Name: "Fornecedor Teste"}
		mockRepoSupplier.On("CreateTx", ctx, tx, supplierFull.Supplier).Return(createdSupplier, nil)

		createdAddress := &modelsAddress.Address{
			ID:           10,
			Street:       "Rua Teste",
			StreetNumber: "45",
			City:         "Cidade Teste",
			State:        "SP",
			PostalCode:   "03459808",
			IsActive:     true,
			Country:      "Brasil",
			SupplierID:   utils.Int64Ptr(1),
		}
		mockRepoAddress.On("CreateTx", ctx, tx, supplierFull.Address).Return(createdAddress, nil)

		createdContact := &modelsContact.Contact{
			ID:          20,
			ContactName: "Contato Teste",
			Email:       "contato@example.com",
			Phone:       "1234567895",
		}
		mockRepoContact.On("CreateTx", ctx, tx, supplierFull.Contact).Return(createdContact, nil)

		tx.On("Rollback", ctx).Return(nil)

		result, err := service.CreateFull(ctx, supplierFull)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "relação fornecedor-categoria inválida")
	})

	t.Run("CategoryRelation: falha ao criar relação no banco", func(t *testing.T) {
		mockRepoSupplier := new(mockSupplier.MockSupplierFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockRepoContact := new(mockContact.MockContactRepository)
		mockRepoCatRel := new(mockCatRel.MockSupplierCategoryRelationRepo)
		tx := new(mockTX.MockTx)

		service := NewSupplierFull(
			mockRepoSupplier,
			mockRepoAddress,
			mockRepoContact,
			mockRepoCatRel,
			nil,
		)

		supplierFull := &modelsSupplierFull.SupplierFull{
			Supplier: &modelsSupplier.Supplier{
				ID:   1,
				Name: "Fornecedor Teste",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelsSupplierCategories.SupplierCategory{
				{ID: 100},
			},
		}

		mockRepoSupplier.On("BeginTx", ctx).Return(tx, nil)
		createdSupplier := &modelsSupplier.Supplier{ID: 1, Name: "Fornecedor Teste"}
		mockRepoSupplier.On("CreateTx", ctx, tx, supplierFull.Supplier).Return(createdSupplier, nil)

		createdAddress := &modelsAddress.Address{
			ID:           10,
			Street:       "Rua Teste",
			StreetNumber: "45",
			City:         "Cidade Teste",
			State:        "SP",
			PostalCode:   "03459808",
			IsActive:     true,
			Country:      "Brasil",
			SupplierID:   utils.Int64Ptr(1),
		}
		mockRepoAddress.On("CreateTx", ctx, tx, supplierFull.Address).Return(createdAddress, nil)

		createdContact := &modelsContact.Contact{
			ID:          20,
			ContactName: "Contato Teste",
			Email:       "contato@example.com",
			Phone:       "1234567895",
		}
		mockRepoContact.On("CreateTx", ctx, tx, supplierFull.Contact).Return(createdContact, nil)

		relationErr := errors.New("erro ao criar relação no banco")
		mockRepoCatRel.On("CreateTx", ctx, tx, mock.AnythingOfType("*model.SupplierCategoryRelation")).
			Return(nil, relationErr)

		tx.On("Rollback", ctx).Return(nil)

		result, err := service.CreateFull(ctx, supplierFull)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, relationErr, err)
		assert.Contains(t, err.Error(), "erro ao criar relação no banco")
	})

	t.Run("CategoryRelation: falha ao criar relação e rollback também falha", func(t *testing.T) {
		mockRepoSupplier := new(mockSupplier.MockSupplierFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockRepoContact := new(mockContact.MockContactRepository)
		mockRepoCatRel := new(mockCatRel.MockSupplierCategoryRelationRepo)
		tx := new(mockTX.MockTx)

		service := NewSupplierFull(
			mockRepoSupplier,
			mockRepoAddress,
			mockRepoContact,
			mockRepoCatRel,
			nil,
		)

		supplierFull := &modelsSupplierFull.SupplierFull{
			Supplier: &modelsSupplier.Supplier{
				ID:   1,
				Name: "Fornecedor Teste",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelsSupplierCategories.SupplierCategory{
				{ID: 100},
			},
		}

		mockRepoSupplier.On("BeginTx", ctx).Return(tx, nil)
		createdSupplier := &modelsSupplier.Supplier{ID: 1, Name: "Fornecedor Teste"}
		mockRepoSupplier.On("CreateTx", ctx, tx, supplierFull.Supplier).Return(createdSupplier, nil)

		createdAddress := &modelsAddress.Address{
			ID:           10,
			Street:       "Rua Teste",
			StreetNumber: "45",
			City:         "Cidade Teste",
			State:        "SP",
			PostalCode:   "03459808",
			IsActive:     true,
			Country:      "Brasil",
			SupplierID:   utils.Int64Ptr(1),
		}
		mockRepoAddress.On("CreateTx", ctx, tx, supplierFull.Address).Return(createdAddress, nil)

		createdContact := &modelsContact.Contact{
			ID:          20,
			ContactName: "Contato Teste",
			Email:       "contato@example.com",
			Phone:       "1234567895",
		}
		mockRepoContact.On("CreateTx", ctx, tx, supplierFull.Contact).Return(createdContact, nil)

		relationErr := errors.New("erro ao criar relação no banco")
		mockRepoCatRel.On("CreateTx", ctx, tx, mock.AnythingOfType("*model.SupplierCategoryRelation")).
			Return(nil, relationErr)

		rollbackErr := errors.New("erro no rollback")
		tx.On("Rollback", ctx).Return(rollbackErr)

		result, err := service.CreateFull(ctx, supplierFull)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar relação no banco")
		assert.Contains(t, err.Error(), "rollback error")
		assert.Contains(t, err.Error(), rollbackErr.Error())
	})

	t.Run("CategoryRelation: sucesso com múltiplas categorias", func(t *testing.T) {
		mockRepoSupplier := new(mockSupplier.MockSupplierFullRepository)
		mockRepoAddress := new(mockAddress.MockAddressRepository)
		mockRepoContact := new(mockContact.MockContactRepository)
		mockRepoCatRel := new(mockCatRel.MockSupplierCategoryRelationRepo)
		tx := new(mockTX.MockTx)

		service := NewSupplierFull(
			mockRepoSupplier,
			mockRepoAddress,
			mockRepoContact,
			mockRepoCatRel,
			nil,
		)

		supplierFull := &modelsSupplierFull.SupplierFull{
			Supplier: &modelsSupplier.Supplier{
				ID:   1,
				Name: "Fornecedor Teste",
			},
			Address: &modelsAddress.Address{
				ID:           10,
				Street:       "Rua Teste",
				StreetNumber: "45",
				City:         "Cidade Teste",
				State:        "SP",
				PostalCode:   "03459808",
				IsActive:     true,
				Country:      "Brasil",
			},
			Contact: &modelsContact.Contact{
				ID:          20,
				ContactName: "Contato Teste",
				Email:       "contato@example.com",
				Phone:       "1234567895",
			},
			Categories: []modelsSupplierCategories.SupplierCategory{
				{ID: 100},
				{ID: 200},
				{ID: 300},
			},
		}

		mockRepoSupplier.On("BeginTx", ctx).Return(tx, nil)
		createdSupplier := &modelsSupplier.Supplier{ID: 1, Name: "Fornecedor Teste"}
		mockRepoSupplier.On("CreateTx", ctx, tx, supplierFull.Supplier).Return(createdSupplier, nil)

		createdAddress := &modelsAddress.Address{
			ID:           10,
			Street:       "Rua Teste",
			StreetNumber: "45",
			City:         "Cidade Teste",
			State:        "SP",
			PostalCode:   "03459808",
			IsActive:     true,
			Country:      "Brasil",
			SupplierID:   utils.Int64Ptr(1),
		}
		mockRepoAddress.On("CreateTx", ctx, tx, supplierFull.Address).Return(createdAddress, nil)

		createdContact := &modelsContact.Contact{
			ID:          20,
			ContactName: "Contato Teste",
			Email:       "contato@example.com",
			Phone:       "1234567895",
		}
		mockRepoContact.On("CreateTx", ctx, tx, supplierFull.Contact).Return(createdContact, nil)

		// Mock para múltiplas categorias
		for i := 0; i < len(supplierFull.Categories); i++ {
			mockRepoCatRel.On("CreateTx", ctx, tx, mock.AnythingOfType("*model.SupplierCategoryRelation")).
				Return(&modelsCatRel.SupplierCategoryRelation{
					SupplierID: 1,
					CategoryID: int64(supplierFull.Categories[i].ID),
					CreatedAt:  time.Now(),
				}, nil)
		}

		tx.On("Commit", ctx).Return(nil)

		result, err := service.CreateFull(ctx, supplierFull)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Categories, 3)
	})
}
