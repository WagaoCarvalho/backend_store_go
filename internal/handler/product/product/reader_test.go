package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockProduct "github.com/WagaoCarvalho/backend_store_go/infra/mock/product"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
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

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logAdapter)

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
		t.Parallel()
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodGet, "/products/abc/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetVersionByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		mockService.AssertNotCalled(t, "GetVersionByID")
	})

	t.Run("Service error", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logAdapter)

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
}

func TestProductHandler_GetByID(t *testing.T) {

	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logAdapter)

		expected := &models.Product{
			ID:            1,
			ProductName:   "Produto A",
			Manufacturer:  "Marca X",
			CostPrice:     10.0,
			SalePrice:     15.0,
			StockQuantity: 100,
		}

		mockService.On("GetByID", mock.Anything, int64(1)).Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Produto recuperado com sucesso", response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidIDParam", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodGet, "/products/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logAdapter)

		mockService.On("GetByID", mock.Anything, int64(1)).Return(nil, errors.New("erro interno"))

		req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.NotEmpty(t, response["message"])

		mockService.AssertExpectations(t)
	})
}
