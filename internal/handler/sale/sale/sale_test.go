package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mocksale "github.com/WagaoCarvalho/backend_store_go/infra/mock/service/sale"
	dtoSale "github.com/WagaoCarvalho/backend_store_go/internal/dto/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupHandler() (*mocksale.MockSaleService, *SaleHandler) {
	mockService := new(mocksale.MockSaleService)
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	loggerAdapter := logger.NewLoggerAdapter(baseLogger)
	h := NewSaleHandler(mockService, loggerAdapter)
	return mockService, h
}

func TestSaleHandler_Create(t *testing.T) {
	now := time.Now().Format(time.RFC3339)
	t.Run("erro JSON inválido", func(t *testing.T) {
		_, h := setupHandler()
		req := httptest.NewRequest(http.MethodPost, "/sale", bytes.NewBuffer([]byte("{invalid json}")))
		w := httptest.NewRecorder()

		h.Create(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("erro do serviço", func(t *testing.T) {
		mockService, h := setupHandler()
		saleDTO := dtoSale.SaleDTO{UserID: 1, SaleDate: &now, TotalAmount: 100.0, PaymentType: "cash", Status: "active"}
		body, _ := json.Marshal(saleDTO)
		req := httptest.NewRequest(http.MethodPost, "/sale", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		mockService.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("erro serviço"))

		h.Create(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockService, h := setupHandler()
		saleDTO := dtoSale.SaleDTO{UserID: 1, SaleDate: &now, TotalAmount: 100.0, PaymentType: "cash", Status: "active"}
		body, _ := json.Marshal(saleDTO)
		req := httptest.NewRequest(http.MethodPost, "/sale", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		saleModel := dtoSale.ToSaleModel(saleDTO)
		mockService.On("Create", mock.Anything, saleModel).Return(saleModel, nil)

		h.Create(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "Venda criada com sucesso", resp["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("erro de foreign key", func(t *testing.T) {
		mockService, h := setupHandler()
		saleDTO := dtoSale.SaleDTO{
			UserID:      1,
			SaleDate:    &now,
			TotalAmount: 100.0,
			PaymentType: "cash",
			Status:      "active",
		}
		body, _ := json.Marshal(saleDTO)
		req := httptest.NewRequest(http.MethodPost, "/sale", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		saleModel := dtoSale.ToSaleModel(saleDTO)

		// mock retorna erro de foreign key
		mockService.On("Create", mock.Anything, saleModel).Return(nil, errMsg.ErrInvalidForeignKey)

		h.Create(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, errMsg.ErrInvalidForeignKey.Error(), resp["message"])
		mockService.AssertExpectations(t)
	})

}

func TestSaleHandler_GetByID(t *testing.T) {
	t.Run("id inválido", func(t *testing.T) {
		_, h := setupHandler()
		req := httptest.NewRequest(http.MethodGet, "/sale?id=abc", nil)
		w := httptest.NewRecorder()

		h.GetByID(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("erro do serviço", func(t *testing.T) {
		mockService, h := setupHandler()
		req := httptest.NewRequest(http.MethodGet, "/sale?id=1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("GetByID", mock.Anything, int64(1)).Return(nil, errors.New("erro serviço"))

		h.GetByID(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockService, h := setupHandler()
		saleModel := &dtoSale.SaleDTO{ID: utils.Int64Ptr(1), UserID: 1}
		req := httptest.NewRequest(http.MethodGet, "/sale?id=1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("GetByID", mock.Anything, int64(1)).Return(dtoSale.ToSaleModel(*saleModel), nil)

		h.GetByID(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestSaleHandler_GetByClientID(t *testing.T) {
	t.Run("clientID inválido", func(t *testing.T) {
		_, h := setupHandler()
		req := httptest.NewRequest(http.MethodGet, "/sale/client?client_id=0", nil)
		w := httptest.NewRecorder()

		h.GetByClientID(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("erro do serviço", func(t *testing.T) {
		mockService, h := setupHandler()
		req := httptest.NewRequest(http.MethodGet, "/sale/client/1", nil)
		req = mux.SetURLVars(req, map[string]string{"client_id": "1"})
		w := httptest.NewRecorder()

		mockService.On("GetByClientID", mock.Anything, int64(1), mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("erro"))

		h.GetByClientID(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockService, h := setupHandler()
		saleModel := []*dtoSale.SaleDTO{{ID: utils.Int64Ptr(1), UserID: 1}}
		req := httptest.NewRequest(http.MethodGet, "/sale/client/1", nil)
		req = mux.SetURLVars(req, map[string]string{"client_id": "1"})
		w := httptest.NewRecorder()

		mockService.On("GetByClientID", mock.Anything, int64(1), mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(dtoSale.SaleDTOListToModelList(saleModel), nil)

		h.GetByClientID(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})
}

func TestSaleHandler_GetByUserID(t *testing.T) {
	now := time.Now().Format(time.RFC3339)
	t.Run("userID inválido", func(t *testing.T) {
		_, h := setupHandler()
		req := httptest.NewRequest(http.MethodGet, "/sale/user/0", nil)
		w := httptest.NewRecorder()

		h.GetByUserID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("erro do serviço", func(t *testing.T) {
		mockService, h := setupHandler()
		req := httptest.NewRequest(http.MethodGet, "/sale/user/1", nil)
		req = mux.SetURLVars(req, map[string]string{"user_id": "1"})
		w := httptest.NewRecorder()

		mockService.On("GetByUserID", mock.Anything, int64(1), mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("erro serviço"))

		h.GetByUserID(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockService, h := setupHandler()
		saleModel := []*dtoSale.SaleDTO{
			{ID: utils.Int64Ptr(1), UserID: 1, SaleDate: &now, TotalAmount: 100.0, PaymentType: "cash", Status: "active"},
		}

		req := httptest.NewRequest(http.MethodGet, "/sale/user/1", nil)
		req = mux.SetURLVars(req, map[string]string{"user_id": "1"})
		w := httptest.NewRecorder()

		// Converte DTOs para Models para o mock
		saleModels := dtoSale.SaleDTOListToModelList(saleModel)
		mockService.On("GetByUserID", mock.Anything, int64(1), mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(saleModels, nil)

		h.GetByUserID(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "Vendas do usuário recuperadas", resp["message"])
		mockService.AssertExpectations(t)
	})
}

func TestSaleHandler_GetByStatus(t *testing.T) {
	now := time.Now().Format(time.RFC3339)

	t.Run("sucesso", func(t *testing.T) {
		mockService, h := setupHandler()

		// DTO de exemplo
		saleDTOs := []*dtoSale.SaleDTO{
			{ID: utils.Int64Ptr(1), UserID: 1, SaleDate: &now, TotalAmount: 100.0, PaymentType: "cash", Status: "active"},
		}

		// Path param
		req := httptest.NewRequest(http.MethodGet, "/sale/status/active", nil)
		req = mux.SetURLVars(req, map[string]string{"status": "active"})
		w := httptest.NewRecorder()

		// Converte DTO para Model
		saleModels := dtoSale.SaleDTOListToModelList(saleDTOs)
		mockService.On("GetByStatus", mock.Anything, "active", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(saleModels, nil)

		h.GetByStatus(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "Vendas por status recuperadas", resp["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("erro do serviço", func(t *testing.T) {
		mockService, h := setupHandler()
		req := httptest.NewRequest(http.MethodGet, "/sale/status/active", nil)
		req = mux.SetURLVars(req, map[string]string{"status": "active"})
		w := httptest.NewRecorder()

		mockService.On("GetByStatus", mock.Anything, "active", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New("erro serviço"))

		h.GetByStatus(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("status vazio", func(t *testing.T) {
		_, h := setupHandler()
		req := httptest.NewRequest(http.MethodGet, "/sale/status/", nil)
		req = mux.SetURLVars(req, map[string]string{"status": ""})
		w := httptest.NewRecorder()

		h.GetByStatus(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestSaleHandler_GetByDateRange(t *testing.T) {
	now := time.Now().Format(time.RFC3339)
	start := now
	end := now

	t.Run("datas inválidas", func(t *testing.T) {
		_, h := setupHandler()
		req := httptest.NewRequest(http.MethodGet, "/sale/daterange/invalid/invalid", nil)
		req = mux.SetURLVars(req, map[string]string{"start": "invalid", "end": "invalid"})
		w := httptest.NewRecorder()

		h.GetByDateRange(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("erro do serviço", func(t *testing.T) {
		mockService, h := setupHandler()
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/sale/daterange/%s/%s", start, end), nil)
		req = mux.SetURLVars(req, map[string]string{"start": start, "end": end})
		w := httptest.NewRecorder()

		mockService.On("GetByDateRange",
			mock.Anything, // ctx
			mock.Anything, // start time
			mock.Anything, // end time
			mock.Anything, // limit
			mock.Anything, // offset
			mock.Anything, // orderBy
			mock.Anything, // orderDir
		).Return(nil, errors.New("erro serviço"))

		h.GetByDateRange(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("datas ausentes ou inválidas", func(t *testing.T) {
		_, h := setupHandler()

		// start vazio
		req1 := httptest.NewRequest(http.MethodGet, "/sale/daterange//2025-09-12T00:00:00Z", nil)
		req1 = mux.SetURLVars(req1, map[string]string{"start": "", "end": "2025-09-12T00:00:00Z"})
		w1 := httptest.NewRecorder()
		h.GetByDateRange(w1, req1)
		assert.Equal(t, http.StatusBadRequest, w1.Code)

		// end vazio
		req2 := httptest.NewRequest(http.MethodGet, "/sale/daterange/2025-09-11T00:00:00Z/", nil)
		req2 = mux.SetURLVars(req2, map[string]string{"start": "2025-09-11T00:00:00Z", "end": ""})
		w2 := httptest.NewRecorder()
		h.GetByDateRange(w2, req2)
		assert.Equal(t, http.StatusBadRequest, w2.Code)

		// datas inválidas
		req3 := httptest.NewRequest(http.MethodGet, "/sale/daterange/invalid/invalid", nil)
		req3 = mux.SetURLVars(req3, map[string]string{"start": "invalid", "end": "invalid"})
		w3 := httptest.NewRecorder()
		h.GetByDateRange(w3, req3)
		assert.Equal(t, http.StatusBadRequest, w3.Code)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockService, h := setupHandler()

		salesDTO := []*dtoSale.SaleDTO{
			{ID: utils.Int64Ptr(1), UserID: 1, SaleDate: &now, TotalAmount: 100.0, PaymentType: "cash", Status: "active"},
		}

		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/sale/daterange/%s/%s", start, end), nil)
		req = mux.SetURLVars(req, map[string]string{"start": start, "end": end})
		w := httptest.NewRecorder()

		saleModels := dtoSale.SaleDTOListToModelList(salesDTO)
		mockService.On("GetByDateRange",
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
			mock.Anything,
		).Return(saleModels, nil)

		h.GetByDateRange(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, "Vendas por período recuperadas", resp["message"])
		mockService.AssertExpectations(t)
	})
}

func TestSaleHandler_Update(t *testing.T) {
	now := time.Now().Format(time.RFC3339)

	t.Run("erro JSON inválido", func(t *testing.T) {
		_, h := setupHandler()
		req := httptest.NewRequest(http.MethodPut, "/sale/update/1", bytes.NewBuffer([]byte("{invalid json}")))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("erro do serviço", func(t *testing.T) {
		mockService, h := setupHandler()

		saleDTO := dtoSale.SaleDTO{
			UserID:      1,
			SaleDate:    &now,
			TotalAmount: 100.0,
			PaymentType: "cash",
			Status:      "active",
		}

		body, _ := json.Marshal(saleDTO)
		req := httptest.NewRequest(http.MethodPut, "/sale/update/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("Update", mock.Anything, mock.Anything).Return(errors.New("erro serviço"))

		h.Update(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockService, h := setupHandler()

		saleDTO := dtoSale.SaleDTO{
			UserID:      1,
			SaleDate:    &now,
			TotalAmount: 100.0,
			PaymentType: "cash",
			Status:      "active",
		}

		body, _ := json.Marshal(saleDTO)
		req := httptest.NewRequest(http.MethodPut, "/sale/update/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("Update", mock.Anything, mock.Anything).Return(nil)

		h.Update(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("ID inválido ou menor ou igual a zero", func(t *testing.T) {
		_, h := setupHandler()

		now := time.Now().Format(time.RFC3339)
		validSale := dtoSale.SaleDTO{
			UserID:      1,
			SaleDate:    &now,
			TotalAmount: 100.0,
			PaymentType: "cash",
			Status:      "active",
		}
		body, _ := json.Marshal(validSale)

		// ID não numérico
		req1 := httptest.NewRequest(http.MethodPut, "/sale/update/abc", bytes.NewBuffer(body))
		req1 = mux.SetURLVars(req1, map[string]string{"id": "abc"})
		w1 := httptest.NewRecorder()
		h.Update(w1, req1)
		assert.Equal(t, http.StatusBadRequest, w1.Code)

		// ID igual a zero
		req2 := httptest.NewRequest(http.MethodPut, "/sale/update/0", bytes.NewBuffer(body))
		req2 = mux.SetURLVars(req2, map[string]string{"id": "0"})
		w2 := httptest.NewRecorder()
		h.Update(w2, req2)
		assert.Equal(t, http.StatusBadRequest, w2.Code)

		// ID negativo
		req3 := httptest.NewRequest(http.MethodPut, "/sale/update/-5", bytes.NewBuffer(body))
		req3 = mux.SetURLVars(req3, map[string]string{"id": "-5"})
		w3 := httptest.NewRecorder()
		h.Update(w3, req3)
		assert.Equal(t, http.StatusBadRequest, w3.Code)
	})

}

func TestSaleHandler_Delete(t *testing.T) {
	t.Run("ID inválido", func(t *testing.T) {
		_, h := setupHandler()
		req := httptest.NewRequest(http.MethodDelete, "/sale/delete/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.Delete(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("erro do serviço", func(t *testing.T) {
		mockService, h := setupHandler()
		req := httptest.NewRequest(http.MethodDelete, "/sale/delete/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("Delete", mock.Anything, int64(1)).Return(errors.New("erro serviço"))

		h.Delete(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		mockService, h := setupHandler()
		req := httptest.NewRequest(http.MethodDelete, "/sale/delete/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("Delete", mock.Anything, int64(1)).Return(nil)

		h.Delete(w, req)
		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		assert.Empty(t, w.Body.String())

		mockService.AssertExpectations(t)
	})
}
