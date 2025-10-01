package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mockProduct "github.com/WagaoCarvalho/backend_store_go/infra/mock/service/product"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/product"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
)

func TestProductHandler_Create(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("Sucesso - Criar Produto", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		input := dto.ProductDTO{
			ProductName:   "Produto OK",
			Manufacturer:  "Marca X",
			CostPrice:     10.0,
			SalePrice:     15.0,
			StockQuantity: 5,
		}

		expectedModel := dto.ToProductModel(input)
		expectedModel.ID = 123 // simula ID atribuído pelo banco

		mockService.
			On("Create", mock.Anything, mock.MatchedBy(func(m *models.Product) bool {
				return m.ProductName == expectedModel.ProductName &&
					m.Manufacturer == expectedModel.Manufacturer &&
					m.CostPrice == expectedModel.CostPrice &&
					m.SalePrice == expectedModel.SalePrice &&
					m.StockQuantity == expectedModel.StockQuantity
			})).
			Return(expectedModel, nil).
			Once()

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response struct {
			Status  int            `json:"status"`
			Message string         `json:"message"`
			Data    dto.ProductDTO `json:"data"`
		}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Produto criado com sucesso", response.Message)
		assert.NotNil(t, response.Data.ID)
		assert.Equal(t, expectedModel.ID, *response.Data.ID)

		mockService.AssertExpectations(t)
	})

	t.Run("JSON inválido deve retornar 400", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer([]byte(`{invalid`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ForeignKey inválida deve retornar 400", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		input := dto.ProductDTO{
			ProductName:   "Produto FK",
			Manufacturer:  "Marca FK",
			CostPrice:     50.0,
			SalePrice:     60.0,
			StockQuantity: 20,
		}

		mockService.
			On("Create", mock.Anything, mock.Anything).
			Return((*models.Product)(nil), errMsg.ErrInvalidForeignKey).
			Once()

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro inesperado no service deve retornar 500", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		input := dto.ProductDTO{
			ProductName:   "Produto Erro",
			Manufacturer:  "Marca",
			CostPrice:     20.0,
			SalePrice:     30.0,
			StockQuantity: 8,
		}

		expectedModel := dto.ToProductModel(input)

		mockService.
			On("Create", mock.Anything, mock.MatchedBy(func(m *models.Product) bool {
				return m.ProductName == expectedModel.ProductName &&
					m.Manufacturer == expectedModel.Manufacturer &&
					m.CostPrice == expectedModel.CostPrice &&
					m.SalePrice == expectedModel.SalePrice &&
					m.StockQuantity == expectedModel.StockQuantity
			})).
			Return((*models.Product)(nil), errors.New("erro inesperado")).
			Once()

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestProductHandler_GetAll(t *testing.T) {

	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

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

		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

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
		mockService := new(mockProduct.ProductServiceMock)
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
		mockService := new(mockProduct.ProductServiceMock)
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
		mockService := new(mockProduct.ProductServiceMock)
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

func TestProductHandler_EnableProduct(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	setup := func() (*mockProduct.ProductServiceMock, *ProductHandler) {
		mockService := new(mockProduct.ProductServiceMock)
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

	setup := func() (*mockProduct.ProductServiceMock, *ProductHandler) {
		mockService := new(mockProduct.ProductServiceMock)
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

func TestProductHandler_GetByID(t *testing.T) {

	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockProduct.ProductServiceMock)
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
		mockService := new(mockProduct.ProductServiceMock)
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
		mockService := new(mockProduct.ProductServiceMock)
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

func TestProductHandler_GetByName(t *testing.T) {

	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

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
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

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
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

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
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

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
	setup := func() (*mockProduct.ProductServiceMock, *ProductHandler) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)
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

func TestProductHandler_Update(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	mockService := new(mockProduct.ProductServiceMock)
	handler := NewProductHandler(mockService, logAdapter)

	t.Run("ID inválido deve retornar 400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/products/abc", bytes.NewBufferString(`{}`))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"}) // força id inválido
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("JSON inválido deve retornar 400", func(t *testing.T) {
		t.Parallel()
		req := httptest.NewRequest(http.MethodPut, "/products/123", bytes.NewBufferString(`{invalid-json}`))
		req = mux.SetURLVars(req, map[string]string{"id": "123"})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Erro no service deve retornar 500", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		input := dto.ProductDTO{ProductName: "Produto Erro"}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPut, "/products/123", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req = mux.SetURLVars(req, map[string]string{"id": "123"})
		w := httptest.NewRecorder()

		mockService.On("Update", mock.Anything, mock.Anything).Return(nil, errors.New("erro genérico"))

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("ForeignKey inválida deve retornar 400", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		input := dto.ProductDTO{ProductName: "Produto FK"}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPut, "/products/123", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req = mux.SetURLVars(req, map[string]string{"id": "123"})
		w := httptest.NewRecorder()

		mockService.On("Update", mock.Anything, mock.Anything).Return(nil, errMsg.ErrInvalidForeignKey)

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Sucesso - Atualizar Produto", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		id := int64(101)
		input := dto.ProductDTO{
			ProductName:   "Produto Atualizado",
			Manufacturer:  "Marca Y",
			CostPrice:     25.0,
			SalePrice:     40.0,
			StockQuantity: 12,
		}

		expectedModel := dto.ToProductModel(input)
		expectedModel.ID = id

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(m *models.Product) bool {
			return m.ID == id &&
				m.ProductName == input.ProductName &&
				m.Manufacturer == input.Manufacturer &&
				m.CostPrice == input.CostPrice &&
				m.SalePrice == input.SalePrice &&
				m.StockQuantity == input.StockQuantity
		})).Return(expectedModel, nil).Once()

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/products/%d", id), bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": strconv.FormatInt(id, 10)})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Status  int            `json:"status"`
			Message string         `json:"message"`
			Data    dto.ProductDTO `json:"data"`
		}

		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Produto atualizado com sucesso", response.Message)
		assert.NotNil(t, response.Data.ID)
		assert.Equal(t, expectedModel.ID, *response.Data.ID)
		assert.Equal(t, input.ProductName, response.Data.ProductName)
		assert.Equal(t, input.Manufacturer, response.Data.Manufacturer)
		assert.Equal(t, input.CostPrice, response.Data.CostPrice)
		assert.Equal(t, input.SalePrice, response.Data.SalePrice)
		assert.Equal(t, input.StockQuantity, response.Data.StockQuantity)

		mockService.AssertExpectations(t)
	})

}

func TestProductHandler_Delete(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)
		productID := int64(1)

		mockService.On("Delete", mock.Anything, productID).
			Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		// Corrigido: Delete retorna 204 No Content, não 200 OK
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)

		// Verifica que não há corpo de resposta (como esperado para 204)
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		assert.Empty(t, body)

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodDelete, "/products/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		mockService.On("Delete", mock.Anything, int64(1)).Return(errors.New("erro interno"))

		req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

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

