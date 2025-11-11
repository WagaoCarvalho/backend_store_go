package services

import (
	"context"
	"errors"
	"testing"
	"time"

	mockContact "github.com/WagaoCarvalho/backend_store_go/infra/mock/contact"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestContactService_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mockRepo := new(mockContact.MockContact)
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
		mockRepo := new(mockContact.MockContact)
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
		mockRepo := new(mockContact.MockContact)
		service := NewContactService(mockRepo)

		contact, err := service.GetByID(context.Background(), 0)

		assert.Error(t, err)
		assert.Nil(t, contact)
		assert.EqualError(t, err, errMsg.ErrZeroID.Error())
		mockRepo.AssertNotCalled(t, "GetByID", mock.Anything, mock.Anything)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepo := new(mockContact.MockContact)
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
