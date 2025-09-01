package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/login"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/login"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockLoginService struct {
	mock.Mock
}

func (m *MockLoginService) Login(ctx context.Context, credentials models.LoginCredentials) (*models.AuthResponse, error) {
	args := m.Called(ctx, credentials)
	authResp, _ := args.Get(0).(*models.AuthResponse)
	return authResp, args.Error(1)
}

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
		mockService := new(MockLoginService)
		handler := NewLoginHandler(mockService, logAdapter)

		creds := models.LoginCredentials{
			Email:    "user@example.com",
			Password: "password123",
		}

		// retornar AuthResponse correto
		mockService.On("Login", mock.Anything, creds).Return(&models.AuthResponse{
			AccessToken: "valid_token",
			TokenType:   "Bearer",
			ExpiresIn:   3600,
		}, nil)

		body, _ := json.Marshal(dto.LoginCredentialsDTO{
			Email:    creds.Email,
			Password: creds.Password,
		})
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
		assert.Equal(t, float64(3600), response.Data["expires_in"]) // json numbers viram float64
		assert.Equal(t, "Bearer", response.Data["token_type"])

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidCredentials", func(t *testing.T) {
		mockService := new(MockLoginService)
		handler := NewLoginHandler(mockService, logAdapter)

		creds := models.LoginCredentials{
			Email:    "user@example.com",
			Password: "wrongpassword",
		}
		mockService.On("Login", mock.Anything, creds).Return(nil, errors.New("credenciais inválidas"))

		body, _ := json.Marshal(creds)
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
		err := json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

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
		err := json.NewDecoder(resp.Body).Decode(&response)
		require.NoError(t, err)

		assert.Contains(t, response["message"], "método GET não permitido")
	})
}
