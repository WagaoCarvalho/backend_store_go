package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/login"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockLoginService representa um mock do LoginService
type MockLoginService struct {
	mock.Mock
}

func (m *MockLoginService) Login(ctx context.Context, credentials models.LoginCredentials) (string, error) {
	args := m.Called(ctx, credentials)
	return args.String(0), args.Error(1)
}

func TestLoginHandler_Success(t *testing.T) {
	mockService := new(MockLoginService)
	handler := NewLoginHandler(mockService)

	credentials := models.LoginCredentials{
		Email:    "user@example.com",
		Password: "password123",
	}
	mockService.On("Login", mock.Anything, credentials).Return("valid_token", nil)

	body, _ := json.Marshal(credentials)
	r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var responseData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseData)
	assert.Equal(t, "Login realizado com sucesso", responseData["message"])
	assert.Equal(t, "valid_token", responseData["data"].(map[string]interface{})["token"])

	mockService.AssertExpectations(t)
}

func TestLoginHandler_InvalidCredentials(t *testing.T) {
	mockService := new(MockLoginService)
	handler := NewLoginHandler(mockService)

	credentials := models.LoginCredentials{
		Email:    "user@example.com",
		Password: "wrongpassword",
	}
	mockService.On("Login", mock.Anything, credentials).Return("", errors.New("credenciais inválidas"))

	body, _ := json.Marshal(credentials)
	r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

	var responseData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseData)
	assert.Equal(t, "credenciais inválidas", responseData["message"]) // Corrigido para "message"

	mockService.AssertExpectations(t)
}

func TestLoginHandler_InvalidJSON(t *testing.T) {
	mockService := new(MockLoginService)
	handler := NewLoginHandler(mockService)

	// JSON malformado
	invalidJSON := []byte(`{email: "user@example.com", password: }`)

	r := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(invalidJSON))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Login(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var responseData map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&responseData)

	assert.Equal(t, "dados inválidos", responseData["message"])
}

func TestLoginHandler_InvalidMethod(t *testing.T) {
	mockService := new(MockLoginService)
	handler := NewLoginHandler(mockService)

	r := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	handler.Login(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}
