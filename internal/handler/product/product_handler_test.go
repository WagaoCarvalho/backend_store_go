package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/product"
	service_mock "github.com/WagaoCarvalho/backend_store_go/internal/service/product/mocks"
	"github.com/WagaoCarvalho/backend_store_go/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/pkg/utils"
)

func TestProductHandler_Create(t *testing.T) {
	// Silenciar logs
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockService := new(service_mock.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		input := &models.Product{
			ProductName:   "Produto A",
			Manufacturer:  "Marca X",
			CostPrice:     10.0,
			SalePrice:     15.0,
			StockQuantity: 100,
		}
		expected := *input
		expected.ID = 1

		mockService.On("Create", mock.Anything, input).Return(&expected, nil)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Produto criado com sucesso", response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		t.Parallel()
		mockService := new(service_mock.ProductServiceMock)
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
		mockService := new(service_mock.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		input := &models.Product{
			ProductName:   "Produto FK",
			Manufacturer:  "Marca FK",
			CostPrice:     10.0,
			SalePrice:     15.0,
			StockQuantity: 10,
		}

		mockService.On("Create", mock.Anything, input).Return((*models.Product)(nil), repo.ErrInvalidForeignKey)

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

	t.Run("ServiceError", func(t *testing.T) {
		t.Parallel()
		mockService := new(service_mock.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		input := &models.Product{
			ProductName:   "Produto Erro",
			Manufacturer:  "Marca Erro",
			CostPrice:     20.0,
			SalePrice:     25.0,
			StockQuantity: 5,
		}

		mockService.On("Create", mock.Anything, input).Return((*models.Product)(nil), errors.New("falha interna"))

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
	// Silenciar logs
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		mockService := new(service_mock.ProductServiceMock)
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

		mockService := new(service_mock.ProductServiceMock)
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
	// Silenciar logs
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockService := new(service_mock.ProductServiceMock)
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
		mockService := new(service_mock.ProductServiceMock)
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
		mockService := new(service_mock.ProductServiceMock)
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

func TestProductHandler_Enable(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	setup := func() (*service_mock.ProductServiceMock, *ProductHandler) {
		mockService := new(service_mock.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)
		return mockService, handler
	}

	t.Run("Success", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(1)

		mockService.On("Enable", mock.Anything, productID).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPatch, "/product/enable/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Enable(w, req)

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

		handler.Enable(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertNotCalled(t, "Enable")
	})

	t.Run("Product not found", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(2)

		mockService.On("Enable", mock.Anything, productID).Return(repo.ErrProductNotFound).Once()

		req := httptest.NewRequest(http.MethodPatch, "/product/enable/2", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		w := httptest.NewRecorder()

		handler.Enable(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Version conflict error", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(3)

		mockService.On("Enable", mock.Anything, productID).Return(repo.ErrVersionConflict).Once()

		req := httptest.NewRequest(http.MethodPatch, "/product/enable/3", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "3"})
		w := httptest.NewRecorder()

		handler.Enable(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Internal server error", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(4)
		mockErr := errors.New("erro interno")

		mockService.On("Enable", mock.Anything, productID).Return(mockErr).Once()

		req := httptest.NewRequest(http.MethodPatch, "/product/enable/4", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "4"})
		w := httptest.NewRecorder()

		handler.Enable(w, req)

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

		handler.Enable(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
		mockService.AssertNotCalled(t, "Enable")
	})
}

func TestProductHandler_Disable(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	setup := func() (*service_mock.ProductServiceMock, *ProductHandler) {
		mockService := new(service_mock.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)
		return mockService, handler
	}

	t.Run("Success", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(1)

		mockService.On("Disable", mock.Anything, productID).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPatch, "/product/disable/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Disable(w, req)

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

		handler.Disable(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertNotCalled(t, "Disable")
	})

	t.Run("Product not found", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(2)

		mockService.On("Disable", mock.Anything, productID).Return(repo.ErrProductNotFound).Once()

		req := httptest.NewRequest(http.MethodPatch, "/product/disable/2", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		w := httptest.NewRecorder()

		handler.Disable(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Version conflict error", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(3)

		mockService.On("Disable", mock.Anything, productID).Return(repo.ErrVersionConflict).Once()

		req := httptest.NewRequest(http.MethodPatch, "/product/disable/3", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "3"})
		w := httptest.NewRecorder()

		handler.Disable(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Internal server error", func(t *testing.T) {
		mockService, handler := setup()
		productID := int64(4)
		mockErr := errors.New("erro interno")

		mockService.On("Disable", mock.Anything, productID).Return(mockErr).Once()

		req := httptest.NewRequest(http.MethodPatch, "/product/disable/4", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "4"})
		w := httptest.NewRecorder()

		handler.Disable(w, req)

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

		handler.Disable(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
		mockService.AssertNotCalled(t, "Disable")
	})
}

func TestProductHandler_GetById(t *testing.T) {
	// Silenciar logs
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockService := new(service_mock.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		expected := &models.Product{
			ID:            1,
			ProductName:   "Produto A",
			Manufacturer:  "Marca X",
			CostPrice:     10.0,
			SalePrice:     15.0,
			StockQuantity: 100,
		}

		mockService.On("GetById", mock.Anything, int64(1)).Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetById(w, req)

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
		mockService := new(service_mock.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodGet, "/products/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetById(w, req)

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
		mockService := new(service_mock.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		mockService.On("GetById", mock.Anything, int64(1)).Return(nil, errors.New("erro interno"))

		req := httptest.NewRequest(http.MethodGet, "/products/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetById(w, req)

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
	// Silenciar logs
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockService := new(service_mock.ProductServiceMock)
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
		assert.Equal(t, "Produtos recuperados com sucesso", response["message"])
		assert.NotNil(t, response["data"])

		mockService.AssertExpectations(t)
	})

	t.Run("MissingNameParam", func(t *testing.T) {
		t.Parallel()
		mockService := new(service_mock.ProductServiceMock)
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
		mockService := new(service_mock.ProductServiceMock)
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
}

func TestProductHandler_GetByManufacturer(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)
	setup := func() (*service_mock.ProductServiceMock, *ProductHandler) {
		mockService := new(service_mock.ProductServiceMock)
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
			Return(nil, errors.New("produtos não encontrados")).Once()

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

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockService := new(service_mock.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		input := &models.Product{
			ProductName:   "Produto Atualizado",
			Manufacturer:  "Marca Nova",
			CostPrice:     20.0,
			SalePrice:     30.0,
			StockQuantity: 50,
		}

		expected := *input
		expected.ID = 1

		mockService.On("Update", mock.Anything, &expected).Return(&expected, nil)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Produto atualizado com sucesso", response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		t.Parallel()
		mockService := new(service_mock.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		body := `{"product_name":"Produto"}`
		req := httptest.NewRequest(http.MethodPut, "/products/abc", bytes.NewBufferString(body))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		t.Parallel()
		mockService := new(service_mock.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewBufferString("invalid-json"))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		t.Parallel()
		mockService := new(service_mock.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		input := &models.Product{
			ProductName:   "Produto X",
			Manufacturer:  "Marca X",
			CostPrice:     12,
			SalePrice:     18,
			StockQuantity: 30,
		}

		expected := *input
		expected.ID = 1

		mockService.On("Update", mock.Anything, &expected).Return(nil, errors.New("erro interno"))

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Update(w, req)

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

func TestProductHandler_Delete(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockService := new(service_mock.ProductServiceMock)
		handler := NewProductHandler(mockService, logAdapter)

		mockService.On("Delete", mock.Anything, int64(1)).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Produto deletado com sucesso", response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		t.Parallel()
		mockService := new(service_mock.ProductServiceMock)
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
		mockService := new(service_mock.ProductServiceMock)
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
