package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	models_user "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	"github.com/WagaoCarvalho/backend_store_go/utils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func muxSetVars(req *http.Request, vars map[string]string) *http.Request {
	return mux.SetURLVars(req, vars)
}

// Mock do serviço de usuário
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) GetAll(ctx context.Context) ([]models_user.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models_user.User), args.Error(1)
}

func (m *MockUserService) GetById(ctx context.Context, uid int64) (models_user.User, error) {
	args := m.Called(ctx, uid)
	return args.Get(0).(models_user.User), args.Error(1)
}

func (m *MockUserService) GetByEmail(ctx context.Context, email string) (models_user.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(models_user.User), args.Error(1)
}

func (m *MockUserService) Delete(ctx context.Context, uid int64) error {
	args := m.Called(ctx, uid)
	return args.Error(0)
}

func (m *MockUserService) Update(ctx context.Context, user *models_user.User) (models_user.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(models_user.User), args.Error(1)
}

func (m *MockUserService) Create(
	ctx context.Context,
	user *models_user.User,
	categoryIDs []int64,
	address *models_address.Address,
	contact *models_contact.Contact,
) (models_user.User, error) {
	args := m.Called(ctx, user, categoryIDs, address, contact)
	return args.Get(0).(models_user.User), args.Error(1)
}

func TestUserHandler_Create(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	t.Run("Sucesso ao criar usuário", func(t *testing.T) {
		expectedUser := models_user.User{
			UID:      1,
			Username: "testuser",
			Email:    "test@example.com",
		}

		address := &models_address.Address{Street: "Rua 1"}
		contact := &models_contact.Contact{Phone: "12345"}

		requestBody := map[string]interface{}{
			"user": map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
			},
			"category_id": []int64{1, 2},
			"address":     address,
			"contact":     contact,
		}

		body, _ := json.Marshal(requestBody)

		// Corrigido aqui com MatchedBy para o *User
		mockService.On("Create",
			mock.Anything,
			mock.MatchedBy(func(u *models_user.User) bool {
				return u.Username == "testuser" && u.Email == "test@example.com"
			}),
			[]int64{1, 2},
			address,
			contact,
		).Return(expectedUser, nil)

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro método não permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("Erro ao decodificar JSON inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte("{invalid json")))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro ao criar usuário no service", func(t *testing.T) {
		userData := &models_user.User{Username: "failuser", Email: "fail@example.com"}
		categoryIDs := []int64{1}
		address := &models_address.Address{Street: "Fail St"}
		contact := &models_contact.Contact{Phone: "99999"}

		requestBody := map[string]interface{}{
			"user":        userData,
			"category_id": categoryIDs,
			"address":     address,
			"contact":     contact,
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("Create", mock.Anything, userData, categoryIDs, address, contact).
			Return(models_user.User{}, errors.New("erro ao criar usuário"))

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestGetUsers_Success(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	users := []models_user.User{{UID: 1, Username: "João"}}
	mockService.On("GetAll", mock.Anything).Return(users, nil)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()

	handler.GetUsers(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp utils.DefaultResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	assert.Equal(t, "Usuários encontrados", resp.Message)
	assert.NotNil(t, resp.Data)
	mockService.AssertExpectations(t)
}

func TestGetUserById_Success(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	expectedUser := models_user.User{UID: 10, Username: "Carlos"}
	mockService.On("GetById", mock.Anything, int64(10)).Return(expectedUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/users/10", nil)
	req = muxSetVars(req, map[string]string{"id": "10"})
	rec := httptest.NewRecorder()

	handler.GetUserById(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)
	var resp utils.DefaultResponse
	json.NewDecoder(rec.Body).Decode(&resp)
	actualUID := int64(resp.Data.(map[string]interface{})["uid"].(float64))
	assert.Equal(t, expectedUser.UID, actualUID)
}

func TestGetUserByEmail_Success(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	expectedUser := models_user.User{UID: 20, Email: "carlos@email.com"}
	mockService.On("GetByEmail", mock.Anything, "carlos@email.com").Return(expectedUser, nil)

	req := httptest.NewRequest(http.MethodGet, "/users/email/carlos@email.com", nil)
	req = muxSetVars(req, map[string]string{"email": "carlos@email.com"})
	rec := httptest.NewRecorder()

	handler.GetUserByEmail(rec, req)

	assert.Equal(t, http.StatusOK, rec.Code)

	var resp utils.DefaultResponse
	err := json.NewDecoder(rec.Body).Decode(&resp)
	assert.NoError(t, err)

	actualUID := int64(resp.Data.(map[string]interface{})["uid"].(float64))
	assert.Equal(t, expectedUser.UID, actualUID)
}

func TestGetUsers_Error(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	mockService.On("GetAll", mock.Anything).Return([]models_user.User(nil), assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	rec := httptest.NewRecorder()

	handler.GetUsers(rec, req)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestGetUserById_InvalidID(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/users/abc", nil)
	req = muxSetVars(req, map[string]string{"id": "abc"})
	rec := httptest.NewRecorder()

	handler.GetUserById(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)

	// Usando a função ParseErrorResponse do utils
	var resp utils.DefaultResponse
	resp, err := utils.ParseErrorResponse(rec.Body.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.Status)
}

func TestGetUserById_NotFound(t *testing.T) {
	mockService := new(MockUserService)
	handler := NewUserHandler(mockService)

	mockService.On("GetById", mock.Anything, int64(99)).Return(models_user.User{}, fmt.Errorf("usuário não encontrado"))

	req := httptest.NewRequest(http.MethodGet, "/users/99", nil)
	req = muxSetVars(req, map[string]string{"id": "99"})
	rec := httptest.NewRecorder()

	handler.GetUserById(rec, req)

	assert.Equal(t, http.StatusNotFound, rec.Code)

	// Usando a função ParseErrorResponse do utils
	var resp utils.DefaultResponse
	resp, err := utils.ParseErrorResponse(rec.Body.Bytes())
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Status)
}
