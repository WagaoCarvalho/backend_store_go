package handler

import (
	"bytes"
	"encoding/json"
	"errors"
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
	"github.com/stretchr/testify/require"
)

func TestSupplierHandler_GetAll(t *testing.T) {
	t.Run("Sucesso ao obter todos os fornecedores", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		logger := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, logger)

		expectedSuppliers := []*models.Supplier{
			{
				ID:        1,
				Name:      "Fornecedor A",
				CNPJ:      utils.StrToPtr("12345678000199"),
				Version:   1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        2,
				Name:      "Fornecedor B",
				CPF:       utils.StrToPtr("12345678901"),
				Version:   1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		}

		mockService.On("GetAll", mock.Anything).Return(expectedSuppliers, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/suppliers", nil)
		rec := httptest.NewRecorder()

		handler.GetAll(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)

		var response utils.DefaultResponse
		err := utils.FromJSON(rec.Body, &response)
		assert.NoError(t, err)
		assert.Equal(t, "Fornecedores encontrados", response.Message)
		assert.Len(t, response.Data.([]any), 2)
	})

	t.Run("Erro ao obter fornecedores no service", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		logger := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, logger)

		mockService.On("GetAll", mock.Anything).Return([]*models.Supplier(nil), errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodGet, "/suppliers", nil)
		rec := httptest.NewRecorder()

		handler.GetAll(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)

		var response utils.DefaultResponse
		err := utils.FromJSON(rec.Body, &response)
		assert.NoError(t, err)
		assert.Contains(t, response.Message, "erro ao buscar fornecedores")
	})
}

func TestSupplierHandler_GetByID(t *testing.T) {
	mockSvc := new(mockSupplier.MockSupplier)
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	log := logger.NewLoggerAdapter(baseLogger)
	handler := NewSupplierHandler(mockSvc, log)

	t.Run("Sucesso ao obter fornecedor por ID", func(t *testing.T) {
		expected := &models.Supplier{
			ID:   1,
			Name: "Fornecedor A",
		}

		mockSvc.On("GetByID", mock.Anything, int64(1)).Return(expected, nil)

		req := httptest.NewRequest("GET", "/suppliers/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		itemBytes, _ := json.Marshal(resp.Data)
		var result models.Supplier
		err = json.Unmarshal(itemBytes, &result)
		require.NoError(t, err)

		assert.Equal(t, expected.ID, result.ID)
		assert.Equal(t, expected.Name, result.Name)
		assert.Equal(t, "Fornecedor encontrado", resp.Message)

		mockSvc.AssertExpectations(t)
	})

	t.Run("ID inválido", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/suppliers/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "ID inválido", resp["message"])
	})

	t.Run("Fornecedor não encontrado", func(t *testing.T) {
		mockSvc := new(mockSupplier.MockSupplier)
		mockSvc.On("GetByID", mock.Anything, int64(99)).Return(nil, errors.New("fornecedor não encontrado"))

		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockSvc, log)

		req := httptest.NewRequest("GET", "/suppliers/99", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "99"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var resp map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "fornecedor não encontrado", resp["message"])

		mockSvc.AssertExpectations(t)
	})

	t.Run("Erro interno ao buscar fornecedor", func(t *testing.T) {
		mockSvc := new(mockSupplier.MockSupplier)
		mockSvc.On("GetByID", mock.Anything, int64(2)).Return(nil, errors.New("erro inesperado"))

		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockSvc, log)

		req := httptest.NewRequest("GET", "/suppliers/2", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp map[string]any
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "erro inesperado", resp["message"])

		mockSvc.AssertExpectations(t)
	})
}

func TestSupplierHandler_GetByName(t *testing.T) {
	mockService := new(mockSupplier.MockSupplier)
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)
	handler := NewSupplierHandler(mockService, logger)

	t.Run("Sucesso ao buscar fornecedores por nome parcial", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		suppliers := []*models.Supplier{
			{
				ID:   1,
				Name: "Fornecedor A",
				CNPJ: utils.StrToPtr("111"),
			},
			{
				ID:   2,
				Name: "Fornecedor AB",
				CNPJ: utils.StrToPtr("222"),
			},
		}

		mockService.On("GetByName", mock.Anything, "fornecedor").Return(suppliers, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/suppliers/name/fornecedor", nil)
		req = mux.SetURLVars(req, map[string]string{"name": "fornecedor"})
		rec := httptest.NewRecorder()

		handler.GetByName(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro fornecedor não encontrado", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByName", mock.Anything, "notfound").Return(nil, errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodGet, "/suppliers/name/notfound", nil)
		req = mux.SetURLVars(req, map[string]string{"name": "notfound"})
		rec := httptest.NewRecorder()

		handler.GetByName(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro genérico ao buscar fornecedor por nome", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByName", mock.Anything, "error").Return(nil, errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodGet, "/suppliers/name/error", nil)
		req = mux.SetURLVars(req, map[string]string{"name": "error"})
		rec := httptest.NewRecorder()

		handler.GetByName(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestSupplierHandler_GetVersionByID(t *testing.T) {
	mockService := new(mockSupplier.MockSupplier)
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)
	handler := NewSupplierHandler(mockService, logger)

	t.Run("Sucesso ao obter versão do fornecedor", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetVersionByID", mock.Anything, int64(1)).Return(int64(3), nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/suppliers/1/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.GetVersionByID(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		dataMap, ok := resp.Data.(map[string]interface{})
		assert.True(t, ok)
		assert.EqualValues(t, 3, dataMap["version"])
		assert.Equal(t, "Versão do fornecedor obtida com sucesso", resp.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("ID inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/suppliers/abc/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()

		handler.GetVersionByID(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "ID inválido", resp.Message)
	})

	t.Run("Fornecedor não encontrado", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetVersionByID", mock.Anything, int64(2)).
			Return(int64(0), errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodGet, "/suppliers/2/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		rec := httptest.NewRecorder()

		handler.GetVersionByID(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, errMsg.ErrNotFound.Error(), resp.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("Erro interno ao buscar versão", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetVersionByID", mock.Anything, int64(3)).
			Return(int64(0), errors.New("erro inesperado")).Once()

		req := httptest.NewRequest(http.MethodGet, "/suppliers/3/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "3"})
		rec := httptest.NewRecorder()

		handler.GetVersionByID(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "erro inesperado", resp.Message)

		mockService.AssertExpectations(t)
	})
}
