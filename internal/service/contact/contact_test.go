package services

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"

	mock_contact "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/contact"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	convert "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreate(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(log)

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
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
		mockRepo := new(mock_contact.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		userID := int64(1)
		invalidContact := &models.Contact{
			UserID: &userID,
			Email:  "test@example.com",
		}

		_, err := service.Create(context.Background(), invalidContact)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "contact_name")
		assert.ErrorContains(t, err, "campo obrigatório")
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
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
		assert.ErrorIs(t, err, err_msg.ErrCreate)
		assert.Contains(t, err.Error(), "repository error")
		mockRepo.AssertExpectations(t)
	})

	t.Run("validation error - no user, client or supplier ID", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		invalidContact := &models.Contact{
			ContactName: "No Relation",
			Email:       "valid@example.com",
		}
		_, err := service.Create(context.Background(), invalidContact)

		assert.Error(t, err)
		assert.ErrorContains(t, err, "UserID/ClientID/SupplierID")
		assert.ErrorContains(t, err, "exatamente um deve ser informado")
		mockRepo.AssertNotCalled(t, "Create")
	})

	t.Run("validation error - invalid email", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
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
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(log)

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
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
		mockRepo := new(mock_contact.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		mockRepo.On("GetByID", mock.Anything, int64(1)).
			Return((*models.Contact)(nil), err_msg.ErrNotFound)

		contact, err := service.GetByID(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, contact)
		assert.Contains(t, err.Error(), "não encontrado")
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		contact, err := service.GetByID(context.Background(), 0)

		assert.Error(t, err)
		assert.Nil(t, contact)
		assert.Contains(t, err.Error(), "erro ID")
		mockRepo.AssertNotCalled(t, "GetByID")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		mockRepo.On("GetByID", mock.Anything, int64(2)).
			Return((*models.Contact)(nil), errors.New("erro inesperado"))

		contact, err := service.GetByID(context.Background(), 2)

		assert.Error(t, err)
		assert.Nil(t, contact)
		assert.Contains(t, err.Error(), "erro ao buscar")
		mockRepo.AssertExpectations(t)
	})
}

func Test_GetByUserID(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(log)

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		expectedContacts := []*models.Contact{
			{ID: 1, UserID: convert.Int64Ptr(1), ContactName: "User 1"},
			{ID: 2, UserID: convert.Int64Ptr(1), ContactName: "User 2"},
		}

		mockRepo.On("GetByUserID", mock.Anything, int64(1)).Return(expectedContacts, nil)

		contacts, err := service.GetByUserID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedContacts, contacts)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		mockRepo.On("GetByUserID", mock.Anything, int64(1)).Return([]*models.Contact{}, nil)

		contacts, err := service.GetByUserID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Empty(t, contacts)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid user id", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		contacts, err := service.GetByUserID(context.Background(), 0)

		assert.Error(t, err)
		assert.Nil(t, contacts)
		assert.Contains(t, err.Error(), "ID inválido")
		mockRepo.AssertNotCalled(t, "GetByUserID")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		mockRepo.On("GetByUserID", mock.Anything, int64(1)).
			Return(([]*models.Contact)(nil), errors.New("erro de banco"))

		contacts, err := service.GetByUserID(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, contacts)
		assert.Contains(t, err.Error(), "erro ao buscar")
		mockRepo.AssertExpectations(t)
	})
}

func Test_GetByClientID(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(log)

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		expectedContacts := []*models.Contact{
			{ID: 1, ClientID: convert.Int64Ptr(1), ContactName: "Cliente 1"},
			{ID: 2, ClientID: convert.Int64Ptr(1), ContactName: "Cliente 2"},
		}

		mockRepo.On("GetByClientID", mock.Anything, int64(1)).
			Return(expectedContacts, nil)

		contacts, err := service.GetByClientID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, expectedContacts, contacts)
		mockRepo.AssertExpectations(t)
	})

	t.Run("empty list", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		mockRepo.On("GetByClientID", mock.Anything, int64(1)).
			Return([]*models.Contact{}, nil)

		contacts, err := service.GetByClientID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Empty(t, contacts)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid client id", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		contacts, err := service.GetByClientID(context.Background(), 0)

		assert.Error(t, err)
		assert.Nil(t, contacts)
		assert.Contains(t, err.Error(), "ID inválido")
		mockRepo.AssertNotCalled(t, "GetByClientID")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		mockRepo.On("GetByClientID", mock.Anything, int64(1)).
			Return(([]*models.Contact)(nil), errors.New("erro ao buscar"))

		contacts, err := service.GetByClientID(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, contacts)
		assert.Contains(t, err.Error(), "erro ao buscar")
		mockRepo.AssertExpectations(t)
	})
}

