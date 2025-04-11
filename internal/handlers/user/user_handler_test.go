package handlers

import (
	"context"
	"encoding/json"
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

func (m *MockUserService) GetById(ctx context.Context, id int64) (models_user.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models_user.User), args.Error(1)
}

func (m *MockUserService) GetByEmail(ctx context.Context, email string) (models_user.User, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(models_user.User), args.Error(1)
}

func (m *MockUserService) Create(ctx context.Context, user models_user.User, categoryID int64, address models_address.Address, contact models_contact.Contact) (models_user.User, error) {
	args := m.Called(ctx, user, categoryID, address, contact)
	return args.Get(0).(models_user.User), args.Error(1)
}

func (m *MockUserService) Update(ctx context.Context, user models_user.User, contact *models_contact.Contact) (models_user.User, error) {
	args := m.Called(ctx, user, contact)
	return args.Get(0).(models_user.User), args.Error(1)
}

func (m *MockUserService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
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
