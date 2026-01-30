package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockAddress "github.com/WagaoCarvalho/backend_store_go/infra/mock/address"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
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

func TestAddressHandler_GetByUserID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("sucesso", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		expectedAddresses := []*models.Address{
			{ID: 1, Street: "Rua A", UserID: utils.Int64Ptr(1)},
			{ID: 2, Street: "Rua B", UserID: utils.Int64Ptr(1)},
		}

		mockService.
			On("GetByUserID", mock.Anything, int64(1)).
			Return(expectedAddresses, nil)

		w := exec(h.GetByUserID, "/addresses/user/1", "1")

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var response utils.DefaultResponse
		decode(w, &response)

		assert.Equal(t, http.StatusOK, response.Status)
		assert.Equal(t, "Endereços do usuário encontrados", response.Message)
		assert.NotNil(t, response.Data)

		addresses, ok := response.Data.([]interface{})
		assert.True(t, ok)
		assert.Len(t, addresses, 2)
	})

	t.Run("ID inválido retorna 400", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodGet, "/addresses/user/invalid", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
		w := httptest.NewRecorder()

		h.GetByUserID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ErrNotFound retorna 404", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		mockService.
			On("GetByUserID", mock.Anything, int64(1)).
			Return(nil, errMsg.ErrNotFound)

		w := exec(h.GetByUserID, "/addresses/user/1", "1")

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})

	t.Run("erro genérico retorna 500", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		mockService.
			On("GetByUserID", mock.Anything, int64(1)).
			Return(nil, errors.New("db error"))

		w := exec(h.GetByUserID, "/addresses/user/1", "1")

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})
}

func TestAddressHandler_GetByClientCpfID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("sucesso", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		expectedAddresses := []*models.Address{
			{ID: 1, Street: "Rua Cliente", ClientCpfID: utils.Int64Ptr(12345678900)},
		}

		mockService.
			On("GetByClientCpfID", mock.Anything, int64(1)).
			Return(expectedAddresses, nil)

		w := exec(h.GetByClientCpfID, "/addresses/client/1", "1")

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var response utils.DefaultResponse
		decode(w, &response)

		assert.Equal(t, http.StatusOK, response.Status)
		assert.Equal(t, "Endereços do cliente encontrados", response.Message)
		assert.NotNil(t, response.Data)
	})

	t.Run("ID inválido retorna 400", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodGet, "/addresses/client/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.GetByClientCpfID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ErrNotFound retorna 404", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		mockService.
			On("GetByClientCpfID", mock.Anything, int64(1)).
			Return(nil, errMsg.ErrNotFound)

		w := exec(h.GetByClientCpfID, "/addresses/client/1", "1")

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})

	t.Run("erro genérico retorna 500", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		mockService.
			On("GetByClientCpfID", mock.Anything, int64(1)).
			Return(nil, errors.New("db error"))

		w := exec(h.GetByClientCpfID, "/addresses/client/1", "1")

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})
}

func TestAddressHandler_GetBySupplierID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("sucesso", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		expectedAddresses := []*models.Address{
			{ID: 1, Street: "Rua Fornecedor", SupplierID: utils.Int64Ptr(1)},
			{ID: 2, Street: "Rua Fornecedor 2", SupplierID: utils.Int64Ptr(1)},
		}

		mockService.
			On("GetBySupplierID", mock.Anything, int64(1)).
			Return(expectedAddresses, nil)

		w := exec(h.GetBySupplierID, "/addresses/supplier/1", "1")

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var response utils.DefaultResponse
		decode(w, &response)

		assert.Equal(t, http.StatusOK, response.Status)
		assert.Equal(t, "Endereços do fornecedor encontrados", response.Message)
		assert.NotNil(t, response.Data)
	})

	t.Run("ID inválido retorna 400", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodGet, "/addresses/supplier/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.GetBySupplierID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ErrNotFound retorna 404", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		mockService.
			On("GetBySupplierID", mock.Anything, int64(1)).
			Return(nil, errMsg.ErrNotFound)

		w := exec(h.GetBySupplierID, "/addresses/supplier/1", "1")

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})

	t.Run("erro genérico retorna 500", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		mockService.
			On("GetBySupplierID", mock.Anything, int64(1)).
			Return(nil, errors.New("db error"))

		w := exec(h.GetBySupplierID, "/addresses/supplier/1", "1")

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
	})
}

func TestAddressHandler_handleGetAddresses(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("lista vazia retorna sucesso com array vazio", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		mockService.
			On("GetByUserID", mock.Anything, int64(1)).
			Return([]*models.Address{}, nil)

		w := exec(h.GetByUserID, "/addresses/user/1", "1")

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var response utils.DefaultResponse
		decode(w, &response)

		assert.Equal(t, http.StatusOK, response.Status)
		assert.Equal(t, "Endereços do usuário encontrados", response.Message)

		addresses, ok := response.Data.([]interface{})
		assert.True(t, ok)
		assert.Empty(t, addresses)
	})

	t.Run("sem parâmetro ID retorna erro", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodGet, "/addresses/user/", nil)
		// Não seta variáveis de URL
		w := httptest.NewRecorder()

		h.GetByUserID(w, req)

		// Espera-se que GetIDParam retorne erro quando não há parâmetro "id"
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})
}

/* helpers */

func exec(handler http.HandlerFunc, url, id string) *httptest.ResponseRecorder {
	req := newRequest(url, id)
	w := httptest.NewRecorder()
	handler(w, req)
	return w
}

func newRequest(url, id string) *http.Request {
	req := httptest.NewRequest(http.MethodGet, url, nil)
	return mux.SetURLVars(req, map[string]string{"id": id})
}

func decode[T any](w *httptest.ResponseRecorder, dest *T) {
	require.NoError(nil, json.Unmarshal(w.Body.Bytes(), dest))
}
