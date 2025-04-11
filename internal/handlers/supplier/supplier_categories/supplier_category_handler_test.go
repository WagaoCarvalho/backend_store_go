package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_categories"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock Service
type MockSupplierCategoryService struct {
	mock.Mock
}

func (m *MockSupplierCategoryService) Create(ctx context.Context, category *models.SupplierCategory) (int64, error) {
	args := m.Called(ctx, category)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockSupplierCategoryService) GetByID(ctx context.Context, id int64) (*models.SupplierCategory, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*models.SupplierCategory), args.Error(1)
}

func (m *MockSupplierCategoryService) GetAll(ctx context.Context) ([]*models.SupplierCategory, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models.SupplierCategory), args.Error(1)
}

func (m *MockSupplierCategoryService) Update(ctx context.Context, category *models.SupplierCategory) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockSupplierCategoryService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateCategory_Success(t *testing.T) {
	mockSvc := new(MockSupplierCategoryService)
	handler := NewSupplierCategoryHandler(mockSvc)

	input := models.SupplierCategory{Name: "Eletrônicos"}
	mockSvc.On("Create", mock.Anything, &input).Return(int64(1), nil)

	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Create(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCreateCategory_InvalidJSON(t *testing.T) {
	handler := NewSupplierCategoryHandler(new(MockSupplierCategoryService))
	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer([]byte(`{invalid`)))
	w := httptest.NewRecorder()

	handler.Create(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreateCategory_ErrorFromService(t *testing.T) {
	mockSvc := new(MockSupplierCategoryService)
	handler := NewSupplierCategoryHandler(mockSvc)

	input := models.SupplierCategory{Name: ""}
	mockSvc.On("Create", mock.Anything, &input).Return(int64(0), errors.New("nome da categoria é obrigatório"))

	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Create(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetByID_Success(t *testing.T) {
	mockSvc := new(MockSupplierCategoryService)
	handler := NewSupplierCategoryHandler(mockSvc)

	expected := &models.SupplierCategory{ID: 1, Name: "Informática"}
	mockSvc.On("GetByID", mock.Anything, int64(1)).Return(expected, nil)

	req := mux.SetURLVars(httptest.NewRequest(http.MethodGet, "/categories/1", nil), map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetByID_InvalidID(t *testing.T) {
	handler := NewSupplierCategoryHandler(new(MockSupplierCategoryService))
	req := mux.SetURLVars(httptest.NewRequest(http.MethodGet, "/categories/abc", nil), map[string]string{"id": "abc"})
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetByID_NotFound(t *testing.T) {
	mockSvc := new(MockSupplierCategoryService)
	handler := NewSupplierCategoryHandler(mockSvc)

	mockSvc.On("GetByID", mock.Anything, int64(999)).Return((*models.SupplierCategory)(nil), errors.New("não encontrada"))

	req := mux.SetURLVars(httptest.NewRequest(http.MethodGet, "/categories/999", nil), map[string]string{"id": "999"})
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetAll_Success(t *testing.T) {
	mockSvc := new(MockSupplierCategoryService)
	handler := NewSupplierCategoryHandler(mockSvc)

	mockSvc.On("GetAll", mock.Anything).Return([]*models.SupplierCategory{
		{ID: 1, Name: "Insumos"},
	}, nil)

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestGetAll_Error(t *testing.T) {
	mockSvc := new(MockSupplierCategoryService)
	handler := NewSupplierCategoryHandler(mockSvc)

	mockSvc.On("GetAll", mock.Anything).Return(([]*models.SupplierCategory)(nil), errors.New("falha ao buscar"))

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	w := httptest.NewRecorder()

	handler.GetAll(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateCategory_Success(t *testing.T) {
	mockSvc := new(MockSupplierCategoryService)
	handler := NewSupplierCategoryHandler(mockSvc)

	input := models.SupplierCategory{ID: 1, Name: "Atualizado"}
	mockSvc.On("Update", mock.Anything, &input).Return(nil)

	body, _ := json.Marshal(input)
	req := mux.SetURLVars(httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewBuffer(body)), map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.Update(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestUpdateCategory_InvalidID(t *testing.T) {
	handler := NewSupplierCategoryHandler(new(MockSupplierCategoryService))

	req := mux.SetURLVars(httptest.NewRequest(http.MethodPut, "/categories/abc", bytes.NewBuffer([]byte(`{}`))), map[string]string{"id": "abc"})
	w := httptest.NewRecorder()

	handler.Update(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateCategory_InvalidJSON(t *testing.T) {
	handler := NewSupplierCategoryHandler(new(MockSupplierCategoryService))

	req := mux.SetURLVars(httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewBuffer([]byte(`{invalid`))), map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.Update(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateCategory_ServiceError(t *testing.T) {
	mockSvc := new(MockSupplierCategoryService)
	handler := NewSupplierCategoryHandler(mockSvc)

	input := models.SupplierCategory{ID: 1, Name: ""}
	mockSvc.On("Update", mock.Anything, &input).Return(errors.New("nome da categoria é obrigatório"))

	body, _ := json.Marshal(input)
	req := mux.SetURLVars(httptest.NewRequest(http.MethodPut, "/categories/1", bytes.NewBuffer(body)), map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.Update(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteCategory_Success(t *testing.T) {
	mockSvc := new(MockSupplierCategoryService)
	handler := NewSupplierCategoryHandler(mockSvc)

	mockSvc.On("Delete", mock.Anything, int64(1)).Return(nil)

	req := mux.SetURLVars(httptest.NewRequest(http.MethodDelete, "/categories/1", nil), map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestDeleteCategory_InvalidID(t *testing.T) {
	handler := NewSupplierCategoryHandler(new(MockSupplierCategoryService))

	req := mux.SetURLVars(httptest.NewRequest(http.MethodDelete, "/categories/abc", nil), map[string]string{"id": "abc"})
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDeleteCategory_ServiceError(t *testing.T) {
	mockSvc := new(MockSupplierCategoryService)
	handler := NewSupplierCategoryHandler(mockSvc)

	mockSvc.On("Delete", mock.Anything, int64(1)).Return(errors.New("erro ao deletar"))

	req := mux.SetURLVars(httptest.NewRequest(http.MethodDelete, "/categories/1", nil), map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}
