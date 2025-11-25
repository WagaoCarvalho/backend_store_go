package handler

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockUserCatRel "github.com/WagaoCarvalho/backend_store_go/infra/mock/user"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/category_relation"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserCategoryRelationHandler_GetAllRelationsByUserID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)
	t.Run("success - retorna todas as relações do usuário", func(t *testing.T) {
		mockService := new(mockUserCatRel.MockUserCategoryRelation)
		handler := NewUserCategoryRelationHandler(mockService, logger)

		expected := []*models.UserCategoryRelation{
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
		assert.Contains(t, rr.Body.String(), "\"category_id\":1")
		assert.Contains(t, rr.Body.String(), "Relações recuperadas com sucesso")

		mockService.AssertExpectations(t)
	})

	t.Run("error - ID inválido", func(t *testing.T) {
		mockService := new(mockUserCatRel.MockUserCategoryRelation)
		handler := NewUserCategoryRelationHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodGet, "/user-category-relations/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"user_id": "abc"})
		rr := httptest.NewRecorder()

		handler.GetAllRelationsByUserID(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
		assert.Contains(t, rr.Body.String(), "ID de usuário inválido")
	})

	t.Run("error - falha no serviço", func(t *testing.T) {
		mockService := new(mockUserCatRel.MockUserCategoryRelation)
		handler := NewUserCategoryRelationHandler(mockService, logger)

		mockService.
			On("GetAllRelationsByUserID", mock.Anything, int64(1)).
			Return(([]*models.UserCategoryRelation)(nil), errors.New("erro interno"))

		req := httptest.NewRequest(http.MethodGet, "/user-category-relations/1", nil)
		req = mux.SetURLVars(req, map[string]string{"user_id": "1"})
		rr := httptest.NewRecorder()

		handler.GetAllRelationsByUserID(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "erro interno")

		mockService.AssertExpectations(t)
	})
}
