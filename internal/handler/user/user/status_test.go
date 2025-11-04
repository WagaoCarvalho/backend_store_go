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

func TestUserHandler_Disable(t *testing.T) {
	mockService := new(mockUser.MockUser)
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)
	handler := NewUser(mockService, logger)

	t.Run("Sucesso ao desabilitar usuário", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("Disable", mock.Anything, int64(1)).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPatch, "/users/1/disable", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro ID zero", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		// Configura o mock para retornar ErrZeroID
		mockService.On("Disable", mock.Anything, int64(0)).Return(errMsg.ErrZeroID).Once()

		req := httptest.NewRequest(http.MethodPatch, "/users/0/disable", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "0"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro método não permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/1/disable", nil)
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("Erro ID inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/users/abc/disable", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro usuário não encontrado", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("Disable", mock.Anything, int64(999)).Return(errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodPatch, "/users/999/disable", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro genérico ao desabilitar usuário", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("Disable", mock.Anything, int64(2)).Return(errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodPatch, "/users/2/disable", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_Enable(t *testing.T) {
	mockService := new(mockUser.MockUser)
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)
	handler := NewUser(mockService, logger)

	t.Run("Sucesso ao habilitar usuário", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("Enable", mock.Anything, int64(1)).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPatch, "/users/1/enable", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro ID zero", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		// Configura o mock para retornar ErrZeroID
		mockService.On("Enable", mock.Anything, int64(0)).Return(errMsg.ErrZeroID).Once()

		req := httptest.NewRequest(http.MethodPatch, "/users/0/enable", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "0"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro método não permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/1/enable", nil)
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("Erro ID inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/users/abc/enable", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro usuário não encontrado", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("Enable", mock.Anything, int64(999)).Return(errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodPatch, "/users/999/enable", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro genérico ao habilitar usuário", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("Enable", mock.Anything, int64(2)).Return(errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodPatch, "/users/2/enable", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}
