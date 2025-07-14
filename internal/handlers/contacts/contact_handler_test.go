package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	model_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
	contact_services_mock "github.com/WagaoCarvalho/backend_store_go/internal/services/contacts/contact_services_mock"
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
	logAdapter := logger.NewLoggerAdapter(logrus.New())

	t.Run("Success", func(t *testing.T) {
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc, logAdapter)

		cont := &model_contact.Contact{ContactName: "Fulano"}
		mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(c *model_contact.Contact) bool {
			return c.ContactName == "Fulano"
		})).Return(cont, nil)

		body, _ := json.Marshal(cont)
		req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		type createResponse struct {
			Status  int                   `json:"status"`
			Message string                `json:"message"`
			Data    model_contact.Contact `json:"data"`
		}

		var response createResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Fulano", response.Data.ContactName)

		mockSvc.AssertExpectations(t)
	})

	t.Run("ForeignKeyInvalid", func(t *testing.T) {
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc, logAdapter)

		cont := &model_contact.Contact{ContactName: "Contato FK Inválida"}
		mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(c *model_contact.Contact) bool {
			return c.ContactName == "Contato FK Inválida"
		})).Return((*model_contact.Contact)(nil), repositories.ErrInvalidForeignKey)

		body, _ := json.Marshal(cont)
		req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc, logAdapter)

		req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBuffer([]byte("{invalid")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc, logAdapter)

		cont := &model_contact.Contact{ContactName: "Erro"}
		mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(c *model_contact.Contact) bool {
			return c.ContactName == "Erro"
		})).Return(&model_contact.Contact{}, errors.New("erro interno"))

		body, _ := json.Marshal(cont)
		req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockSvc.AssertExpectations(t)
	})

	t.Run("ServiceInternalError", func(t *testing.T) {
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc, logAdapter)

		cont := &model_contact.Contact{ContactName: "Erro Interno"}
		mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(c *model_contact.Contact) bool {
			return c.ContactName == "Erro Interno"
		})).Return((*model_contact.Contact)(nil), errors.New("erro inesperado"))

		body, _ := json.Marshal(cont)
		req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("EmptyContactName", func(t *testing.T) {
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc, logAdapter)

		cont := &model_contact.Contact{ContactName: ""}
		mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(c *model_contact.Contact) bool {
			return c.ContactName == ""
		})).Return(&model_contact.Contact{}, errors.New("nome obrigatório"))

		body, _ := json.Marshal(cont)
		req := httptest.NewRequest(http.MethodPost, "/contacts", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Create(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		mockSvc.AssertExpectations(t)
	})
}

