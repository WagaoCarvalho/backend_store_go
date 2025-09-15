package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockSupplier "github.com/WagaoCarvalho/backend_store_go/infra/mock/service/supplier"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/supplier_category"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_categories"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSupplierCategoryHandler_Create(t *testing.T) {
	mockSvc := new(mockSupplier.MockSupplierCategoryService)
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)
	handler := NewSupplierCategoryHandler(mockSvc, logger)

	t.Run("Sucesso", func(t *testing.T) {
		categoryDTO := dto.SupplierCategoryDTO{Name: "Alimentos"}
		modelCategory := dto.ToSupplierCategoryModel(categoryDTO)

		mockSvc.On("Create", mock.Anything, modelCategory).Return(modelCategory, nil)

		body, _ := json.Marshal(categoryDTO)
		req := httptest.NewRequest(http.MethodPost, "/supplier-categories", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		var result models.SupplierCategory
		itemBytes, _ := json.Marshal(response.Data)
		err = json.Unmarshal(itemBytes, &result)
		require.NoError(t, err)

		assert.Equal(t, categoryDTO.Name, result.Name)
		assert.Equal(t, "Categoria de fornecedor criada com sucesso", response.Message)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Erro parse JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/supplier-categories", bytes.NewBuffer([]byte("invalid")))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Status)
	})

	t.Run("Erro ao criar categoria", func(t *testing.T) {
		categoryDTO := dto.SupplierCategoryDTO{Name: "Equipamentos"}
		modelCategory := dto.ToSupplierCategoryModel(categoryDTO)

		mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(c *models.SupplierCategory) bool {
			return c.Name == modelCategory.Name
		})).Return(nil, errors.New("erro ao criar categoria"))

		body, _ := json.Marshal(categoryDTO)
		req := httptest.NewRequest("POST", "/aupplier-categories", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Status)
		assert.Equal(t, "erro ao criar categoria", response.Message)

		mockSvc.AssertExpectations(t)
	})

}

func TestSupplierCategoryHandler_GetByID(t *testing.T) {
	mockSvc := new(mockSupplier.MockSupplierCategoryService)
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
		mockSvc.On("GetByID", mock.Anything, int64(2)).Return(nil, errors.New("categoria não encontrada"))

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
		mockSvc := new(mockSupplier.MockSupplierCategoryService)
		handler := NewSupplierCategoryHandler(mockSvc, log)

		req := mux.SetURLVars(httptest.NewRequest("GET", "/supplier-categories/999", nil), map[string]string{"id": "999"})
		w := httptest.NewRecorder()

		mockSvc.On("GetByID", mock.Anything, int64(999)).Return(nil, errMsg.ErrNotFound)

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
	mockSvc := new(mockSupplier.MockSupplierCategoryService)
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
		mockSvc := new(mockSupplier.MockSupplierCategoryService)
		mockSvc.On("GetAll", mock.Anything).Return(nil, errors.New("erro inesperado"))

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

func TestSupplierCategoryHandler_Update(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mockSupplier.MockSupplierCategoryService)
		handler := NewSupplierCategoryHandler(mockSvc, logger)

		categoryDTO := dto.SupplierCategoryDTO{
			Name: "Categoria Atualizada",
		}
		modelCategory := dto.ToSupplierCategoryModel(categoryDTO)
		modelCategory.ID = 123

		body, _ := json.Marshal(categoryDTO)

		mockSvc.On("Update", mock.Anything, mock.MatchedBy(func(c *models.SupplierCategory) bool {
			return c.ID == 123 && c.Name == "Categoria Atualizada"
		})).Return(nil)

		req := mux.SetURLVars(
			httptest.NewRequest("PUT", "/supplier-categories/123", bytes.NewBuffer(body)),
			map[string]string{"id": "123"},
		)

		w := httptest.NewRecorder()
		handler.Update(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Status)
		assert.Equal(t, "Categoria atualizada com sucesso", resp.Message)

		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockSvc := new(mockSupplier.MockSupplierCategoryService)
		handler := NewSupplierCategoryHandler(mockSvc, logger)

		categoryDTO := dto.SupplierCategoryDTO{Name: "Categoria Inválida"}
		body, _ := json.Marshal(categoryDTO)

		req := mux.SetURLVars(
			httptest.NewRequest("PUT", "/supplier-categories/abc", bytes.NewBuffer(body)),
			map[string]string{"id": "abc"},
		)
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.NotEmpty(t, resp.Message)

		mockSvc.AssertExpectations(t) // não deve ter chamadas
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		mockSvc := new(mockSupplier.MockSupplierCategoryService)
		handler := NewSupplierCategoryHandler(mockSvc, logger)

		req := mux.SetURLVars(
			httptest.NewRequest("PUT", "/supplier-categories/123", bytes.NewBuffer([]byte(`{invalid-json}`))),
			map[string]string{"id": "123"},
		)
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.NotEmpty(t, resp.Message)

		mockSvc.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockSvc := new(mockSupplier.MockSupplierCategoryService)
		handler := NewSupplierCategoryHandler(mockSvc, logger)

		categoryDTO := dto.SupplierCategoryDTO{Name: "Não Existe"}
		modelCategory := dto.ToSupplierCategoryModel(categoryDTO)
		modelCategory.ID = 999

		body, _ := json.Marshal(categoryDTO)

		req := mux.SetURLVars(
			httptest.NewRequest("PUT", "/supplier-categories/999", bytes.NewBuffer(body)),
			map[string]string{"id": "999"},
		)
		w := httptest.NewRecorder()

		mockSvc.On("Update", mock.Anything, mock.MatchedBy(func(c *models.SupplierCategory) bool {
			return c.ID == 999
		})).Return(errMsg.ErrNotFound)

		handler.Update(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.Status)
		assert.Contains(t, resp.Message, "não encontrado")

		mockSvc.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockSvc := new(mockSupplier.MockSupplierCategoryService)
		handler := NewSupplierCategoryHandler(mockSvc, logger)

		categoryDTO := dto.SupplierCategoryDTO{Name: "Erro Serviço"}
		modelCategory := dto.ToSupplierCategoryModel(categoryDTO)
		modelCategory.ID = 321

		body, _ := json.Marshal(categoryDTO)

		req := mux.SetURLVars(
			httptest.NewRequest("PUT", "/supplier-categories/321", bytes.NewBuffer(body)),
			map[string]string{"id": "321"},
		)
		w := httptest.NewRecorder()

		mockSvc.On("Update", mock.Anything, mock.MatchedBy(func(c *models.SupplierCategory) bool {
			return c.ID == 321
		})).Return(errors.New("falha interna"))

		handler.Update(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)
		assert.Contains(t, resp.Message, "falha interna")

		mockSvc.AssertExpectations(t)
	})
}

