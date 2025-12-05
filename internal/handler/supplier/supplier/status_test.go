package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockSupplier "github.com/WagaoCarvalho/backend_store_go/infra/mock/supplier"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSupplierHandler_Enable(t *testing.T) {
	mockService := new(mockSupplier.MockSupplier)
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)
	handler := NewSupplierHandler(mockService, logger)

	t.Run("Sucesso ao habilitar fornecedor", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		supplierID := int64(1)
		requestBody := map[string]interface{}{
			"version": 2,
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("GetByID", mock.Anything, supplierID).
			Return(&models.Supplier{
				ID:      supplierID,
				Status:  false,
				Version: 2,
			}, nil).Once()

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
			return s.ID == supplierID && s.Status && s.Version == 2
		})).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/enable", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Método não permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/suppliers/1/enable", nil)
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("ID inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/suppliers/abc/enable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Versão inválida", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/enable", bytes.NewReader([]byte(`{"version":0}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro ao obter fornecedor - não encontrado", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByID", mock.Anything, int64(1)).
			Return(nil, errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/enable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro genérico ao obter fornecedor", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByID", mock.Anything, int64(1)).
			Return(nil, errors.New("erro banco")).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/enable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Conflito de versão ao atualizar fornecedor", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByID", mock.Anything, int64(1)).
			Return(&models.Supplier{
				ID:      1,
				Status:  false,
				Version: 1,
			}, nil).Once()

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
			return s.ID == 1 && s.Version == 1
		})).Return(errMsg.ErrVersionConflict).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/enable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro genérico ao habilitar fornecedor", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByID", mock.Anything, int64(1)).
			Return(&models.Supplier{
				ID:      1,
				Status:  false,
				Version: 1,
			}, nil).Once()

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
			return s.ID == 1 && s.Version == 1
		})).Return(errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/enable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestSupplierHandler_Disable(t *testing.T) {
	mockService := new(mockSupplier.MockSupplier)
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)
	handler := NewSupplierHandler(mockService, logger)

	t.Run("Sucesso ao desabilitar fornecedor", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		supplierID := int64(1)
		requestBody := map[string]interface{}{"version": 2}
		body, _ := json.Marshal(requestBody)

		mockService.On("GetByID", mock.Anything, supplierID).Return(&models.Supplier{
			ID:      supplierID,
			Status:  true,
			Version: 2,
		}, nil).Once()

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
			return s.ID == supplierID && !s.Status && s.Version == 2
		})).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/disable", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Método não permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/suppliers/1/disable", nil)
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("ID inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/suppliers/abc/disable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Versão inválida", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/disable", bytes.NewReader([]byte(`{"version":0}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro ao obter fornecedor - não encontrado", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByID", mock.Anything, int64(1)).Return(nil, errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/disable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro ao obter fornecedor - genérico", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByID", mock.Anything, int64(1)).Return(nil, errors.New("erro banco")).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/disable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Conflito de versão ao atualizar fornecedor", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByID", mock.Anything, int64(1)).Return(&models.Supplier{
			ID:      1,
			Status:  true,
			Version: 1,
		}, nil).Once()

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
			return s.ID == 1 && !s.Status && s.Version == 1
		})).Return(errMsg.ErrVersionConflict).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/disable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro genérico ao atualizar fornecedor", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByID", mock.Anything, int64(1)).Return(&models.Supplier{
			ID:      1,
			Status:  true,
			Version: 1,
		}, nil).Once()

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
			return s.ID == 1 && !s.Status && s.Version == 1
		})).Return(errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/disable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}
