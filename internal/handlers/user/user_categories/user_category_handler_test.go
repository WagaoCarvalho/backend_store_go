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
	models_user_categories "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_categories"
	"github.com/WagaoCarvalho/backend_store_go/utils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserCategoryService struct {
	mock.Mock
}

func (m *MockUserCategoryService) GetAll(ctx context.Context) ([]*models_user_categories.UserCategory, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*models_user_categories.UserCategory), args.Error(1)
}

func (m *MockUserCategoryService) GetByID(ctx context.Context, id int64) (*models_user_categories.UserCategory, error) {
	args := m.Called(ctx, id)
	if obj := args.Get(0); obj != nil {
		return obj.(*models_user_categories.UserCategory), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserCategoryService) Create(ctx context.Context, cat *models_user_categories.UserCategory) (*models_user_categories.UserCategory, error) {
	args := m.Called(ctx, cat)
	if obj := args.Get(0); obj != nil {
		return obj.(*models_user_categories.UserCategory), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserCategoryService) Update(ctx context.Context, cat *models_user_categories.UserCategory) (*models_user_categories.UserCategory, error) {
	args := m.Called(ctx, cat)
	if obj := args.Get(0); obj != nil {
		return obj.(*models_user_categories.UserCategory), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *MockUserCategoryService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestUserCategoryHandler_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(MockUserCategoryService)
		handler := handlers.NewUserCategoryHandler(mockSvc)

		category := &models_user_categories.UserCategory{Name: "Nova"}
		mockSvc.On("Create", mock.Anything, category).Return(category, nil)

		body, _ := json.Marshal(category)
		req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Create(w, req)
		assert.Equal(t, http.StatusCreated, w.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Convertendo o data para o tipo correto
		itemBytes, _ := json.Marshal(response.Data)
		var result models_user_categories.UserCategory
		json.Unmarshal(itemBytes, &result)

		assert.Equal(t, category.Name, result.Name)
		assert.Equal(t, "Categoria criada com sucesso", response.Message)

		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		handler := handlers.NewUserCategoryHandler(new(MockUserCategoryService))

		req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer([]byte("{invalid")))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "erro ao decodificar JSON")
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockSvc := new(MockUserCategoryService)
		handler := handlers.NewUserCategoryHandler(mockSvc)

		input := models_user_categories.UserCategory{Name: "Erro"}

		mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(c *models_user_categories.UserCategory) bool {
			return c.Name == input.Name
		})).Return(nil, errors.New("erro ao criar categoria")) // ✅ CORREÇÃO AQUI

		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/categories", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Create(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)
		assert.Equal(t, "erro ao criar categoria", resp.Message)

		mockSvc.AssertExpectations(t)
	})

}

func TestUserCategoryHandler_GetAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(MockUserCategoryService)
		handler := handlers.NewUserCategoryHandler(mockSvc)

		expected := []*models_user_categories.UserCategory{{ID: 1, Name: "Categoria"}}
		mockSvc.On("GetAll", mock.Anything).Return(expected, nil)

		req := httptest.NewRequest("GET", "/categories", nil)
		w := httptest.NewRecorder()

		handler.GetAll(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Converter o campo Data para slice correto
		rawData, ok := response.Data.([]interface{})
		assert.True(t, ok)

		var result []*models_user_categories.UserCategory
		for _, item := range rawData {
			itemBytes, _ := json.Marshal(item)
			var cat models_user_categories.UserCategory
			json.Unmarshal(itemBytes, &cat)
			result = append(result, &cat)
		}

		assert.Equal(t, expected, result)
		assert.Equal(t, "Categorias recuperadas com sucesso", response.Message)

		mockSvc.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockSvc := new(MockUserCategoryService)
		handler := handlers.NewUserCategoryHandler(mockSvc)

		mockSvc.On("GetAll", mock.Anything).Return([]*models_user_categories.UserCategory{}, errors.New("erro de banco"))

		req := httptest.NewRequest("GET", "/categories", nil)
		w := httptest.NewRecorder()

		handler.GetAll(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, response.Status)
		assert.Contains(t, response.Message, "erro de banco")
		assert.Nil(t, response.Data)

		mockSvc.AssertExpectations(t)
	})

}

