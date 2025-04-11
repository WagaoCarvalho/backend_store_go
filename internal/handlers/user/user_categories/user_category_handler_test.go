package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	handlers "github.com/WagaoCarvalho/backend_store_go/internal/handlers/user/user_categories"
	user_categories "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_categories"
	"github.com/WagaoCarvalho/backend_store_go/utils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserCategoryService struct {
	mock.Mock
}

func (m *MockUserCategoryService) GetAll(ctx context.Context) ([]user_categories.UserCategory, error) {
	args := m.Called(ctx)
	return args.Get(0).([]user_categories.UserCategory), args.Error(1)
}

func (m *MockUserCategoryService) GetById(ctx context.Context, id int64) (user_categories.UserCategory, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(user_categories.UserCategory), args.Error(1)
}

func (m *MockUserCategoryService) Create(ctx context.Context, cat user_categories.UserCategory) (user_categories.UserCategory, error) {
	args := m.Called(ctx, cat)
	return args.Get(0).(user_categories.UserCategory), args.Error(1)
}

func (m *MockUserCategoryService) Update(ctx context.Context, cat user_categories.UserCategory) (user_categories.UserCategory, error) {
	args := m.Called(ctx, cat)
	return args.Get(0).(user_categories.UserCategory), args.Error(1)
}

