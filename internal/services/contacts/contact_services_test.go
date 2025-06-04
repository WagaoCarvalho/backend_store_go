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

func TestCreateContact(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		userID := int64(1)
		inputContact := &models.Contact{
			UserID:          &userID,
			ContactName:     "Test User",
			ContactPosition: "Developer",
			Email:           "test@example.com",
		}

		expectedContact := inputContact
		expectedContact.ID = 1

		mockRepo.On("Create", mock.Anything, inputContact).Return(expectedContact, nil)

		created, err := service.Create(context.Background(), inputContact)

		assert.NoError(t, err)
		assert.Equal(t, expectedContact, created)
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation error - missing name", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		userID := int64(1)
		invalidContact := &models.Contact{
			UserID: &userID,
			Email:  "test@example.com",
		}

		_, err := service.Create(context.Background(), invalidContact)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nome do contato é obrigatório")
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		userID := int64(1)
		inputContact := &models.Contact{
			UserID:          &userID,
			ContactName:     "Test User",
			ContactPosition: "Developer",
		}

		mockRepo.On("Create", mock.Anything, inputContact).Return(&models.Contact{}, errors.New("repository error"))

		_, err := service.Create(context.Background(), inputContact)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao criar contato")
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation error - no user, client or supplier ID", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		invalidContact := &models.Contact{
			ContactName: "No Relation",
			Email:       "valid@example.com",
		}

		_, err := service.Create(context.Background(), invalidContact)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "o contato deve estar associado")
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("validation error - invalid email", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		clientID := int64(2)
		invalidContact := &models.Contact{
			ClientID:        &clientID,
			ContactName:     "Invalid Email",
			Email:           "invalid-email",
			ContactPosition: "Manager",
		}

		_, err := service.Create(context.Background(), invalidContact)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "email inválido")
		mockRepo.AssertNotCalled(t, "Create")
	})

}

func TestGetContactByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		expectedContact := &models.Contact{
			ID:          1,
			ContactName: "Test User",
			Email:       "test@example.com",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expectedContact, nil)

		contact, err := service.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedContact, contact)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return((*models.Contact)(nil), repositories.ErrContactNotFound)

		contact, err := service.GetByID(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, contact)
		assert.Contains(t, err.Error(), "contato não encontrado")
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		contact, err := service.GetByID(context.Background(), 0)

		assert.Error(t, err)
		assert.Nil(t, contact)
		assert.Contains(t, err.Error(), "ID inválido")
		mockRepo.AssertNotCalled(t, "GetContactByID")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		mockRepo.On("GetByID", mock.Anything, int64(2)).
			Return((*models.Contact)(nil), errors.New("erro inesperado"))

		contact, err := service.GetByID(context.Background(), 2)

		assert.Error(t, err)
		assert.Nil(t, contact)
		assert.Contains(t, err.Error(), "erro ao verificar contato")
		mockRepo.AssertExpectations(t)
	})

}

func TestGetContactsBySupplier(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		expectedContacts := []*models.Contact{
			{
				ID:          1,
				SupplierID:  ptrInt64(1),
				ContactName: "Fornecedor 1",
			},
			{
				ID:          2,
				SupplierID:  ptrInt64(1),
				ContactName: "Fornecedor 2",
			},
		}

		mockRepo.On("GetBySupplierID", mock.Anything, int64(1)).
			Return(expectedContacts, nil)

		contacts, err := service.GetBySupplier(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedContacts, contacts)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		mockRepo.On("GetBySupplierID", mock.Anything, int64(1)).
			Return([]*models.Contact{}, nil)

		contacts, err := service.GetBySupplier(context.Background(), 1)

		assert.NoError(t, err)
		assert.Empty(t, contacts)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid supplier id", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		contacts, err := service.GetBySupplier(context.Background(), 0)

		assert.Error(t, err)
		assert.Nil(t, contacts)
		assert.Contains(t, err.Error(), "ID de fornecedor inválido")
		mockRepo.AssertNotCalled(t, "GetBySupplierID")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		mockRepo.On("GetBySupplierID", mock.Anything, int64(1)).
			Return([]*models.Contact(nil), errors.New("erro de banco"))

		contacts, err := service.GetBySupplier(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, contacts)
		assert.Contains(t, err.Error(), "erro ao listar contatos do fornecedor")
		mockRepo.AssertExpectations(t)
	})

}

