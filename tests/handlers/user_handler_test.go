package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WagaoCarvalho/backend_store_go/internal/handlers"
	"github.com/WagaoCarvalho/backend_store_go/internal/models"
	"github.com/WagaoCarvalho/backend_store_go/utils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type mockUserService struct{}

func (m *mockUserService) GetUsers(ctx context.Context) ([]models.User, error) {
	// Mockando a resposta de usuários
	return []models.User{
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
	}, nil
}

func TestGetUsers(t *testing.T) {
	// Criando o serviço mock
	mockService := &mockUserService{}
	userHandler := handlers.NewUserHandler(mockService)

	// Criando um roteador e registrando a rota de GetUsers
	r := mux.NewRouter()
	r.HandleFunc("/users", userHandler.GetUsers).Methods("GET")

	// Criando a requisição HTTP simulada
	req := httptest.NewRequest("GET", "/users", nil)
	rr := httptest.NewRecorder()

	// Chamando o handler para a requisição simulada
	r.ServeHTTP(rr, req)

	// Verificando o código de status HTTP
	assert.Equal(t, http.StatusOK, rr.Code)

	// Verificando o corpo da resposta
	var response utils.DefaultResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Status)

	// Verificando se os dados retornados são os esperados
	users := response.Data.([]interface{})
	assert.Equal(t, 2, len(users))

	// Verificando os dados dos usuários retornados
	assert.Equal(t, float64(1), users[0].(map[string]interface{})["uid"])
	assert.Equal(t, "user1", users[0].(map[string]interface{})["nickname"])
	assert.Equal(t, "user1@example.com", users[0].(map[string]interface{})["email"])

	assert.Equal(t, float64(2), users[1].(map[string]interface{})["uid"])
	assert.Equal(t, "user2", users[1].(map[string]interface{})["nickname"])
	assert.Equal(t, "user2@example.com", users[1].(map[string]interface{})["email"])
}
