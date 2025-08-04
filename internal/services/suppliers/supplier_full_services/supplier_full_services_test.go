package services

import (
	"context"
	"errors"
	"testing"

	model_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	model_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	models_supplier "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier"
	models_supplier_categories "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_categories"
	models_supplier_cat_rel "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_category_relations"
	models_full "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_full"
	repo_address "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
	repo_contact "github.com/WagaoCarvalho/backend_store_go/internal/repositories/contacts"
	repo_tx "github.com/WagaoCarvalho/backend_store_go/internal/repositories/mocks"
	repo_relation "github.com/WagaoCarvalho/backend_store_go/internal/repositories/suppliers/supplier_category_relations"
	repo_supplier "github.com/WagaoCarvalho/backend_store_go/internal/repositories/suppliers/supplier_full_repositories"
	"github.com/WagaoCarvalho/backend_store_go/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreateFull_Validation(t *testing.T) {
	// Setup dos mocks
	mockSupplierRepo := new(repo_supplier.MockSupplierFullRepository)
	mockAddressRepo := new(repo_address.MockAddressRepository)
	mockContactRepo := new(repo_contact.MockContactRepository)
	mockRelationRepo := new(repo_relation.MockSupplierCategoryRelationRepo)
	logger := logger.NewLoggerAdapter(logrus.New())

	service := NewSupplierFullService(
		mockSupplierRepo,
		mockAddressRepo,
		mockContactRepo,
		mockRelationRepo,
		logger,
	)

	ctx := context.Background()

	t.Run("deve_falhar_quando_supplier_nil", func(t *testing.T) {
		invalidSupplier := &models_full.SupplierFull{
			Address:    &model_address.Address{Street: "Rua Teste"},
			Contact:    &model_contact.Contact{Phone: "1112345678"},
			Categories: []models_supplier_categories.SupplierCategory{{ID: 1}},
		}

		result, err := service.CreateFull(ctx, invalidSupplier)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "fornecedor é obrigatório", err.Error())
	})

	t.Run("deve_falhar_quando_address_nil", func(t *testing.T) {
		invalidSupplier := &models_full.SupplierFull{
			Supplier: &models_supplier.Supplier{
				Name:   "Fornecedor Teste",
				Status: true,
			},
			Contact:    &model_contact.Contact{Phone: "1112345678"},
			Categories: []models_supplier_categories.SupplierCategory{{ID: 1}},
		}

		result, err := service.CreateFull(ctx, invalidSupplier)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "endereço é obrigatório", err.Error())
	})

	t.Run("deve_falhar_quando_contact_nil", func(t *testing.T) {
		invalidSupplier := &models_full.SupplierFull{
			Supplier: &models_supplier.Supplier{
				Name:   "Fornecedor Teste",
				Status: true,
			},
			Address:    &model_address.Address{Street: "Rua Teste"},
			Categories: []models_supplier_categories.SupplierCategory{{ID: 1}},
		}

		result, err := service.CreateFull(ctx, invalidSupplier)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "contato é obrigatório", err.Error())
	})

	t.Run("deve_falhar_quando_sem_categorias", func(t *testing.T) {
		invalidSupplier := &models_full.SupplierFull{
			Supplier: &models_supplier.Supplier{
				Name:   "Fornecedor Teste",
				Status: true,
			},
			Address: &model_address.Address{Street: "Rua Teste"},
			Contact: &model_contact.Contact{Phone: "1112345678"},
		}

		result, err := service.CreateFull(ctx, invalidSupplier)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Equal(t, "pelo menos uma categoria é obrigatória", err.Error())
	})

	t.Run("deve_falhar_quando_supplier_invalido", func(t *testing.T) {
		invalidSupplier := &models_full.SupplierFull{
			Supplier: &models_supplier.Supplier{
				Name: "", // nome inválido
			},
			Address:    &model_address.Address{Street: "Rua Teste"},
			Contact:    &model_contact.Contact{Phone: "1112345678"},
			Categories: []models_supplier_categories.SupplierCategory{{ID: 1}},
		}

		result, err := service.CreateFull(ctx, invalidSupplier)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "fornecedor inválido")
	})

	t.Run("nao_deve_iniciar_transacao_quando_validacao_falhar", func(t *testing.T) {
		mockSupplierRepo.AssertNotCalled(t, "BeginTx")
	})

}

