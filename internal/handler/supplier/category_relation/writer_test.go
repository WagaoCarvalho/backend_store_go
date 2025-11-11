package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockSupplier "github.com/WagaoCarvalho/backend_store_go/infra/mock/supplier"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/supplier/category_relation"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setup() (*mockSupplier.MockSupplierCategoryRelation, *supplierCategoryRelationHandler) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	mockService := new(mockSupplier.MockSupplierCategoryRelation)
	handler := NewSupplierCategoryRelationHandler(mockService, logAdapter)

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

		mockService.On("Create", mock.Anything, mock.MatchedBy(func(r *model.SupplierCategoryRelation) bool {
			return r.SupplierID == 1 && r.CategoryID == 2
		})).Return(expectedModel, nil).Once()

		body, _ := json.Marshal(struct {
			Relation *dto.SupplierCategoryRelationsDTO `json:"relation"`
		}{Relation: &relationDTO})

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

	t.Run("error - método não permitido", func(t *testing.T) {
		mockService, handler := setup()

		// Cria uma requisição com método GET (não permitido)
		req := httptest.NewRequest(http.MethodGet, "/supplier-category-relations", nil)
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Contains(t, response.Message, "método GET não permitido")

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

	t.Run("error - relação já existente", func(t *testing.T) {
		mockService, handler := setup()

		relationDTO := dto.SupplierCategoryRelationsDTO{
			SupplierID: utils.Int64Ptr(1),
			CategoryID: utils.Int64Ptr(2),
		}

		mockService.On("Create", mock.Anything, mock.MatchedBy(func(r *model.SupplierCategoryRelation) bool {
			return r.SupplierID == 1 && r.CategoryID == 2
		})).Return(nil, errMsg.ErrRelationExists).Once()

		body, _ := json.Marshal(struct {
			Relation *dto.SupplierCategoryRelationsDTO `json:"relation"`
		}{Relation: &relationDTO})

		req := httptest.NewRequest(http.MethodPost, "/supplier-category-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)

		var response utils.DefaultResponse
		_ = json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, http.StatusConflict, response.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("error - falha interna", func(t *testing.T) {
		mockService, handler := setup()

		relationDTO := dto.SupplierCategoryRelationsDTO{
			SupplierID: utils.Int64Ptr(1),
			CategoryID: utils.Int64Ptr(2),
		}

		mockService.On("Create", mock.Anything, mock.MatchedBy(func(r *model.SupplierCategoryRelation) bool {
			return r.SupplierID == 1 && r.CategoryID == 2
		})).Return(nil, errors.New("erro inesperado")).Once()

		body, _ := json.Marshal(struct {
			Relation *dto.SupplierCategoryRelationsDTO `json:"relation"`
		}{Relation: &relationDTO})

		req := httptest.NewRequest(http.MethodPost, "/supplier-category-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var response utils.DefaultResponse
		_ = json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, http.StatusInternalServerError, response.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("error - chave estrangeira inválida", func(t *testing.T) {
		mockService, handler := setup()

		relationDTO := dto.SupplierCategoryRelationsDTO{
			SupplierID: utils.Int64Ptr(99),
			CategoryID: utils.Int64Ptr(88),
		}

		mockService.On("Create", mock.Anything, mock.MatchedBy(func(r *model.SupplierCategoryRelation) bool {
			return r.SupplierID == 99 && r.CategoryID == 88
		})).Return(nil, errMsg.ErrDBInvalidForeignKey).Once()

		body, _ := json.Marshal(struct {
			Relation *dto.SupplierCategoryRelationsDTO `json:"relation"`
		}{Relation: &relationDTO})

		req := httptest.NewRequest(http.MethodPost, "/supplier-category-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var response utils.DefaultResponse
		_ = json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, http.StatusBadRequest, response.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("error - relation não fornecida", func(t *testing.T) {
		mockService, handler := setup()

		// Corpo da requisição com relation = null
		body, _ := json.Marshal(struct {
			Relation *dto.SupplierCategoryRelationsDTO `json:"relation"`
		}{Relation: nil})

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
		assert.Equal(t, "relation não fornecida", response.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("error - IDs zerados", func(t *testing.T) {
		mockService, handler := setup()

		// DTO com IDs zerados
		relationDTO := dto.SupplierCategoryRelationsDTO{
			SupplierID: utils.Int64Ptr(0),
			CategoryID: utils.Int64Ptr(0),
		}

		mockService.On("Create", mock.Anything, mock.MatchedBy(func(r *model.SupplierCategoryRelation) bool {
			return r.SupplierID == 0 && r.CategoryID == 0
		})).Return(nil, errMsg.ErrZeroID).Once()

		body, _ := json.Marshal(struct {
			Relation *dto.SupplierCategoryRelationsDTO `json:"relation"`
		}{Relation: &relationDTO})

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
		assert.Equal(t, errMsg.ErrZeroID.Error(), response.Message)

		mockService.AssertExpectations(t)
	})

}

func TestSupplierCategoryRelationHandler_Delete(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("erro - ids inválidos", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplierCategoryRelation)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodDelete, "/supplier-category-relations/invalid/invalid", nil)
		req = mux.SetURLVars(req, map[string]string{
			"supplier_id": "invalid",
			"category_id": "invalid",
		})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("erro - erro no serviço", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplierCategoryRelation)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

		mockService.
			On("Delete", mock.Anything, int64(123), int64(456)).
			Return(errors.New("erro ao deletar"))

		req := httptest.NewRequest(http.MethodDelete, "/supplier-category-relations/123/456", nil)
		req = mux.SetURLVars(req, map[string]string{
			"supplier_id": "123",
			"category_id": "456",
		})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - relação excluída", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplierCategoryRelation)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

		mockService.
			On("Delete", mock.Anything, int64(123), int64(456)).
			Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/supplier-category-relations/123/456", nil)
		req = mux.SetURLVars(req, map[string]string{
			"supplier_id": "123",
			"category_id": "456",
		})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

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
		mockService := new(mockSupplier.MockSupplierCategoryRelation)
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
		mockService := new(mockSupplier.MockSupplierCategoryRelation)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

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
		mockService := new(mockSupplier.MockSupplierCategoryRelation)
		handler := NewSupplierCategoryRelationHandler(mockService, logger)

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
