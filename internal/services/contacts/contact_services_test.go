package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/contacts"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New()) // Logger real (sem necessidade de mock)

	t.Run("success", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		userID := int64(1)
		inputContact := &models.Contact{
			UserID:          &userID,
			ContactName:     "Test User",
			ContactPosition: "Developer",
			Email:           "test@example.com",
		}
		expectedContact := *inputContact
		expectedContact.ID = 1

		mockRepo.On("Create", mock.Anything, inputContact).Return(&expectedContact, nil)

		created, err := service.Create(context.Background(), inputContact)
		assert.NoError(t, err)
		assert.Equal(t, &expectedContact, created)
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation error - missing name", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		userID := int64(1)
		invalidContact := &models.Contact{
			UserID: &userID,
			Email:  "test@example.com",
		}

		_, err := service.Create(context.Background(), invalidContact)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "ContactName")
		assert.ErrorContains(t, err, "campo obrigatório")
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		userID := int64(1)
		inputContact := &models.Contact{
			UserID:          &userID,
			ContactName:     "Test User",
			ContactPosition: "Developer",
			Email:           "test@example.com",
		}
		repoErr := errors.New("repository error")
		mockRepo.On("Create", mock.Anything, inputContact).Return((*models.Contact)(nil), repoErr)

		_, err := service.Create(context.Background(), inputContact)

		assert.Error(t, err)
		assert.ErrorIs(t, err, ErrCreateContact)
		assert.Contains(t, err.Error(), "repository error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation error - no user, client or supplier ID", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		invalidContact := &models.Contact{
			ContactName: "No Relation",
			Email:       "valid@example.com",
		}
		_, err := service.Create(context.Background(), invalidContact)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "UserID/ClientID/SupplierID")
		assert.ErrorContains(t, err, "pelo menos um deve ser informado")
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("validation error - invalid email", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		clientID := int64(2)
		invalidContact := &models.Contact{
			ClientID:        &clientID,
			ContactName:     "Invalid Email",
			Email:           "invalid-email",
			ContactPosition: "Manager",
		}
		_, err := service.Create(context.Background(), invalidContact)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "formato inválido")
		mockRepo.AssertNotCalled(t, "Create")
	})

}

func Test_GetByID(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New()) // logger real ou mock

	t.Run("success", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo, logger)

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
		service := NewContactService(mockRepo, logger)

		mockRepo.On("GetByID", mock.Anything, int64(1)).
			Return((*models.Contact)(nil), repositories.ErrContactNotFound)

		contact, err := service.GetByID(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, contact)
		assert.Contains(t, err.Error(), "contato não encontrado")
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		contact, err := service.GetByID(context.Background(), 0)

		assert.Error(t, err)
		assert.Nil(t, contact)
		assert.Contains(t, err.Error(), "ID inválido")
		mockRepo.AssertNotCalled(t, "GetByID")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		mockRepo.On("GetByID", mock.Anything, int64(2)).
			Return((*models.Contact)(nil), errors.New("erro inesperado"))

		contact, err := service.GetByID(context.Background(), 2)

		assert.Error(t, err)
		assert.Nil(t, contact)
		assert.Contains(t, err.Error(), "erro ao verificar contato")
		mockRepo.AssertExpectations(t)
	})
}

func Test_GetByUserID(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	t.Run("success", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		expectedContacts := []*models.Contact{
			{ID: 1, UserID: ptrInt64(1), ContactName: "User 1"},
			{ID: 2, UserID: ptrInt64(1), ContactName: "User 2"},
		}

		mockRepo.On("GetByUserID", mock.Anything, int64(1)).Return(expectedContacts, nil)

		contacts, err := service.GetByUserID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedContacts, contacts)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		mockRepo.On("GetByUserID", mock.Anything, int64(1)).Return([]*models.Contact{}, nil)

		contacts, err := service.GetByUserID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Empty(t, contacts)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid user id", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		contacts, err := service.GetByUserID(context.Background(), 0)

		assert.Error(t, err)
		assert.Nil(t, contacts)
		assert.Contains(t, err.Error(), "ID de usuário inválido")
		mockRepo.AssertNotCalled(t, "GetByUserID")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		mockRepo.On("GetByUserID", mock.Anything, int64(1)).
			Return(([]*models.Contact)(nil), errors.New("erro de banco"))

		contacts, err := service.GetByUserID(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, contacts)
		assert.Contains(t, err.Error(), "erro ao listar contatos do usuário")
		mockRepo.AssertExpectations(t)
	})
}

