package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/WagaoCarvalho/backend_store_go/internal/handlers"
	"github.com/WagaoCarvalho/backend_store_go/internal/models"
	"github.com/WagaoCarvalho/backend_store_go/utils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockUserService struct {
	mock.Mock
}

func (m *mockUserService) GetUsers(ctx context.Context) ([]models.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *mockUserService) GetUserById(ctx context.Context, uid int64) (models.User, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *mockUserService) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(models.User), args.Error(1)
}

func (m *mockUserService) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(models.User), args.Error(1)
}

func TestGetUsers(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := handlers.NewUserHandler(mockService)

	r := mux.NewRouter()
	r.HandleFunc("/users", userHandler.GetUsers).Methods("GET")

	// Simulando a resposta do mock
	mockService.On("GetUsers", mock.Anything).Return([]models.User{
		{
			UID:      1,
			Username: "user1",
			Email:    "user1@example.com",
			Status:   true,
		},
		{
			UID:      2,
			Username: "user2",
			Email:    "user2@example.com",
			Status:   false,
		},
	}, nil)

	req := httptest.NewRequest("GET", "/users", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response utils.DefaultResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Status)

	users := response.Data.([]interface{})
	assert.Equal(t, 2, len(users))

	assert.Equal(t, float64(1), users[0].(map[string]interface{})["uid"])
	assert.Equal(t, "user1", users[0].(map[string]interface{})["nickname"])
	assert.Equal(t, "user1@example.com", users[0].(map[string]interface{})["email"])

	assert.Equal(t, float64(2), users[1].(map[string]interface{})["uid"])
	assert.Equal(t, "user2", users[1].(map[string]interface{})["nickname"])
	assert.Equal(t, "user2@example.com", users[1].(map[string]interface{})["email"])
}

func TestGetUserById(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := handlers.NewUserHandler(mockService)

	r := mux.NewRouter()
	r.HandleFunc("/user/{id}", userHandler.GetUserById).Methods("GET")

	mockService.On("GetUserById", mock.Anything, int64(1)).Return(models.User{
		UID:      1,
		Username: "user1",
		Email:    "user1@example.com",
		Status:   true,
	}, nil)

	req := httptest.NewRequest("GET", "/user/1", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response utils.DefaultResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Status)

	user := response.Data.(map[string]interface{})
	assert.Equal(t, float64(1), user["uid"])
	assert.Equal(t, "user1", user["nickname"])
	assert.Equal(t, "user1@example.com", user["email"])
}

func TestGetUserById_UserNotFound(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := handlers.NewUserHandler(mockService)

	r := mux.NewRouter()
	r.HandleFunc("/user/{id}", userHandler.GetUserById).Methods("GET")

	mockService.On("GetUserById", mock.Anything, int64(999)).Return(models.User{}, fmt.Errorf("usuário não encontrado"))

	req := httptest.NewRequest("GET", "/user/999", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var response utils.DefaultResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, response.Status)

	assert.Contains(t, response.Message, "usuário não encontrado")
}

func TestCreateUser(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := handlers.NewUserHandler(mockService)

	// Definindo o usuário esperado com dados dinâmicos
	now := time.Now()
	expectedUser := models.User{
		UID:       1,
		Username:  "user1",
		Email:     "user1@example.com",
		Password:  "hashedpassword",
		Status:    true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Configurando o mock para retornar o usuário criado
	mockService.On("CreateUser", mock.Anything, mock.MatchedBy(func(user models.User) bool {
		// Verifica os campos principais
		return user.Username == "user1" && user.Email == "user1@example.com" && user.Password == "hashedpassword" && user.Status == true
	})).Return(expectedUser, nil)

	// Criando um roteador e registrando a rota de CreateUser
	r := mux.NewRouter()
	r.HandleFunc("/user", userHandler.CreateUser).Methods("POST")

	// Criando o corpo da requisição com os dados do usuário a ser criado
	userData := `{
		"username": "user1",
		"email": "user1@example.com",
		"password": "hashedpassword",
		"status": true
	}`
	req := httptest.NewRequest("POST", "/user", strings.NewReader(userData))
	rr := httptest.NewRecorder()

	// Chamando o handler para a requisição simulada
	r.ServeHTTP(rr, req)

	// Verificando o código de status HTTP
	assert.Equal(t, http.StatusCreated, rr.Code)

	// Verificando o corpo da resposta
	var response utils.DefaultResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.Status)

	// Verificando os dados retornados
	user := response.Data.(map[string]interface{})
	assert.Equal(t, float64(1), user["uid"])
	assert.Equal(t, "user1", user["username"])
	assert.Equal(t, "user1@example.com", user["email"])

	// Verificando se o método CreateUser foi chamado no mock
	mockService.AssertCalled(t, "CreateUser", mock.Anything, mock.MatchedBy(func(user models.User) bool {
		return user.Username == "user1" && user.Email == "user1@example.com" && user.Password == "hashedpassword" && user.Status == true
	}))
}

func TestCreateUser_Failure(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := handlers.NewUserHandler(mockService)

	// Definindo o usuário a ser criado
	newUser := models.User{
		Username: "user2",
		Email:    "user2@example.com",
		Password: "hashedpassword",
		Status:   true,
	}

	// Configurando o mock para simular erro no serviço
	mockService.On("CreateUser", mock.Anything, newUser).Return(models.User{}, fmt.Errorf("erro ao criar usuário"))

	// Criando um roteador e registrando a rota de CreateUser
	r := mux.NewRouter()
	r.HandleFunc("/user", userHandler.CreateUser).Methods("POST")

	// Criando o corpo da requisição com os dados do usuário a ser criado
	userData := `{
		"username": "user2",
		"email": "user2@example.com",
		"password": "hashedpassword",
		"status": true
	}`
	req := httptest.NewRequest("POST", "/user", strings.NewReader(userData))
	rr := httptest.NewRecorder()

	// Chamando o handler para a requisição simulada
	r.ServeHTTP(rr, req)

	// Verificando o código de status HTTP
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	// Verificando o corpo da resposta
	var response utils.DefaultResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, response.Status)

	// Verificando a mensagem de erro
	assert.Contains(t, strings.ToLower(response.Message), "erro ao criar usuário")

	// Verificando se o método CreateUser foi chamado no mock
	mockService.AssertCalled(t, "CreateUser", mock.Anything, newUser)
}
