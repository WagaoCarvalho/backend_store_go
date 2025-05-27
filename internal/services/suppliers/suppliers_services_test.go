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
	repository "github.com/WagaoCarvalho/backend_store_go/internal/repositories/suppliers"
	mock_supplier "github.com/WagaoCarvalho/backend_store_go/internal/services/suppliers/supplier_service_mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewSupplierService(t *testing.T) {
	mockRepo := new(mock_supplier.MockSupplierRepository)
	mockRelation := new(mock_supplier.MockSupplierCategoryRelationService)
	mockAddress := new(mock_supplier.MockAddressService)
	mockContact := new(mock_supplier.MockContactService)
	mockCategory := new(mock_supplier.MockSupplierCategoryService)

	service := NewSupplierService(
		mockRepo,
		mockRelation,
		mockAddress,
		mockContact,
		mockCategory,
	)

	assert.NotNil(t, service)
}

func TestCreateSupplierById(t *testing.T) {
	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
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
		supplierWithID.ID = 123
		address := &models_address.Address{Street: "Rua A", City: "Cidade", State: "Estado", PostalCode: "12345"}
		contact := &models_contact.Contact{ContactName: "Contato A", Email: "contato@example.com"}
		categoryID := int64(1)

		mockSupplierRepo.On("Create", ctx, supplier).Return(supplierWithID, nil)
		mockRelationService.On("HasRelation", ctx, supplierWithID.ID, categoryID).Return(false, nil)
		mockRelationService.On("Create", ctx, supplierWithID.ID, categoryID).Return(&supplier_category_relations.SupplierCategoryRelations{}, nil)
		mockAddressService.On("Create", ctx, mock.Anything).Return(address, nil)
		mockContactService.On("Create", ctx, mock.Anything).Return(contact, nil)
		mockCategoryService.On("GetByID", ctx, categoryID).Return(&supplier_category.SupplierCategory{ID: categoryID, Name: "Categoria A"}, nil)

		id, err := service.Create(ctx, &supplier, categoryID, address, contact)

		assert.NoError(t, err)
		assert.Equal(t, supplierWithID.ID, id)
		mockSupplierRepo.AssertExpectations(t)
		mockRelationService.AssertExpectations(t)
		mockAddressService.AssertExpectations(t)
		mockContactService.AssertExpectations(t)
	})

	t.Run("InvalidName", func(t *testing.T) {
		service := &supplierService{}
		_, err := service.Create(context.TODO(), &models.Supplier{Name: ""}, 0, nil, nil)

		assert.Error(t, err)
		assert.Equal(t, "nome do fornecedor é obrigatório", err.Error())
	})

	t.Run("RepoError", func(t *testing.T) {
		mockSupplierRepo := new(mock_supplier.MockSupplierRepository)
		service := &supplierService{repo: mockSupplierRepo}
		supplier := models.Supplier{Name: "Fornecedor A"}
		mockSupplierRepo.On("Create", ctx, supplier).Return(models.Supplier{}, fmt.Errorf("erro db"))

		_, err := service.Create(ctx, &supplier, 0, nil, nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar fornecedor")
		mockSupplierRepo.AssertExpectations(t)
	})

	t.Run("HasRelationError", func(t *testing.T) {
		mockSupplierRepo := new(mock_supplier.MockSupplierRepository)
		mockRelationService := new(mock_supplier.MockSupplierCategoryRelationService)
		service := &supplierService{
			repo:            mockSupplierRepo,
			relationService: mockRelationService,
		}

		supplier := models.Supplier{Name: "Fornecedor A"}
		supplierWithID := supplier
		supplierWithID.ID = 42
		categoryID := int64(2)

		mockSupplierRepo.On("Create", ctx, supplier).Return(supplierWithID, nil)
		mockRelationService.On("HasRelation", ctx, supplierWithID.ID, categoryID).Return(false, fmt.Errorf("erro ao verificar"))

		_, err := service.Create(ctx, &supplier, categoryID, nil, nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao verificar existência da relação")
		mockSupplierRepo.AssertExpectations(t)
		mockRelationService.AssertExpectations(t)
	})

	t.Run("RelationAlreadyExists", func(t *testing.T) {
		mockSupplierRepo := new(mock_supplier.MockSupplierRepository)
		mockRelationService := new(mock_supplier.MockSupplierCategoryRelationService)
		service := &supplierService{
			repo:            mockSupplierRepo,
			relationService: mockRelationService,
		}

		supplier := models.Supplier{Name: "Fornecedor A"}
		supplierWithID := supplier
		supplierWithID.ID = 10
		categoryID := int64(1)

		mockSupplierRepo.On("Create", ctx, supplier).Return(supplierWithID, nil)
		mockRelationService.On("HasRelation", ctx, supplierWithID.ID, categoryID).Return(true, nil)

		_, err := service.Create(ctx, &supplier, categoryID, nil, nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "relação fornecedor-categoria já existe")
		mockSupplierRepo.AssertExpectations(t)
		mockRelationService.AssertExpectations(t)
	})

	t.Run("RelationCreateError", func(t *testing.T) {
		mockSupplierRepo := new(mock_supplier.MockSupplierRepository)
		mockRelationService := new(mock_supplier.MockSupplierCategoryRelationService)
		service := &supplierService{
			repo:            mockSupplierRepo,
			relationService: mockRelationService,
		}

		supplier := models.Supplier{Name: "Fornecedor B"}
		supplierWithID := supplier
		supplierWithID.ID = 55
		categoryID := int64(3)

		mockSupplierRepo.On("Create", ctx, supplier).Return(supplierWithID, nil)
		mockRelationService.On("HasRelation", ctx, supplierWithID.ID, categoryID).Return(false, nil)
		mockRelationService.On("Create", ctx, supplierWithID.ID, categoryID).Return(&supplier_category_relations.SupplierCategoryRelations{}, fmt.Errorf("falha ao criar"))

		_, err := service.Create(ctx, &supplier, categoryID, nil, nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar relação fornecedor-categoria")
		mockSupplierRepo.AssertExpectations(t)
		mockRelationService.AssertExpectations(t)
	})

	t.Run("AddressError", func(t *testing.T) {
		mockSupplierRepo := new(mock_supplier.MockSupplierRepository)
		mockRelationService := new(mock_supplier.MockSupplierCategoryRelationService)
		mockAddressService := new(mock_supplier.MockAddressService)
		mockCategoryService := new(mock_supplier.MockSupplierCategoryService)

		service := &supplierService{
			repo:            mockSupplierRepo,
			relationService: mockRelationService,
			addressService:  mockAddressService,
			categoryService: mockCategoryService,
		}

		supplier := models.Supplier{Name: "Fornecedor A"}
		supplierWithID := supplier
		supplierWithID.ID = 99
		categoryID := int64(1)
		address := &models_address.Address{Street: "X"}

		mockSupplierRepo.On("Create", ctx, supplier).Return(supplierWithID, nil)
		mockRelationService.On("HasRelation", ctx, supplierWithID.ID, categoryID).Return(false, nil)
		mockRelationService.On("Create", ctx, supplierWithID.ID, categoryID).Return(&supplier_category_relations.SupplierCategoryRelations{}, nil)
		mockCategoryService.On("GetByID", ctx, categoryID).Return(nil, nil)
		mockAddressService.On("Create", ctx, mock.Anything).Return(&models_address.Address{}, fmt.Errorf("erro ao criar endereço"))

		_, err := service.Create(ctx, &supplier, categoryID, address, nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar endereço")
		mockSupplierRepo.AssertExpectations(t)
		mockRelationService.AssertExpectations(t)
		mockCategoryService.AssertExpectations(t)
		mockAddressService.AssertExpectations(t)
	})

	t.Run("CategoryFetchError", func(t *testing.T) {
		mockSupplierRepo := new(mock_supplier.MockSupplierRepository)
		mockRelationService := new(mock_supplier.MockSupplierCategoryRelationService)
		mockCategoryService := new(mock_supplier.MockSupplierCategoryService)

		service := &supplierService{
			repo:            mockSupplierRepo,
			relationService: mockRelationService,
			categoryService: mockCategoryService,
		}

		supplier := models.Supplier{Name: "Fornecedor Y"}
		supplierWithID := supplier
		supplierWithID.ID = 200
		categoryID := int64(5)

		mockSupplierRepo.On("Create", ctx, supplier).Return(supplierWithID, nil)
		mockRelationService.On("HasRelation", ctx, supplierWithID.ID, categoryID).Return(false, nil)
		mockRelationService.On("Create", ctx, supplierWithID.ID, categoryID).Return(&supplier_category_relations.SupplierCategoryRelations{}, nil)
		mockCategoryService.On("GetByID", ctx, categoryID).Return(nil, fmt.Errorf("erro ao buscar categoria"))

		_, err := service.Create(ctx, &supplier, categoryID, nil, nil)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao buscar categoria")
		mockSupplierRepo.AssertExpectations(t)
		mockRelationService.AssertExpectations(t)
		mockCategoryService.AssertExpectations(t)
	})

	t.Run("ContactError", func(t *testing.T) {
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

		supplier := models.Supplier{Name: "Fornecedor C"}
		supplierWithID := supplier
		supplierWithID.ID = 110
		categoryID := int64(2)
		address := &models_address.Address{Street: "Rua C"}
		contact := &models_contact.Contact{ContactName: "Contato X", Email: "x@example.com"}

		mockSupplierRepo.On("Create", ctx, supplier).Return(supplierWithID, nil)
		mockRelationService.On("HasRelation", ctx, supplierWithID.ID, categoryID).Return(false, nil)
		mockRelationService.On("Create", ctx, supplierWithID.ID, categoryID).Return(&supplier_category_relations.SupplierCategoryRelations{}, nil)
		mockCategoryService.On("GetByID", ctx, categoryID).Return(nil, nil)
		mockAddressService.On("Create", ctx, mock.Anything).Return(&models_address.Address{}, nil)
		mockContactService.On("Create", ctx, mock.Anything).Return(&models_contact.Contact{}, fmt.Errorf("erro ao criar contato"))

		_, err := service.Create(ctx, &supplier, categoryID, address, contact)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar contato")
		mockSupplierRepo.AssertExpectations(t)
		mockRelationService.AssertExpectations(t)
		mockCategoryService.AssertExpectations(t)
		mockAddressService.AssertExpectations(t)
		mockContactService.AssertExpectations(t)
	})
}

func TestGetSupplierByID(t *testing.T) {
	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplierRepository)

		service := &supplierService{
			repo: mockRepo,
		}

		expectedSupplier := &models.Supplier{
			ID:   1,
			Name: "Fornecedor Teste",
		}

		mockRepo.On("GetByID", ctx, int64(1)).Return(expectedSupplier, nil)

		result, err := service.GetByID(ctx, 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedSupplier, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplierRepository)

		service := &supplierService{
			repo: mockRepo,
		}

		mockRepo.On("GetByID", ctx, int64(999)).Return((*models.Supplier)(nil), nil)

		result, err := service.GetByID(ctx, 999)

		assert.NoError(t, err)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplierRepository)

		service := &supplierService{
			repo: mockRepo,
		}

		mockRepo.On("GetByID", ctx, int64(2)).Return((*models.Supplier)(nil), fmt.Errorf("erro no banco"))

		result, err := service.GetByID(ctx, 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "erro no banco")
		mockRepo.AssertExpectations(t)
	})
}

