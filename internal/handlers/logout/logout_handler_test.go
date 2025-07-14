package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockLogoutService struct {
	mock.Mock
}

func (m *MockLogoutService) Logout(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func TestLogoutHandler_Logout(t *testing.T) {
	logAdapter := logger.NewLoggerAdapter(logrus.New())

	t.Run("Success", func(t *testing.T) {
		mockSvc := new(MockLogoutService)
		handler := NewLogoutHandler(mockSvc, logAdapter)

		token := "valid_token"
		mockSvc.On("Logout", mock.Anything, token).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		handler.Logout(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidMethod", func(t *testing.T) {
		mockSvc := new(MockLogoutService)
		handler := NewLogoutHandler(mockSvc, logAdapter)

		req := httptest.NewRequest(http.MethodGet, "/logout", nil)
		w := httptest.NewRecorder()

		handler.Logout(w, req)
		assert.Equal(t, http.StatusMethodNotAllowed, w.Result().StatusCode)
	})

	t.Run("MissingAuthorizationHeader", func(t *testing.T) {
		mockSvc := new(MockLogoutService)
		handler := NewLogoutHandler(mockSvc, logAdapter)

		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		w := httptest.NewRecorder()

		handler.Logout(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})

	t.Run("InvalidAuthorizationFormat", func(t *testing.T) {
		mockSvc := new(MockLogoutService)
		handler := NewLogoutHandler(mockSvc, logAdapter)

		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		req.Header.Set("Authorization", "InvalidFormat")
		w := httptest.NewRecorder()

		handler.Logout(w, req)
		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockSvc := new(MockLogoutService)
		handler := NewLogoutHandler(mockSvc, logAdapter)

		token := "some_token"
		errorExpected := assert.AnError
		mockSvc.On("Logout", mock.Anything, token).Return(errorExpected)

		req := httptest.NewRequest(http.MethodPost, "/logout", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		handler.Logout(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})
}
