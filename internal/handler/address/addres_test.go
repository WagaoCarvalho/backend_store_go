package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	mockAddress "github.com/WagaoCarvalho/backend_store_go/infra/mock/service/address"
	dtoAddress "github.com/WagaoCarvalho/backend_store_go/internal/dto/address"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
)

func newRequestWithVars(method, url string, body []byte, vars map[string]string) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	return mux.SetURLVars(req, vars)
}

func TestAddressHandler_Create(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success - Create Address", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logAdapter)

		uid := int64(10)
		inputDTO := &dtoAddress.AddressDTO{
			ID:         &uid,
			UserID:     &uid,
			Street:     "Rua Exemplo",
			City:       "Cidade",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "12345-678",
		}

		// Mock do service retornando o mesmo DTO
		mockService.On("Create", mock.Anything, inputDTO).Return(inputDTO, nil)

		// Serializa o request
		body, _ := json.Marshal(inputDTO)
		req := httptest.NewRequest(http.MethodPost, "/addresses", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		// Desserializa o response em struct compatível com DTO
		var response struct {
			Status  int                   `json:"status"`
			Message string                `json:"message"`
			Data    dtoAddress.AddressDTO `json:"data"`
		}

		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Endereço criado com sucesso", response.Message)

		// Valida campos do DTO retornado
		assert.NotNil(t, response.Data.ID)
		assert.Equal(t, *inputDTO.ID, *response.Data.ID)
		assert.NotNil(t, response.Data.UserID)
		assert.Equal(t, *inputDTO.UserID, *response.Data.UserID)
		assert.Equal(t, inputDTO.Street, response.Data.Street)
		assert.Equal(t, inputDTO.City, response.Data.City)
		assert.Equal(t, inputDTO.State, response.Data.State)
		assert.Equal(t, inputDTO.Country, response.Data.Country)
		assert.Equal(t, inputDTO.PostalCode, response.Data.PostalCode)

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logAdapter)

		uid := int64(42)
		input := &dtoAddress.AddressDTO{
			UserID:     &uid,
			Street:     "Rua Falha",
			City:       "ErroCity",
			State:      "MG",
			Country:    "Brasil",
			PostalCode: "02460-000",
		}

		// Mock do service agora recebe e retorna DTO
		mockService.On("Create", mock.Anything, input).Return((*dtoAddress.AddressDTO)(nil), assert.AnError)

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

	t.Run("InvalidJSON", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockAddress.MockAddressService)
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
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logAdapter)

		uid := int64(99)
		input := &dtoAddress.AddressDTO{
			UserID:     &uid,
			Street:     "Rua FK",
			City:       "Cidade FK",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "99999-999",
		}

		// Mock do service retorna erro de foreign key
		mockService.On("Create", mock.Anything, input).Return((*dtoAddress.AddressDTO)(nil), errMsg.ErrInvalidForeignKey)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/addresses", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Opcional: verificar que a mensagem de log foi escrita
		// Isso depende de como você configurou o logAdapter/mocks

		mockService.AssertExpectations(t)
	})

}

