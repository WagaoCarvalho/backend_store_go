package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockUser "github.com/WagaoCarvalho/backend_store_go/infra/mock/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserHandler_GetVersionByID(t *testing.T) {
	mockService := new(mockUser.MockUser)
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)
	handler := NewUserHandler(mockService, logger)

	t.Run("Sucesso ao obter versão por ID", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetVersionByID", mock.Anything, int64(1)).Return(int64(5), nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/users/1/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.GetVersionByID(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro ID inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/abc/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()

		handler.GetVersionByID(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro usuário não encontrado", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetVersionByID", mock.Anything, int64(999)).Return(int64(0), errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodGet, "/users/999/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		rec := httptest.NewRecorder()

		handler.GetVersionByID(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro genérico ao obter versão", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetVersionByID", mock.Anything, int64(2)).Return(int64(0), errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodGet, "/users/2/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		rec := httptest.NewRecorder()

		handler.GetVersionByID(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}
