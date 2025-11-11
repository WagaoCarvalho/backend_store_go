package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockUserCatRel "github.com/WagaoCarvalho/backend_store_go/infra/mock/user"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/category_relation"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/user/category_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserCategoryRelationHandler_Create(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("success - relação criada", func(t *testing.T) {
		mockService := new(mockUserCatRel.MockUserCategoryRelation)
		handler := NewUserCategoryRelationHandler(mockService, logger)

		dtoRel := dto.UserCategoryRelationsDTO{
			UserID:     *utils.Int64Ptr(1),
			CategoryID: *utils.Int64Ptr(2),
		}
		modelRel := dto.ToUserCategoryRelationsModel(dtoRel)

		mockService.On("Create", mock.Anything, mock.MatchedBy(func(rel *models.UserCategoryRelation) bool {
			return rel.UserID == 1 && rel.CategoryID == 2
		})).Return(modelRel, nil).Once()

		body, _ := json.Marshal(dtoRel)
		req := httptest.NewRequest(http.MethodPost, "/relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, "Relação criada com sucesso", resp.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("success - relação já existente", func(t *testing.T) {
		mockService := new(mockUserCatRel.MockUserCategoryRelation)
		handler := NewUserCategoryRelationHandler(mockService, logger)

		dtoRel := dto.UserCategoryRelationsDTO{
			UserID:     *utils.Int64Ptr(1),
			CategoryID: *utils.Int64Ptr(2),
		}

		mockService.On("Create", mock.Anything, mock.MatchedBy(func(rel *models.UserCategoryRelation) bool {
			return rel.UserID == 1 && rel.CategoryID == 2
		})).Return(nil, errMsg.ErrRelationExists).Once()

		body, _ := json.Marshal(dtoRel)
		req := httptest.NewRequest(http.MethodPost, "/relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, "Relação já existente", resp.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("error - corpo inválido (JSON parse)", func(t *testing.T) {
		handler := NewUserCategoryRelationHandler(new(mockUserCatRel.MockUserCategoryRelation), logger)

		req := httptest.NewRequest(http.MethodPost, "/relations", bytes.NewBufferString("invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, http.StatusBadRequest, resp.Status)
	})

	t.Run("error - chave estrangeira inválida", func(t *testing.T) {
		mockService := new(mockUserCatRel.MockUserCategoryRelation)
		handler := NewUserCategoryRelationHandler(mockService, logger)

		dtoRel := dto.UserCategoryRelationsDTO{
			UserID:     *utils.Int64Ptr(99),
			CategoryID: *utils.Int64Ptr(88),
		}

		mockService.On("Create", mock.Anything, mock.MatchedBy(func(rel *models.UserCategoryRelation) bool {
			return rel.UserID == 99 && rel.CategoryID == 88
		})).Return(nil, errMsg.ErrDBInvalidForeignKey).Once()

		body, _ := json.Marshal(dtoRel)
		req := httptest.NewRequest(http.MethodPost, "/relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, http.StatusBadRequest, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("error - falha interna do serviço", func(t *testing.T) {
		mockService := new(mockUserCatRel.MockUserCategoryRelation)
		handler := NewUserCategoryRelationHandler(mockService, logger)

		dtoRel := dto.UserCategoryRelationsDTO{
			UserID:     *utils.Int64Ptr(1),
			CategoryID: *utils.Int64Ptr(2),
		}

		mockService.On("Create", mock.Anything, mock.MatchedBy(func(rel *models.UserCategoryRelation) bool {
			return rel.UserID == 1 && rel.CategoryID == 2
		})).Return(nil, errors.New("erro inesperado")).Once()

		body, _ := json.Marshal(dtoRel)
		req := httptest.NewRequest(http.MethodPost, "/relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, http.StatusInternalServerError, resp.Status)
		assert.Contains(t, resp.Message, "erro ao criar relação")

		mockService.AssertExpectations(t)
	})

	t.Run("error - modelo nulo ou ID inválido", func(t *testing.T) {
		baseLogger := logrus.New()
		baseLogger.Out = &bytes.Buffer{}

		handler := NewUserCategoryRelationHandler(new(mockUserCatRel.MockUserCategoryRelation), logger)

		// Testa JSON válido, mas com IDs zerados
		dtoRel := dto.UserCategoryRelationsDTO{
			UserID:     *utils.Int64Ptr(0),
			CategoryID: *utils.Int64Ptr(0),
		}

		body, _ := json.Marshal(dtoRel)
		req := httptest.NewRequest(http.MethodPost, "/relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "modelo nulo ou ID inválido")
	})

}

func TestUserCategoryRelationHandler_Delete(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("success - relação deletada com sucesso", func(t *testing.T) {
		mockService := new(mockUserCatRel.MockUserCategoryRelation)
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
		mockService := new(mockUserCatRel.MockUserCategoryRelation)
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
		mockService := new(mockUserCatRel.MockUserCategoryRelation)
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
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("success - todas as relações deletadas com sucesso", func(t *testing.T) {
		mockService := new(mockUserCatRel.MockUserCategoryRelation)
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
		mockService := new(mockUserCatRel.MockUserCategoryRelation)
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
		mockService := new(mockUserCatRel.MockUserCategoryRelation)
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
