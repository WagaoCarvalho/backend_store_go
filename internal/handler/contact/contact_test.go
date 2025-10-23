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
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
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
		handler := NewContact(mockService, logAdapter)

		inputDTO := &dtoContact.ContactDTO{
			ContactName: "Contato Teste",
			Email:       "teste@email.com",
			Phone:       "123456789",
		}

		// Converte para model para mock
		inputModel := dtoContact.ToContactModel(*inputDTO)

		mockService.On("Create", mock.Anything, inputModel).Return(inputModel, nil)

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

		assert.Equal(t, inputDTO.ContactName, response.Data.ContactName)
		assert.Equal(t, inputDTO.Email, response.Data.Email)
		assert.Equal(t, inputDTO.Phone, response.Data.Phone)

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockContact.MockContactService)
		handler := NewContact(mockService, logAdapter)

		inputDTO := &dtoContact.ContactDTO{
			ContactName: "Contato Erro",
			Email:       "erro@email.com",
			Phone:       "987654321",
		}
		inputModel := dtoContact.ToContactModel(*inputDTO)

		mockService.On("Create", mock.Anything, inputModel).Return((*models.Contact)(nil), assert.AnError)

		body, _ := json.Marshal(inputDTO)
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
		handler := NewContact(mockService, logAdapter)

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
		handler := NewContact(mockService, logAdapter)

		inputDTO := &dtoContact.ContactDTO{
			ContactName: "Contato FK",
			Email:       "fk@email.com",
			Phone:       "999999999",
		}
		inputModel := dtoContact.ToContactModel(*inputDTO)

		mockService.On("Create", mock.Anything, inputModel).Return((*models.Contact)(nil), errMsg.ErrDBInvalidForeignKey)

		body, _ := json.Marshal(inputDTO)
		req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro - InvalidData", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContact(mockService, logAdapter)

		inputDTO := &dtoContact.ContactDTO{
			ContactName: "",
			Email:       "teste@email.com",
		}

		inputModel := dtoContact.ToContactModel(*inputDTO)

		mockService.On("Create", mock.Anything, inputModel).Return(nil, errMsg.ErrInvalidData)

		body, _ := json.Marshal(inputDTO)
		req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro - Duplicate", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContact(mockService, logAdapter)

		inputDTO := &dtoContact.ContactDTO{
			ContactName: "Contato Teste",
			Email:       "teste@email.com",
		}

		inputModel := dtoContact.ToContactModel(*inputDTO)

		mockService.On("Create", mock.Anything, inputModel).Return(nil, errMsg.ErrDuplicate)

		body, _ := json.Marshal(inputDTO)
		req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusConflict, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro - NotFound", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContact(mockService, logAdapter)

		inputDTO := &dtoContact.ContactDTO{
			ContactName: "Contato Teste",
			Email:       "teste@email.com",
		}

		inputModel := dtoContact.ToContactModel(*inputDTO)

		mockService.On("Create", mock.Anything, inputModel).Return(nil, errMsg.ErrNotFound)

		body, _ := json.Marshal(inputDTO)
		req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestContactHandler_GetByID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("deve retornar contato com sucesso", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContact(mockService, logger)

		expectedID := int64(1)
		expectedModel := &models.Contact{
			ID:          expectedID,
			ContactName: "Contato Teste",
			Email:       "teste@email.com",
			Phone:       "123456789",
		}

		mockService.On("GetByID", mock.Anything, int64(1)).Return(expectedModel, nil)

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
		assert.Equal(t, float64(expectedID), dataMap["id"])
		assert.Equal(t, expectedModel.ContactName, dataMap["contact_name"])
		assert.Equal(t, expectedModel.Email, dataMap["email"])
		assert.Equal(t, expectedModel.Phone, dataMap["phone"])

		mockService.AssertExpectations(t)
	})

	t.Run("deve retornar erro 400 se o ID for inválido", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContact(mockService, logger)

		req := httptest.NewRequest(http.MethodGet, "/contacts/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("deve retornar erro 404 se o serviço retornar erro", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContact(mockService, logger)

		mockService.On("GetByID", mock.Anything, int64(1)).Return((*models.Contact)(nil), assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/contacts/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

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
		handler := NewContact(mockService, logger)

		id := int64(1)
		dto := dtoContact.ContactDTO{
			ID:          &id,
			ContactName: "Nome Teste",
			Email:       "teste@email.com",
			Phone:       "11999999999",
		}

		modelContact := dtoContact.ToContactModel(dto)
		mockService.On("Update", mock.Anything, modelContact).Return(nil)

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
		handler := NewContact(mockService, logger)

		dto := dtoContact.ContactDTO{ContactName: "Nome Teste"}
		req := httptest.NewRequest(http.MethodPut, "/contacts/abc", toReader(dto))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.Update(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ParseJSONError", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContact(mockService, logger)

		req := httptest.NewRequest(http.MethodPut, "/contacts/1", strings.NewReader("invalid-json"))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("ErrInvalidData", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContact(mockService, logger)

		id := int64(1)
		dto := dtoContact.ContactDTO{ID: &id, ContactName: ""} // falha de validação
		modelContact := dtoContact.ToContactModel(dto)
		mockService.On("Update", mock.Anything, modelContact).Return(errMsg.ErrInvalidData)

		req := httptest.NewRequest(http.MethodPut, "/contacts/1", toReader(dto))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("ErrDuplicate", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContact(mockService, logger)

		id := int64(1)
		dto := dtoContact.ContactDTO{ID: &id, ContactName: "Nome Teste"}
		modelContact := dtoContact.ToContactModel(dto)
		mockService.On("Update", mock.Anything, modelContact).Return(errMsg.ErrDuplicate)

		req := httptest.NewRequest(http.MethodPut, "/contacts/1", toReader(dto))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)
		assert.Equal(t, http.StatusConflict, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("ErrNotFound", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContact(mockService, logger)

		id := int64(1)
		dto := dtoContact.ContactDTO{ID: &id, ContactName: "Nome Teste"}
		modelContact := dtoContact.ToContactModel(dto)
		mockService.On("Update", mock.Anything, modelContact).Return(errMsg.ErrNotFound)

		req := httptest.NewRequest(http.MethodPut, "/contacts/1", toReader(dto))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)
		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("UnexpectedError", func(t *testing.T) {
		mockService := new(mockContact.MockContactService)
		handler := NewContact(mockService, logger)

		id := int64(1)
		dto := dtoContact.ContactDTO{ID: &id, ContactName: "Nome Teste"}
		modelContact := dtoContact.ToContactModel(dto)
		mockService.On("Update", mock.Anything, modelContact).Return(fmt.Errorf("erro inesperado"))

		req := httptest.NewRequest(http.MethodPut, "/contacts/1", toReader(dto))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
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
		handler := NewContact(mockSvc, logAdapter)

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
		handler := NewContact(mockSvc, logAdapter)

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
		handler := NewContact(mockSvc, logAdapter)

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
