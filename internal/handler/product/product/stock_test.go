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

	setup := func() (*mockProduct.ProductMock, *productHandler) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, logAdapter)
		return mockService, handler
	}

	t.Run("Deve retornar erro quando o método não for PATCH", func(t *testing.T) {
		mockService, handler := setup()

		req := httptest.NewRequest(http.MethodGet, "/products/1/stock", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.UpdateStock(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
		mockService.AssertNotCalled(t, "UpdateStock")
	})

	t.Run("Deve retornar erro quando o body for inválido", func(t *testing.T) {
		mockService, handler := setup()

		req := httptest.NewRequest(http.MethodPatch, "/products/1/stock", strings.NewReader("invalid-json"))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.UpdateStock(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertNotCalled(t, "UpdateStock")
	})

	t.Run("Deve retornar erro quando o ID for inválido", func(t *testing.T) {
		mockService, handler := setup()

		payload := `{"quantity": 10}`
		req := httptest.NewRequest(http.MethodPatch, "/products/abc/stock", strings.NewReader(payload))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.UpdateStock(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertNotCalled(t, "UpdateStock")
	})

	t.Run("Deve atualizar estoque com sucesso", func(t *testing.T) {
		mockService, handler := setup()

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
		mockService, handler := setup()

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

	t.Run("Deve retornar 400 quando ID inválido do serviço", func(t *testing.T) {
		mockService, handler := setup()

		payload := `{"quantity": 10}`
		req := httptest.NewRequest(http.MethodPatch, "/products/0/stock", strings.NewReader(payload))
		req = mux.SetURLVars(req, map[string]string{"id": "0"})
		w := httptest.NewRecorder()

		mockService.On("UpdateStock", mock.Anything, int64(0), 10).Return(errMsg.ErrZeroID).Once()

		handler.UpdateStock(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar 400 quando quantidade inválida", func(t *testing.T) {
		mockService, handler := setup()

		payload := `{"quantity": -5}`
		req := httptest.NewRequest(http.MethodPatch, "/products/1/stock", strings.NewReader(payload))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("UpdateStock", mock.Anything, int64(1), -5).Return(errMsg.ErrInvalidQuantity).Once()

		handler.UpdateStock(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar 500 quando houver erro interno", func(t *testing.T) {
		mockService, handler := setup()

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

	t.Run("Deve retornar 500 quando houver conflito de versão (não tratado)", func(t *testing.T) {
		mockService, handler := setup()

		payload := `{"quantity": 10}`
		req := httptest.NewRequest(http.MethodPatch, "/products/1/stock", strings.NewReader(payload))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("UpdateStock", mock.Anything, int64(1), 10).Return(errMsg.ErrVersionConflict).Once()

		handler.UpdateStock(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		// Handler revisado não trata ErrVersionConflict, retorna 500
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}
func TestProductHandler_IncreaseStock(t *testing.T) {
	newLogger := func() *logger.LogAdapter {
		log := logrus.New()
		log.Out = &bytes.Buffer{}
		return logger.NewLoggerAdapter(log)
	}

	setup := func() (*mockProduct.ProductMock, *productHandler) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, newLogger())
		return mockService, handler
	}

	t.Run("Deve retornar erro quando o método for inválido", func(t *testing.T) {
		mockService, handler := setup()

		req := httptest.NewRequest(http.MethodGet, "/products/1/increase-stock", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.IncreaseStock(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		mockService.AssertNotCalled(t, "IncreaseStock")
	})

	t.Run("Deve retornar erro quando o ID for inválido", func(t *testing.T) {
		mockService, handler := setup()

		req := httptest.NewRequest(http.MethodPatch, "/products/abc/increase-stock", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.IncreaseStock(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "IncreaseStock")
	})

	t.Run("Deve retornar erro quando o body for inválido", func(t *testing.T) {
		mockService, handler := setup()

		req := httptest.NewRequest(http.MethodPatch, "/products/1/increase-stock", strings.NewReader("{invalid-json}"))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.IncreaseStock(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "IncreaseStock")
	})

	t.Run("Deve aumentar estoque com sucesso", func(t *testing.T) {
		mockService, handler := setup()

		body := `{"amount": 5}` // Campo corrigido
		req := httptest.NewRequest(http.MethodPatch, "/products/1/increase-stock", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("IncreaseStock", mock.Anything, int64(1), 5).Return(nil).Once()

		handler.IncreaseStock(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar erro quando produto não encontrado", func(t *testing.T) {
		mockService, handler := setup()

		body := `{"amount": 5}`
		req := httptest.NewRequest(http.MethodPatch, "/products/1/increase-stock", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("IncreaseStock", mock.Anything, int64(1), 5).Return(errMsg.ErrNotFound).Once()

		handler.IncreaseStock(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar 400 quando ID inválido do serviço", func(t *testing.T) {
		mockService, handler := setup()

		body := `{"amount": 5}`
		req := httptest.NewRequest(http.MethodPatch, "/products/0/increase-stock", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "0"})
		w := httptest.NewRecorder()

		mockService.On("IncreaseStock", mock.Anything, int64(0), 5).Return(errMsg.ErrZeroID).Once()

		handler.IncreaseStock(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar 400 quando quantidade inválida", func(t *testing.T) {
		mockService, handler := setup()

		body := `{"amount": 0}` // Quantidade zero é inválida
		req := httptest.NewRequest(http.MethodPatch, "/products/1/increase-stock", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("IncreaseStock", mock.Anything, int64(1), 0).Return(errMsg.ErrInvalidQuantity).Once()

		handler.IncreaseStock(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar 500 quando houver conflito de versão (não tratado)", func(t *testing.T) {
		mockService, handler := setup()

		body := `{"amount": 5}`
		req := httptest.NewRequest(http.MethodPatch, "/products/1/increase-stock", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("IncreaseStock", mock.Anything, int64(1), 5).Return(errMsg.ErrVersionConflict).Once()

		handler.IncreaseStock(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code) // Não mais 409
		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar erro interno para falhas inesperadas", func(t *testing.T) {
		mockService, handler := setup()

		body := `{"amount": 5}`
		req := httptest.NewRequest(http.MethodPatch, "/products/1/increase-stock", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("IncreaseStock", mock.Anything, int64(1), 5).Return(fmt.Errorf("erro inesperado")).Once()

		handler.IncreaseStock(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestProductHandler_DecreaseStock(t *testing.T) {
	newLogger := func() *logger.LogAdapter {
		log := logrus.New()
		log.Out = &bytes.Buffer{}
		return logger.NewLoggerAdapter(log)
	}

	setup := func() (*mockProduct.ProductMock, *productHandler) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, newLogger())
		return mockService, handler
	}

	t.Run("Deve retornar erro quando o método for inválido", func(t *testing.T) {
		mockService, handler := setup()

		req := httptest.NewRequest(http.MethodGet, "/products/1/decrease-stock", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.DecreaseStock(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		mockService.AssertNotCalled(t, "DecreaseStock")
	})

	t.Run("Deve retornar erro quando o ID for inválido", func(t *testing.T) {
		mockService, handler := setup()

		req := httptest.NewRequest(http.MethodPatch, "/products/abc/decrease-stock", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.DecreaseStock(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "DecreaseStock")
	})

	t.Run("Deve retornar erro quando o body for inválido", func(t *testing.T) {
		mockService, handler := setup()

		req := httptest.NewRequest(http.MethodPatch, "/products/1/decrease-stock", strings.NewReader("{invalid-json}"))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.DecreaseStock(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "DecreaseStock")
	})

	t.Run("Deve diminuir estoque com sucesso", func(t *testing.T) {
		mockService, handler := setup()

		body := `{"amount": 5}` // Campo corrigido
		req := httptest.NewRequest(http.MethodPatch, "/products/1/decrease-stock", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("DecreaseStock", mock.Anything, int64(1), 5).Return(nil).Once()

		handler.DecreaseStock(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar erro quando produto não encontrado", func(t *testing.T) {
		mockService, handler := setup()

		body := `{"amount": 5}`
		req := httptest.NewRequest(http.MethodPatch, "/products/1/decrease-stock", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("DecreaseStock", mock.Anything, int64(1), 5).Return(errMsg.ErrNotFound).Once()

		handler.DecreaseStock(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar 400 quando ID inválido do serviço", func(t *testing.T) {
		mockService, handler := setup()

		body := `{"amount": 5}`
		req := httptest.NewRequest(http.MethodPatch, "/products/0/decrease-stock", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "0"})
		w := httptest.NewRecorder()

		mockService.On("DecreaseStock", mock.Anything, int64(0), 5).Return(errMsg.ErrZeroID).Once()

		handler.DecreaseStock(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar 400 quando quantidade inválida", func(t *testing.T) {
		mockService, handler := setup()

		body := `{"amount": 0}` // Quantidade zero é inválida
		req := httptest.NewRequest(http.MethodPatch, "/products/1/decrease-stock", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("DecreaseStock", mock.Anything, int64(1), 0).Return(errMsg.ErrInvalidQuantity).Once()

		handler.DecreaseStock(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar 400 quando estoque insuficiente", func(t *testing.T) {
		mockService, handler := setup()

		body := `{"amount": 100}` // Mais do que tem no estoque
		req := httptest.NewRequest(http.MethodPatch, "/products/1/decrease-stock", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("DecreaseStock", mock.Anything, int64(1), 100).Return(errMsg.ErrInsufficientStock).Once()

		handler.DecreaseStock(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar 500 quando houver conflito de versão (não tratado)", func(t *testing.T) {
		mockService, handler := setup()

		body := `{"amount": 5}`
		req := httptest.NewRequest(http.MethodPatch, "/products/1/decrease-stock", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("DecreaseStock", mock.Anything, int64(1), 5).Return(errMsg.ErrVersionConflict).Once()

		handler.DecreaseStock(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code) // Não mais 409
		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar erro interno para falhas inesperadas", func(t *testing.T) {
		mockService, handler := setup()

		body := `{"amount": 5}`
		req := httptest.NewRequest(http.MethodPatch, "/products/1/decrease-stock", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("DecreaseStock", mock.Anything, int64(1), 5).Return(fmt.Errorf("erro inesperado")).Once()

		handler.DecreaseStock(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestProductHandler_GetStock(t *testing.T) {
	newLogger := func() *logger.LogAdapter {
		log := logrus.New()
		log.Out = &bytes.Buffer{}
		return logger.NewLoggerAdapter(log)
	}

	setup := func() (*mockProduct.ProductMock, *productHandler) {
		mockService := new(mockProduct.ProductMock)
		handler := NewProductHandler(mockService, newLogger())
		return mockService, handler
	}

	t.Run("Deve retornar erro quando o método for inválido", func(t *testing.T) {
		mockService, handler := setup()

		req := httptest.NewRequest(http.MethodPost, "/products/1/stock", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetStock(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		mockService.AssertNotCalled(t, "GetStock")
	})

	t.Run("Deve retornar erro quando o ID for inválido", func(t *testing.T) {
		mockService, handler := setup()

		req := httptest.NewRequest(http.MethodGet, "/products/abc/stock", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetStock(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "GetStock")
	})

	t.Run("Deve retornar 400 quando ID inválido do serviço", func(t *testing.T) {
		mockService, handler := setup()

		req := httptest.NewRequest(http.MethodGet, "/products/0/stock", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "0"})
		w := httptest.NewRecorder()

		mockService.On("GetStock", mock.Anything, int64(0)).Return(0, errMsg.ErrZeroID).Once()

		handler.GetStock(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar erro quando produto não encontrado", func(t *testing.T) {
		mockService, handler := setup()

		req := httptest.NewRequest(http.MethodGet, "/products/1/stock", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("GetStock", mock.Anything, int64(1)).Return(0, errMsg.ErrNotFound).Once()

		handler.GetStock(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar erro interno para falhas inesperadas", func(t *testing.T) {
		mockService, handler := setup()

		req := httptest.NewRequest(http.MethodGet, "/products/1/stock", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("GetStock", mock.Anything, int64(1)).Return(0, fmt.Errorf("erro inesperado")).Once()

		handler.GetStock(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar estoque com sucesso", func(t *testing.T) {
		mockService, handler := setup()

		req := httptest.NewRequest(http.MethodGet, "/products/1/stock", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("GetStock", mock.Anything, int64(1)).Return(20, nil).Once()

		handler.GetStock(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp utils.DefaultResponse
		err := json.NewDecoder(w.Body).Decode(&resp)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.Status)
		assert.Equal(t, "Estoque recuperado com sucesso", resp.Message) // Mensagem corrigida

		data, ok := resp.Data.(map[string]interface{})
		if !ok {
			t.Fatalf("esperava map[string]interface{} em Data, mas veio %T", resp.Data)
		}

		assert.Equal(t, float64(1), data["product_id"])
		assert.Equal(t, float64(20), data["stock_quantity"])

		mockService.AssertExpectations(t)
	})

	t.Run("Deve retornar estoque zero com sucesso", func(t *testing.T) {
		mockService, handler := setup()

		req := httptest.NewRequest(http.MethodGet, "/products/1/stock", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("GetStock", mock.Anything, int64(1)).Return(0, nil).Once()

		handler.GetStock(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp utils.DefaultResponse
		err := json.NewDecoder(w.Body).Decode(&resp)
		assert.NoError(t, err)

		data, ok := resp.Data.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(0), data["stock_quantity"])

		mockService.AssertExpectations(t)
	})
}