func TestSupplierService_CreateFull(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	setup := func() (
		*repo_supplier.MockSupplierFullRepository,
		*repo_address.MockAddressRepository,
		*repo_contact.MockContactRepository,
		*repo_relation.MockSupplierCategoryRelationRepo,
		*repo_tx.MockTx,
		SupplierFullService,
	) {
		mockSupplierRepo := new(repo_supplier.MockSupplierFullRepository)
		mockAddressRepo := new(repo_address.MockAddressRepository)
		mockContactRepo := new(repo_contact.MockContactRepository)
		mockRelationRepo := new(repo_relation.MockSupplierCategoryRelationRepo)
		mockTx := new(repo_tx.MockTx)

		supplierService := NewSupplierFullService(
			mockSupplierRepo,
			mockAddressRepo,
			mockContactRepo,
			mockRelationRepo,
			logger,
		)

		return mockSupplierRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockTx, supplierService
	}

	t.Run("transacao_nil", func(t *testing.T) {
		mockSupplierRepo, _, _, _, _, supplierService := setup()

		mockSupplierRepo.On("BeginTx", mock.Anything).Return(nil, nil)

		cpf := "12345678900"

		supplierFull := &models_full.SupplierFull{
			Supplier: &models_supplier.Supplier{
				Name:   "Fornecedor Teste",
				CPF:    &cpf,
				Status: true,
			},
			Address: &model_address.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &model_contact.Contact{
				ContactName: "João da Silva",
				Phone:       "11999999999",
				Email:       "joao@teste.com",
			},
			Categories: []models_supplier_categories.SupplierCategory{
				{ID: 1},
				{ID: 2},
			},
		}

		_, err := supplierService.CreateFull(context.Background(), supplierFull)

		assert.Error(t, err)
		assert.EqualError(t, err, "transação inválida")

		mockSupplierRepo.AssertExpectations(t)
	})

	t.Run("erro ao iniciar transação", func(t *testing.T) {
		mockSupplierRepo, _, _, _, _, supplierService := setup()

		cpf := "12345678900"

		supplierFull := &models_full.SupplierFull{
			Supplier: &models_supplier.Supplier{
				Name:   "Fornecedor Teste",
				CPF:    &cpf,
				Status: true,
			},
			Address: &model_address.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &model_contact.Contact{
				ContactName: "João da Silva",
				Phone:       "11999999999",
				Email:       "joao@teste.com",
			},
			Categories: []models_supplier_categories.SupplierCategory{
				{ID: 1},
				{ID: 2},
			},
		}

		mockSupplierRepo.On("BeginTx", mock.Anything).Return(nil, errors.New("falha na transação")).Once()

		_, err := supplierService.CreateFull(context.Background(), supplierFull)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao iniciar transação")

		mockSupplierRepo.AssertExpectations(t)
	})

	t.Run("erro_ao_fazer_rollback", func(t *testing.T) {
		mockSupplierRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockTx, supplierService := setup()

		// Configuração dos mocks
		mockSupplierRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)

		// Mock das criações bem-sucedidas
		mockSupplierRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&models_supplier.Supplier{ID: 1}, nil)

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_address.Address{ID: 1}, nil)

		mockContactRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_contact.Contact{ID: 1}, nil)

		// Mock com falha na criação da relação
		mockRelationRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(nil, errors.New("erro ao criar relação"))

		// Mock com falha no rollback
		mockTx.On("Rollback", mock.Anything).Return(errors.New("erro ao dar rollback"))

		// Dados de entrada
		cpf := "12345678900"
		supplierFull := &models_full.SupplierFull{
			Supplier: &models_supplier.Supplier{
				Name:   "Fornecedor Teste",
				CPF:    &cpf,
				Status: true,
			},
			Address: &model_address.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &model_contact.Contact{
				ContactName: "João da Silva",
				Phone:       "1112345678",
				Email:       "joao@teste.com",
			},
			Categories: []models_supplier_categories.SupplierCategory{
				{ID: 1},
			},
		}

		// Execução
		_, err := supplierService.CreateFull(context.Background(), supplierFull)

		// Verificações
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar relação")
		assert.Contains(t, err.Error(), "rollback error")

		mockSupplierRepo.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
		mockContactRepo.AssertExpectations(t)
		mockRelationRepo.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("erro_ao_commitar_transacao", func(t *testing.T) {
		mockSupplierRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockTx, supplierService := setup()

		// Mock da transação
		mockSupplierRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)

		// Mock da criação do fornecedor
		mockSupplierRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&models_supplier.Supplier{
				ID:   1,
				Name: "Fornecedor Teste",
			}, nil)

		// Mock da criação do endereço
		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_address.Address{ID: 1}, nil)

		// Mock da criação do contato
		mockContactRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_contact.Contact{ID: 1}, nil)

		// Mock da criação da relação com categoria
		mockRelationRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&models_supplier_cat_rel.SupplierCategoryRelations{SupplierID: 1, CategoryID: 1}, nil)

		// Mock do commit que falha
		mockTx.On("Commit", mock.Anything).Return(errors.New("erro ao commitar transação"))

		// Mock do rollback que deve ser chamado após erro no commit
		mockTx.On("Rollback", mock.Anything).Return(nil).Once()

		// Dados de entrada
		cpf := "12345678900"
		supplierFull := &models_full.SupplierFull{
			Supplier: &models_supplier.Supplier{
				Name:   "Fornecedor Teste",
				CPF:    &cpf,
				Status: true,
			},
			Address: &model_address.Address{
				Street:     "Rua Teste",
				City:       "Cidade Teste",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &model_contact.Contact{
				ContactName: "Contato Teste",
				Phone:       "1112345678",
				Email:       "contato@teste.com",
			},
			Categories: []models_supplier_categories.SupplierCategory{
				{ID: 1},
			},
		}

		// Execução
		_, err := supplierService.CreateFull(context.Background(), supplierFull)

		// Verificações
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao commitar transação")

		// Verificações de expectativa
		mockSupplierRepo.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
		mockContactRepo.AssertExpectations(t)
		mockRelationRepo.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("erro_ao_fazer_rollback_apos_commit_falhar", func(t *testing.T) {
		// Setup
		mockSupplierRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockTx, supplierService := setup()

		// Dados válidos
		cpf := "12345678900"
		supplierFull := &models_full.SupplierFull{
			Supplier: &models_supplier.Supplier{
				Name:   "Fornecedor Teste",
				CPF:    &cpf,
				Status: true,
			},
			Address: &model_address.Address{
				Street:     "Rua Teste",
				City:       "Cidade Teste",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &model_contact.Contact{
				ContactName: "Contato Teste",
				Phone:       "1112345678",
				Email:       "contato@teste.com",
			},
			Categories: []models_supplier_categories.SupplierCategory{
				{ID: 1},
			},
		}

		// Configuração dos mocks
		mockSupplierRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)

		mockSupplierRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&models_supplier.Supplier{ID: 1}, nil)

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_address.Address{ID: 1}, nil)

		mockContactRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_contact.Contact{ID: 1}, nil)

		// Aqui erro na criação da relação, mas retorna nil erro para simular sucesso
		mockRelationRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(nil, nil)

		// Mock do commit que falha
		commitError := errors.New("erro no commit")
		mockTx.On("Commit", mock.Anything).Return(commitError)

		// Mock do rollback que também falha
		rollbackError := errors.New("erro no rollback")
		mockTx.On("Rollback", mock.Anything).Return(rollbackError)

		// Execução
		_, err := supplierService.CreateFull(context.Background(), supplierFull)

		// Verificações
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao commitar transação")
		assert.Contains(t, err.Error(), "erro no commit")
		assert.Contains(t, err.Error(), "rollback error")
		assert.Contains(t, err.Error(), "erro no rollback")

		mockSupplierRepo.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
		mockContactRepo.AssertExpectations(t)
		mockRelationRepo.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("erro_ao_criar_endereco", func(t *testing.T) {
		mockSupplierRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockTx, supplierService := setup()

		mockSupplierRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)

		mockSupplierRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&models_supplier.Supplier{
				ID:     1,
				Name:   "Fornecedor Teste",
				Status: true,
			}, nil)

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(nil, errors.New("erro ao criar endereço"))

		// Relação pode ou não ser chamada, pode usar Maybe()
		mockRelationRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(nil, nil).
			Maybe()

		mockTx.On("Rollback", mock.Anything).Return(nil)

		cpf := "12345678900"
		supplierFull := &models_full.SupplierFull{
			Supplier: &models_supplier.Supplier{
				Name:   "Fornecedor Teste",
				CPF:    &cpf,
				Status: true,
			},
			Address: &model_address.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &model_contact.Contact{
				ContactName: "Contato Teste",
				Phone:       "1112345678",
				Email:       "contato@teste.com",
			},
			Categories: []models_supplier_categories.SupplierCategory{
				{ID: 1},
			},
		}

		_, err := supplierService.CreateFull(context.Background(), supplierFull)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar endereço")

		mockSupplierRepo.AssertExpectations(t)
		mockRelationRepo.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
		mockTx.AssertExpectations(t)
		mockContactRepo.AssertExpectations(t)
	})

	t.Run("erro_ao_criar_contato", func(t *testing.T) {
		mockSupplierRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockTx, supplierService := setup()

		mockSupplierRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)

		mockSupplierRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&models_supplier.Supplier{
				ID:     1,
				Name:   "Fornecedor Teste",
				Status: true,
			}, nil)

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_address.Address{
				ID:         1,
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			}, nil)

		mockContactRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(nil, errors.New("erro ao criar contato"))

		mockRelationRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(nil, nil).
			Maybe()

		mockTx.On("Rollback", mock.Anything).Return(nil)

		cpf := "12345678900"
		supplierFull := &models_full.SupplierFull{
			Supplier: &models_supplier.Supplier{
				Name:   "Fornecedor Teste",
				CPF:    &cpf,
				Status: true,
			},
			Address: &model_address.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &model_contact.Contact{
				ContactName: "Ari",
				Phone:       "1234567895",
				Email:       "contato@example.com",
			},
			Categories: []models_supplier_categories.SupplierCategory{
				{ID: 1},
			},
		}

		_, err := supplierService.CreateFull(context.Background(), supplierFull)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar contato")

		mockSupplierRepo.AssertExpectations(t)
		mockRelationRepo.AssertExpectations(t)
		mockContactRepo.AssertExpectations(t)
		mockTx.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
	})
	t.Run("erro_ao_criar_fornecedor_na_transacao", func(t *testing.T) {
		mockSupplierRepo, _, _, _, mockTx, supplierService := setup()

		mockSupplierRepo.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()

		mockSupplierRepo.On("CreateTx", mock.Anything, mockTx, mock.MatchedBy(func(supplier *models_supplier.Supplier) bool {
			return supplier.Name == "Fornecedor Teste"
		})).Return(nil, errors.New("falha ao criar fornecedor")).Once()

		mockTx.On("Rollback", mock.Anything).Return(nil).Once()

		cpf := "12345678900"
		supplierFull := &models_full.SupplierFull{
			Supplier: &models_supplier.Supplier{
				Name:   "Fornecedor Teste",
				CPF:    &cpf,
				Status: true,
			},
			Categories: []models_supplier_categories.SupplierCategory{
				{ID: 1},
			},
			Contact: &model_contact.Contact{
				ContactName: "Ari",
				Phone:       "1234567895",
				Email:       "contato@example.com",
			},
			Address: &model_address.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "99999999",
			},
		}

		_, err := supplierService.CreateFull(context.Background(), supplierFull)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "falha ao criar fornecedor")

		mockSupplierRepo.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})
	t.Run("sucesso_na_criacao_completa", func(t *testing.T) {
		mockSupplierRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockTx, supplierService := setup()

		cpf := "12345678900"
		supplierFull := &models_full.SupplierFull{
			Supplier: &models_supplier.Supplier{
				Name:   "Fornecedor Teste",
				CPF:    &cpf,
				Status: true,
			},
			Address: &model_address.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &model_contact.Contact{
				ContactName: "Ari",
				Phone:       "1234567895",
				Email:       "contato@example.com",
			},
			Categories: []models_supplier_categories.SupplierCategory{
				{ID: 1},
			},
		}

		expectedSupplier := &models_supplier.Supplier{
			ID:     1,
			Name:   "Fornecedor Teste",
			CPF:    &cpf,
			Status: true,
		}

		mockSupplierRepo.On("BeginTx", mock.Anything).
			Return(mockTx, nil).Once()

		mockSupplierRepo.On("CreateTx", mock.Anything, mockTx, mock.MatchedBy(func(s *models_supplier.Supplier) bool {
			return s.Name == "Fornecedor Teste" && s.CPF != nil && *s.CPF == cpf
		})).Return(expectedSupplier, nil).Once()

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.MatchedBy(func(addr *model_address.Address) bool {
			return addr.City == "Cidade B" && addr.PostalCode == "12345678"
		})).Return(&model_address.Address{ID: 1}, nil).Once()

		mockContactRepo.On("CreateTx", mock.Anything, mockTx, mock.MatchedBy(func(c *model_contact.Contact) bool {
			return c.Phone == "1234567895" && c.ContactName == "Ari"
		})).Return(&model_contact.Contact{ID: 1}, nil).Once()

		mockRelationRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&models_supplier_cat_rel.SupplierCategoryRelations{SupplierID: 1, CategoryID: 1}, nil).Once()

		mockTx.On("Commit", mock.Anything).Return(nil).Once()

		result, err := supplierService.CreateFull(context.Background(), supplierFull)

		assert.NoError(t, err)
		assert.Equal(t, expectedSupplier.ID, result.Supplier.ID)
		assert.Equal(t, expectedSupplier.Name, result.Supplier.Name)
		assert.Equal(t, expectedSupplier.Status, result.Supplier.Status)

		mockSupplierRepo.AssertExpectations(t)
		mockTx.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
		mockContactRepo.AssertExpectations(t)
		mockRelationRepo.AssertExpectations(t)
	})
	t.Run("falha_validacao_do_endereco_faz_rollback", func(t *testing.T) {
		mockSupplierRepo, _, _, _, mockTx, supplierService := setup()

		cpf := "12345678900"
		supplierFull := &models_full.SupplierFull{
			Supplier: &models_supplier.Supplier{
				Name:   "Fornecedor Teste",
				CPF:    &cpf,
				Status: true,
			},
			Address: &model_address.Address{
				Street: "", // força falha de validação
			},
			Contact: &model_contact.Contact{
				ContactName: "Ari",
				Phone:       "1234567895",
				Email:       "contato@example.com",
			},
			Categories: []models_supplier_categories.SupplierCategory{
				{ID: 1},
			},
		}

		mockSupplierRepo.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
		mockSupplierRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&models_supplier.Supplier{ID: 1, Name: "Fornecedor Teste", CPF: &cpf}, nil).Once()

		mockTx.On("Rollback", mock.Anything).Return(nil).Once()

		_, err := supplierService.CreateFull(context.Background(), supplierFull)
		assert.ErrorContains(t, err, "endereço inválido")

		mockSupplierRepo.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("falha_validacao_do_contato_faz_rollback", func(t *testing.T) {
		mockSupplierRepo, mockAddressRepo, _, _, mockTx, supplierService := setup()

		cpf := "12345678900"
		supplierFull := &models_full.SupplierFull{
			Supplier: &models_supplier.Supplier{
				Name:   "Fornecedor Teste",
				CPF:    &cpf,
				Status: true,
			},
			Categories: []models_supplier_categories.SupplierCategory{
				{ID: 1},
			},
			Address: &model_address.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &model_contact.Contact{
				Phone: "invalido", // força erro de validação
			},
		}

		mockSupplierRepo.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
		mockSupplierRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&models_supplier.Supplier{ID: 1, Name: "Fornecedor Teste", CPF: &cpf}, nil).Once()

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_address.Address{ID: 1}, nil).Once()

		mockTx.On("Rollback", mock.Anything).Return(nil).Once()

		_, err := supplierService.CreateFull(context.Background(), supplierFull)
		assert.ErrorContains(t, err, "contato inválido")

		mockSupplierRepo.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})
	t.Run("falha_validacao_relacao_supplier_categoria_faz_rollback", func(t *testing.T) {
		mockSupplierRepo, mockAddressRepo, mockContactRepo, _, mockTx, supplierService := setup()

		cpf := "12345678900"
		supplierFull := &models_full.SupplierFull{
			Supplier: &models_supplier.Supplier{
				Name:   "Fornecedor Teste",
				CPF:    &cpf,
				Status: true,
			},
			Categories: []models_supplier_categories.SupplierCategory{
				{ID: 0}, // força falha na relação categoria
			},
			Address: &model_address.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &model_contact.Contact{
				ContactName: "João",
				Phone:       "1234567890",
				Email:       "joao@example.com",
			},
		}

		mockSupplierRepo.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
		mockSupplierRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&models_supplier.Supplier{ID: 1, Name: "Fornecedor Teste", CPF: &cpf}, nil).Once()

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_address.Address{ID: 1}, nil).Once()

		mockContactRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_contact.Contact{ID: 1}, nil).Once()

		mockTx.On("Rollback", mock.Anything).Return(nil).Once()

		_, err := supplierService.CreateFull(context.Background(), supplierFull)
		assert.ErrorContains(t, err, "relação fornecedor-categoria inválida")

		mockSupplierRepo.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
		mockContactRepo.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})
	t.Run("panic_faz_rollback", func(t *testing.T) {
		mockSupplierRepo, mockAddressRepo, mockContactRepo, _, mockTx, supplierService := setup()

		cpf := "12345678900"
		supplierFull := &models_full.SupplierFull{
			Supplier: &models_supplier.Supplier{
				Name:   "Fornecedor Teste",
				CPF:    &cpf,
				Status: true,
			},
			Categories: []models_supplier_categories.SupplierCategory{
				{ID: 1},
			},
			Address: &model_address.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &model_contact.Contact{
				ContactName: "Ari",
				Phone:       "1234567895",
				Email:       "contato@example.com",
			},
		}

		mockSupplierRepo.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
		mockSupplierRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&models_supplier.Supplier{ID: 1, Name: "Fornecedor Teste", CPF: &cpf}, nil).Once()

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&model_address.Address{ID: 1}, nil).Once()

		mockContactRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Run(func(args mock.Arguments) {
				panic("panic simulado")
			}).Once()

		mockTx.On("Rollback", mock.Anything).Return(nil).Once()

		assert.Panics(t, func() {
			_, _ = supplierService.CreateFull(context.Background(), supplierFull)
		})

		mockTx.AssertExpectations(t)
		mockSupplierRepo.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
		mockContactRepo.AssertExpectations(t)
	})

}
