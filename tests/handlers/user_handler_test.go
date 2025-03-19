package handlers

import (
	"context"
	"encoding/json"
	"fmt"
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

func (m *mockUserService) GetUserById(ctx context.Context, uid int64) (models.User, error) {
	if uid == 1 {
		return models.User{
			UID:      1,
			Username: "user1",
			Email:    "user1@example.com",
			Status:   true,
		}, nil
	}
	return models.User{}, fmt.Errorf("usuário não encontrado")
}

func TestGetUsers(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := handlers.NewUserHandler(mockService)

	r := mux.NewRouter()
	r.HandleFunc("/users", userHandler.GetUsers).Methods("GET")

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

func (m *mockUserService) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	if email == "user1@example.com" {
		return models.User{
			UID:      1,
			Username: "user1",
			Email:    "user1@example.com",
			Status:   true,
		}, nil
	}
	return models.User{}, fmt.Errorf("usuário não encontrado")
}

func TestGetUserByEmail(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := handlers.NewUserHandler(mockService)

	r := mux.NewRouter()
	r.HandleFunc("/user/{email}", userHandler.GetUserByEmail).Methods("GET")

	req := httptest.NewRequest("GET", "/user/user1@example.com", nil)
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

func TestGetUserByEmail_UserNotFound(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := handlers.NewUserHandler(mockService)

	r := mux.NewRouter()
	r.HandleFunc("/user/{email}", userHandler.GetUserByEmail).Methods("GET")

	req := httptest.NewRequest("GET", "/user/notfound@example.com", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)

	var response utils.DefaultResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, response.Status)

	assert.Contains(t, response.Message, "usuário não encontrado")
}
