package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockProduct "github.com/WagaoCarvalho/backend_store_go/infra/mock/product"
	modelFilter "github.com/WagaoCarvalho/backend_store_go/internal/model/product/filter"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductHandler_Filter(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("erro - falha no serviço", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logger)

		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return(nil, errors.New("db error"))

		req := httptest.NewRequest(http.MethodGet, "/products/filter?limit=10&offset=0", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("erro - filtro inválido", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logger)

		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return(nil, errMsg.ErrInvalidFilter)

		req := httptest.NewRequest(http.MethodGet, "/products/filter?limit=-1", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - retorna lista de produtos", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logger)

		mockProducts := []*model.Product{
			{
				ID:                 1,
				ProductName:        "Produto X",
				Manufacturer:       "Fabricante A",
				CostPrice:          10.5,
				SalePrice:          15.0,
				StockQuantity:      100,
				AllowDiscount:      true,
				MinDiscountPercent: 5.0,
				MaxDiscountPercent: 20.0,
				Status:             true,
				Version:            2,
			},
			{
				ID:                 2,
				ProductName:        "Produto Y",
				Manufacturer:       "Fabricante B",
				CostPrice:          20.0,
				SalePrice:          25.0,
				StockQuantity:      200,
				AllowDiscount:      false,
				MinDiscountPercent: 0.0,
				MaxDiscountPercent: 10.0,
				Status:             true,
				Version:            3,
			},
		}

		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return(mockProducts, nil)

		req := httptest.NewRequest(http.MethodGet, "/products/filter?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Status)
		assert.Equal(t, "Produtos listados com sucesso", resp.Message)

		data, ok := resp.Data.(map[string]any)
		assert.True(t, ok)
		assert.Equal(t, float64(2), data["total"])
		assert.NotEmpty(t, data["items"])

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro com supplier_id válido", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logger)

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(filter *modelFilter.ProductFilter) bool {
				return filter.SupplierID != nil && *filter.SupplierID == 123
			})).
			Return([]*model.Product{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/products/filter?supplier_id=123", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro com supplier_id inválido (ignorado)", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logger)

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(filter *modelFilter.ProductFilter) bool {
				return filter.SupplierID == nil
			})).
			Return([]*model.Product{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/products/filter?supplier_id=abc", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro com version válido", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logger)

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(filter *modelFilter.ProductFilter) bool {
				return filter.Version != nil && *filter.Version == 2
			})).
			Return([]*model.Product{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/products/filter?version=2", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro com version inválido (ignorado)", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logger)

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(filter *modelFilter.ProductFilter) bool {
				return filter.Version == nil
			})).
			Return([]*model.Product{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/products/filter?version=abc", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro com status true", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logger)

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(filter *modelFilter.ProductFilter) bool {
				return filter.Status != nil && *filter.Status == true
			})).
			Return([]*model.Product{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/products/filter?status=true", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro com status false", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logger)

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(filter *modelFilter.ProductFilter) bool {
				return filter.Status != nil && *filter.Status == false
			})).
			Return([]*model.Product{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/products/filter?status=false", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro com status inválido (ignorado)", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logger)

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(filter *modelFilter.ProductFilter) bool {
				return filter.Status == nil
			})).
			Return([]*model.Product{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/products/filter?status=abc", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro com allow_discount true", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logger)

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(filter *modelFilter.ProductFilter) bool {
				return filter.AllowDiscount != nil && *filter.AllowDiscount == true
			})).
			Return([]*model.Product{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/products/filter?allow_discount=true", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro com allow_discount false", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logger)

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(filter *modelFilter.ProductFilter) bool {
				return filter.AllowDiscount != nil && *filter.AllowDiscount == false
			})).
			Return([]*model.Product{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/products/filter?allow_discount=false", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro com allow_discount inválido (ignorado)", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logger)

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(filter *modelFilter.ProductFilter) bool {
				return filter.AllowDiscount == nil
			})).
			Return([]*model.Product{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/products/filter?allow_discount=abc", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - múltiplos filtros combinados", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logger)

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(filter *modelFilter.ProductFilter) bool {
				return filter.SupplierID != nil && *filter.SupplierID == 456 &&
					filter.Version != nil && *filter.Version == 3 &&
					filter.Status != nil && *filter.Status == true &&
					filter.AllowDiscount != nil && *filter.AllowDiscount == false
			})).
			Return([]*model.Product{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/products/filter?supplier_id=456&version=3&status=true&allow_discount=false", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

}
