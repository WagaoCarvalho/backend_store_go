package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockService "github.com/WagaoCarvalho/backend_store_go/infra/mock/sale"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/sale/item"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/item"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaleItemHandler_GetByID(t *testing.T) {
	ctx := context.Background()
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	log := logger.NewLoggerAdapter(baseLogger)
	mockService := new(mockService.MockSaleItem)
	handler := NewSaleItemHandler(mockService, log)

	item := &model.SaleItem{
		ID:          1,
		SaleID:      10,
		ProductID:   20,
		Quantity:    2,
		UnitPrice:   50.0,
		Discount:    5.0,
		Tax:         2.5,
		Subtotal:    97.5,
		Description: "Produto teste",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	t.Run("sucesso - item encontrado", func(t *testing.T) {
		mockService.On("GetByID", mock.Anything, int64(1)).Return(item, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/sale-items/1", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)

		// Converter o Data para SaleItemDTO
		itemData, err := json.Marshal(respBody.Data)
		assert.NoError(t, err)

		var itemDTO dto.SaleItemDTO
		err = json.Unmarshal(itemData, &itemDTO)
		assert.NoError(t, err)

		assert.Equal(t, int64(1), *itemDTO.ID)
		mockService.AssertExpectations(t)
	})

	t.Run("erro - ID inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/sale-items/0", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "0"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("erro - método não permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/sale-items/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("erro - serviço retorna erro", func(t *testing.T) {
		expectedErr := fmt.Errorf("erro ao buscar item no banco de dados")
		mockService.On("GetByID", mock.Anything, int64(1)).Return(nil, expectedErr).Once()

		req := httptest.NewRequest(http.MethodGet, "/sale-items/1", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		// Verificar se a resposta de erro está correta
		var errorResp map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&errorResp)
		assert.NoError(t, err)

		// A estrutura exata depende da sua utils.ErrorResponse
		// Normalmente retorna status e message
		assert.Contains(t, errorResp, "status")
		assert.Contains(t, errorResp, "message")

		mockService.AssertExpectations(t)
	})

	t.Run("erro - ID não numérico", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/sale-items/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("erro - ID negativo", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/sale-items/-1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "-1"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestSaleItemHandler_GetBySaleID(t *testing.T) {
	ctx := context.Background()
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	log := logger.NewLoggerAdapter(baseLogger)
	mockService := new(mockService.MockSaleItem)
	handler := NewSaleItemHandler(mockService, log)

	items := []*model.SaleItem{
		{
			ID:          1,
			SaleID:      10,
			ProductID:   20,
			Quantity:    2,
			UnitPrice:   50.0,
			Discount:    5.0,
			Tax:         2.5,
			Subtotal:    97.5,
			Description: "Produto teste",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Primeiro, vamos descobrir quais são os valores padrão da função GetPaginationParams
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	defaultLimit, defaultOffset := utils.GetPaginationParams(req)

	t.Run("sucesso - listar itens por venda", func(t *testing.T) {
		mockService.On("GetBySaleID", mock.Anything, int64(10), defaultLimit, defaultOffset).Return(items, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/sale-items?sale_id=10", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"sale_id": "10"})

		w := httptest.NewRecorder()
		handler.GetBySaleID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, respBody.Status)
		assert.Equal(t, "Itens da venda recuperados com sucesso", respBody.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("erro - sale_id inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/sale-items?sale_id=0", nil)
		req = mux.SetURLVars(req, map[string]string{"sale_id": "0"})
		w := httptest.NewRecorder()

		handler.GetBySaleID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("erro - sale_id não numérico", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/sale-items?sale_id=abc", nil)
		req = mux.SetURLVars(req, map[string]string{"sale_id": "abc"})
		w := httptest.NewRecorder()

		handler.GetBySaleID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("erro - sale_id negativo", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/sale-items?sale_id=-1", nil)
		req = mux.SetURLVars(req, map[string]string{"sale_id": "-1"})
		w := httptest.NewRecorder()

		handler.GetBySaleID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("erro - método não permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/sale-items?sale_id=10", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"sale_id": "10"})
		w := httptest.NewRecorder()

		handler.GetBySaleID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

		var errorResp map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Contains(t, errorResp, "status")
		assert.Contains(t, errorResp, "message")
	})

	t.Run("erro - serviço retorna erro interno", func(t *testing.T) {
		expectedErr := fmt.Errorf("erro de conexão com o banco de dados")
		mockService.On("GetBySaleID", mock.Anything, int64(10), defaultLimit, defaultOffset).Return(nil, expectedErr).Once()

		req := httptest.NewRequest(http.MethodGet, "/sale-items?sale_id=10", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"sale_id": "10"})

		w := httptest.NewRecorder()
		handler.GetBySaleID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var errorResp map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Contains(t, errorResp, "status")
		assert.Contains(t, errorResp, "message")

		mockService.AssertExpectations(t)
	})

	t.Run("erro - serviço retorna erro com sale_id específico", func(t *testing.T) {
		expectedErr := fmt.Errorf("venda não encontrada")
		mockService.On("GetBySaleID", mock.Anything, int64(999), defaultLimit, defaultOffset).Return(nil, expectedErr).Once()

		req := httptest.NewRequest(http.MethodGet, "/sale-items?sale_id=999", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"sale_id": "999"})

		w := httptest.NewRecorder()
		handler.GetBySaleID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - com paginação personalizada", func(t *testing.T) {
		customLimit := 50
		customOffset := 10
		mockService.On("GetBySaleID", mock.Anything, int64(10), customLimit, customOffset).Return(items, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/sale-items?sale_id=10", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"sale_id": "10"})

		q := req.URL.Query()
		q.Add("limit", "50")
		q.Add("offset", "10")
		req.URL.RawQuery = q.Encode()

		w := httptest.NewRecorder()
		handler.GetBySaleID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - sem parâmetros de paginação (valores padrão)", func(t *testing.T) {
		// Usar os valores padrão que descobrimos no início
		mockService.On("GetBySaleID", mock.Anything, int64(10), defaultLimit, defaultOffset).Return(items, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/sale-items?sale_id=10", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"sale_id": "10"})
		// Sem parâmetros limit e offset

		w := httptest.NewRecorder()
		handler.GetBySaleID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		mockService.AssertExpectations(t)
	})
}

func TestSaleItemHandler_GetByProductID(t *testing.T) {
	ctx := context.Background()
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	log := logger.NewLoggerAdapter(baseLogger)
	mockService := new(mockService.MockSaleItem)
	handler := NewSaleItemHandler(mockService, log)

	items := []*model.SaleItem{
		{
			ID:          1,
			SaleID:      10,
			ProductID:   20,
			Quantity:    2,
			UnitPrice:   50.0,
			Discount:    5.0,
			Tax:         2.5,
			Subtotal:    97.5,
			Description: "Produto teste",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	// Descobrir os valores padrão da função GetPaginationParams
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	defaultLimit, defaultOffset := utils.GetPaginationParams(req)

	t.Run("sucesso - listar itens por produto", func(t *testing.T) {
		mockService.On("GetByProductID", mock.Anything, int64(20), defaultLimit, defaultOffset).Return(items, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/sale-items?product_id=20", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"product_id": "20"})

		w := httptest.NewRecorder()
		handler.GetByProductID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, respBody.Status)
		assert.Equal(t, "Itens do produto recuperados com sucesso", respBody.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("erro - product_id inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/sale-items?product_id=0", nil)
		req = mux.SetURLVars(req, map[string]string{"product_id": "0"})
		w := httptest.NewRecorder()

		handler.GetByProductID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("erro - product_id não numérico", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/sale-items?product_id=abc", nil)
		req = mux.SetURLVars(req, map[string]string{"product_id": "abc"})
		w := httptest.NewRecorder()

		handler.GetByProductID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("erro - product_id negativo", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/sale-items?product_id=-1", nil)
		req = mux.SetURLVars(req, map[string]string{"product_id": "-1"})
		w := httptest.NewRecorder()

		handler.GetByProductID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	// NOVOS TESTES PARA COBRIR OS TRECHOS MENCIONADOS

	t.Run("erro - método não permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/sale-items?product_id=20", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"product_id": "20"})
		w := httptest.NewRecorder()

		handler.GetByProductID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

		// Verificar se a resposta de erro está correta
		var errorResp map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Contains(t, errorResp, "status")
		assert.Contains(t, errorResp, "message")
	})

	t.Run("erro - serviço retorna erro interno", func(t *testing.T) {
		expectedErr := fmt.Errorf("erro de conexão com o banco de dados")
		mockService.On("GetByProductID", mock.Anything, int64(20), defaultLimit, defaultOffset).Return(nil, expectedErr).Once()

		req := httptest.NewRequest(http.MethodGet, "/sale-items?product_id=20", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"product_id": "20"})

		w := httptest.NewRecorder()
		handler.GetByProductID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		// Verificar se a resposta de erro está correta
		var errorResp map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Contains(t, errorResp, "status")
		assert.Contains(t, errorResp, "message")

		mockService.AssertExpectations(t)
	})

	t.Run("erro - serviço retorna erro com product_id específico", func(t *testing.T) {
		expectedErr := fmt.Errorf("produto não encontrado")
		mockService.On("GetByProductID", mock.Anything, int64(999), defaultLimit, defaultOffset).Return(nil, expectedErr).Once()

		req := httptest.NewRequest(http.MethodGet, "/sale-items?product_id=999", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"product_id": "999"})

		w := httptest.NewRecorder()
		handler.GetByProductID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - com paginação personalizada", func(t *testing.T) {
		customLimit := 50
		customOffset := 10
		mockService.On("GetByProductID", mock.Anything, int64(20), customLimit, customOffset).Return(items, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/sale-items?product_id=20", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"product_id": "20"})

		q := req.URL.Query()
		q.Add("limit", "50")
		q.Add("offset", "10")
		req.URL.RawQuery = q.Encode()

		w := httptest.NewRecorder()
		handler.GetByProductID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - sem parâmetros de paginação (valores padrão)", func(t *testing.T) {
		// Usar os valores padrão que descobrimos no início
		mockService.On("GetByProductID", mock.Anything, int64(20), defaultLimit, defaultOffset).Return(items, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/sale-items?product_id=20", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"product_id": "20"})
		// Sem parâmetros limit e offset

		w := httptest.NewRecorder()
		handler.GetByProductID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - com diferentes métodos HTTP não-GET devem falhar", func(t *testing.T) {
		methods := []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch}

		for _, method := range methods {
			t.Run(fmt.Sprintf("método %s não permitido", method), func(t *testing.T) {
				req := httptest.NewRequest(method, "/sale-items?product_id=20", nil).WithContext(ctx)
				req = mux.SetURLVars(req, map[string]string{"product_id": "20"})
				w := httptest.NewRecorder()

				handler.GetByProductID(w, req)

				resp := w.Result()
				defer resp.Body.Close()

				assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
			})
		}
	})

	t.Run("erro - serviço retorna diferentes tipos de erro", func(t *testing.T) {
		errorCases := []struct {
			name        string
			expectedErr error
			productID   int64
		}{
			{
				name:        "erro de banco de dados",
				expectedErr: fmt.Errorf("timeout na consulta ao banco"),
				productID:   30,
			},
			{
				name:        "erro de validação",
				expectedErr: fmt.Errorf("dados inválidos"),
				productID:   40,
			},
			{
				name:        "erro genérico",
				expectedErr: fmt.Errorf("erro interno do sistema"),
				productID:   50,
			},
		}

		for _, tc := range errorCases {
			t.Run(tc.name, func(t *testing.T) {
				mockService.On("GetByProductID", mock.Anything, tc.productID, defaultLimit, defaultOffset).Return(nil, tc.expectedErr).Once()

				req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/sale-items?product_id=%d", tc.productID), nil).WithContext(ctx)
				req = mux.SetURLVars(req, map[string]string{"product_id": fmt.Sprintf("%d", tc.productID)})

				w := httptest.NewRecorder()
				handler.GetByProductID(w, req)

				resp := w.Result()
				defer resp.Body.Close()

				assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

				mockService.AssertExpectations(t)
			})
		}
	})
}
