package services

import (
	"context"
	"errors"
	"testing"
	"time"

	mockUser "github.com/WagaoCarvalho/backend_store_go/infra/mock/user"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/common/filter"
	userFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/user/filter"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserService_Filter(t *testing.T) {

	t.Run("falha quando filtro é nulo", func(t *testing.T) {
		mockRepo := new(mockUser.MockUser)
		service := NewUserFilterService(mockRepo)

		result, err := service.Filter(context.Background(), nil)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidFilter)
		mockRepo.AssertNotCalled(t, "Filter", mock.Anything, mock.Anything)
	})

	t.Run("falha na validação do filtro", func(t *testing.T) {
		mockRepo := new(mockUser.MockUser)
		service := NewUserFilterService(mockRepo)

		invalidFilter := &userFilter.UserFilter{
			BaseFilter: filter.BaseFilter{
				Limit: -1, // inválido
			},
		}

		result, err := service.Filter(context.Background(), invalidFilter)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrInvalidFilter)
		mockRepo.AssertNotCalled(t, "Filter", mock.Anything, mock.Anything)
	})

	t.Run("falha ao buscar no repositório", func(t *testing.T) {
		mockRepo := new(mockUser.MockUser)
		service := NewUserFilterService(mockRepo)

		validFilter := &userFilter.UserFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		dbErr := errors.New("falha no banco de dados")

		mockRepo.
			On("Filter", mock.Anything, validFilter).
			Return(nil, dbErr).
			Once()

		result, err := service.Filter(context.Background(), validFilter)

		assert.Nil(t, result)
		assert.ErrorIs(t, err, errMsg.ErrGet)
		assert.ErrorContains(t, err, dbErr.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso ao retornar lista de usuários", func(t *testing.T) {
		mockRepo := new(mockUser.MockUser)
		service := NewUserFilterService(mockRepo)

		validFilter := &userFilter.UserFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
		}

		now := time.Now()
		description1 := "Usuário administrador"
		description2 := ""

		mockUsers := []*model.User{
			{
				UID:         1,
				Username:    "admin",
				Email:       "admin@example.com",
				Password:    "hashed_password_1",
				Description: description1,
				Status:      true,
				Version:     1,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				UID:         2,
				Username:    "guest",
				Email:       "guest@example.com",
				Password:    "hashed_password_2",
				Description: description2,
				Status:      false,
				Version:     1,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}

		mockRepo.
			On("Filter", mock.Anything, validFilter).
			Return(mockUsers, nil).
			Once()

		result, err := service.Filter(context.Background(), validFilter)

		assert.NoError(t, err)
		assert.Len(t, result, 2)

		// Verificar primeiro usuário
		assert.Equal(t, int64(1), result[0].UID)
		assert.Equal(t, "admin", result[0].Username)
		assert.Equal(t, "admin@example.com", result[0].Email)
		assert.Equal(t, "hashed_password_1", result[0].Password)
		assert.NotNil(t, result[0].Description)
		assert.Equal(t, "Usuário administrador", result[0].Description)
		assert.True(t, result[0].Status)
		assert.Equal(t, 1, result[0].Version)
		assert.WithinDuration(t, now, result[0].CreatedAt, time.Second)
		assert.WithinDuration(t, now, result[0].UpdatedAt, time.Second)

		// Verificar segundo usuário
		assert.Equal(t, int64(2), result[1].UID)
		assert.Equal(t, "guest", result[1].Username)
		assert.Equal(t, "guest@example.com", result[1].Email)
		assert.Equal(t, "hashed_password_2", result[1].Password)
		assert.NotNil(t, result[1].Description)
		assert.Equal(t, "", result[1].Description)
		assert.False(t, result[1].Status)
		assert.Equal(t, 1, result[1].Version)
		assert.WithinDuration(t, now, result[1].CreatedAt, time.Second)
		assert.WithinDuration(t, now, result[1].UpdatedAt, time.Second)

		mockRepo.AssertExpectations(t)
	})

	t.Run("sucesso com filtros específicos de username e email", func(t *testing.T) {
		mockRepo := new(mockUser.MockUser)
		service := NewUserFilterService(mockRepo)

		validFilter := &userFilter.UserFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  5,
				Offset: 0,
			},
			Username: "john",
			Email:    "john@example.com",
		}

		now := time.Now()
		description := "Usuário específico"

		mockUsers := []*model.User{
			{
				UID:         10,
				Username:    "john_doe",
				Email:       "john@example.com",
				Password:    "hashed_password_john",
				Description: description,
				Status:      true,
				Version:     2,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}

		mockRepo.
			On("Filter", mock.Anything, validFilter).
			Return(mockUsers, nil).
			Once()

		result, err := service.Filter(context.Background(), validFilter)

		assert.NoError(t, err)
		assert.Len(t, result, 1)

		assert.Equal(t, int64(10), result[0].UID)
		assert.Equal(t, "john_doe", result[0].Username)
		assert.Equal(t, "john@example.com", result[0].Email)
		assert.Equal(t, "hashed_password_john", result[0].Password)
		assert.True(t, result[0].Status)
		assert.Equal(t, 2, result[0].Version)

		mockRepo.AssertExpectations(t)
	})

	t.Run("retorna lista vazia quando nenhum usuário encontrado", func(t *testing.T) {
		mockRepo := new(mockUser.MockUser)
		service := NewUserFilterService(mockRepo)

		validFilter := &userFilter.UserFilter{
			BaseFilter: filter.BaseFilter{
				Limit:  10,
				Offset: 0,
			},
			Username: "nonexistent",
		}

		mockRepo.
			On("Filter", mock.Anything, validFilter).
			Return([]*model.User{}, nil).
			Once()

		result, err := service.Filter(context.Background(), validFilter)

		assert.NoError(t, err)
		assert.Len(t, result, 0)
		assert.NotNil(t, result) // Verifica que não é nil

		mockRepo.AssertExpectations(t)
	})
}
