package services

import (
	"context"
	"fmt"
	"testing"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier"
	supplier_category_relations "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_category_relations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSupplierService struct {
	CreateFn  func(ctx context.Context, supplier *models.Supplier, categoryID int64, address *models_address.Address, contact *models_contact.Contact) (int64, error)
	GetByIDFn func(ctx context.Context, id int64) (*models.Supplier, error)
	GetAllFn  func(ctx context.Context) ([]*models.Supplier, error)
	UpdateFn  func(ctx context.Context, supplier *models.Supplier) error
	DeleteFn  func(ctx context.Context, id int64) error
}

func (m *MockSupplierService) Create(ctx context.Context, supplier *models.Supplier, categoryID int64, address *models_address.Address, contact *models_contact.Contact) (int64, error) {
	return m.CreateFn(ctx, supplier, categoryID, address, contact)
}

func (m *MockSupplierService) GetByID(ctx context.Context, id int64) (*models.Supplier, error) {
	return m.GetByIDFn(ctx, id)
}

func (m *MockSupplierService) GetAll(ctx context.Context) ([]*models.Supplier, error) {
	return m.GetAllFn(ctx)
}

func (m *MockSupplierService) Update(ctx context.Context, supplier *models.Supplier) error {
	return m.UpdateFn(ctx, supplier)
}

func (m *MockSupplierService) Delete(ctx context.Context, id int64) error {
	return m.DeleteFn(ctx, id)
}

type mockRepo struct{ mock.Mock }

func (m *mockRepo) Create(ctx context.Context, supplier *models.Supplier) (int64, error) {
	args := m.Called(ctx, supplier)
	return args.Get(0).(int64), args.Error(1)
}

func (m *mockRepo) GetByID(ctx context.Context, id int64) (*models.Supplier, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Supplier), args.Error(1)
}

func (m *mockRepo) GetAll(ctx context.Context) ([]*models.Supplier, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Supplier), args.Error(1)
}

func (m *mockRepo) Update(ctx context.Context, supplier *models.Supplier) error {
	args := m.Called(ctx, supplier)
	return args.Error(0)
}

func (m *mockRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockRelationService struct {
	mock.Mock
}

func (m *mockRelationService) Create(ctx context.Context, supplierID, categoryID int64) (*supplier_category_relations.SupplierCategoryRelations, error) {
	args := m.Called(ctx, supplierID, categoryID)
	if args.Get(0) != nil {
		return args.Get(0).(*supplier_category_relations.SupplierCategoryRelations), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRelationService) GetBySupplier(ctx context.Context, supplierID int64) ([]*supplier_category_relations.SupplierCategoryRelations, error) {
	args := m.Called(ctx, supplierID)
	return args.Get(0).([]*supplier_category_relations.SupplierCategoryRelations), args.Error(1)
}

func (m *mockRelationService) GetByCategory(ctx context.Context, categoryID int64) ([]*supplier_category_relations.SupplierCategoryRelations, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).([]*supplier_category_relations.SupplierCategoryRelations), args.Error(1)
}

func (m *mockRelationService) Delete(ctx context.Context, supplierID, categoryID int64) error {
	args := m.Called(ctx, supplierID, categoryID)
	return args.Error(0)
}

func (m *mockRelationService) DeleteAll(ctx context.Context, supplierID int64) error {
	args := m.Called(ctx, supplierID)
	return args.Error(0)
}

func (m *mockRelationService) HasRelation(ctx context.Context, supplierID, categoryID int64) (bool, error) {
	args := m.Called(ctx, supplierID, categoryID)
	return args.Bool(0), args.Error(1)
}

type mockAddressService struct{ mock.Mock }

func (m *mockAddressService) Create(ctx context.Context, address models_address.Address) (models_address.Address, error) {
	args := m.Called(ctx, address)
	return args.Get(0).(models_address.Address), args.Error(1)
}

func (m *mockAddressService) GetByID(ctx context.Context, id int) (models_address.Address, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models_address.Address), args.Error(1)
}

func (m *mockAddressService) Update(ctx context.Context, address models_address.Address) error {
	args := m.Called(ctx, address)
	return args.Error(0)
}

