package services

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/WagaoCarvalho/backend_store_go/internal/models"
	"github.com/WagaoCarvalho/backend_store_go/internal/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock do UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUsers(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) GetUserById(ctx context.Context, uid int64) (models.User, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(models.User), args.Error(1)
}

func TestUserService_GetUsers(t *testing.T) {
	// Criar o mock do repositório
	mockRepo := new(MockUserRepository)

	// Dados simulados que o mock deve retornar
	expectedUsers := []models.User{
		{
			UID:       1,
			Username:  "user1",
			Email:     "user1@example.com",
			Password:  "hash1",
			Status:    true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			UID:       2,
			Username:  "user2",
			Email:     "user2@example.com",
			Password:  "hash2",
			Status:    false,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	// Configurar o mock para retornar os usuários simulados quando chamado
	mockRepo.On("GetUsers", mock.Anything).Return(expectedUsers, nil)

	// Criar a instância do serviço passando o mock do repositório
	userService := services.NewUserService(mockRepo)

	// Chamar o método que estamos testando
	users, err := userService.GetUsers(context.Background())

	// Verificar se não houve erro
	assert.NoError(t, err)

	// Verificar se os usuários retornados são os esperados
	assert.Equal(t, expectedUsers, users)

	// Verificar se o método GetUsers foi chamado corretamente no mock
	mockRepo.AssertCalled(t, "GetUsers", mock.Anything)
}

func TestUserService_GetUserById(t *testing.T) {
	// Criar o mock do repositório
	mockRepo := new(MockUserRepository)

	// Dados simulados que o mock deve retornar
	expectedUser := models.User{
		UID:       1,
		Username:  "user1",
		Email:     "user1@example.com",
		Password:  "hash1",
		Status:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Configurar o mock para retornar o usuário simulado quando GetUserById for chamado com o ID 1
	mockRepo.On("GetUserById", mock.Anything, int64(1)).Return(expectedUser, nil)

	// Criar a instância do serviço passando o mock do repositório
	userService := services.NewUserService(mockRepo)

	// Chamar o método que estamos testando
	user, err := userService.GetUserById(context.Background(), 1)

	// Verificar se não houve erro
	assert.NoError(t, err)

	// Verificar se o usuário retornado é o esperado
	assert.Equal(t, expectedUser, user)

	// Verificar se o método GetUserById foi chamado corretamente no mock
	mockRepo.AssertCalled(t, "GetUserById", mock.Anything, int64(1))
}

func TestUserService_GetUserById_UserNotFound(t *testing.T) {
	// Criar o mock do repositório
	mockRepo := new(MockUserRepository)

	// Configurar o mock para retornar erro quando GetUserById for chamado com um ID que não existe
	mockRepo.On("GetUserById", mock.Anything, int64(999)).Return(models.User{}, fmt.Errorf("usuário não encontrado"))

	// Criar a instância do serviço passando o mock do repositório
	userService := services.NewUserService(mockRepo)

	// Chamar o método que estamos testando
	user, err := userService.GetUserById(context.Background(), 999)

	// Verificar se houve erro
	assert.Error(t, err)

	// Verificar se o usuário retornado está vazio
	assert.Equal(t, models.User{}, user)

	// Verificar se o método GetUserById foi chamado corretamente no mock
	mockRepo.AssertCalled(t, "GetUserById", mock.Anything, int64(999))
}
