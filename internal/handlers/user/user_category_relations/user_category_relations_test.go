package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	user_category_relations "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_category_relations"
	user_category_relations_mock "github.com/WagaoCarvalho/backend_store_go/internal/services/user/user_category_relations/user_category_relations_mock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserCategoryRelationHandler_Create(t *testing.T) {
	t.Run("success - relação criada", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService)

		relation := &user_category_relations.UserCategoryRelations{
			UserID:     1,
			CategoryID: 2,
		}

		mockService.
			On("Create", mock.Anything, int64(1), int64(2)).
			Return(relation, nil)

		body, _ := json.Marshal(relation)
		req := httptest.NewRequest(http.MethodPost, "/relations", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("error - corpo inválido", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService)

		req := httptest.NewRequest(http.MethodPost, "/relations", bytes.NewBufferString("invalid-json"))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("error - falha no serviço", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService)

		relation := &user_category_relations.UserCategoryRelations{
			UserID:     1,
			CategoryID: 2,
		}

		mockService.
			On("Create", mock.Anything, int64(1), int64(2)).
			Return(nil, errors.New("erro ao criar relação"))

		body, _ := json.Marshal(relation)
		req := httptest.NewRequest(http.MethodPost, "/relations", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserCategoryRelationHandler_GetByUserID(t *testing.T) {
	t.Run("success - relações recuperadas com sucesso", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService)

		expected := []*user_category_relations.UserCategoryRelations{
			{UserID: 1, CategoryID: 100},
			{UserID: 1, CategoryID: 101},
		}

		mockService.
			On("GetByUserID", mock.Anything, int64(1)).
			Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/relations/user/1", nil)
		req = mux.SetURLVars(req, map[string]string{"user_id": "1"})
		rr := httptest.NewRecorder()

		handler.GetByUserID(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), `"user_id":1`)
		assert.Contains(t, rr.Body.String(), `"category_id":100`)
		assert.Contains(t, rr.Body.String(), `"Relações recuperadas com sucesso`)

		mockService.AssertExpectations(t)
	})

	t.Run("error - ID inválido", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/relations/user/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"user_id": "abc"})
		rr := httptest.NewRecorder()

		handler.GetByUserID(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "ID de usuário inválido")
	})

	t.Run("error - serviço retorna erro", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService)

		mockService.
			On("GetByUserID", mock.Anything, int64(99)).
			Return(nil, errors.New("falha no banco"))

		req := httptest.NewRequest(http.MethodGet, "/relations/user/99", nil)
		req = mux.SetURLVars(req, map[string]string{"user_id": "99"})
		rr := httptest.NewRecorder()

		handler.GetByUserID(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "falha no banco")

		mockService.AssertExpectations(t)
	})
}

func TestUserCategoryRelationHandler_GetByCategoryID(t *testing.T) {
	t.Run("success - relações recuperadas com sucesso", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService)

		expected := []*user_category_relations.UserCategoryRelations{
			{UserID: 1, CategoryID: 100},
			{UserID: 2, CategoryID: 100},
		}

		mockService.
			On("GetAll", mock.Anything, int64(100)).
			Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/relations/category/100", nil)
		req = mux.SetURLVars(req, map[string]string{"category_id": "100"})
		rr := httptest.NewRecorder()

		handler.GetByCategoryID(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), `"user_id":1`)
		assert.Contains(t, rr.Body.String(), `"user_id":2`)
		assert.Contains(t, rr.Body.String(), `"category_id":100`)
		assert.Contains(t, rr.Body.String(), `"Relações recuperadas com sucesso`)

		mockService.AssertExpectations(t)
	})

	t.Run("bad request - ID da categoria inválido", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/relations/category/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"category_id": "abc"})
		rr := httptest.NewRecorder()

		handler.GetByCategoryID(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "ID da categoria inválido")
	})

	t.Run("internal error - falha ao recuperar relações", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService)

		mockService.
			On("GetAll", mock.Anything, int64(100)).
			Return(nil, fmt.Errorf("erro interno"))

		req := httptest.NewRequest(http.MethodGet, "/relations/category/100", nil)
		req = mux.SetURLVars(req, map[string]string{"category_id": "100"})
		rr := httptest.NewRecorder()

		handler.GetByCategoryID(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "erro interno")

		mockService.AssertExpectations(t)
	})
}

func TestUserCategoryRelationHandler_Delete(t *testing.T) {
	t.Run("success - relação deletada com sucesso", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService)

		userID := int64(1)
		categoryID := int64(100)

		mockService.
			On("Delete", mock.Anything, userID, categoryID).
			Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/relations/1/100", nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id":     "1",
			"category_id": "100",
		})
		rr := httptest.NewRecorder()

		handler.Delete(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("bad request - IDs inválidos", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService)

		req := httptest.NewRequest(http.MethodDelete, "/relations/abc/xyz", nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id":     "abc",
			"category_id": "xyz",
		})
		rr := httptest.NewRecorder()

		handler.Delete(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "IDs inválidos")
	})

	t.Run("internal error - erro ao deletar relação", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService)

		userID := int64(2)
		categoryID := int64(200)

		mockService.
			On("Delete", mock.Anything, userID, categoryID).
			Return(fmt.Errorf("erro ao deletar"))

		req := httptest.NewRequest(http.MethodDelete, "/relations/2/200", nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id":     "2",
			"category_id": "200",
		})
		rr := httptest.NewRecorder()

		handler.Delete(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "erro ao deletar")
		mockService.AssertExpectations(t)
	})
}

func TestUserCategoryRelationHandler_DeleteAll(t *testing.T) {
	t.Run("success - todas as relações deletadas com sucesso", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService)

		userID := int64(1)

		mockService.
			On("DeleteAll", mock.Anything, userID).
			Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/relations/user/1", nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id": "1",
		})
		rr := httptest.NewRecorder()

		handler.DeleteAll(rr, req)

		assert.Equal(t, http.StatusNoContent, rr.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("bad request - ID de usuário inválido", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService)

		req := httptest.NewRequest(http.MethodDelete, "/relations/user/abc", nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id": "abc",
		})
		rr := httptest.NewRecorder()

		handler.DeleteAll(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "ID de usuário inválido")
	})

	t.Run("internal error - erro ao deletar todas as relações", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService)

		userID := int64(2)

		mockService.
			On("DeleteAll", mock.Anything, userID).
			Return(fmt.Errorf("erro ao deletar todas as relações"))

		req := httptest.NewRequest(http.MethodDelete, "/relations/user/2", nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id": "2",
		})
		rr := httptest.NewRecorder()

		handler.DeleteAll(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "erro ao deletar todas as relações")
		mockService.AssertExpectations(t)
	})
}
