package services

import (
	"context"
	"fmt"
	"testing"

	mockUser "github.com/WagaoCarvalho/backend_store_go/infra/mock/user"
	modelUser "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_Disable(t *testing.T) {

	setup := func() (
		*mockUser.MockUser,
		*MockHasher,
		User,
	) {
		mockUserRepo := new(mockUser.MockUser)
		mockHasher := new(MockHasher)

		userService := NewUserService(mockUserRepo, mockHasher)

		return mockUserRepo, mockHasher, userService
	}

	t.Run("falha: ID inválido", func(t *testing.T) {
		_, _, service := setup()

		err := service.Disable(context.Background(), 0)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("falha: usuário não encontrado", func(t *testing.T) {
		mockRepo, _, service := setup()

		mockRepo.On("GetByID", mock.Anything, int64(1)).
			Return(nil, errMsg.ErrNotFound).Once()

		err := service.Disable(context.Background(), 1)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso: usuário já desativado (não deve chamar Disable)", func(t *testing.T) {
		mockRepo, _, service := setup()

		user := &modelUser.User{UID: 1, Status: false}

		mockRepo.On("GetByID", mock.Anything, int64(1)).
			Return(user, nil).Once()

		err := service.Disable(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Disable", mock.Anything, int64(1))
	})

	t.Run("sucesso: desativa usuário com sucesso", func(t *testing.T) {
		mockRepo, _, service := setup()

		user := &modelUser.User{UID: 1, Status: true}

		mockRepo.On("GetByID", mock.Anything, int64(1)).
			Return(user, nil).Once()
		mockRepo.On("Disable", mock.Anything, int64(1)).
			Return(nil).Once()

		err := service.Disable(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha: erro inesperado ao buscar usuário", func(t *testing.T) {
		mockRepo, _, service := setup()

		mockRepo.On("GetByID", mock.Anything, int64(1)).
			Return(nil, fmt.Errorf("erro no banco")).Once()

		err := service.Disable(context.Background(), 1)

		assert.ErrorContains(t, err, "erro ao buscar")
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha: erro ao desativar usuário", func(t *testing.T) {
		mockRepo, _, service := setup()

		user := &modelUser.User{UID: 1, Status: true}

		mockRepo.On("GetByID", mock.Anything, int64(1)).
			Return(user, nil).Once()
		mockRepo.On("Disable", mock.Anything, int64(1)).
			Return(fmt.Errorf("erro ao desativar")).Once()

		err := service.Disable(context.Background(), 1)

		assert.ErrorContains(t, err, "erro ao desativar")
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_Enable(t *testing.T) {

	setup := func() (
		*mockUser.MockUser,
		*MockHasher,
		User,
	) {
		mockUserRepo := new(mockUser.MockUser)
		mockHasher := new(MockHasher)

		userService := NewUserService(mockUserRepo, mockHasher)

		return mockUserRepo, mockHasher, userService
	}

	t.Run("falha: ID inválido", func(t *testing.T) {
		_, _, service := setup()

		err := service.Enable(context.Background(), 0)

		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("falha: usuário não encontrado", func(t *testing.T) {
		mockRepo, _, service := setup()

		mockRepo.On("GetByID", mock.Anything, int64(99)).
			Return(nil, errMsg.ErrNotFound).Once()

		err := service.Enable(context.Background(), 99)

		assert.ErrorIs(t, err, errMsg.ErrNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso: usuário já ativo (não deve chamar Enable)", func(t *testing.T) {
		mockRepo, _, service := setup()

		user := &modelUser.User{UID: 1, Status: true}

		mockRepo.On("GetByID", mock.Anything, int64(1)).
			Return(user, nil).Once()

		err := service.Enable(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertNotCalled(t, "Enable", mock.Anything, int64(1))
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso: ativa usuário com sucesso", func(t *testing.T) {
		mockRepo, _, service := setup()

		user := &modelUser.User{UID: 1, Status: false}

		mockRepo.On("GetByID", mock.Anything, int64(1)).
			Return(user, nil).Once()
		mockRepo.On("Enable", mock.Anything, int64(1)).
			Return(nil).Once()

		err := service.Enable(context.Background(), 1)

		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha: erro inesperado ao buscar usuário", func(t *testing.T) {
		mockRepo, _, service := setup()

		mockRepo.On("GetByID", mock.Anything, int64(1)).
			Return(nil, fmt.Errorf("falha na consulta")).Once()

		err := service.Enable(context.Background(), 1)

		assert.ErrorContains(t, err, "falha na consulta")
		mockRepo.AssertExpectations(t)
	})

	t.Run("falha: erro ao ativar usuário", func(t *testing.T) {
		mockRepo, _, service := setup()

		user := &modelUser.User{UID: 1, Status: false}

		mockRepo.On("GetByID", mock.Anything, int64(1)).
			Return(user, nil).Once()
		mockRepo.On("Enable", mock.Anything, int64(1)).
			Return(fmt.Errorf("erro ao ativar")).Once()

		err := service.Enable(context.Background(), 1)

		assert.ErrorContains(t, err, "erro ao ativar")
		mockRepo.AssertExpectations(t)
	})
}
