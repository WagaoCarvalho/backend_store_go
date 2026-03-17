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
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	setup := func() (*mockSale.MockSale, *saleFilterHandler) {
		mockService := new(mockSale.MockSale)
		handler := NewSaleFilterHandler(mockService, logAdapter)
		return mockService, handler
	}

	// TESTES DE VALIDAÇÃO DE PARÂMETROS DESCONHECIDOS
	t.Run("erro - parâmetro desconhecido na query", func(t *testing.T) {
		mockService, handler := setup()
		req := httptest.NewRequest(http.MethodGet, "/sales/filter?parametro_invalido=valor", nil)
		rec := httptest.NewRecorder()
		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "parâmetro de consulta inválido")

		mockService.AssertNotCalled(t, "Filter")
	})

	// TESTES DE VALIDAÇÃO DE VALORES INVÁLIDOS
	t.Run("erro - client_id inválido (não numérico)", func(t *testing.T) {
		mockService, handler := setup()
		req := httptest.NewRequest(http.MethodGet, "/sales/filter?client_id=abc", nil)
		rec := httptest.NewRecorder()
		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "client_id deve ser um número inteiro")

		mockService.AssertNotCalled(t, "Filter")
	})

	t.Run("erro - user_id inválido (não numérico)", func(t *testing.T) {
		mockService, handler := setup()
		req := httptest.NewRequest(http.MethodGet, "/sales/filter?user_id=xyz", nil)
		rec := httptest.NewRecorder()
		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "user_id deve ser um número inteiro")

		mockService.AssertNotCalled(t, "Filter")
	})

	// TESTES DE ERRO DO SERVIÇO
	t.Run("erro - falha no serviço (erro genérico)", func(t *testing.T) {
		mockService, handler := setup()
		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return(nil, errors.New("db error")).
			Once()

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

	t.Run("erro - filtro inválido retornado pelo serviço", func(t *testing.T) {
		mockService, handler := setup()
		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return(nil, errMsg.ErrInvalidFilter).
			Once()

		req := httptest.NewRequest(http.MethodGet, "/sales/filter?limit=10&offset=0", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)

		mockService.AssertExpectations(t)
	})

	// TESTES DE SUCESSO
	t.Run("sucesso - retorna lista de vendas", func(t *testing.T) {
		mockService, handler := setup()

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
			Return(mockSales, nil).
			Once()

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
		mockService, handler := setup()
		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return([]*model.Sale{}, nil).
			Once()

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
		mockService, handler := setup()
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
			Return(mockSales, nil).
			Once()

		req := httptest.NewRequest(http.MethodGet, "/sales/filter?payment_type=credit", nil)
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

	t.Run("sucesso - filtro por client_id válido", func(t *testing.T) {
		mockService, handler := setup()
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
			Return(mockSales, nil).
			Once()

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

	t.Run("sucesso - filtro por user_id válido", func(t *testing.T) {
		mockService, handler := setup()
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
			Return(mockSales, nil).
			Once()

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
		mockService, handler := setup()
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
			Return(mockSales, nil).
			Once()

		req := httptest.NewRequest(http.MethodGet, "/sales/filter?status=pending", nil)
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

	t.Run("sucesso - múltiplos filtros combinados", func(t *testing.T) {
		mockService, handler := setup()
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
			Return(mockSales, nil).
			Once()

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

	t.Run("sucesso - paginação com valores padrão quando não informados", func(t *testing.T) {
		mockService, handler := setup()
		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f *filter.SaleFilter) bool {
				return f.Limit == 10 && f.Offset == 0
			})).
			Return([]*model.Sale{}, nil).
			Once()

		req := httptest.NewRequest(http.MethodGet, "/sales/filter", nil)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})
}
