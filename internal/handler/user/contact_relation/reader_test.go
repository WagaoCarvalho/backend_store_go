package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/contact_relation"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/user/contact_relation"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserContactRelationHandler_GetAllByUserID(t *testing.T) {
	t.Run("success - retorna relações", func(t *testing.T) {
		mockService, handler := setupUserContact()

		expectedRelations := []*model.UserContactRelation{
			{
				UserID:    1,
				ContactID: 10,
				CreatedAt: time.Now(),
			},
			{
				UserID:    1,
				ContactID: 20,
				CreatedAt: time.Now(),
			},
		}

		mockService.On("GetAllRelationsByUserID", mock.Anything, int64(1)).
			Return(expectedRelations, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/user-contact-relations/1", nil)
		req = mux.SetURLVars(req, map[string]string{"user_id": "1"})
		rec := httptest.NewRecorder()

		handler.GetAllByUserID(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Relações encontradas", response.Message)

		dataBytes, _ := json.Marshal(response.Data)
		var got []dto.UserContactRelationDTO
		_ = json.Unmarshal(dataBytes, &got)

		assert.Len(t, got, 2)
		assert.Equal(t, int64(10), got[0].ContactID)
		assert.Equal(t, int64(20), got[1].ContactID)

		mockService.AssertExpectations(t)
	})

	t.Run("error - ID inválido", func(t *testing.T) {
		_, handler := setupUserContact()

		req := httptest.NewRequest(http.MethodGet, "/user-contact-relations/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"user_id": "abc"})
		rec := httptest.NewRecorder()

		handler.GetAllByUserID(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("error - service falha", func(t *testing.T) {
		mockService, handler := setupUserContact()

		mockService.On("GetAllRelationsByUserID", mock.Anything, int64(1)).
			Return(nil, errors.New("db error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/user-contact-relations/1", nil)
		req = mux.SetURLVars(req, map[string]string{"user_id": "1"})
		rec := httptest.NewRecorder()

		handler.GetAllByUserID(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockService.AssertExpectations(t)
	})
}

func TestUserContactRelationHandler_HasRelation(t *testing.T) {
	t.Run("success - relação existe", func(t *testing.T) {
		mockService, handler := setupUserContact()

		mockService.On("HasUserContactRelation", mock.Anything, int64(1), int64(2)).
			Return(true, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/user-contact-relations/has?user_id=1&contact_id=2", nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id":    "1",
			"contact_id": "2",
		})
		rec := httptest.NewRecorder()

		handler.HasRelation(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Verificação concluída", response.Message)

		dataMap := response.Data.(map[string]any)
		assert.Equal(t, true, dataMap["exists"])

		mockService.AssertExpectations(t)
	})

	t.Run("success - relação não existe", func(t *testing.T) {
		mockService, handler := setupUserContact()

		mockService.On("HasUserContactRelation", mock.Anything, int64(1), int64(3)).
			Return(false, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/user-contact-relations/has?user_id=1&contact_id=3", nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id":    "1",
			"contact_id": "3",
		})
		rec := httptest.NewRecorder()

		handler.HasRelation(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response utils.DefaultResponse
		_ = json.NewDecoder(resp.Body).Decode(&response)

		dataMap := response.Data.(map[string]any)
		assert.Equal(t, false, dataMap["exists"])

		mockService.AssertExpectations(t)
	})

	t.Run("error - user_id inválido", func(t *testing.T) {
		_, handler := setupUserContact()

		req := httptest.NewRequest(http.MethodGet, "/user-contact-relations/has?user_id=abc&contact_id=2", nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id":    "abc",
			"contact_id": "2",
		})
		rec := httptest.NewRecorder()

		handler.HasRelation(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("error - contact_id inválido", func(t *testing.T) {
		_, handler := setupUserContact()

		req := httptest.NewRequest(http.MethodGet, "/user-contact-relations/has?user_id=1&contact_id=abc", nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id":    "1",
			"contact_id": "abc",
		})
		rec := httptest.NewRecorder()

		handler.HasRelation(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("error - service falha", func(t *testing.T) {
		mockService, handler := setupUserContact()

		mockService.On("HasUserContactRelation", mock.Anything, int64(1), int64(2)).
			Return(false, errors.New("db error")).Once()

		req := httptest.NewRequest(http.MethodGet, "/user-contact-relations/has?user_id=1&contact_id=2", nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id":    "1",
			"contact_id": "2",
		})
		rec := httptest.NewRecorder()

		handler.HasRelation(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockService.AssertExpectations(t)
	})
}
