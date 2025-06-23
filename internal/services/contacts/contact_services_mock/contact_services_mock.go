package services

import (
	"context"

	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	"github.com/stretchr/testify/mock"
)

type MockContactService struct {
	mock.Mock
}

func (m *MockContactService) Create(ctx context.Context, c *models_contact.Contact) (*models_contact.Contact, error) {
	args := m.Called(ctx, c)
	return args.Get(0).(*models_contact.Contact), args.Error(1)
}

func (m *MockContactService) GetByID(ctx context.Context, id int64) (*models_contact.Contact, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models_contact.Contact), args.Error(1)
}

func (m *MockContactService) GetVersionByID(ctx context.Context, id int64) (int, error) {
	args := m.Called(ctx, id)
	return args.Int(0), args.Error(1)
}

func (m *MockContactService) GetByUserID(ctx context.Context, userID int64) ([]*models_contact.Contact, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*models_contact.Contact), args.Error(1)
}

func (m *MockContactService) GetByClientID(ctx context.Context, clientID int64) ([]*models_contact.Contact, error) {
	args := m.Called(ctx, clientID)
	return args.Get(0).([]*models_contact.Contact), args.Error(1)
}

func (m *MockContactService) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models_contact.Contact, error) {
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
