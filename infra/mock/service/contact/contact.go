package services

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	"github.com/stretchr/testify/mock"
)

type MockContactService struct {
	mock.Mock
}

func (m *MockContactService) Create(ctx context.Context, contact *models.Contact) (*models.Contact, error) {
	args := m.Called(ctx, contact)

	var created *models.Contact
	if args.Get(0) != nil {
		created = args.Get(0).(*models.Contact)
	}

	return created, args.Error(1)
}

func (m *MockContactService) GetByID(ctx context.Context, id int64) (*models.Contact, error) {
	args := m.Called(ctx, id)

	var contact *models.Contact
	if args.Get(0) != nil {
		contact = args.Get(0).(*models.Contact)
	}

	return contact, args.Error(1)
}

func (m *MockContactService) GetByUserID(ctx context.Context, userID int64) ([]*models.Contact, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Contact), args.Error(1)
}

func (m *MockContactService) GetByClientID(ctx context.Context, clientID int64) ([]*models.Contact, error) {
	args := m.Called(ctx, clientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Contact), args.Error(1)
}

func (m *MockContactService) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.Contact, error) {
	args := m.Called(ctx, supplierID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Contact), args.Error(1)
}

func (m *MockContactService) Update(ctx context.Context, contact *models.Contact) error {
	args := m.Called(ctx, contact)
	return args.Error(0)
}

func (m *MockContactService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
