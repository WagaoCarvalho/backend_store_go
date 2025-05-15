package services

import (
	"context"
	"fmt"
	"testing"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier"
	supplier_category "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_categories"
	supplier_category_relations "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_category_relations"
	mock_supplier "github.com/WagaoCarvalho/backend_store_go/internal/services/suppliers/supplier_service_mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

//
// Casos de sucesso
//

func TestCreateSupplier_Success(t *testing.T) {
	ctx := context.TODO()

	mockSupplierRepo := new(mock_supplier.MockSupplierRepository)
	mockRelationService := new(mock_supplier.MockSupplierCategoryRelationService)
	mockAddressService := new(mock_supplier.MockAddressService)
	mockContactService := new(mock_supplier.MockContactService)
	mockCategoryService := new(mock_supplier.MockSupplierCategoryService)

	service := &supplierService{
		repo:            mockSupplierRepo,
		relationService: mockRelationService,
		addressService:  mockAddressService,
		contactService:  mockContactService,
		categoryService: mockCategoryService,
	}

	supplier := models.Supplier{Name: "Fornecedor A"}
	supplierWithID := supplier
	supplierWithID.ID = 123 // ID que será retornado pelo mock

	address := &models_address.Address{Street: "Rua A", City: "Cidade", State: "Estado", PostalCode: "12345"}
	contact := &models_contact.Contact{ContactName: "Contato A", Email: "contato@example.com"}
	categoryID := int64(1)

	// Mock ajustado para retornar o supplier completo com ID
	mockSupplierRepo.On("Create", ctx, supplier).Return(supplierWithID, nil)

	// Os outros mocks permanecem iguais
	mockRelationService.On("HasRelation", ctx, supplierWithID.ID, categoryID).Return(false, nil)
	mockRelationService.On("Create", ctx, supplierWithID.ID, categoryID).Return(&supplier_category_relations.SupplierCategoryRelations{}, nil)
	mockAddressService.On("Create", ctx, mock.Anything).Return(*address, nil)
	mockContactService.On("Create", ctx, mock.Anything).Return(*contact, nil)
	mockCategoryService.On("GetByID", ctx, categoryID).Return(&supplier_category.SupplierCategory{ID: categoryID, Name: "Categoria A"}, nil)

	// Supondo que service.Create retorna o ID do supplier criado
	id, err := service.Create(ctx, &supplier, categoryID, address, contact)

	assert.NoError(t, err)
	assert.Equal(t, supplierWithID.ID, id) // Verifica se o ID retornado é o esperado
	mockSupplierRepo.AssertExpectations(t)
	mockRelationService.AssertExpectations(t)
	mockAddressService.AssertExpectations(t)
	mockContactService.AssertExpectations(t)
}

//
// Casos de erro de entrada
//

func TestCreateSupplier_InvalidName(t *testing.T) {
	service := &supplierService{}
	_, err := service.Create(context.TODO(), &models.Supplier{Name: ""}, 0, nil, nil)

	assert.Error(t, err)
	assert.Equal(t, "nome do fornecedor é obrigatório", err.Error())
}

func TestCreateSupplier_RepoError(t *testing.T) {
	ctx := context.TODO()

	// Mocks
	mockSupplierRepo := new(mock_supplier.MockSupplierRepository)
	mockRelationService := new(mock_supplier.MockSupplierCategoryRelationService)
	mockAddressService := new(mock_supplier.MockAddressService)
	mockContactService := new(mock_supplier.MockContactService)
	mockCategoryService := new(mock_supplier.MockSupplierCategoryService)

	// Serviço com todos os mocks inicializados
	service := &supplierService{
		repo:            mockSupplierRepo,
		relationService: mockRelationService,
		addressService:  mockAddressService,
		contactService:  mockContactService,
		categoryService: mockCategoryService,
	}

	// Dados simulados - agora usando o valor (não ponteiro)
	supplier := models.Supplier{Name: "Fornecedor A"}

	// Mock ajustado para esperar models.Supplier (não *models.Supplier)
	mockSupplierRepo.On("Create", ctx, supplier).Return(models.Supplier{}, fmt.Errorf("erro db"))

	// Execução - note que passamos &supplier para o service
	_, err := service.Create(ctx, &supplier, 0, nil, nil)

	// Verificações
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao criar fornecedor")
	mockSupplierRepo.AssertExpectations(t)
}
func TestCreateSupplier_HasRelationError(t *testing.T) {
	ctx := context.TODO()

	// Mocks
	mockSupplierRepo := new(mock_supplier.MockSupplierRepository)
	mockRelationService := new(mock_supplier.MockSupplierCategoryRelationService)
	mockAddressService := new(mock_supplier.MockAddressService)
	mockContactService := new(mock_supplier.MockContactService)
	mockCategoryService := new(mock_supplier.MockSupplierCategoryService)

	// Serviço com todos os mocks
	service := &supplierService{
		repo:            mockSupplierRepo,
		relationService: mockRelationService,
		addressService:  mockAddressService,
		contactService:  mockContactService,
		categoryService: mockCategoryService,
	}

	// Dados simulados - agora usando o valor (não ponteiro) para o mock
	supplier := models.Supplier{Name: "Fornecedor A"}
	supplierWithID := supplier
	supplierWithID.ID = 42

	categoryID := int64(2)

	// Configuração dos mocks - note que o mock espera models.Supplier (não ponteiro)
	mockSupplierRepo.On("Create", ctx, supplier).Return(supplierWithID, nil)
	mockRelationService.On("HasRelation", ctx, supplierWithID.ID, categoryID).Return(false, fmt.Errorf("erro ao verificar"))

	// Execução - passamos &supplier para o service (como parece ser a assinatura esperada)
	_, err := service.Create(ctx, &supplier, categoryID, nil, nil)

	// Verificações
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao verificar existência da relação")
	mockSupplierRepo.AssertExpectations(t)
	mockRelationService.AssertExpectations(t)
}

func TestCreateSupplier_RelationAlreadyExists(t *testing.T) {
	ctx := context.TODO()

	// Mocks
	mockSupplierRepo := new(mock_supplier.MockSupplierRepository)
	mockRelationService := new(mock_supplier.MockSupplierCategoryRelationService)
	mockAddressService := new(mock_supplier.MockAddressService)
	mockContactService := new(mock_supplier.MockContactService)
	mockCategoryService := new(mock_supplier.MockSupplierCategoryService)

	// Serviço com todos os mocks
	service := &supplierService{
		repo:            mockSupplierRepo,
		relationService: mockRelationService,
		addressService:  mockAddressService,
		contactService:  mockContactService,
		categoryService: mockCategoryService,
	}

	// Dados simulados
	inputSupplier := models.Supplier{Name: "Fornecedor A"} // Agora é um valor, não ponteiro
	createdSupplier := models.Supplier{
		ID:   10,
		Name: "Fornecedor A",
	}
	categoryID := int64(1)

	// Configuração dos mocks
	mockSupplierRepo.On("Create", ctx, inputSupplier).Return(createdSupplier, nil)
	mockRelationService.On("HasRelation", ctx, createdSupplier.ID, categoryID).Return(true, nil)

	// Execução - note que agora passamos o endereço do inputSupplier
	_, err := service.Create(ctx, &inputSupplier, categoryID, nil, nil)

	// Verificações
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "relação fornecedor-categoria já existe")

	// Verificar se os mocks foram chamados como esperado
	mockSupplierRepo.AssertExpectations(t)
	mockRelationService.AssertExpectations(t)
}