func TestGetAllSuppliers(t *testing.T) {
	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplierRepository)

		service := &supplierService{
			repo: mockRepo,
		}

		expectedSuppliers := []*models.Supplier{
			{ID: 1, Name: "Fornecedor A"},
			{ID: 2, Name: "Fornecedor B"},
		}

		mockRepo.On("GetAll", ctx).Return(expectedSuppliers, nil)

		result, err := service.GetAll(ctx)

		assert.NoError(t, err)
		assert.Equal(t, expectedSuppliers, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("EmptyList", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplierRepository)

		service := &supplierService{
			repo: mockRepo,
		}

		mockRepo.On("GetAll", ctx).Return([]*models.Supplier{}, nil)

		result, err := service.GetAll(ctx)

		assert.NoError(t, err)
		assert.Empty(t, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplierRepository)

		service := &supplierService{
			repo: mockRepo,
		}

		mockRepo.On("GetAll", ctx).Return(([]*models.Supplier)(nil), fmt.Errorf("erro no banco"))

		result, err := service.GetAll(ctx)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "erro no banco")
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateSupplier(t *testing.T) {
	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplierRepository)
		service := &supplierService{repo: mockRepo}

		supplier := &models.Supplier{
			ID:      1,
			Name:    "Fornecedor Atualizado",
			Version: 1,
		}

		mockRepo.On("Update", ctx, supplier).Return(nil)

		err := service.Update(ctx, supplier)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("MissingName", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplierRepository)
		service := &supplierService{repo: mockRepo}

		supplier := &models.Supplier{
			ID:      2,
			Name:    "",
			Version: 1,
		}

		err := service.Update(ctx, supplier)

		assert.ErrorIs(t, err, ErrSupplierNameRequired)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("MissingVersion", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplierRepository)
		service := &supplierService{repo: mockRepo}

		supplier := &models.Supplier{
			ID:   3,
			Name: "Fornecedor Sem Versão",
			// Version ausente ou zero
		}

		err := service.Update(ctx, supplier)

		assert.ErrorIs(t, err, ErrSupplierVersionRequired)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplierRepository)
		service := &supplierService{repo: mockRepo}

		supplier := &models.Supplier{
			ID:      0,
			Name:    "Fornecedor Inválido",
			Version: 1,
		}

		err := service.Update(ctx, supplier)

		assert.ErrorIs(t, err, ErrInvalidSupplierID)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("VersionConflict", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplierRepository)
		service := &supplierService{repo: mockRepo}

		supplier := &models.Supplier{
			ID:      4,
			Name:    "Fornecedor Conflito",
			Version: 2,
		}

		mockRepo.On("Update", ctx, supplier).Return(repository.ErrVersionConflict)

		err := service.Update(ctx, supplier)

		assert.ErrorIs(t, err, ErrSupplierVersionConflict)
		mockRepo.AssertExpectations(t)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplierRepository)
		service := &supplierService{repo: mockRepo}

		supplier := &models.Supplier{
			ID:      3,
			Name:    "Fornecedor com Erro",
			Version: 1,
		}

		mockRepo.On("Update", ctx, supplier).Return(fmt.Errorf("erro ao atualizar"))

		err := service.Update(ctx, supplier)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao atualizar fornecedor")
		mockRepo.AssertExpectations(t)
	})

	t.Run("SupplierNotFound", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplierRepository)
		service := &supplierService{repo: mockRepo}

		supplier := &models.Supplier{
			ID:      10,
			Name:    "Fornecedor Inexistente",
			Version: 1,
		}

		mockRepo.On("Update", ctx, supplier).Return(repository.ErrSupplierNotFound)

		err := service.Update(ctx, supplier)

		assert.ErrorIs(t, err, ErrSupplierNotFound)
		mockRepo.AssertExpectations(t)
	})

}

func TestDeleteSupplier(t *testing.T) {
	ctx := context.TODO()

	t.Run("Success", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplierRepository)

		service := &supplierService{
			repo: mockRepo,
		}

		id := int64(1)

		mockRepo.On("Delete", ctx, id).Return(nil)

		err := service.Delete(ctx, id)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplierRepository)

		service := &supplierService{
			repo: mockRepo,
		}

		id := int64(0)

		err := service.Delete(ctx, id)

		assert.ErrorIs(t, err, ErrInvalidSupplierIDForDeletion)
		mockRepo.AssertNotCalled(t, "Delete", mock.Anything, mock.Anything)
	})

	t.Run("RepositoryError", func(t *testing.T) {
		mockRepo := new(mock_supplier.MockSupplierRepository)

		service := &supplierService{
			repo: mockRepo,
		}

		id := int64(2)

		mockRepo.On("Delete", ctx, id).Return(fmt.Errorf("erro ao deletar"))

		err := service.Delete(ctx, id)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao deletar")
		mockRepo.AssertExpectations(t)
	})
}
