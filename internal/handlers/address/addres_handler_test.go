package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	models "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
)

// Mock do serviço
type MockAddressService struct {
	mock.Mock
}

func (m *MockAddressService) CreateAddress(ctx context.Context, address models.Address) (models.Address, error) {
	args := m.Called(ctx, address)
	return args.Get(0).(models.Address), args.Error(1)
}

func (m *MockAddressService) GetAddressByID(ctx context.Context, id int) (models.Address, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Address), args.Error(1)
}

func (m *MockAddressService) UpdateAddress(ctx context.Context, address models.Address) error {
	args := m.Called(ctx, address)
	return args.Error(0)
}

func (m *MockAddressService) DeleteAddress(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestCreateAddressHandler_Success(t *testing.T) {
	mockService := new(MockAddressService)
	handler := NewAddressHandler(mockService)

	input := models.Address{
		Street: "Rua Exemplo", City: "Cidade", State: "Estado", PostalCode: "12345",
	}
	expected := input
	expected.ID = 1

	mockService.On("CreateAddress", mock.Anything, input).Return(expected, nil)

	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/addresses", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateAddress(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Endereço criado com sucesso", response["message"])
	mockService.AssertExpectations(t)
}

func TestCreateAddressHandler_InvalidJSON(t *testing.T) {
	mockService := new(MockAddressService)
	handler := NewAddressHandler(mockService)

	req := httptest.NewRequest(http.MethodPost, "/addresses", bytes.NewBuffer([]byte(`{invalid`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateAddress(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateAddressHandler_ServiceError(t *testing.T) {
	mockService := new(MockAddressService)
	handler := NewAddressHandler(mockService)

	input := models.Address{
		Street: "Rua Falha", City: "ErroCity", State: "Estado", PostalCode: "00000",
	}

	mockService.On("CreateAddress", mock.Anything, input).Return(models.Address{}, assert.AnError)

	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/addresses", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.CreateAddress(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestGetAddressHandler_Success(t *testing.T) {
	mockService := new(MockAddressService)
	handler := NewAddressHandler(mockService)

	expected := models.Address{ID: 1, Street: "Rua", City: "Cidade"}

	mockService.On("GetAddressByID", mock.Anything, 1).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/addresses/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.GetAddress(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestGetAddressHandler_ServiceError(t *testing.T) {
	mockService := new(MockAddressService)
	handler := NewAddressHandler(mockService)

	mockService.On("GetAddressByID", mock.Anything, 1).Return(models.Address{}, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/addresses/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.GetAddress(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	mockService.AssertExpectations(t)
}

func TestGetAddressHandler_InvalidID(t *testing.T) {
	mockService := new(MockAddressService)
	handler := NewAddressHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/addresses/abc", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	w := httptest.NewRecorder()

	handler.GetAddress(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestUpdateAddressHandler_Success(t *testing.T) {
	mockService := new(MockAddressService)
	handler := NewAddressHandler(mockService)

	input := models.Address{Street: "Nova Rua"}
	inputWithID := input
	inputWithID.ID = 2

	mockService.On("UpdateAddress", mock.Anything, inputWithID).Return(nil)

	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPut, "/addresses/2", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": "2"})
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.UpdateAddress(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestUpdateAddressHandler_ServiceError(t *testing.T) {
	mockService := new(MockAddressService)
	handler := NewAddressHandler(mockService)

	input := models.Address{Street: "Nova Rua"}
	inputWithID := input
	inputWithID.ID = 2

	mockService.On("UpdateAddress", mock.Anything, inputWithID).Return(assert.AnError)

	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPut, "/addresses/2", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": "2"})
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.UpdateAddress(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestUpdateAddressHandler_InvalidJSON(t *testing.T) {
	mockService := new(MockAddressService)
	handler := NewAddressHandler(mockService)

	req := httptest.NewRequest(http.MethodPut, "/addresses/2", bytes.NewBuffer([]byte(`{invalid-json}`)))
	req = mux.SetURLVars(req, map[string]string{"id": "2"})
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.UpdateAddress(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestDeleteAddressHandler_Success(t *testing.T) {
	mockService := new(MockAddressService)
	handler := NewAddressHandler(mockService)

	mockService.On("DeleteAddress", mock.Anything, 1).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/addresses/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.DeleteAddress(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestDeleteAddressHandler_InvalidID(t *testing.T) {
	mockService := new(MockAddressService)
	handler := NewAddressHandler(mockService)

	req := httptest.NewRequest(http.MethodDelete, "/addresses/abc", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	w := httptest.NewRecorder()

	handler.DeleteAddress(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestDeleteAddressHandler_ServiceError(t *testing.T) {
	mockService := new(MockAddressService)
	handler := NewAddressHandler(mockService)

	mockService.On("DeleteAddress", mock.Anything, 1).Return(assert.AnError)

	req := httptest.NewRequest(http.MethodDelete, "/addresses/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.DeleteAddress(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	mockService.AssertExpectations(t)
}