func TestCreateSupplier_RelationCreateError(t *testing.T) {
	ctx := context.TODO()

	// Mocks
	mockSupplierRepo := new(mock_supplier.MockSupplierRepository)
	mockRelationService := new(mock_supplier.MockSupplierCategoryRelationService)
	mockAddressService := new(mock_supplier.MockAddressService)
	mockContactService := new(mock_supplier.MockContactService)
	mockCategoryService := new(mock_supplier.MockSupplierCategoryService)

	// Serviço com todos os mocks
	service := &supplierService{
		repo:            mockSupplierRepo,
		relationService: mockRelationService,
		addressService:  mockAddressService,
		contactService:  mockContactService,
		categoryService: mockCategoryService,
	}

	// Dados simulados
	supplierInput := models.Supplier{Name: "Fornecedor B"} // Agora é um valor
	supplierPtr := &supplierInput                          // Ponteiro para passar para o serviço
	categoryID := int64(3)
	createdSupplier := models.Supplier{
		ID:   55,
		Name: "Fornecedor B",
	}

	// Configuração dos mocks
	mockSupplierRepo.On("Create", ctx, supplierInput).Return(createdSupplier, nil)
	mockRelationService.On("HasRelation", ctx, createdSupplier.ID, categoryID).Return(false, nil)
	mockRelationService.On("Create", ctx, createdSupplier.ID, categoryID).
		Return(&supplier_category_relations.SupplierCategoryRelations{}, fmt.Errorf("falha ao criar"))

	// Execução - passamos o ponteiro para o serviço, mas o mock espera o valor
	_, err := service.Create(ctx, supplierPtr, categoryID, nil, nil)

	// Verificações
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao criar relação fornecedor-categoria")

	// Verificar se os mocks foram chamados como esperado
	mockSupplierRepo.AssertExpectations(t)
	mockRelationService.AssertExpectations(t)
}
func TestCreateSupplier_AddressError(t *testing.T) {
	ctx := context.TODO()

	// 1. Criar todos os mocks necessários
	mockSupplierRepo := new(mock_supplier.MockSupplierRepository)
	mockRelationService := new(mock_supplier.MockSupplierCategoryRelationService)
	mockAddressService := new(mock_supplier.MockAddressService)
	mockContactService := new(mock_supplier.MockContactService)
	mockCategoryService := new(mock_supplier.MockSupplierCategoryService)

	// 2. Configurar o serviço com todos os mocks
	service := &supplierService{
		repo:            mockSupplierRepo,
		relationService: mockRelationService,
		addressService:  mockAddressService,
		contactService:  mockContactService,
		categoryService: mockCategoryService,
	}

	// 3. Dados de teste
	supplierInput := models.Supplier{Name: "Fornecedor A"} // Valor para o mock
	supplierPtr := &supplierInput                          // Ponteiro para o serviço
	categoryID := int64(1)
	supplierID := int64(99)
	createdSupplier := models.Supplier{
		ID:   supplierID,
		Name: "Fornecedor A",
	}
	address := &models_address.Address{Street: "X"}

	// 4. Configurar os mocks
	mockSupplierRepo.On("Create", ctx, supplierInput).Return(createdSupplier, nil)
	mockRelationService.On("HasRelation", ctx, supplierID, categoryID).Return(false, nil)
	mockRelationService.On("Create", ctx, supplierID, categoryID).Return(&supplier_category_relations.SupplierCategoryRelations{}, nil)

	// IMPORTANTE: Configurar para NÃO retornar erro na busca da categoria
	mockCategoryService.On("GetByID", ctx, categoryID).Return(nil, nil)

	// Configurar o erro no endereço
	mockAddressService.On("Create", ctx, mock.Anything).Return(models_address.Address{}, fmt.Errorf("erro no endereço"))

	// 5. Executar o teste
	_, err := service.Create(ctx, supplierPtr, categoryID, address, nil)

	// 6. Verificações
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao criar endereço")

	// 7. Verificar se os mocks foram chamados como esperado
	mockSupplierRepo.AssertExpectations(t)
	mockRelationService.AssertExpectations(t)
	mockCategoryService.AssertExpectations(t)
	mockAddressService.AssertExpectations(t)
}

