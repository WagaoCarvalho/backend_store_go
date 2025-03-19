package repositories

import (
	"context"
	"testing"
	"time"

	"github.com/WagaoCarvalho/backend_store_go/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetUsers(t *testing.T) {
	// Cria uma instância do mock
	mockRepo := new(MockUserRepository)

	// Define os dados que serão retornados pelo mock
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

	// Configura o mock para retornar os dados esperados
	mockRepo.On("GetUsers", mock.Anything).Return(expectedUsers, nil)

	// Chama o método que queremos testar
	users, err := mockRepo.GetUsers(context.Background())

	// Verifica se o erro é nulo
	assert.NoError(t, err)

	// Verifica se os usuários retornados são os esperados
	assert.Equal(t, expectedUsers, users)

	// Verifica se o método foi chamado com os parâmetros corretos
	mockRepo.AssertCalled(t, "GetUsers", mock.Anything)
}

func TestGetUserById(t *testing.T) {
	// Cria uma instância do mock do repositório
	mockRepo := new(MockUserRepository)

	// Define os dados que serão retornados pelo mock para o GetUserById
	expectedUser := models.User{
		UID:       1,
		Username:  "user1",
		Email:     "user1@example.com",
		Password:  "hash1",
		Status:    true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Configura o mock para retornar o usuário esperado quando GetUserById for chamado
	mockRepo.On("GetUserById", mock.Anything, int64(1)).Return(expectedUser, nil)

	// Chama o método GetUserById para testar
	user, err := mockRepo.GetUserById(context.Background(), 1)

	// Verifica se não houve erro
	assert.NoError(t, err)

	// Verifica se o usuário retornado é o esperado
	assert.Equal(t, expectedUser, user)

	// Verifica se o método GetUserById foi chamado com o uid correto
	mockRepo.AssertCalled(t, "GetUserById", mock.Anything, int64(1))
}
