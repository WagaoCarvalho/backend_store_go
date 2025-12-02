package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockService "github.com/WagaoCarvalho/backend_store_go/infra/mock/sale"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/sale/item"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/item"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaleItemCreate(t *testing.T) {
	ctx := context.Background()
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	log := logger.NewLoggerAdapter(baseLogger)
	mockService := new(mockService.MockSaleItem)

	h := NewSaleItemHandler(mockService, log)

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

	dtoItem := dto.ToSaleItemDTO(item)

	body, _ := json.Marshal(dtoItem)

	t.Run("sucesso", func(t *testing.T) {
		mockService.On("Create", mock.Anything, mock.AnythingOfType("*model.SaleItem")).
			Return(item, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/sale-items", bytes.NewBuffer(body)).WithContext(ctx)
		w := httptest.NewRecorder()

		h.Create(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("erro método não permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/sale-items", bytes.NewBuffer(body)).WithContext(ctx)
		w := httptest.NewRecorder()

		h.Create(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
	})

	t.Run("erro parse JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/sale-items", bytes.NewBufferString("invalid-json")).WithContext(ctx)
		w := httptest.NewRecorder()

		h.Create(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("erro de foreign key (400)", func(t *testing.T) {
		mockService.On("Create", mock.Anything, mock.AnythingOfType("*model.SaleItem")).
			Return(nil, errMsg.ErrDBInvalidForeignKey).Once()

		req := httptest.NewRequest(http.MethodPost, "/sale-items", bytes.NewBuffer(body)).WithContext(ctx)
		w := httptest.NewRecorder()

		h.Create(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("erro interno (500)", func(t *testing.T) {
		mockService.On("Create", mock.Anything, mock.AnythingOfType("*model.SaleItem")).
			Return(nil, errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodPost, "/sale-items", bytes.NewBuffer(body)).WithContext(ctx)
		w := httptest.NewRecorder()

		h.Create(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestSaleItemHandler_Update(t *testing.T) {
	ctx := context.Background()
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	log := logger.NewLoggerAdapter(baseLogger)
	mockService := new(mockService.MockSaleItem)
	handler := NewSaleItemHandler(mockService, log)

	t.Run("sucesso - item atualizado", func(t *testing.T) {
		mockService.On("Update", mock.Anything, mock.AnythingOfType("*model.SaleItem")).Return(nil).Once()

		itemDTO := dto.SaleItemDTO{
			ID:          utils.Int64Ptr(int64(1)),
			SaleID:      10,
			ProductID:   int64(20),
			Quantity:    2,
			UnitPrice:   50.0,
			Discount:    5.0,
			Tax:         2.5,
			Subtotal:    97.5,
			Description: "Produto teste atualizado",
		}

		body, _ := json.Marshal(itemDTO)
		req := httptest.NewRequest(http.MethodPut, "/sale-items/1", bytes.NewBuffer(body)).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&respBody)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, respBody.Status)
		assert.Equal(t, "Item de venda atualizado com sucesso", respBody.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("erro - JSON inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/sale-items/1", bytes.NewBuffer([]byte("{invalid json}")))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("erro - método não permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/sale-items/1", nil) // GET em vez de PUT
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

		var errorResp map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Contains(t, errorResp["message"], "método GET não permitido")
	})

	t.Run("erro - ID inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/sale-items/0", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "0"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("erro - ID não numérico", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/sale-items/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("erro - serviço retorna ErrInvalidData", func(t *testing.T) {
		mockService.On("Update", mock.Anything, mock.AnythingOfType("*model.SaleItem")).Return(errMsg.ErrInvalidData).Once()

		itemDTO := dto.SaleItemDTO{
			SaleID:    10,
			ProductID: int64(20),
			Quantity:  2,
			UnitPrice: 50.0,
		}

		body, _ := json.Marshal(itemDTO)
		req := httptest.NewRequest(http.MethodPut, "/sale-items/1", bytes.NewBuffer(body)).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("erro - serviço retorna ErrNotFound", func(t *testing.T) {
		mockService.On("Update", mock.Anything, mock.AnythingOfType("*model.SaleItem")).Return(errMsg.ErrNotFound).Once()

		itemDTO := dto.SaleItemDTO{
			SaleID:    int64(10),
			ProductID: int64(20),
			Quantity:  2,
			UnitPrice: 50.0,
		}

		body, _ := json.Marshal(itemDTO)
		req := httptest.NewRequest(http.MethodPut, "/sale-items/999", bytes.NewBuffer(body)).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("erro - serviço retorna ErrVersionConflict", func(t *testing.T) {
		mockService.On("Update", mock.Anything, mock.AnythingOfType("*model.SaleItem")).Return(errMsg.ErrZeroVersion).Once()

		itemDTO := dto.SaleItemDTO{
			SaleID:    int64(10),
			ProductID: int64(20),
			Quantity:  2,
			UnitPrice: 50.0,
		}

		body, _ := json.Marshal(itemDTO)
		req := httptest.NewRequest(http.MethodPut, "/sale-items/1", bytes.NewBuffer(body)).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("erro - serviço retorna erro interno", func(t *testing.T) {
		expectedErr := fmt.Errorf("erro interno do banco")
		mockService.On("Update", mock.Anything, mock.AnythingOfType("*model.SaleItem")).Return(expectedErr).Once()

		itemDTO := dto.SaleItemDTO{
			SaleID:    int64(10),
			ProductID: int64(20),
			Quantity:  2,
			UnitPrice: 50.0,
		}

		body, _ := json.Marshal(itemDTO)
		req := httptest.NewRequest(http.MethodPut, "/sale-items/1", bytes.NewBuffer(body)).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("erro - diferentes métodos HTTP não permitidos", func(t *testing.T) {
		methods := []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPatch}

		for _, method := range methods {
			t.Run(fmt.Sprintf("método %s não permitido", method), func(t *testing.T) {
				req := httptest.NewRequest(method, "/sale-items/1", nil).WithContext(ctx)
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				w := httptest.NewRecorder()

				handler.Update(w, req)

				resp := w.Result()
				defer resp.Body.Close()

				assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
			})
		}
	})

	t.Run("sucesso - campos parciais no DTO", func(t *testing.T) {
		mockService.On("Update", mock.Anything, mock.AnythingOfType("*model.SaleItem")).Return(nil).Once()

		// DTO com apenas alguns campos preenchidos
		itemDTO := dto.SaleItemDTO{
			Quantity: 5, // Apenas quantidade alterada
			Discount: 10.0,
		}

		body, _ := json.Marshal(itemDTO)
		req := httptest.NewRequest(http.MethodPut, "/sale-items/1", bytes.NewBuffer(body)).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestSaleItemHandler_Delete(t *testing.T) {
	ctx := context.Background()
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	log := logger.NewLoggerAdapter(baseLogger)
	mockService := new(mockService.MockSaleItem)
	handler := NewSaleItemHandler(mockService, log)

	t.Run("sucesso - item deletado", func(t *testing.T) {
		mockService.On("Delete", mock.Anything, int64(1)).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/sale-items/1", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		assert.Equal(t, 0, w.Body.Len()) // No Content não deve ter corpo

		mockService.AssertExpectations(t)
	})

	t.Run("erro - método não permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/sale-items/1", nil) // GET em vez de DELETE
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

		var errorResp map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Contains(t, errorResp["message"], "método GET não permitido")
	})

	t.Run("erro - ID inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/sale-items/0", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "0"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("erro - ID não numérico", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/sale-items/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("erro - ID negativo", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/sale-items/-1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "-1"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("erro - serviço retorna erro", func(t *testing.T) {
		expectedErr := fmt.Errorf("erro ao deletar item")
		mockService.On("Delete", mock.Anything, int64(999)).Return(expectedErr).Once()

		req := httptest.NewRequest(http.MethodDelete, "/sale-items/999", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

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

	t.Run("erro - diferentes métodos HTTP não permitidos", func(t *testing.T) {
		methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch}

		for _, method := range methods {
			t.Run(fmt.Sprintf("método %s não permitido", method), func(t *testing.T) {
				req := httptest.NewRequest(method, "/sale-items/1", nil).WithContext(ctx)
				req = mux.SetURLVars(req, map[string]string{"id": "1"})
				w := httptest.NewRecorder()

				handler.Delete(w, req)

				resp := w.Result()
				defer resp.Body.Close()

				assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
			})
		}
	})

	t.Run("sucesso - IDs diferentes", func(t *testing.T) {
		testCases := []struct {
			id   string
			want int64
		}{
			{"1", 1},
			{"100", 100},
			{"9999", 9999},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("ID %s", tc.id), func(t *testing.T) {
				mockService.On("Delete", mock.Anything, tc.want).Return(nil).Once()

				req := httptest.NewRequest(http.MethodDelete, "/sale-items/"+tc.id, nil).WithContext(ctx)
				req = mux.SetURLVars(req, map[string]string{"id": tc.id})
				w := httptest.NewRecorder()

				handler.Delete(w, req)

				resp := w.Result()
				defer resp.Body.Close()

				assert.Equal(t, http.StatusNoContent, resp.StatusCode)
				assert.Equal(t, 0, w.Body.Len())

				mockService.AssertExpectations(t)
			})
		}
	})
}

func TestSaleItemHandler_DeleteBySaleID(t *testing.T) {
	ctx := context.Background()
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	log := logger.NewLoggerAdapter(baseLogger)
	mockService := new(mockService.MockSaleItem)
	handler := NewSaleItemHandler(mockService, log)

	t.Run("sucesso - itens da venda deletados", func(t *testing.T) {
		mockService.On("DeleteBySaleID", mock.Anything, int64(10)).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/sale-items/sale/10", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"sale_id": "10"})
		w := httptest.NewRecorder()

		handler.DeleteBySaleID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		assert.Equal(t, 0, w.Body.Len()) // No Content não deve ter corpo

		mockService.AssertExpectations(t)
	})

	t.Run("erro - método não permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/sale-items/sale/10", nil) // GET em vez de DELETE
		req = mux.SetURLVars(req, map[string]string{"sale_id": "10"})
		w := httptest.NewRecorder()

		handler.DeleteBySaleID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)

		var errorResp map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&errorResp)
		assert.NoError(t, err)
		assert.Contains(t, errorResp["message"], "método GET não permitido")
	})

	t.Run("erro - sale_id inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/sale-items/sale/0", nil)
		req = mux.SetURLVars(req, map[string]string{"sale_id": "0"})
		w := httptest.NewRecorder()

		handler.DeleteBySaleID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("erro - sale_id não numérico", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/sale-items/sale/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"sale_id": "abc"})
		w := httptest.NewRecorder()

		handler.DeleteBySaleID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("erro - sale_id negativo", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/sale-items/sale/-1", nil)
		req = mux.SetURLVars(req, map[string]string{"sale_id": "-1"})
		w := httptest.NewRecorder()

		handler.DeleteBySaleID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("erro - serviço retorna erro", func(t *testing.T) {
		expectedErr := fmt.Errorf("erro ao deletar itens da venda")
		mockService.On("DeleteBySaleID", mock.Anything, int64(999)).Return(expectedErr).Once()

		req := httptest.NewRequest(http.MethodDelete, "/sale-items/sale/999", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"sale_id": "999"})
		w := httptest.NewRecorder()

		handler.DeleteBySaleID(w, req)

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

	t.Run("erro - diferentes métodos HTTP não permitidos", func(t *testing.T) {
		methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch}

		for _, method := range methods {
			t.Run(fmt.Sprintf("método %s não permitido", method), func(t *testing.T) {
				req := httptest.NewRequest(method, "/sale-items/sale/10", nil).WithContext(ctx)
				req = mux.SetURLVars(req, map[string]string{"sale_id": "10"})
				w := httptest.NewRecorder()

				handler.DeleteBySaleID(w, req)

				resp := w.Result()
				defer resp.Body.Close()

				assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode)
			})
		}
	})

	t.Run("sucesso - diferentes sale_ids", func(t *testing.T) {
		testCases := []struct {
			saleID string
			want   int64
		}{
			{"1", 1},
			{"100", 100},
			{"9999", 9999},
		}

		for _, tc := range testCases {
			t.Run(fmt.Sprintf("sale_id %s", tc.saleID), func(t *testing.T) {
				mockService.On("DeleteBySaleID", mock.Anything, tc.want).Return(nil).Once()

				req := httptest.NewRequest(http.MethodDelete, "/sale-items/sale/"+tc.saleID, nil).WithContext(ctx)
				req = mux.SetURLVars(req, map[string]string{"sale_id": tc.saleID})
				w := httptest.NewRecorder()

				handler.DeleteBySaleID(w, req)

				resp := w.Result()
				defer resp.Body.Close()

				assert.Equal(t, http.StatusNoContent, resp.StatusCode)
				assert.Equal(t, 0, w.Body.Len())

				mockService.AssertExpectations(t)
			})
		}
	})

	t.Run("erro - serviço retorna diferentes tipos de erro", func(t *testing.T) {
		errorCases := []struct {
			name        string
			expectedErr error
			saleID      int64
		}{
			{
				name:        "erro de banco de dados",
				expectedErr: fmt.Errorf("timeout na conexão com o banco"),
				saleID:      50,
			},
			{
				name:        "erro de constraint",
				expectedErr: fmt.Errorf("violação de chave estrangeira"),
				saleID:      60,
			},
			{
				name:        "erro genérico",
				expectedErr: fmt.Errorf("erro interno do sistema"),
				saleID:      70,
			},
		}

		for _, tc := range errorCases {
			t.Run(tc.name, func(t *testing.T) {
				mockService.On("DeleteBySaleID", mock.Anything, tc.saleID).Return(tc.expectedErr).Once()

				req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/sale-items/sale/%d", tc.saleID), nil).WithContext(ctx)
				req = mux.SetURLVars(req, map[string]string{"sale_id": fmt.Sprintf("%d", tc.saleID)})
				w := httptest.NewRecorder()

				handler.DeleteBySaleID(w, req)

				resp := w.Result()
				defer resp.Body.Close()

				assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

				mockService.AssertExpectations(t)
			})
		}
	})
}
