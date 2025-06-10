package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	model_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/contacts"
	contact_services_mock "github.com/WagaoCarvalho/backend_store_go/internal/services/contacts/contact_services_mock"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Helpers
func newRequestWithVars(method, url string, body []byte, vars map[string]string) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	return mux.SetURLVars(req, vars)
}

func TestContactHandler_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc)

		cont := &model_contact.Contact{ContactName: "Fulano"}
		mockSvc.On("Create", mock.Anything, cont).Return(cont, nil)

		body, _ := json.Marshal(cont)
		req := httptest.NewRequest("POST", "/contacts", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		handler := NewContactHandler(new(contact_services_mock.MockContactService))

		req := httptest.NewRequest("POST", "/contacts", bytes.NewBuffer([]byte("{invalid")))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc)

		cont := &model_contact.Contact{ContactName: "Erro"}
		mockSvc.On("Create", mock.Anything, cont).Return(&model_contact.Contact{}, errors.New("erro interno"))

		body, _ := json.Marshal(cont)
		req := httptest.NewRequest("POST", "/contacts", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("EmptyContactName", func(t *testing.T) {
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc)

		cont := &model_contact.Contact{ContactName: ""}
		mockSvc.On("Create", mock.Anything, cont).Return(&model_contact.Contact{}, errors.New("nome obrigatório"))

		body, _ := json.Marshal(cont)
		req := httptest.NewRequest("POST", "/contacts", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestContactHandler_GetByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc)

		mockSvc.On("GetByID", mock.Anything, int64(1)).Return(&model_contact.Contact{
			ID:          1,
			ContactName: "Contato",
		}, nil)

		req := newRequestWithVars("GET", "/contacts/1", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		handler := NewContactHandler(new(contact_services_mock.MockContactService))

		req := newRequestWithVars("GET", "/contacts/abc", nil, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc)

		mockSvc.On("GetByID", mock.Anything, int64(2)).Return(&model_contact.Contact{}, errors.New("não encontrado"))

		req := newRequestWithVars("GET", "/contacts/2", nil, map[string]string{"id": "2"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestContactHandler_GetVersionByID(t *testing.T) {
	mockSvc := new(contact_services_mock.MockContactService)
	handler := NewContactHandler(mockSvc)

	t.Run("success", func(t *testing.T) {
		mockSvc.On("GetVersionByID", mock.Anything, int64(1)).Return(3, nil).Once()

		req := newRequestWithVars("GET", "/contacts/version/1", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetVersionByID(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)

		var resp utils.DefaultResponse
		err := json.NewDecoder(w.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Status)
		assert.Equal(t, "Versão do contato encontrada", resp.Message)
		dataMap, ok := resp.Data.(map[string]interface{})
		assert.True(t, ok)
		assert.Equal(t, float64(3), dataMap["version"])
	})

	t.Run("invalid ID format", func(t *testing.T) {
		req := newRequestWithVars("GET", "/contacts/version/abc", nil, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetVersionByID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service returns ErrInvalidID", func(t *testing.T) {
		mockSvc.On("GetVersionByID", mock.Anything, int64(2)).Return(0, services.ErrInvalidID).Once()

		req := newRequestWithVars("GET", "/contacts/version/2", nil, map[string]string{"id": "2"})
		w := httptest.NewRecorder()

		handler.GetVersionByID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("contact not found", func(t *testing.T) {
		mockSvc.On("GetVersionByID", mock.Anything, int64(3)).Return(0, services.ErrContactNotFound).Once()

		req := newRequestWithVars("GET", "/contacts/version/3", nil, map[string]string{"id": "3"})
		w := httptest.NewRecorder()

		handler.GetVersionByID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		mockSvc.On("GetVersionByID", mock.Anything, int64(4)).Return(0, errors.New("erro inesperado")).Once()

		req := newRequestWithVars("GET", "/contacts/version/4", nil, map[string]string{"id": "4"})
		w := httptest.NewRecorder()

		handler.GetVersionByID(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestContactHandler_GetByUser(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc)

		mockSvc.On("GetByUser", mock.Anything, int64(1)).Return([]*model_contact.Contact{}, nil)

		req := newRequestWithVars("GET", "/contacts/user/1", nil, map[string]string{"userID": "1"})
		w := httptest.NewRecorder()

		handler.GetByUser(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc)

		mockSvc.On("GetByUser", mock.Anything, int64(2)).Return([]*model_contact.Contact{}, errors.New("erro ao buscar"))

		req := newRequestWithVars("GET", "/contacts/user/2", nil, map[string]string{"userID": "2"})
		w := httptest.NewRecorder()

		handler.GetByUser(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		handler := NewContactHandler(new(contact_services_mock.MockContactService))

		req := newRequestWithVars("GET", "/contacts/user/abc", nil, map[string]string{"userID": "abc"})
		w := httptest.NewRecorder()

		handler.GetByUser(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestContactHandler_GetByClient(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc)

		mockSvc.On("GetByClient", mock.Anything, int64(10)).Return([]*model_contact.Contact{}, nil)

		req := newRequestWithVars("GET", "/contacts/client/10", nil, map[string]string{"clientID": "10"})
		w := httptest.NewRecorder()

		handler.GetByClient(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc)

		req := newRequestWithVars("GET", "/contacts/client/1", nil, map[string]string{"clientID": "1"})
		w := httptest.NewRecorder()

		// Corrigido: retorno com tipo correto
		mockSvc.On("GetByClient", mock.Anything, int64(1)).Return([]*model_contact.Contact(nil), errors.New("erro ao buscar contatos"))

		handler.GetByClient(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "erro ao buscar contatos")
		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidClientID", func(t *testing.T) {
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc)

		req := newRequestWithVars("GET", "/contacts/client/abc", nil, map[string]string{"clientID": "abc"})
		w := httptest.NewRecorder()

		handler.GetByClient(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestContactHandler_GetBySupplier(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockSvc)

		mockSvc.On("GetBySupplier", mock.Anything, int64(5)).Return([]*model_contact.Contact{}, nil)

		req := newRequestWithVars("GET", "/contacts/supplier/5", nil, map[string]string{"supplierID": "5"})
		w := httptest.NewRecorder()

		handler.GetBySupplier(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(contact_services_mock.MockContactService)
		handler := NewContactHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/contacts/supplier/1", nil)
		req = mux.SetURLVars(req, map[string]string{"supplierID": "1"})
		rr := httptest.NewRecorder()

		var contatos []*model_contact.Contact = nil
		mockService.
			On("GetBySupplier", mock.Anything, int64(1)).
			Return(contatos, errors.New("erro ao buscar contatos do fornecedor"))

		handler.GetBySupplier(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "erro ao buscar contatos do fornecedor")
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		handler := NewContactHandler(new(contact_services_mock.MockContactService))

		req := newRequestWithVars("GET", "/contacts/supplier/abc", nil, map[string]string{"supplierID": "abc"})
		w := httptest.NewRecorder()

		handler.GetBySupplier(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestContactHandler_Update(t *testing.T) {
	mockSvc := new(contact_services_mock.MockContactService)
	handler := NewContactHandler(mockSvc)

	t.Run("success", func(t *testing.T) {
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
		req := newRequestWithVars("PUT", "/contacts/1", []byte("{invalid"), map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
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
		req := newRequestWithVars("PUT", "/contacts/abc", []byte("{}"), map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestContactHandler_Delete(t *testing.T) {
	mockSvc := new(contact_services_mock.MockContactService)
	handler := NewContactHandler(mockSvc)

	t.Run("success", func(t *testing.T) {
		mockSvc.On("Delete", mock.Anything, int64(1)).Return(nil).Once()

		req := newRequestWithVars("DELETE", "/contacts/1", nil, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("invalid ID", func(t *testing.T) {
		req := newRequestWithVars("DELETE", "/contacts/abc", nil, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		mockSvc.On("Delete", mock.Anything, int64(99)).Return(errors.New("não encontrado")).Once()

		req := newRequestWithVars("DELETE", "/contacts/99", nil, map[string]string{"id": "99"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockSvc.AssertExpectations(t)
	})
}
