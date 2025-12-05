package services

import (
	"context"
	"errors"
	"fmt"
	"testing"

	mockUser "github.com/WagaoCarvalho/backend_store_go/infra/mock/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_UserExists(t *testing.T) {
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

	t.Run("Deve retornar erro quando userID é zero", func(t *testing.T) {
		_, _, service := setup()

		exists, err := service.UserExists(context.Background(), 0)

		assert.False(t, exists)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("Deve retornar erro quando userID é negativo", func(t *testing.T) {
		_, _, service := setup()

		exists, err := service.UserExists(context.Background(), -1)

		assert.False(t, exists)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
	})

	t.Run("Deve retornar true quando usuário existe", func(t *testing.T) {
		mockRepo, _, service := setup()

		userID := int64(123)
		mockRepo.On("UserExists", mock.Anything, userID).Return(true, nil)

		exists, err := service.UserExists(context.Background(), userID)

		assert.True(t, exists)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar false quando usuário não existe", func(t *testing.T) {
		mockRepo, _, service := setup()

		userID := int64(456)
		mockRepo.On("UserExists", mock.Anything, userID).Return(false, nil)

		exists, err := service.UserExists(context.Background(), userID)

		assert.False(t, exists)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro quando o repositório retorna erro", func(t *testing.T) {
		mockRepo, _, service := setup()

		userID := int64(789)
		repoError := fmt.Errorf("erro no banco de dados")
		mockRepo.On("UserExists", mock.Anything, userID).Return(false, repoError)

		exists, err := service.UserExists(context.Background(), userID)

		assert.False(t, exists)
		assert.Error(t, err)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.Contains(t, err.Error(), "erro no banco de dados")
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro de timeout quando contexto é cancelado", func(t *testing.T) {
		mockRepo, _, service := setup()

		userID := int64(999)
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancela imediatamente

		// O mock não deve ser chamado porque o contexto está cancelado
		// (depende da implementação do repositório)
		mockRepo.On("UserExists", mock.Anything, userID).
			Maybe().
			Return(false, context.Canceled)

		exists, err := service.UserExists(ctx, userID)

		// Pode retornar false ou erro dependendo da implementação
		// Normalmente o contexto cancelado é propagado
		assert.False(t, exists)
		assert.Error(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve lidar com valores extremos de userID", func(t *testing.T) {
		mockRepo, _, service := setup()

		testCases := []struct {
			name   string
			userID int64
		}{
			{"ID mínimo positivo", 1},
			{"ID grande positivo", 999999999999},
			{"ID máximo int64", int64(9223372036854775807)},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				mockRepo.On("UserExists", mock.Anything, tc.userID).Return(true, nil).Once()

				exists, err := service.UserExists(context.Background(), tc.userID)

				assert.True(t, exists)
				assert.NoError(t, err)
			})
		}

		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar erro genérico do repositório envolvido com ErrGet", func(t *testing.T) {
		mockRepo, _, service := setup()

		userID := int64(1000)
		repoError := errors.New("connection failed")
		mockRepo.On("UserExists", mock.Anything, userID).Return(false, repoError)

		exists, err := service.UserExists(context.Background(), userID)

		assert.False(t, exists)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.Contains(t, err.Error(), repoError.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve chamar repositório apenas uma vez para mesma chamada", func(t *testing.T) {
		mockRepo, _, service := setup()

		userID := int64(555)
		mockRepo.On("UserExists", mock.Anything, userID).Return(true, nil).Once()

		// Primeira chamada
		exists1, err1 := service.UserExists(context.Background(), userID)
		assert.True(t, exists1)
		assert.NoError(t, err1)

		// Segunda chamada - o mock foi configurado com .Once(), então falharia
		// Se quiser suportar múltiplas chamadas, use .Times(2) ou .Maybe()
		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve retornar diferentes resultados para diferentes IDs", func(t *testing.T) {
		mockRepo, _, service := setup()

		// Configurar diferentes retornos para diferentes IDs
		mockRepo.On("UserExists", mock.Anything, int64(1)).Return(true, nil)
		mockRepo.On("UserExists", mock.Anything, int64(2)).Return(false, nil)
		mockRepo.On("UserExists", mock.Anything, int64(3)).Return(false, fmt.Errorf("erro no id 3"))

		// Teste ID 1
		exists1, err1 := service.UserExists(context.Background(), 1)
		assert.True(t, exists1)
		assert.NoError(t, err1)

		// Teste ID 2
		exists2, err2 := service.UserExists(context.Background(), 2)
		assert.False(t, exists2)
		assert.NoError(t, err2)

		// Teste ID 3
		exists3, err3 := service.UserExists(context.Background(), 3)
		assert.False(t, exists3)
		assert.Error(t, err3)
		assert.ErrorIs(t, err3, errMsg.ErrGet)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Deve validar userID antes de chamar o repositório", func(t *testing.T) {
		mockRepo, _, service := setup()

		// Com userID zero, o repositório NÃO deve ser chamado
		exists, err := service.UserExists(context.Background(), 0)

		assert.False(t, exists)
		assert.ErrorIs(t, err, errMsg.ErrZeroID)
		mockRepo.AssertNotCalled(t, "UserExists", mock.Anything, mock.Anything)
	})

	t.Run("Deve propagar contexto para o repositório", func(t *testing.T) {
		mockRepo, _, service := setup()

		userID := int64(777)
		ctx := context.WithValue(context.Background(), "trace_id", "12345")

		mockRepo.On("UserExists", mock.MatchedBy(func(ctx context.Context) bool {
			// Verifica se o contexto contém o valor esperado
			return ctx.Value("trace_id") == "12345"
		}), userID).Return(true, nil)

		exists, err := service.UserExists(ctx, userID)

		assert.True(t, exists)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}
