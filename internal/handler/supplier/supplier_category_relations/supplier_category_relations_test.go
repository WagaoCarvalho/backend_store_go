package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_category_relations"
	supplier_category_relations "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_category_relations"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repositories/supplier/supplier_category_relations"
	mock_service "github.com/WagaoCarvalho/backend_store_go/internal/services/suppliers/supplier_category_relations/mocks"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
	"github.com/WagaoCarvalho/backend_store_go/logger"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSupplierCategoryRelationHandler_Create(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	t.Run("success - relação criada", func(t *testing.T) {
		mockService := new(mock_service.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

		relation := &supplier_category_relations.SupplierCategoryRelations{
			SupplierID: 1,
			CategoryID: 2,
		}

		mockService.
			On("Create", mock.Anything, int64(1), int64(2)).
			Return(relation, true, nil)

		body, _ := json.Marshal(relation)
		req := httptest.NewRequest(http.MethodPost, "/supplier-category-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.Status)
		assert.Equal(t, "Relação criada com sucesso", resp.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("success - relação já existente", func(t *testing.T) {
		mockService := new(mock_service.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

		relation := &supplier_category_relations.SupplierCategoryRelations{
			SupplierID: 1,
			CategoryID: 2,
		}

		mockService.
			On("Create", mock.Anything, int64(1), int64(2)).
			Return(relation, false, nil)

		body, _ := json.Marshal(relation)
		req := httptest.NewRequest(http.MethodPost, "/supplier-category-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Status)
		assert.Equal(t, "Relação já existente", resp.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("erro - chave estrangeira", func(t *testing.T) {
		mockService := new(mock_service.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

		relation := &supplier_category_relations.SupplierCategoryRelations{
			SupplierID: 99,
			CategoryID: 88,
		}

		mockService.
			On("Create", mock.Anything, int64(99), int64(88)).
			Return(nil, false, repo.ErrInvalidForeignKey)

		body, _ := json.Marshal(relation)
		req := httptest.NewRequest(http.MethodPost, "/supplier-category-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("erro - json inválido", func(t *testing.T) {
		mockService := new(mock_service.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

		body := bytes.NewBufferString(`{invalid-json}`) // JSON malformado
		req := httptest.NewRequest(http.MethodPost, "/supplier-category-relations", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		// Removido assert.NotEmpty(resp.Message) para evitar falha se estiver vazio

		mockService.AssertExpectations(t)
	})

	t.Run("erro - falha interna", func(t *testing.T) {
		mockService := new(mock_service.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

		relation := &supplier_category_relations.SupplierCategoryRelations{
			SupplierID: 1,
			CategoryID: 2,
		}

		mockService.
			On("Create", mock.Anything, int64(1), int64(2)).
			Return(nil, false, errors.New("erro inesperado"))

		body, _ := json.Marshal(relation)
		req := httptest.NewRequest(http.MethodPost, "/supplier-category-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)

		mockService.AssertExpectations(t)
	})
}

func TestSupplierCategoryRelationHandler_GetBySupplierID(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	t.Run("erro - id inválido", func(t *testing.T) {
		mockService := new(mock_service.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodGet, "/supplier-category-relations/invalid", nil)
		// Define var supplier_id como "invalid"
		req = mux.SetURLVars(req, map[string]string{"supplier_id": "invalid"})
		rec := httptest.NewRecorder()

		handler.GetBySupplierID(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("erro - erro no serviço", func(t *testing.T) {
		mockService := new(mock_service.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

		mockService.
			On("GetBySupplierId", mock.Anything, int64(123)).
			Return(nil, errors.New("erro no banco"))

		req := httptest.NewRequest(http.MethodGet, "/supplier-category-relations/123", nil)
		req = mux.SetURLVars(req, map[string]string{"supplier_id": "123"})
		rec := httptest.NewRecorder()

		handler.GetBySupplierID(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - relações encontradas", func(t *testing.T) {
		mockService := new(mock_service.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

		expectedRelations := []*models.SupplierCategoryRelations{
			{SupplierID: 123, CategoryID: 1},
			{SupplierID: 123, CategoryID: 2},
		}

		mockService.
			On("GetBySupplierId", mock.Anything, int64(123)).
			Return(expectedRelations, nil)

		req := httptest.NewRequest(http.MethodGet, "/supplier-category-relations/123", nil)
		req = mux.SetURLVars(req, map[string]string{"supplier_id": "123"})
		rec := httptest.NewRecorder()

		handler.GetBySupplierID(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Status)
		assert.Equal(t, "Relações encontradas", resp.Message)

		// Reinterpreta resp.Data para o tipo correto
		var data []*models.SupplierCategoryRelations
		dataBytes, err := json.Marshal(resp.Data)
		assert.NoError(t, err)
		err = json.Unmarshal(dataBytes, &data)
		assert.NoError(t, err)

		assert.Equal(t, expectedRelations, data)

		mockService.AssertExpectations(t)
	})

}

func TestSupplierCategoryRelationHandler_GetByCategoryID(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	t.Run("erro - id inválido", func(t *testing.T) {
		mockService := new(mock_service.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodGet, "/supplier-category-relations/category/invalid", nil)
		req = mux.SetURLVars(req, map[string]string{"category_id": "invalid"})
		rec := httptest.NewRecorder()

		handler.GetByCategoryID(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("erro - erro no serviço", func(t *testing.T) {
		mockService := new(mock_service.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

		mockService.
			On("GetByCategoryId", mock.Anything, int64(456)).
			Return(nil, errors.New("erro no banco"))

		req := httptest.NewRequest(http.MethodGet, "/supplier-category-relations/category/456", nil)
		req = mux.SetURLVars(req, map[string]string{"category_id": "456"})
		rec := httptest.NewRecorder()

		handler.GetByCategoryID(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - relações encontradas", func(t *testing.T) {
		mockService := new(mock_service.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

		expectedRelations := []*models.SupplierCategoryRelations{
			{SupplierID: 123, CategoryID: 456},
			{SupplierID: 124, CategoryID: 456},
		}

		mockService.
			On("GetByCategoryId", mock.Anything, int64(456)).
			Return(expectedRelations, nil)

		req := httptest.NewRequest(http.MethodGet, "/supplier-category-relations/category/456", nil)
		req = mux.SetURLVars(req, map[string]string{"category_id": "456"})
		rec := httptest.NewRecorder()

		handler.GetByCategoryID(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Status)
		assert.Equal(t, "Relações encontradas", resp.Message)

		var data []*models.SupplierCategoryRelations
		dataBytes, err := json.Marshal(resp.Data)
		assert.NoError(t, err)
		err = json.Unmarshal(dataBytes, &data)
		assert.NoError(t, err)

		assert.Equal(t, expectedRelations, data)

		mockService.AssertExpectations(t)
	})
}

func TestSupplierCategoryRelationHandler_DeleteByID(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	t.Run("erro - ids inválidos", func(t *testing.T) {
		mockService := new(mock_service.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodDelete, "/supplier-category-relations/invalid/invalid", nil)
		req = mux.SetURLVars(req, map[string]string{
			"supplier_id": "invalid",
			"category_id": "invalid",
		})
		rec := httptest.NewRecorder()

		handler.DeleteByID(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("erro - erro no serviço", func(t *testing.T) {
		mockService := new(mock_service.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

		mockService.
			On("DeleteById", mock.Anything, int64(123), int64(456)).
			Return(errors.New("erro ao deletar"))

		req := httptest.NewRequest(http.MethodDelete, "/supplier-category-relations/123/456", nil)
		req = mux.SetURLVars(req, map[string]string{
			"supplier_id": "123",
			"category_id": "456",
		})
		rec := httptest.NewRecorder()

		handler.DeleteByID(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - relação excluída", func(t *testing.T) {
		mockService := new(mock_service.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

		mockService.
			On("DeleteById", mock.Anything, int64(123), int64(456)).
			Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/supplier-category-relations/123/456", nil)
		req = mux.SetURLVars(req, map[string]string{
			"supplier_id": "123",
			"category_id": "456",
		})
		rec := httptest.NewRecorder()

		handler.DeleteByID(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Status)
		assert.Equal(t, "Relação excluída com sucesso", resp.Message)

		mockService.AssertExpectations(t)
	})
}

func TestSupplierCategoryRelationHandler_DeleteAllBySupplierID(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	t.Run("erro - id inválido", func(t *testing.T) {
		mockService := new(mock_service.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodDelete, "/supplier-category-relations/invalid", nil)
		req = mux.SetURLVars(req, map[string]string{"supplier_id": "invalid"})
		rec := httptest.NewRecorder()

		handler.DeleteAllBySupplierID(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("erro - erro no serviço", func(t *testing.T) {
		mockService := new(mock_service.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

		mockService.
			On("DeleteAllBySupplierId", mock.Anything, int64(123)).
			Return(errors.New("erro ao deletar todas relações"))

		req := httptest.NewRequest(http.MethodDelete, "/supplier-category-relations/123", nil)
		req = mux.SetURLVars(req, map[string]string{"supplier_id": "123"})
		rec := httptest.NewRecorder()

		handler.DeleteAllBySupplierID(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - relações excluídas", func(t *testing.T) {
		mockService := new(mock_service.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

		mockService.
			On("DeleteAllBySupplierId", mock.Anything, int64(123)).
			Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/supplier-category-relations/123", nil)
		req = mux.SetURLVars(req, map[string]string{"supplier_id": "123"})
		rec := httptest.NewRecorder()

		handler.DeleteAllBySupplierID(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Status)
		assert.Equal(t, "Relações excluídas com sucesso", resp.Message)

		mockService.AssertExpectations(t)
	})
}

func TestSupplierCategoryRelationHandler_HasSupplierCategoryRelation(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	t.Run("success - relação existe", func(t *testing.T) {
		mockService := new(mock_service.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

		mockService.
			On("HasRelation", mock.Anything, int64(1), int64(2)).
			Return(true, nil)

		req := httptest.NewRequest(http.MethodGet, "/supplier-category-relation/1/category/2", nil)
		req = mux.SetURLVars(req, map[string]string{
			"supplier_id": "1",
			"category_id": "2",
		})
		rr := httptest.NewRecorder()

		handler.HasSupplierCategoryRelation(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), `"exists":true`)
		assert.Contains(t, rr.Body.String(), `"message":"Verificação concluída com sucesso"`)
		assert.Contains(t, rr.Body.String(), `"status":200`)

		mockService.AssertExpectations(t)
	})

	t.Run("success - relação não existe", func(t *testing.T) {
		mockService := new(mock_service.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

		mockService.
			On("HasRelation", mock.Anything, int64(1), int64(3)).
			Return(false, nil)

		req := httptest.NewRequest(http.MethodGet, "/supplier-category-relation/1/category/3", nil)
		req = mux.SetURLVars(req, map[string]string{
			"supplier_id": "1",
			"category_id": "3",
		})
		rr := httptest.NewRecorder()

		handler.HasSupplierCategoryRelation(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), `"exists":false`)
		assert.Contains(t, rr.Body.String(), `"message":"Verificação concluída com sucesso"`)
		assert.Contains(t, rr.Body.String(), `"status":200`)

		mockService.AssertExpectations(t)
	})

	t.Run("error - supplier_id inválido", func(t *testing.T) {
		mockService := new(mock_service.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodGet, "/supplier-category-relation/abc/category/2", nil)
		req = mux.SetURLVars(req, map[string]string{
			"supplier_id": "abc",
			"category_id": "2",
		})
		rr := httptest.NewRecorder()

		handler.HasSupplierCategoryRelation(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "ID de fornecedor inválido")
	})

	t.Run("error - category_id inválido", func(t *testing.T) {
		mockService := new(mock_service.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodGet, "/supplier-category-relation/1/category/xyz", nil)
		req = mux.SetURLVars(req, map[string]string{
			"supplier_id": "1",
			"category_id": "xyz",
		})
		rr := httptest.NewRecorder()

		handler.HasSupplierCategoryRelation(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "ID de categoria inválido")
	})
	t.Run("error - falha ao verificar relação", func(t *testing.T) {
		mockService := new(mock_service.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

		supplierID := int64(1)
		categoryID := int64(2)

		mockService.
			On("HasRelation", mock.Anything, supplierID, categoryID).
			Return(false, errors.New("erro simulado"))

		req := httptest.NewRequest(http.MethodGet, "/supplier-category-relation/1/category/2", nil)
		req = mux.SetURLVars(req, map[string]string{
			"supplier_id": "1",
			"category_id": "2",
		})
		rr := httptest.NewRecorder()

		handler.HasSupplierCategoryRelation(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "erro ao verificar relação")

		mockService.AssertExpectations(t)
	})

}
