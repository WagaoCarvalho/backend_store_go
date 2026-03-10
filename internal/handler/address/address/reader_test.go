package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockAddress "github.com/WagaoCarvalho/backend_store_go/infra/mock/address"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address/address"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAddressHandler_GetByID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("sucesso", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		expectedAddress := &models.Address{
			ID:           1,
			Street:       "Rua Teste",
			StreetNumber: "123",
			City:         "São Paulo",
		}

		mockService.
			On("GetByID", mock.Anything, int64(1)).
			Return(expectedAddress, nil)

		req := newRequest("/addresses/1", "1")
		w := httptest.NewRecorder()

		h.GetByID(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var response utils.DefaultResponse
		decode(w, &response)

		assert.Equal(t, http.StatusOK, response.Status)
		assert.Equal(t, "Endereço encontrado", response.Message)
		assert.NotNil(t, response.Data)
	})

	t.Run("ID inválido retorna 400", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodGet, "/addresses/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.GetByID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ErrNotFound retorna 404", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		mockService.
			On("GetByID", mock.Anything, int64(1)).
			Return((*models.Address)(nil), errMsg.ErrNotFound)

		req := newRequest("/addresses/1", "1")
		w := httptest.NewRecorder()

		h.GetByID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)

		var response utils.DefaultResponse
		decode(w, &response)

		assert.Equal(t, http.StatusNotFound, response.Status)
	})

	t.Run("erro genérico retorna 500", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		mockService.
			On("GetByID", mock.Anything, int64(1)).
			Return((*models.Address)(nil), errors.New("db error"))

		req := newRequest("/addresses/1", "1")
		w := httptest.NewRecorder()

		h.GetByID(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})
}

/* helpers */

func newRequest(url, id string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, url, nil)
	return mux.SetURLVars(req, map[string]string{"id": id})
}

func decode[T any](w *httptest.ResponseRecorder, dest *T) {
	require.NoError(nil, json.Unmarshal(w.Body.Bytes(), dest))
}
