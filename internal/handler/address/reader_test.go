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
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		expected := []*models.Address{
			{ID: 1, Street: "Rua 1", City: "Cidade A"},
			{ID: 2, Street: "Rua 2", City: "Cidade B"},
		}

		mockService.On("GetByUserID", mock.Anything, int64(1)).Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/addresses/user/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.GetByUserID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Endereços do usuário encontrados", response["message"])
		assert.Equal(t, float64(http.StatusOK), response["status"])

		data, ok := response["data"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, data, 2)

		first := data[0].(map[string]interface{})
		assert.Equal(t, float64(expected[0].ID), first["id"])
		assert.Equal(t, expected[0].Street, first["street"])
		assert.Equal(t, expected[0].City, first["city"])

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodGet, "/addresses/user/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.GetByUserID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		mockService.On("GetByUserID", mock.Anything, int64(1)).Return(nil, assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/addresses/user/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.GetByUserID(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		// Mock do service retornando ErrNotFound
		mockService.On("GetByUserID", mock.Anything, int64(1)).
			Return(nil, errMsg.ErrNotFound)

		req := httptest.NewRequest(http.MethodGet, "/addresses/user/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.GetByUserID(w, req)

		// Verifica se retornou 404
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)

		// Verifica se a mensagem contém "usuário não encontrado"
		var resp utils.DefaultResponse
		require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
		assert.Contains(t, resp.Message, "usuário não encontrado")

		mockService.AssertExpectations(t)
	})

}

func TestAddressHandler_GetByClientID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		expected := []*models.Address{
			{ID: 1, ClientID: utils.Int64Ptr(1), Street: "Rua 1", City: "Cidade A"},
			{ID: 2, ClientID: utils.Int64Ptr(1), Street: "Rua 2", City: "Cidade B"},
		}

		mockService.On("GetByClientID", mock.Anything, int64(1)).Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/addresses/client/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.GetByClientID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Endereços do cliente encontrados", response["message"])
		assert.Equal(t, float64(http.StatusOK), response["status"])

		data, ok := response["data"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, data, 2)

		first := data[0].(map[string]interface{})
		assert.Equal(t, float64(expected[0].ID), first["id"])
		assert.Equal(t, expected[0].Street, first["street"])
		assert.Equal(t, expected[0].City, first["city"])

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodGet, "/addresses/client/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.GetByClientID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		mockService.On("GetByClientID", mock.Anything, int64(1)).Return(nil, assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/addresses/client/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.GetByClientID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestAddressHandler_GetBySupplierID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		expected := []*models.Address{
			{ID: 1, SupplierID: utils.Int64Ptr(1), Street: "Rua 1", City: "Cidade A"},
			{ID: 2, SupplierID: utils.Int64Ptr(1), Street: "Rua 2", City: "Cidade B"},
		}

		mockService.On("GetBySupplierID", mock.Anything, int64(1)).Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/addresses/supplier/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.GetBySupplierID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Endereços do fornecedor encontrados", response["message"])
		assert.Equal(t, float64(http.StatusOK), response["status"])

		data, ok := response["data"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, data, 2)

		first := data[0].(map[string]interface{})
		assert.Equal(t, float64(expected[0].ID), first["id"])
		assert.Equal(t, expected[0].Street, first["street"])
		assert.Equal(t, expected[0].City, first["city"])

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodGet, "/addresses/supplier/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.GetBySupplierID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(mockAddress.MockAddress)
		h := NewAddressHandler(mockService, logAdapter)

		mockService.On("GetBySupplierID", mock.Anything, int64(1)).Return(nil, assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/addresses/supplier/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.GetBySupplierID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})
}
