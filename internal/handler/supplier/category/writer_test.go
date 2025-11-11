package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockSupplier "github.com/WagaoCarvalho/backend_store_go/infra/mock/supplier"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/category"
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

func TestSupplierCategoryHandler_Create(t *testing.T) {
	mockSvc := new(mockSupplier.MockSupplierCategory)
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

func TestSupplierCategoryHandler_Update(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success", func(t *testing.T) {
		mockSvc := new(mockSupplier.MockSupplierCategory)
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
		mockSvc := new(mockSupplier.MockSupplierCategory)
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
		mockSvc := new(mockSupplier.MockSupplierCategory)
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
		mockSvc := new(mockSupplier.MockSupplierCategory)
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
		mockSvc := new(mockSupplier.MockSupplierCategory)
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
		mockSvc := new(mockSupplier.MockSupplierCategory)
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
		mockSvc := new(mockSupplier.MockSupplierCategory)
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
		mockSvc := new(mockSupplier.MockSupplierCategory)
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
		mockSvc := new(mockSupplier.MockSupplierCategory)
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
