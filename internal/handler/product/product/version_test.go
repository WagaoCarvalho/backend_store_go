package handler

import (
	"bytes"
	"encoding/json"
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

func TestProductHandler_GetVersionByID(t *testing.T) {
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
		version := int64(5)

		mockService.On("GetVersionByID", mock.Anything, productID).Return(version, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/products/1/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetVersionByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, float64(http.StatusOK), response["status"])
		assert.Equal(t, "Versão do produto recuperada com sucesso", response["message"])

		data, ok := response["data"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(productID), data["product_id"])
		assert.Equal(t, float64(version), data["version"])

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID parameter", func(t *testing.T) {
		mockService, handler := setup()

		req := httptest.NewRequest(http.MethodGet, "/products/abc/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetVersionByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertNotCalled(t, "GetVersionByID")
	})

	t.Run("Product not found - returns 404", func(t *testing.T) {
		mockService, handler := setup()

		productID := int64(999)

		mockService.On("GetVersionByID", mock.Anything, productID).Return(int64(0), errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodGet, "/products/999/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		w := httptest.NewRecorder()

		handler.GetVersionByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID from service - returns 400", func(t *testing.T) {
		mockService, handler := setup()

		productID := int64(0)

		mockService.On("GetVersionByID", mock.Anything, productID).Return(int64(0), errMsg.ErrZeroID).Once()

		req := httptest.NewRequest(http.MethodGet, "/products/0/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "0"})
		w := httptest.NewRecorder()

		handler.GetVersionByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Service error - returns 500", func(t *testing.T) {
		mockService, handler := setup()

		productID := int64(1)
		mockErr := errors.New("erro interno")

		mockService.On("GetVersionByID", mock.Anything, productID).Return(int64(0), mockErr).Once()

		req := httptest.NewRequest(http.MethodGet, "/products/1/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetVersionByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Method not allowed", func(t *testing.T) {
		mockService, handler := setup()

		req := httptest.NewRequest(http.MethodPost, "/products/1/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetVersionByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
		mockService.AssertNotCalled(t, "GetVersionByID")
	})

	t.Run("Returns version zero for valid product", func(t *testing.T) {
		mockService, handler := setup()

		productID := int64(2)
		version := int64(0)

		mockService.On("GetVersionByID", mock.Anything, productID).Return(version, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/products/2/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		w := httptest.NewRecorder()

		handler.GetVersionByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		data, ok := response["data"].(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(0), data["version"])

		mockService.AssertExpectations(t)
	})
}
