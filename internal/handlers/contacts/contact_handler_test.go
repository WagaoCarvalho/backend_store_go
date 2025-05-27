package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	model_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock do serviço
type MockContactService struct {
	mock.Mock
}

func (m *MockContactService) Create(ctx context.Context, c *model_contact.Contact) (*model_contact.Contact, error) {
	args := m.Called(ctx, c) // captura os argumentos retornados do mock
	contact, _ := args.Get(0).(*model_contact.Contact)
	return contact, args.Error(1)
}

func (m *MockContactService) GetByID(ctx context.Context, id int64) (*model_contact.Contact, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*model_contact.Contact), args.Error(1)
}

func (m *MockContactService) GetByUser(ctx context.Context, userID int64) ([]*model_contact.Contact, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model_contact.Contact), args.Error(1)
}

func (m *MockContactService) GetByClient(ctx context.Context, clientID int64) ([]*model_contact.Contact, error) {
	args := m.Called(ctx, clientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model_contact.Contact), args.Error(1)
}

func (m *MockContactService) GetBySupplier(ctx context.Context, supplierID int64) ([]*model_contact.Contact, error) {
	args := m.Called(ctx, supplierID)
	return args.Get(0).([]*model_contact.Contact), args.Error(1)
}

func (m *MockContactService) Update(ctx context.Context, c *model_contact.Contact) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockContactService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Helpers
func newRequestWithVars(method, url string, body []byte, vars map[string]string) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	return mux.SetURLVars(req, vars)
}

func TestContactHandler_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(MockContactService)
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
		handler := NewContactHandler(new(MockContactService))

		req := httptest.NewRequest("POST", "/contacts", bytes.NewBuffer([]byte("{invalid")))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockSvc := new(MockContactService)
		handler := NewContactHandler(mockSvc)

		cont := &model_contact.Contact{ContactName: "Erro"}
		mockSvc.On("Create", mock.Anything, cont).Return(model_contact.Contact{}, errors.New("erro interno"))

		body, _ := json.Marshal(cont)
		req := httptest.NewRequest("POST", "/contacts", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})

	t.Run("EmptyContactName", func(t *testing.T) {
		mockSvc := new(MockContactService)
		handler := NewContactHandler(mockSvc)

		cont := &model_contact.Contact{ContactName: ""}
		mockSvc.On("Create", mock.Anything, cont).Return(model_contact.Contact{}, errors.New("nome obrigatório"))

		body, _ := json.Marshal(cont)
		req := httptest.NewRequest("POST", "/contacts", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockSvc.AssertExpectations(t)
	})
}

