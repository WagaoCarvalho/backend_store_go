package services

import (
	"context"
	"errors"
	"testing"

	login "github.com/WagaoCarvalho/backend_store_go/internal/models"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockLoginService struct {
	mock.Mock
}

func (m *MockLoginService) Login(ctx context.Context, email string, password string) (string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}

func TestLoginService_Login_Success(t *testing.T) {
	mockRepo := new(services.MockUserRepository)
	service := NewLoginService(mockRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("senha123"), bcrypt.DefaultCost)

	user := models.User{
		UID:      1,
		Email:    "teste@email.com",
		Password: string(hashedPassword),
	}

	mockRepo.On("GetUserByEmail", mock.Anything, "teste@email.com").Return(user, nil)

	token, err := service.Login(context.Background(), login.LoginCredentials{
		Email:    "teste@email.com",
		Password: "senha123",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestLoginService_Login_Failure_WrongPassword(t *testing.T) {
	mockRepo := new(services.MockUserRepository)
	service := NewLoginService(mockRepo)

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("outrasenha"), bcrypt.DefaultCost)

	user := models.User{
		UID:      1,
		Email:    "teste@email.com",
		Password: string(hashedPassword),
	}

	mockRepo.On("GetUserByEmail", mock.Anything, "teste@email.com").Return(user, nil)

	token, err := service.Login(context.Background(), login.LoginCredentials{
		Email:    "teste@email.com",
		Password: "senhaErrada",
	})

	assert.Error(t, err)
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestLoginService_Login_Failure_UserNotFound(t *testing.T) {
	mockRepo := new(services.MockUserRepository)
	service := NewLoginService(mockRepo)

	mockRepo.On("GetUserByEmail", mock.Anything, "naoexiste@email.com").Return(models.User{}, errors.New("usuário não encontrado"))

	token, err := service.Login(context.Background(), login.LoginCredentials{
		Email:    "naoexiste@email.com",
		Password: "senha123",
	})

	assert.Error(t, err)
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}
