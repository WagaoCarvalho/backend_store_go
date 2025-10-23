package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockSupplier "github.com/WagaoCarvalho/backend_store_go/infra/mock/service/supplier"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/supplier_category_relation"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier_category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setup() (*mockSupplier.MockSupplierCategoryRelationService, *SupplierCategoryRelation) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	mockService := new(mockSupplier.MockSupplierCategoryRelationService)
	handler := NewSupplierCategoryRelation(mockService, logAdapter)

	return mockService, handler
}

func TestSupplierCategoryRelationHandler_Create(t *testing.T) {
	t.Run("success - relação criada", func(t *testing.T) {
		mockService, handler := setup()

		relationDTO := dto.SupplierCategoryRelationsDTO{
			SupplierID: utils.Int64Ptr(1),
			CategoryID: utils.Int64Ptr(2),
		}
		expectedModel := &model.SupplierCategoryRelation{
			SupplierID: 1,
			CategoryID: 2,
		}

		mockService.On("Create", mock.Anything, int64(1), int64(2)).
			Return(expectedModel, true, nil).Once()

		body, _ := json.Marshal(relationDTO)
		req := httptest.NewRequest(http.MethodPost, "/supplier-category-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Relação criada com sucesso", response.Message)

		dataBytes, _ := json.Marshal(response.Data)
		var got dto.SupplierCategoryRelationsDTO
		_ = json.Unmarshal(dataBytes, &got)

		assert.Equal(t, int64(1), *got.SupplierID)
		assert.Equal(t, int64(2), *got.CategoryID)

		mockService.AssertExpectations(t)
	})

	t.Run("success - relação já existente", func(t *testing.T) {
		mockService, handler := setup()

		relationDTO := dto.SupplierCategoryRelationsDTO{
			SupplierID: utils.Int64Ptr(1),
			CategoryID: utils.Int64Ptr(2),
		}
		expectedModel := &model.SupplierCategoryRelation{
			SupplierID: 1,
			CategoryID: 2,
		}

		mockService.On("Create", mock.Anything, int64(1), int64(2)).
			Return(expectedModel, false, nil).Once()

		body, _ := json.Marshal(relationDTO)
		req := httptest.NewRequest(http.MethodPost, "/supplier-category-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Relação já existente", response.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("error - JSON inválido", func(t *testing.T) {
		mockService, handler := setup()

		req := httptest.NewRequest(http.MethodPost, "/supplier-category-relations", bytes.NewBuffer([]byte("{invalid-json")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("error - chave estrangeira inválida", func(t *testing.T) {
		mockService, handler := setup()

		relationDTO := dto.SupplierCategoryRelationsDTO{
			SupplierID: utils.Int64Ptr(99),
			CategoryID: utils.Int64Ptr(88),
		}

		mockService.On("Create", mock.Anything, int64(99), int64(88)).
			Return(nil, false, errMsg.ErrDBInvalidForeignKey).Once()

		body, _ := json.Marshal(relationDTO)
		req := httptest.NewRequest(http.MethodPost, "/supplier-category-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, response.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("error - falha interna", func(t *testing.T) {
		mockService, handler := setup()

		relationDTO := dto.SupplierCategoryRelationsDTO{
			SupplierID: utils.Int64Ptr(1),
			CategoryID: utils.Int64Ptr(2),
		}

		mockService.On("Create", mock.Anything, int64(1), int64(2)).
			Return(nil, false, errors.New("erro inesperado")).Once()

		body, _ := json.Marshal(relationDTO)
		req := httptest.NewRequest(http.MethodPost, "/supplier-category-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Status)

		mockService.AssertExpectations(t)
	})
}

func TestSupplierCategoryRelationHandler_GetBySupplierID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("erro - id inválido", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelation(mockService, logger)

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
		mockService := new(mockSupplier.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelation(mockService, logger)

		mockService.
			On("GetBySupplierID", mock.Anything, int64(123)).
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
		mockService := new(mockSupplier.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelation(mockService, logger)

		expectedRelations := []*model.SupplierCategoryRelation{
			{SupplierID: 123, CategoryID: 1},
			{SupplierID: 123, CategoryID: 2},
		}

		mockService.
			On("GetBySupplierID", mock.Anything, int64(123)).
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

		// Reinterpreta resp.Data para o tipo de DTO
		var data []dto.SupplierCategoryRelationsDTO
		dataBytes, err := json.Marshal(resp.Data)
		assert.NoError(t, err)
		err = json.Unmarshal(dataBytes, &data)
		assert.NoError(t, err)

		assert.Len(t, data, 2)
		assert.Equal(t, int64(123), *data[0].SupplierID)
		assert.Equal(t, int64(1), *data[0].CategoryID)
		assert.Equal(t, int64(123), *data[1].SupplierID)
		assert.Equal(t, int64(2), *data[1].CategoryID)

		mockService.AssertExpectations(t)
	})

}

func TestSupplierCategoryRelationHandler_GetByCategoryID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("erro - id inválido", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelation(mockService, logger)

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
		mockService := new(mockSupplier.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelation(mockService, logger)

		mockService.
			On("GetByCategoryID", mock.Anything, int64(456)).
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
		mockService := new(mockSupplier.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelation(mockService, logger)

		expectedRelations := []*model.SupplierCategoryRelation{
			{SupplierID: 123, CategoryID: 456},
			{SupplierID: 124, CategoryID: 456},
		}

		mockService.
			On("GetByCategoryID", mock.Anything, int64(456)).
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

		// Reinterpreta resp.Data como DTO
		var data []dto.SupplierCategoryRelationsDTO
		dataBytes, err := json.Marshal(resp.Data)
		assert.NoError(t, err)
		err = json.Unmarshal(dataBytes, &data)
		assert.NoError(t, err)

		assert.Len(t, data, 2)
		assert.Equal(t, int64(123), *data[0].SupplierID)
		assert.Equal(t, int64(456), *data[0].CategoryID)
		assert.Equal(t, int64(124), *data[1].SupplierID)
		assert.Equal(t, int64(456), *data[1].CategoryID)

		mockService.AssertExpectations(t)
	})

}

func TestSupplierCategoryRelationHandler_DeleteByID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("erro - ids inválidos", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelation(mockService, logger)

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
		mockService := new(mockSupplier.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelation(mockService, logger)

		mockService.
			On("DeleteByID", mock.Anything, int64(123), int64(456)).
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
		mockService := new(mockSupplier.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelation(mockService, logger)

		mockService.
			On("DeleteByID", mock.Anything, int64(123), int64(456)).
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
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("erro - id inválido", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelation(mockService, logger)

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
		mockService := new(mockSupplier.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelation(mockService, logger)

		mockService.
			On("DeleteAllBySupplierID", mock.Anything, int64(123)).
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
		mockService := new(mockSupplier.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelation(mockService, logger)

		mockService.
			On("DeleteAllBySupplierID", mock.Anything, int64(123)).
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
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("success - relação existe", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelation(mockService, logger)

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
		mockService := new(mockSupplier.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelation(mockService, logger)

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
		mockService := new(mockSupplier.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelation(mockService, logger)

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
		mockService := new(mockSupplier.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelation(mockService, logger)

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
		mockService := new(mockSupplier.MockSupplierCategoryRelationService)
		handler := NewSupplierCategoryRelation(mockService, logger)

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
