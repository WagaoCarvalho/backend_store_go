package services

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	mock_contact "github.com/WagaoCarvalho/backend_store_go/infra/mock/repo/contact"
	dtoContact "github.com/WagaoCarvalho/backend_store_go/internal/dto/contact"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestContactService_Create(t *testing.T) {
	t.Run("sucesso na criação do contato", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		log := logrus.New()
		log.Out = &bytes.Buffer{}
		logger := logger.NewLoggerAdapter(log)

		service := NewContactService(mockRepo, logger)

		userID := int64(1)
		contactDTO := &dtoContact.ContactDTO{
			UserID:      &userID,
			ContactName: "Contato Teste",
			Email:       "teste@email.com",
			Phone:       "1234567898",
		}

		contactModel := dtoContact.ToContactModel(*contactDTO)

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(m *model.Contact) bool {
			return m.UserID != nil &&
				*m.UserID == *contactDTO.UserID &&
				m.ContactName == contactDTO.ContactName &&
				m.Email == contactDTO.Email &&
				m.Phone == contactDTO.Phone
		})).Return(contactModel, nil)

		createdContact, err := service.Create(context.Background(), contactDTO)

		assert.NoError(t, err)
		assert.Equal(t, contactDTO.ContactName, createdContact.ContactName)
		assert.Equal(t, contactDTO.Email, createdContact.Email)
		assert.Equal(t, contactDTO.Phone, createdContact.Phone)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha na validação do contato ContactName obrigatório", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		log := logrus.New()
		log.Out = &bytes.Buffer{}
		logger := logger.NewLoggerAdapter(log)

		service := NewContactService(mockRepo, logger)

		// Preenchendo UserID, ClientID e SupplierID para testar apenas ContactName
		userID := int64(1)
		clientID := int64(1)
		supplierID := int64(1)

		contactDTO := &dtoContact.ContactDTO{
			UserID:     &userID,
			ClientID:   &clientID,
			SupplierID: &supplierID,
			Email:      "teste@email.com",
			Phone:      "1234567898",
			// ContactName deixado vazio para gerar erro
		}

		createdContact, err := service.Create(context.Background(), contactDTO)

		assert.Nil(t, createdContact)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "erro no campo 'contact_name'") // mensagem real do Validate
		mockRepo.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("falha no repositório ao criar contato", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		log := logrus.New()
		log.Out = &bytes.Buffer{}
		logger := logger.NewLoggerAdapter(log)

		service := NewContactService(mockRepo, logger)

		userID := int64(1)
		contactDTO := &dtoContact.ContactDTO{
			UserID:      &userID,
			ContactName: "Contato Teste",
			Email:       "teste@email.com",
			Phone:       "1234567898",
		}

		expectedErr := errors.New("erro no banco")

		mockRepo.On("Create", mock.Anything, mock.MatchedBy(func(m *model.Contact) bool {
			return m.UserID != nil &&
				*m.UserID == *contactDTO.UserID &&
				m.ContactName == contactDTO.ContactName &&
				m.Email == contactDTO.Email &&
				m.Phone == contactDTO.Phone
		})).Return((*model.Contact)(nil), expectedErr)

		createdContact, err := service.Create(context.Background(), contactDTO)

		assert.Nil(t, createdContact)
		assert.Equal(t, expectedErr, err)
		mockRepo.AssertExpectations(t)
	})
}

