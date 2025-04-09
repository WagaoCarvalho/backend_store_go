package services

import (
	"context"
	"errors"
	"testing"
	"time"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/contacts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockContactRepository struct {
	mock.Mock
}

func (m *MockContactRepository) CreateContact(ctx context.Context, contact *models.Contact) error {
	args := m.Called(ctx, contact)
	return args.Error(0)
}

func (m *MockContactRepository) GetContactByID(ctx context.Context, id int64) (*models.Contact, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Contact), args.Error(1)
}

func (m *MockContactRepository) GetContactByUserID(ctx context.Context, userID int64) ([]*models.Contact, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*models.Contact), args.Error(1)
}

func (m *MockContactRepository) GetContactByClientID(ctx context.Context, clientID int64) ([]*models.Contact, error) {
	args := m.Called(ctx, clientID)
	return args.Get(0).([]*models.Contact), args.Error(1)
}

func (m *MockContactRepository) GetContactBySupplierID(ctx context.Context, supplierID int64) ([]*models.Contact, error) {
	args := m.Called(ctx, supplierID)
	return args.Get(0).([]*models.Contact), args.Error(1)
}

func (m *MockContactRepository) Updatecontac(ctx context.Context, contact *models.Contact) error {
	args := m.Called(ctx, contact)
	return args.Error(0)
}

func (m *MockContactRepository) Deletecontact(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateContact(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockContactRepository)
		service := NewContactService(mockRepo)

		userID := int64(1)
		newContact := &models.Contact{
			UserID:          &userID,
			ContactName:     "Test User",
			ContactPosition: "Developer",
			Email:           "test@example.com",
		}

		mockRepo.On("CreateContact", mock.Anything, newContact).Return(nil)

		err := service.CreateContact(context.Background(), newContact)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation error - missing name", func(t *testing.T) {
		mockRepo := new(MockContactRepository)
		service := NewContactService(mockRepo)

		userID := int64(1)
		invalidContact := &models.Contact{
			UserID: &userID,
			Email:  "test@example.com",
		}

		err := service.CreateContact(context.Background(), invalidContact)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nome do contato é obrigatório")
		mockRepo.AssertNotCalled(t, "CreateContact")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(MockContactRepository)
		service := NewContactService(mockRepo)

		userID := int64(1)
		newContact := &models.Contact{
			UserID:          &userID,
			ContactName:     "Test User",
			ContactPosition: "Developer",
		}

		mockRepo.On("CreateContact", mock.Anything, newContact).Return(errors.New("repository error"))

		err := service.CreateContact(context.Background(), newContact)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar contato")
		mockRepo.AssertExpectations(t)
	})
}

func TestGetContactByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockContactRepository)
		service := NewContactService(mockRepo)

		expectedContact := &models.Contact{
			ID:          1,
			ContactName: "Test User",
			Email:       "test@example.com",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockRepo.On("GetContactByID", mock.Anything, int64(1)).Return(expectedContact, nil)

		contact, err := service.GetContactByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedContact, contact)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(MockContactRepository)
		service := NewContactService(mockRepo)

		mockRepo.On("GetContactByID", mock.Anything, int64(1)).Return((*models.Contact)(nil), repositories.ErrContactNotFound)

		contact, err := service.GetContactByID(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, contact)
		assert.Contains(t, err.Error(), "contato não encontrado")
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockRepo := new(MockContactRepository)
		service := NewContactService(mockRepo)

		contact, err := service.GetContactByID(context.Background(), 0)

		assert.Error(t, err)
		assert.Nil(t, contact)
		assert.Contains(t, err.Error(), "ID inválido")
		mockRepo.AssertNotCalled(t, "GetContactByID")
	})
}

