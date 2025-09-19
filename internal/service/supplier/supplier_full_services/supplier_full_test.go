package services

import (
	"context"
	"errors"
	"testing"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"

	mockTX "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo"
	mockAddress "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/address"
	mockContact "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/contact"
	mockSupplier "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/supplier"
	mockSupplierCatRel "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/supplier"
	modelAddress "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	modelContact "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	modelSupplier "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	modelSupplierCategories "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_categories"
	modelSupplierCatRel "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_category_relations"
	modelFull "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_full"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSupplierService_CreateFull(t *testing.T) {

	setup := func() (
		*mockSupplier.MockSupplierFullRepository,
		*mockAddress.MockAddressRepository,
		*mockContact.MockContactRepository,
		*mockSupplierCatRel.MockSupplierCategoryRelationRepo,
		*mockTX.MockTx,
		SupplierFullService,
	) {
		mockSupplierRepo := new(mockSupplier.MockSupplierFullRepository)
		mockAddressRepo := new(mockAddress.MockAddressRepository)
		mockContactRepo := new(mockContact.MockContactRepository)
		mockRelationRepo := new(mockSupplierCatRel.MockSupplierCategoryRelationRepo)
		mockTx := new(mockTX.MockTx)

		supplierService := NewSupplierFullService(
			mockSupplierRepo,
			mockAddressRepo,
			mockContactRepo,
			mockRelationRepo,
		)

		return mockSupplierRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockTx, supplierService
	}

	t.Run("transacao_nil", func(t *testing.T) {
		mockSupplierRepo, _, _, _, _, supplierService := setup()

		mockSupplierRepo.On("BeginTx", mock.Anything).Return(nil, nil)

		cpf := "12345678900"

		supplierFull := &modelFull.SupplierFull{
			Supplier: &modelSupplier.Supplier{
				Name:   "Fornecedor Teste",
				CPF:    &cpf,
				Status: true,
			},
			Address: &modelAddress.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &modelContact.Contact{
				ContactName: "João da Silva",
				Phone:       "11999999999",
				Email:       "joao@teste.com",
			},
			Categories: []modelSupplierCategories.SupplierCategory{
				{ID: 1},
				{ID: 2},
			},
		}

		_, err := supplierService.CreateFull(context.Background(), supplierFull)

		assert.Error(t, err)
		assert.EqualError(t, err, "transação inválida")

		mockSupplierRepo.AssertExpectations(t)
	})

	t.Run("supplierFull_ou_supplier_nil_e_validacao", func(t *testing.T) {
		_, _, _, _, _, supplierService := setup()

		// Caso 1: supplierFull é nil
		result, err := supplierService.CreateFull(context.Background(), nil)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)

		// Caso 2: supplierFull existe, mas Supplier é nil
		invalidSupplierFull := &modelFull.SupplierFull{}
		result, err = supplierService.CreateFull(context.Background(), invalidSupplierFull)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)

		// Caso 3: Supplier presente mas inválido (ex: nome vazio)
		invalidSupplierFull = &modelFull.SupplierFull{
			Supplier: &modelSupplier.Supplier{Name: ""},
		}
		result, err = supplierService.CreateFull(context.Background(), invalidSupplierFull)
		assert.Nil(t, result)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
	})

	t.Run("erro ao iniciar transação", func(t *testing.T) {
		mockSupplierRepo, _, _, _, _, supplierService := setup()

		cpf := "12345678900"

		supplierFull := &modelFull.SupplierFull{
			Supplier: &modelSupplier.Supplier{
				Name:   "Fornecedor Teste",
				CPF:    &cpf,
				Status: true,
			},
			Address: &modelAddress.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &modelContact.Contact{
				ContactName: "João da Silva",
				Phone:       "11999999999",
				Email:       "joao@teste.com",
			},
			Categories: []modelSupplierCategories.SupplierCategory{
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

		mockSupplierRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)

		mockSupplierRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&modelSupplier.Supplier{ID: 1}, nil)

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&modelAddress.Address{ID: 1}, nil)

		mockContactRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&modelContact.Contact{ID: 1}, nil)

		mockRelationRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(nil, errors.New("erro ao criar relação"))

		mockTx.On("Rollback", mock.Anything).Return(errors.New("erro ao dar rollback"))

		cpf := "12345678900"
		supplierFull := &modelFull.SupplierFull{
			Supplier: &modelSupplier.Supplier{
				Name:   "Fornecedor Teste",
				CPF:    &cpf,
				Status: true,
			},
			Address: &modelAddress.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &modelContact.Contact{
				ContactName: "João da Silva",
				Phone:       "1112345678",
				Email:       "joao@teste.com",
			},
			Categories: []modelSupplierCategories.SupplierCategory{
				{ID: 1},
			},
		}

		_, err := supplierService.CreateFull(context.Background(), supplierFull)

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

		mockSupplierRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)

		mockSupplierRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&modelSupplier.Supplier{
				ID:   1,
				Name: "Fornecedor Teste",
			}, nil)

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&modelAddress.Address{ID: 1}, nil)

		mockContactRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&modelContact.Contact{ID: 1}, nil)

		mockRelationRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&modelSupplierCatRel.SupplierCategoryRelations{SupplierID: 1, CategoryID: 1}, nil)

		mockTx.On("Commit", mock.Anything).Return(errors.New("erro ao commitar transação"))

		mockTx.On("Rollback", mock.Anything).Return(nil).Once()

		cpf := "12345678900"
		supplierFull := &modelFull.SupplierFull{
			Supplier: &modelSupplier.Supplier{
				Name:   "Fornecedor Teste",
				CPF:    &cpf,
				Status: true,
			},
			Address: &modelAddress.Address{
				Street:     "Rua Teste",
				City:       "Cidade Teste",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &modelContact.Contact{
				ContactName: "Contato Teste",
				Phone:       "1112345678",
				Email:       "contato@teste.com",
			},
			Categories: []modelSupplierCategories.SupplierCategory{
				{ID: 1},
			},
		}

		_, err := supplierService.CreateFull(context.Background(), supplierFull)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao commitar transação")

		mockSupplierRepo.AssertExpectations(t)
		mockAddressRepo.AssertExpectations(t)
		mockContactRepo.AssertExpectations(t)
		mockRelationRepo.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("erro_ao_fazer_rollback_apos_commit_falhar", func(t *testing.T) {

		mockSupplierRepo, mockAddressRepo, mockContactRepo, mockRelationRepo, mockTx, supplierService := setup()

		cpf := "12345678900"
		supplierFull := &modelFull.SupplierFull{
			Supplier: &modelSupplier.Supplier{
				Name:   "Fornecedor Teste",
				CPF:    &cpf,
				Status: true,
			},
			Address: &modelAddress.Address{
				Street:     "Rua Teste",
				City:       "Cidade Teste",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &modelContact.Contact{
				ContactName: "Contato Teste",
				Phone:       "1112345678",
				Email:       "contato@teste.com",
			},
			Categories: []modelSupplierCategories.SupplierCategory{
				{ID: 1},
			},
		}

		mockSupplierRepo.On("BeginTx", mock.Anything).Return(mockTx, nil)

		mockSupplierRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&modelSupplier.Supplier{ID: 1}, nil)

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&modelAddress.Address{ID: 1}, nil)

		mockContactRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&modelContact.Contact{ID: 1}, nil)

		mockRelationRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(nil, nil)

		commitError := errors.New("erro no commit")
		mockTx.On("Commit", mock.Anything).Return(commitError)

		rollbackError := errors.New("erro no rollback")
		mockTx.On("Rollback", mock.Anything).Return(rollbackError)

		_, err := supplierService.CreateFull(context.Background(), supplierFull)

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
			Return(&modelSupplier.Supplier{
				ID:     1,
				Name:   "Fornecedor Teste",
				Status: true,
			}, nil)

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(nil, errors.New("erro ao criar endereço"))

		mockRelationRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(nil, nil).
			Maybe()

		mockTx.On("Rollback", mock.Anything).Return(nil)

		cpf := "12345678900"
		supplierFull := &modelFull.SupplierFull{
			Supplier: &modelSupplier.Supplier{
				Name:   "Fornecedor Teste",
				CPF:    &cpf,
				Status: true,
			},
			Address: &modelAddress.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &modelContact.Contact{
				ContactName: "Contato Teste",
				Phone:       "1112345678",
				Email:       "contato@teste.com",
			},
			Categories: []modelSupplierCategories.SupplierCategory{
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
			Return(&modelSupplier.Supplier{
				ID:     1,
				Name:   "Fornecedor Teste",
				Status: true,
			}, nil)

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&modelAddress.Address{
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
		supplierFull := &modelFull.SupplierFull{
			Supplier: &modelSupplier.Supplier{
				Name:   "Fornecedor Teste",
				CPF:    &cpf,
				Status: true,
			},
			Address: &modelAddress.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &modelContact.Contact{
				ContactName: "Ari",
				Phone:       "1234567895",
				Email:       "contato@example.com",
			},
			Categories: []modelSupplierCategories.SupplierCategory{
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

		mockSupplierRepo.On("CreateTx", mock.Anything, mockTx, mock.MatchedBy(func(supplier *modelSupplier.Supplier) bool {
			return supplier.Name == "Fornecedor Teste"
		})).Return(nil, errors.New("falha ao criar fornecedor")).Once()

		mockTx.On("Rollback", mock.Anything).Return(nil).Once()

		cpf := "12345678900"
		supplierFull := &modelFull.SupplierFull{
			Supplier: &modelSupplier.Supplier{
				Name:   "Fornecedor Teste",
				CPF:    &cpf,
				Status: true,
			},
			Categories: []modelSupplierCategories.SupplierCategory{
				{ID: 1},
			},
			Contact: &modelContact.Contact{
				ContactName: "Ari",
				Phone:       "1234567895",
				Email:       "contato@example.com",
			},
			Address: &modelAddress.Address{
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
		supplierFull := &modelFull.SupplierFull{
			Supplier: &modelSupplier.Supplier{
				Name:   "Fornecedor Teste",
				CPF:    &cpf,
				Status: true,
			},
			Address: &modelAddress.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &modelContact.Contact{
				ContactName: "Ari",
				Phone:       "1234567895",
				Email:       "contato@example.com",
			},
			Categories: []modelSupplierCategories.SupplierCategory{
				{ID: 1},
			},
		}

		expectedSupplier := &modelSupplier.Supplier{
			ID:     1,
			Name:   "Fornecedor Teste",
			CPF:    &cpf,
			Status: true,
		}

		mockSupplierRepo.On("BeginTx", mock.Anything).
			Return(mockTx, nil).Once()

		mockSupplierRepo.On("CreateTx", mock.Anything, mockTx, mock.MatchedBy(func(s *modelSupplier.Supplier) bool {
			return s.Name == "Fornecedor Teste" && s.CPF != nil && *s.CPF == cpf
		})).Return(expectedSupplier, nil).Once()

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.MatchedBy(func(addr *modelAddress.Address) bool {
			return addr.City == "Cidade B" && addr.PostalCode == "12345678"
		})).Return(&modelAddress.Address{ID: 1}, nil).Once()

		mockContactRepo.On("CreateTx", mock.Anything, mockTx, mock.MatchedBy(func(c *modelContact.Contact) bool {
			return c.Phone == "1234567895" && c.ContactName == "Ari"
		})).Return(&modelContact.Contact{ID: 1}, nil).Once()

		mockRelationRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&modelSupplierCatRel.SupplierCategoryRelations{SupplierID: 1, CategoryID: 1}, nil).Once()

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
		supplierFull := &modelFull.SupplierFull{
			Supplier: &modelSupplier.Supplier{
				Name:   "Fornecedor Teste",
				CPF:    &cpf,
				Status: true,
			},
			Address: &modelAddress.Address{
				Street: "",
			},
			Contact: &modelContact.Contact{
				ContactName: "Ari",
				Phone:       "1234567895",
				Email:       "contato@example.com",
			},
			Categories: []modelSupplierCategories.SupplierCategory{
				{ID: 1},
			},
		}

		mockSupplierRepo.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
		mockSupplierRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&modelSupplier.Supplier{ID: 1, Name: "Fornecedor Teste", CPF: &cpf}, nil).Once()

		mockTx.On("Rollback", mock.Anything).Return(nil).Once()

		_, err := supplierService.CreateFull(context.Background(), supplierFull)
		assert.ErrorContains(t, err, "endereço inválido")

		mockSupplierRepo.AssertExpectations(t)
		mockTx.AssertExpectations(t)
	})

	t.Run("falha_validacao_do_contato_faz_rollback", func(t *testing.T) {
		mockSupplierRepo, mockAddressRepo, _, _, mockTx, supplierService := setup()

		cpf := "12345678900"
		supplierFull := &modelFull.SupplierFull{
			Supplier: &modelSupplier.Supplier{
				Name:   "Fornecedor Teste",
				CPF:    &cpf,
				Status: true,
			},
			Categories: []modelSupplierCategories.SupplierCategory{
				{ID: 1},
			},
			Address: &modelAddress.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &modelContact.Contact{
				Phone: "invalido",
			},
		}

		mockSupplierRepo.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
		mockSupplierRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&modelSupplier.Supplier{ID: 1, Name: "Fornecedor Teste", CPF: &cpf}, nil).Once()

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&modelAddress.Address{ID: 1}, nil).Once()

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
		supplierFull := &modelFull.SupplierFull{
			Supplier: &modelSupplier.Supplier{
				Name:   "Fornecedor Teste",
				CPF:    &cpf,
				Status: true,
			},
			Categories: []modelSupplierCategories.SupplierCategory{
				{ID: 0},
			},
			Address: &modelAddress.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &modelContact.Contact{
				ContactName: "João",
				Phone:       "1234567890",
				Email:       "joao@example.com",
			},
		}

		mockSupplierRepo.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
		mockSupplierRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&modelSupplier.Supplier{ID: 1, Name: "Fornecedor Teste", CPF: &cpf}, nil).Once()

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&modelAddress.Address{ID: 1}, nil).Once()

		mockContactRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&modelContact.Contact{ID: 1}, nil).Once()

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
		supplierFull := &modelFull.SupplierFull{
			Supplier: &modelSupplier.Supplier{
				Name:   "Fornecedor Teste",
				CPF:    &cpf,
				Status: true,
			},
			Categories: []modelSupplierCategories.SupplierCategory{
				{ID: 1},
			},
			Address: &modelAddress.Address{
				Street:     "Rua A",
				City:       "Cidade B",
				State:      "SP",
				Country:    "Brasil",
				PostalCode: "12345678",
			},
			Contact: &modelContact.Contact{
				ContactName: "Ari",
				Phone:       "1234567895",
				Email:       "contato@example.com",
			},
		}

		mockSupplierRepo.On("BeginTx", mock.Anything).Return(mockTx, nil).Once()
		mockSupplierRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&modelSupplier.Supplier{ID: 1, Name: "Fornecedor Teste", CPF: &cpf}, nil).Once()

		mockAddressRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Return(&modelAddress.Address{ID: 1}, nil).Once()

		mockContactRepo.On("CreateTx", mock.Anything, mockTx, mock.Anything).
			Run(func(_ mock.Arguments) {
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
