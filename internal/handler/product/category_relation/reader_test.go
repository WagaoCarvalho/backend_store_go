package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockProductCatRel "github.com/WagaoCarvalho/backend_store_go/infra/mock/product"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category_relation"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductCategoryRelationHandler_GetAllRelationsByProductID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)
	t.Run("success - retorna todas as relações do usuário", func(t *testing.T) {
		mockService := new(mockProductCatRel.MockProductCategoryRelation)
		handler := NewProductCategoryRelationHandler(mockService, logger)

		expected := []*models.ProductCategoryRelation{
			{ProductID: 1, CategoryID: 10},
			{ProductID: 1, CategoryID: 20},
		}

		mockService.
			On("GetAllRelationsByProductID", mock.Anything, int64(1)).
			Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/product-category-relations/1", nil)
		req = mux.SetURLVars(req, map[string]string{"product_id": "1"})
		rr := httptest.NewRecorder()

		handler.GetAllRelationsByProductID(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "\"product_id\":1")
		assert.Contains(t, rr.Body.String(), "\"category_id\":1")
		assert.Contains(t, rr.Body.String(), "Relações recuperadas com sucesso")

		mockService.AssertExpectations(t)
	})

	t.Run("error - ID inválido", func(t *testing.T) {
		mockService := new(mockProductCatRel.MockProductCategoryRelation)
		handler := NewProductCategoryRelationHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodGet, "/product-category-relations/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"product_id": "abc"})
		rr := httptest.NewRecorder()

		handler.GetAllRelationsByProductID(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "ID de usuário inválido")
	})

	t.Run("error - falha no serviço", func(t *testing.T) {
		mockService := new(mockProductCatRel.MockProductCategoryRelation)
		handler := NewProductCategoryRelationHandler(mockService, logger)

		mockService.
			On("GetAllRelationsByProductID", mock.Anything, int64(1)).
			Return([]*models.ProductCategoryRelation(nil), errors.New("erro interno"))

		req := httptest.NewRequest(http.MethodGet, "/product-category-relations/1", nil)
		req = mux.SetURLVars(req, map[string]string{"product_id": "1"})
		rr := httptest.NewRecorder()

		handler.GetAllRelationsByProductID(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "erro interno")

		mockService.AssertExpectations(t)
	})
}
