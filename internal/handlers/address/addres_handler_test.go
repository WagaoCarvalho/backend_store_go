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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	addresses_services "github.com/WagaoCarvalho/backend_store_go/internal/services/addresses"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/addresses"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
)

func TestAddressHandler_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService)

		uid := int64(1)
		input := &models.Address{
			UserID:     &uid,
			Street:     "Rua Exemplo",
			City:       "Cidade",
			State:      "Estado",
			Country:    "Brasil",
			PostalCode: "12345",
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
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Endereço criado com sucesso", response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("CreateValidationError", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService)

		userID := int64(1)
		input := &models.Address{
			UserID:     &userID,
			Street:     "",
			City:       "São Paulo",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "01000-000",
		}

		mockService.On("Create", mock.Anything, mock.Anything).Return(nil)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/addresses", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		expected := `{
		"status": 400,
		"message": "Erro no campo 'Street': campo obrigatório",
		"data": null
	}`

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())

		mockService.AssertNotCalled(t, "Create", mock.Anything, mock.Anything)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService)
		req := httptest.NewRequest(http.MethodPost, "/addresses", bytes.NewBuffer([]byte(`{invalid`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		handler.Create(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService)

		uid := int64(42)
		input := &models.Address{
			UserID:     &uid,
			Street:     "Rua Falha",
			City:       "ErroCity",
			State:      "Estado",
			Country:    "Brasil",
			PostalCode: "00000",
		}

		mockService.On("Create", mock.Anything, input).Return(&models.Address{}, assert.AnError)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/addresses", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})

}

func TestAddressHandler_GetByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService)
		expected := &models.Address{ID: int64(0), Street: "Rua", City: "Cidade"}
		mockService.On("GetByID", mock.Anything, int64(1)).Return(expected, nil)
		req := httptest.NewRequest(http.MethodGet, "/addresses/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		handler.GetByID(w, req)
		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Endereço encontrado", response["message"])
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService)
		req := httptest.NewRequest(http.MethodGet, "/addresses/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()
		handler.GetByID(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService)
		mockService.On("GetByID", mock.Anything, int64(1)).Return(&models.Address{}, assert.AnError)
		req := httptest.NewRequest(http.MethodGet, "/addresses/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()
		handler.GetByID(w, req)
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
	})
}

func TestAddressHandler_Update(t *testing.T) {
	t.Run("InvalidIDParam", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService)

		req := httptest.NewRequest(http.MethodPut, "/addresses/abc", strings.NewReader(`{}`))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rr := httptest.NewRecorder()

		handler.Update(rr, req)

		expected := `{
			"status": 400,
			"message": "strconv.ParseInt: parsing \"abc\": invalid syntax",
			"data": null
		}`

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.JSONEq(t, expected, rr.Body.String())
	})

	t.Run("ValidateReturnsGenericError", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService)

		userID := int64(1)
		input := &models.Address{
			UserID:     &userID,
			Street:     "cause_generic_error", // força erro genérico
			City:       "Cidade",
			State:      "ST",
			Country:    "Pais",
			PostalCode: "12345-678",
		}

		body, _ := json.Marshal(input)

		req := httptest.NewRequest(http.MethodPut, "/addresses/3", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "3"})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Update(w, req)

		expected := `{
		"status": 400,
		"message": "erro genérico na validação",
		"data": null
	}`

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())

		mockService.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService)

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
		handler := NewAddressHandler(mockService)

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

	t.Run("ValidationErrorFromValidateMethod", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService)

		userID := int64(1)
		input := &models.Address{
			ID:         2,
			UserID:     &userID,
			Street:     "", // campo obrigatório vazio -> erro na validação local
			City:       "São Paulo",
			State:      "SP",
			Country:    "Brasil",
			PostalCode: "01000-000",
		}

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

		// Update do serviço NÃO deve ser chamado pois erro ocorreu na validação local
		mockService.AssertNotCalled(t, "Update", mock.Anything, mock.Anything)
	})

	t.Run("ValidationErrorFromService", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService)

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
		handler := NewAddressHandler(mockService)

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
	t.Run("Success", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService)

		mockService.On("Delete", mock.Anything, int64(1), 2).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/addresses/1?version=2", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})

		w := httptest.NewRecorder()
		handler.Delete(w, req)

		expected := `{
			"status": 200,
			"message": "Endereço deletado com sucesso",
			"data": null
		}`

		assert.Equal(t, http.StatusOK, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService)

		req := httptest.NewRequest(http.MethodDelete, "/addresses/abc?version=1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		expected := `{
			"status": 400,
			"message": "ID inválido (esperado número inteiro)",
			"data": null
		}`

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("InvalidVersion", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService)

		req := httptest.NewRequest(http.MethodDelete, "/addresses/1?version=abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		expected := `{
			"status": 400,
			"message": "versão inválida (esperado número inteiro)",
			"data": null
		}`

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
	})

	t.Run("MissingIDError", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService)

		mockService.On("Delete", mock.Anything, int64(0), 1).Return(services.ErrAddressIDRequired)

		req := httptest.NewRequest(http.MethodDelete, "/addresses/0?version=1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "0"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		expected := `{
			"status": 400,
			"message": "endereço ID é obrigatório",
			"data": null
		}`

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("MissingVersionError", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService)

		mockService.On("Delete", mock.Anything, int64(1), 0).Return(services.ErrVersionRequired)

		req := httptest.NewRequest(http.MethodDelete, "/addresses/1?version=0", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		expected := `{
			"status": 400,
			"message": "versão é obrigatória",
			"data": null
		}`

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("VersionConflict", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService)

		mockService.On("Delete", mock.Anything, int64(1), 2).Return(services.ErrVersionConflict)

		req := httptest.NewRequest(http.MethodDelete, "/addresses/1?version=2", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		expected := `{
			"status": 409,
			"message": "conflito de versão",
			"data": null
		}`

		assert.Equal(t, http.StatusConflict, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService)

		mockService.On("Delete", mock.Anything, int64(1), 1).Return(errors.New("erro inesperado"))

		req := httptest.NewRequest(http.MethodDelete, "/addresses/1?version=1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})

		w := httptest.NewRecorder()
		handler.Delete(w, req)

		expected := `{
			"status": 500,
			"message": "erro inesperado",
			"data": null
		}`

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("DeleteNotFoundError", func(t *testing.T) {
		mockService := new(addresses_services.MockAddressService)
		handler := NewAddressHandler(mockService)

		id := int64(99)

		mockService.On("Delete", mock.Anything, id, 1).Return(utils.ErrNotFound)

		req := httptest.NewRequest(http.MethodDelete, "/addresses/99?version=1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "99"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		expected := fmt.Sprintf(`{
			"status": 404,
			"message": "%s",
			"data": null
		}`, utils.ErrNotFound.Error())

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, expected, w.Body.String())
		mockService.AssertExpectations(t)
	})
}
