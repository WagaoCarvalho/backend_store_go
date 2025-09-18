package repositories

import (
	"context"

	"github.com/stretchr/testify/mock"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client_credit"
)

type MockClientCreditRepository struct {
	mock.Mock
}

func (m *MockClientCreditRepository) Create(ctx context.Context, credit *models.ClientCredit) (*models.ClientCredit, error) {
	args := m.Called(ctx, credit)
	var result *models.ClientCredit
	if args.Get(0) != nil {
		result = args.Get(0).(*models.ClientCredit)
	}
	return result, args.Error(1)
}

func (m *MockClientCreditRepository) GetByID(ctx context.Context, id int64) (*models.ClientCredit, error) {
	args := m.Called(ctx, id)
	var result *models.ClientCredit
	if args.Get(0) != nil {
		result = args.Get(0).(*models.ClientCredit)
	}
	return result, args.Error(1)
}

func (m *MockClientCreditRepository) GetByClientID(ctx context.Context, clientID int64) (*models.ClientCredit, error) {
	args := m.Called(ctx, clientID)
	var result *models.ClientCredit
	if args.Get(0) != nil {
		result = args.Get(0).(*models.ClientCredit)
	}
	return result, args.Error(1)
}

func (m *MockClientCreditRepository) Update(ctx context.Context, credit *models.ClientCredit) error {
	args := m.Called(ctx, credit)
	return args.Error(0)
}

func (m *MockClientCreditRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
