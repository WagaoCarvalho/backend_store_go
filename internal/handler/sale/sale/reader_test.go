package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	dtoSale "github.com/WagaoCarvalho/backend_store_go/internal/dto/sale/sale"
)

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

	t.Run("method not allowed", func(t *testing.T) {
		_, h := setupHandler()
		req := httptest.NewRequest(http.MethodPost, "/sale/1", nil) // POST em vez de GET
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.GetByID(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		assert.Contains(t, w.Body.String(), "método POST não permitido")
	})

	t.Run("sucesso", func(t *testing.T) {
		mockService, h := setupHandler()
		saleModel := &dtoSale.SaleDTO{ID: utils.Int64Ptr(1), UserID: utils.Int64Ptr(1)}
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

	t.Run("method not allowed", func(t *testing.T) {
		_, h := setupHandler()
		req := httptest.NewRequest(http.MethodPost, "/sale/client/1", nil) // POST em vez de GET
		req = mux.SetURLVars(req, map[string]string{"client_id": "1"})
		w := httptest.NewRecorder()

		h.GetByClientID(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		assert.Contains(t, w.Body.String(), "método POST não permitido")
	})

	t.Run("sucesso", func(t *testing.T) {
		mockService, h := setupHandler()
		saleModel := []*dtoSale.SaleDTO{{ID: utils.Int64Ptr(1), UserID: utils.Int64Ptr(1)}}
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

	t.Run("method not allowed", func(t *testing.T) {
		_, h := setupHandler()
		req := httptest.NewRequest(http.MethodPost, "/sale/user/1", nil) // POST em vez de GET
		req = mux.SetURLVars(req, map[string]string{"user_id": "1"})
		w := httptest.NewRecorder()

		h.GetByUserID(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		assert.Contains(t, w.Body.String(), "método POST não permitido")
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
			{ID: utils.Int64Ptr(1), UserID: utils.Int64Ptr(1), SaleDate: &now, TotalAmount: 100.0, PaymentType: "cash", Status: "active"},
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

	t.Run("method not allowed", func(t *testing.T) {
		_, h := setupHandler()
		req := httptest.NewRequest(http.MethodPost, "/sale/date-range/2025-01-01/2025-01-31", nil) // POST em vez de GET
		req = mux.SetURLVars(req, map[string]string{
			"start": "2025-01-01",
			"end":   "2025-01-31",
		})
		w := httptest.NewRecorder()

		h.GetByDateRange(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		assert.Contains(t, w.Body.String(), "método POST não permitido")
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
			{ID: utils.Int64Ptr(1), UserID: utils.Int64Ptr(1), SaleDate: &now, TotalAmount: 100.0, PaymentType: "cash", Status: "active"},
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
