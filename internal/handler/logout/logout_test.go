package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mockLogout "github.com/WagaoCarvalho/backend_store_go/infra/mock/service/logout"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestLogoutHandler_Logout(t *testing.T) {
	// Configuração do logger
	baseLogger := logrus.New()
	baseLogger.Out = &strings.Builder{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success", func(t *testing.T) {
		mockService := new(mockLogout.MockLogoutService)
		handler := NewLogoutHandler(mockService, logAdapter)

		token := "valid_token"
		mockService.On("Logout", mock.Anything, token).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		handler.Logout(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidMethod", func(t *testing.T) {
		mockService := new(mockLogout.MockLogoutService)
		handler := NewLogoutHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodGet, "/logout", nil)
		w := httptest.NewRecorder()

		handler.Logout(w, req)
		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("MissingAuthorizationHeader", func(t *testing.T) {
		mockService := new(mockLogout.MockLogoutService)
		handler := NewLogoutHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		w := httptest.NewRecorder()

		handler.Logout(w, req)
		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("InvalidAuthorizationHeader", func(t *testing.T) {
		mockService := new(mockLogout.MockLogoutService)
		handler := NewLogoutHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		req.Header.Set("Authorization", "InvalidTokenFormat")
		w := httptest.NewRecorder()

		handler.Logout(w, req)
		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(mockLogout.MockLogoutService)
		handler := NewLogoutHandler(mockService, logAdapter)

		token := "token_error"
		mockService.On("Logout", mock.Anything, token).Return(errors.New("service failure"))

		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		handler.Logout(w, req)
		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockService.AssertExpectations(t)
	})
}
