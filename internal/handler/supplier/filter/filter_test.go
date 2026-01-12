package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockSupplier "github.com/WagaoCarvalho/backend_store_go/infra/mock/supplier"
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/filter"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/supplier/supplier"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSupplierHandler_Filter_Success(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	log := logger.NewLoggerAdapter(baseLogger)

	t.Run("sucesso - retorna lista de fornecedores", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		handler := NewSupplierFilterHandler(mockService, log)

		now := time.Now()

		mockSuppliers := []*model.Supplier{
			{
				ID:        1,
				Name:      "Fornecedor A",
				CNPJ:      utils.StrToPtr("12.345.678/0001-90"),
				Status:    true,
				Version:   1,
				CreatedAt: now,
				UpdatedAt: now,
			},
			{
				ID:        2,
				Name:      "Fornecedor B",
				CPF:       utils.StrToPtr("123.456.789-00"),
				Status:    false,
				Version:   1,
				CreatedAt: now,
				UpdatedAt: now,
			},
		}

		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return(mockSuppliers, nil).
			Once()

		req := httptest.NewRequest(
			http.MethodGet,
			"/suppliers/filter?limit=10&offset=0",
			nil,
		)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.Status)

		data := resp.Data.(map[string]any)
		assert.Equal(t, float64(2), data["total"])

		items := data["items"].([]any)
		assert.Len(t, items, 2)

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - parse do status booleano", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		handler := NewSupplierFilterHandler(mockService, log)

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f any) bool {
				filter, ok := f.(*filter.SupplierFilter)
				if !ok {
					return false
				}
				return filter.Status != nil && *filter.Status == true
			})).
			Return([]*model.Supplier{}, nil).
			Once()

		req := httptest.NewRequest(
			http.MethodGet,
			"/suppliers/filter?status=true&limit=10&offset=0",
			nil,
		)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("erro - falha ao converter filtro DTO para model", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		handler := NewSupplierFilterHandler(mockService, log)

		// Força erro no ToModel: CPF e CNPJ juntos
		req := httptest.NewRequest(
			http.MethodGet,
			"/suppliers/filter?cpf=123.456.789-00&cnpj=12.345.678/0001-90",
			nil,
		)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)

		mockService.AssertNotCalled(t, "Filter", mock.Anything, mock.Anything)
	})

	t.Run("erro - filtro inválido retornado pelo serviço", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		handler := NewSupplierFilterHandler(mockService, log)

		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return(nil, errMsg.ErrInvalidFilter).
			Once()

		req := httptest.NewRequest(
			http.MethodGet,
			"/suppliers/filter?limit=10&offset=0",
			nil,
		)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("erro - falha interna no serviço", func(t *testing.T) {
		mockService := new(mockSupplier.MockSupplier)
		handler := NewSupplierFilterHandler(mockService, log)

		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return(nil, errors.New("erro interno")).
			Once()

		req := httptest.NewRequest(
			http.MethodGet,
			"/suppliers/filter?limit=10&offset=0",
			nil,
		)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)

		mockService.AssertExpectations(t)
	})

}
