package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock do serviço
type MockContactService struct {
	mock.Mock
}

func (m *MockContactService) CreateContact(ctx context.Context, c *contact.Contact) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockContactService) GetContactByID(ctx context.Context, id int64) (*contact.Contact, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*contact.Contact), args.Error(1)
}

func (m *MockContactService) GetContactsByUser(ctx context.Context, userID int64) ([]*contact.Contact, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*contact.Contact), args.Error(1)
}

func (m *MockContactService) GetContactsByClient(ctx context.Context, clientID int64) ([]*contact.Contact, error) {
	args := m.Called(ctx, clientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*contact.Contact), args.Error(1)
}

func (m *MockContactService) GetContactsBySupplier(ctx context.Context, supplierID int64) ([]*contact.Contact, error) {
	args := m.Called(ctx, supplierID)
	return args.Get(0).([]*contact.Contact), args.Error(1)
}

func (m *MockContactService) UpdateContact(ctx context.Context, c *contact.Contact) error {
	args := m.Called(ctx, c)
	return args.Error(0)
}

func (m *MockContactService) DeleteContact(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Helpers
func newRequestWithVars(method, url string, body []byte, vars map[string]string) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	return mux.SetURLVars(req, vars)
}

func TestCreateContact_Success(t *testing.T) {
	mockSvc := new(MockContactService)
	handler := NewContactHandler(mockSvc)

	cont := contact.Contact{ContactName: "Fulano"}
	mockSvc.On("CreateContact", mock.Anything, &cont).Return(nil)

	body, _ := json.Marshal(cont)
	req := httptest.NewRequest("POST", "/contacts", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.CreateContact(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCreateContact_InvalidJSON(t *testing.T) {
	handler := NewContactHandler(new(MockContactService))

	req := httptest.NewRequest("POST", "/contacts", bytes.NewBuffer([]byte("{invalid")))
	w := httptest.NewRecorder()

	handler.CreateContact(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateContact_ServiceError(t *testing.T) {
	mockSvc := new(MockContactService)
	handler := NewContactHandler(mockSvc)

	cont := contact.Contact{ContactName: "Erro"}
	mockSvc.On("CreateContact", mock.Anything, &cont).Return(errors.New("erro interno"))

	body, _ := json.Marshal(cont)
	req := httptest.NewRequest("POST", "/contacts", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.CreateContact(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetContactByID_Success(t *testing.T) {
	mockSvc := new(MockContactService)
	handler := NewContactHandler(mockSvc)

	mockSvc.On("GetContactByID", mock.Anything, int64(1)).Return(&contact.Contact{ID: 1, ContactName: "Contato"}, nil)

	req := newRequestWithVars("GET", "/contacts/1", nil, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.GetContactByID(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetContactByID_InvalidID(t *testing.T) {
	handler := NewContactHandler(new(MockContactService))

	req := newRequestWithVars("GET", "/contacts/abc", nil, map[string]string{"id": "abc"})
	w := httptest.NewRecorder()

	handler.GetContactByID(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetContactByID_NotFound(t *testing.T) {
	mockSvc := new(MockContactService)
	handler := NewContactHandler(mockSvc)

	mockSvc.On("GetContactByID", mock.Anything, int64(2)).Return(&contact.Contact{}, errors.New("não encontrado"))

	req := newRequestWithVars("GET", "/contacts/2", nil, map[string]string{"id": "2"})
	w := httptest.NewRecorder()

	handler.GetContactByID(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetContactsByUser_Success(t *testing.T) {
	mockSvc := new(MockContactService)
	handler := NewContactHandler(mockSvc)

	mockSvc.On("GetContactsByUser", mock.Anything, int64(1)).Return([]*contact.Contact{}, nil)

	req := newRequestWithVars("GET", "/contacts/user/1", nil, map[string]string{"userID": "1"})
	w := httptest.NewRecorder()

	handler.GetContactsByUser(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetContactsByUser_InvalidID(t *testing.T) {
	handler := NewContactHandler(new(MockContactService))

	req := newRequestWithVars("GET", "/contacts/user/abc", nil, map[string]string{"userID": "abc"})
	w := httptest.NewRecorder()

	handler.GetContactsByUser(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetContactsByClient_Success(t *testing.T) {
	mockSvc := new(MockContactService)
	handler := NewContactHandler(mockSvc)

	mockSvc.On("GetContactsByClient", mock.Anything, int64(10)).Return([]*contact.Contact{}, nil)

	req := newRequestWithVars("GET", "/contacts/client/10", nil, map[string]string{"clientID": "10"})
	w := httptest.NewRecorder()

	handler.GetContactsByClient(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetContactsByClient_InvalidID(t *testing.T) {
	handler := NewContactHandler(new(MockContactService))

	req := newRequestWithVars("GET", "/contacts/client/abc", nil, map[string]string{"clientID": "abc"})
	w := httptest.NewRecorder()

	handler.GetContactsByClient(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetContactsBySupplier_Success(t *testing.T) {
	mockSvc := new(MockContactService)
	handler := NewContactHandler(mockSvc)

	mockSvc.On("GetContactsBySupplier", mock.Anything, int64(5)).Return([]*contact.Contact{}, nil)

	req := newRequestWithVars("GET", "/contacts/supplier/5", nil, map[string]string{"supplierID": "5"})
	w := httptest.NewRecorder()

	handler.GetContactsBySupplier(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetContactsBySupplier_InvalidID(t *testing.T) {
	handler := NewContactHandler(new(MockContactService))

	req := newRequestWithVars("GET", "/contacts/supplier/abc", nil, map[string]string{"supplierID": "abc"})
	w := httptest.NewRecorder()

	handler.GetContactsBySupplier(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetContactsByUser_Error(t *testing.T) {
	mockSvc := new(MockContactService)
	handler := NewContactHandler(mockSvc)

	mockSvc.On("GetContactsByUser", mock.Anything, int64(2)).Return(nil, errors.New("erro ao buscar"))

	req := newRequestWithVars("GET", "/contacts/user/2", nil, map[string]string{"userID": "2"})
	w := httptest.NewRecorder()

	handler.GetContactsByUser(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateContact_Success(t *testing.T) {
	mockSvc := new(MockContactService)
	handler := NewContactHandler(mockSvc)

	cont := contact.Contact{ID: 1, ContactName: "Atualizado"}
	mockSvc.On("UpdateContact", mock.Anything, &cont).Return(nil)

	body, _ := json.Marshal(cont)
	req := newRequestWithVars("PUT", "/contacts/1", body, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.UpdateContact(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateContact_InvalidJSON(t *testing.T) {
	handler := NewContactHandler(new(MockContactService))
	req := newRequestWithVars("PUT", "/contacts/1", []byte("{invalid"), map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.UpdateContact(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateContact_ServiceError(t *testing.T) {
	mockSvc := new(MockContactService)
	handler := NewContactHandler(mockSvc)

	cont := contact.Contact{ID: 1, ContactName: "Falha"}
	mockSvc.On("UpdateContact", mock.Anything, &cont).Return(errors.New("erro ao atualizar"))

	body, _ := json.Marshal(cont)
	req := newRequestWithVars("PUT", "/contacts/1", body, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.UpdateContact(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestUpdateContact_InvalidID(t *testing.T) {
	handler := NewContactHandler(new(MockContactService))
	req := newRequestWithVars("PUT", "/contacts/abc", []byte("{}"), map[string]string{"id": "abc"})
	w := httptest.NewRecorder()

	handler.UpdateContact(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteContact_Success(t *testing.T) {
	mockSvc := new(MockContactService)
	handler := NewContactHandler(mockSvc)

	mockSvc.On("DeleteContact", mock.Anything, int64(1)).Return(nil)

	req := newRequestWithVars("DELETE", "/contacts/1", nil, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.DeleteContact(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestDeleteContact_InvalidID(t *testing.T) {
	handler := NewContactHandler(new(MockContactService))
	req := newRequestWithVars("DELETE", "/contacts/abc", nil, map[string]string{"id": "abc"})
	w := httptest.NewRecorder()

	handler.DeleteContact(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteContact_NotFound(t *testing.T) {
	mockSvc := new(MockContactService)
	handler := NewContactHandler(mockSvc)

	mockSvc.On("DeleteContact", mock.Anything, int64(99)).Return(errors.New("não encontrado"))

	req := newRequestWithVars("DELETE", "/contacts/99", nil, map[string]string{"id": "99"})
	w := httptest.NewRecorder()

	handler.DeleteContact(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}