func TestProductHandler_UpdateStock(t *testing.T) {

	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	t.Run("Deve retornar erro quando o método não for PATCH", func(t *testing.T) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodGet, "/products/1/stock", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.UpdateStock(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("Deve retornar erro quando o body for inválido", func(t *testing.T) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodPatch, "/products/1/stock", strings.NewReader("invalid-json"))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.UpdateStock(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Deve retornar erro quando o ID for inválido", func(t *testing.T) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		payload := `{"quantity": 10}`
		req := httptest.NewRequest(http.MethodPatch, "/products/abc/stock", strings.NewReader(payload))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.UpdateStock(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Deve retornar erro quando o service falhar", func(t *testing.T) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		payload := `{"quantity": 10}`
		req := httptest.NewRequest(http.MethodPatch, "/products/1/stock", strings.NewReader(payload))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("UpdateStock", mock.Anything, int64(1), 10).Return(fmt.Errorf("erro do service")).Once()

		handler.UpdateStock(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve atualizar estoque com sucesso", func(t *testing.T) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		payload := `{"quantity": 10}`
		req := httptest.NewRequest(http.MethodPatch, "/products/1/stock", strings.NewReader(payload))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("UpdateStock", mock.Anything, int64(1), 10).Return(nil).Once()

		handler.UpdateStock(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar 404 quando o produto não for encontrado", func(t *testing.T) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		payload := `{"quantity": 10}`
		req := httptest.NewRequest(http.MethodPatch, "/products/1/stock", strings.NewReader(payload))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("UpdateStock", mock.Anything, int64(1), 10).Return(errMsg.ErrNotFound).Once()

		handler.UpdateStock(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar 409 quando houver conflito de versão", func(t *testing.T) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		payload := `{"quantity": 10}`
		req := httptest.NewRequest(http.MethodPatch, "/products/1/stock", strings.NewReader(payload))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("UpdateStock", mock.Anything, int64(1), 10).Return(errMsg.ErrVersionConflict).Once()

		handler.UpdateStock(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

}

func TestProductHandler_IncreaseStock(t *testing.T) {

	newLogger := func() *logger.LogAdapter {
		log := logrus.New()
		log.Out = &bytes.Buffer{}
		return logger.NewLoggerAdapter(log)
	}

	t.Run("Deve retornar erro quando o método for inválido", func(t *testing.T) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, newLogger())

		req := httptest.NewRequest(http.MethodGet, "/products/1/increase-stock", nil)
		w := httptest.NewRecorder()

		handler.IncreaseStock(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("Deve retornar erro quando o ID for inválido", func(t *testing.T) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, newLogger())

		req := httptest.NewRequest(http.MethodPatch, "/products/abc/increase-stock", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.IncreaseStock(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Deve retornar erro quando o body for inválido", func(t *testing.T) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, newLogger())

		req := httptest.NewRequest(http.MethodPatch, "/products/1/increase-stock", strings.NewReader("{invalid-json}"))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.IncreaseStock(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Deve retornar erro quando produto não encontrado", func(t *testing.T) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, newLogger())

		body := `{"stock_quantity": 5}`
		req := httptest.NewRequest(http.MethodPatch, "/products/1/increase-stock", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("IncreaseStock", mock.Anything, int64(1), 5).Return(errMsg.ErrNotFound)

		handler.IncreaseStock(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar erro de conflito de versão", func(t *testing.T) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, newLogger())

		body := `{"stock_quantity": 5}`
		req := httptest.NewRequest(http.MethodPatch, "/products/1/increase-stock", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("IncreaseStock", mock.Anything, int64(1), 5).Return(errMsg.ErrVersionConflict)

		handler.IncreaseStock(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar erro interno para falhas inesperadas", func(t *testing.T) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, newLogger())

		body := `{"stock_quantity": 5}`
		req := httptest.NewRequest(http.MethodPatch, "/products/1/increase-stock", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("IncreaseStock", mock.Anything, int64(1), 5).Return(fmt.Errorf("erro inesperado"))

		handler.IncreaseStock(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve aumentar estoque com sucesso", func(t *testing.T) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, newLogger())

		body := `{"stock_quantity": 5}`
		req := httptest.NewRequest(http.MethodPatch, "/products/1/increase-stock", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("IncreaseStock", mock.Anything, int64(1), 5).Return(nil)

		handler.IncreaseStock(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestProductHandler_DecreaseStock(t *testing.T) {
	newLogger := func() *logger.LogAdapter {
		log := logrus.New()
		log.Out = &bytes.Buffer{}
		return logger.NewLoggerAdapter(log)
	}

	t.Run("Deve retornar erro quando o método for inválido", func(t *testing.T) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, newLogger())

		req := httptest.NewRequest(http.MethodGet, "/products/1/decrease-stock", nil)
		w := httptest.NewRecorder()

		handler.DecreaseStock(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("Deve retornar erro quando o ID for inválido", func(t *testing.T) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, newLogger())

		req := httptest.NewRequest(http.MethodPatch, "/products/abc/decrease-stock", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.DecreaseStock(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Deve retornar erro quando o body for inválido", func(t *testing.T) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, newLogger())

		req := httptest.NewRequest(http.MethodPatch, "/products/1/decrease-stock", strings.NewReader("{invalid-json}"))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.DecreaseStock(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Deve retornar erro quando produto não encontrado", func(t *testing.T) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, newLogger())

		body := `{"stock_quantity": 5}`
		req := httptest.NewRequest(http.MethodPatch, "/products/1/decrease-stock", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("DecreaseStock", mock.Anything, int64(1), 5).Return(errMsg.ErrNotFound)

		handler.DecreaseStock(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar erro de conflito de versão", func(t *testing.T) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, newLogger())

		body := `{"stock_quantity": 5}`
		req := httptest.NewRequest(http.MethodPatch, "/products/1/decrease-stock", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("DecreaseStock", mock.Anything, int64(1), 5).Return(errMsg.ErrVersionConflict)

		handler.DecreaseStock(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar erro interno para falhas inesperadas", func(t *testing.T) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, newLogger())

		body := `{"stock_quantity": 5}`
		req := httptest.NewRequest(http.MethodPatch, "/products/1/decrease-stock", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("DecreaseStock", mock.Anything, int64(1), 5).Return(fmt.Errorf("erro inesperado"))

		handler.DecreaseStock(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve diminuir estoque com sucesso", func(t *testing.T) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, newLogger())

		body := `{"stock_quantity": 5}`
		req := httptest.NewRequest(http.MethodPatch, "/products/1/decrease-stock", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("DecreaseStock", mock.Anything, int64(1), 5).Return(nil)

		handler.DecreaseStock(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestProductHandler_GetStock(t *testing.T) {

	newLogger := func() *logger.LogAdapter {
		log := logrus.New()
		log.Out = &bytes.Buffer{}
		return logger.NewLoggerAdapter(log)
	}

	t.Run("Deve retornar erro quando o método for inválido", func(t *testing.T) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, newLogger())

		req := httptest.NewRequest(http.MethodPost, "/products/1/stock", nil)
		w := httptest.NewRecorder()

		handler.GetStock(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("Deve retornar erro quando o ID for inválido", func(t *testing.T) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, newLogger())

		req := httptest.NewRequest(http.MethodGet, "/products/abc/stock", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetStock(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Deve retornar erro quando produto não encontrado", func(t *testing.T) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, newLogger())

		req := httptest.NewRequest(http.MethodGet, "/products/1/stock", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("GetStock", mock.Anything, int64(1)).Return(0, errMsg.ErrNotFound)

		handler.GetStock(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar erro interno para falhas inesperadas", func(t *testing.T) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, newLogger())

		req := httptest.NewRequest(http.MethodGet, "/products/1/stock", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("GetStock", mock.Anything, int64(1)).Return(0, fmt.Errorf("erro inesperado"))

		handler.GetStock(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar estoque com sucesso", func(t *testing.T) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, newLogger())

		req := httptest.NewRequest(http.MethodGet, "/products/1/stock", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("GetStock", mock.Anything, int64(1)).Return(20, nil)

		handler.GetStock(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp utils.DefaultResponse
		err := json.NewDecoder(w.Body).Decode(&resp)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.Status)
		assert.Equal(t, "Produtos listados com sucesso", resp.Message)

		data, ok := resp.Data.(map[string]interface{})
		if !ok {
			t.Fatalf("esperava map[string]interface{} em Data, mas veio %T", resp.Data)
		}

		assert.Equal(t, float64(1), data["product_id"])
		assert.Equal(t, float64(20), data["stock_quantity"])

		mockService.AssertExpectations(t)
	})
}

func TestProductHandler_EnableDiscount(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	setup := func() (*mockProduct.ProductServiceMock, *ProductHandler) {
		mockService := new(mockProduct.ProductServiceMock)
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

	setup := func() (*mockProduct.ProductServiceMock, *ProductHandler) {
		mockService := new(mockProduct.ProductServiceMock)
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

	setup := func() (*mockProduct.ProductServiceMock, *ProductHandler) {
		mockService := new(mockProduct.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)
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
