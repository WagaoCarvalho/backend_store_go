package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockLogin "github.com/WagaoCarvalho/backend_store_go/infra/mock/login"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/login"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func newLoginRequest(method, url string, body []byte) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	return mux.SetURLVars(req, map[string]string{})
}

func TestLoginHandler_Login(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success", func(t *testing.T) {
		mockService := new(mockLogin.MockLoginService)
		handler := NewLoginHandler(mockService, logAdapter)

		email := "user@example.com"
		password := "password123"

		mockService.On("Login", mock.Anything, email, password).Return(&models.AuthResponse{
			AccessToken: "valid_token",
			TokenType:   "Bearer",
			ExpiresIn:   3600,
		}, nil)

		body, _ := json.Marshal(map[string]string{"email": email, "password": password})
		req := newLoginRequest(http.MethodPost, "/login", body)
		w := httptest.NewRecorder()

		handler.Login(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Status  int                    `json:"status"`
			Message string                 `json:"message"`
			Data    map[string]interface{} `json:"data"`
		}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Login realizado com sucesso", response.Message)
		assert.Equal(t, "valid_token", response.Data["access_token"])
		assert.Equal(t, "Bearer", response.Data["token_type"])
		assert.Equal(t, float64(3600), response.Data["expires_in"])

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidMethod", func(t *testing.T) {
		mockService := new(mockLogin.MockLoginService)
		handler := NewLoginHandler(mockService, logAdapter)

		req := newLoginRequest(http.MethodGet, "/login", nil)
		w := httptest.NewRecorder()

		handler.Login(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Contains(t, response["message"], "método GET não permitido")
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		mockService := new(mockLogin.MockLoginService)
		handler := NewLoginHandler(mockService, logAdapter)

		invalidJSON := []byte(`{email: "user@example.com", password: }`)
		req := newLoginRequest(http.MethodPost, "/login", invalidJSON)
		w := httptest.NewRecorder()

		handler.Login(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, "dados inválidos", response["message"])
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		mockService := new(mockLogin.MockLoginService)
		handler := NewLoginHandler(mockService, logAdapter)

		email := "user@example.com"
		password := "wrongpassword"

		mockService.On("Login", mock.Anything, email, password).Return(nil, errors.New("credenciais inválidas"))

		body, _ := json.Marshal(map[string]string{"email": email, "password": password})
		req := newLoginRequest(http.MethodPost, "/login", body)
		w := httptest.NewRecorder()

		handler.Login(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, "credenciais inválidas", response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(mockLogin.MockLoginService)
		handler := NewLoginHandler(mockService, logAdapter)

		email := "user@example.com"
		password := "password123"

		mockService.On("Login", mock.Anything, email, password).Return(nil, errors.New("erro inesperado"))

		body, _ := json.Marshal(map[string]string{"email": email, "password": password})
		req := newLoginRequest(http.MethodPost, "/login", body)
		w := httptest.NewRecorder()

		handler.Login(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)
		assert.Equal(t, "erro inesperado", response["message"])

		mockService.AssertExpectations(t)
	})
}
