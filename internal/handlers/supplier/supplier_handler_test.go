package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	models_supplier "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier"
	models_supplier_realiations "github.com/WagaoCarvalho/backend_store_go/internal/models/supplier/supplier_category_relations"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/suppliers"
	"github.com/WagaoCarvalho/backend_store_go/utils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newRequestWithVars(method, url string, body []byte, vars map[string]string) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	return mux.SetURLVars(req, vars)
}

func TestCreate_ValidationError(t *testing.T) {
	mockSvc := new(MockSupplierService)
	handler := NewSupplierHandler(mockSvc)

	input := models_supplier.Supplier{}
	categoryID := int64(1)

	mockSvc.On("Create", mock.Anything, &input, categoryID, mock.Anything, mock.Anything).
		Return(int64(0), errors.New("nome do fornecedor é obrigatório"))

	body, _ := json.Marshal(map[string]interface{}{
		"supplier":    input,
		"category_id": categoryID,
	})
	req := httptest.NewRequest(http.MethodPost, "/suppliers/with-category", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Create(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCreate_InvalidSupplierData(t *testing.T) {
	mockSvc := new(MockSupplierService)
	handler := NewSupplierHandler(mockSvc)

	input := models_supplier.Supplier{}
	categoryID := int64(1)

	mockSvc.On("Create", mock.Anything, &input, categoryID, (*models_address.Address)(nil), (*models_contact.Contact)(nil)).
		Return(int64(0), errors.New("fornecedor inválido"))

	requestBody := map[string]interface{}{
		"supplier":    input,
		"category_id": categoryID,
	}
	body, _ := json.Marshal(requestBody)
	req := httptest.NewRequest(http.MethodPost, "/suppliers/with-category", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Create(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	mockSvc.AssertExpectations(t)
}

func TestCreate_Success(t *testing.T) {
	mockSvc := new(MockSupplierService)
	handler := NewSupplierHandler(mockSvc)

	input := models_supplier.Supplier{Name: "Fornecedor X"}
	categoryID := int64(1)
	expectedID := int64(42)
	address := &models_address.Address{Street: "Rua A"}
	contact := &models_contact.Contact{ContactName: "Fulano"}

	// Atualize o corpo da requisição
	body, _ := json.Marshal(map[string]interface{}{
		"supplier":    input,
		"category_id": categoryID,
		"address":     address,
		"contact":     contact,
	})

	mockSvc.On("Create", mock.Anything, &input, categoryID, address, contact).
		Return(expectedID, nil)

	req := httptest.NewRequest(http.MethodPost, "/suppliers/with-category", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.Create(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	var resp utils.DefaultResponse
	err := json.NewDecoder(w.Body).Decode(&resp)
	assert.NoError(t, err)
	assert.Equal(t, "Fornecedor com categoria criado com sucesso", resp.Message)
	assert.Equal(t, float64(expectedID), resp.Data.(map[string]interface{})["id"].(float64))

	mockSvc.AssertExpectations(t)
}

func TestCreate_InvalidJSON(t *testing.T) {
	handler := NewSupplierHandler(new(MockSupplierService))

	req := httptest.NewRequest(http.MethodPost, "/suppliers/with-category", bytes.NewBufferString("invalid"))
	w := httptest.NewRecorder()

	handler.Create(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreate_CreateSupplierError(t *testing.T) {
	mockRepo := new(MockSupplierRepo)
	mockRelationService := new(MockSupplierCategoryRelationService)
	mockAddressService := new(MockAddressService)
	mockContactService := new(MockContactService)
	service := services.NewSupplierService(mockRepo, mockRelationService, mockAddressService, mockContactService)

	input := &models_supplier.Supplier{Name: "Fornecedor Y"}
	categoryID := int64(1)

	mockRepo.On("Create", mock.Anything, input).Return(int64(0), fmt.Errorf("erro ao criar fornecedor"))

	resultID, err := service.Create(context.Background(), input, categoryID, nil, nil)

	assert.Error(t, err)
	assert.Equal(t, int64(0), resultID)

	mockRepo.AssertExpectations(t)
	mockRelationService.AssertNotCalled(t, "Create")
}

func TestCreate_RelationExistsError(t *testing.T) {
	mockRepo := new(MockSupplierRepo)
	mockRelationService := new(MockSupplierCategoryRelationService)
	mockAddressService := new(MockAddressService)
	mockContactService := new(MockContactService)
	service := services.NewSupplierService(mockRepo, mockRelationService, mockAddressService, mockContactService)

	input := &models_supplier.Supplier{Name: "Fornecedor Z"}
	categoryID := int64(1)

	mockRepo.On("Create", mock.Anything, input).Return(int64(1), nil)
	mockRelationService.On("HasRelation", mock.Anything, int64(1), categoryID).Return(true, nil)

	resultID, err := service.Create(context.Background(), input, categoryID, nil, nil)

	assert.Error(t, err)
	assert.Equal(t, int64(0), resultID)

	mockRepo.AssertExpectations(t)
	mockRelationService.AssertExpectations(t)
}

func TestCreate_CreateRelationError(t *testing.T) {
	mockRepo := new(MockSupplierRepo)
	mockRelationService := new(MockSupplierCategoryRelationService)
	mockAddressService := new(MockAddressService)
	mockContactService := new(MockContactService)
	service := services.NewSupplierService(mockRepo, mockRelationService, mockAddressService, mockContactService)

	input := &models_supplier.Supplier{Name: "Fornecedor A"}
	categoryID := int64(1)

	mockRepo.On("Create", mock.Anything, input).Return(int64(1), nil)
	mockRelationService.On("HasRelation", mock.Anything, int64(1), categoryID).Return(false, nil)
	mockRelationService.On("Create", mock.Anything, int64(1), categoryID).Return(
		&models_supplier_realiations.SupplierCategoryRelations{},
		fmt.Errorf("erro ao criar relação"),
	)

	resultID, err := service.Create(context.Background(), input, categoryID, nil, nil)

	assert.Error(t, err)
	assert.Equal(t, int64(0), resultID)

	mockRepo.AssertExpectations(t)
	mockRelationService.AssertExpectations(t)
}

func TestCreate_InvalidCategoryID(t *testing.T) {
	mockRepo := new(MockSupplierRepo)
	mockRelationService := new(MockSupplierCategoryRelationService)
	mockAddressService := new(MockAddressService)
	mockContactService := new(MockContactService)
	service := services.NewSupplierService(mockRepo, mockRelationService, mockAddressService, mockContactService)

	input := &models_supplier.Supplier{Name: "Fornecedor B"}
	categoryID := int64(0)

	mockRepo.On("Create", mock.Anything, input).Return(int64(0), errors.New("categoria inválida"))

	resultID, err := service.Create(context.Background(), input, categoryID, nil, nil)

	assert.Error(t, err)
	assert.Equal(t, int64(0), resultID)

	mockRelationService.AssertNotCalled(t, "Create")
	mockRelationService.AssertNotCalled(t, "HasRelation")
}

func TestGetSupplierByID_Success(t *testing.T) {
	mockSvc := new(MockSupplierService)
	handler := NewSupplierHandler(mockSvc)

	expected := &models_supplier.Supplier{ID: 1, Name: "Fornecedor"}
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

	mockSvc.On("GetByID", mock.Anything, int64(999)).Return((*models_supplier.Supplier)(nil), errors.New("não encontrado"))

	req := newRequestWithVars("GET", "/suppliers/999", nil, map[string]string{"id": "999"})
	w := httptest.NewRecorder()

	handler.GetSupplierByID(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetAllSuppliers_Success(t *testing.T) {
	mockSvc := new(MockSupplierService)
	handler := NewSupplierHandler(mockSvc)

	mockSvc.On("GetAll", mock.Anything).Return([]*models_supplier.Supplier{{ID: 1}}, nil)

	req := httptest.NewRequest("GET", "/suppliers", nil)
	w := httptest.NewRecorder()

	handler.GetAllSuppliers(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetAllSuppliers_Error(t *testing.T) {
	mockSvc := new(MockSupplierService)
	handler := NewSupplierHandler(mockSvc)

	mockSvc.On("GetAll", mock.Anything).Return([]*models_supplier.Supplier(nil), errors.New("erro de banco"))

	req := httptest.NewRequest("GET", "/suppliers", nil)
	w := httptest.NewRecorder()

	handler.GetAllSuppliers(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateSupplier_Success(t *testing.T) {
	mockSvc := new(MockSupplierService)
	handler := NewSupplierHandler(mockSvc)

	input := &models_supplier.Supplier{ID: 1, Name: "Atualizado"}
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

	input := &models_supplier.Supplier{ID: 1, Name: "Erro"}
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
