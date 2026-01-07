package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockProduct "github.com/WagaoCarvalho/backend_store_go/infra/mock/product"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/product/filter"
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
		handler := NewProductFilterHandler(mockService, logger)

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
		handler := NewProductFilterHandler(mockService, logger)

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
		handler := NewProductFilterHandler(mockService, logger)

		modelProducts := []*model.Product{
			{
				ID:           1,
				ProductName:  "Produto A",
				Manufacturer: "Marca A",
				SalePrice:    100,
				Status:       true,
				Version:      1,
			},
			{
				ID:           2,
				ProductName:  "Produto B",
				Manufacturer: "Marca B",
				SalePrice:    200,
				Status:       false,
				Version:      1,
			},
		}

		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return(modelProducts, nil)

		req := httptest.NewRequest(http.MethodGet, "/products/filter?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Status)
		assert.Equal(t, "Produtos listados com sucesso", resp.Message)

		data := resp.Data.(map[string]any)
		assert.Equal(t, float64(2), data["total"])
		assert.NotEmpty(t, data["items"])

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro com supplier_id", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductFilterHandler(mockService, logger)

		mockProducts := []*model.Product{
			{
				ID:          1,
				SupplierID:  utils.Int64Ptr(10),
				ProductName: "Produto fornecedor",
			},
		}

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f *filter.ProductFilter) bool {
				return f.SupplierID != nil && *f.SupplierID == int64(10)
			})).
			Return(mockProducts, nil)

		req := httptest.NewRequest(http.MethodGet, "/products/filter?supplier_id=10", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		data := resp.Data.(map[string]any)
		assert.Equal(t, float64(1), data["total"])

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro com allow_discount true", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductFilterHandler(mockService, logger)

		mockProducts := []*model.Product{
			{
				ID:            1,
				ProductName:   "Produto com desconto",
				AllowDiscount: true,
			},
		}

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f *filter.ProductFilter) bool {
				return f.AllowDiscount != nil && *f.AllowDiscount == true
			})).
			Return(mockProducts, nil)

		req := httptest.NewRequest(http.MethodGet, "/products/filter?allow_discount=true", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		data := resp.Data.(map[string]any)
		assert.Equal(t, float64(1), data["total"])

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - retorna lista vazia", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductFilterHandler(mockService, logger)

		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return([]*model.Product{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/products/filter?limit=5&offset=0", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Status)

		data := resp.Data.(map[string]any)
		assert.Equal(t, float64(0), data["total"])

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro com status true", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductFilterHandler(mockService, logger)

		mockProducts := []*model.Product{
			{
				ID:          1,
				ProductName: "Produto Ativo",
				Status:      true,
			},
		}

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f *filter.ProductFilter) bool {
				return f.Status != nil && *f.Status
			})).
			Return(mockProducts, nil)

		req := httptest.NewRequest(http.MethodGet, "/products/filter?status=true", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		data := resp.Data.(map[string]any)
		assert.Equal(t, float64(1), data["total"])

		items := data["items"].([]any)
		assert.Equal(t, "Produto Ativo", items[0].(map[string]any)["product_name"])

		mockService.AssertExpectations(t)
	})
}