func TestUpdateContact(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockContactRepository)
		service := NewContactService(mockRepo)

		existingContact := &models.Contact{
			ID:          1,
			ContactName: "Old Name",
			Email:       "old@example.com",
		}

		updatedContact := &models.Contact{
			ID:          1,
			ContactName: "New Name",
			Email:       "new@example.com",
		}

		mockRepo.On("GetContactByID", mock.Anything, int64(1)).Return(existingContact, nil)

		mockRepo.On("Updatecontac", mock.Anything, updatedContact).Return(nil)

		err := service.UpdateContact(context.Background(), updatedContact)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("contact not found", func(t *testing.T) {
		mockRepo := new(MockContactRepository)
		service := NewContactService(mockRepo)

		updatedContact := &models.Contact{
			ID:          1,
			ContactName: "New Name",
		}

		mockRepo.On("GetContactByID", mock.Anything, int64(1)).Return((*models.Contact)(nil), repositories.ErrContactNotFound)

		err := service.UpdateContact(context.Background(), updatedContact)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "contato não encontrado")
		mockRepo.AssertNotCalled(t, "Updatecontac")
	})

	t.Run("validation error", func(t *testing.T) {
		mockRepo := new(MockContactRepository)
		service := NewContactService(mockRepo)

		invalidContact := &models.Contact{
			ID: 1,
		}

		err := service.UpdateContact(context.Background(), invalidContact)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nome do contato é obrigatório")
		mockRepo.AssertNotCalled(t, "GetContactByID")
		mockRepo.AssertNotCalled(t, "Updatecontac")
	})
}

func TestDeleteContact(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockContactRepository)
		service := NewContactService(mockRepo)

		existingContact := &models.Contact{
			ID:          1,
			ContactName: "Test User",
		}

		mockRepo.On("GetContactByID", mock.Anything, int64(1)).Return(existingContact, nil)
		mockRepo.On("Deletecontact", mock.Anything, int64(1)).Return(nil)

		err := service.DeleteContact(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("contact not found", func(t *testing.T) {
		mockRepo := new(MockContactRepository)
		service := NewContactService(mockRepo)

		mockRepo.On("GetContactByID", mock.Anything, int64(1)).Return((*models.Contact)(nil), repositories.ErrContactNotFound)

		err := service.DeleteContact(context.Background(), 1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "contato não encontrado")
		mockRepo.AssertNotCalled(t, "Deletecontact")
	})

	t.Run("invalid id", func(t *testing.T) {
		mockRepo := new(MockContactRepository)
		service := NewContactService(mockRepo)

		err := service.DeleteContact(context.Background(), 0)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ID inválido")
		mockRepo.AssertNotCalled(t, "GetContactByID")
		mockRepo.AssertNotCalled(t, "Deletecontact")
	})
}

func TestGetContactsByUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(MockContactRepository)
		service := NewContactService(mockRepo)

		expectedContacts := []*models.Contact{
			{
				ID:          1,
				UserID:      ptrInt64(1),
				ContactName: "User 1",
			},
			{
				ID:          2,
				UserID:      ptrInt64(1),
				ContactName: "User 2",
			},
		}

		mockRepo.On("GetContactByUserID", mock.Anything, int64(1)).Return(expectedContacts, nil)

		contacts, err := service.GetContactsByUser(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedContacts, contacts)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		mockRepo := new(MockContactRepository)
		service := NewContactService(mockRepo)

		mockRepo.On("GetContactByUserID", mock.Anything, int64(1)).Return([]*models.Contact{}, nil)

		contacts, err := service.GetContactsByUser(context.Background(), 1)

		assert.NoError(t, err)
		assert.Empty(t, contacts)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid user id", func(t *testing.T) {
		mockRepo := new(MockContactRepository)
		service := NewContactService(mockRepo)

		contacts, err := service.GetContactsByUser(context.Background(), 0)

		assert.Error(t, err)
		assert.Nil(t, contacts)
		assert.Contains(t, err.Error(), "ID de usuário inválido")
		mockRepo.AssertNotCalled(t, "GetContactByUserID")
	})
}

func ptrInt64(i int64) *int64 {
	return &i
}
