package handler

import (
	"bytes"
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

func TestProductFilterHandler_Filter(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	setup := func() (*mockProduct.ProductMock, *productFilterHandler) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductFilterHandler(mockService, logAdapter)
		return mockService, handler
	}

	t.Run("erro - método não permitido", func(t *testing.T) {
		mockService, handler := setup()
		req := httptest.NewRequest(http.MethodPost, "/products/filter", nil)
		rec := httptest.NewRecorder()
		handler.Filter(rec, req)
		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
		mockService.AssertNotCalled(t, "Filter")
	})

	t.Run("erro - supplier_id inválido (não numérico) - continua sem erro", func(t *testing.T) {
		mockService, handler := setup()
		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f *filter.ProductFilter) bool {
				return f.SupplierID == nil
			})).
			Return([]*model.Product{}, nil).
			Once()
		req := httptest.NewRequest(http.MethodGet, "/products/filter?supplier_id=abc", nil)
		rec := httptest.NewRecorder()
		handler.Filter(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("erro - status inválido (não booleano) - continua sem erro", func(t *testing.T) {
		mockService, handler := setup()
		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f *filter.ProductFilter) bool {
				return f.Status == nil
			})).
			Return([]*model.Product{}, nil).
			Once()
		req := httptest.NewRequest(http.MethodGet, "/products/filter?status=abc", nil)
		rec := httptest.NewRecorder()
		handler.Filter(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("erro - allow_discount inválido (não booleano) - continua sem erro", func(t *testing.T) {
		mockService, handler := setup()
		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f *filter.ProductFilter) bool {
				return f.AllowDiscount == nil
			})).
			Return([]*model.Product{}, nil).
			Once()
		req := httptest.NewRequest(http.MethodGet, "/products/filter?allow_discount=abc", nil)
		rec := httptest.NewRecorder()
		handler.Filter(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("erro - falha no serviço (erro genérico)", func(t *testing.T) {
		mockService, handler := setup()
		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return(nil, errors.New("db error")).
			Once()
		req := httptest.NewRequest(http.MethodGet, "/products/filter?limit=10&offset=0", nil)
		rec := httptest.NewRecorder()
		handler.Filter(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("erro - filtro inválido (serviço retorna ErrInvalidFilter)", func(t *testing.T) {
		mockService, handler := setup()
		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return(nil, errMsg.ErrInvalidFilter).
			Once()
		req := httptest.NewRequest(http.MethodGet, "/products/filter?limit=10&offset=0", nil)
		rec := httptest.NewRecorder()
		handler.Filter(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("erro - ID zero (serviço retorna ErrZeroID)", func(t *testing.T) {
		mockService, handler := setup()
		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f *filter.ProductFilter) bool {
				return f.SupplierID != nil && *f.SupplierID == 0
			})).
			Return(nil, errMsg.ErrZeroID).
			Once()
		req := httptest.NewRequest(http.MethodGet, "/products/filter?supplier_id=0", nil)
		rec := httptest.NewRecorder()
		handler.Filter(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - retorna lista de produtos", func(t *testing.T) {
		mockService, handler := setup()
		modelProducts := []*model.Product{
			{ID: 1, ProductName: "Produto A", Manufacturer: "Marca A", SalePrice: 100, Status: true, Version: 1},
			{ID: 2, ProductName: "Produto B", Manufacturer: "Marca B", SalePrice: 200, Status: false, Version: 1},
		}
		mockService.On("Filter", mock.Anything, mock.Anything).Return(modelProducts, nil).Once()
		req := httptest.NewRequest(http.MethodGet, "/products/filter?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()
		handler.Filter(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - retorna lista vazia", func(t *testing.T) {
		mockService, handler := setup()
		mockService.On("Filter", mock.Anything, mock.Anything).Return([]*model.Product{}, nil).Once()
		req := httptest.NewRequest(http.MethodGet, "/products/filter?limit=5&offset=0", nil)
		rec := httptest.NewRecorder()
		handler.Filter(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro com supplier_id", func(t *testing.T) {
		mockService, handler := setup()
		mockProducts := []*model.Product{{ID: 1, SupplierID: utils.Int64Ptr(10), ProductName: "Produto fornecedor"}}
		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f *filter.ProductFilter) bool {
				return f.SupplierID != nil && *f.SupplierID == int64(10)
			})).
			Return(mockProducts, nil).
			Once()
		req := httptest.NewRequest(http.MethodGet, "/products/filter?supplier_id=10", nil)
		rec := httptest.NewRecorder()
		handler.Filter(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro com allow_discount true", func(t *testing.T) {
		mockService, handler := setup()
		mockProducts := []*model.Product{{ID: 1, AllowDiscount: true, ProductName: "Produto com desconto"}}
		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f *filter.ProductFilter) bool {
				return f.AllowDiscount != nil && *f.AllowDiscount == true
			})).
			Return(mockProducts, nil).
			Once()
		req := httptest.NewRequest(http.MethodGet, "/products/filter?allow_discount=true", nil)
		rec := httptest.NewRecorder()
		handler.Filter(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro com status true", func(t *testing.T) {
		mockService, handler := setup()
		mockProducts := []*model.Product{{ID: 1, Status: true, ProductName: "Produto Ativo"}}
		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f *filter.ProductFilter) bool {
				return f.Status != nil && *f.Status
			})).
			Return(mockProducts, nil).
			Once()
		req := httptest.NewRequest(http.MethodGet, "/products/filter?status=true", nil)
		rec := httptest.NewRecorder()
		handler.Filter(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro com product_name e manufacturer", func(t *testing.T) {
		mockService, handler := setup()
		mockProducts := []*model.Product{{ID: 1, ProductName: "Notebook Dell", Manufacturer: "Dell"}}
		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f *filter.ProductFilter) bool {
				return f.ProductName == "Notebook" && f.Manufacturer == "Dell"
			})).
			Return(mockProducts, nil).
			Once()
		req := httptest.NewRequest(http.MethodGet, "/products/filter?product_name=Notebook&manufacturer=Dell", nil)
		rec := httptest.NewRecorder()
		handler.Filter(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro com barcode", func(t *testing.T) {
		mockService, handler := setup()
		mockProducts := []*model.Product{{ID: 1, Barcode: utils.StrToPtr("123456789")}}
		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f *filter.ProductFilter) bool {
				return f.Barcode == "123456789"
			})).
			Return(mockProducts, nil).
			Once()
		req := httptest.NewRequest(http.MethodGet, "/products/filter?barcode=123456789", nil)
		rec := httptest.NewRecorder()
		handler.Filter(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - paginação com valores padrão quando não informados", func(t *testing.T) {
		mockService, handler := setup()
		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f *filter.ProductFilter) bool {
				return f.BaseFilter.Limit == 10 && f.BaseFilter.Offset == 0
			})).
			Return([]*model.Product{}, nil).
			Once()
		req := httptest.NewRequest(http.MethodGet, "/products/filter", nil)
		rec := httptest.NewRecorder()
		handler.Filter(rec, req)
		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})
}
