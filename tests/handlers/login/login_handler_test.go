package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	handlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/login"
	"github.com/WagaoCarvalho/backend_store_go/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLoginHandler_Success(t *testing.T) {
	mockService := new(MockLoginService)
	handler := handlers.NewLoginHandler(mockService)

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
	handler := handlers.NewLoginHandler(mockService)

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

func TestLoginHandler_InvalidMethod(t *testing.T) {
	mockService := new(MockLoginService)
	handler := handlers.NewLoginHandler(mockService)

	r := httptest.NewRequest(http.MethodGet, "/login", nil)
	w := httptest.NewRecorder()

	handler.Login(w, r)

	resp := w.Result()
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
}
