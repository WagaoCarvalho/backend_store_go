package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

func (m *MockAddressService) Create(ctx context.Context, address models.Address) (models.Address, error) {
	args := m.Called(ctx, address)
	return args.Get(0).(models.Address), args.Error(1)
}

func (m *MockAddressService) GetByID(ctx context.Context, id int) (models.Address, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(models.Address), args.Error(1)
}

func (m *MockAddressService) Update(ctx context.Context, address models.Address) error {
	args := m.Called(ctx, address)
	return args.Error(0)
}

func (m *MockAddressService) Delete(ctx context.Context, id int) error {
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
	expected.ID = int64(0)

	mockService.On("Create", mock.Anything, input).Return(expected, nil)

	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/addresses", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

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

	handler.Create(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateAddressHandler_ServiceError(t *testing.T) {
	mockService := new(MockAddressService)
	handler := NewAddressHandler(mockService)

	input := models.Address{
		Street: "Rua Falha", City: "ErroCity", State: "Estado", PostalCode: "00000",
	}

	mockService.On("Create", mock.Anything, input).Return(models.Address{}, assert.AnError)

	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPost, "/addresses", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Create(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

func TestGetAddressByID_Success(t *testing.T) {
	mockService := new(MockAddressService)
	handler := NewAddressHandler(mockService)

	expected := models.Address{ID: int64(0), Street: "Rua", City: "Cidade"}

	mockService.On("GetByID", mock.Anything, 1).Return(expected, nil)

	req := httptest.NewRequest(http.MethodGet, "/addresses/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, "Endereço encontrado", response["message"])
	mockService.AssertExpectations(t)
}

func TestGetAddressByID_ServiceError(t *testing.T) {
	mockService := new(MockAddressService)
	handler := NewAddressHandler(mockService)

	mockService.On("GetByID", mock.Anything, 1).Return(models.Address{}, assert.AnError)

	req := httptest.NewRequest(http.MethodGet, "/addresses/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)

	var response map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["message"])
	mockService.AssertExpectations(t)
}

func TestGetAddressByID_InvalidID(t *testing.T) {
	mockService := new(MockAddressService)
	handler := NewAddressHandler(mockService)

	req := httptest.NewRequest(http.MethodGet, "/addresses/abc", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	w := httptest.NewRecorder()

	handler.GetByID(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	err := json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.NotEmpty(t, response["message"])
}

func TestAddressHandler_Update_InvalidIDParam(t *testing.T) {
	// mocks
	mockService := new(MockAddressService)
	handler := NewAddressHandler(mockService)

	// requisição com ID inválido
	req := httptest.NewRequest(http.MethodPut, "/addresses/abc", strings.NewReader(`{}`))
	req = mux.SetURLVars(req, map[string]string{"id": "abc"})
	rr := httptest.NewRecorder()

	handler.Update(rr, req)

	expected := `{
		"data": null,
		"message": "strconv.ParseInt: parsing \"abc\": invalid syntax",
		"status": 400
	}`

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.JSONEq(t, expected, rr.Body.String())
}

func TestUpdateAddressHandler_Success(t *testing.T) {
	mockService := new(MockAddressService)
	handler := NewAddressHandler(mockService)

	input := models.Address{Street: "Nova Rua"}
	id := int64(2)
	inputWithID := input
	inputWithID.ID = id

	mockService.On("Update", mock.Anything, inputWithID).Return(nil)

	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPut, "/addresses/2", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": "2"})
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Update(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	mockService.AssertExpectations(t)
}

func TestUpdateAddressHandler_ServiceError(t *testing.T) {
	mockService := new(MockAddressService)
	handler := NewAddressHandler(mockService)

	input := models.Address{Street: "Nova Rua"}
	id := int64(2)
	inputWithID := input
	inputWithID.ID = id

	mockService.On("Update", mock.Anything, inputWithID).Return(assert.AnError)

	body, _ := json.Marshal(input)
	req := httptest.NewRequest(http.MethodPut, "/addresses/2", bytes.NewBuffer(body))
	req = mux.SetURLVars(req, map[string]string{"id": "2"})
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.Update(w, req)

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

	handler.Update(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestDeleteAddressHandler_Success(t *testing.T) {
	mockService := new(MockAddressService)
	handler := NewAddressHandler(mockService)

	mockService.On("Delete", mock.Anything, 1).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, "/addresses/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.Delete(w, req)

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

	handler.Delete(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestDeleteAddressHandler_ServiceError(t *testing.T) {
	mockService := new(MockAddressService)
	handler := NewAddressHandler(mockService)

	mockService.On("Delete", mock.Anything, 1).Return(assert.AnError)

	req := httptest.NewRequest(http.MethodDelete, "/addresses/1", nil)
	req = mux.SetURLVars(req, map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.Delete(w, req)

	resp := w.Result()
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	mockService.AssertExpectations(t)
}
