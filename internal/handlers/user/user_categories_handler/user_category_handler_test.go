package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	models_user_categories "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_categories"
	user_category_services_mock "github.com/WagaoCarvalho/backend_store_go/internal/services/user/user_categories/user_categories_services_mock"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func newRequestWithVars(method, url string, body []byte, vars map[string]string) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewReader(body))
	return mux.SetURLVars(req, vars)
}

func TestUserCategoryHandler_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(user_category_services_mock.MockUserCategoryService)
		handler := NewUserCategoryHandler(mockSvc)

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

		itemBytes, _ := json.Marshal(response.Data)
		var result models_user_categories.UserCategory
		json.Unmarshal(itemBytes, &result)

		assert.Equal(t, category.Name, result.Name)
		assert.Equal(t, "Categoria criada com sucesso", response.Message)

		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
		handler := NewUserCategoryHandler(new(user_category_services_mock.MockUserCategoryService))

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
		mockSvc := new(user_category_services_mock.MockUserCategoryService)
		handler := NewUserCategoryHandler(mockSvc)

		input := models_user_categories.UserCategory{Name: "Erro"}

		mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(c *models_user_categories.UserCategory) bool {
			return c.Name == input.Name
		})).Return(nil, errors.New("erro ao criar categoria"))

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
		mockSvc := new(user_category_services_mock.MockUserCategoryService)
		handler := NewUserCategoryHandler(mockSvc)

		expected := []*models_user_categories.UserCategory{{ID: 1, Name: "Categoria"}}
		mockSvc.On("GetAll", mock.Anything).Return(expected, nil)

		req := httptest.NewRequest("GET", "/categories", nil)
		w := httptest.NewRecorder()

		handler.GetAll(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

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
		mockSvc := new(user_category_services_mock.MockUserCategoryService)
		handler := NewUserCategoryHandler(mockSvc)

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
		mockSvc := new(user_category_services_mock.MockUserCategoryService)
		handler := NewUserCategoryHandler(mockSvc)

		expected := &models_user_categories.UserCategory{ID: 1, Name: "Teste"}
		mockSvc.On("GetByID", mock.Anything, int64(1)).Return(expected, nil)

		req := mux.SetURLVars(httptest.NewRequest("GET", "/categories/1", nil), map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetById(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		itemBytes, _ := json.Marshal(response.Data)
		var result models_user_categories.UserCategory
		json.Unmarshal(itemBytes, &result)

		assert.Equal(t, *expected, result)
		assert.Equal(t, "Categoria recuperada com sucesso", response.Message)

		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		handler := NewUserCategoryHandler(new(user_category_services_mock.MockUserCategoryService))
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
		mockSvc := new(user_category_services_mock.MockUserCategoryService)
		handler := NewUserCategoryHandler(mockSvc)

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
		mockSvc := new(user_category_services_mock.MockUserCategoryService)
		handler := NewUserCategoryHandler(mockSvc)

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
		handler := NewUserCategoryHandler(new(user_category_services_mock.MockUserCategoryService))

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
		handler := NewUserCategoryHandler(new(user_category_services_mock.MockUserCategoryService))

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
		mockSvc := new(user_category_services_mock.MockUserCategoryService)
		handler := NewUserCategoryHandler(mockSvc)

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
		assert.Contains(t, resp.Message, "erro ao atualizar")

		mockSvc.AssertExpectations(t)
	})
}

func TestUserCategoryHandler_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(user_category_services_mock.MockUserCategoryService)
		handler := NewUserCategoryHandler(mockSvc)

		mockSvc.On("Delete", mock.Anything, int64(1)).Return(nil)

		req := mux.SetURLVars(httptest.NewRequest("DELETE", "/categories/1", nil), map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Equal(t, 0, w.Body.Len())

		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		handler := NewUserCategoryHandler(new(user_category_services_mock.MockUserCategoryService))

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
		mockSvc := new(user_category_services_mock.MockUserCategoryService)
		handler := NewUserCategoryHandler(mockSvc)

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