func TestGetByID_Success(t *testing.T) {
	mockSvc := new(MockContactService)
	handler := NewContactHandler(mockSvc)

	mockSvc.On("GetByID", mock.Anything, int64(1)).Return(&model_contact.Contact{ID: 1, ContactName: "Contato"}, nil)

	req := newRequestWithVars("GET", "/contacts/1", nil, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetByID_InvalidID(t *testing.T) {
	handler := NewContactHandler(new(MockContactService))

	req := newRequestWithVars("GET", "/contacts/abc", nil, map[string]string{"id": "abc"})
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetByID_NotFound(t *testing.T) {
	mockSvc := new(MockContactService)
	handler := NewContactHandler(mockSvc)

	mockSvc.On("GetByID", mock.Anything, int64(2)).Return(&model_contact.Contact{}, errors.New("não encontrado"))

	req := newRequestWithVars("GET", "/contacts/2", nil, map[string]string{"id": "2"})
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetByUser_Success(t *testing.T) {
	mockSvc := new(MockContactService)
	handler := NewContactHandler(mockSvc)

	mockSvc.On("GetByUser", mock.Anything, int64(1)).Return([]*model_contact.Contact{}, nil)

	req := newRequestWithVars("GET", "/contacts/user/1", nil, map[string]string{"userID": "1"})
	w := httptest.NewRecorder()

	handler.GetByUser(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetByUser_InvalidID(t *testing.T) {
	handler := NewContactHandler(new(MockContactService))

	req := newRequestWithVars("GET", "/contacts/user/abc", nil, map[string]string{"userID": "abc"})
	w := httptest.NewRecorder()

	handler.GetByUser(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetByClient_Success(t *testing.T) {
	mockSvc := new(MockContactService)
	handler := NewContactHandler(mockSvc)

	mockSvc.On("GetByClient", mock.Anything, int64(10)).Return([]*model_contact.Contact{}, nil)

	req := newRequestWithVars("GET", "/contacts/client/10", nil, map[string]string{"clientID": "10"})
	w := httptest.NewRecorder()

	handler.GetByClient(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestContactHandler_GetByClient_ServiceError(t *testing.T) {
	mockService := new(MockContactService)
	handler := NewContactHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/contacts/client/1", nil)
	req = mux.SetURLVars(req, map[string]string{"clientID": "1"})
	rr := httptest.NewRecorder()

	// Simula erro do serviço
	mockService.On("GetByClient", mock.Anything, int64(1)).Return(nil, errors.New("erro ao buscar contatos"))

	handler.GetByClient(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "erro ao buscar contatos")
	mockService.AssertExpectations(t)
}

func TestGetByClient_InvalidID(t *testing.T) {
	handler := NewContactHandler(new(MockContactService))

	req := newRequestWithVars("GET", "/contacts/client/abc", nil, map[string]string{"clientID": "abc"})
	w := httptest.NewRecorder()

	handler.GetByClient(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetBySupplier_Success(t *testing.T) {
	mockSvc := new(MockContactService)
	handler := NewContactHandler(mockSvc)

	mockSvc.On("GetBySupplier", mock.Anything, int64(5)).Return([]*model_contact.Contact{}, nil)

	req := newRequestWithVars("GET", "/contacts/supplier/5", nil, map[string]string{"supplierID": "5"})
	w := httptest.NewRecorder()

	handler.GetBySupplier(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestContactHandler_GetBySupplier_ServiceError(t *testing.T) {
	mockService := new(MockContactService)
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
}

func TestGetBySupplier_InvalidID(t *testing.T) {
	handler := NewContactHandler(new(MockContactService))

	req := newRequestWithVars("GET", "/contacts/supplier/abc", nil, map[string]string{"supplierID": "abc"})
	w := httptest.NewRecorder()

	handler.GetBySupplier(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetByUser_Error(t *testing.T) {
	mockSvc := new(MockContactService)
	handler := NewContactHandler(mockSvc)

	mockSvc.On("GetByUser", mock.Anything, int64(2)).Return(nil, errors.New("erro ao buscar"))

	req := newRequestWithVars("GET", "/contacts/user/2", nil, map[string]string{"userID": "2"})
	w := httptest.NewRecorder()

	handler.GetByUser(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdate_Success(t *testing.T) {
	mockSvc := new(MockContactService)
	handler := NewContactHandler(mockSvc)

	cont := model_contact.Contact{ID: 1, ContactName: "Atualizado"}
	mockSvc.On("Update", mock.Anything, &cont).Return(nil)

	body, _ := json.Marshal(cont)
	req := newRequestWithVars("PUT", "/contacts/1", body, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.Update(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdate_InvalidJSON(t *testing.T) {
	handler := NewContactHandler(new(MockContactService))
	req := newRequestWithVars("PUT", "/contacts/1", []byte("{invalid"), map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.Update(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdate_ServiceError(t *testing.T) {
	mockSvc := new(MockContactService)
	handler := NewContactHandler(mockSvc)

	cont := model_contact.Contact{ID: 1, ContactName: "Falha"}
	mockSvc.On("Update", mock.Anything, &cont).Return(errors.New("erro ao atualizar"))

	body, _ := json.Marshal(cont)
	req := newRequestWithVars("PUT", "/contacts/1", body, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.Update(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestUpdate_InvalidID(t *testing.T) {
	handler := NewContactHandler(new(MockContactService))
	req := newRequestWithVars("PUT", "/contacts/abc", []byte("{}"), map[string]string{"id": "abc"})
	w := httptest.NewRecorder()

	handler.Update(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDelete_Success(t *testing.T) {
	mockSvc := new(MockContactService)
	handler := NewContactHandler(mockSvc)

	mockSvc.On("Delete", mock.Anything, int64(1)).Return(nil)

	req := newRequestWithVars("DELETE", "/contacts/1", nil, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.Delete(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDelete_InvalidID(t *testing.T) {
	handler := NewContactHandler(new(MockContactService))
	req := newRequestWithVars("DELETE", "/contacts/abc", nil, map[string]string{"id": "abc"})
	w := httptest.NewRecorder()

	handler.Delete(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDelete_NotFound(t *testing.T) {
	mockSvc := new(MockContactService)
	handler := NewContactHandler(mockSvc)

	mockSvc.On("Delete", mock.Anything, int64(99)).Return(errors.New("não encontrado"))

	req := newRequestWithVars("DELETE", "/contacts/99", nil, map[string]string{"id": "99"})
	w := httptest.NewRecorder()

	handler.Delete(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func ptrInt64(i int64) *int64 {
	return &i
}
