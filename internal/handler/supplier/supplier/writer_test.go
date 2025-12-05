package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockSupplier "github.com/WagaoCarvalho/backend_store_go/infra/mock/supplier"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSupplierHandler_Create(t *testing.T) {
	now := time.Now()

	t.Run("Sucesso ao criar fornecedor", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		logger := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, logger)

		expectedSupplier := &models.Supplier{
			ID:        1,
			Name:      "Fornecedor Teste",
			CNPJ:      utils.StrToPtr("12345678000199"),
			Version:   1,
			CreatedAt: now,
			UpdatedAt: now,
		}

		requestBody := map[string]interface{}{
			"supplier": map[string]interface{}{
				"name": "Fornecedor Teste",
				"cnpj": "12345678000199",
			},
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("Create",
			mock.Anything,
			mock.MatchedBy(func(s *models.Supplier) bool {
				return s.Name == "Fornecedor Teste" && s.CNPJ != nil && *s.CNPJ == "12345678000199"
			}),
		).Return(expectedSupplier, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Falha ao criar fornecedor - supplier nil", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		logger := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, logger)

		requestBody := map[string]interface{}{}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.NewDecoder(rec.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.Equal(t, "supplier não fornecido", resp.Message)

		mockService.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("Erro método não permitido", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		logger := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodGet, "/suppliers", nil)
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
		mockService.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("Erro ao decodificar JSON inválido", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		logger := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader([]byte("{invalid json")))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockService.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("Erro ao criar fornecedor no service", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		logger := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, logger)

		supplierData := &models.Supplier{
			Name: "Fornecedor Falso",
			CNPJ: utils.StrToPtr("98765432000188"),
		}

		requestBody := map[string]interface{}{
			"supplier": supplierData,
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("Create",
			mock.Anything,
			mock.MatchedBy(func(s *models.Supplier) bool {
				return s.Name == supplierData.Name && s.CNPJ != nil && *s.CNPJ == *supplierData.CNPJ
			}),
		).Return(nil, errors.New("erro ao criar fornecedor")).Once()

		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestSupplierHandler_Update(t *testing.T) {
	mockService := new(mockSupplier.MockSupplier)
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)
	handler := NewSupplierHandler(mockService, logger)

	t.Run("Sucesso ao atualizar fornecedor", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		supplierID := int64(1)
		requestBody := map[string]interface{}{
			"supplier": map[string]interface{}{
				"name":    "Fornecedor Atualizado",
				"version": 2,
				"cnpj":    "12345678000199",
			},
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
			return s.ID == supplierID && s.Name == "Fornecedor Atualizado" && s.Version == 2
		})).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPut, "/suppliers/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro chave estrangeira inválida", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		supplierID := int64(1)
		requestBody := map[string]interface{}{
			"supplier": map[string]interface{}{
				"name":    "Fornecedor X",
				"version": 1,
				"cnpj":    "12345678000199",
			},
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
			return s.ID == supplierID
		})).Return(errMsg.ErrDBInvalidForeignKey).Once()

		req := httptest.NewRequest(http.MethodPut, "/suppliers/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro fornecedor duplicado", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		supplierID := int64(1)
		requestBody := map[string]interface{}{
			"supplier": map[string]interface{}{
				"name":    "Fornecedor Y",
				"version": 1,
				"cnpj":    "12345678000199",
			},
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
			return s.ID == supplierID
		})).Return(errMsg.ErrDuplicate).Once()

		req := httptest.NewRequest(http.MethodPut, "/suppliers/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro método não permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/suppliers/1", nil)
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("Erro ID inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/suppliers/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro JSON inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/suppliers/1", bytes.NewReader([]byte("{invalid json")))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro dados do fornecedor ausentes", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"supplier": nil,
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPut, "/suppliers/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro validação (ErrInvalidData)", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		body := `{"supplier": {"name": "", "version": 1}}`

		mockService.On("Update", mock.Anything, mock.Anything).Return(errMsg.ErrInvalidData).Once()

		req := httptest.NewRequest(http.MethodPut, "/suppliers/1", bytes.NewReader([]byte(body)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro ID zero (ErrZeroID)", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		body := `{"supplier": {"name": "Fornecedor", "version": 1}}`

		mockService.On("Update", mock.Anything, mock.Anything).Return(errMsg.ErrZeroID).Once()

		req := httptest.NewRequest(http.MethodPut, "/suppliers/1", bytes.NewReader([]byte(body)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro conflito de versão", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		body := `{"supplier": {"name": "Fornecedor", "version": 2}}`

		mockService.On("Update", mock.Anything, mock.Anything).Return(errMsg.ErrVersionConflict).Once()

		req := httptest.NewRequest(http.MethodPut, "/suppliers/1", bytes.NewReader([]byte(body)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro não encontrado", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		body := `{"supplier": {"name": "Fornecedor Inexistente", "version": 2}}`

		mockService.On("Update", mock.Anything, mock.Anything).Return(errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodPut, "/suppliers/1", bytes.NewReader([]byte(body)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro interno genérico", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		body := `{"supplier": {"name": "Fornecedor X", "version": 1}}`

		mockService.On("Update", mock.Anything, mock.Anything).Return(errors.New("erro inesperado")).Once()

		req := httptest.NewRequest(http.MethodPut, "/suppliers/1", bytes.NewReader([]byte(body)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestSupplierHandler_Delete(t *testing.T) {
	mockService := new(mockSupplier.MockSupplier)
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)
	handler := NewSupplierHandler(mockService, logger)

	t.Run("Sucesso ao deletar fornecedor", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("Delete", mock.Anything, int64(1)).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/suppliers/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Método não permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/suppliers/1", nil)
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("ID inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/suppliers/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Fornecedor não encontrado", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("Delete", mock.Anything, int64(1)).Return(fmt.Errorf("fornecedor não encontrado")).Once()

		req := httptest.NewRequest(http.MethodDelete, "/suppliers/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro genérico ao deletar fornecedor", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("Delete", mock.Anything, int64(1)).Return(errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodDelete, "/suppliers/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}
