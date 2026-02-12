package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockProductCatRel "github.com/WagaoCarvalho/backend_store_go/infra/mock/product"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductCategoryRelationHandler_GetAllRelationsByProductID(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	setup := func() (*mockProductCatRel.MockProductCategoryRelation, *productCategoryRelationHandler) {
		mockService := new(mockProductCatRel.MockProductCategoryRelation)
		handler := NewProductCategoryRelationHandler(mockService, logAdapter)
		return mockService, handler
	}

	t.Run("success - retorna todas as relações do produto", func(t *testing.T) {
		mockService, handler := setup()

		expected := []*models.ProductCategoryRelation{
			{ProductID: 1, CategoryID: 10},
			{ProductID: 1, CategoryID: 20},
		}

		mockService.
			On("GetAllRelationsByProductID", mock.Anything, int64(1)).
			Return(expected, nil).
			Once()

		req := httptest.NewRequest(http.MethodGet, "/product-category-relations/1", nil)
		req = mux.SetURLVars(req, map[string]string{"product_id": "1"})
		rr := httptest.NewRecorder()

		handler.GetAllRelationsByProductID(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "\"product_id\":1")
		assert.Contains(t, rr.Body.String(), "\"category_id\":10")
		assert.Contains(t, rr.Body.String(), "\"category_id\":20")
		assert.Contains(t, rr.Body.String(), "Relações recuperadas com sucesso")

		mockService.AssertExpectations(t)
	})

	t.Run("success - produto sem relações retorna slice vazio", func(t *testing.T) {
		mockService, handler := setup()

		emptySlice := []*models.ProductCategoryRelation{}
		mockService.
			On("GetAllRelationsByProductID", mock.Anything, int64(1)).
			Return(emptySlice, nil).
			Once()

		req := httptest.NewRequest(http.MethodGet, "/product-category-relations/1", nil)
		req = mux.SetURLVars(req, map[string]string{"product_id": "1"})
		rr := httptest.NewRecorder()

		handler.GetAllRelationsByProductID(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "\"data\":[]")
		assert.Contains(t, rr.Body.String(), "Relações recuperadas com sucesso")
		assert.NotContains(t, rr.Body.String(), "\"category_id\"")

		mockService.AssertExpectations(t)
	})

	t.Run("error - ID inválido (não numérico)", func(t *testing.T) {
		mockService, handler := setup()

		req := httptest.NewRequest(http.MethodGet, "/product-category-relations/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"product_id": "abc"})
		rr := httptest.NewRecorder()

		handler.GetAllRelationsByProductID(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "ID de produto inválido")
		mockService.AssertNotCalled(t, "GetAllRelationsByProductID")
	})

	t.Run("error - ID inválido (zero) - serviço retorna ErrZeroID", func(t *testing.T) {
		mockService, handler := setup()

		mockService.
			On("GetAllRelationsByProductID", mock.Anything, int64(0)).
			Return(([]*models.ProductCategoryRelation)(nil), errMsg.ErrZeroID).
			Once()

		req := httptest.NewRequest(http.MethodGet, "/product-category-relations/0", nil)
		req = mux.SetURLVars(req, map[string]string{"product_id": "0"})
		rr := httptest.NewRecorder()

		handler.GetAllRelationsByProductID(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "ID de produto inválido")
		mockService.AssertExpectations(t)
	})

	t.Run("error - produto não encontrado (ErrNotFound)", func(t *testing.T) {
		mockService, handler := setup()

		mockService.
			On("GetAllRelationsByProductID", mock.Anything, int64(999)).
			Return(([]*models.ProductCategoryRelation)(nil), errMsg.ErrNotFound).
			Once()

		req := httptest.NewRequest(http.MethodGet, "/product-category-relations/999", nil)
		req = mux.SetURLVars(req, map[string]string{"product_id": "999"})
		rr := httptest.NewRecorder()

		handler.GetAllRelationsByProductID(rr, req)

		assert.Equal(t, http.StatusNotFound, rr.Code)
		assert.Contains(t, rr.Body.String(), "produto não encontrado")
		mockService.AssertExpectations(t)
	})

	t.Run("error - falha genérica no serviço", func(t *testing.T) {
		mockService, handler := setup()

		mockService.
			On("GetAllRelationsByProductID", mock.Anything, int64(1)).
			Return(([]*models.ProductCategoryRelation)(nil), errors.New("erro interno do banco")).
			Once()

		req := httptest.NewRequest(http.MethodGet, "/product-category-relations/1", nil)
		req = mux.SetURLVars(req, map[string]string{"product_id": "1"})
		rr := httptest.NewRecorder()

		handler.GetAllRelationsByProductID(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		// Handler retorna mensagem padronizada, não o erro bruto
		assert.Contains(t, rr.Body.String(), "erro ao buscar relações do produto")
		assert.NotContains(t, rr.Body.String(), "erro interno do banco")
		mockService.AssertExpectations(t)
	})

}