func TestCreateSupplier_ContactError(t *testing.T) {
	ctx := context.TODO()

	// 1. Criar todos os mocks necessários
	mockSupplierRepo := new(mock_supplier.MockSupplierRepository)
	mockRelationService := new(mock_supplier.MockSupplierCategoryRelationService)
	mockAddressService := new(mock_supplier.MockAddressService)
	mockContactService := new(mock_supplier.MockContactService)
	mockCategoryService := new(mock_supplier.MockSupplierCategoryService)

	// 2. Configurar o serviço com todos os mocks
	service := &supplierService{
		repo:            mockSupplierRepo,
		relationService: mockRelationService,
		addressService:  mockAddressService,
		contactService:  mockContactService,
		categoryService: mockCategoryService,
	}

	// 3. Dados de teste
	supplierInput := models.Supplier{Name: "Fornecedor C"} // Valor para o mock
	supplierPtr := &supplierInput                          // Ponteiro para o serviço
	categoryID := int64(2)
	supplierID := int64(110)
	createdSupplier := models.Supplier{
		ID:   supplierID,
		Name: "Fornecedor C",
	}
	address := &models_address.Address{Street: "Rua C"}
	contact := &models_contact.Contact{ContactName: "Contato X", Email: "x@example.com"}

	// 4. Configurar os mocks
	mockSupplierRepo.On("Create", ctx, supplierInput).Return(createdSupplier, nil)
	mockRelationService.On("HasRelation", ctx, supplierID, categoryID).Return(false, nil)
	mockRelationService.On("Create", ctx, supplierID, categoryID).Return(&supplier_category_relations.SupplierCategoryRelations{}, nil)

	// IMPORTANTE: Configurar para NÃO retornar erro na busca da categoria
	mockCategoryService.On("GetByID", ctx, categoryID).Return(nil, nil)

	// Configurar sucesso no endereço
	mockAddressService.On("Create", ctx, mock.Anything).Return(models_address.Address{}, nil)

	// Configurar o erro no contato
	mockContactService.On("Create", ctx, mock.Anything).Return(models_contact.Contact{}, fmt.Errorf("erro ao criar contato"))

	// 5. Executar o teste
	_, err := service.Create(ctx, supplierPtr, categoryID, address, contact)

	// 6. Verificações
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao criar contato")

	// 7. Verificar se os mocks foram chamados como esperado
	mockSupplierRepo.AssertExpectations(t)
	mockRelationService.AssertExpectations(t)
	mockCategoryService.AssertExpectations(t)
	mockAddressService.AssertExpectations(t)
	mockContactService.AssertExpectations(t)
}
