package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockSale "github.com/WagaoCarvalho/backend_store_go/infra/mock/sale"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/filter"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/sale/sale"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaleHandler_Filter(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("erro - falha no serviço", func(t *testing.T) {
		mockService := new(mockSale.MockSale)
		handler := NewSaleFilterHandler(mockService, logger)

		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return(nil, errors.New("db error"))

		req := httptest.NewRequest(http.MethodGet, "/sales/filter?limit=10&offset=0", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("erro - filtro inválido", func(t *testing.T) {
		mockService := new(mockSale.MockSale)
		handler := NewSaleFilterHandler(mockService, logger)

		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return(nil, errMsg.ErrInvalidFilter)

		req := httptest.NewRequest(http.MethodGet, "/sales/filter?limit=-1", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - retorna lista de vendas", func(t *testing.T) {
		mockService := new(mockSale.MockSale)
		handler := NewSaleFilterHandler(mockService, logger)

		now := time.Now()

		mockSales := []*model.Sale{
			{
				ID:          1,
				PaymentType: "credit",
				Status:      "paid",
				TotalAmount: 100,
				SaleDate:    now,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				ID:          2,
				PaymentType: "debit",
				Status:      "pending",
				TotalAmount: 50,
				SaleDate:    now,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}

		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return(mockSales, nil)

		req := httptest.NewRequest(http.MethodGet, "/sales/filter?limit=2&offset=0", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Status)
		assert.Equal(t, "Vendas listadas com sucesso", resp.Message)

		data := resp.Data.(map[string]any)
		assert.Equal(t, float64(2), data["total"])
		assert.NotEmpty(t, data["items"])

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - retorna lista vazia", func(t *testing.T) {
		mockService := new(mockSale.MockSale)
		handler := NewSaleFilterHandler(mockService, logger)

		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return([]*model.Sale{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/sales/filter?limit=5&offset=0", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Status)

		data := resp.Data.(map[string]any)
		assert.Equal(t, float64(0), data["total"])

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro por payment_type", func(t *testing.T) {
		mockService := new(mockSale.MockSale)
		handler := NewSaleFilterHandler(mockService, logger)

		mockSales := []*model.Sale{
			{
				ID:          1,
				PaymentType: "credit",
				Status:      "paid",
				TotalAmount: 200,
			},
		}

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f *filter.SaleFilter) bool {
				return f.PaymentType == "credit"
			})).
			Return(mockSales, nil)

		req := httptest.NewRequest(http.MethodGet, "/sales/filter?payment_type=credit", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		data := resp.Data.(map[string]any)
		assert.Equal(t, float64(1), data["total"])
		items := data["items"].([]any)
		assert.Equal(t, "credit", items[0].(map[string]any)["payment_type"])

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro por client_id", func(t *testing.T) {
		mockService := new(mockSale.MockSale)
		handler := NewSaleFilterHandler(mockService, logger)

		mockSales := []*model.Sale{
			{
				ID:          1,
				ClientID:    utils.Int64Ptr(int64(100)),
				PaymentType: "credit",
				Status:      "paid",
				TotalAmount: 150.0,
			},
		}

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f *filter.SaleFilter) bool {
				return f.ClientID != nil && *f.ClientID == 100
			})).
			Return(mockSales, nil)

		req := httptest.NewRequest(http.MethodGet, "/sales/filter?client_id=100", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		data := resp.Data.(map[string]any)
		assert.Equal(t, float64(1), data["total"])

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro por user_id", func(t *testing.T) {
		mockService := new(mockSale.MockSale)
		handler := NewSaleFilterHandler(mockService, logger)

		mockSales := []*model.Sale{
			{
				ID:          1,
				UserID:      utils.Int64Ptr(int64(50)),
				PaymentType: "pix",
				Status:      "completed",
				TotalAmount: 75.0,
			},
		}

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f *filter.SaleFilter) bool {
				return f.UserID != nil && *f.UserID == 50
			})).
			Return(mockSales, nil)

		req := httptest.NewRequest(http.MethodGet, "/sales/filter?user_id=50", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		data := resp.Data.(map[string]any)
		assert.Equal(t, float64(1), data["total"])

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro por status", func(t *testing.T) {
		mockService := new(mockSale.MockSale)
		handler := NewSaleFilterHandler(mockService, logger)

		mockSales := []*model.Sale{
			{
				ID:          1,
				PaymentType: "cash",
				Status:      "pending",
				TotalAmount: 200.0,
			},
		}

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f *filter.SaleFilter) bool {
				return f.Status == "pending"
			})).
			Return(mockSales, nil)

		req := httptest.NewRequest(http.MethodGet, "/sales/filter?status=pending", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		data := resp.Data.(map[string]any)
		assert.Equal(t, float64(1), data["total"])
		items := data["items"].([]any)
		assert.Equal(t, "pending", items[0].(map[string]any)["status"])

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - múltiplos filtros combinados", func(t *testing.T) {
		mockService := new(mockSale.MockSale)
		handler := NewSaleFilterHandler(mockService, logger)

		mockSales := []*model.Sale{
			{
				ID:          1,
				ClientID:    utils.Int64Ptr(int64(100)),
				UserID:      utils.Int64Ptr(int64(50)),
				PaymentType: "credit",
				Status:      "completed",
				TotalAmount: 300.0,
			},
		}

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f *filter.SaleFilter) bool {
				return f.ClientID != nil && *f.ClientID == 100 &&
					f.UserID != nil && *f.UserID == 50 &&
					f.PaymentType == "credit" &&
					f.Status == "completed"
			})).
			Return(mockSales, nil)

		req := httptest.NewRequest(http.MethodGet,
			"/sales/filter?client_id=100&user_id=50&payment_type=credit&status=completed",
			nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		data := resp.Data.(map[string]any)
		assert.Equal(t, float64(1), data["total"])

		mockService.AssertExpectations(t)
	})

	t.Run("erro - ToModel retorna erro de validação para status inválido", func(t *testing.T) {
		mockService := new(mockSale.MockSale)
		handler := NewSaleFilterHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodGet, "/sales/filter?status=INVALID_STATUS_VALUE", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code, "Deve retornar 400 quando status é inválido")

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "filtro inválido", "Mensagem deve conter 'filtro inválido'")
		assert.Contains(t, resp.Message, "status", "Mensagem deve mencionar o campo 'status'")
		assert.Contains(t, resp.Message, "INVALID_STATUS_VALUE", "Mensagem deve conter o valor inválido")

		mockService.AssertNotCalled(t, "Filter", mock.Anything, mock.Anything,
			"Serviço não deve ser chamado quando ToModel retorna erro")
	})

}
