package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	mockProduct "github.com/WagaoCarvalho/backend_store_go/infra/mock/product"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/product/product"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestProductHandler_GetAll(t *testing.T) {

	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockService := new(mockProduct.ProductMock)
		handler := NewProduct(mockService, logAdapter)

		expectedProducts := []*models.Product{
			{ID: 1, ProductName: "Produto 1"},
			{ID: 2, ProductName: "Produto 2"},
		}

		mockService.On("GetAll", mock.Anything, 10, 0).Return(expectedProducts, nil)

		req := httptest.NewRequest(http.MethodGet, "/products?limit=10&offset=0", nil)
		w := httptest.NewRecorder()

		handler.GetAll(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Produtos listados com sucesso", response["message"])
		assert.NotNil(t, response["data"])

		mockService.AssertExpectations(t)
	})

	t.Run("Internal_service_error", func(t *testing.T) {
		t.Parallel()

		mockService := new(mockProduct.ProductMock)
		handler := NewProduct(mockService, logAdapter)

		mockErr := errors.New("erro interno")
		mockService.On("GetAll", mock.Anything, 10, 0).Return(nil, mockErr)

		req := httptest.NewRequest(http.MethodGet, "/products?limit=10&offset=0", nil)
		w := httptest.NewRecorder()

		handler.GetAll(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, response.Status)
		assert.Equal(t, "erro interno", response.Message)
		assert.Nil(t, response.Data)

		mockService.AssertExpectations(t)
	})
}

func TestProductHandler_GetVersionByID(t *testing.T) {

	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockProduct.ProductMock)
		handler := NewProduct(mockService, logAdapter)

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
		handler := NewProduct(mockService, logAdapter)

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
		handler := NewProduct(mockService, logAdapter)

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
		handler := NewProduct(mockService, logAdapter)

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
		handler := NewProduct(mockService, logAdapter)

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
		handler := NewProduct(mockService, logAdapter)

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

func TestProductHandler_GetByName(t *testing.T) {

	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockProduct.ProductMock)
		handler := NewProduct(mockService, logAdapter)

		expectedProducts := []*models.Product{
			{ID: 1, ProductName: "Produto A"},
			{ID: 2, ProductName: "Produto A"},
		}

		mockService.On("GetByName", mock.Anything, "Produto A").Return(expectedProducts, nil)

		nameParam := url.PathEscape("Produto A")
		req := httptest.NewRequest(http.MethodGet, "/product/name/"+nameParam, nil)
		req = mux.SetURLVars(req, map[string]string{"name": "Produto A"})
		w := httptest.NewRecorder()

		handler.GetByName(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Produtos encontrados", response["message"])
		assert.NotNil(t, response["data"])

		mockService.AssertExpectations(t)
	})

	t.Run("MissingNameParam", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockProduct.ProductMock)
		handler := NewProduct(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodGet, "/product/name/", nil)
		req = mux.SetURLVars(req, map[string]string{"name": ""})
		w := httptest.NewRecorder()

		handler.GetByName(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Contains(t, response["message"], "obrigatório")

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockProduct.ProductMock)
		handler := NewProduct(mockService, logAdapter)

		mockService.On("GetByName", mock.Anything, "Inexistente").Return(nil, errors.New("erro interno"))

		req := httptest.NewRequest(http.MethodGet, "/product/name/Inexistente", nil)
		req = mux.SetURLVars(req, map[string]string{"name": "Inexistente"})
		w := httptest.NewRecorder()

		handler.GetByName(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "erro interno", response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("No products found", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockProduct.ProductMock)
		handler := NewProduct(mockService, logAdapter)

		// Mock retornando lista vazia
		mockService.On("GetByName", mock.Anything, "Inexistente").
			Return([]*models.Product{}, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/product/name/Inexistente", nil)
		req = mux.SetURLVars(req, map[string]string{"name": "Inexistente"})
		w := httptest.NewRecorder()

		handler.GetByName(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Status  int              `json:"status"`
			Message string           `json:"message"`
			Data    []dto.ProductDTO `json:"data"`
		}

		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)

		assert.Equal(t, "Nenhum produto encontrado", response.Message)
		assert.Len(t, response.Data, 0)

		mockService.AssertExpectations(t)
	})

}

func TestProductHandler_GetByManufacturer(t *testing.T) {
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
		manufacturer := "Samsung"
		expectedProducts := []models.Product{
			{ID: 1, ProductName: "Galaxy S22", Manufacturer: manufacturer},
			{ID: 2, ProductName: "Galaxy Tab S8", Manufacturer: manufacturer},
		}

		mockService.On("GetByManufacturer", mock.Anything, manufacturer).
			Return(expectedProducts, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/products/manufacturer/"+manufacturer, nil)
		req = mux.SetURLVars(req, map[string]string{"manufacturer": manufacturer})
		w := httptest.NewRecorder()

		handler.GetByManufacturer(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)

	})

	t.Run("Missing manufacturer param", func(t *testing.T) {
		_, handler := setup()

		req := httptest.NewRequest(http.MethodGet, "/products/manufacturer/", nil)
		req = mux.SetURLVars(req, map[string]string{"manufacturer": ""})
		w := httptest.NewRecorder()

		handler.GetByManufacturer(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Products not found", func(t *testing.T) {
		mockService, handler := setup()
		manufacturer := "UnknownBrand"

		mockService.On("GetByManufacturer", mock.Anything, manufacturer).
			Return(nil, errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodGet, "/products/manufacturer/"+manufacturer, nil)
		req = mux.SetURLVars(req, map[string]string{"manufacturer": manufacturer})
		w := httptest.NewRecorder()

		handler.GetByManufacturer(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Internal server error", func(t *testing.T) {
		mockService, handler := setup()
		manufacturer := "LG"

		mockService.On("GetByManufacturer", mock.Anything, manufacturer).
			Return(nil, errors.New("falha no banco")).Once()

		req := httptest.NewRequest(http.MethodGet, "/products/manufacturer/"+manufacturer, nil)
		req = mux.SetURLVars(req, map[string]string{"manufacturer": manufacturer})
		w := httptest.NewRecorder()

		handler.GetByManufacturer(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}
