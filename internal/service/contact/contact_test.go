package services

import (
	"context"
	"errors"
	"testing"
	"time"

	mockContact "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/contact"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestContactService_Create(t *testing.T) {
	t.Run("sucesso na criação do contato", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		userID := int64(1)
		contact := &model.Contact{
			UserID:      &userID,
			ContactName: "Contato Teste",
			Email:       "teste@email.com",
			Phone:       "1234567898",
		}

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(m *model.Contact) bool {
			return m.UserID != nil &&
				*m.UserID == *contact.UserID &&
				m.ContactName == contact.ContactName &&
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

	t.Run("falha na validação do contato ContactName obrigatório", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		userID := int64(1)
		clientID := int64(1)
		supplierID := int64(1)

		contact := &model.Contact{
			UserID:     &userID,
			ClientID:   &clientID,
			SupplierID: &supplierID,
			Email:      "teste@email.com",
			Phone:      "1234567898",
			// ContactName vazio para gerar erro
		}

		createdContact, err := service.Create(context.Background(), contact)

		assert.Nil(t, createdContact)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("falha no repositório ao criar contato", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		userID := int64(1)
		contact := &model.Contact{
			UserID:      &userID,
			ContactName: "Contato Teste",
			Email:       "teste@email.com",
			Phone:       "1234567898",
		}

		expectedErr := errors.New("erro no banco")

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(m *model.Contact) bool {
			return m.UserID != nil &&
				*m.UserID == *contact.UserID &&
				m.ContactName == contact.ContactName &&
				m.Email == contact.Email &&
				m.Phone == contact.Phone
		})).Return((*model.Contact)(nil), expectedErr)

		createdContact, err := service.Create(context.Background(), contact)

		assert.Nil(t, createdContact)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestContactService_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

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
		assert.EqualError(t, err, errMsg.ErrID.Error())
		mockRepo.AssertNotCalled(t, "GetByID")
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

func TestContactService_GetByUserID(t *testing.T) {
	t.Run("sucesso ao buscar contatos por UserID", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		contactModels := []*model.Contact{
			{
				ID:          1,
				UserID:      utils.Int64Ptr(1),
				ContactName: "Contato Teste",
				Email:       "teste@email.com",
				Phone:       "123456789",
			},
		}

		mockRepo.On("GetByUserID", mock.Anything, int64(1)).
			Return(contactModels, nil)

		result, err := service.GetByUserID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1), result[0].ID)
		assert.Equal(t, "Contato Teste", result[0].ContactName)
		assert.Equal(t, "teste@email.com", result[0].Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha ao buscar contatos com UserID inválido", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		result, err := service.GetByUserID(context.Background(), 0)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, errMsg.ErrID.Error())
		mockRepo.AssertNotCalled(t, "GetByUserID")
	})

	t.Run("nenhum contato encontrado por UserID", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		mockRepo.On("GetByUserID", mock.Anything, int64(2)).
			Return(([]*model.Contact)(nil), errMsg.ErrNotFound)

		result, err := service.GetByUserID(context.Background(), 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, errMsg.ErrNotFound.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestContactService_GetByClientID(t *testing.T) {
	t.Run("sucesso ao buscar contatos por ClientID", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		contactModels := []*model.Contact{
			{
				ID:          1,
				ClientID:    utils.Int64Ptr(1),
				ContactName: "Cliente Teste",
				Email:       "cliente@email.com",
				Phone:       "987654321",
			},
		}

		mockRepo.On("GetByClientID", mock.Anything, int64(1)).
			Return(contactModels, nil)

		result, err := service.GetByClientID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1), result[0].ID)
		assert.Equal(t, "Cliente Teste", result[0].ContactName)
		assert.Equal(t, "cliente@email.com", result[0].Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha ao buscar contatos com ClientID inválido", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		result, err := service.GetByClientID(context.Background(), 0)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, errMsg.ErrID.Error())
		mockRepo.AssertNotCalled(t, "GetByClientID")
	})

	t.Run("nenhum contato encontrado por ClientID", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		mockRepo.On("GetByClientID", mock.Anything, int64(2)).
			Return(([]*model.Contact)(nil), errMsg.ErrNotFound)

		result, err := service.GetByClientID(context.Background(), 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, errMsg.ErrNotFound.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestContactService_GetBySupplierID(t *testing.T) {
	t.Run("sucesso ao buscar contatos por SupplierID", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		contactModels := []*model.Contact{
			{
				ID:          1,
				SupplierID:  utils.Int64Ptr(1),
				ContactName: "Fornecedor Teste",
				Email:       "fornecedor@email.com",
				Phone:       "1122334455",
			},
		}

		mockRepo.On("GetBySupplierID", mock.Anything, int64(1)).
			Return(contactModels, nil)

		result, err := service.GetBySupplierID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1), result[0].ID)
		assert.Equal(t, "Fornecedor Teste", result[0].ContactName)
		assert.Equal(t, "fornecedor@email.com", result[0].Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha ao buscar contatos com SupplierID inválido", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		result, err := service.GetBySupplierID(context.Background(), 0)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, errMsg.ErrID.Error())
		mockRepo.AssertNotCalled(t, "GetBySupplierID")
	})

	t.Run("nenhum contato encontrado por SupplierID", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		mockRepo.On("GetBySupplierID", mock.Anything, int64(2)).
			Return(([]*model.Contact)(nil), errMsg.ErrNotFound)

		result, err := service.GetBySupplierID(context.Background(), 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, errMsg.ErrNotFound.Error())
		mockRepo.AssertExpectations(t)
	})
}

func TestContactService_UpdateContact(t *testing.T) {
	makeContactModel := func() *model.Contact {
		userID := int64(1)
		return &model.Contact{
			ID:          1,
			UserID:      &userID,
			ContactName: "Contato Teste",
			Email:       "teste@email.com",
			Phone:       "1234567898",
		}
	}

	t.Run("sucesso na atualização do contato", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		contact := makeContactModel()

		mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(c *model.Contact) bool {
			return c != nil && c.ID == contact.ID && *c.UserID == *contact.UserID
		})).Return(nil)

		err := service.Update(context.Background(), contact)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro ao atualizar contato com ID inválido", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		userID := int64(1)
		contact := &model.Contact{
			UserID:      &userID,
			ContactName: "Nome Teste",
			Email:       "teste@email.com",
		}

		err := service.Update(context.Background(), contact)
		assert.ErrorIs(t, err, errMsg.ErrID)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("falha na validação do contato no update", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		contact := &model.Contact{
			ID:    1,
			Email: "teste@email.com",
		}

		err := service.Update(context.Background(), contact)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrInvalidData)
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("erro genérico ao atualizar contato", func(t *testing.T) {
		mockRepo := new(mockContact.MockContactRepository)
		service := NewContactService(mockRepo)

		contact := &model.Contact{
			ID:          1,
			UserID:      utils.Int64Ptr(1),
			ContactName: "Nome Teste",
			Email:       "teste@email.com",
			Phone:       "1112345678",
		}

		mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(c *model.Contact) bool {
			return c.ID == contact.ID
		})).Return(errors.New("erro inesperado no banco"))

		err := service.Update(context.Background(), contact)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "erro ao atualizar")
		assert.ErrorContains(t, err, "erro inesperado no banco")
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