func TestContactHandler_GetByID(t *testing.T) {
	logAdapter := logger.NewLoggerAdapter(logrus.New())

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc, logAdapter)

		expectedContact := &model_contact.Contact{
			ID:          1,
			ContactName: "Contato",
		}
		mockSvc.On("GetByID", mock.Anything, int64(1)).Return(expectedContact, nil)

		req := newRequestWithVars("GET", "/contacts/1", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp model_contact.Contact
		err := json.NewDecoder(w.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.Equal(t, expectedContact.ID, resp.ID)
		assert.Equal(t, expectedContact.ContactName, resp.ContactName)

		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		t.Parallel()
		handler := NewContactHandler(new(contact_services_mock.MockContactService), logAdapter)

		req := newRequestWithVars("GET", "/contacts/abc", nil, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc, logAdapter)

		mockSvc.On("GetByID", mock.Anything, int64(2)).Return((*model_contact.Contact)(nil), errors.New("não encontrado"))

		req := newRequestWithVars("GET", "/contacts/2", nil, map[string]string{"id": "2"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		mockSvc.AssertExpectations(t)
	})
}

func TestContactHandler_GetByUserID(t *testing.T) {
	logAdapter := logger.NewLoggerAdapter(logrus.New())

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc, logAdapter)

		mockSvc.On("GetByUserID", mock.Anything, int64(1)).Return([]*model_contact.Contact{}, nil)

		req := newRequestWithVars("GET", "/contacts/user/1", nil, map[string]string{"userID": "1"})
		w := httptest.NewRecorder()

		handler.GetByUserID(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc, logAdapter)

		mockSvc.On("GetByUserID", mock.Anything, int64(2)).Return([]*model_contact.Contact{}, errors.New("erro ao buscar"))

		req := newRequestWithVars("GET", "/contacts/user/2", nil, map[string]string{"userID": "2"})
		w := httptest.NewRecorder()

		handler.GetByUserID(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		t.Parallel()
		handler := NewContactHandler(new(contact_services_mock.MockContactService), logAdapter)

		req := newRequestWithVars("GET", "/contacts/user/abc", nil, map[string]string{"userID": "abc"})
		w := httptest.NewRecorder()

		handler.GetByUserID(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestContactHandler_GetByClientID(t *testing.T) {
	logAdapter := logger.NewLoggerAdapter(logrus.New())

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc, logAdapter)

		mockSvc.On("GetByClientID", mock.Anything, int64(10)).Return([]*model_contact.Contact{}, nil)

		req := newRequestWithVars("GET", "/contacts/client/10", nil, map[string]string{"clientID": "10"})
		w := httptest.NewRecorder()

		handler.GetByClientID(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc, logAdapter)

		mockSvc.On("GetByClientID", mock.Anything, int64(1)).Return([]*model_contact.Contact(nil), errors.New("erro ao buscar contatos"))

		req := newRequestWithVars("GET", "/contacts/client/1", nil, map[string]string{"clientID": "1"})
		w := httptest.NewRecorder()

		handler.GetByClientID(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "erro ao buscar contatos")
		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidClientID", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc, logAdapter)

		req := newRequestWithVars("GET", "/contacts/client/abc", nil, map[string]string{"clientID": "abc"})
		w := httptest.NewRecorder()

		handler.GetByClientID(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestContactHandler_GetBySupplierID(t *testing.T) {
	logAdapter := logger.NewLoggerAdapter(logrus.New())

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc, logAdapter)

		mockSvc.On("GetBySupplierID", mock.Anything, int64(5)).Return([]*model_contact.Contact{}, nil)

		req := newRequestWithVars("GET", "/contacts/supplier/5", nil, map[string]string{"supplierID": "5"})
		w := httptest.NewRecorder()

		handler.GetBySupplierID(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc, logAdapter)

		mockSvc.On("GetBySupplierID", mock.Anything, int64(1)).Return(([]*model_contact.Contact)(nil), errors.New("erro ao buscar contatos do fornecedor"))

		req := newRequestWithVars("GET", "/contacts/supplier/1", nil, map[string]string{"supplierID": "1"})
		w := httptest.NewRecorder()

		handler.GetBySupplierID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "erro ao buscar contatos do fornecedor")
		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc, logAdapter)

		req := newRequestWithVars("GET", "/contacts/supplier/abc", nil, map[string]string{"supplierID": "abc"})
		w := httptest.NewRecorder()

		handler.GetBySupplierID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestContactHandler_Update(t *testing.T) {
	logAdapter := logger.NewLoggerAdapter(logrus.New())

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc, logAdapter)

		cont := model_contact.Contact{ID: 1, ContactName: "Atualizado"}
		mockSvc.On("Update", mock.Anything, &cont).Return(nil).Once()

		body, _ := json.Marshal(cont)
		req := newRequestWithVars("PUT", "/contacts/1", body, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc, logAdapter)

		req := newRequestWithVars("PUT", "/contacts/1", []byte("{invalid"), map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc, logAdapter)

		cont := model_contact.Contact{ID: 1, ContactName: "Falha"}
		mockSvc.On("Update", mock.Anything, &cont).Return(errors.New("erro ao atualizar")).Once()

		body, _ := json.Marshal(cont)
		req := newRequestWithVars("PUT", "/contacts/1", body, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc, logAdapter)

		req := newRequestWithVars("PUT", "/contacts/abc", []byte("{}"), map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestContactHandler_Delete(t *testing.T) {
	logAdapter := logger.NewLoggerAdapter(logrus.New())

	t.Run("success", func(t *testing.T) {
		t.Parallel()
		mockSvc := new(contact_services_mock.MockContactService)
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
		mockSvc := new(contact_services_mock.MockContactService)
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
		mockSvc := new(contact_services_mock.MockContactService)
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
