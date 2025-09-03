package services

import (
	"context"

	dtoContact "github.com/WagaoCarvalho/backend_store_go/internal/dto/contact"
	"github.com/stretchr/testify/mock"
)

type MockContactService struct {
	mock.Mock
}

func (m *MockContactService) Create(ctx context.Context, dto *dtoContact.ContactDTO) (*dtoContact.ContactDTO, error) {
	args := m.Called(ctx, dto)

	var contactDTO *dtoContact.ContactDTO
	if args.Get(0) != nil {
		contactDTO = args.Get(0).(*dtoContact.ContactDTO)
	}

	return contactDTO, args.Error(1)
}

func (m *MockContactService) GetByID(ctx context.Context, id int64) (*dtoContact.ContactDTO, error) {
	args := m.Called(ctx, id)

	var contactDTO *dtoContact.ContactDTO
	if args.Get(0) != nil {
		contactDTO = args.Get(0).(*dtoContact.ContactDTO)
	}

	return contactDTO, args.Error(1)
}

func (m *MockContactService) GetByUserID(ctx context.Context, userID int64) ([]*dtoContact.ContactDTO, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*dtoContact.ContactDTO), args.Error(1)
}

func (m *MockContactService) GetByClientID(ctx context.Context, clientID int64) ([]*dtoContact.ContactDTO, error) {
	args := m.Called(ctx, clientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*dtoContact.ContactDTO), args.Error(1)
}

func (m *MockContactService) GetBySupplierID(ctx context.Context, supplierID int64) ([]*dtoContact.ContactDTO, error) {
	args := m.Called(ctx, supplierID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*dtoContact.ContactDTO), args.Error(1)
}

func (m *MockContactService) Update(ctx context.Context, contact *dtoContact.ContactDTO) error {
	args := m.Called(ctx, contact)
	return args.Error(0)
}

func (m *MockContactService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