func (m *mockAddressService) Delete(ctx context.Context, supplierID int) error {
	args := m.Called(ctx, supplierID)
	return args.Error(0)
}

type mockContactService struct{ mock.Mock }

func (m *mockContactService) Create(ctx context.Context, contact *models_contact.Contact) error {
	args := m.Called(ctx, contact)
	return args.Error(0)
}

func (m *mockContactService) GetByID(ctx context.Context, id int64) (*models_contact.Contact, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models_contact.Contact), args.Error(1)
}

func (m *mockContactService) GetByUser(ctx context.Context, userID int64) ([]*models_contact.Contact, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*models_contact.Contact), args.Error(1)
}

func (m *mockContactService) GetByClient(ctx context.Context, clientID int64) ([]*models_contact.Contact, error) {
	args := m.Called(ctx, clientID)
	return args.Get(0).([]*models_contact.Contact), args.Error(1)
}

func (m *mockContactService) GetBySupplier(ctx context.Context, supplierID int64) ([]*models_contact.Contact, error) {
	args := m.Called(ctx, supplierID)
	return args.Get(0).([]*models_contact.Contact), args.Error(1)
}

func (m *mockContactService) Update(ctx context.Context, contact *models_contact.Contact) error {
	args := m.Called(ctx, contact)
	return args.Error(0)
}