func TestUserCategoryHandler_GetById(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(MockUserCategoryService)
		handler := handlers.NewUserCategoryHandler(mockSvc)

		expected := &models_user_categories.UserCategory{ID: 1, Name: "Teste"}
		mockSvc.On("GetByID", mock.Anything, int64(1)).Return(expected, nil)

		req := mux.SetURLVars(httptest.NewRequest("GET", "/categories/1", nil), map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetById(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		// Convertendo o data para o tipo correto
		itemBytes, _ := json.Marshal(response.Data)
		var result models_user_categories.UserCategory
		json.Unmarshal(itemBytes, &result)

		assert.Equal(t, *expected, result)
		assert.Equal(t, "Categoria recuperada com sucesso", response.Message)

		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		handler := handlers.NewUserCategoryHandler(new(MockUserCategoryService))
		req := mux.SetURLVars(httptest.NewRequest("GET", "/categories/abc", nil), map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetById(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "ID inválido")
	})

	t.Run("NotFound", func(t *testing.T) {
		mockSvc := new(MockUserCategoryService)
		handler := handlers.NewUserCategoryHandler(mockSvc)

		mockSvc.On("GetByID", mock.Anything, int64(999)).Return(nil, errors.New("categoria não encontrada"))

		req := mux.SetURLVars(httptest.NewRequest("GET", "/categories/999", nil), map[string]string{"id": "999"})
		w := httptest.NewRecorder()

		handler.GetById(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.Status)
		assert.Equal(t, "categoria não encontrada", resp.Message)

		mockSvc.AssertExpectations(t)
	})
}

func TestUserCategoryHandler_Update(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(MockUserCategoryService)
		handler := handlers.NewUserCategoryHandler(mockSvc)

		category := &models_user_categories.UserCategory{ID: 1, Name: "Atualizada"}
		mockSvc.On("Update", mock.Anything, category).Return(category, nil)

		body, _ := json.Marshal(category)
		req := mux.SetURLVars(httptest.NewRequest("PUT", "/categories/1", bytes.NewBuffer(body)), map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		itemBytes, _ := json.Marshal(response.Data)
		var result models_user_categories.UserCategory
		json.Unmarshal(itemBytes, &result)

		assert.Equal(t, category.Name, result.Name)
		assert.Equal(t, "Categoria atualizada com sucesso", response.Message)

		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		handler := handlers.NewUserCategoryHandler(new(MockUserCategoryService))

		req := mux.SetURLVars(httptest.NewRequest("PUT", "/categories/abc", bytes.NewBuffer([]byte("{}"))), map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "ID inválido")
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		handler := handlers.NewUserCategoryHandler(new(MockUserCategoryService))

		req := mux.SetURLVars(httptest.NewRequest("PUT", "/categories/1", bytes.NewBuffer([]byte("{invalid"))), map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "erro ao decodificar JSON")
	})

	t.Run("UpdateError", func(t *testing.T) {
		mockSvc := new(MockUserCategoryService)
		handler := handlers.NewUserCategoryHandler(mockSvc)

		category := &models_user_categories.UserCategory{ID: 2, Name: "Falha"}
		mockSvc.On("Update", mock.Anything, category).Return(nil, errors.New("erro ao atualizar"))

		body, _ := json.Marshal(category)
		req := mux.SetURLVars(httptest.NewRequest("PUT", "/categories/2", bytes.NewBuffer(body)), map[string]string{"id": "2"})
		w := httptest.NewRecorder()

		handler.Update(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)
		assert.Equal(t, "erro ao atualizar", resp.Message)

		mockSvc.AssertExpectations(t)
	})
}

func TestUserCategoryHandler_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(MockUserCategoryService)
		handler := handlers.NewUserCategoryHandler(mockSvc)

		mockSvc.On("Delete", mock.Anything, int64(1)).Return(nil)

		req := mux.SetURLVars(httptest.NewRequest("DELETE", "/categories/1", nil), map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, 0, w.Body.Len())

		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		handler := handlers.NewUserCategoryHandler(new(MockUserCategoryService))

		req := mux.SetURLVars(httptest.NewRequest("DELETE", "/categories/abc", nil), map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "ID inválido")
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockSvc := new(MockUserCategoryService)
		handler := handlers.NewUserCategoryHandler(mockSvc)

		mockSvc.On("Delete", mock.Anything, int64(10)).Return(errors.New("erro ao deletar"))

		req := mux.SetURLVars(httptest.NewRequest("DELETE", "/categories/10", nil), map[string]string{"id": "10"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)
		assert.Equal(t, "erro ao deletar", resp.Message)

		mockSvc.AssertExpectations(t)
	})
}
