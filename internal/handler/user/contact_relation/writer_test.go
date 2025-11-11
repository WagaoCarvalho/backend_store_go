package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	mockUser "github.com/WagaoCarvalho/backend_store_go/infra/mock/user"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/contact_relation"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/user/contact_relation"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupUserContact() (*mockUser.MockUserContactRelation, *userContactRelationHandler) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	mockService := new(mockUser.MockUserContactRelation)
	handler := NewUserContactRelationHandler(mockService, logAdapter)

	return mockService, handler
}

func TestUserContactRelationHandler_Create(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("success - relação criada", func(t *testing.T) {
		mockService := new(mockUser.MockUserContactRelation)
		handler := NewUserContactRelationHandler(mockService, logger)

		dtoRel := dto.UserContactRelationDTO{
			UserID:    *utils.Int64Ptr(1),
			ContactID: *utils.Int64Ptr(2),
		}
		createdRel := &model.UserContactRelation{
			UserID:    1,
			ContactID: 2,
			CreatedAt: time.Now(),
		}

		mockService.On("Create", mock.Anything, mock.MatchedBy(func(rel *model.UserContactRelation) bool {
			return rel.UserID == 1 && rel.ContactID == 2
		})).Return(createdRel, nil).Once()

		body, _ := json.Marshal(dtoRel)
		req := httptest.NewRequest(http.MethodPost, "/contact-relations", bytes.NewBuffer(body))
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
		mockService := new(mockUser.MockUserContactRelation)
		handler := NewUserContactRelationHandler(mockService, logger)

		dtoRel := dto.UserContactRelationDTO{
			UserID:    *utils.Int64Ptr(1),
			ContactID: *utils.Int64Ptr(2),
		}
		existingRel := &model.UserContactRelation{
			UserID:    1,
			ContactID: 2,
			CreatedAt: time.Now(),
		}

		mockService.On("Create", mock.Anything, mock.MatchedBy(func(rel *model.UserContactRelation) bool {
			return rel.UserID == 1 && rel.ContactID == 2
		})).Return(existingRel, errMsg.ErrRelationExists).Once()

		body, _ := json.Marshal(dtoRel)
		req := httptest.NewRequest(http.MethodPost, "/contact-relations", bytes.NewBuffer(body))
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
		handler := NewUserContactRelationHandler(new(mockUser.MockUserContactRelation), logger)

		req := httptest.NewRequest(http.MethodPost, "/contact-relations", bytes.NewBufferString("invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "erro ao decodificar JSON")
	})

	t.Run("error - chave estrangeira inválida", func(t *testing.T) {
		mockService := new(mockUser.MockUserContactRelation)
		handler := NewUserContactRelationHandler(mockService, logger)

		dtoRel := dto.UserContactRelationDTO{
			UserID:    *utils.Int64Ptr(99),
			ContactID: *utils.Int64Ptr(88),
		}

		mockService.On("Create", mock.Anything, mock.MatchedBy(func(rel *model.UserContactRelation) bool {
			return rel.UserID == 99 && rel.ContactID == 88
		})).Return(nil, errMsg.ErrDBInvalidForeignKey).Once()

		body, _ := json.Marshal(dtoRel)
		req := httptest.NewRequest(http.MethodPost, "/contact-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "chave estrangeira inválida")

		mockService.AssertExpectations(t)
	})

	t.Run("error - falha interna do serviço", func(t *testing.T) {
		mockService := new(mockUser.MockUserContactRelation)
		handler := NewUserContactRelationHandler(mockService, logger)

		dtoRel := dto.UserContactRelationDTO{
			UserID:    *utils.Int64Ptr(1),
			ContactID: *utils.Int64Ptr(2),
		}

		mockService.On("Create", mock.Anything, mock.MatchedBy(func(rel *model.UserContactRelation) bool {
			return rel.UserID == 1 && rel.ContactID == 2
		})).Return(nil, errors.New("erro inesperado")).Once()

		body, _ := json.Marshal(dtoRel)
		req := httptest.NewRequest(http.MethodPost, "/contact-relations", bytes.NewBuffer(body))
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

	t.Run("error - modelo nulo ou ID inválido (UserID zero)", func(t *testing.T) {
		handler := NewUserContactRelationHandler(new(mockUser.MockUserContactRelation), logger)

		dtoRel := dto.UserContactRelationDTO{
			UserID:    *utils.Int64Ptr(0), // UserID zero
			ContactID: *utils.Int64Ptr(2),
		}

		body, _ := json.Marshal(dtoRel)
		req := httptest.NewRequest(http.MethodPost, "/contact-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "modelo nulo ou ID inválido")
	})

	t.Run("error - modelo nulo ou ID inválido (ContactID zero)", func(t *testing.T) {
		handler := NewUserContactRelationHandler(new(mockUser.MockUserContactRelation), logger)

		dtoRel := dto.UserContactRelationDTO{
			UserID:    *utils.Int64Ptr(1),
			ContactID: *utils.Int64Ptr(0), // ContactID zero
		}

		body, _ := json.Marshal(dtoRel)
		req := httptest.NewRequest(http.MethodPost, "/contact-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "modelo nulo ou ID inválido")
	})

	t.Run("error - modelo nulo ou ID inválido (ambos IDs zero)", func(t *testing.T) {
		handler := NewUserContactRelationHandler(new(mockUser.MockUserContactRelation), logger)

		dtoRel := dto.UserContactRelationDTO{
			UserID:    *utils.Int64Ptr(0), // Ambos IDs zero
			ContactID: *utils.Int64Ptr(0),
		}

		body, _ := json.Marshal(dtoRel)
		req := httptest.NewRequest(http.MethodPost, "/contact-relations", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "modelo nulo ou ID inválido")
	})

	t.Run("error - JSON vazio", func(t *testing.T) {
		handler := NewUserContactRelationHandler(new(mockUser.MockUserContactRelation), logger)

		req := httptest.NewRequest(http.MethodPost, "/contact-relations", bytes.NewBufferString("{}"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.Equal(t, http.StatusBadRequest, resp.Status)
	})

	t.Run("success - diferentes combinações de IDs", func(t *testing.T) {
		testCases := []struct {
			name      string
			userID    int64
			contactID int64
		}{
			{"IDs pequenos", 1, 2},
			{"IDs grandes", 999, 888},
			{"IDs diferentes", 123, 456},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				mockService := new(mockUser.MockUserContactRelation)
				handler := NewUserContactRelationHandler(mockService, logger)

				dtoRel := dto.UserContactRelationDTO{
					UserID:    *utils.Int64Ptr(tc.userID),
					ContactID: *utils.Int64Ptr(tc.contactID),
				}
				createdRel := &model.UserContactRelation{
					UserID:    tc.userID,
					ContactID: tc.contactID,
					CreatedAt: time.Now(),
				}

				mockService.On("Create", mock.Anything, mock.MatchedBy(func(rel *model.UserContactRelation) bool {
					return rel.UserID == tc.userID && rel.ContactID == tc.contactID
				})).Return(createdRel, nil).Once()

				body, _ := json.Marshal(dtoRel)
				req := httptest.NewRequest(http.MethodPost, "/contact-relations", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()

				handler.Create(rec, req)

				assert.Equal(t, http.StatusCreated, rec.Code)

				var resp utils.DefaultResponse
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				assert.Equal(t, "Relação criada com sucesso", resp.Message)

				mockService.AssertExpectations(t)
			})
		}
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
