package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockUser "github.com/WagaoCarvalho/backend_store_go/infra/mock/service/user"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/user_contact_relation"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_contact_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupUserContact() (*mockUser.MockUserContactRelationService, *UserContactRelation) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	mockService := new(mockUser.MockUserContactRelationService)
	handler := NewUserContactRelation(mockService, logAdapter)

	return mockService, handler
}

func TestUserContactRelationHandler_Create(t *testing.T) {
	t.Run("success - relação criada", func(t *testing.T) {
		mockService, handler := setupUserContact()

		relationDTO := dto.UserContactRelationDTO{
			UserID:    1,
			ContactID: 2,
		}
		expectedModel := &model.UserContactRelation{
			UserID:    1,
			ContactID: 2,
			CreatedAt: time.Now(),
		}

		mockService.On("Create", mock.Anything, int64(1), int64(2)).
			Return(expectedModel, true, nil).Once()

		body, _ := json.Marshal(relationDTO)
		req := httptest.NewRequest(http.MethodPost, "/user-contact-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Relação criada com sucesso", response.Message)

		dataBytes, _ := json.Marshal(response.Data)
		var got dto.UserContactRelationDTO
		_ = json.Unmarshal(dataBytes, &got)

		assert.Equal(t, int64(1), got.UserID)
		assert.Equal(t, int64(2), got.ContactID)

		mockService.AssertExpectations(t)
	})

	t.Run("success - relação já existente", func(t *testing.T) {
		mockService, handler := setupUserContact()

		relationDTO := dto.UserContactRelationDTO{
			UserID:    1,
			ContactID: 2,
		}
		expectedModel := &model.UserContactRelation{
			UserID:    1,
			ContactID: 2,
			CreatedAt: time.Now(),
		}

		mockService.On("Create", mock.Anything, int64(1), int64(2)).
			Return(expectedModel, false, nil).Once()

		body, _ := json.Marshal(relationDTO)
		req := httptest.NewRequest(http.MethodPost, "/user-contact-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response utils.DefaultResponse
		_ = json.NewDecoder(resp.Body).Decode(&response)
		assert.Equal(t, "Relação já existente", response.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("error - chave estrangeira inválida", func(t *testing.T) {
		mockService, handler := setupUserContact()

		relationDTO := dto.UserContactRelationDTO{
			UserID:    1,
			ContactID: 99,
		}

		mockService.On("Create", mock.Anything, int64(1), int64(99)).
			Return(nil, false, errMsg.ErrDBInvalidForeignKey).Once()

		body, _ := json.Marshal(relationDTO)
		req := httptest.NewRequest(http.MethodPost, "/user-contact-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		mockService.AssertExpectations(t)
	})

	t.Run("error - interno no serviço", func(t *testing.T) {
		mockService, handler := setupUserContact()

		relationDTO := dto.UserContactRelationDTO{
			UserID:    1,
			ContactID: 2,
		}

		mockService.On("Create", mock.Anything, int64(1), int64(2)).
			Return(nil, false, errors.New("db error")).Once()

		body, _ := json.Marshal(relationDTO)
		req := httptest.NewRequest(http.MethodPost, "/user-contact-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockService.AssertExpectations(t)
	})

	t.Run("error - JSON inválido", func(t *testing.T) {
		_, handler := setupUserContact()

		req := httptest.NewRequest(http.MethodPost, "/user-contact-relations", bytes.NewBuffer([]byte("{invalid json")))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

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

func TestUserContactRelationHandler_Delete(t *testing.T) {
	t.Run("success - relação deletada", func(t *testing.T) {
		mockService, handler := setupUserContact()

		mockService.On("Delete", mock.Anything, int64(1), int64(2)).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/user-contact-relations?user_id=1&contact_id=2", nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id":    "1",
			"contact_id": "2",
		})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Relação deletada com sucesso", response.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("error - user_id inválido", func(t *testing.T) {
		_, handler := setupUserContact()

		req := httptest.NewRequest(http.MethodDelete, "/user-contact-relations?user_id=abc&contact_id=2", nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id":    "abc",
			"contact_id": "2",
		})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("error - contact_id inválido", func(t *testing.T) {
		_, handler := setupUserContact()

		req := httptest.NewRequest(http.MethodDelete, "/user-contact-relations?user_id=1&contact_id=abc", nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id":    "1",
			"contact_id": "abc",
		})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("error - service falha", func(t *testing.T) {
		mockService, handler := setupUserContact()

		mockService.On("Delete", mock.Anything, int64(1), int64(2)).Return(errors.New("db error")).Once()

		req := httptest.NewRequest(http.MethodDelete, "/user-contact-relations?user_id=1&contact_id=2", nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id":    "1",
			"contact_id": "2",
		})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockService.AssertExpectations(t)
	})
}

func TestUserContactRelationHandler_DeleteAll(t *testing.T) {
	t.Run("success - todas relações deletadas", func(t *testing.T) {
		mockService, handler := setupUserContact()

		mockService.On("DeleteAll", mock.Anything, int64(1)).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/user-contact-relations?user_id=1", nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id": "1",
		})
		rec := httptest.NewRecorder()

		handler.DeleteAll(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Relações deletadas com sucesso", response.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("error - user_id inválido", func(t *testing.T) {
		_, handler := setupUserContact()

		req := httptest.NewRequest(http.MethodDelete, "/user-contact-relations?user_id=abc", nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id": "abc",
		})
		rec := httptest.NewRecorder()

		handler.DeleteAll(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("error - service falha", func(t *testing.T) {
		mockService, handler := setupUserContact()

		mockService.On("DeleteAll", mock.Anything, int64(1)).Return(errors.New("db error")).Once()

		req := httptest.NewRequest(http.MethodDelete, "/user-contact-relations?user_id=1", nil)
		req = mux.SetURLVars(req, map[string]string{
			"user_id": "1",
		})
		rec := httptest.NewRecorder()

		handler.DeleteAll(rec, req)

		resp := rec.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockService.AssertExpectations(t)
	})
}
