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

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier"
	supplier "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/supplier/supplier"
	service_mock "github.com/WagaoCarvalho/backend_store_go/internal/service/supplier/supplier_services_mock"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
	"github.com/WagaoCarvalho/backend_store_go/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func strPtr(s string) *string {
	return &s
}

func TestSupplierHandler_Create(t *testing.T) {
	mockService := new(service_mock.MockSupplierService)
	logger := logger.NewLoggerAdapter(logrus.New())
	handler := NewSupplierHandler(mockService, logger)

	t.Run("Sucesso ao criar fornecedor", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		expectedSupplier := &supplier.Supplier{
			ID:        1,
			Name:      "Fornecedor Teste",
			CNPJ:      strPtr("12345678000199"),
			CPF:       nil,
			Version:   1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
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
			mock.MatchedBy(func(s *supplier.Supplier) bool {
				return s.Name == "Fornecedor Teste" && s.CNPJ != nil && *s.CNPJ == "12345678000199"
			}),
		).Return(expectedSupplier, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro método não permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/suppliers", nil)
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("Erro ao decodificar JSON inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewReader([]byte("{invalid json")))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro ao criar fornecedor no service", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		supplierData := &supplier.Supplier{
			Name: "Fornecedor Falso",
			CNPJ: strPtr("98765432000188"),
		}

		requestBody := map[string]interface{}{
			"supplier": supplierData,
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("Create",
			mock.Anything,
			mock.MatchedBy(func(s *supplier.Supplier) bool {
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

func TestSupplierHandler_GetAll(t *testing.T) {
	t.Run("Sucesso ao obter todos os fornecedores", func(t *testing.T) {
		mockService := new(service_mock.MockSupplierService)
		logger := logger.NewLoggerAdapter(logrus.New())
		handler := NewSupplierHandler(mockService, logger)

		expectedSuppliers := []*supplier.Supplier{
			{
				ID:        1,
				Name:      "Fornecedor A",
				CNPJ:      strPtr("12345678000199"),
				Version:   1,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        2,
				Name:      "Fornecedor B",
				CPF:       strPtr("12345678901"),
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
		err := utils.FromJson(rec.Body, &response)
		assert.NoError(t, err)
		assert.Equal(t, "Fornecedores encontrados", response.Message)
		assert.Len(t, response.Data.([]any), 2)
	})

	t.Run("Erro ao obter fornecedores no service", func(t *testing.T) {
		mockService := new(service_mock.MockSupplierService)
		logger := logger.NewLoggerAdapter(logrus.New())
		handler := NewSupplierHandler(mockService, logger)

		mockService.On("GetAll", mock.Anything).Return(nil, errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodGet, "/suppliers", nil)
		rec := httptest.NewRecorder()

		handler.GetAll(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)

		var response utils.DefaultResponse
		err := utils.FromJson(rec.Body, &response)
		assert.NoError(t, err)
		assert.Contains(t, response.Message, "erro ao buscar fornecedores")
	})
}

func TestSupplierHandler_GetByID(t *testing.T) {
	mockSvc := new(service_mock.MockSupplierService)
	log := logger.NewLoggerAdapter(logrus.New())
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
		json.Unmarshal(itemBytes, &result)

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
		mockSvc := new(service_mock.MockSupplierService)
		mockSvc.On("GetByID", mock.Anything, int64(99)).Return(nil, errors.New("fornecedor não encontrado"))

		log := logger.NewLoggerAdapter(logrus.New())
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
		mockSvc := new(service_mock.MockSupplierService)
		mockSvc.On("GetByID", mock.Anything, int64(2)).Return(nil, errors.New("erro inesperado"))

		log := logger.NewLoggerAdapter(logrus.New())
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
	mockService := new(service_mock.MockSupplierService)
	logger := logger.NewLoggerAdapter(logrus.New())
	handler := NewSupplierHandler(mockService, logger)

	t.Run("Sucesso ao buscar fornecedores por nome parcial", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		suppliers := []*models.Supplier{
			{
				ID:   1,
				Name: "Fornecedor A",
				CNPJ: strPtr("111"),
			},
			{
				ID:   2,
				Name: "Fornecedor AB",
				CNPJ: strPtr("222"),
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

		mockService.On("GetByName", mock.Anything, "notfound").Return(nil, repo.ErrSupplierNotFound).Once()

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
	mockService := new(service_mock.MockSupplierService)
	logger := logger.NewLoggerAdapter(logrus.New())
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
			Return(int64(0), repo.ErrSupplierNotFound).Once()

		req := httptest.NewRequest(http.MethodGet, "/suppliers/2/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		rec := httptest.NewRecorder()

		handler.GetVersionByID(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, repo.ErrSupplierNotFound.Error(), resp.Message)

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

func TestSupplierHandler_Update(t *testing.T) {
	mockService := new(service_mock.MockSupplierService)
	logger := logger.NewLoggerAdapter(logrus.New())
	handler := NewSupplierHandler(mockService, logger)

	t.Run("Sucesso ao atualizar fornecedor", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		supplierID := int64(1)
		supplierToUpdate := &models.Supplier{
			ID:      supplierID,
			Name:    "Fornecedor Atualizado",
			Version: 2,
			CNPJ:    strPtr("12345678000199"),
		}

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
		})).Return(supplierToUpdate, nil).Once()

		req := httptest.NewRequest(http.MethodPut, "/suppliers/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
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

	t.Run("Erro conflito de versão", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		supplierID := int64(1)

		requestBody := map[string]interface{}{
			"supplier": map[string]interface{}{
				"version": 2,
			},
		}

		body, _ := json.Marshal(requestBody)

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
			return s.ID == supplierID && s.Version == 2
		})).Return(nil, repo.ErrVersionConflict).Once()

		req := httptest.NewRequest(http.MethodPut, "/suppliers/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro genérico ao atualizar fornecedor", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		supplierID := int64(1)

		requestBody := map[string]interface{}{
			"supplier": map[string]interface{}{
				"version": 2,
			},
		}

		body, _ := json.Marshal(requestBody)

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
			return s.ID == supplierID && s.Version == 2
		})).Return(nil, errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodPut, "/suppliers/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestSupplierHandler_Enable(t *testing.T) {
	mockService := new(service_mock.MockSupplierService)
	logger := logger.NewLoggerAdapter(logrus.New())
	handler := NewSupplierHandler(mockService, logger)

	t.Run("Sucesso ao habilitar fornecedor", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		supplierID := int64(1)
		supplierToUpdate := &models.Supplier{
			ID:      supplierID,
			Status:  true,
			Version: 2,
		}

		requestBody := map[string]interface{}{
			"version": 2,
		}

		body, _ := json.Marshal(requestBody)

		mockService.On("GetByID", mock.Anything, supplierID).Return(&models.Supplier{
			ID:      supplierID,
			Status:  false,
			Version: 2,
		}, nil).Once()

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
			return s.ID == supplierID && s.Status == true && s.Version == 2
		})).Return(supplierToUpdate, nil).Once()

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

	t.Run("Erro ao obter fornecedor", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByID", mock.Anything, int64(1)).Return(nil, fmt.Errorf("fornecedor não encontrado")).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/enable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Conflito de versão ao atualizar fornecedor", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByID", mock.Anything, int64(1)).Return(&models.Supplier{
			ID:      1,
			Status:  false,
			Version: 1,
		}, nil).Once()
		mockService.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
			return s.ID == 1 && s.Version == 1
		})).Return(nil, repo.ErrVersionConflict).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/enable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro genérico ao obter fornecedor", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByID", mock.Anything, int64(1)).Return(nil, errors.New("erro banco")).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/enable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro genérico ao habilitar fornecedor", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		// Simula retorno válido do GetByID para seguir no fluxo
		mockService.On("GetByID", mock.Anything, int64(1)).Return(&models.Supplier{
			ID:      1,
			Status:  false,
			Version: 1,
		}, nil).Once()

		// Simula erro genérico no Update, que deve disparar o h.logger.Error e HTTP 500
		mockService.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
			return s.ID == 1 && s.Version == 1
		})).Return((*models.Supplier)(nil), errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/enable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})

}

func TestSupplierHandler_Disable(t *testing.T) {
	mockService := new(service_mock.MockSupplierService)
	logger := logger.NewLoggerAdapter(logrus.New())
	handler := NewSupplierHandler(mockService, logger)

	t.Run("Sucesso ao desabilitar fornecedor", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		supplierID := int64(1)
		supplierToUpdate := &models.Supplier{
			ID:      supplierID,
			Status:  false,
			Version: 2,
		}

		requestBody := map[string]interface{}{
			"version": 2,
		}

		body, _ := json.Marshal(requestBody)

		mockService.On("GetByID", mock.Anything, supplierID).Return(&models.Supplier{
			ID:      supplierID,
			Status:  true,
			Version: 2,
		}, nil).Once()

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
			return s.ID == supplierID && s.Status == false && s.Version == 2
		})).Return(supplierToUpdate, nil).Once()

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

	t.Run("Erro ao obter fornecedor", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByID", mock.Anything, int64(1)).Return(nil, fmt.Errorf("fornecedor não encontrado")).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/disable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
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
			return s.ID == 1 && s.Version == 1
		})).Return(nil, repo.ErrVersionConflict).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/disable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro genérico ao obter fornecedor", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByID", mock.Anything, int64(1)).Return(nil, errors.New("erro banco")).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/disable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro genérico ao desabilitar fornecedor", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByID", mock.Anything, int64(1)).Return(&models.Supplier{
			ID:      1,
			Status:  true,
			Version: 1,
		}, nil).Once()

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(s *models.Supplier) bool {
			return s.ID == 1 && s.Version == 1
		})).Return((*models.Supplier)(nil), errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodPatch, "/suppliers/1/disable", bytes.NewReader([]byte(`{"version":1}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestSupplierHandler_Delete(t *testing.T) {
	mockService := new(service_mock.MockSupplierService)
	logger := logger.NewLoggerAdapter(logrus.New())
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