func TestGetContactsByClient(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		expectedContacts := []*models.Contact{
			{
				ID:          1,
				ClientID:    ptrInt64(1),
				ContactName: "Cliente 1",
			},
			{
				ID:          2,
				ClientID:    ptrInt64(1),
				ContactName: "Cliente 2",
			},
		}

		mockRepo.On("GetByClientID", mock.Anything, int64(1)).
			Return(expectedContacts, nil)

		contacts, err := service.GetByClient(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedContacts, contacts)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		mockRepo.On("GetByClientID", mock.Anything, int64(1)).
			Return([]*models.Contact{}, nil)

		contacts, err := service.GetByClient(context.Background(), 1)

		assert.NoError(t, err)
		assert.Empty(t, contacts)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid client id", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		contacts, err := service.GetByClient(context.Background(), 0)

		assert.Error(t, err)
		assert.Nil(t, contacts)
		assert.Contains(t, err.Error(), "ID de cliente inválido")
		mockRepo.AssertNotCalled(t, "GetByClientID")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		mockRepo.On("GetByClientID", mock.Anything, int64(1)).
			Return([]*models.Contact(nil), errors.New("erro de banco"))

		contacts, err := service.GetByClient(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, contacts)
		assert.Contains(t, err.Error(), "erro ao listar contatos do cliente")
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateContact(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		existingContact := &models.Contact{
			ID:          1,
			ContactName: "Old Name",
			Email:       "old@example.com",
			Version:     1,
		}

		updatedContact := &models.Contact{
			ID:          1,
			ContactName: "New Name",
			Email:       "new@example.com",
			Version:     1,
		}

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(existingContact, nil)
		mockRepo.On("Update", mock.Anything, updatedContact).Return(nil)

		err := service.Update(context.Background(), updatedContact)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("contact not found", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		updatedContact := &models.Contact{
			ID:          1,
			ContactName: "New Name",
			Version:     1,
		}

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return((*models.Contact)(nil), repositories.ErrContactNotFound)

		err := service.Update(context.Background(), updatedContact)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "contato não encontrado")
		mockRepo.AssertNotCalled(t, "Updatecontac")
	})

	t.Run("validation error", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		invalidContact := &models.Contact{
			ID: 1,
		}

		err := service.Update(context.Background(), invalidContact)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "nome do contato é obrigatório")
		mockRepo.AssertNotCalled(t, "GetByID")
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("validation error - invalid ID", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		invalidContact := &models.Contact{
			ID:          0,
			ContactName: "Nome válido",
		}

		err := service.Update(context.Background(), invalidContact)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ID inválido")
		mockRepo.AssertNotCalled(t, "GetByID")
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("repository error on GetByID", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		contact := &models.Contact{
			ID:          1,
			ContactName: "Valid Name",
			Version:     1,
		}

		mockRepo.On("GetByID", mock.Anything, int64(1)).
			Return((*models.Contact)(nil), errors.New("erro inesperado"))

		err := service.Update(context.Background(), contact)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao verificar existência do contato antes da atualização")

		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("repository error on Update", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		contact := &models.Contact{
			ID:          1,
			ContactName: "Valid Name",
			Version:     1,
		}

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(contact, nil)
		mockRepo.On("Update", mock.Anything, contact).Return(errors.New("falha ao atualizar"))

		err := service.Update(context.Background(), contact)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao atualizar contato")
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation error - missing version", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		contact := &models.Contact{
			ID:          1,
			ContactName: "Nome válido",
			Version:     0, // versão ausente
		}

		err := service.Update(context.Background(), contact)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "versão do contato é obrigatória")
		mockRepo.AssertNotCalled(t, "GetByID")
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("version conflict error", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		contact := &models.Contact{
			ID:          1,
			ContactName: "Nome válido",
			Version:     1,
		}

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(contact, nil)
		mockRepo.On("Update", mock.Anything, contact).Return(repositories.ErrVersionConflict)

		err := service.Update(context.Background(), contact)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "conflito de versão")
		mockRepo.AssertExpectations(t)
	})

}

func TestDeleteContact(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		existingContact := &models.Contact{
			ID:          1,
			ContactName: "Test User",
		}

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(existingContact, nil)
		mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

		err := service.Delete(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("contact not found", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return((*models.Contact)(nil), repositories.ErrContactNotFound)

		err := service.Delete(context.Background(), 1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "contato não encontrado")
		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("invalid id", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		err := service.Delete(context.Background(), 0)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "ID inválido")
		mockRepo.AssertNotCalled(t, "GetByID")
		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("repository error on GetByID", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		mockRepo.On("GetByID", mock.Anything, int64(1)).
			Return((*models.Contact)(nil), errors.New("erro inesperado"))

		err := service.Delete(context.Background(), 1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao verificar contato")
		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("repository error on Delete", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		existingContact := &models.Contact{
			ID:          1,
			ContactName: "Test User",
		}

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(existingContact, nil)
		mockRepo.On("Delete", mock.Anything, int64(1)).Return(errors.New("falha ao deletar"))

		err := service.Delete(context.Background(), 1)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao deletar contato")
		mockRepo.AssertExpectations(t)
	})

}

func TestGetContactsByUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
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

		mockRepo.On("GetByUserID", mock.Anything, int64(1)).Return(expectedContacts, nil)

		contacts, err := service.GetByUser(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedContacts, contacts)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		mockRepo.On("GetByUserID", mock.Anything, int64(1)).Return([]*models.Contact{}, nil)

		contacts, err := service.GetByUser(context.Background(), 1)

		assert.NoError(t, err)
		assert.Empty(t, contacts)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid user id", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		contacts, err := service.GetByUser(context.Background(), 0)

		assert.Error(t, err)
		assert.Nil(t, contacts)
		assert.Contains(t, err.Error(), "ID de usuário inválido")
		mockRepo.AssertNotCalled(t, "GetByUserID")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo)

		mockRepo.On("GetByUserID", mock.Anything, int64(1)).
			Return([]*models.Contact(nil), errors.New("erro de banco"))

		contacts, err := service.GetByUser(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, contacts)
		assert.Contains(t, err.Error(), "erro ao listar contatos do usuário")
		mockRepo.AssertExpectations(t)
	})

}

func ptrInt64(i int64) *int64 {
	return &i
}
