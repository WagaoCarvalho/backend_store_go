package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockService "github.com/WagaoCarvalho/backend_store_go/infra/mock/sale"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/sale/item"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSaleItemHandler_ItemExists(t *testing.T) {
	ctx := context.Background()
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	log := logger.NewLoggerAdapter(baseLogger)
	mockService := new(mockService.MockSaleItem)
	h := NewSaleItemHandler(mockService, log)

	t.Run("sucesso - item existe", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("ItemExists", mock.Anything, int64(1)).
			Return(true, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/sale-items/1", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		h.ItemExists(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp struct {
			Status  int                       `json:"status"`
			Message string                    `json:"message"`
			Data    dto.ItemExistsResponseDTO `json:"data"`
		}
		err := json.NewDecoder(rec.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.True(t, resp.Data.Exists)
		mockService.AssertExpectations(t)
	})

	t.Run("erro - método não permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/sale-items/1", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		h.ItemExists(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("erro - id inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/sale-items/0", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"id": "0"})
		rec := httptest.NewRecorder()

		h.ItemExists(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("erro - foreign key violation", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("ItemExists", mock.Anything, int64(2)).
			Return(false, errMsg.ErrDBInvalidForeignKey).Once()

		req := httptest.NewRequest(http.MethodGet, "/sale-items/2", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		rec := httptest.NewRecorder()

		h.ItemExists(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("erro - interno no service", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		mockService.On("ItemExists", mock.Anything, int64(3)).
			Return(false, errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodGet, "/sale-items/3", nil).WithContext(ctx)
		req = mux.SetURLVars(req, map[string]string{"id": "3"})
		rec := httptest.NewRecorder()

		h.ItemExists(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}
