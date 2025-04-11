package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier"
	"github.com/WagaoCarvalho/backend_store_go/utils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSupplierService para testes
type MockSupplierService struct {
	mock.Mock
}

func (m *MockSupplierService) Create(ctx context.Context, s *models.Supplier) (int64, error) {
	args := m.Called(ctx, s)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSupplierService) GetByID(ctx context.Context, id int64) (*models.Supplier, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.Supplier), args.Error(1)
}

func (m *MockSupplierService) GetAll(ctx context.Context) ([]*models.Supplier, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.Supplier), args.Error(1)
}

func (m *MockSupplierService) Update(ctx context.Context, s *models.Supplier) error {
	args := m.Called(ctx, s)
	return args.Error(0)
}

func (m *MockSupplierService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

// Helpers
func newRequestWithVars(method, url string, body []byte, vars map[string]string) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	return mux.SetURLVars(req, vars)
}

func TestCreateSupplier_Success(t *testing.T) {
	mockSvc := new(MockSupplierService)
	handler := NewSupplierHandler(mockSvc)

	input := &models.Supplier{Name: "Fornecedor X"}
	mockSvc.On("Create", mock.Anything, input).Return(int64(1), nil)

	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.CreateSupplier(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var resp utils.DefaultResponse
	json.NewDecoder(w.Body).Decode(&resp)
	assert.Equal(t, "Fornecedor criado com sucesso", resp.Message)
}

func TestCreateSupplier_InvalidJSON(t *testing.T) {
	handler := NewSupplierHandler(new(MockSupplierService))

	req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewBufferString("invalid"))
	w := httptest.NewRecorder()

	handler.CreateSupplier(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateSupplier_ValidationError(t *testing.T) {
	mockSvc := new(MockSupplierService)
	handler := NewSupplierHandler(mockSvc)

	input := &models.Supplier{}
	mockSvc.On("Create", mock.Anything, input).Return(int64(0), errors.New("nome do fornecedor é obrigatório"))

	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/suppliers", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.CreateSupplier(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetSupplierByID_Success(t *testing.T) {
	mockSvc := new(MockSupplierService)
	handler := NewSupplierHandler(mockSvc)

	expected := &models.Supplier{ID: 1, Name: "Fornecedor"}
	mockSvc.On("GetByID", mock.Anything, int64(1)).Return(expected, nil)

	req := newRequestWithVars("GET", "/suppliers/1", nil, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.GetSupplierByID(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetSupplierByID_InvalidID(t *testing.T) {
	handler := NewSupplierHandler(new(MockSupplierService))

	req := newRequestWithVars("GET", "/suppliers/abc", nil, map[string]string{"id": "abc"})
	w := httptest.NewRecorder()

	handler.GetSupplierByID(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetSupplierByID_NotFound(t *testing.T) {
	mockSvc := new(MockSupplierService)
	handler := NewSupplierHandler(mockSvc)

	mockSvc.On("GetByID", mock.Anything, int64(999)).Return((*models.Supplier)(nil), errors.New("não encontrado"))

	req := newRequestWithVars("GET", "/suppliers/999", nil, map[string]string{"id": "999"})
	w := httptest.NewRecorder()

	handler.GetSupplierByID(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetAllSuppliers_Success(t *testing.T) {
	mockSvc := new(MockSupplierService)
	handler := NewSupplierHandler(mockSvc)

	mockSvc.On("GetAll", mock.Anything).Return([]*models.Supplier{{ID: 1}}, nil)

	req := httptest.NewRequest("GET", "/suppliers", nil)
	w := httptest.NewRecorder()

	handler.GetAllSuppliers(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetAllSuppliers_Error(t *testing.T) {
	mockSvc := new(MockSupplierService)
	handler := NewSupplierHandler(mockSvc)

	mockSvc.On("GetAll", mock.Anything).Return([]*models.Supplier(nil), errors.New("erro de banco"))

	req := httptest.NewRequest("GET", "/suppliers", nil)
	w := httptest.NewRecorder()

	handler.GetAllSuppliers(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateSupplier_Success(t *testing.T) {
	mockSvc := new(MockSupplierService)
	handler := NewSupplierHandler(mockSvc)

	input := &models.Supplier{ID: 1, Name: "Atualizado"}
	mockSvc.On("Update", mock.Anything, input).Return(nil)

	body, _ := json.Marshal(input)
	req := newRequestWithVars("PUT", "/suppliers/1", body, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.UpdateSupplier(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestUpdateSupplier_InvalidID(t *testing.T) {
	handler := NewSupplierHandler(new(MockSupplierService))

	req := newRequestWithVars("PUT", "/suppliers/abc", nil, map[string]string{"id": "abc"})
	w := httptest.NewRecorder()

	handler.UpdateSupplier(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateSupplier_InvalidJSON(t *testing.T) {
	handler := NewSupplierHandler(new(MockSupplierService))

	req := newRequestWithVars("PUT", "/suppliers/1", []byte("{invalid"), map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.UpdateSupplier(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateSupplier_Error(t *testing.T) {
	mockSvc := new(MockSupplierService)
	handler := NewSupplierHandler(mockSvc)

	input := &models.Supplier{ID: 1, Name: "Erro"}
	mockSvc.On("Update", mock.Anything, input).Return(errors.New("erro"))

	body, _ := json.Marshal(input)
	req := newRequestWithVars("PUT", "/suppliers/1", body, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.UpdateSupplier(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteSupplier_Success(t *testing.T) {
	mockSvc := new(MockSupplierService)
	handler := NewSupplierHandler(mockSvc)

	mockSvc.On("Delete", mock.Anything, int64(1)).Return(nil)

	req := newRequestWithVars("DELETE", "/suppliers/1", nil, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.DeleteSupplier(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestDeleteSupplier_InvalidID(t *testing.T) {
	handler := NewSupplierHandler(new(MockSupplierService))

	req := newRequestWithVars("DELETE", "/suppliers/abc", nil, map[string]string{"id": "abc"})
	w := httptest.NewRecorder()

	handler.DeleteSupplier(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteSupplier_Error(t *testing.T) {
	mockSvc := new(MockSupplierService)
	handler := NewSupplierHandler(mockSvc)

	mockSvc.On("Delete", mock.Anything, int64(999)).Return(errors.New("não encontrado"))

	req := newRequestWithVars("DELETE", "/suppliers/999", nil, map[string]string{"id": "999"})
	w := httptest.NewRecorder()

	handler.DeleteSupplier(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}
