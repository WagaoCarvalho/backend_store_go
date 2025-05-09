package handlers

import (
	"context"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	models_supplier "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier"
	models_supplier_realiations "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_category_relations"
	"github.com/stretchr/testify/mock"
)

type MockSupplierService struct {
	mock.Mock
}

type MockSupplierRepo struct {
	mock.Mock
}

type MockSupplierCategoryRelationService struct {
	mock.Mock
}

type MockAddressService struct {
	mock.Mock
}

type MockContactService struct {
	mock.Mock
}

// MockSupplierRepo
func (m *MockSupplierRepo) Create(ctx context.Context, s *models_supplier.Supplier) (int64, error) {
	args := m.Called(ctx, s)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSupplierCategoryRelationService) HasRelation(ctx context.Context, supplierID, categoryID int64) (bool, error) {
	args := m.Called(ctx, supplierID, categoryID)
	return args.Bool(0), args.Error(1)
}

func (m *MockSupplierRepo) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockSupplierRepo) GetByID(ctx context.Context, id int64) (*models_supplier.Supplier, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models_supplier.Supplier), args.Error(1)
}

func (m *MockSupplierRepo) GetAll(ctx context.Context) ([]*models_supplier.Supplier, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models_supplier.Supplier), args.Error(1)
}

func (m *MockSupplierRepo) Update(ctx context.Context, s *models_supplier.Supplier) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

// MockSupplierCategoryRelationService
func (m *MockSupplierCategoryRelationService) Create(ctx context.Context, supplierID, categoryID int64) (*models_supplier_realiations.SupplierCategoryRelations, error) {
	args := m.Called(ctx, supplierID, categoryID)
	return args.Get(0).(*models_supplier_realiations.SupplierCategoryRelations), args.Error(1)
}

func (m *MockSupplierCategoryRelationService) Delete(ctx context.Context, supplierID, categoryID int64) error {
	args := m.Called(ctx, supplierID, categoryID)
	return args.Error(0)
}

func (m *MockSupplierCategoryRelationService) DeleteAll(ctx context.Context, supplierID int64) error {
	args := m.Called(ctx, supplierID)
	return args.Error(0)
}

func (m *MockSupplierCategoryRelationService) GetBySupplier(ctx context.Context, supplierID int64) ([]*models_supplier_realiations.SupplierCategoryRelations, error) {
	args := m.Called(ctx, supplierID)
	return args.Get(0).([]*models_supplier_realiations.SupplierCategoryRelations), args.Error(1)
}

func (m *MockSupplierCategoryRelationService) GetByCategory(ctx context.Context, categoryID int64) ([]*models_supplier_realiations.SupplierCategoryRelations, error) {
	args := m.Called(ctx, categoryID)
	return args.Get(0).([]*models_supplier_realiations.SupplierCategoryRelations), args.Error(1)
}

// MockSupplierService
func (m *MockSupplierService) Create(ctx context.Context, s *models_supplier.Supplier, categoryID int64, address *models_address.Address, contact *models_contact.Contact) (int64, error) {
	args := m.Called(ctx, s, categoryID, address, contact)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSupplierService) GetByID(ctx context.Context, id int64) (*models_supplier.Supplier, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models_supplier.Supplier), args.Error(1)
}

func (m *MockSupplierService) GetAll(ctx context.Context) ([]*models_supplier.Supplier, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models_supplier.Supplier), args.Error(1)
}

func (m *MockSupplierService) Update(ctx context.Context, s *models_supplier.Supplier) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *MockSupplierService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockAddressService
func (m *MockAddressService) Create(ctx context.Context, address models_address.Address) (models_address.Address, error) {
	args := m.Called(ctx, address)
	return args.Get(0).(models_address.Address), args.Error(1)
}

func (m *MockAddressService) GetByID(ctx context.Context, id int) (models_address.Address, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models_address.Address), args.Error(1)
}

func (m *MockAddressService) Update(ctx context.Context, address models_address.Address) error {
	args := m.Called(ctx, address)
	return args.Error(0)
}

func (m *MockAddressService) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockContactService
func (m *MockContactService) Create(ctx context.Context, contact *models_contact.Contact) error {
	args := m.Called(ctx, contact)
	return args.Error(0)
}

func (m *MockContactService) GetByID(ctx context.Context, id int64) (*models_contact.Contact, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models_contact.Contact), args.Error(1)
}

func (m *MockContactService) GetByUser(ctx context.Context, userID int64) ([]*models_contact.Contact, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*models_contact.Contact), args.Error(1)
}

func (m *MockContactService) GetByClient(ctx context.Context, clientID int64) ([]*models_contact.Contact, error) {
	args := m.Called(ctx, clientID)
	return args.Get(0).([]*models_contact.Contact), args.Error(1)
}

func (m *MockContactService) GetBySupplier(ctx context.Context, supplierID int64) ([]*models_contact.Contact, error) {
	args := m.Called(ctx, supplierID)
	return args.Get(0).([]*models_contact.Contact), args.Error(1)
}

func (m *MockContactService) Update(ctx context.Context, contact *models_contact.Contact) error {
	args := m.Called(ctx, contact)
	return args.Error(0)
}

func (m *MockContactService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
