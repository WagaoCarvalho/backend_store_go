package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockProduct "github.com/WagaoCarvalho/backend_store_go/infra/mock/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductHandler_EnableProduct(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	setup := func() (*mockProduct.ProductMock, *productHandler) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logAdapter)
		return mockService, handler
	}

	t.Run("Success", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(1)

		mockService.On("EnableProduct", mock.Anything, productID).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPatch, "/product/enable/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.EnableProduct(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID parameter", func(t *testing.T) {
		mockService, handler := setup()

		req := httptest.NewRequest(http.MethodPatch, "/product/enable/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.EnableProduct(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertNotCalled(t, "EnableProduct")
	})

	t.Run("Product not found", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(2)

		mockService.On("EnableProduct", mock.Anything, productID).Return(errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodPatch, "/product/enable/2", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		w := httptest.NewRecorder()

		handler.EnableProduct(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Version conflict error", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(3)

		mockService.On("EnableProduct", mock.Anything, productID).Return(errMsg.ErrVersionConflict).Once()

		req := httptest.NewRequest(http.MethodPatch, "/product/enable/3", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "3"})
		w := httptest.NewRecorder()

		handler.EnableProduct(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Internal server error", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(4)
		mockErr := errors.New("erro interno")

		mockService.On("EnableProduct", mock.Anything, productID).Return(mockErr).Once()

		req := httptest.NewRequest(http.MethodPatch, "/product/enable/4", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "4"})
		w := httptest.NewRecorder()

		handler.EnableProduct(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Method not allowed", func(t *testing.T) {
		mockService, handler := setup()

		req := httptest.NewRequest(http.MethodGet, "/product/enable/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.EnableProduct(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
		mockService.AssertNotCalled(t, "EnableProduct")
	})
}

func TestProductHandler_DisableProduct(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	setup := func() (*mockProduct.ProductMock, *productHandler) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logAdapter)
		return mockService, handler
	}

	t.Run("Success", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(1)

		mockService.On("DisableProduct", mock.Anything, productID).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPatch, "/product/disable/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.DisableProduct(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID parameter", func(t *testing.T) {
		mockService, handler := setup()

		req := httptest.NewRequest(http.MethodPatch, "/product/disable/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.DisableProduct(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertNotCalled(t, "DisableProduct")
	})

	t.Run("Product not found", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(2)

		mockService.On("DisableProduct", mock.Anything, productID).Return(errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodPatch, "/product/disable/2", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		w := httptest.NewRecorder()

		handler.DisableProduct(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Version conflict error", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(3)

		mockService.On("DisableProduct", mock.Anything, productID).Return(errMsg.ErrVersionConflict).Once()

		req := httptest.NewRequest(http.MethodPatch, "/product/disable/3", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "3"})
		w := httptest.NewRecorder()

		handler.DisableProduct(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Internal server error", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(4)
		mockErr := errors.New("erro interno")

		mockService.On("DisableProduct", mock.Anything, productID).Return(mockErr).Once()

		req := httptest.NewRequest(http.MethodPatch, "/product/disable/4", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "4"})
		w := httptest.NewRecorder()

		handler.DisableProduct(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Method not allowed", func(t *testing.T) {
		mockService, handler := setup()

		req := httptest.NewRequest(http.MethodGet, "/product/disable/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.DisableProduct(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
		mockService.AssertNotCalled(t, "DisableProduct")
	})
}