func Test_GetBySupplierID(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(log)

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		expectedContacts := []*models.Contact{
			{
				ID:          1,
				SupplierID:  convert.Int64Ptr(1),
				ContactName: "Fornecedor 1",
			},
			{
				ID:          2,
				SupplierID:  convert.Int64Ptr(1),
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
		mockRepo := new(mock_contact.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		mockRepo.On("GetBySupplierID", mock.Anything, int64(1)).Return([]*models.Contact{}, nil)

		contacts, err := service.GetBySupplierID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Empty(t, contacts)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid supplier id", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		contacts, err := service.GetBySupplierID(context.Background(), 0)

		assert.Error(t, err)
		assert.Nil(t, contacts)
		assert.Contains(t, err.Error(), "ID inválido")
		mockRepo.AssertNotCalled(t, "GetBySupplierID")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		mockRepo.On("GetBySupplierID", mock.Anything, int64(1)).
			Return([]*models.Contact(nil), errors.New("erro ao buscar"))

		contacts, err := service.GetBySupplierID(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, contacts)
		assert.Contains(t, err.Error(), "erro ao buscar")
		mockRepo.AssertExpectations(t)
	})
}

func TestContactService_Update(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(log)

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		userID := int64(1)
		existingContact := &models.Contact{
			ID:          1,
			ContactName: "Old Name",
			UserID:      &userID,
		}
		updatedContact := &models.Contact{
			ID:          1,
			ContactName: "New Name",
			UserID:      &userID,
		}

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(existingContact, nil).Once()
		mockRepo.On("Update", mock.Anything, updatedContact).Return(nil).Once()

		err := service.Update(context.Background(), updatedContact)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("contact not found", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		userID := int64(1)
		updatedContact := &models.Contact{
			ID:          1,
			ContactName: "New Name",
			UserID:      &userID,
		}

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return((*models.Contact)(nil), err_msg.ErrNotFound).Once()

		err := service.Update(context.Background(), updatedContact)
		assert.ErrorIs(t, err, err_msg.ErrNotFound)

		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("validation error - missing contact name", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		userID := int64(1)
		invalidContact := &models.Contact{
			ID:     1,
			UserID: &userID,
		}

		err := service.Update(context.Background(), invalidContact)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "campo obrigatório")

		mockRepo.AssertNotCalled(t, "GetByID")
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("validation error - invalid ID", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		invalidContact := &models.Contact{ID: 0, ContactName: "Válido"}

		err := service.Update(context.Background(), invalidContact)
		assert.ErrorIs(t, err, err_msg.ErrID)
		mockRepo.AssertNotCalled(t, "GetByID")
		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("repository error on GetByID", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
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
		assert.ErrorContains(t, err, "erro ao buscar")

		mockRepo.AssertNotCalled(t, "Update")
	})

	t.Run("repository error on Update", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		userID := int64(1)
		contact := &models.Contact{
			ID:          1,
			ContactName: "Válido",
			UserID:      &userID,
		}

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(contact, nil).Once()
		mockRepo.On("Update", mock.Anything, contact).Return(errors.New("erro ao atualizar")).Once()

		err := service.Update(context.Background(), contact)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "erro ao atualizar")

		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteContact(t *testing.T) {
	for name, tt := range map[string]func(t *testing.T){
		"success": func(t *testing.T) {
			mockRepo := new(mock_contact.MockContactRepository)
			log := logrus.New()
			log.Out = &bytes.Buffer{}
			logger := logger.NewLoggerAdapter(log)
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
			mockRepo := new(mock_contact.MockContactRepository)
			log := logrus.New()
			log.Out = &bytes.Buffer{}
			logger := logger.NewLoggerAdapter(log)
			service := NewContactService(mockRepo, logger)

			mockRepo.On("GetByID", mock.Anything, int64(1)).
				Return((*models.Contact)(nil), err_msg.ErrNotFound)

			err := service.Delete(context.Background(), 1)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "não encontrado")

			mockRepo.AssertNotCalled(t, "Delete")
		},
		"invalid id": func(t *testing.T) {
			mockRepo := new(mock_contact.MockContactRepository)
			log := logrus.New()
			log.Out = &bytes.Buffer{}
			logger := logger.NewLoggerAdapter(log)
			service := NewContactService(mockRepo, logger)

			err := service.Delete(context.Background(), 0)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "erro ID")

			mockRepo.AssertNotCalled(t, "GetByID")
			mockRepo.AssertNotCalled(t, "Delete")
		},
		"repository error on GetByID": func(t *testing.T) {
			mockRepo := new(mock_contact.MockContactRepository)
			log := logrus.New()
			log.Out = &bytes.Buffer{}
			logger := logger.NewLoggerAdapter(log)
			service := NewContactService(mockRepo, logger)

			mockRepo.On("GetByID", mock.Anything, int64(1)).
				Return((*models.Contact)(nil), errors.New("erro inesperado"))

			err := service.Delete(context.Background(), 1)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "erro ao buscar")

			mockRepo.AssertNotCalled(t, "Delete")
		},
		"repository error on Delete": func(t *testing.T) {
			mockRepo := new(mock_contact.MockContactRepository)
			log := logrus.New()
			log.Out = &bytes.Buffer{}
			logger := logger.NewLoggerAdapter(log)
			service := NewContactService(mockRepo, logger)

			existingContact := &models.Contact{
				ID:          1,
				ContactName: "Test User",
			}
			mockRepo.On("GetByID", mock.Anything, int64(1)).Return(existingContact, nil)
			mockRepo.On("Delete", mock.Anything, int64(1)).Return(errors.New("falha ao deletar"))

			err := service.Delete(context.Background(), 1)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "erro ao deletar")

			mockRepo.AssertExpectations(t)
		},
	} {
		t.Run(name, tt)
	}
}
