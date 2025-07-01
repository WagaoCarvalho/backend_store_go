package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	models "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/addresses"
	addresses_services "github.com/WagaoCarvalho/backend_store_go/internal/services/addresses/address_services_mock"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
)

func TestAddressHandler_Create(t *testing.T) {
	logAdapter := logger.NewLoggerAdapter(logrus.New()) // evita conflito de nome com pacote

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockService := new(addresses_services.MockAddressService)
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
		mockService := new(addresses_services.MockAddressService)
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
		mockService := new(addresses_services.MockAddressService)
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

		mockService.On("Create", mock.Anything, input).Return((*models.Address)(nil), repositories.ErrInvalidForeignKey)

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
		mockService := new(addresses_services.MockAddressService)
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
		mockService := new(addresses_services.MockAddressService)
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
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodGet, "/addresses/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("deve retornar erro 404 se o serviço retornar erro", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
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
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService, logger) // passa logger no handler
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
		assert.Equal(t, float64(http.StatusOK), response["status"]) // JSON decodifica ints como float64
		data, ok := response["data"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, data, 2)
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService, logger)
		req := httptest.NewRequest(http.MethodGet, "/addresses/user/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()
		handler.GetByUserID(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
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
		mockService := new(addresses_services.MockAddressService)
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
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodGet, "/addresses/client/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetByClientID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		t.Parallel()
		mockService := new(addresses_services.MockAddressService)
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
	logger := logger.NewLoggerAdapter(logrus.New()) // logger real ou mock

	t.Run("Success", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
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
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService, logger)
		req := httptest.NewRequest(http.MethodGet, "/addresses/supplier/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetBySupplierID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
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
	logAdapter := logger.NewLoggerAdapter(logrus.New()) // Logger real ou mock

	t.Run("InvalidIDParam", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService, logAdapter)
		req := httptest.NewRequest(http.MethodPut, "/addresses/abc", strings.NewReader(`{}`))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rr := httptest.NewRecorder()
		handler.Update(rr, req)

		expected := `{
			"status": 400,
			"message": "invalid ID format: abc",
			"data": null
		}`
		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService, logAdapter)
		req := httptest.NewRequest(http.MethodPut, "/addresses/2", bytes.NewBuffer([]byte(`{invalid-json}`)))
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.Update(w, req)

		expected := `{
			"status": 400,
			"message": "erro ao decodificar JSON: invalid character 'i' looking for beginning of object key string",
			"data": null
		}`
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("Success", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService, logAdapter)
		var userID int64 = 1
		input := &models.Address{
			ID:         0,
			UserID:     &userID,
			ClientID:   nil,
			SupplierID: nil,
			Street:     "Nova Rua",
			City:       "São Paulo",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "01000-000",
		}
		id := int64(2)
		inputWithID := *input
		inputWithID.ID = id

		mockService.On("Update", mock.Anything, &inputWithID).Return(nil)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPut, "/addresses/2", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Update(w, req)

		expected := `{
			"status": 200,
			"message": "Endereço atualizado com sucesso",
			"data": null
		}`
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("ValidationErrorFromService", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService, logAdapter)
		userID := int64(1)
		input := &models.Address{
			ID:         0,
			UserID:     &userID,
			ClientID:   nil,
			SupplierID: nil,
			Street:     "Nova Rua",
			City:       "São Paulo",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "01000-000",
		}
		id := int64(2)
		inputWithID := *input
		inputWithID.ID = id

		validationErr := &utils.ValidationError{
			Field:   "Street",
			Message: "campo obrigatório",
		}

		mockService.On("Update", mock.Anything, &inputWithID).Return(validationErr)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPut, "/addresses/2", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Update(w, req)

		expected := `{
			"status": 400,
			"message": "Erro no campo 'Street': campo obrigatório",
			"data": null
		}`
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService, logAdapter)
		userID := int64(1)
		input := &models.Address{
			ID:         0,
			UserID:     &userID,
			ClientID:   nil,
			SupplierID: nil,
			Street:     "Nova Rua",
			City:       "São Paulo",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "01000-000",
		}
		id := int64(2)
		inputWithID := *input
		inputWithID.ID = id

		mockService.On("Update", mock.Anything, &inputWithID).Return(assert.AnError)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPut, "/addresses/2", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Update(w, req)

		expected := fmt.Sprintf(`{
			"status": 500,
			"message": "%s",
			"data": null
		}`, assert.AnError.Error())

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
		mockService.AssertExpectations(t)
	})
}

func TestAddressHandler_Delete(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New()) // Logger simples para os testes

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService, logger)
		mockService.On("Delete", mock.Anything, int64(1)).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/addresses/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		expected := `{"status":200,"message":"Endereço deletado com sucesso","data":null}`
		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		t.Parallel()
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodDelete, "/addresses/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		expected := `{"status":400,"message":"ID inválido (esperado número inteiro)","data":null}`
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("MissingIDError", func(t *testing.T) {
		t.Parallel()
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService, logger)

		mockService.On("Delete", mock.Anything, int64(0)).Return(services.ErrAddressIDRequired)

		req := httptest.NewRequest(http.MethodDelete, "/addresses/0", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "0"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		expected := `{"status":400,"message":"endereço ID é obrigatório","data":null}`
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		t.Parallel()
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService, logger)

		mockService.On("Delete", mock.Anything, int64(1)).Return(errors.New("erro inesperado"))

		req := httptest.NewRequest(http.MethodDelete, "/addresses/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		expected := `{"status":500,"message":"erro inesperado","data":null}`
		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("DeleteNotFoundError", func(t *testing.T) {
		t.Parallel()
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService, logger)

		id := int64(99)
		mockService.On("Delete", mock.Anything, id).Return(utils.ErrNotFound)

		req := httptest.NewRequest(http.MethodDelete, "/addresses/99", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "99"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		expected := fmt.Sprintf(`{"status":404,"message":"%s","data":null}`, utils.ErrNotFound.Error())
		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
		mockService.AssertExpectations(t)
	})
}