func (m *MockUserCategoryService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestGetCategories_Success(t *testing.T) {
	mockSvc := new(MockUserCategoryService)
	handler := handlers.NewUserCategoryHandler(mockSvc)

	expected := []user_categories.UserCategory{{ID: 1, Name: "Categoria"}}
	mockSvc.On("GetAll", mock.Anything).Return(expected, nil)

	req := httptest.NewRequest("GET", "/categories", nil)
	w := httptest.NewRecorder()

	handler.GetCategories(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response utils.DefaultResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Convertendo o data para o tipo correto
	data, ok := response.Data.([]interface{})
	assert.True(t, ok)

	var categories []user_categories.UserCategory
	for _, item := range data {
		itemBytes, _ := json.Marshal(item)
		var cat user_categories.UserCategory
		json.Unmarshal(itemBytes, &cat)
		categories = append(categories, cat)
	}

	assert.Equal(t, expected, categories)
	assert.Equal(t, "Categorias recuperadas com sucesso", response.Message)
}

func TestGetCategoryById_Success(t *testing.T) {
	mockSvc := new(MockUserCategoryService)
	handler := handlers.NewUserCategoryHandler(mockSvc)

	expected := user_categories.UserCategory{ID: 1, Name: "Teste"}
	mockSvc.On("GetById", mock.Anything, int64(1)).Return(expected, nil)

	req := mux.SetURLVars(httptest.NewRequest("GET", "/categories/1", nil), map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.GetCategoryById(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response utils.DefaultResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Convertendo o data para o tipo correto
	itemBytes, _ := json.Marshal(response.Data)
	var result user_categories.UserCategory
	json.Unmarshal(itemBytes, &result)

	assert.Equal(t, expected, result)
	assert.Equal(t, "Categoria recuperada com sucesso", response.Message)
}

func TestGetCategoryById_InvalidID(t *testing.T) {
	handler := handlers.NewUserCategoryHandler(new(MockUserCategoryService))
	req := mux.SetURLVars(httptest.NewRequest("GET", "/categories/abc", nil), map[string]string{"id": "abc"})
	w := httptest.NewRecorder()

	handler.GetCategoryById(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp utils.DefaultResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.Status)
	assert.Contains(t, resp.Message, "ID inválido")
}

func TestGetCategoryById_NotFound(t *testing.T) {
	mockSvc := new(MockUserCategoryService)
	handler := handlers.NewUserCategoryHandler(mockSvc)

	mockSvc.On("GetById", mock.Anything, int64(999)).Return(user_categories.UserCategory{}, errors.New("categoria não encontrada"))

	req := mux.SetURLVars(httptest.NewRequest("GET", "/categories/999", nil), map[string]string{"id": "999"})
	w := httptest.NewRecorder()

	handler.GetCategoryById(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	var resp utils.DefaultResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.Status)
	assert.Equal(t, "categoria não encontrada", resp.Message)
}

func TestCreateCategory_Success(t *testing.T) {
	mockSvc := new(MockUserCategoryService)
	handler := handlers.NewUserCategoryHandler(mockSvc)

	category := user_categories.UserCategory{Name: "Nova"}
	mockSvc.On("Create", mock.Anything, category).Return(category, nil)

	body, _ := json.Marshal(category)
	req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	handler.CreateCategory(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var response utils.DefaultResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Convertendo o data para o tipo correto
	itemBytes, _ := json.Marshal(response.Data)
	var result user_categories.UserCategory
	json.Unmarshal(itemBytes, &result)

	assert.Equal(t, category.Name, result.Name)
	assert.Equal(t, "Categoria criada com sucesso", response.Message)
}

func TestCreateCategory_InvalidJSON(t *testing.T) {
	handler := handlers.NewUserCategoryHandler(new(MockUserCategoryService))
	req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer([]byte("{invalid")))
	w := httptest.NewRecorder()

	handler.CreateCategory(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp utils.DefaultResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.Status)
	assert.Contains(t, resp.Message, "erro ao decodificar JSON")
}

func TestUpdateCategory_Success(t *testing.T) {
	mockSvc := new(MockUserCategoryService)
	handler := handlers.NewUserCategoryHandler(mockSvc)

	category := user_categories.UserCategory{ID: 1, Name: "Atualizada"}
	mockSvc.On("Update", mock.Anything, category).Return(category, nil)

	body, _ := json.Marshal(category)
	req := mux.SetURLVars(httptest.NewRequest("PUT", "/categories/1", bytes.NewBuffer(body)), map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.UpdateCategory(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response utils.DefaultResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)

	// Convertendo o data para o tipo correto
	itemBytes, _ := json.Marshal(response.Data)
	var result user_categories.UserCategory
	json.Unmarshal(itemBytes, &result)

	assert.Equal(t, category.Name, result.Name)
	assert.Equal(t, "Categoria atualizada com sucesso", response.Message)
}

func TestUpdateCategory_InvalidID(t *testing.T) {
	handler := handlers.NewUserCategoryHandler(new(MockUserCategoryService))

	req := mux.SetURLVars(httptest.NewRequest("PUT", "/categories/abc", bytes.NewBuffer([]byte("{}"))), map[string]string{"id": "abc"})
	w := httptest.NewRecorder()

	handler.UpdateCategory(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp utils.DefaultResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.Status)
	assert.Contains(t, resp.Message, "ID inválido")
}

func TestUpdateCategory_InvalidJSON(t *testing.T) {
	handler := handlers.NewUserCategoryHandler(new(MockUserCategoryService))

	req := mux.SetURLVars(httptest.NewRequest("PUT", "/categories/1", bytes.NewBuffer([]byte("{invalid"))), map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.UpdateCategory(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp utils.DefaultResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.Status)
	assert.Contains(t, resp.Message, "erro ao decodificar JSON")
}

func TestDeleteCategoryById_Success(t *testing.T) {
	mockSvc := new(MockUserCategoryService)
	handler := handlers.NewUserCategoryHandler(mockSvc)

	mockSvc.On("Delete", mock.Anything, int64(1)).Return(nil)

	req := mux.SetURLVars(httptest.NewRequest("DELETE", "/categories/1", nil), map[string]string{"id": "1"})
	w := httptest.NewRecorder()

	handler.DeleteCategoryById(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, 0, w.Body.Len())
}

func TestDeleteCategoryById_InvalidID(t *testing.T) {
	handler := handlers.NewUserCategoryHandler(new(MockUserCategoryService))

	req := mux.SetURLVars(httptest.NewRequest("DELETE", "/categories/abc", nil), map[string]string{"id": "abc"})
	w := httptest.NewRecorder()

	handler.DeleteCategoryById(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var resp utils.DefaultResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.Status)
	assert.Contains(t, resp.Message, "ID inválido")
}

func TestDeleteCategoryById_Error(t *testing.T) {
	mockSvc := new(MockUserCategoryService)
	handler := handlers.NewUserCategoryHandler(mockSvc)

	mockSvc.On("Delete", mock.Anything, int64(10)).Return(errors.New("erro ao deletar"))

	req := mux.SetURLVars(httptest.NewRequest("DELETE", "/categories/10", nil), map[string]string{"id": "10"})
	w := httptest.NewRecorder()

	handler.DeleteCategoryById(w, req)
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var resp utils.DefaultResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.Status)
	assert.Equal(t, "erro ao deletar", resp.Message)
}