func Test_GetByClientID(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	t.Run("success", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		expectedContacts := []*models.Contact{
			{ID: 1, ClientID: ptrInt64(1), ContactName: "Cliente 1"},
			{ID: 2, ClientID: ptrInt64(1), ContactName: "Cliente 2"},
		}

		mockRepo.On("GetByClientID", mock.Anything, int64(1)).
			Return(expectedContacts, nil)

		contacts, err := service.GetByClientID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedContacts, contacts)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		mockRepo.On("GetByClientID", mock.Anything, int64(1)).
			Return([]*models.Contact{}, nil)

		contacts, err := service.GetByClientID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Empty(t, contacts)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid client id", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		contacts, err := service.GetByClientID(context.Background(), 0)

		assert.Error(t, err)
		assert.Nil(t, contacts)
		assert.Contains(t, err.Error(), "ID de cliente inválido")
		mockRepo.AssertNotCalled(t, "GetByClientID")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		mockRepo.On("GetByClientID", mock.Anything, int64(1)).
			Return(([]*models.Contact)(nil), errors.New("erro de banco"))

		contacts, err := service.GetByClientID(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, contacts)
		assert.Contains(t, err.Error(), "erro ao listar contatos do cliente")
		mockRepo.AssertExpectations(t)
	})
}

func Test_GetBySupplierID(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New()) // ou um mock se preferir

	t.Run("success", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo, logger)

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

		mockRepo.On("GetBySupplierID", mock.Anything, int64(1)).Return(expectedContacts, nil)

		contacts, err := service.GetBySupplierID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedContacts, contacts)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		mockRepo.On("GetBySupplierID", mock.Anything, int64(1)).Return([]*models.Contact{}, nil)

		contacts, err := service.GetBySupplierID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Empty(t, contacts)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid supplier id", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		contacts, err := service.GetBySupplierID(context.Background(), 0)

		assert.Error(t, err)
		assert.Nil(t, contacts)
		assert.Contains(t, err.Error(), "ID de fornecedor inválido")
		mockRepo.AssertNotCalled(t, "GetBySupplierID")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		mockRepo.On("GetBySupplierID", mock.Anything, int64(1)).
			Return([]*models.Contact(nil), errors.New("erro de banco"))

		contacts, err := service.GetBySupplierID(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, contacts)
		assert.Contains(t, err.Error(), "erro ao listar contatos do fornecedor")
		mockRepo.AssertExpectations(t)
	})
}