func TestAddressHandler_GetByID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("deve retornar endereço com sucesso", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logger)

		expectedID := int64(1)
		expected := &dtoAddress.AddressDTO{
			ID:     &expectedID,
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

		dataMap := response.Data.(map[string]interface{})
		assert.Equal(t, float64(*expected.ID), dataMap["id"])
		assert.Equal(t, expected.Street, dataMap["street"])
		assert.Equal(t, expected.City, dataMap["city"])

		mockService.AssertExpectations(t)
	})

	t.Run("deve retornar erro 400 se o ID for inválido", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodGet, "/addresses/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("deve retornar erro 404 se o serviço retornar erro", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logger)

		mockService.On("GetByID", mock.Anything, int64(1)).Return((*dtoAddress.AddressDTO)(nil), assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/addresses/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestAddressHandler_GetByUserID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logger)

		expected := []*dtoAddress.AddressDTO{
			{ID: utils.Int64Ptr(1), Street: "Rua 1", City: "Cidade A"},
			{ID: utils.Int64Ptr(2), Street: "Rua 2", City: "Cidade B"},
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
		assert.Len(t, data, 2) // Agora deve passar

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodGet, "/addresses/user/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetByUserID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logger)

		mockService.On("GetByUserID", mock.Anything, int64(1)).Return(nil, assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/addresses/user/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByUserID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestAddressHandler_GetByClientID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logger)

		expected := []*dtoAddress.AddressDTO{
			{ID: utils.Int64Ptr(1), ClientID: utils.Int64Ptr(1), Street: "Rua 1", City: "Cidade A"},
			{ID: utils.Int64Ptr(2), ClientID: utils.Int64Ptr(1), Street: "Rua 2", City: "Cidade B"},
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
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodGet, "/addresses/client/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetByClientID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logger)

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
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logger)

		expected := []*dtoAddress.AddressDTO{
			{ID: utils.Int64Ptr(1), SupplierID: utils.Int64Ptr(1), Street: "Rua 1", City: "Cidade A"},
			{ID: utils.Int64Ptr(2), SupplierID: utils.Int64Ptr(1), Street: "Rua 2", City: "Cidade B"},
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
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodGet, "/addresses/supplier/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetBySupplierID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logger)

		mockService.On("GetBySupplierID", mock.Anything, int64(1)).Return(nil, assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/addresses/supplier/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetBySupplierID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestAddressHandler_Update(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	toReader := func(v any) io.Reader {
		b, _ := json.Marshal(v)
		return bytes.NewReader(b)
	}

	t.Run("Success", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logger)

		id := int64(1)
		dto := dtoAddress.AddressDTO{
			ID:         &id,
			Street:     "Rua Nova",
			City:       "Cidade Nova",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "12345678",
		}

		mockService.On("Update", mock.Anything, &dto).Return(nil)

		body := toReader(dto)
		req := httptest.NewRequest(http.MethodPut, "/addresses/1", body)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Endereço atualizado com sucesso", response["message"])
		assert.Equal(t, float64(http.StatusOK), response["status"])

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logger)

		dto := dtoAddress.AddressDTO{Street: "Rua Teste"}
		body := toReader(dto)
		req := httptest.NewRequest(http.MethodPut, "/addresses/abc", body)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ParseJSONError", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logger)

		// Corpo inválido que não é JSON
		body := strings.NewReader("invalid-json")
		req := httptest.NewRequest(http.MethodPut, "/addresses/1", body)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ValidationError", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logger)

		id := int64(1)
		dto := dtoAddress.AddressDTO{ID: &id, Street: ""} // inválido
		mockService.On("Update", mock.Anything, &dto).Return(&validators.ValidationError{Message: "campo obrigatório"})

		body := toReader(dto)
		req := httptest.NewRequest(http.MethodPut, "/addresses/1", body)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logger)

		id := int64(1)
		dto := dtoAddress.AddressDTO{ID: &id, Street: "Rua Erro"}
		mockService.On("Update", mock.Anything, &dto).Return(fmt.Errorf("erro inesperado"))

		body := toReader(dto)
		req := httptest.NewRequest(http.MethodPut, "/addresses/1", body)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("ServiceReturnsErrID", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logger)

		dto := dtoAddress.AddressDTO{
			Street:     "Rua Teste",
			City:       "Cidade",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "12345678",
		}

		body, _ := json.Marshal(dto)
		req := httptest.NewRequest(http.MethodPut, "/addresses/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		// Mock do service para retornar ErrID
		mockService.On("Update", mock.Anything, mock.MatchedBy(func(a *dtoAddress.AddressDTO) bool {
			return a != nil && a.ID != nil && *a.ID == 1
		})).Return(errMsg.ErrID)

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

}

func TestAddressHandler_Delete(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(mockAddress.MockAddressService)
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
		mockSvc := new(mockAddress.MockAddressService)
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
		mockSvc := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockSvc, logAdapter)

		mockSvc.On("Delete", mock.Anything, int64(99)).Return(errMsg.ErrNotFound).Once()

		req := newRequestWithVars("DELETE", "/addresses/99", nil, map[string]string{"id": "99"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.Contains(t, w.Body.String(), errMsg.ErrNotFound.Error())

		mockSvc.AssertExpectations(t)
	})

	t.Run("validation error", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockSvc, logAdapter)

		mockSvc.On("Delete", mock.Anything, int64(2)).Return(errMsg.ErrID).Once()

		req := newRequestWithVars("DELETE", "/addresses/2", nil, map[string]string{"id": "2"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Contains(t, w.Body.String(), "erro ID inválido")

		mockSvc.AssertExpectations(t)
	})

	t.Run("internal server error", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(mockAddress.MockAddressService)
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
