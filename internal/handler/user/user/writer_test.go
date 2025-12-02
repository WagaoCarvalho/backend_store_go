package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockUser "github.com/WagaoCarvalho/backend_store_go/infra/mock/user"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/user"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserHandler_Create(t *testing.T) {
	mockService := new(mockUser.MockUser)
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)
	handler := NewUserHandler(mockService, logger)

	t.Run("Sucesso ao criar usuário", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		expectedUser := &model.User{
			UID:      1,
			Username: "testuser",
			Email:    "test@example.com",
		}

		requestDTO := &dto.UserDTO{
			Username: "testuser",
			Email:    "test@example.com",
		}

		requestBody := map[string]interface{}{
			"user": requestDTO,
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("Create",
			mock.Anything,
			mock.MatchedBy(func(u *model.User) bool {
				return u.Username == requestDTO.Username && u.Email == requestDTO.Email
			}),
		).Return(expectedUser, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro método não permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("Erro ao decodificar JSON inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte("{invalid json")))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro ao criar usuário no service", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		requestDTO := &dto.UserDTO{
			Username: "failuser",
			Email:    "fail@example.com",
		}

		requestBody := map[string]interface{}{
			"user": requestDTO,
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("Create",
			mock.Anything,
			mock.MatchedBy(func(u *model.User) bool {
				return u.Username == requestDTO.Username && u.Email == requestDTO.Email
			}),
		).Return(nil, errors.New("erro ao criar usuário")).Once()

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_Update(t *testing.T) {
	mockService := new(mockUser.MockUser)
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)
	handler := NewUserHandler(mockService, logger)

	t.Run("Sucesso ao atualizar usuário", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		userID := int64(1)
		requestDTO := dto.UserDTO{
			Username: "updatedUser",
			Email:    "updated@example.com",
			Version:  2,
		}

		body, _ := json.Marshal(map[string]interface{}{"user": requestDTO})

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.UID == userID
		})).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro método não permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/users/1", nil)
		rec := httptest.NewRecorder()

		handler.Update(rec, req)
		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("Erro ID inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/users/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro JSON inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader([]byte("{invalid json")))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro dados do usuário ausentes", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		body, _ := json.Marshal(map[string]interface{}{"user": nil})
		req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro conflito de versão", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		userID := int64(1)
		requestDTO := dto.UserDTO{Version: 2}
		body, _ := json.Marshal(map[string]interface{}{"user": requestDTO})

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.UID == userID && u.Version == 2
		})).Return(errMsg.ErrZeroVersion).Once()

		req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)
		assert.Equal(t, http.StatusConflict, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro dados inválidos (email incorreto)", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		requestDTO := dto.UserDTO{
			Email:   "invalid-email",
			Version: 2,
		}
		body, _ := json.Marshal(map[string]interface{}{"user": requestDTO})

		mockService.On("Update", mock.Anything, mock.Anything).Return(errMsg.ErrInvalidData).Once()

		req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro usuário não encontrado", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		userID := int64(999)
		requestDTO := dto.UserDTO{Version: 1}
		body, _ := json.Marshal(map[string]interface{}{"user": requestDTO})

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.UID == userID
		})).Return(errMsg.ErrNotFound).Once()

		req := httptest.NewRequest(http.MethodPut, "/users/999", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)
		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro genérico ao atualizar usuário", func(t *testing.T) {
		mockService.ExpectedCalls = nil
		requestDTO := dto.UserDTO{Version: 2}
		body, _ := json.Marshal(map[string]interface{}{"user": requestDTO})

		mockService.On("Update", mock.Anything, mock.Anything).Return(errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_Delete(t *testing.T) {
	mockService := new(mockUser.MockUser)
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)
	handler := NewUserHandler(mockService, logger)

	t.Run("Sucesso ao deletar usuário", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("Delete", mock.Anything, int64(1)).Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.Empty(t, rec.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("Erro método não permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("Erro ID inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/users/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro usuário não encontrado", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("Delete", mock.Anything, int64(999)).Return(errors.New("usuário não encontrado")).Once()

		req := httptest.NewRequest(http.MethodDelete, "/users/999", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro genérico ao deletar usuário", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("Delete", mock.Anything, int64(2)).Return(errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodDelete, "/users/2", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}
