package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockProduct "github.com/WagaoCarvalho/backend_store_go/infra/mock/product"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductHandler_EnableDiscount(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	setup := func() (*mockProduct.ProductMock, *Product) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProduct(mockService, logAdapter)
		return mockService, handler
	}

	t.Run("Success", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(1)

		mockService.On("EnableDiscount", mock.Anything, productID).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPatch, "/product/discount/enable/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.EnableDiscount(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		_, handler := setup()

		req := httptest.NewRequest(http.MethodPatch, "/product/discount/enable/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.EnableDiscount(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Method Not Allowed", func(t *testing.T) {
		_, handler := setup()

		req := httptest.NewRequest(http.MethodGet, "/product/discount/enable/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.EnableDiscount(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("Service returns not found", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(2)

		mockService.On("EnableDiscount", mock.Anything, productID).Return(errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodPatch, "/product/discount/enable/2", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		w := httptest.NewRecorder()

		handler.EnableDiscount(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Produto não encontrado", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(1)

		mockService.On("EnableDiscount", mock.Anything, productID).Return(errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodPatch, "/product/discount/enable/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.EnableDiscount(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro inesperado", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(1)

		mockService.On("EnableDiscount", mock.Anything, productID).Return(fmt.Errorf("erro inesperado")).Once()

		req := httptest.NewRequest(http.MethodPatch, "/product/discount/enable/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.EnableDiscount(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

}

func TestProductHandler_DisableDiscount(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	setup := func() (*mockProduct.ProductMock, *Product) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProduct(mockService, logAdapter)
		return mockService, handler
	}

	t.Run("Success", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(1)

		mockService.On("DisableDiscount", mock.Anything, productID).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPatch, "/product/discount/disable/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.DisableDiscount(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("MethodNotAllowed", func(t *testing.T) {
		_, handler := setup()

		req := httptest.NewRequest(http.MethodGet, "/product/discount/disable/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.DisableDiscount(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("InvalidID", func(t *testing.T) {
		_, handler := setup()

		req := httptest.NewRequest(http.MethodPatch, "/product/discount/disable/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.DisableDiscount(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ProductNotFound", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(99)

		mockService.On("DisableDiscount", mock.Anything, productID).Return(errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodPatch, "/product/discount/disable/99", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "99"})
		w := httptest.NewRecorder()

		handler.DisableDiscount(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("InternalServerError", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(2)

		mockService.On("DisableDiscount", mock.Anything, productID).Return(errors.New("db error")).Once()

		req := httptest.NewRequest(http.MethodPatch, "/product/discount/disable/2", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		w := httptest.NewRecorder()

		handler.DisableDiscount(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestProductHandler_ApplyDiscount(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	setup := func() (*mockProduct.ProductMock, *Product) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProduct(mockService, logAdapter)
		return mockService, handler
	}

	t.Run("Success", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(1)
		percent := 10.0
		expectedProduct := &models.Product{ID: productID, SalePrice: 90.0}

		mockService.On("ApplyDiscount", mock.Anything, productID, percent).
			Return(expectedProduct, nil).Once()

		body := bytes.NewBufferString(`{"percent": 10}`)
		req := httptest.NewRequest(http.MethodPatch, "/product/discount/1", body)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.ApplyDiscount(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var product models.Product
		_ = json.NewDecoder(resp.Body).Decode(&product)
		assert.Equal(t, expectedProduct.ID, product.ID)
		mockService.AssertExpectations(t)
	})

	t.Run("Method not allowed", func(t *testing.T) {
		_, handler := setup()

		req := httptest.NewRequest(http.MethodGet, "/product/discount/1", nil)
		w := httptest.NewRecorder()

		handler.ApplyDiscount(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("Invalid ID", func(t *testing.T) {
		_, handler := setup()

		body := bytes.NewBufferString(`{"percent": 10}`)
		req := httptest.NewRequest(http.MethodPatch, "/product/discount/abc", body)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.ApplyDiscount(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Invalid payload", func(t *testing.T) {
		_, handler := setup()

		body := bytes.NewBufferString(`{"percent": "x"}`)
		req := httptest.NewRequest(http.MethodPatch, "/product/discount/1", body)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.ApplyDiscount(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Product not found", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(1)
		percent := 10.0

		mockService.On("ApplyDiscount", mock.Anything, productID, percent).
			Return(nil, errMsg.ErrNotFound).Once()

		body := bytes.NewBufferString(`{"percent": 10}`)
		req := httptest.NewRequest(http.MethodPatch, "/product/discount/1", body)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.ApplyDiscount(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Internal server error", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(1)
		percent := 10.0

		mockService.On("ApplyDiscount", mock.Anything, productID, percent).
			Return(nil, errors.New("db error")).Once()

		body := bytes.NewBufferString(`{"percent": 10}`)
		req := httptest.NewRequest(http.MethodPatch, "/product/discount/1", body)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.ApplyDiscount(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}
