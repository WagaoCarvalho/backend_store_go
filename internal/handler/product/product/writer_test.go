package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mockProduct "github.com/WagaoCarvalho/backend_store_go/infra/mock/product"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/product/product"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/product/product"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
)

func TestProductHandler_Create(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("Sucesso - Criar Produto", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockProduct.ProductMock)
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
		mockService := new(mockProduct.ProductMock)
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
		mockService := new(mockProduct.ProductMock)
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
			Return((*models.Product)(nil), errMsg.ErrDBInvalidForeignKey).
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
		mockService := new(mockProduct.ProductMock)
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

func TestProductHandler_Update(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	t.Parallel()

	validDTO := dto.ProductDTO{
		ProductName:   "Produto Teste",
		Manufacturer:  "Marca X",
		CostPrice:     10,
		SalePrice:     15,
		StockQuantity: 5,
	}

	t.Run("ID inválido deve retornar 400", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodPut, "/products/abc", bytes.NewBufferString(`{}`))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("falha: JSON inválido deve retornar 400", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logAdapter)

		// Corpo JSON inválido
		req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewBufferString(`{ invalido `))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("dados inválidos (ErrInvalidData) deve retornar 400", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logAdapter)

		reqBody, _ := json.Marshal(validDTO)
		req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewBuffer(reqBody))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockService.On("Update", mock.Anything, mock.Anything).Return(errMsg.ErrInvalidData).Once()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("foreign key inválida (ErrInvalidForeignKey) deve retornar 400", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logAdapter)

		reqBody, _ := json.Marshal(validDTO)
		req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewBuffer(reqBody))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockService.On("Update", mock.Anything, mock.Anything).Return(errMsg.ErrDBInvalidForeignKey).Once()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("ID zero (ErrZeroID) deve retornar 400", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logAdapter)

		reqBody, _ := json.Marshal(validDTO)
		req := httptest.NewRequest(http.MethodPut, "/products/0", bytes.NewBuffer(reqBody))
		req = mux.SetURLVars(req, map[string]string{"id": "0"})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockService.On("Update", mock.Anything, mock.Anything).Return(errMsg.ErrZeroID).Once()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("produto não encontrado (ErrNotFound) deve retornar 404", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logAdapter)

		reqBody, _ := json.Marshal(validDTO)
		req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewBuffer(reqBody))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockService.On("Update", mock.Anything, mock.Anything).Return(errMsg.ErrNotFound).Once()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("conflito de versão (ErrVersionConflict) deve retornar 409", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logAdapter)

		reqBody, _ := json.Marshal(validDTO)
		req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewBuffer(reqBody))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockService.On("Update", mock.Anything, mock.Anything).Return(errMsg.ErrVersionConflict).Once()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusConflict, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("erro genérico do service deve retornar 500", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logAdapter)

		reqBody, _ := json.Marshal(validDTO)
		req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewBuffer(reqBody))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockService.On("Update", mock.Anything, mock.Anything).Return(errors.New("erro genérico")).Once()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - atualizar produto deve retornar 200", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logAdapter)

		id := int64(1)
		dtoInput := validDTO
		expectedModel := dto.ToProductModel(dtoInput)
		expectedModel.ID = id

		reqBody, _ := json.Marshal(dtoInput)
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/products/%d", id), bytes.NewBuffer(reqBody))
		req = mux.SetURLVars(req, map[string]string{"id": fmt.Sprintf("%d", id)})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(p *models.Product) bool {
			return p.ID == id &&
				p.ProductName == dtoInput.ProductName &&
				p.Manufacturer == dtoInput.Manufacturer
		})).Return(nil).Once()

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
		assert.Equal(t, dtoInput.ProductName, response.Data.ProductName)
		assert.Equal(t, dtoInput.Manufacturer, response.Data.Manufacturer)

		mockService.AssertExpectations(t)
	})
}

func TestProductHandler_Delete(t *testing.T) {
	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockProduct.ProductMock)
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
		mockService := new(mockProduct.ProductMock)
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
		mockService := new(mockProduct.ProductMock)
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
