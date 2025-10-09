package services

import (
	"context"
	"errors"
	"testing"
	"time"

	mockContact "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/contact"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestContactService_Create(t *testing.T) {
	t.Run("sucesso na criação do contato", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		contact := &model.Contact{
			ContactName: "Contato Teste",
			Email:       "teste@email.com",
			Phone:       "1234567898",
		}

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(m *model.Contact) bool {
			return m.ContactName == contact.ContactName &&
				m.Email == contact.Email &&
				m.Phone == contact.Phone
		})).Return(contact, nil)

		createdContact, err := service.Create(context.Background(), contact)

		assert.NoError(t, err)
		assert.NotNil(t, createdContact)
		assert.Equal(t, contact.ContactName, createdContact.ContactName)
		assert.Equal(t, contact.Email, createdContact.Email)
		assert.Equal(t, contact.Phone, createdContact.Phone)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha na validação do contato", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		contact := &model.Contact{} // campos obrigatórios ausentes

		createdContact, err := service.Create(context.Background(), contact)

		assert.Nil(t, createdContact)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("erro do repositório ao criar contato", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		contact := &model.Contact{
			ContactName: "Contato Teste",
			Email:       "teste@email.com",
			Phone:       "1234567898",
		}

		expectedErr := errors.New("erro no banco")
		mockRepo.On("Create", mock.Anything, mock.Anything).Return((*model.Contact)(nil), expectedErr)

		createdContact, err := service.Create(context.Background(), contact)

		assert.Nil(t, createdContact)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrCreate)
		assert.Contains(t, err.Error(), expectedErr.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestContactService_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo) // apenas contact

		expectedContact := &model.Contact{
			ID:          1,
			ContactName: "Test User",
			Email:       "test@example.com",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockRepo.On("GetByID", mock.Anything, int64(1)).
			Return(expectedContact, nil)

		contact, err := service.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, contact)
		assert.Equal(t, expectedContact.ID, contact.ID)
		assert.Equal(t, expectedContact.ContactName, contact.ContactName)
		assert.Equal(t, expectedContact.Email, contact.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		mockRepo.On("GetByID", mock.Anything, int64(1)).
			Return((*model.Contact)(nil), errMsg.ErrNotFound)

		contact, err := service.GetByID(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, contact)
		assert.EqualError(t, err, errMsg.ErrNotFound.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		contact, err := service.GetByID(context.Background(), 0)

		assert.Error(t, err)
		assert.Nil(t, contact)
		assert.EqualError(t, err, errMsg.ErrZeroID.Error())
		mockRepo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		mockRepo.On("GetByID", mock.Anything, int64(2)).
			Return((*model.Contact)(nil), errors.New("erro inesperado"))

		contact, err := service.GetByID(context.Background(), 2)

		assert.Error(t, err)
		assert.Nil(t, contact)
		assert.Contains(t, err.Error(), "erro ao buscar")
		mockRepo.AssertExpectations(t)
	})
}

func TestContactService_Update(t *testing.T) {
	t.Run("sucesso na atualização do contato", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		contact := &model.Contact{
			ID:          1,
			ContactName: "Contato Atualizado",
			Email:       "teste@atualizado.com",
			Phone:       "1234567898",
		}

		mockRepo.On("Update", mock.Anything, contact).Return(nil)

		err := service.Update(context.Background(), contact)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha por ID zero", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		contact := &model.Contact{
			ID:          0,
			ContactName: "Contato Teste",
		}

		err := service.Update(context.Background(), contact)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("falha na validação do contato", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		contact := &model.Contact{
			ID: 1,
			// Campos obrigatórios ausentes para falhar na validação
		}

		err := service.Update(context.Background(), contact)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		mockRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("erro do repositório ao atualizar", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		contact := &model.Contact{
			ID:          1,
			ContactName: "Contato Teste",
			Email:       "teste@email.com",
		}

		expectedErr := errors.New("erro no banco")
		mockRepo.On("Update", mock.Anything, contact).Return(expectedErr)

		err := service.Update(context.Background(), contact)

		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrUpdate)
		assert.Contains(t, err.Error(), expectedErr.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestContactService_DeleteContact(t *testing.T) {
	t.Run("sucesso ao deletar contato", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		existingContact := &model.Contact{ID: 1, ContactName: "Test User"}
		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(existingContact, nil)
		mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

		err := service.Delete(context.Background(), 1)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("contato não encontrado", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		mockRepo.On("GetByID", mock.Anything, int64(1)).
			Return((*model.Contact)(nil), errMsg.ErrNotFound)

		err := service.Delete(context.Background(), 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "não encontrado")
		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("id inválido", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		err := service.Delete(context.Background(), 0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ID")

		mockRepo.AssertNotCalled(t, "GetByID")
		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("erro do repositório no GetByID", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		mockRepo.On("GetByID", mock.Anything, int64(1)).
			Return((*model.Contact)(nil), errors.New("erro inesperado"))

		err := service.Delete(context.Background(), 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao buscar")
		mockRepo.AssertNotCalled(t, "Delete")
	})

	t.Run("erro do repositório no Delete", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		existingContact := &model.Contact{ID: 1, ContactName: "Test User"}
		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(existingContact, nil)
		mockRepo.On("Delete", mock.Anything, int64(1)).Return(errors.New("falha ao deletar"))

		err := service.Delete(context.Background(), 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "erro ao deletar")
		mockRepo.AssertExpectations(t)
	})
}
