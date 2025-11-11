package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockSupplier "github.com/WagaoCarvalho/backend_store_go/infra/mock/supplier"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/category"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSupplierCategoryHandler_GetByID(t *testing.T) {
	mockSvc := new(mockSupplier.MockSupplierCategory)
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	log := logger.NewLoggerAdapter(baseLogger)
	handler := NewSupplierCategoryHandler(mockSvc, log)

	t.Run("Sucesso ao buscar por ID", func(t *testing.T) {
		expected := &models.SupplierCategory{ID: 1, Name: "Fornecedor X"}

		mockSvc.On("GetByID", mock.Anything, int64(1)).Return(expected, nil)

		req := httptest.NewRequest("GET", "/supplier-categories/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		itemBytes, _ := json.Marshal(resp.Data)
		var result models.SupplierCategory
		err = json.Unmarshal(itemBytes, &result)
		require.NoError(t, err)

		assert.Equal(t, expected.ID, result.ID)
		assert.Equal(t, expected.Name, result.Name)
		assert.Equal(t, "Categoria encontrada com sucesso", resp.Message)

		mockSvc.AssertExpectations(t)
	})

	t.Run("ID inválido (não numérico)", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/supplier-categories/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
	})

	t.Run("Erro ao buscar categoria", func(t *testing.T) {
		mockSvc.On("GetByID", mock.Anything, int64(2)).Return((*models.SupplierCategory)(nil), errors.New("categoria não encontrada"))

		req := httptest.NewRequest("GET", "/supplier-categories/2", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		mockSvc.AssertExpectations(t)
	})

	t.Run("GetByID_NotFound", func(t *testing.T) {
		mockSvc := new(mockSupplier.MockSupplierCategory)
		handler := NewSupplierCategoryHandler(mockSvc, log)

		req := mux.SetURLVars(httptest.NewRequest("GET", "/supplier-categories/999", nil), map[string]string{"id": "999"})
		w := httptest.NewRecorder()

		mockSvc.On("GetByID", mock.Anything, int64(999)).Return((*models.SupplierCategory)(nil), errMsg.ErrNotFound)

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Contains(t, resp.Message, "não encontrado")

		mockSvc.AssertExpectations(t)
	})
}

func TestSupplierCategoryHandler_GetAll(t *testing.T) {
	mockSvc := new(mockSupplier.MockSupplierCategory)
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	log := logger.NewLoggerAdapter(baseLogger)
	handler := NewSupplierCategoryHandler(mockSvc, log)

	t.Run("Sucesso ao buscar todas as categorias", func(t *testing.T) {
		expected := []*models.SupplierCategory{
			{ID: 1, Name: "Cat A"},
			{ID: 2, Name: "Cat B"},
		}

		mockSvc.On("GetAll", mock.Anything).Return(expected, nil)

		req := httptest.NewRequest("GET", "/supplier-categories", nil)
		w := httptest.NewRecorder()

		handler.GetAll(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		itemBytes, _ := json.Marshal(resp.Data)
		var result []*models.SupplierCategory
		err = json.Unmarshal(itemBytes, &result)
		require.NoError(t, err)

		assert.Len(t, result, 2)
		assert.Equal(t, expected[0].ID, result[0].ID)
		assert.Equal(t, "Categorias encontradas com sucesso", resp.Message)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Erro ao buscar categorias", func(t *testing.T) {
		mockSvc := new(mockSupplier.MockSupplierCategory)
		mockSvc.On("GetAll", mock.Anything).Return([]*models.SupplierCategory(nil), errors.New("erro inesperado"))

		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}
		logger := logger.NewLoggerAdapter(baseLogger)
		handler := NewSupplierCategoryHandler(mockSvc, logger)

		req := httptest.NewRequest("GET", "/supplier-categories", nil)
		w := httptest.NewRecorder()

		handler.GetAll(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)

		assert.Equal(t, "erro inesperado", resp.Message)

		mockSvc.AssertExpectations(t)
	})
}