func TestContactService_Update(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	t.Run("success", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		userID := int64(1)
		existingContact := &models.Contact{
			ID:          1,
			ContactName: "Old Name",
			UserID:      &userID, // essencial para validação passar
		}
		updatedContact := &models.Contact{
			ID:          1,
			ContactName: "New Name",
			UserID:      &userID, // essencial para validação passar
		}

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(existingContact, nil).Once()
		mockRepo.On("Update", mock.Anything, updatedContact).Return(nil).Once()

		err := service.Update(context.Background(), updatedContact)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("contact not found", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		userID := int64(1)
		updatedContact := &models.Contact{
			ID:          1,
			ContactName: "New Name",
			UserID:      &userID,
		}

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return((*models.Contact)(nil), repositories.ErrContactNotFound).Once()

		err := service.Update(context.Background(), updatedContact)
		assert.ErrorIs(t, err, ErrContactNotFound)

		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("validation error - missing contact name", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		userID := int64(1)
		invalidContact := &models.Contact{
			ID:     1,
			UserID: &userID, // Necessário para passar da validação do relacionamento
		}

		err := service.Update(context.Background(), invalidContact)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "campo obrigatório") // Valida o erro do nome

		mockRepo.AssertNotCalled(t, "GetByID")
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("validation error - invalid ID", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		invalidContact := &models.Contact{ID: 0, ContactName: "Válido"}

		err := service.Update(context.Background(), invalidContact)
		assert.ErrorIs(t, err, ErrInvalidID)
		mockRepo.AssertNotCalled(t, "GetByID")
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("repository error on GetByID", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		userID := int64(1)
		contact := &models.Contact{
			ID:          1,
			ContactName: "Válido",
			UserID:      &userID,
		}

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return((*models.Contact)(nil), errors.New("erro inesperado")).Once()

		err := service.Update(context.Background(), contact)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "erro ao verificar existência do contato antes da atualização")

		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("repository error on Update", func(t *testing.T) {
		mockRepo := new(repositories.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		userID := int64(1)
		contact := &models.Contact{
			ID:          1,
			ContactName: "Válido",
			UserID:      &userID,
		}

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(contact, nil).Once()
		mockRepo.On("Update", mock.Anything, contact).Return(errors.New("falha ao atualizar")).Once()

		err := service.Update(context.Background(), contact)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "erro ao atualizar contato")

		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteContact(t *testing.T) {
	for name, tt := range map[string]func(t *testing.T){
		"success": func(t *testing.T) {
			mockRepo := new(repositories.MockContactRepository)
			logger := logger.NewLoggerAdapter(logrus.New())
			service := NewContactService(mockRepo, logger)

			existingContact := &models.Contact{
				ID:          1,
				ContactName: "Test User",
			}
			mockRepo.On("GetByID", mock.Anything, int64(1)).Return(existingContact, nil)
			mockRepo.On("Delete", mock.Anything, int64(1)).Return(nil)

			err := service.Delete(context.Background(), 1)
			assert.NoError(t, err)

			mockRepo.AssertExpectations(t)
		},
		"contact not found": func(t *testing.T) {
			mockRepo := new(repositories.MockContactRepository)
			logger := logger.NewLoggerAdapter(logrus.New())
			service := NewContactService(mockRepo, logger)

			mockRepo.On("GetByID", mock.Anything, int64(1)).
				Return((*models.Contact)(nil), repositories.ErrContactNotFound)

			err := service.Delete(context.Background(), 1)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "contato não encontrado")

			mockRepo.AssertNotCalled(t, "Delete")
		},
		"invalid id": func(t *testing.T) {
			mockRepo := new(repositories.MockContactRepository)
			logger := logger.NewLoggerAdapter(logrus.New())
			service := NewContactService(mockRepo, logger)

			err := service.Delete(context.Background(), 0)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "ID inválido")

			mockRepo.AssertNotCalled(t, "GetByID")
			mockRepo.AssertNotCalled(t, "Delete")
		},
		"repository error on GetByID": func(t *testing.T) {
			mockRepo := new(repositories.MockContactRepository)
			logger := logger.NewLoggerAdapter(logrus.New())
			service := NewContactService(mockRepo, logger)

			mockRepo.On("GetByID", mock.Anything, int64(1)).
				Return((*models.Contact)(nil), errors.New("erro inesperado"))

			err := service.Delete(context.Background(), 1)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "erro ao verificar contato")

			mockRepo.AssertNotCalled(t, "Delete")
		},
		"repository error on Delete": func(t *testing.T) {
			mockRepo := new(repositories.MockContactRepository)
			logger := logger.NewLoggerAdapter(logrus.New())
			service := NewContactService(mockRepo, logger)

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
		},
	} {
		t.Run(name, tt)
	}
}

func ptrInt64(i int64) *int64 {
	return &i
}
