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

	mockAddress "github.com/WagaoCarvalho/backend_store_go/infra/mock/service/address"
	dtoAddress "github.com/WagaoCarvalho/backend_store_go/internal/dto/address"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
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
			UserID:     &uid,
			Street:     "Rua Exemplo",
			City:       "Cidade",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "12345-678",
		}

		// modelo esperado no service
		expectedModel := dtoAddress.ToAddressModel(*inputDTO)
		// simula retorno já com ID preenchido
		expectedModel.ID = uid

		mockService.
			On("Create", mock.Anything, mock.MatchedBy(func(m *models.Address) bool {
				return m.UserID != nil &&
					*m.UserID == *inputDTO.UserID &&
					m.Street == inputDTO.Street &&
					m.City == inputDTO.City &&
					m.State == inputDTO.State &&
					m.Country == inputDTO.Country &&
					m.PostalCode == inputDTO.PostalCode
			})).
			Return(expectedModel, nil)

		body, _ := json.Marshal(inputDTO)
		req := httptest.NewRequest(http.MethodPost, "/addresses", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response struct {
			Status  int                   `json:"status"`
			Message string                `json:"message"`
			Data    dtoAddress.AddressDTO `json:"data"`
		}

		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Endereço criado com sucesso", response.Message)

		assert.NotNil(t, response.Data.ID)
		assert.Equal(t, expectedModel.ID, *response.Data.ID)
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

		mockService.On("Create", mock.Anything, mock.Anything).
			Return((*models.Address)(nil), assert.AnError)

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

		mockService.On("Create", mock.Anything, mock.Anything).
			Return((*models.Address)(nil), errMsg.ErrInvalidForeignKey)

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
}

func TestAddressHandler_GetByID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("deve retornar endereço com sucesso", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
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
		mockService := new(mockAddress.MockAddressService)
		h := NewAddressHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodGet, "/addresses/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.GetByID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("deve retornar erro 404 se o serviço retornar erro", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
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
		mockService := new(mockAddress.MockAddressService)
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
		mockService := new(mockAddress.MockAddressService)
		h := NewAddressHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodGet, "/addresses/user/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.GetByUserID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
		h := NewAddressHandler(mockService, logAdapter)

		mockService.On("GetByUserID", mock.Anything, int64(1)).Return(nil, assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/addresses/user/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		h.GetByUserID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestAddressHandler_GetByClientID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
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
		mockService := new(mockAddress.MockAddressService)
		h := NewAddressHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodGet, "/addresses/client/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.GetByClientID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
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
		mockService := new(mockAddress.MockAddressService)
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
		mockService := new(mockAddress.MockAddressService)
		h := NewAddressHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodGet, "/addresses/supplier/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		h.GetBySupplierID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
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

func TestAddressHandler_Update(t *testing.T) {
	addressID := int64(1)
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success - Update Address", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logAdapter)

		uid := int64(10)
		inputDTO := &dtoAddress.AddressDTO{
			UserID:     &uid,
			Street:     "Rua Nova",
			City:       "Cidade Nova",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "12345678",
		}

		expectedModel := dtoAddress.ToAddressModel(*inputDTO)
		expectedModel.ID = addressID

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(m *models.Address) bool {
			return m.ID == addressID &&
				m.Street == inputDTO.Street &&
				m.City == inputDTO.City &&
				m.State == inputDTO.State &&
				m.Country == inputDTO.Country &&
				m.PostalCode == inputDTO.PostalCode &&
				m.UserID != nil && *m.UserID == uid
		})).Return(nil)

		body, _ := json.Marshal(inputDTO)
		req := httptest.NewRequest(http.MethodPut, "/addresses/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		// Set ID in URL params for Gorilla Mux
		req = mux.SetURLVars(req, map[string]string{"id": "1"})

		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Status  int                   `json:"status"`
			Message string                `json:"message"`
			Data    dtoAddress.AddressDTO `json:"data"`
		}

		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Endereço atualizado com sucesso", response.Message)
		assert.Equal(t, addressID, *response.Data.ID)
		assert.Equal(t, uid, *response.Data.UserID)

		mockService.AssertExpectations(t)
	})

	t.Run("ValidationError", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logAdapter)

		inputDTO := &dtoAddress.AddressDTO{}
		body, _ := json.Marshal(inputDTO)
		req := httptest.NewRequest(http.MethodPut, "/addresses/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		// Set ID in URL params
		req = mux.SetURLVars(req, map[string]string{"id": "1"})

		w := httptest.NewRecorder()

		mockService.On("Update", mock.Anything, mock.Anything).Return(&validators.ValidationError{
			Field:   "user_id/client_id/supplier_id",
			Message: "campo obrigatório",
		})

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("ForeignKey inválida deve retornar 400", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logAdapter)

		uid := int64(99)
		inputDTO := &dtoAddress.AddressDTO{UserID: &uid}
		body, _ := json.Marshal(inputDTO)
		req := httptest.NewRequest(http.MethodPut, "/addresses/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		req = mux.SetURLVars(req, map[string]string{"id": "1"})

		w := httptest.NewRecorder()

		mockService.On("Update", mock.Anything, mock.Anything).Return(errMsg.ErrInvalidForeignKey)

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logAdapter)

		uid := int64(42)
		inputDTO := &dtoAddress.AddressDTO{UserID: &uid, Street: "Rua Erro"}
		body, _ := json.Marshal(inputDTO)
		req := httptest.NewRequest(http.MethodPut, "/addresses/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		req = mux.SetURLVars(req, map[string]string{"id": "1"})

		w := httptest.NewRecorder()

		mockService.On("Update", mock.Anything, mock.Anything).Return(errors.New("erro genérico"))

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID error deve retornar 400", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logAdapter)

		uid := int64(99)
		inputDTO := &dtoAddress.AddressDTO{UserID: &uid}
		body, _ := json.Marshal(inputDTO)
		req := httptest.NewRequest(http.MethodPut, "/addresses/1", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		req = mux.SetURLVars(req, map[string]string{"id": "1"})

		w := httptest.NewRecorder()

		// Mock retornando ErrID específico
		mockService.On("Update", mock.Anything, mock.Anything).Return(errMsg.ErrID)

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		// Verifique também a mensagem de erro se necessário
		var errorResponse struct {
			Status  int    `json:"status"`
			Message string `json:"message"`
			Error   string `json:"error"`
		}

		err := json.NewDecoder(resp.Body).Decode(&errorResponse)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, errorResponse.Status)
		assert.Contains(t, errorResponse.Message, "ID") // ou verifique a mensagem específica

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodPut, "/addresses/1", bytes.NewBuffer([]byte("{invalid")))
		req.Header.Set("Content-Type", "application/json")

		req = mux.SetURLVars(req, map[string]string{"id": "1"})

		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Invalid ID param", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logAdapter)

		// Teste com ID inválido na URL
		req := httptest.NewRequest(http.MethodPut, "/addresses/abc", nil)
		req.Header.Set("Content-Type", "application/json")

		// IMPORTANTE: Não defina parâmetros de URL aqui - deixe o utils.GetIDParam falhar
		// req = mux.SetURLVars(req, map[string]string{"id": "abc"}) // REMOVER ESTA LINHA

		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Missing ID param", func(t *testing.T) {
		mockService := new(mockAddress.MockAddressService)
		handler := NewAddressHandler(mockService, logAdapter)

		// Teste sem parâmetro ID na URL
		req := httptest.NewRequest(http.MethodPut, "/addresses/", nil)
		req.Header.Set("Content-Type", "application/json")

		// Não defina nenhum parâmetro de URL
		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
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
