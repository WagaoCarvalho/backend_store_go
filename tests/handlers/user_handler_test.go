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

func (m *mockUserService) UpdateUser(ctx context.Context, user models.User) (models.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(models.User), args.Error(1)
}

func TestGetUsers(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := handlers.NewUserHandler(mockService)

	r := mux.NewRouter()
	r.HandleFunc("/users", userHandler.GetUsers).Methods("GET")

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
	assert.Equal(t, "user1", users[0].(map[string]interface{})["username"])
	assert.Equal(t, "user1@example.com", users[0].(map[string]interface{})["email"])

	assert.Equal(t, float64(2), users[1].(map[string]interface{})["uid"])
	assert.Equal(t, "user2", users[1].(map[string]interface{})["username"])
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
	assert.Equal(t, "user1", user["username"])
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

	mockService.On("CreateUser", mock.Anything, mock.MatchedBy(func(user models.User) bool {

		return user.Username == "user1" && user.Email == "user1@example.com" && user.Password == "hashedpassword" && user.Status == true
	})).Return(expectedUser, nil)

	r := mux.NewRouter()
	r.HandleFunc("/user", userHandler.CreateUser).Methods("POST")

	userData := `{
		"username": "user1",
		"email": "user1@example.com",
		"password": "hashedpassword",
		"status": true
	}`
	req := httptest.NewRequest("POST", "/user", strings.NewReader(userData))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)

	var response utils.DefaultResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, response.Status)

	user := response.Data.(map[string]interface{})
	assert.Equal(t, float64(1), user["uid"])
	assert.Equal(t, "user1", user["username"])
	assert.Equal(t, "user1@example.com", user["email"])

	mockService.AssertCalled(t, "CreateUser", mock.Anything, mock.MatchedBy(func(user models.User) bool {
		return user.Username == "user1" && user.Email == "user1@example.com" && user.Password == "hashedpassword" && user.Status == true
	}))
}

func TestCreateUser_Failure(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := handlers.NewUserHandler(mockService)

	newUser := models.User{
		Username: "user2",
		Email:    "user2@example.com",
		Password: "hashedpassword",
		Status:   true,
	}

	mockService.On("CreateUser", mock.Anything, newUser).Return(models.User{}, fmt.Errorf("erro ao criar usuário"))

	r := mux.NewRouter()
	r.HandleFunc("/user", userHandler.CreateUser).Methods("POST")

	userData := `{
		"username": "user2",
		"email": "user2@example.com",
		"password": "hashedpassword",
		"status": true
	}`
	req := httptest.NewRequest("POST", "/user", strings.NewReader(userData))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var response utils.DefaultResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, response.Status)

	assert.Contains(t, strings.ToLower(response.Message), "erro ao criar usuário")

	mockService.AssertCalled(t, "CreateUser", mock.Anything, newUser)
}

func TestUpdateUser(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := handlers.NewUserHandler(mockService)

	now := time.Now()
	updatedUser := models.User{
		UID:       1,
		Username:  "updatedUser",
		Email:     "updated@example.com",
		Password:  "newhashedpassword",
		Status:    false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	mockService.On("UpdateUser", mock.Anything, mock.MatchedBy(func(user models.User) bool {
		return user.Username == "updatedUser" && user.Email == "updated@example.com" && user.Password == "newhashedpassword" && user.Status == false
	})).Return(updatedUser, nil)

	r := mux.NewRouter()
	r.HandleFunc("/user/{id}", userHandler.UpdateUser).Methods("PUT")

	userData := `{
		"username": "updatedUser",
		"email": "updated@example.com",
		"password": "newhashedpassword",
		"status": false
	}`
	req := httptest.NewRequest("PUT", "/user/1", strings.NewReader(userData))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var response utils.DefaultResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, response.Status)

	user := response.Data.(map[string]interface{})
	assert.Equal(t, float64(1), user["uid"])
	assert.Equal(t, "updatedUser", user["username"])
	assert.Equal(t, "updated@example.com", user["email"])

	mockService.AssertCalled(t, "UpdateUser", mock.Anything, mock.MatchedBy(func(user models.User) bool {
		return user.Username == "updatedUser" && user.Email == "updated@example.com" && user.Password == "newhashedpassword" && user.Status == false
	}))
}

func TestUpdateUser_Failure(t *testing.T) {
	mockService := &mockUserService{}
	userHandler := handlers.NewUserHandler(mockService)

	mockService.On("UpdateUser", mock.Anything, mock.MatchedBy(func(user models.User) bool {
		return user.Username == "errorUser" && user.Email == "error@example.com" && user.Password == "hashedpassword" && user.Status == true
	})).Return(models.User{}, fmt.Errorf("erro ao atualizar usuário"))

	r := mux.NewRouter()
	r.HandleFunc("/user/{id}", userHandler.UpdateUser).Methods("PUT")

	userData := `{
		"username": "errorUser",
		"email": "error@example.com",
		"password": "hashedpassword",
		"status": true
	}`
	req := httptest.NewRequest("PUT", "/user/1", strings.NewReader(userData))
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	var response utils.DefaultResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, response.Status)

	assert.Contains(t, strings.ToLower(response.Message), "erro ao atualizar usuário")

	mockService.AssertCalled(t, "UpdateUser", mock.Anything, mock.MatchedBy(func(user models.User) bool {
		return user.Username == "errorUser" && user.Email == "error@example.com" && user.Password == "hashedpassword" && user.Status == true
	}))
}