func (m *mockContactService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateSupplier_Success(t *testing.T) {
	ctx := context.TODO()
	mockRepo := new(mockRepo)
	mockRelation := new(mockRelationService)
	mockAddr := new(mockAddressService)
	mockContact := new(mockContactService)

	service := &supplierService{
		repo:            mockRepo,
		relationService: mockRelation,
		addressService:  mockAddr,
		contactService:  mockContact,
	}

	supplier := &models.Supplier{Name: "Fornecedor A"}
	address := &models_address.Address{Street: "Rua A", City: "Cidade", State: "Estado", PostalCode: "12345"}
	contact := &models_contact.Contact{ContactName: "Contato A", Email: "contato@example.com"}
	categoryID := int64(1)
	supplierID := int64(123)

	mockRepo.On("Create", ctx, supplier).Return(supplierID, nil)
	mockRelation.On("HasRelation", ctx, supplierID, categoryID).Return(false, nil)
	mockRelation.On("Create", ctx, supplierID, categoryID).Return(&supplier_category_relations.SupplierCategoryRelations{}, nil)
	mockAddr.On("Create", ctx, mock.Anything).Return(*address, nil)
	mockContact.On("Create", ctx, mock.Anything).Return(nil)

	id, err := service.Create(ctx, supplier, categoryID, address, contact)

	assert.NoError(t, err)
	assert.Equal(t, supplierID, id)
	mockRepo.AssertExpectations(t)
	mockRelation.AssertExpectations(t)
	mockAddr.AssertExpectations(t)
	mockContact.AssertExpectations(t)
}

func TestCreateSupplier_InvalidName(t *testing.T) {
	service := &supplierService{}
	_, err := service.Create(context.TODO(), &models.Supplier{Name: ""}, 0, nil, nil)
	assert.Error(t, err)
	assert.Equal(t, "nome do fornecedor é obrigatório", err.Error())
}

func TestCreateSupplier_RepoError(t *testing.T) {
	ctx := context.TODO()
	mockRepo := new(mockRepo)
	service := &supplierService{repo: mockRepo}

	supplier := &models.Supplier{Name: "Fornecedor A"}
	mockRepo.On("Create", ctx, supplier).Return(int64(0), fmt.Errorf("erro db"))

	_, err := service.Create(ctx, supplier, 0, nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao criar fornecedor")
}

func TestCreateSupplier_RelationAlreadyExists(t *testing.T) {
	ctx := context.TODO()
	mockRepo := new(mockRepo)
	mockRelation := new(mockRelationService)

	service := &supplierService{
		repo:            mockRepo,
		relationService: mockRelation,
	}

	supplier := &models.Supplier{Name: "Fornecedor A"}
	categoryID := int64(1)
	supplierID := int64(10)

	mockRepo.On("Create", ctx, supplier).Return(supplierID, nil)
	mockRelation.On("HasRelation", ctx, supplierID, categoryID).Return(true, nil)

	_, err := service.Create(ctx, supplier, categoryID, nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "relação fornecedor-categoria já existe")
}

func TestCreateSupplier_AddressError(t *testing.T) {
	ctx := context.TODO()
	mockRepo := new(mockRepo)
	mockRelation := new(mockRelationService)
	mockAddr := new(mockAddressService)

	service := &supplierService{
		repo:            mockRepo,
		relationService: mockRelation,
		addressService:  mockAddr,
	}

	supplier := &models.Supplier{Name: "Fornecedor A"}
	categoryID := int64(1)
	supplierID := int64(99)
	address := &models_address.Address{Street: "X"}

	mockRepo.On("Create", ctx, supplier).Return(supplierID, nil)
	mockRelation.On("HasRelation", ctx, supplierID, categoryID).Return(false, nil)
	mockRelation.On("Create", ctx, supplierID, categoryID).Return(&supplier_category_relations.SupplierCategoryRelations{}, nil)
	mockAddr.On("Create", ctx, mock.Anything).Return(models_address.Address{}, fmt.Errorf("erro no endereço"))

	_, err := service.Create(ctx, supplier, categoryID, address, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao criar endereço")
}

func TestCreateSupplier_ContactError(t *testing.T) {
	ctx := context.TODO()
	mockRepo := new(mockRepo)
	mockRelation := new(mockRelationService)
	mockAddr := new(mockAddressService)
	mockContact := new(mockContactService)

	service := &supplierService{
		repo:            mockRepo,
		relationService: mockRelation,
		addressService:  mockAddr,
		contactService:  mockContact,
	}

	supplier := &models.Supplier{Name: "Fornecedor A"}
	categoryID := int64(1)
	supplierID := int64(99)
	address := &models_address.Address{Street: "X"}
	contact := &models_contact.Contact{Email: "test@a.com"}

	mockRepo.On("Create", ctx, supplier).Return(supplierID, nil)
	mockRelation.On("HasRelation", ctx, supplierID, categoryID).Return(false, nil)
	mockRelation.On("Create", ctx, supplierID, categoryID).Return(&supplier_category_relations.SupplierCategoryRelations{}, nil)
	mockAddr.On("Create", ctx, mock.Anything).Return(*address, nil)
	mockContact.On("Create", ctx, contact).Return(fmt.Errorf("erro contato"))

	_, err := service.Create(ctx, supplier, categoryID, address, contact)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao criar contato")
}

func TestCreateSupplier_HasRelationError(t *testing.T) {
	ctx := context.TODO()
	mockRepo := new(mockRepo)
	mockRelation := new(mockRelationService)

	service := &supplierService{
		repo:            mockRepo,
		relationService: mockRelation,
	}

	supplier := &models.Supplier{Name: "Fornecedor A"}
	categoryID := int64(2)
	supplierID := int64(42)

	mockRepo.On("Create", ctx, supplier).Return(supplierID, nil)
	mockRelation.On("HasRelation", ctx, supplierID, categoryID).Return(false, fmt.Errorf("erro ao verificar"))

	_, err := service.Create(ctx, supplier, categoryID, nil, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao verificar existência da relação")
}

func TestCreateSupplier_RelationCreateError(t *testing.T) {
	ctx := context.TODO()
	mockRepo := new(mockRepo)
	mockRelation := new(mockRelationService)

	service := &supplierService{
		repo:            mockRepo,
		relationService: mockRelation,
	}

	supplier := &models.Supplier{Name: "Fornecedor B"}
	categoryID := int64(3)
	supplierID := int64(55)

	mockRepo.On("Create", ctx, supplier).Return(supplierID, nil)
	mockRelation.On("HasRelation", ctx, supplierID, categoryID).Return(false, nil)
	mockRelation.On("Create", ctx, supplierID, categoryID).Return(nil, fmt.Errorf("falha ao criar"))

	_, err := service.Create(ctx, supplier, categoryID, nil, nil)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "erro ao criar relação fornecedor-categoria")
}
