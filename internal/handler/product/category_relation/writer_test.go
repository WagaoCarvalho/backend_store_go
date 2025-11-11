package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockProductCatRel "github.com/WagaoCarvalho/backend_store_go/infra/mock/product"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/product/category_relation"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestProductCategoryRelationHandler_Create(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("success - relação criada", func(t *testing.T) {
		mockService := new(mockProductCatRel.MockProductCategoryRelation)
		handler := NewProductCategoryRelationHandler(mockService, logger)

		dtoRel := dto.ProductCategoryRelationsDTO{
			ProductID:  *utils.Int64Ptr(1),
			CategoryID: *utils.Int64Ptr(2),
		}
		modelRel := dto.ToProductCategoryRelationsModel(dtoRel)

		mockService.
			On("Create", mock.Anything, mock.MatchedBy(func(rel *models.ProductCategoryRelation) bool {
				return rel.ProductID == 1 && rel.CategoryID == 2
			})).
			Return(modelRel, nil).Once()

		body, _ := json.Marshal(dtoRel)
		req := httptest.NewRequest(http.MethodPost, "/relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, "Relação criada com sucesso", resp.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("error - corpo JSON inválido", func(t *testing.T) {
		mockService := new(mockProductCatRel.MockProductCategoryRelation)
		handler := NewProductCategoryRelationHandler(mockService, logger)

		// JSON inválido para forçar erro no parse
		req := httptest.NewRequest(http.MethodPost, "/relations", bytes.NewBufferString("{invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "erro ao decodificar JSON")
	})

	t.Run("error - modelo nulo ou IDs inválidos", func(t *testing.T) {
		mockService := new(mockProductCatRel.MockProductCategoryRelation)
		handler := NewProductCategoryRelationHandler(mockService, logger)

		// DTO com IDs zerados para forçar erro na validação
		dtoRel := dto.ProductCategoryRelationsDTO{
			ProductID:  *utils.Int64Ptr(0),  // ID inválido
			CategoryID: *utils.Int64Ptr(-1), // ID inválido
		}

		body, _ := json.Marshal(dtoRel)
		req := httptest.NewRequest(http.MethodPost, "/relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "modelo nulo ou ID inválido")
	})

	t.Run("error - chave estrangeira inválida", func(t *testing.T) {
		mockService := new(mockProductCatRel.MockProductCategoryRelation)
		handler := NewProductCategoryRelationHandler(mockService, logger)

		dtoRel := dto.ProductCategoryRelationsDTO{
			ProductID:  *utils.Int64Ptr(1),
			CategoryID: *utils.Int64Ptr(999), // ID que não existe
		}

		mockService.
			On("Create", mock.Anything, mock.MatchedBy(func(rel *models.ProductCategoryRelation) bool {
				return rel.ProductID == 1 && rel.CategoryID == 999
			})).
			Return(nil, errMsg.ErrDBInvalidForeignKey).Once()

		body, _ := json.Marshal(dtoRel)
		req := httptest.NewRequest(http.MethodPost, "/relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "chave estrangeira inválida")

		mockService.AssertExpectations(t)
	})

	t.Run("error - relação já existente", func(t *testing.T) {
		mockService := new(mockProductCatRel.MockProductCategoryRelation)
		handler := NewProductCategoryRelationHandler(mockService, logger)

		dtoRel := dto.ProductCategoryRelationsDTO{
			ProductID:  *utils.Int64Ptr(1),
			CategoryID: *utils.Int64Ptr(2),
		}
		modelRel := dto.ToProductCategoryRelationsModel(dtoRel)

		mockService.
			On("Create", mock.Anything, mock.MatchedBy(func(rel *models.ProductCategoryRelation) bool {
				return rel.ProductID == 1 && rel.CategoryID == 2
			})).
			Return(modelRel, errMsg.ErrRelationExists).Once()

		body, _ := json.Marshal(dtoRel)
		req := httptest.NewRequest(http.MethodPost, "/relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, "Relação já existente", resp.Message)
		assert.Equal(t, http.StatusOK, resp.Status)
		assert.NotNil(t, resp.Data)

		mockService.AssertExpectations(t)
	})

	t.Run("error - falha interna no serviço", func(t *testing.T) {
		mockService := new(mockProductCatRel.MockProductCategoryRelation)
		handler := NewProductCategoryRelationHandler(mockService, logger)

		dtoRel := dto.ProductCategoryRelationsDTO{
			ProductID:  *utils.Int64Ptr(1),
			CategoryID: *utils.Int64Ptr(2),
		}

		mockService.
			On("Create", mock.Anything, mock.Anything).
			Return(nil, errors.New("erro inesperado no banco")).Once()

		body, _ := json.Marshal(dtoRel)
		req := httptest.NewRequest(http.MethodPost, "/relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, http.StatusInternalServerError, resp.Status)
		assert.Contains(t, resp.Message, "erro ao criar relação")

		mockService.AssertExpectations(t)
	})

}

func TestProductCategoryRelationHandler_Delete(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("success - relação deletada com sucesso", func(t *testing.T) {
		mockService := new(mockProductCatRel.MockProductCategoryRelation)
		handler := NewProductCategoryRelationHandler(mockService, logger)

		productID := int64(1)
		categoryID := int64(100)

		mockService.
			On("Delete", mock.Anything, productID, categoryID).
			Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/relations/1/100", nil)
		req = mux.SetURLVars(req, map[string]string{
			"product_id":  "1",
			"category_id": "100",
		})
		rr := httptest.NewRecorder()

		handler.Delete(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
		assert.Empty(t, rr.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("bad request - IDs inválidos", func(t *testing.T) {
		mockService := new(mockProductCatRel.MockProductCategoryRelation)
		handler := NewProductCategoryRelationHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodDelete, "/relations/abc/xyz", nil)
		req = mux.SetURLVars(req, map[string]string{
			"product_id":  "abc",
			"category_id": "xyz",
		})
		rr := httptest.NewRecorder()

		handler.Delete(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "IDs inválidos")
	})

	t.Run("internal error - erro ao deletar relação", func(t *testing.T) {
		mockService := new(mockProductCatRel.MockProductCategoryRelation)
		handler := NewProductCategoryRelationHandler(mockService, logger)

		productID := int64(2)
		categoryID := int64(200)

		mockService.
			On("Delete", mock.Anything, productID, categoryID).
			Return(fmt.Errorf("erro ao deletar"))

		req := httptest.NewRequest(http.MethodDelete, "/relations/2/200", nil)
		req = mux.SetURLVars(req, map[string]string{
			"product_id":  "2",
			"category_id": "200",
		})
		rr := httptest.NewRecorder()

		handler.Delete(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "erro ao deletar")
		mockService.AssertExpectations(t)
	})
}

func TestProductCategoryRelationHandler_DeleteAll(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("success - todas as relações deletadas com sucesso", func(t *testing.T) {
		mockService := new(mockProductCatRel.MockProductCategoryRelation)
		handler := NewProductCategoryRelationHandler(mockService, logger)

		productID := int64(1)

		mockService.
			On("DeleteAll", mock.Anything, productID).
			Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/relations/product/1", nil)
		req = mux.SetURLVars(req, map[string]string{
			"product_id": "1",
		})
		rr := httptest.NewRecorder()

		handler.DeleteAll(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
		assert.Empty(t, rr.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("bad request - ID de usuário inválido", func(t *testing.T) {
		mockService := new(mockProductCatRel.MockProductCategoryRelation)
		handler := NewProductCategoryRelationHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodDelete, "/relations/product/abc", nil)
		req = mux.SetURLVars(req, map[string]string{
			"product_id": "abc",
		})
		rr := httptest.NewRecorder()

		handler.DeleteAll(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "ID de usuário inválido")
	})

	t.Run("internal error - erro ao deletar todas as relações", func(t *testing.T) {
		mockService := new(mockProductCatRel.MockProductCategoryRelation)
		handler := NewProductCategoryRelationHandler(mockService, logger)

		productID := int64(2)

		mockService.
			On("DeleteAll", mock.Anything, productID).
			Return(fmt.Errorf("erro ao deletar todas as relações"))

		req := httptest.NewRequest(http.MethodDelete, "/relations/product/2", nil)
		req = mux.SetURLVars(req, map[string]string{
			"product_id": "2",
		})
		rr := httptest.NewRecorder()

		handler.DeleteAll(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "erro ao deletar todas as relações")
		mockService.AssertExpectations(t)
	})
}
