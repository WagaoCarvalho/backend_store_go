package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	model "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_categories"
	err_msg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	service_mock "github.com/WagaoCarvalho/backend_store_go/internal/service/user/user_categories/user_categories_services_mock"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserCategoryHandler_Create(t *testing.T) {
	mockSvc := new(service_mock.MockUserCategoryService)
	logger := logger.NewLoggerAdapter(logrus.New())
	handler := NewUserCategoryHandler(mockSvc, logger)

	t.Run("Success", func(t *testing.T) {
		category := &model.UserCategory{Name: "Nova"}
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
		var result model.UserCategory
		json.Unmarshal(itemBytes, &result)

		assert.Equal(t, category.Name, result.Name)
		assert.Equal(t, "Categoria criada com sucesso", response.Message)

		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidJSON", func(t *testing.T) {
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
		input := model.UserCategory{Name: "Erro"}

		mockSvc.On("Create", mock.Anything, mock.MatchedBy(func(c *model.UserCategory) bool {
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

func TestUserCategoryHandler_GetById(t *testing.T) {
	mockSvc := new(service_mock.MockUserCategoryService)
	logger := logger.NewLoggerAdapter(logrus.New())
	handler := NewUserCategoryHandler(mockSvc, logger)

	t.Run("Success", func(t *testing.T) {
		expected := &model.UserCategory{ID: 1, Name: "Teste"}
		mockSvc.On("GetByID", mock.Anything, int64(1)).Return(expected, nil)

		req := mux.SetURLVars(httptest.NewRequest("GET", "/categories/1", nil), map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetById(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)

		itemBytes, _ := json.Marshal(response.Data)
		var result model.UserCategory
		json.Unmarshal(itemBytes, &result)

		assert.Equal(t, *expected, result)
		assert.Equal(t, "Categoria recuperada com sucesso", response.Message)

		mockSvc.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockSvc := new(service_mock.MockUserCategoryService)
		handler := NewUserCategoryHandler(mockSvc, logger)

		mockSvc.On("GetByID", mock.Anything, int64(42)).Return(nil, errors.New("erro inesperado"))

		req := mux.SetURLVars(httptest.NewRequest("GET", "/categories/42", nil), map[string]string{"id": "42"})
		w := httptest.NewRecorder()

		handler.GetById(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)
		assert.Equal(t, "erro inesperado", resp.Message)

		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
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

func TestUserCategoryHandler_GetAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockSvc := new(service_mock.MockUserCategoryService)
		logger := logger.NewLoggerAdapter(logrus.New())
		handler := NewUserCategoryHandler(mockSvc, logger)

		expected := []*model.UserCategory{{ID: 1, Name: "Categoria"}}
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

		var result []*model.UserCategory
		for _, item := range rawData {
			itemBytes, _ := json.Marshal(item)
			var cat model.UserCategory
			json.Unmarshal(itemBytes, &cat)
			result = append(result, &cat)
		}

		assert.Equal(t, expected, result)
		assert.Equal(t, "Categorias recuperadas com sucesso", response.Message)

		mockSvc.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockSvc := new(service_mock.MockUserCategoryService)
		logger := logger.NewLoggerAdapter(logrus.New())
		handler := NewUserCategoryHandler(mockSvc, logger)

		mockSvc.On("GetAll", mock.Anything).Return(([]*model.UserCategory)(nil), errors.New("erro de banco"))

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

func TestUserCategoryHandler_Update(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	t.Run("Success", func(t *testing.T) {
		mockSvc := new(service_mock.MockUserCategoryService)
		handler := NewUserCategoryHandler(mockSvc, logger)

		category := &model.UserCategory{ID: 1, Name: "Atualizada"}
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
		var result model.UserCategory
		json.Unmarshal(itemBytes, &result)

		assert.Equal(t, category.Name, result.Name)
		assert.Equal(t, "Categoria atualizada com sucesso", response.Message)

		mockSvc.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockSvc := new(service_mock.MockUserCategoryService)
		handler := NewUserCategoryHandler(mockSvc, logger)

		category := &model.UserCategory{Name: "Inexistente"}
		body, _ := json.Marshal(category)

		req := mux.SetURLVars(httptest.NewRequest("PUT", "/categories/999", bytes.NewBuffer(body)), map[string]string{"id": "999"})
		w := httptest.NewRecorder()

		mockSvc.On("Update", mock.Anything, mock.MatchedBy(func(c *model.UserCategory) bool {
			return c.ID == 999
		})).Return(nil, err_msg.ErrCategoryNotFound)

		handler.Update(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(w.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.Status)
		assert.Equal(t, "categoria não encontrada", resp.Message)

		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		handler := NewUserCategoryHandler(new(service_mock.MockUserCategoryService), logger)

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
		handler := NewUserCategoryHandler(new(service_mock.MockUserCategoryService), logger)

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
		mockSvc := new(service_mock.MockUserCategoryService)
		handler := NewUserCategoryHandler(mockSvc, logger)

		category := &model.UserCategory{ID: 2, Name: "Falha"}
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
	logger := logger.NewLoggerAdapter(logrus.New())

	t.Run("Success", func(t *testing.T) {
		mockSvc := new(service_mock.MockUserCategoryService)
		handler := NewUserCategoryHandler(mockSvc, logger)

		mockSvc.On("Delete", mock.Anything, int64(1)).Return(nil)

		req := mux.SetURLVars(httptest.NewRequest("DELETE", "/categories/1", nil), map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.Delete(w, req)

		assert.Equal(t, http.StatusNoContent, w.Code)
		assert.Empty(t, w.Body.String())

		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		handler := NewUserCategoryHandler(new(service_mock.MockUserCategoryService), logger)

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
		mockSvc := new(service_mock.MockUserCategoryService)
		handler := NewUserCategoryHandler(mockSvc, logger)

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
