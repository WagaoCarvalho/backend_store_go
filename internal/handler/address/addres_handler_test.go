package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	err_msg_pg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/db"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators"
	service_mock "github.com/WagaoCarvalho/backend_store_go/internal/service/address/address_services_mock"
)

func newRequestWithVars(method, url string, body []byte, vars map[string]string) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	return mux.SetURLVars(req, vars)
}

func TestAddressHandler_Create(t *testing.T) {
	logAdapter := logger.NewLoggerAdapter(logrus.New())

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockService := new(service_mock.MockAddressService)
		handler := NewAddressHandler(mockService, logAdapter)

		uid := int64(1)
		input := &models.Address{
			UserID:     &uid,
			Street:     "Rua Exemplo",
			City:       "Cidade",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "12345-678",
		}
		expected := *input
		expected.ID = 1

		mockService.On("Create", mock.Anything, input).Return(&expected, nil)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/addresses", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Endereço criado com sucesso", response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		t.Parallel()
		mockService := new(service_mock.MockAddressService)
		handler := NewAddressHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodPost, "/addresses", bytes.NewBuffer([]byte(`{invalid`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ForeignKey inválida deve retornar 400", func(t *testing.T) {
		t.Parallel()
		mockService := new(service_mock.MockAddressService)
		handler := NewAddressHandler(mockService, logAdapter)

		userID := int64(99)
		input := &models.Address{
			UserID:     &userID,
			Street:     "Rua FK",
			City:       "Cidade FK",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "99999-999",
		}

		mockService.On("Create", mock.Anything, input).Return((*models.Address)(nil), err_msg_pg.ErrInvalidForeignKey)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/addresses", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		t.Parallel()
		mockService := new(service_mock.MockAddressService)
		handler := NewAddressHandler(mockService, logAdapter)

		uid := int64(42)
		input := &models.Address{
			UserID:     &uid,
			Street:     "Rua Falha",
			City:       "ErroCity",
			State:      "MG",
			Country:    "Brasil",
			PostalCode: "02460-000",
		}

		mockService.On("Create", mock.Anything, input).Return((*models.Address)(nil), assert.AnError)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/addresses", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestAddressHandler_GetByID(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	t.Run("deve retornar endereço com sucesso", func(t *testing.T) {
		mockService := new(service_mock.MockAddressService)
		handler := NewAddressHandler(mockService, logger)

		expected := &models.Address{
			ID:     1,
			Street: "Rua Exemplo",
			City:   "Cidade",
		}

		mockService.On("GetByID", mock.Anything, int64(1)).Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/addresses/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Endereço encontrado", response.Message)
		assert.Equal(t, float64(1), response.Data.(map[string]interface{})["id"])

		mockService.AssertExpectations(t)
	})

	t.Run("deve retornar erro 400 se o ID for inválido", func(t *testing.T) {
		mockService := new(service_mock.MockAddressService)
		handler := NewAddressHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodGet, "/addresses/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("deve retornar erro 404 se o serviço retornar erro", func(t *testing.T) {
		mockService := new(service_mock.MockAddressService)
		handler := NewAddressHandler(mockService, logger)

		mockService.On("GetByID", mock.Anything, int64(1)).Return((*models.Address)(nil), assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/addresses/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestAddressHandler_GetByUserID(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	t.Run("Success", func(t *testing.T) {
		mockService := new(service_mock.MockAddressService)
		handler := NewAddressHandler(mockService, logger)
		expected := []*models.Address{
			{ID: int64(1), Street: "Rua 1", City: "Cidade A"},
			{ID: int64(2), Street: "Rua 2", City: "Cidade B"},
		}
		mockService.On("GetByUserID", mock.Anything, int64(1)).Return(expected, nil)
		req := httptest.NewRequest(http.MethodGet, "/addresses/user/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		handler.GetByUserID(w, req)
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
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockService := new(service_mock.MockAddressService)
		handler := NewAddressHandler(mockService, logger)
		req := httptest.NewRequest(http.MethodGet, "/addresses/user/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()
		handler.GetByUserID(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(service_mock.MockAddressService)
		handler := NewAddressHandler(mockService, logger)
		mockService.On("GetByUserID", mock.Anything, int64(1)).Return(nil, assert.AnError)
		req := httptest.NewRequest(http.MethodGet, "/addresses/user/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		handler.GetByUserID(w, req)
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})
}

func TestAddressHandler_GetByClientID(t *testing.T) {
	logAdapter := logger.NewLoggerAdapter(logrus.New())

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockService := new(service_mock.MockAddressService)
		handler := NewAddressHandler(mockService, logAdapter)

		expected := []*models.Address{
			{ID: int64(1), Street: "Rua Cliente 1", City: "Cidade 1"},
			{ID: int64(2), Street: "Rua Cliente 2", City: "Cidade 2"},
		}
		mockService.On("GetByClientID", mock.Anything, int64(1)).Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/addresses/client/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByClientID(w, req)

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

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		t.Parallel()
		mockService := new(service_mock.MockAddressService)
		handler := NewAddressHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodGet, "/addresses/client/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetByClientID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		t.Parallel()
		mockService := new(service_mock.MockAddressService)
		handler := NewAddressHandler(mockService, logAdapter)

		mockService.On("GetByClientID", mock.Anything, int64(1)).Return(nil, assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/addresses/client/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByClientID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestAddressHandler_GetBySupplierID(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	t.Run("Success", func(t *testing.T) {
		mockService := new(service_mock.MockAddressService)
		handler := NewAddressHandler(mockService, logger)
		expected := []*models.Address{
			{ID: int64(1), Street: "Rua Fornecedor 1", City: "Cidade 1"},
			{ID: int64(2), Street: "Rua Fornecedor 2", City: "Cidade 2"},
		}
		mockService.On("GetBySupplierID", mock.Anything, int64(1)).Return(expected, nil)
		req := httptest.NewRequest(http.MethodGet, "/addresses/supplier/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetBySupplierID(w, req)

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

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockService := new(service_mock.MockAddressService)
		handler := NewAddressHandler(mockService, logger)
		req := httptest.NewRequest(http.MethodGet, "/addresses/supplier/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetBySupplierID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(service_mock.MockAddressService)
		handler := NewAddressHandler(mockService, logger)
		mockService.On("GetBySupplierID", mock.Anything, int64(1)).Return(nil, assert.AnError)
		req := httptest.NewRequest(http.MethodGet, "/addresses/supplier/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetBySupplierID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})
}

func TestAddressHandler_Update(t *testing.T) {
	logAdapter := logger.NewLoggerAdapter(logrus.New())

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(service_mock.MockAddressService)
		handler := NewAddressHandler(mockSvc, logAdapter)

		addr := models.Address{ID: 1, Street: "Rua Atualizada"}
		mockSvc.On("Update", mock.Anything, &addr).Return(nil).Once()

		body, _ := json.Marshal(addr)
		req := newRequestWithVars("PUT", "/addresses/1", body, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(service_mock.MockAddressService)
		handler := NewAddressHandler(mockSvc, logAdapter)

		req := newRequestWithVars("PUT", "/addresses/1", []byte("{invalid"), map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(service_mock.MockAddressService)
		handler := NewAddressHandler(mockSvc, logAdapter)

		addr := models.Address{ID: 1, Street: ""}
		validationErr := &validators.ValidationError{Field: "Street", Message: "campo obrigatório"}

		mockSvc.On("Update", mock.Anything, &addr).Return(validationErr).Once()

		body, _ := json.Marshal(addr)
		req := newRequestWithVars("PUT", "/addresses/1", body, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(service_mock.MockAddressService)
		handler := NewAddressHandler(mockSvc, logAdapter)

		addr := models.Address{ID: 1, Street: "Rua Falha"}
		mockSvc.On("Update", mock.Anything, &addr).Return(errors.New("erro ao atualizar")).Once()

		body, _ := json.Marshal(addr)
		req := newRequestWithVars("PUT", "/addresses/1", body, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(service_mock.MockAddressService)
		handler := NewAddressHandler(mockSvc, logAdapter)

		req := newRequestWithVars("PUT", "/addresses/abc", []byte("{}"), map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestAddressHandler_Delete(t *testing.T) {
	logAdapter := logger.NewLoggerAdapter(logrus.New())

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(service_mock.MockAddressService)
		handler := NewAddressHandler(mockSvc, logAdapter)

		mockSvc.On("Delete", mock.Anything, int64(1)).Return(nil).Once()

		req := newRequestWithVars("DELETE", "/addresses/1", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		assert.Empty(t, w.Body.String())

		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(service_mock.MockAddressService)
		handler := NewAddressHandler(mockSvc, logAdapter)

		req := newRequestWithVars("DELETE", "/addresses/abc", nil, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Contains(t, w.Body.String(), "invalid ID format: abc")

	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(service_mock.MockAddressService)
		handler := NewAddressHandler(mockSvc, logAdapter)

		mockSvc.On("Delete", mock.Anything, int64(99)).Return(validators.ErrNotFound).Once()

		req := newRequestWithVars("DELETE", "/addresses/99", nil, map[string]string{"id": "99"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.Contains(t, w.Body.String(), validators.ErrNotFound.Error())

		mockSvc.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(service_mock.MockAddressService)
		handler := NewAddressHandler(mockSvc, logAdapter)

		mockSvc.On("Delete", mock.Anything, int64(2)).Return(err_msg.ErrAddressIDRequired).Once()

		req := newRequestWithVars("DELETE", "/addresses/2", nil, map[string]string{"id": "2"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Contains(t, w.Body.String(), "erro id deve ser maior que 0")

		mockSvc.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(service_mock.MockAddressService)
		handler := NewAddressHandler(mockSvc, logAdapter)

		mockSvc.On("Delete", mock.Anything, int64(3)).Return(assert.AnError).Once()

		req := newRequestWithVars("DELETE", "/addresses/3", nil, map[string]string{"id": "3"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		assert.Contains(t, w.Body.String(), assert.AnError.Error())

		mockSvc.AssertExpectations(t)
	})
}