func TestSupplierCategoryHandler_Delete(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mockSupplier.MockSupplierCategoryService)
		handler := NewSupplierCategoryHandler(mockSvc, logger)

		req := mux.SetURLVars(httptest.NewRequest("DELETE", "/supplier-categories/123", nil), map[string]string{"id": "123"})
		w := httptest.NewRecorder()

		mockSvc.On("Delete", mock.Anything, int64(123)).Return(nil)

		handler.Delete(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Empty(t, w.Body.String())

		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockSvc := new(mockSupplier.MockSupplierCategoryService)
		handler := NewSupplierCategoryHandler(mockSvc, logger)

		req := mux.SetURLVars(httptest.NewRequest("DELETE", "/supplier-categories/abc", nil), map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.NotEmpty(t, resp.Message)

		mockSvc.AssertExpectations(t) // Delete não deve ser chamado
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockSvc := new(mockSupplier.MockSupplierCategoryService)
		handler := NewSupplierCategoryHandler(mockSvc, logger)

		req := mux.SetURLVars(httptest.NewRequest("DELETE", "/supplier-categories/456", nil), map[string]string{"id": "456"})
		w := httptest.NewRecorder()

		mockSvc.On("Delete", mock.Anything, int64(456)).Return(errors.New("erro inesperado"))

		handler.Delete(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)
		assert.NotEmpty(t, resp.Message)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Delete_NotFound", func(t *testing.T) {
		mockSvc := new(mockSupplier.MockSupplierCategoryService)
		handler := NewSupplierCategoryHandler(mockSvc, logger)

		req := mux.SetURLVars(httptest.NewRequest("DELETE", "/supplier-categories/999", nil), map[string]string{"id": "999"})
		w := httptest.NewRecorder()

		mockSvc.On("Delete", mock.Anything, int64(999)).Return(errMsg.ErrNotFound)

		handler.Delete(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.Status)
		assert.Contains(t, resp.Message, "não encontrado")

		mockSvc.AssertExpectations(t)
	})
}
