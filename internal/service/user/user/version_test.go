package services

import (
	"context"
	"fmt"
	"testing"

	mockUser "github.com/WagaoCarvalho/backend_store_go/infra/mock/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_GetVersionByID(t *testing.T) {

	setup := func() (
		*mockUser.MockUser,
		*MockHasher,
		User,
	) {
		mockUserRepo := new(mockUser.MockUser)
		mockHasher := new(MockHasher)

		userService := NewUserService(
			mockUserRepo,

			mockHasher,
		)

		return mockUserRepo, mockHasher, userService
	}

	t.Run("falha: ID inválido", func(t *testing.T) {
		_, _, service := setup()

		version, err := service.GetVersionByID(context.Background(), 0)

		assert.Equal(t, int64(0), version)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("Deve retornar versão quando usuário for encontrado", func(t *testing.T) {
		mockRepo, _, service := setup()

		mockRepo.On("GetVersionByID", mock.Anything, int64(1)).Return(int64(5), nil)

		version, err := service.GetVersionByID(context.Background(), 1)

		assert.NoError(t, err)
		assert.Equal(t, int64(5), version)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro de usuário não encontrado", func(t *testing.T) {
		mockRepo, _, service := setup()

		mockRepo.On("GetVersionByID", mock.Anything, int64(999)).Return(
			int64(0),
			errMsg.ErrNotFound,
		)

		version, err := service.GetVersionByID(context.Background(), 999)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		assert.Equal(t, int64(0), version)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro genérico quando falhar no repositório", func(t *testing.T) {
		mockRepo, _, service := setup()

		mockRepo.On("GetVersionByID", mock.Anything, int64(2)).Return(
			int64(0),
			fmt.Errorf("erro no banco de dados"),
		)

		version, err := service.GetVersionByID(context.Background(), 2)

		assert.ErrorContains(t, err, "conflito de versão")
		assert.Equal(t, int64(0), version)
		mockRepo.AssertExpectations(t)
	})
}
