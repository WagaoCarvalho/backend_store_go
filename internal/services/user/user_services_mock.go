package services

import (
	"context"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	models_user "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	user_category_relations "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_category_relations"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

type MockUserCategoryRelationRepositories struct {
	mock.Mock
}

type MockAddressRepository struct {
	mock.Mock
}

type MockContactRepository struct {
	mock.Mock
}

// MockUserRepository
func (m *MockUserRepository) Create(
	ctx context.Context,
	user *models_user.User,
) (models_user.User, error) {

	args := m.Called(ctx, user)
	if createdUser, ok := args.Get(0).(models_user.User); ok {
		return createdUser, args.Error(1)
	}
	return models_user.User{}, args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user models_user.User) (models_user.User, error) {
	args := m.Called(ctx, user)
	if updatedUser, ok := args.Get(0).(models_user.User); ok {
		return updatedUser, args.Error(1)
	}
	return models_user.User{}, args.Error(1)
}

func (m *MockUserRepository) Delete(ctx context.Context, uid int64) error {
	args := m.Called(ctx, uid)
	return args.Error(0)
}

func (m *MockUserRepository) GetAll(ctx context.Context) ([]models_user.User, error) {
	args := m.Called(ctx)
	if users, ok := args.Get(0).([]models_user.User); ok {
		return users, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserRepository) GetById(ctx context.Context, uid int64) (models_user.User, error) {
	args := m.Called(ctx, uid)
	if user, ok := args.Get(0).(models_user.User); ok {
		return user, args.Error(1)
	}
	return models_user.User{}, args.Error(1)
}

func (m *MockUserRepository) GetByEmail(ctx context.Context, email string) (models_user.User, error) {
	args := m.Called(ctx, email)
	if user, ok := args.Get(0).(models_user.User); ok {
		return user, args.Error(1)
	}
	return models_user.User{}, args.Error(1)
}

// MockUserCategoryRelationRepositories
func (m *MockUserCategoryRelationRepositories) Create(ctx context.Context, relation *user_category_relations.UserCategoryRelations) (*user_category_relations.UserCategoryRelations, error) {
	args := m.Called(ctx, relation)
	if createdRelation, ok := args.Get(0).(*user_category_relations.UserCategoryRelations); ok {
		return createdRelation, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserCategoryRelationRepositories) GetAll(ctx context.Context, userID int64) ([]user_category_relations.UserCategoryRelations, error) {
	args := m.Called(ctx, userID)
	if relations, ok := args.Get(0).([]user_category_relations.UserCategoryRelations); ok {
		return relations, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserCategoryRelationRepositories) GetRelations(ctx context.Context, categoryID int64) ([]user_category_relations.UserCategoryRelations, error) {
	args := m.Called(ctx, categoryID)
	if relations, ok := args.Get(0).([]user_category_relations.UserCategoryRelations); ok {
		return relations, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserCategoryRelationRepositories) Delete(ctx context.Context, userID, categoryID int64) error {
	args := m.Called(ctx, userID, categoryID)
	return args.Error(0)
}

func (m *MockUserCategoryRelationRepositories) DeleteAll(ctx context.Context, userID int64) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *MockUserCategoryRelationRepositories) GetByCategoryID(ctx context.Context, categoryID int64) ([]user_category_relations.UserCategoryRelations, error) {
	args := m.Called(ctx, categoryID)
	if relations, ok := args.Get(0).([]user_category_relations.UserCategoryRelations); ok {
		return relations, args.Error(1)
	}
	return nil, args.Error(1)
}

// No seu arquivo de mocks (services/mocks_test.go)

// GetByUserID mock para retornar relações de categorias por usuário
func (m *MockUserCategoryRelationRepositories) GetByUserID(ctx context.Context, userID int64) ([]user_category_relations.UserCategoryRelations, error) {
	args := m.Called(ctx, userID)

	// Retorna as relações mockadas ou erro conforme configurado nos testes
	if relations, ok := args.Get(0).([]user_category_relations.UserCategoryRelations); ok {
		return relations, args.Error(1)
	}
	return nil, args.Error(1)
}

// MockAddressRepository
func (m *MockAddressRepository) Create(ctx context.Context, address models_address.Address) (models_address.Address, error) {
	args := m.Called(ctx, address)
	if addr, ok := args.Get(0).(models_address.Address); ok {
		return addr, args.Error(1)
	}
	return models_address.Address{}, args.Error(1)
}

func (m *MockAddressRepository) GetByID(ctx context.Context, id int) (models_address.Address, error) {
	args := m.Called(ctx, id)
	if addr, ok := args.Get(0).(models_address.Address); ok {
		return addr, args.Error(1)
	}
	return models_address.Address{}, args.Error(1)
}

func (m *MockAddressRepository) Update(ctx context.Context, address models_address.Address) error {
	args := m.Called(ctx, address)
	return args.Error(0)
}

func (m *MockAddressRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// MockContactRepository

func (m *MockContactRepository) Create(ctx context.Context, c models_contact.Contact) (models_contact.Contact, error) {
	args := m.Called(ctx, c)

	if contact, ok := args.Get(0).(models_contact.Contact); ok {
		return contact, args.Error(1)
	}

	return models_contact.Contact{}, args.Error(1)
}

func (m *MockContactRepository) GetByID(ctx context.Context, id int64) (*models_contact.Contact, error) {
	args := m.Called(ctx, id)
	if contact, ok := args.Get(0).(*models_contact.Contact); ok {
		return contact, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockContactRepository) GetByUserID(ctx context.Context, userID int64) ([]*models_contact.Contact, error) {
	args := m.Called(ctx, userID)
	if contacts, ok := args.Get(0).([]*models_contact.Contact); ok {
		return contacts, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockContactRepository) GetByClientID(ctx context.Context, clientID int64) ([]*models_contact.Contact, error) {
	args := m.Called(ctx, clientID)
	if contacts, ok := args.Get(0).([]*models_contact.Contact); ok {
		return contacts, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockContactRepository) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models_contact.Contact, error) {
	args := m.Called(ctx, supplierID)
	if contacts, ok := args.Get(0).([]*models_contact.Contact); ok {
		return contacts, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockContactRepository) Update(ctx context.Context, contact *models_contact.Contact) error {
	args := m.Called(ctx, contact)
	return args.Error(0)
}

func (m *MockContactRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