func Test_GetByID(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(log)

	t.Run("success", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		expectedContact := &model.Contact{
			ID:          1,
			ContactName: "Test User",
			Email:       "test@example.com",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockRepo.On("GetByID", mock.Anything, int64(1)).Return(expectedContact, nil)

		contact, err := service.GetByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.NotNil(t, contact)
		assert.Equal(t, expectedContact.ID, *contact.ID)
		assert.Equal(t, expectedContact.ContactName, contact.ContactName)
		assert.Equal(t, expectedContact.Email, contact.Email)
		mockRepo.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		mockRepo.On("GetByID", mock.Anything, int64(1)).
			Return((*model.Contact)(nil), errMsg.ErrNotFound)

		contact, err := service.GetByID(context.Background(), 1)

		assert.Error(t, err)
		assert.Nil(t, contact)
		assert.EqualError(t, err, errMsg.ErrNotFound.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid id", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		service := NewContactService(mockRepo, logger)

		contact, err := service.GetByID(context.Background(), 0)

		assert.Error(t, err)
		assert.Nil(t, contact)
		assert.EqualError(t, err, errMsg.ErrID.Error())
		mockRepo.AssertNotCalled(t, "GetByID")
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		service := NewContactService(mockRepo, logger)

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
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	mockLogger := logger.NewLoggerAdapter(log)
	mockRepo := new(mock_contact.MockContactRepository)
	service := NewContactService(mockRepo, mockLogger)

	t.Run("sucesso ao buscar contatos por UserID", func(t *testing.T) {
		contactModels := []*models.Contact{
			{
				ID:          1,
				UserID:      utils.Int64Ptr(1),
				ContactName: "Contato Teste",
				Email:       "teste@email.com",
				Phone:       "123456789",
			},
		}

		mockRepo.On("GetByUserID", mock.Anything, int64(1)).Return(contactModels, nil)

		result, err := service.GetByUserID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1), *result[0].ID)
		assert.Equal(t, "Contato Teste", result[0].ContactName)
		assert.Equal(t, "teste@email.com", result[0].Email)
		mockRepo.AssertExpectations(t)

		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil
	})

	t.Run("falha ao buscar contatos com UserID inválido", func(t *testing.T) {
		service := NewContactService(nil, mockLogger)

		result, err := service.GetByUserID(context.Background(), 0)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, errMsg.ErrID.Error())
	})

	t.Run("nenhum contato encontrado por UserID", func(t *testing.T) {
		mockRepo.On("GetByUserID", mock.Anything, int64(2)).
			Return(([]*models.Contact)(nil), errMsg.ErrNotFound)

		result, err := service.GetByUserID(context.Background(), 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, errMsg.ErrNotFound.Error())

		mockRepo.AssertExpectations(t)
	})

}

func TestContactService_GetByClientID(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	mockLogger := logger.NewLoggerAdapter(log)
	mockRepo := new(mock_contact.MockContactRepository)
	service := NewContactService(mockRepo, mockLogger)

	t.Run("sucesso ao buscar contatos por ClientID", func(t *testing.T) {
		contactModels := []*models.Contact{
			{
				ID:          1,
				ClientID:    utils.Int64Ptr(1),
				ContactName: "Cliente Teste",
				Email:       "cliente@email.com",
				Phone:       "987654321",
			},
		}

		mockRepo.On("GetByClientID", mock.Anything, int64(1)).Return(contactModels, nil)

		result, err := service.GetByClientID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1), *result[0].ID)
		assert.Equal(t, "Cliente Teste", result[0].ContactName)
		assert.Equal(t, "cliente@email.com", result[0].Email)
		mockRepo.AssertExpectations(t)

		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil
	})

	t.Run("falha ao buscar contatos com ClientID inválido", func(t *testing.T) {
		service := NewContactService(nil, mockLogger)

		result, err := service.GetByClientID(context.Background(), 0)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, errMsg.ErrID.Error())
	})

	t.Run("nenhum contato encontrado por ClientID", func(t *testing.T) {
		mockRepo.On("GetByClientID", mock.Anything, int64(2)).
			Return(([]*models.Contact)(nil), errMsg.ErrNotFound)

		result, err := service.GetByClientID(context.Background(), 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, errMsg.ErrNotFound.Error())

		mockRepo.AssertExpectations(t)
	})

}

