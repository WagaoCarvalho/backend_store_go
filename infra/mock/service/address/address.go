package services

import (
	"context"

	dtoAddress "github.com/WagaoCarvalho/backend_store_go/internal/dto/address"
	"github.com/stretchr/testify/mock"
)

type MockAddressService struct {
	mock.Mock
}

func (m *MockAddressService) Create(ctx context.Context, dto *dtoAddress.AddressDTO) (*dtoAddress.AddressDTO, error) {
	args := m.Called(ctx, dto)

	if result := args.Get(0); result != nil {
		return result.(*dtoAddress.AddressDTO), args.Error(1)
	}

	return nil, args.Error(1)
}

func (m *MockAddressService) GetByID(ctx context.Context, id int64) (*dtoAddress.AddressDTO, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dtoAddress.AddressDTO), args.Error(1)
}

func (m *MockAddressService) GetByUserID(ctx context.Context, userID int64) ([]*dtoAddress.AddressDTO, error) {
	args := m.Called(ctx, userID)

	// Precaução para evitar panic se o valor retornado for nil
	if dtos, ok := args.Get(0).([]*dtoAddress.AddressDTO); ok {
		return dtos, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAddressService) GetByClientID(ctx context.Context, id int64) ([]*dtoAddress.AddressDTO, error) {
	args := m.Called(ctx, id)

	if dtos, ok := args.Get(0).([]*dtoAddress.AddressDTO); ok {
		return dtos, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAddressService) GetBySupplierID(ctx context.Context, id int64) ([]*dtoAddress.AddressDTO, error) {
	args := m.Called(ctx, id)

	if dtos, ok := args.Get(0).([]*dtoAddress.AddressDTO); ok {
		return dtos, args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockAddressService) Update(ctx context.Context, address *dtoAddress.AddressDTO) error {
	args := m.Called(ctx, address)
	return args.Error(0)
}

func (m *MockAddressService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
