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
	t.Run("successfully get all suppliers", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		expectedSuppliers := []*models.Supplier{
			{
				ID:        1,
				Name:      "Fornecedor A",
				CNPJ:      utils.StrToPtr("12345678000199"),
				Version:   1,
				Status:    true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        2,
				Name:      "Fornecedor B",
				CPF:       utils.StrToPtr("12345678901"),
				Version:   1,
				Status:    true,
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

	t.Run("return error when service fails", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		mockService.On("GetAll", mock.Anything).Return([]*models.Supplier(nil), errors.New("database error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/suppliers", nil)
		rec := httptest.NewRecorder()

		handler.GetAll(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestSupplierHandler_GetByID(t *testing.T) {
	t.Run("successfully get supplier by id", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		expectedSupplier := &models.Supplier{
			ID:      1,
			Name:    "Fornecedor A",
			CNPJ:    utils.StrToPtr("12345678000199"),
			Status:  true,
			Version: 1,
		}

		mockService.On("GetByID", mock.Anything, int64(1)).Return(expectedSupplier, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/suppliers/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.GetByID(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Fornecedor encontrado", response.Message)

		dataBytes, _ := json.Marshal(response.Data)
		var result models.Supplier
		err = json.Unmarshal(dataBytes, &result)
		require.NoError(t, err)
		assert.Equal(t, expectedSupplier.ID, result.ID)
		assert.Equal(t, expectedSupplier.Name, result.Name)

		mockService.AssertExpectations(t)
	})

	t.Run("return bad request when id is invalid", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		req := httptest.NewRequest(http.MethodGet, "/suppliers/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()

		handler.GetByID(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response.Message, "invalid ID")
	})

	t.Run("return not found when supplier does not exist", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		mockService.On("GetByID", mock.Anything, int64(999)).Return(nil, errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodGet, "/suppliers/999", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		rec := httptest.NewRecorder()

		handler.GetByID(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, errMsg.ErrNotFound.Error(), response.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("return internal server error when service fails", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		mockService.On("GetByID", mock.Anything, int64(1)).Return(nil, errors.New("unexpected error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/suppliers/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.GetByID(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "unexpected error", response.Message)

		mockService.AssertExpectations(t)
	})
}

func TestSupplierHandler_GetByName(t *testing.T) {
	t.Run("successfully get suppliers by name", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		expectedSuppliers := []*models.Supplier{
			{
				ID:   1,
				Name: "Fornecedor A",
				CNPJ: utils.StrToPtr("12345678000199"),
			},
			{
				ID:   2,
				Name: "Fornecedor AB",
				CPF:  utils.StrToPtr("12345678901"),
			},
		}

		mockService.On("GetByName", mock.Anything, "fornecedor").Return(expectedSuppliers, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/suppliers/name/fornecedor", nil)
		req = mux.SetURLVars(req, map[string]string{"name": "fornecedor"})
		rec := httptest.NewRecorder()

		handler.GetByName(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Fornecedores encontrados", response.Message)
		assert.Len(t, response.Data.([]any), 2)

		mockService.AssertExpectations(t)
	})

	t.Run("return bad request when name parameter is missing", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		req := httptest.NewRequest(http.MethodGet, "/suppliers/name/", nil)
		req = mux.SetURLVars(req, map[string]string{})
		rec := httptest.NewRecorder()

		handler.GetByName(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		// Corrigido: mensagem real retornada pelo handler
		assert.Equal(t, "missing or empty param: name", response.Message)

		mockService.AssertNotCalled(t, "GetByName")
	})

	t.Run("return bad request when name parameter is empty", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		req := httptest.NewRequest(http.MethodGet, "/suppliers/name/", nil)
		req = mux.SetURLVars(req, map[string]string{"name": ""})
		rec := httptest.NewRecorder()

		handler.GetByName(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		// Corrigido: mensagem real retornada pelo handler
		assert.Equal(t, "missing or empty param: name", response.Message)

		mockService.AssertNotCalled(t, "GetByName")
	})

	t.Run("return not found when no suppliers found", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		mockService.On("GetByName", mock.Anything, "notfound").Return(nil, errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodGet, "/suppliers/name/notfound", nil)
		req = mux.SetURLVars(req, map[string]string{"name": "notfound"})
		rec := httptest.NewRecorder()

		handler.GetByName(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, errMsg.ErrNotFound.Error(), response.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("return internal server error when service fails", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		mockService.On("GetByName", mock.Anything, "error").Return(nil, errors.New("database error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/suppliers/name/error", nil)
		req = mux.SetURLVars(req, map[string]string{"name": "error"})
		rec := httptest.NewRecorder()

		handler.GetByName(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "database error", response.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("successfully get suppliers by name with special characters", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		expectedSuppliers := []*models.Supplier{
			{
				ID:   1,
				Name: "Fornecedor Especial",
				CNPJ: utils.StrToPtr("12345678000199"),
			},
		}

		mockService.On("GetByName", mock.Anything, "fornecedor@#$").Return(expectedSuppliers, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/suppliers/name/fornecedor@#$", nil)
		req = mux.SetURLVars(req, map[string]string{"name": "fornecedor@#$"})
		rec := httptest.NewRecorder()

		handler.GetByName(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Fornecedores encontrados", response.Message)
		assert.Len(t, response.Data.([]any), 1)

		mockService.AssertExpectations(t)
	})

	t.Run("successfully get suppliers by name with case insensitive search", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		expectedSuppliers := []*models.Supplier{
			{
				ID:   1,
				Name: "FORNECEDOR MAIUSCULO",
				CNPJ: utils.StrToPtr("12345678000199"),
			},
		}

		// O serviço deve receber o nome como foi enviado (pode ser tratado no service)
		mockService.On("GetByName", mock.Anything, "FORNECEDOR").Return(expectedSuppliers, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/suppliers/name/FORNECEDOR", nil)
		req = mux.SetURLVars(req, map[string]string{"name": "FORNECEDOR"})
		rec := httptest.NewRecorder()

		handler.GetByName(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Fornecedores encontrados", response.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("return empty list when service returns empty slice (not ErrNotFound)", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		mockService.On("GetByName", mock.Anything, "empty").Return([]*models.Supplier{}, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/suppliers/name/empty", nil)
		req = mux.SetURLVars(req, map[string]string{"name": "empty"})
		rec := httptest.NewRecorder()

		handler.GetByName(rec, req)

		// Deve retornar OK com lista vazia
		assert.Equal(t, http.StatusOK, rec.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Fornecedores encontrados", response.Message)
		assert.Len(t, response.Data.([]any), 0)

		mockService.AssertExpectations(t)
	})

}

func TestSupplierHandler_GetVersionByID(t *testing.T) {
	t.Run("successfully get version by id", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		mockService.On("GetVersionByID", mock.Anything, int64(1)).Return(int64(5), nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/suppliers/1/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.GetVersionByID(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Versão do fornecedor obtida com sucesso", response.Message)

		dataMap, ok := response.Data.(map[string]interface{})
		assert.True(t, ok)
		assert.EqualValues(t, 5, dataMap["version"])

		mockService.AssertExpectations(t)
	})

	t.Run("return bad request when id is invalid", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		req := httptest.NewRequest(http.MethodGet, "/suppliers/abc/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()

		handler.GetVersionByID(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Contains(t, response.Message, "invalid ID")
	})

	t.Run("return not found when supplier does not exist", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		mockService.On("GetVersionByID", mock.Anything, int64(999)).Return(int64(0), errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodGet, "/suppliers/999/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		rec := httptest.NewRecorder()

		handler.GetVersionByID(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, errMsg.ErrNotFound.Error(), response.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("return internal server error when service fails", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		log := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierHandler(mockService, log)

		mockService.On("GetVersionByID", mock.Anything, int64(1)).Return(int64(0), errors.New("database error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/suppliers/1/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.GetVersionByID(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "database error", response.Message)

		mockService.AssertExpectations(t)
	})
}
