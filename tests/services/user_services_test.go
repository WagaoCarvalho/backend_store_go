package services

import (
	"context"
	"os"
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

func init() {
	// Define o ambiente como "test"
	os.Setenv("GO_ENV", "test")
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
