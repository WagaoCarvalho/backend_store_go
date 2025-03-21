package services

import (
	"context"
	"errors"
	"testing"

	"github.com/WagaoCarvalho/backend_store_go/internal/models"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/login"
	services_test "github.com/WagaoCarvalho/backend_store_go/tests/services/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginService_Login_Success(t *testing.T) {
	mockRepo := new(services_test.MockUserRepository)
	service := services.NewLoginService(mockRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)

	user := models.User{
		UID:      1,
		Email:    "teste@email.com",
		Password: string(hashedPassword),
	}

	mockRepo.On("GetUserByEmail", mock.Anything, "teste@email.com").Return(user, nil)

	token, err := service.Login(context.Background(), models.LoginCredentials{
		Email:    "teste@email.com",
		Password: "senha123",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestLoginService_Login_Failure_WrongPassword(t *testing.T) {
	mockRepo := new(services_test.MockUserRepository)
	service := services.NewLoginService(mockRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("outrasenha"), bcrypt.DefaultCost)

	user := models.User{
		UID:      1,
		Email:    "teste@email.com",
		Password: string(hashedPassword),
	}

	mockRepo.On("GetUserByEmail", mock.Anything, "teste@email.com").Return(user, nil)

	token, err := service.Login(context.Background(), models.LoginCredentials{
		Email:    "teste@email.com",
		Password: "senhaErrada",
	})

	assert.Error(t, err)
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestLoginService_Login_Failure_UserNotFound(t *testing.T) {
	mockRepo := new(services_test.MockUserRepository)
	service := services.NewLoginService(mockRepo)

	mockRepo.On("GetUserByEmail", mock.Anything, "naoexiste@email.com").Return(models.User{}, errors.New("usuário não encontrado"))

	token, err := service.Login(context.Background(), models.LoginCredentials{
		Email:    "naoexiste@email.com",
		Password: "senha123",
	})

	assert.Error(t, err)
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}
