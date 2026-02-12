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
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductHandler_EnableDiscount(t *testing.T) {
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

	setup := func() (*mockProduct.ProductMock, *productHandler) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logAdapter)
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

	setup := func() (*mockProduct.ProductMock, *productHandler) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logAdapter)
		return mockService, handler
	}

	t.Run("Success", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(1)
		percent := 10.0

		// Apenas retorna nil (sucesso), sem produto
		mockService.On("ApplyDiscount", mock.Anything, productID, percent).
			Return(nil).Once()

		body := bytes.NewBufferString(`{"percent": 10}`)
		req := httptest.NewRequest(http.MethodPatch, "/product/discount/1", body)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.ApplyDiscount(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Verifica resposta padrão
		var response utils.DefaultResponse
		_ = json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, http.StatusOK, response.Status)
		assert.Equal(t, "Desconto aplicado com sucesso", response.Message)

		// Verifica dados retornados
		data, ok := response.Data.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(productID), data["product_id"])
		assert.Equal(t, percent, data["percent"])

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

	t.Run("Invalid percent (< 0)", func(t *testing.T) {
		_, handler := setup()

		body := bytes.NewBufferString(`{"percent": -5}`)
		req := httptest.NewRequest(http.MethodPatch, "/product/discount/1", body)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.ApplyDiscount(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Invalid percent (> 100)", func(t *testing.T) {
		_, handler := setup()

		body := bytes.NewBufferString(`{"percent": 150}`)
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
			Return(errMsg.ErrNotFound).Once()

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

	t.Run("Product discount not allowed", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(1)
		percent := 10.0

		mockService.On("ApplyDiscount", mock.Anything, productID, percent).
			Return(errMsg.ErrProductDiscountNotAllowed).Once()

		body := bytes.NewBufferString(`{"percent": 10}`)
		req := httptest.NewRequest(http.MethodPatch, "/product/discount/1", body)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.ApplyDiscount(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Internal server error", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(1)
		percent := 10.0

		mockService.On("ApplyDiscount", mock.Anything, productID, percent).
			Return(errors.New("db error")).Once()

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

	t.Run("Invalid discount percent error from service", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(1)
		percent := 50.0 // Percentual válido, mas serviço retorna erro

		mockService.On("ApplyDiscount", mock.Anything, productID, percent).
			Return(errMsg.ErrInvalidDiscountPercent).Once()

		body := bytes.NewBufferString(`{"percent": 50}`)
		req := httptest.NewRequest(http.MethodPatch, "/product/discount/1", body)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.ApplyDiscount(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

}
