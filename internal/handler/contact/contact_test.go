package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	mockContact "github.com/WagaoCarvalho/backend_store_go/infra/mock/service/contact"
	dtoContact "github.com/WagaoCarvalho/backend_store_go/internal/dto/contact"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	validators "github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils/validators/validator"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newRequestWithVars(method, url string, body []byte, vars map[string]string) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	return mux.SetURLVars(req, vars)
}

func TestContactHandler_Create(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success - Create Contact", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockContact.MockContactService)
		handler := NewContactHandler(mockService, logAdapter)

		userID := int64(10)
		inputDTO := &dtoContact.ContactDTO{
			ID:          &userID,
			UserID:      &userID,
			ContactName: "Contato Teste",
			Email:       "teste@email.com",
			Phone:       "123456789",
		}

		mockService.On("Create", mock.Anything, inputDTO).Return(inputDTO, nil)

		body, _ := json.Marshal(inputDTO)
		req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response struct {
			Status  int                   `json:"status"`
			Message string                `json:"message"`
			Data    dtoContact.ContactDTO `json:"data"`
		}

		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Contato criado com sucesso", response.Message)

		assert.NotNil(t, response.Data.ID)
		assert.Equal(t, *inputDTO.ID, *response.Data.ID)
		assert.NotNil(t, response.Data.UserID)
		assert.Equal(t, *inputDTO.UserID, *response.Data.UserID)
		assert.Equal(t, inputDTO.ContactName, response.Data.ContactName)
		assert.Equal(t, inputDTO.Email, response.Data.Email)
		assert.Equal(t, inputDTO.Phone, response.Data.Phone)

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockContact.MockContactService)
		handler := NewContactHandler(mockService, logAdapter)

		userID := int64(42)
		input := &dtoContact.ContactDTO{
			UserID:      &userID,
			ContactName: "Contato Erro",
			Email:       "erro@email.com",
			Phone:       "987654321",
		}

		mockService.On("Create", mock.Anything, input).Return((*dtoContact.ContactDTO)(nil), assert.AnError)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBuffer(body))
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
		mockService := new(mockContact.MockContactService)
		handler := NewContactHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBuffer([]byte(`{invalid`)))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ForeignKey inválida deve retornar 400", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockContact.MockContactService)
		handler := NewContactHandler(mockService, logAdapter)

		userID := int64(99)
		input := &dtoContact.ContactDTO{
			UserID:      &userID,
			ContactName: "Contato FK",
			Email:       "fk@email.com",
			Phone:       "999999999",
		}

		mockService.On("Create", mock.Anything, input).Return((*dtoContact.ContactDTO)(nil), errMsg.ErrInvalidForeignKey)

		body, _ := json.Marshal(input)
		req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestContactHandler_GetByID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("deve retornar contato com sucesso", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContactHandler(mockService, logger)

		expectedID := int64(1)
		expected := &dtoContact.ContactDTO{
			ID:          &expectedID,
			ContactName: "Contato Teste",
			Email:       "teste@email.com",
			Phone:       "123456789",
		}

		mockService.On("GetByID", mock.Anything, int64(1)).Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/contacts/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Contato encontrado", response.Message)

		dataMap := response.Data.(map[string]interface{})
		assert.Equal(t, float64(*expected.ID), dataMap["id"])
		assert.Equal(t, expected.ContactName, dataMap["contact_name"])
		assert.Equal(t, expected.Email, dataMap["email"])
		assert.Equal(t, expected.Phone, dataMap["phone"])

		mockService.AssertExpectations(t)
	})

	t.Run("deve retornar erro 400 se o ID for inválido", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContactHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodGet, "/contacts/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("deve retornar erro 404 se o serviço retornar erro", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContactHandler(mockService, logger)

		mockService.On("GetByID", mock.Anything, int64(1)).Return((*dtoContact.ContactDTO)(nil), assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/contacts/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestContactHandler_GetByUserID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContactHandler(mockService, logger)

		expected := []*dtoContact.ContactDTO{
			{ID: utils.Int64Ptr(1), ContactName: "Contato 1", Email: "c1@email.com"},
			{ID: utils.Int64Ptr(2), ContactName: "Contato 2", Email: "c2@email.com"},
		}

		mockService.On("GetByUserID", mock.Anything, int64(1)).Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/contacts/user/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByUserID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Contatos do usuário encontrados", response["message"])
		assert.Equal(t, float64(http.StatusOK), response["status"])

		data, ok := response["data"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, data, 2)

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContactHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodGet, "/contacts/user/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetByUserID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContactHandler(mockService, logger)

		mockService.On("GetByUserID", mock.Anything, int64(1)).Return(nil, assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/contacts/user/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByUserID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestContactHandler_GetByClientID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContactHandler(mockService, logger)

		expected := []*dtoContact.ContactDTO{
			{ID: utils.Int64Ptr(1), ContactName: "Cliente 1"},
			{ID: utils.Int64Ptr(2), ContactName: "Cliente 2"},
		}

		mockService.On("GetByClientID", mock.Anything, int64(1)).Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/contacts/client/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByClientID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Contatos do cliente encontrados", response["message"])
		assert.Equal(t, float64(http.StatusOK), response["status"])

		data, ok := response["data"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, data, 2)

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContactHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodGet, "/contacts/client/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetByClientID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContactHandler(mockService, logger)

		mockService.On("GetByClientID", mock.Anything, int64(1)).Return(nil, assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/contacts/client/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByClientID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestContactHandler_GetBySupplierID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContactHandler(mockService, logger)

		expected := []*dtoContact.ContactDTO{
			{ID: utils.Int64Ptr(1), ContactName: "Fornecedor 1"},
			{ID: utils.Int64Ptr(2), ContactName: "Fornecedor 2"},
		}

		mockService.On("GetBySupplierID", mock.Anything, int64(1)).Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/contacts/supplier/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetBySupplierID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Contatos do fornecedor encontrados", response["message"])
		assert.Equal(t, float64(http.StatusOK), response["status"])

		data, ok := response["data"].([]interface{})
		assert.True(t, ok)
		assert.Len(t, data, 2)

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContactHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodGet, "/contacts/supplier/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetBySupplierID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContactHandler(mockService, logger)

		mockService.On("GetBySupplierID", mock.Anything, int64(1)).Return(nil, assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/contacts/supplier/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetBySupplierID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestContactHandler_Update(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	toReader := func(v any) io.Reader {
		b, _ := json.Marshal(v)
		return bytes.NewReader(b)
	}

	t.Run("Success", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContactHandler(mockService, logger)

		id := int64(1)
		dto := dtoContact.ContactDTO{
			ID:          &id,
			ContactName: "Nome Teste",
			Email:       "teste@email.com",
			Phone:       "11999999999",
		}

		mockService.On("Update", mock.Anything, &dto).Return(nil)

		req := httptest.NewRequest(http.MethodPut, "/contacts/1", toReader(dto))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response map[string]interface{}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Contato atualizado com sucesso", response["message"])
		assert.Equal(t, float64(http.StatusOK), response["status"])

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContactHandler(mockService, logger)

		dto := dtoContact.ContactDTO{ContactName: "Nome Teste"}
		req := httptest.NewRequest(http.MethodPut, "/contacts/abc", toReader(dto))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.Update(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ParseJSONError", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContactHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodPut, "/contacts/1", strings.NewReader("invalid-json"))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ValidationError", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContactHandler(mockService, logger)

		id := int64(1)
		dto := dtoContact.ContactDTO{ID: &id, ContactName: ""} // Nome vazio para gerar erro de validação
		mockService.On("Update", mock.Anything, &dto).Return(&validators.ValidationError{Message: "campo obrigatório"})

		req := httptest.NewRequest(http.MethodPut, "/contacts/1", toReader(dto))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContactHandler(mockService, logger)

		id := int64(1)
		dto := dtoContact.ContactDTO{ID: &id, ContactName: "Nome Erro"}
		mockService.On("Update", mock.Anything, &dto).Return(fmt.Errorf("erro inesperado"))

		req := httptest.NewRequest(http.MethodPut, "/contacts/1", toReader(dto))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("ServiceReturnsErrID", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContactHandler(mockService, logger)

		dto := dtoContact.ContactDTO{
			ContactName: "Nome Teste",
			Email:       "teste@email.com",
		}

		req := httptest.NewRequest(http.MethodPut, "/contacts/1", toReader(dto))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(c *dtoContact.ContactDTO) bool {
			return c != nil && c.ID != nil && *c.ID == 1
		})).Return(errMsg.ErrID)

		handler.Update(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestContactHandler_Delete(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(mockContact.MockContactService)
		handler := NewContactHandler(mockSvc, logAdapter)

		mockSvc.On("Delete", mock.Anything, int64(1)).Return(nil).Once()

		req := newRequestWithVars("DELETE", "/contacts/1", nil, map[string]string{"id": "1"})
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
		mockSvc := new(mockContact.MockContactService)
		handler := NewContactHandler(mockSvc, logAdapter)

		req := newRequestWithVars("DELETE", "/contacts/abc", nil, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Contains(t, w.Body.String(), "invalid ID format: abc")

	})

	t.Run("not found", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(mockContact.MockContactService)
		handler := NewContactHandler(mockSvc, logAdapter)

		mockSvc.On("Delete", mock.Anything, int64(99)).Return(errors.New("contato não encontrado")).Once()

		req := newRequestWithVars("DELETE", "/contacts/99", nil, map[string]string{"id": "99"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.Contains(t, w.Body.String(), "contato não encontrado")

		mockSvc.AssertExpectations(t)
	})
}