func TestContactService_GetBySupplierID(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	mockLogger := logger.NewLoggerAdapter(log)
	mockRepo := new(mock_contact.MockContactRepository)
	service := NewContactService(mockRepo, mockLogger)

	t.Run("sucesso ao buscar contatos por SupplierID", func(t *testing.T) {
		contactModels := []*models.Contact{
			{
				ID:          1,
				SupplierID:  utils.Int64Ptr(1),
				ContactName: "Fornecedor Teste",
				Email:       "fornecedor@email.com",
				Phone:       "1122334455",
			},
		}

		mockRepo.On("GetBySupplierID", mock.Anything, int64(1)).Return(contactModels, nil)

		result, err := service.GetBySupplierID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, int64(1), *result[0].ID)
		assert.Equal(t, "Fornecedor Teste", result[0].ContactName)
		assert.Equal(t, "fornecedor@email.com", result[0].Email)
		mockRepo.AssertExpectations(t)

		mockRepo.ExpectedCalls = nil
		mockRepo.Calls = nil
	})

	t.Run("falha ao buscar contatos com SupplierID inválido", func(t *testing.T) {
		service := NewContactService(nil, mockLogger)

		result, err := service.GetBySupplierID(context.Background(), 0)

		assert.Nil(t, result)
		assert.Error(t, err)
		assert.EqualError(t, err, errMsg.ErrID.Error())
	})

	t.Run("nenhum contato encontrado por SupplierID", func(t *testing.T) {
		mockRepo.On("GetBySupplierID", mock.Anything, int64(2)).
			Return(([]*models.Contact)(nil), errMsg.ErrNotFound)

		result, err := service.GetBySupplierID(context.Background(), 2)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.EqualError(t, err, errMsg.ErrNotFound.Error())

		mockRepo.AssertExpectations(t)
	})
}

func TestContactService_UpdateContact(t *testing.T) {
	makeContactDTO := func() dtoContact.ContactDTO {
		userID := int64(1)
		return dtoContact.ContactDTO{
			ID:          utils.Int64Ptr(1),
			UserID:      &userID,
			ContactName: "Contato Teste",
			Email:       "teste@email.com",
			Phone:       "1234567898",
		}
	}

	t.Run("sucesso na atualização do contato", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		log := logrus.New()
		log.Out = &bytes.Buffer{}
		logger := logger.NewLoggerAdapter(log)
		service := NewContactService(mockRepo, logger)

		contactDTO := makeContactDTO()
		contactModel := dtoContact.ToContactModel(contactDTO)

		mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(c *model.Contact) bool {
			return c != nil && c.ID == contactModel.ID && *c.UserID == *contactModel.UserID
		})).Return(nil)

		err := service.Update(context.Background(), &contactDTO)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("erro ao atualizar contato com ID inválido", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		log := logrus.New()
		log.Out = &bytes.Buffer{}
		logger := logger.NewLoggerAdapter(log)
		service := NewContactService(mockRepo, logger)

		userID := int64(1)
		contactDTO := &dtoContact.ContactDTO{
			UserID:      &userID,
			ContactName: "Nome Teste",
			Email:       "teste@email.com",
		}

		err := service.Update(context.Background(), contactDTO)
		assert.ErrorIs(t, err, errMsg.ErrID)
	})

	t.Run("falha na validação do contato no update", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		log := logrus.New()
		log.Out = &bytes.Buffer{}
		logger := logger.NewLoggerAdapter(log)
		service := NewContactService(mockRepo, logger)

		contactDTO := &dtoContact.ContactDTO{
			ID:    utils.Int64Ptr(1),
			Email: "teste@email.com",
		}

		err := service.Update(context.Background(), contactDTO)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "contact_name") // <- ajustar aqui
		mockRepo.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("erro genérico ao atualizar contato", func(t *testing.T) {
		mockRepo := new(mock_contact.MockContactRepository)
		log := logrus.New()
		log.Out = &bytes.Buffer{}
		logger := logger.NewLoggerAdapter(log)
		service := NewContactService(mockRepo, logger)

		contactDTO := &dtoContact.ContactDTO{
			ID:          utils.Int64Ptr(1),
			UserID:      utils.Int64Ptr(1),
			ContactName: "Nome Teste",
			Email:       "teste@email.com",
			Phone:       "1112345678",
		}

		mockRepo.On("Update", mock.Anything, mock.MatchedBy(func(c *model.Contact) bool {
			return c.ID == *contactDTO.ID
		})).Return(fmt.Errorf("erro inesperado no banco"))

		err := service.Update(context.Background(), contactDTO)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "erro ao atualizar")
		assert.ErrorContains(t, err, "erro inesperado no banco")
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

			existingContact := &model.Contact{
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
				Return((*model.Contact)(nil), errMsg.ErrNotFound)

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
				Return((*model.Contact)(nil), errors.New("erro inesperado"))

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

			existingContact := &model.Contact{
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
