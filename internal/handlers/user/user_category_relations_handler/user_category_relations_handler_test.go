package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	user_category_relations_models "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_category_relations"
	user_category_relations_mock "github.com/WagaoCarvalho/backend_store_go/internal/services/user/user_category_relations/user_category_relations_services_mock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserCategoryRelationHandler_Create(t *testing.T) {
	t.Run("success - relação criada", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService)

		relation := &user_category_relations_models.UserCategoryRelations{
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

		relation := &user_category_relations_models.UserCategoryRelations{
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

func TestUserCategoryRelationHandler_GetAllRelationsByUserID(t *testing.T) {
	t.Run("success - retorna todas as relações do usuário", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService)

		expected := []*user_category_relations_models.UserCategoryRelations{
			{UserID: 1, CategoryID: 10},
			{UserID: 1, CategoryID: 20},
		}

		mockService.On("GetAllRelationsByUserID", mock.Anything, int64(1)).Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/user-category-relations/1", nil)
		req = mux.SetURLVars(req, map[string]string{"user_id": "1"})

		rr := httptest.NewRecorder()

		handler.GetAllRelationsByUserID(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "\"user_id\":1")
		assert.Contains(t, rr.Body.String(), "\"category_id\":10")
		assert.Contains(t, rr.Body.String(), "Relações recuperadas com sucesso")

		mockService.AssertExpectations(t)
	})

	t.Run("error - ID inválido", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/user-category-relations/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rr := httptest.NewRecorder()

		handler.GetAllRelationsByUserID(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "ID de usuário inválido")
	})

	t.Run("error - falha no serviço", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService)

		mockService.On("GetAllRelationsByUserID", mock.Anything, int64(1)).Return(nil, errors.New("erro interno"))

		req := httptest.NewRequest(http.MethodGet, "/user-category-relations/1", nil)
		req = mux.SetURLVars(req, map[string]string{"user_id": "1"})

		rr := httptest.NewRecorder()

		handler.GetAllRelationsByUserID(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "erro interno")
	})
}

func TestUserCategoryRelationHandler_GetVersionByUserID(t *testing.T) {
	t.Run("success - versão retornada com sucesso", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService)

		userID := int64(1)
		expectedVersion := 5

		mockService.
			On("GetVersionByUserID", mock.Anything, userID).
			Return(expectedVersion, nil)

		req := httptest.NewRequest(http.MethodGet, "/user-category-relation/version/1", nil)
		req = mux.SetURLVars(req, map[string]string{"user_id": strconv.FormatInt(userID, 10)})
		rr := httptest.NewRecorder()

		handler.GetVersionByUserID(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), "\"Vers\u00e3o recuperada com sucesso")
		assert.Contains(t, rr.Body.String(), fmt.Sprintf("%d", expectedVersion))
		mockService.AssertExpectations(t)
	})

	t.Run("error - ID de usuário inválido", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/user-category-relation/version/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"user_id": "abc"})
		rr := httptest.NewRecorder()

		handler.GetVersionByUserID(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "ID de usu\u00e1rio inv\u00e1lido")
	})

	t.Run("error - falha ao buscar versão", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService)

		userID := int64(1)
		errExpected := errors.New("erro no banco")

		mockService.
			On("GetVersionByUserID", mock.Anything, userID).
			Return(0, errExpected)

		req := httptest.NewRequest(http.MethodGet, "/user-category-relation/version/1", nil)
		req = mux.SetURLVars(req, map[string]string{"user_id": strconv.FormatInt(userID, 10)})
		rr := httptest.NewRecorder()

		handler.GetVersionByUserID(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "erro no banco")
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
