package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestProductHandler_UpdateStock(t *testing.T) {

	log := logrus.New()
	log.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(log)

	t.Run("Deve retornar erro quando o método não for PATCH", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
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
		mockService := new(mockProduct.ProductMock)
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
		mockService := new(mockProduct.ProductMock)
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
		mockService := new(mockProduct.ProductMock)
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
		mockService := new(mockProduct.ProductMock)
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
		mockService := new(mockProduct.ProductMock)
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
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logAdapter)

		payload := `{"quantity": 10}`
		req := httptest.NewRequest(http.MethodPatch, "/products/1/stock", strings.NewReader(payload))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("UpdateStock", mock.Anything, int64(1), 10).Return(errMsg.ErrZeroVersion).Once()

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
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, newLogger())

		req := httptest.NewRequest(http.MethodGet, "/products/1/increase-stock", nil)
		w := httptest.NewRecorder()

		handler.IncreaseStock(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("Deve retornar erro quando o ID for inválido", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, newLogger())

		req := httptest.NewRequest(http.MethodPatch, "/products/abc/increase-stock", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.IncreaseStock(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Deve retornar erro quando o body for inválido", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, newLogger())

		req := httptest.NewRequest(http.MethodPatch, "/products/1/increase-stock", strings.NewReader("{invalid-json}"))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.IncreaseStock(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Deve retornar erro quando produto não encontrado", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
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
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, newLogger())

		body := `{"stock_quantity": 5}`
		req := httptest.NewRequest(http.MethodPatch, "/products/1/increase-stock", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("IncreaseStock", mock.Anything, int64(1), 5).Return(errMsg.ErrZeroVersion)

		handler.IncreaseStock(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar erro interno para falhas inesperadas", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
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
		mockService := new(mockProduct.ProductMock)
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
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, newLogger())

		req := httptest.NewRequest(http.MethodGet, "/products/1/decrease-stock", nil)
		w := httptest.NewRecorder()

		handler.DecreaseStock(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("Deve retornar erro quando o ID for inválido", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, newLogger())

		req := httptest.NewRequest(http.MethodPatch, "/products/abc/decrease-stock", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.DecreaseStock(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Deve retornar erro quando o body for inválido", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, newLogger())

		req := httptest.NewRequest(http.MethodPatch, "/products/1/decrease-stock", strings.NewReader("{invalid-json}"))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.DecreaseStock(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Deve retornar erro quando produto não encontrado", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
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
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, newLogger())

		body := `{"stock_quantity": 5}`
		req := httptest.NewRequest(http.MethodPatch, "/products/1/decrease-stock", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("DecreaseStock", mock.Anything, int64(1), 5).Return(errMsg.ErrZeroVersion)

		handler.DecreaseStock(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar erro interno para falhas inesperadas", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
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
		mockService := new(mockProduct.ProductMock)
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
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, newLogger())

		req := httptest.NewRequest(http.MethodPost, "/products/1/stock", nil)
		w := httptest.NewRecorder()

		handler.GetStock(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
	})

	t.Run("Deve retornar erro quando o ID for inválido", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, newLogger())

		req := httptest.NewRequest(http.MethodGet, "/products/abc/stock", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetStock(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("Deve retornar erro quando produto não encontrado", func(t *testing.T) {
		mockService := new(mockProduct.ProductMock)
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
		mockService := new(mockProduct.ProductMock)
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
		mockService := new(mockProduct.ProductMock)
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
