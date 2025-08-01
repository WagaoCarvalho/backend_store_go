package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	user_category_relations_models "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_category_relations"
	repositories "github.com/WagaoCarvalho/backend_store_go/internal/repositories/addresses"
	user_category_relations_mock "github.com/WagaoCarvalho/backend_store_go/internal/services/users/user_category_relations/user_category_relations_services_mock"
	"github.com/WagaoCarvalho/backend_store_go/logger"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserCategoryRelationHandler_Create(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	t.Run("success - relação criada", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService, logger)

		relation := &user_category_relations_models.UserCategoryRelations{
			UserID:     1,
			CategoryID: 2,
		}

		mockService.
			On("Create", mock.Anything, int64(1), int64(2)).
			Return(relation, true, nil)

		body, _ := json.Marshal(relation)
		req := httptest.NewRequest(http.MethodPost, "/relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("success - relação já existia", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService, logger)

		relation := &user_category_relations_models.UserCategoryRelations{
			UserID:     1,
			CategoryID: 2,
		}

		mockService.
			On("Create", mock.Anything, int64(1), int64(2)).
			Return(relation, false, nil)

		body, _ := json.Marshal(relation)
		req := httptest.NewRequest(http.MethodPost, "/relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("error - corpo inválido", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodPost, "/relations", bytes.NewBufferString("invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("error - falha no serviço", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService, logger)

		relation := &user_category_relations_models.UserCategoryRelations{
			UserID:     1,
			CategoryID: 2,
		}

		mockService.
			On("Create", mock.Anything, int64(1), int64(2)).
			Return(nil, false, errors.New("erro ao criar relação"))

		body, _ := json.Marshal(relation)
		req := httptest.NewRequest(http.MethodPost, "/relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserCategoryRelationHandler_Create_ForeignKeyInvalid(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())
	mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
	handler := NewUserCategoryRelationHandler(mockService, logger)

	body := `{"user_id":1,"category_id":999}`

	mockService.
		On("Create", mock.Anything, int64(1), int64(999)).
		Return(nil, false, repositories.ErrInvalidForeignKey)

	req := httptest.NewRequest(http.MethodPost, "/relations", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Create(rec, req)

	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Contains(t, rec.Body.String(), repositories.ErrInvalidForeignKey.Error())
	mockService.AssertExpectations(t)
}

func TestUserCategoryRelationHandler_GetAllRelationsByUserID(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	t.Run("success - retorna todas as relações do usuário", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService, logger)

		expected := []*user_category_relations_models.UserCategoryRelations{
			{UserID: 1, CategoryID: 10},
			{UserID: 1, CategoryID: 20},
		}

		mockService.
			On("GetAllRelationsByUserID", mock.Anything, int64(1)).
			Return(expected, nil)

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
		handler := NewUserCategoryRelationHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodGet, "/user-category-relations/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"user_id": "abc"})
		rr := httptest.NewRecorder()

		handler.GetAllRelationsByUserID(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "ID de usuário inválido")
	})

	t.Run("error - falha no serviço", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService, logger)

		mockService.
			On("GetAllRelationsByUserID", mock.Anything, int64(1)).
			Return(nil, errors.New("erro interno"))

		req := httptest.NewRequest(http.MethodGet, "/user-category-relations/1", nil)
		req = mux.SetURLVars(req, map[string]string{"user_id": "1"})
		rr := httptest.NewRecorder()

		handler.GetAllRelationsByUserID(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "erro interno")

		mockService.AssertExpectations(t)
	})
}

func TestUserCategoryRelationHandler_HasUserCategoryRelation(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	t.Run("success - relação existe", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService, logger)

		mockService.On("HasUserCategoryRelation", mock.Anything, int64(1), int64(2)).Return(true, nil)

		req := httptest.NewRequest(http.MethodGet, "/user-category-relation/1/category/2", nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id":     "1",
			"category_id": "2",
		})
		rr := httptest.NewRecorder()

		handler.HasUserCategoryRelation(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), `"exists":true`)
		assert.Contains(t, rr.Body.String(), "Verificação concluída com sucesso")

		mockService.AssertExpectations(t)
	})

	t.Run("success - relação não existe", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService, logger)

		mockService.On("HasUserCategoryRelation", mock.Anything, int64(1), int64(3)).Return(false, nil)

		req := httptest.NewRequest(http.MethodGet, "/user-category-relation/1/category/3", nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id":     "1",
			"category_id": "3",
		})
		rr := httptest.NewRecorder()

		handler.HasUserCategoryRelation(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Contains(t, rr.Body.String(), `"exists":false`)
		assert.Contains(t, rr.Body.String(), "Verificação concluída com sucesso")

		mockService.AssertExpectations(t)
	})

	t.Run("error - user_id inválido", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodGet, "/user-category-relation/abc/category/2", nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id":     "abc",
			"category_id": "2",
		})
		rr := httptest.NewRecorder()

		handler.HasUserCategoryRelation(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "ID de usuário inválido")
	})

	t.Run("error - category_id inválido", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodGet, "/user-category-relation/1/category/xyz", nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id":     "1",
			"category_id": "xyz",
		})
		rr := httptest.NewRecorder()

		handler.HasUserCategoryRelation(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "ID de categoria inválido")
	})

	t.Run("error - falha no serviço", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService, logger)

		mockService.On("HasUserCategoryRelation", mock.Anything, int64(1), int64(2)).Return(false, errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodGet, "/user-category-relation/1/category/2", nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id":     "1",
			"category_id": "2",
		})
		rr := httptest.NewRecorder()

		handler.HasUserCategoryRelation(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "erro interno")

		mockService.AssertExpectations(t)
	})
}

func TestUserCategoryRelationHandler_Delete(t *testing.T) {
	logger := logger.NewLoggerAdapter(logrus.New())

	t.Run("success - relação deletada com sucesso", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService, logger)

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
		assert.Empty(t, rr.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("bad request - IDs inválidos", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService, logger)

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
		handler := NewUserCategoryRelationHandler(mockService, logger)

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
	logger := logger.NewLoggerAdapter(logrus.New())

	t.Run("success - todas as relações deletadas com sucesso", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService, logger)

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
		assert.Empty(t, rr.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("bad request - ID de usuário inválido", func(t *testing.T) {
		mockService := new(user_category_relations_mock.MockUserCategoryRelationService)
		handler := NewUserCategoryRelationHandler(mockService, logger)

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
		handler := NewUserCategoryRelationHandler(mockService, logger)

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
