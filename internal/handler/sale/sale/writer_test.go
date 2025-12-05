package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mocksale "github.com/WagaoCarvalho/backend_store_go/infra/mock/sale"
	dtoSale "github.com/WagaoCarvalho/backend_store_go/internal/dto/sale/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupHandler() (*mocksale.MockSale, *saleHandler) {
	mockService := new(mocksale.MockSale)
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

	t.Run("method not allowed", func(t *testing.T) {
		_, h := setupHandler()
		req := httptest.NewRequest(http.MethodGet, "/sale", nil) // GET em vez de POST
		w := httptest.NewRecorder()

		h.Create(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		assert.Contains(t, w.Body.String(), "método GET não permitido")
	})

	t.Run("erro do serviço", func(t *testing.T) {
		mockService, h := setupHandler()
		saleDTO := dtoSale.SaleDTO{UserID: utils.Int64Ptr(1), SaleDate: &now, TotalAmount: 100.0, PaymentType: "cash", Status: "active"}
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
		saleDTO := dtoSale.SaleDTO{UserID: utils.Int64Ptr(1), SaleDate: &now, TotalAmount: 100.0, PaymentType: "cash", Status: "active"}
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
			UserID:      utils.Int64Ptr(1),
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
		mockService.On("Create", mock.Anything, saleModel).Return(nil, errMsg.ErrDBInvalidForeignKey)

		h.Create(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, errMsg.ErrDBInvalidForeignKey.Error(), resp["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("erro de dados inválidos - ErrInvalidData", func(t *testing.T) {
		mockService, h := setupHandler()
		saleDTO := dtoSale.SaleDTO{
			UserID:      utils.Int64Ptr(1),
			SaleDate:    &now,
			TotalAmount: 100.0,
			PaymentType: "cash",
			Status:      "active",
		}
		body, _ := json.Marshal(saleDTO)
		req := httptest.NewRequest(http.MethodPost, "/sale", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		saleModel := dtoSale.ToSaleModel(saleDTO)

		// mock retorna erro de dados inválidos
		mockService.On("Create", mock.Anything, saleModel).Return(nil, errMsg.ErrInvalidData)

		h.Create(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, errMsg.ErrInvalidData.Error(), resp["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("erro de foreign key - ErrDBInvalidForeignKey", func(t *testing.T) {
		mockService, h := setupHandler()
		saleDTO := dtoSale.SaleDTO{
			UserID:      utils.Int64Ptr(1),
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
		mockService.On("Create", mock.Anything, saleModel).Return(nil, errMsg.ErrDBInvalidForeignKey)

		h.Create(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, errMsg.ErrDBInvalidForeignKey.Error(), resp["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("erro de duplicidade - ErrDuplicate", func(t *testing.T) {
		mockService, h := setupHandler()
		saleDTO := dtoSale.SaleDTO{
			UserID:      utils.Int64Ptr(1),
			SaleDate:    &now,
			TotalAmount: 100.0,
			PaymentType: "cash",
			Status:      "active",
		}
		body, _ := json.Marshal(saleDTO)
		req := httptest.NewRequest(http.MethodPost, "/sale", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		saleModel := dtoSale.ToSaleModel(saleDTO)

		// mock retorna erro de duplicidade
		mockService.On("Create", mock.Anything, saleModel).Return(nil, errMsg.ErrDuplicate)

		h.Create(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)

		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, errMsg.ErrDuplicate.Error(), resp["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("erro genérico do serviço (fora do switch)", func(t *testing.T) {
		mockService, h := setupHandler()
		saleDTO := dtoSale.SaleDTO{
			UserID:      utils.Int64Ptr(1),
			SaleDate:    &now,
			TotalAmount: 100.0,
			PaymentType: "cash",
			Status:      "active",
		}
		body, _ := json.Marshal(saleDTO)
		req := httptest.NewRequest(http.MethodPost, "/sale", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		saleModel := dtoSale.ToSaleModel(saleDTO)

		// mock retorna erro genérico (não capturado pelo switch)
		genericErr := errors.New("erro genérico do banco de dados")
		mockService.On("Create", mock.Anything, saleModel).Return(nil, genericErr)

		h.Create(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, genericErr.Error(), resp["message"])
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

	t.Run("method not allowed", func(t *testing.T) {
		_, h := setupHandler()
		req := httptest.NewRequest(http.MethodGet, "/sale/update/1", nil) // GET em vez de PUT
		req = mux.SetURLVars(req, map[string]string{
			"id": "1",
		})
		w := httptest.NewRecorder()

		h.Update(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		assert.Contains(t, w.Body.String(), "método GET não permitido")
	})

	t.Run("erro do serviço", func(t *testing.T) {
		mockService, h := setupHandler()

		saleDTO := dtoSale.SaleDTO{
			UserID:      utils.Int64Ptr(1),
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
			UserID:      utils.Int64Ptr(1),
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
			UserID:      utils.Int64Ptr(1),
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

	t.Run("erro de dados inválidos - ErrInvalidData", func(t *testing.T) {
		mockService, h := setupHandler()

		saleDTO := dtoSale.SaleDTO{
			UserID:      utils.Int64Ptr(1),
			SaleDate:    &now,
			TotalAmount: 100.0,
			PaymentType: "cash",
			Status:      "active",
		}

		body, _ := json.Marshal(saleDTO)
		req := httptest.NewRequest(http.MethodPut, "/sale/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		saleModel := dtoSale.ToSaleModel(saleDTO)
		saleModel.ID = 1 // Set ID from URL param
		saleDTO.ID = utils.Int64Ptr(1)

		mockService.On("Update", mock.Anything, saleModel).Return(errMsg.ErrInvalidData)

		h.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, errMsg.ErrInvalidData.Error(), resp["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("erro de ID zero - ErrZeroID", func(t *testing.T) {
		mockService, h := setupHandler()

		saleDTO := dtoSale.SaleDTO{
			UserID:      utils.Int64Ptr(1),
			SaleDate:    &now,
			TotalAmount: 100.0,
			PaymentType: "cash",
			Status:      "active",
		}

		body, _ := json.Marshal(saleDTO)
		req := httptest.NewRequest(http.MethodPut, "/sale/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		saleModel := dtoSale.ToSaleModel(saleDTO)
		saleModel.ID = 1
		saleDTO.ID = utils.Int64Ptr(1)

		mockService.On("Update", mock.Anything, saleModel).Return(errMsg.ErrZeroID)

		h.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, errMsg.ErrZeroID.Error(), resp["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("erro de versão zero - ErrZeroVersion", func(t *testing.T) {
		mockService, h := setupHandler()

		saleDTO := dtoSale.SaleDTO{
			UserID:      utils.Int64Ptr(1),
			SaleDate:    &now,
			TotalAmount: 100.0,
			PaymentType: "cash",
			Status:      "active",
		}

		body, _ := json.Marshal(saleDTO)
		req := httptest.NewRequest(http.MethodPut, "/sale/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		saleModel := dtoSale.ToSaleModel(saleDTO)
		saleModel.ID = 1
		saleDTO.ID = utils.Int64Ptr(1)

		mockService.On("Update", mock.Anything, saleModel).Return(errMsg.ErrVersionConflict)

		h.Update(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)

		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, errMsg.ErrVersionConflict.Error(), resp["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("erro de foreign key - ErrDBInvalidForeignKey", func(t *testing.T) {
		mockService, h := setupHandler()

		saleDTO := dtoSale.SaleDTO{
			UserID:      utils.Int64Ptr(1),
			SaleDate:    &now,
			TotalAmount: 100.0,
			PaymentType: "cash",
			Status:      "active",
		}

		body, _ := json.Marshal(saleDTO)
		req := httptest.NewRequest(http.MethodPut, "/sale/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		saleModel := dtoSale.ToSaleModel(saleDTO)
		saleModel.ID = 1
		saleDTO.ID = utils.Int64Ptr(1)

		mockService.On("Update", mock.Anything, saleModel).Return(errMsg.ErrDBInvalidForeignKey)

		h.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, errMsg.ErrDBInvalidForeignKey.Error(), resp["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("venda não encontrada - ErrNotFound", func(t *testing.T) {
		mockService, h := setupHandler()

		saleDTO := dtoSale.SaleDTO{
			UserID:      utils.Int64Ptr(1),
			SaleDate:    &now,
			TotalAmount: 100.0,
			PaymentType: "cash",
			Status:      "active",
		}

		body, _ := json.Marshal(saleDTO)
		req := httptest.NewRequest(http.MethodPut, "/sale/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		saleModel := dtoSale.ToSaleModel(saleDTO)
		saleModel.ID = 1
		saleDTO.ID = utils.Int64Ptr(1)

		mockService.On("Update", mock.Anything, saleModel).Return(errMsg.ErrNotFound)

		h.Update(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, errMsg.ErrNotFound.Error(), resp["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("erro genérico do serviço (fora do switch)", func(t *testing.T) {
		mockService, h := setupHandler()

		saleDTO := dtoSale.SaleDTO{
			UserID:      utils.Int64Ptr(1),
			SaleDate:    &now,
			TotalAmount: 100.0,
			PaymentType: "cash",
			Status:      "active",
		}

		body, _ := json.Marshal(saleDTO)
		req := httptest.NewRequest(http.MethodPut, "/sale/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		saleModel := dtoSale.ToSaleModel(saleDTO)
		saleModel.ID = 1
		saleDTO.ID = utils.Int64Ptr(1)

		genericErr := errors.New("erro genérico do banco de dados")
		mockService.On("Update", mock.Anything, saleModel).Return(genericErr)

		h.Update(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp map[string]interface{}
		_ = json.Unmarshal(w.Body.Bytes(), &resp)
		assert.Equal(t, genericErr.Error(), resp["message"])
		mockService.AssertExpectations(t)
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

	t.Run("method not allowed", func(t *testing.T) {
		_, h := setupHandler()
		req := httptest.NewRequest(http.MethodGet, "/sale/delete/1", nil) // GET em vez de DELETE
		req = mux.SetURLVars(req, map[string]string{
			"id": "1",
		})
		w := httptest.NewRecorder()

		h.Delete(w, req)

		assert.Equal(t, http.StatusMethodNotAllowed, w.Code)
		assert.Contains(t, w.Body.String(), "método GET não permitido")
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
