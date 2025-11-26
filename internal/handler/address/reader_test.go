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

	t.Run("deve retornar endereço com sucesso", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		expected := &models.Address{
			ID:     1,
			Street: "Rua Exemplo",
			City:   "Cidade",
		}

		mockService.On("GetByID", mock.Anything, int64(1)).Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/addresses/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.GetByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Endereço encontrado", response.Message)

		data := response.Data.(map[string]interface{})
		assert.Equal(t, float64(expected.ID), data["id"])
		assert.Equal(t, expected.Street, data["street"])
		assert.Equal(t, expected.City, data["city"])

		mockService.AssertExpectations(t)
	})

	t.Run("deve retornar erro 400 se o ID for inválido", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodGet, "/addresses/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.GetByID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("deve retornar erro 404 se o serviço retornar erro", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		mockService.On("GetByID", mock.Anything, int64(1)).Return((*models.Address)(nil), errors.New("not found"))

		req := httptest.NewRequest(http.MethodGet, "/addresses/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.GetByID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestAddressHandler_GetByUserID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	loggerAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, loggerAdapter)

		expected := []*models.Address{
			{ID: 1, Street: "Rua 1", City: "Cidade A"},
			{ID: 2, Street: "Rua 2", City: "Cidade B"},
		}

		mockService.On("GetByUserID", mock.Anything, int64(1)).
			Return(expected, nil)

		w := execRequest(h.GetByUserID, "GET", "/addresses/user/1", map[string]string{"id": "1"})
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var resp utils.DefaultResponse
		require.NoError(t, decodeResponse(w, &resp))

		assert.Equal(t, "Endereços do usuário encontrados", resp.Message)
		assert.Len(t, resp.Data.([]interface{}), 2)

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, loggerAdapter)

		w := execRequest(h.GetByUserID, "GET", "/addresses/user/abc", map[string]string{"id": "abc"})
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, loggerAdapter)

		mockService.On("GetByUserID", mock.Anything, int64(1)).
			Return(nil, assert.AnError)

		w := execRequest(h.GetByUserID, "GET", "/addresses/user/1", map[string]string{"id": "1"})
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)

		mockService.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, loggerAdapter)

		mockService.On("GetByUserID", mock.Anything, int64(1)).
			Return(nil, errMsg.ErrNotFound)

		w := execRequest(h.GetByUserID, "GET", "/addresses/user/1", map[string]string{"id": "1"})

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)

		var resp utils.DefaultResponse
		require.NoError(t, decodeResponse(w, &resp))

		assert.Equal(t, errMsg.ErrNotFound.Error(), resp.Message)

		mockService.AssertExpectations(t)
	})
}

func TestAddressHandler_GetByClientID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	loggerAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, loggerAdapter)

		expected := []*models.Address{
			{ID: 1, ClientID: utils.Int64Ptr(1), Street: "Rua 1", City: "Cidade A"},
			{ID: 2, ClientID: utils.Int64Ptr(1), Street: "Rua 2", City: "Cidade B"},
		}

		mockService.On("GetByClientID", mock.Anything, int64(1)).
			Return(expected, nil)

		w := execRequest(h.GetByClientID, "GET", "/addresses/client/1", map[string]string{"id": "1"})
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, loggerAdapter)

		w := execRequest(h.GetByClientID, "GET", "/addresses/client/abc", map[string]string{"id": "abc"})
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, loggerAdapter)

		mockService.On("GetByClientID", mock.Anything, int64(1)).
			Return(nil, assert.AnError)

		w := execRequest(h.GetByClientID, "GET", "/addresses/client/1", map[string]string{"id": "1"})
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)

		mockService.AssertExpectations(t)
	})
}

func TestAddressHandler_GetBySupplierID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	loggerAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, loggerAdapter)

		expected := []*models.Address{
			{ID: 1, SupplierID: utils.Int64Ptr(1), Street: "Rua 1", City: "Cidade A"},
			{ID: 2, SupplierID: utils.Int64Ptr(1), Street: "Rua 2", City: "Cidade B"},
		}

		mockService.On("GetBySupplierID", mock.Anything, int64(1)).
			Return(expected, nil)

		w := execRequest(h.GetBySupplierID, "GET", "/addresses/supplier/1", map[string]string{"id": "1"})
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, loggerAdapter)

		w := execRequest(h.GetBySupplierID, "GET", "/addresses/supplier/abc", map[string]string{"id": "abc"})
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, loggerAdapter)

		mockService.On("GetBySupplierID", mock.Anything, int64(1)).
			Return(nil, assert.AnError)

		w := execRequest(h.GetBySupplierID, "GET", "/addresses/supplier/1", map[string]string{"id": "1"})
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)

		mockService.AssertExpectations(t)
	})
}

func execRequest(
	handler http.HandlerFunc,
	method string,
	url string,
	vars map[string]string,
) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, url, nil)
	if vars != nil {
		req = mux.SetURLVars(req, vars)
	}
	w := httptest.NewRecorder()
	handler(w, req)
	return w
}

func decodeResponse[T any](w *httptest.ResponseRecorder, dest *T) error {
	return json.Unmarshal(w.Body.Bytes(), dest)
}
