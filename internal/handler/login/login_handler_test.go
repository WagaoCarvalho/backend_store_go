package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/login"
	"github.com/WagaoCarvalho/backend_store_go/logger"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLoginService struct {
	mock.Mock
}

func (m *MockLoginService) Login(ctx context.Context, credentials models.LoginCredentials) (string, error) {
	args := m.Called(ctx, credentials)
	return args.String(0), args.Error(1)
}

func newLoginRequest(method, url string, body []byte) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	return mux.SetURLVars(req, map[string]string{})
}

func TestLoginHandler_Login(t *testing.T) {
	logAdapter := logger.NewLoggerAdapter(logrus.New())

	t.Run("Success", func(t *testing.T) {
		mockService := new(MockLoginService)
		handler := NewLoginHandler(mockService, logAdapter)

		creds := models.LoginCredentials{
			Email:    "user@example.com",
			Password: "password123",
		}
		mockService.On("Login", mock.Anything, creds).Return("valid_token", nil)

		body, _ := json.Marshal(creds)
		req := newLoginRequest(http.MethodPost, "/login", body)
		w := httptest.NewRecorder()

		handler.Login(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		type loginResponse struct {
			Status  int                    `json:"status"`
			Message string                 `json:"message"`
			Data    map[string]interface{} `json:"data"`
		}

		var response loginResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Login realizado com sucesso", response.Message)
		assert.Equal(t, "valid_token", response.Data["token"])

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		mockService := new(MockLoginService)
		handler := NewLoginHandler(mockService, logAdapter)

		creds := models.LoginCredentials{
			Email:    "user@example.com",
			Password: "wrongpassword",
		}
		mockService.On("Login", mock.Anything, creds).Return("", errors.New("credenciais inválidas"))

		body, _ := json.Marshal(creds)
		req := newLoginRequest(http.MethodPost, "/login", body)
		w := httptest.NewRecorder()

		handler.Login(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		assert.Equal(t, "credenciais inválidas", response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		mockService := new(MockLoginService)
		handler := NewLoginHandler(mockService, logAdapter)

		invalidJSON := []byte(`{email: "user@example.com", password: }`)
		req := newLoginRequest(http.MethodPost, "/login", invalidJSON)
		w := httptest.NewRecorder()

		handler.Login(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		assert.Equal(t, "dados inválidos", response["message"])
	})

	t.Run("InvalidMethod", func(t *testing.T) {
		mockService := new(MockLoginService)
		handler := NewLoginHandler(mockService, logAdapter)

		req := newLoginRequest(http.MethodGet, "/login", nil)
		w := httptest.NewRecorder()

		handler.Login(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)

		assert.Contains(t, response["message"], "método GET não permitido")
	})
}
