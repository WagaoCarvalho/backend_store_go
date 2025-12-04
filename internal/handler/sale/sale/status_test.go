package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	dtoSale "github.com/WagaoCarvalho/backend_store_go/internal/dto/sale/sale"
)

func TestSaleHandler_GetByStatus(t *testing.T) {
	now := time.Now().Format(time.RFC3339)

	t.Run("sucesso", func(t *testing.T) {
		mockService, h := setupHandler()

		// DTO de exemplo
		saleDTOs := []*dtoSale.SaleDTO{
			{ID: utils.Int64Ptr(1), UserID: utils.Int64Ptr(1), SaleDate: &now, TotalAmount: 100.0, PaymentType: "cash", Status: "active"},
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

	t.Run("method not allowed", func(t *testing.T) {
		_, h := setupHandler()
		req := httptest.NewRequest(http.MethodPost, "/sale/status/active", nil) // POST em vez de GET
		req = mux.SetURLVars(req, map[string]string{"status": "active"})
		w := httptest.NewRecorder()

		h.GetByStatus(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		assert.Contains(t, w.Body.String(), "método POST não permitido")
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

func TestSaleHandler_Cancel(t *testing.T) {
	t.Run("ID inválido", func(t *testing.T) {
		_, h := setupHandler()
		req := httptest.NewRequest(http.MethodPatch, "/sale/cancel/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.Cancel(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("method not allowed", func(t *testing.T) {
		_, h := setupHandler()
		req := httptest.NewRequest(http.MethodGet, "/sale/cancel/1", nil) // GET em vez de PATCH
		req = mux.SetURLVars(req, map[string]string{
			"id": "1",
		})
		w := httptest.NewRecorder()

		h.Cancel(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		assert.Contains(t, w.Body.String(), "método GET não permitido")
	})

	t.Run("erro do serviço", func(t *testing.T) {
		svc, h := setupHandler()
		svc.On("Cancel", mock.Anything, int64(1)).Return(errors.New("erro do serviço")).Once()

		req := httptest.NewRequest(http.MethodPatch, "/sale/cancel/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.Cancel(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		svc.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		svc, h := setupHandler()
		svc.On("Cancel", mock.Anything, int64(2)).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPatch, "/sale/cancel/2", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		w := httptest.NewRecorder()

		h.Cancel(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
		svc.AssertExpectations(t)
	})
}

func TestSaleHandler_Complete(t *testing.T) {
	t.Run("ID inválido", func(t *testing.T) {
		_, h := setupHandler()
		req := httptest.NewRequest(http.MethodPatch, "/sale/complete/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.Complete(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("method not allowed", func(t *testing.T) {
		_, h := setupHandler()
		req := httptest.NewRequest(http.MethodGet, "/sale/complete/1", nil) // GET em vez de PATCH
		req = mux.SetURLVars(req, map[string]string{
			"id": "1",
		})
		w := httptest.NewRecorder()

		h.Complete(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		assert.Contains(t, w.Body.String(), "método GET não permitido")
	})

	t.Run("erro do serviço", func(t *testing.T) {
		svc, h := setupHandler()
		svc.On("Complete", mock.Anything, int64(1)).Return(errors.New("erro do serviço")).Once()

		req := httptest.NewRequest(http.MethodPatch, "/sale/complete/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.Complete(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		svc.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		svc, h := setupHandler()
		svc.On("Complete", mock.Anything, int64(2)).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPatch, "/sale/complete/2", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		w := httptest.NewRecorder()

		h.Complete(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
		svc.AssertExpectations(t)
	})
}

func TestSaleHandler_Returned(t *testing.T) {
	t.Run("ID inválido", func(t *testing.T) {
		_, h := setupHandler()
		req := httptest.NewRequest(http.MethodPatch, "/sale/returned/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.Returned(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ID zero", func(t *testing.T) {
		_, h := setupHandler()
		req := httptest.NewRequest(http.MethodPatch, "/sale/returned/0", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "0"})
		w := httptest.NewRecorder()

		h.Returned(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("method not allowed", func(t *testing.T) {
		_, h := setupHandler()
		req := httptest.NewRequest(http.MethodGet, "/sale/returned/1", nil) // GET em vez de PATCH
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.Returned(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		assert.Contains(t, w.Body.String(), "método GET não permitido")
	})

	t.Run("erro do serviço - sale não encontrada", func(t *testing.T) {
		svc, h := setupHandler()
		svc.On("Returned", mock.Anything, int64(1)).Return(errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodPatch, "/sale/returned/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.Returned(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		svc.AssertExpectations(t)
	})

	t.Run("erro do serviço - sale não concluída", func(t *testing.T) {
		svc, h := setupHandler()
		svc.On("Returned", mock.Anything, int64(2)).Return(fmt.Errorf("%w: somente vendas concluídas podem ser devolvidas", errMsg.ErrInvalidData)).Once()

		req := httptest.NewRequest(http.MethodPatch, "/sale/returned/2", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		w := httptest.NewRecorder()

		h.Returned(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		svc.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		svc, h := setupHandler()
		svc.On("Returned", mock.Anything, int64(3)).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPatch, "/sale/returned/3", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "3"})
		w := httptest.NewRecorder()

		h.Returned(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Empty(t, w.Body.String())
		svc.AssertExpectations(t)
	})
}

func TestSaleHandler_Activate(t *testing.T) {
	t.Run("ID inválido", func(t *testing.T) {
		_, h := setupHandler()
		req := httptest.NewRequest(http.MethodPatch, "/sale/activate/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.Activate(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ID zero", func(t *testing.T) {
		_, h := setupHandler()
		req := httptest.NewRequest(http.MethodPatch, "/sale/activate/0", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "0"})
		w := httptest.NewRecorder()

		h.Activate(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("method not allowed", func(t *testing.T) {
		_, h := setupHandler()
		req := httptest.NewRequest(http.MethodPost, "/sale/activate/1", nil) // POST em vez de PATCH
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.Activate(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		assert.Contains(t, w.Body.String(), "método POST não permitido")
	})

	t.Run("erro do serviço - sale não encontrada", func(t *testing.T) {
		svc, h := setupHandler()
		svc.On("Activate", mock.Anything, int64(1)).Return(errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodPatch, "/sale/activate/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.Activate(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		svc.AssertExpectations(t)
	})

	t.Run("erro do serviço - sale não cancelada nem devolvida", func(t *testing.T) {
		svc, h := setupHandler()
		svc.On("Activate", mock.Anything, int64(2)).Return(fmt.Errorf("%w: somente vendas canceladas ou devolvidas podem ser reativadas", errMsg.ErrInvalidData)).Once()

		req := httptest.NewRequest(http.MethodPatch, "/sale/activate/2", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		w := httptest.NewRecorder()

		h.Activate(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		svc.AssertExpectations(t)
	})

	t.Run("sucesso", func(t *testing.T) {
		svc, h := setupHandler()
		svc.On("Activate", mock.Anything, int64(3)).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPatch, "/sale/activate/3", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "3"})
		w := httptest.NewRecorder()

		h.Activate(w, req)
		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Empty(t, w.Body.String())
		svc.AssertExpectations(t)
	})
}
