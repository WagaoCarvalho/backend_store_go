package repositories

import (
	"context"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/mock"
)

type MockContactRepository struct {
	mock.Mock
}

func (m *MockContactRepository) Create(ctx context.Context, c *models.Contact) (*models.Contact, error) {
	args := m.Called(ctx, c)
	return args.Get(0).(*models.Contact), args.Error(1)
}

func (m *MockContactRepository) CreateTx(ctx context.Context, tx pgx.Tx, contact *models.Contact) (*models.Contact, error) {
	args := m.Called(ctx, tx, contact) // Isso aciona o .Run() definido no teste

	var result *models.Contact
	if args.Get(0) != nil {
		result = args.Get(0).(*models.Contact)
	}

	return result, args.Error(1)
}

func (m *MockContactRepository) GetByID(ctx context.Context, id int64) (*models.Contact, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Contact), args.Error(1)
}

func (m *MockContactRepository) GetByUserID(ctx context.Context, userID int64) ([]*models.Contact, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*models.Contact), args.Error(1)
}

func (m *MockContactRepository) GetByClientID(ctx context.Context, clientID int64) ([]*models.Contact, error) {
	args := m.Called(ctx, clientID)
	return args.Get(0).([]*models.Contact), args.Error(1)
}

func (m *MockContactRepository) GetBySupplierID(ctx context.Context, supplierID int64) ([]*models.Contact, error) {
	args := m.Called(ctx, supplierID)
	return args.Get(0).([]*models.Contact), args.Error(1)
}

func (m *MockContactRepository) Update(ctx context.Context, contact *models.Contact) error {
	args := m.Called(ctx, contact)
	return args.Error(0)
}

func (m *MockContactRepository) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}
