package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	model "github.com/WagaoCarvalho/backend_store_go/internal/model/user"
	repo "github.com/WagaoCarvalho/backend_store_go/internal/repo/user/user"
	service "github.com/WagaoCarvalho/backend_store_go/internal/service/user/user_services_mock"
	"github.com/WagaoCarvalho/backend_store_go/logger"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserHandler_Create(t *testing.T) {
	mockService := new(service.MockUserService)
	logger := logger.NewLoggerAdapter(logrus.New())
	handler := NewUserHandler(mockService, logger)

	t.Run("Sucesso ao criar usuário", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		expectedUser := &model.User{
			UID:      1,
			Username: "testuser",
			Email:    "test@example.com",
		}

		requestBody := map[string]interface{}{
			"user": map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
			},
		}

		body, _ := json.Marshal(requestBody)

		mockService.On("Create",
			mock.Anything,
			mock.MatchedBy(func(u *model.User) bool {
				return u.Username == "testuser" && u.Email == "test@example.com"
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

		userData := &model.User{
			Username: "failuser",
			Email:    "fail@example.com",
		}

		requestBody := map[string]interface{}{
			"user": userData,
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("Create", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.Username == userData.Username && u.Email == userData.Email
		})).Return(nil, errors.New("erro ao criar usuário")).Once()

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_GetAll(t *testing.T) {
	mockService := new(service.MockUserService)
	logger := logger.NewLoggerAdapter(logrus.New())
	handler := NewUserHandler(mockService, logger)

	t.Run("Sucesso ao buscar todos usuários", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		users := []*model.User{
			{UID: 1, Username: "user1", Email: "user1@example.com"},
			{UID: 2, Username: "user2", Email: "user2@example.com"},
		}

		mockService.On("GetAll", mock.Anything).Return(users, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()

		handler.GetAll(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro ao buscar usuários", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetAll", mock.Anything).Return(nil, errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()

		handler.GetAll(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_GetByID(t *testing.T) {
	mockService := new(service.MockUserService)
	logger := logger.NewLoggerAdapter(logrus.New())
	handler := NewUserHandler(mockService, logger)

	t.Run("Sucesso ao buscar usuário por ID", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		user := &model.User{
			UID:      1,
			Username: "user1",
			Email:    "user1@example.com",
		}

		mockService.On("GetByID", mock.Anything, int64(1)).Return(user, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/users/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.GetByID(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro ID inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()

		handler.GetByID(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro usuário não encontrado", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByID", mock.Anything, int64(999)).Return(nil, errors.New("usuário não encontrado")).Once()

		req := httptest.NewRequest(http.MethodGet, "/users/999", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		rec := httptest.NewRecorder()

		handler.GetByID(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro genérico ao buscar usuário", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByID", mock.Anything, int64(2)).Return(nil, errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodGet, "/users/2", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		rec := httptest.NewRecorder()

		handler.GetByID(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_GetVersionByID(t *testing.T) {
	mockService := new(service.MockUserService)
	logger := logger.NewLoggerAdapter(logrus.New())
	handler := NewUserHandler(mockService, logger)

	t.Run("Sucesso ao obter versão por ID", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetVersionByID", mock.Anything, int64(1)).Return(int64(5), nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/users/1/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.GetVersionByID(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro ID inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/abc/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()

		handler.GetVersionByID(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro usuário não encontrado", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetVersionByID", mock.Anything, int64(999)).Return(int64(0), repo.ErrUserNotFound).Once()

		req := httptest.NewRequest(http.MethodGet, "/users/999/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		rec := httptest.NewRecorder()

		handler.GetVersionByID(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro genérico ao obter versão", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetVersionByID", mock.Anything, int64(2)).Return(int64(0), errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodGet, "/users/2/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		rec := httptest.NewRecorder()

		handler.GetVersionByID(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_GetByEmail(t *testing.T) {
	mockService := new(service.MockUserService)
	logger := logger.NewLoggerAdapter(logrus.New())
	handler := NewUserHandler(mockService, logger)

	t.Run("Sucesso ao buscar usuário por email", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		user := &model.User{
			UID:      1,
			Username: "user1",
			Email:    "user1@example.com",
		}

		mockService.On("GetByEmail", mock.Anything, "user1@example.com").Return(user, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/users/email/user1@example.com", nil)
		req = mux.SetURLVars(req, map[string]string{"email": "user1@example.com"})
		rec := httptest.NewRecorder()

		handler.GetByEmail(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro usuário não encontrado", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByEmail", mock.Anything, "notfound@example.com").Return(nil, errors.New("usuário não encontrado")).Once()

		req := httptest.NewRequest(http.MethodGet, "/users/email/notfound@example.com", nil)
		req = mux.SetURLVars(req, map[string]string{"email": "notfound@example.com"})
		rec := httptest.NewRecorder()

		handler.GetByEmail(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro genérico ao buscar usuário", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByEmail", mock.Anything, "error@example.com").Return(nil, errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodGet, "/users/email/error@example.com", nil)
		req = mux.SetURLVars(req, map[string]string{"email": "error@example.com"})
		rec := httptest.NewRecorder()

		handler.GetByEmail(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_GetByName(t *testing.T) {
	mockService := new(service.MockUserService)
	logger := logger.NewLoggerAdapter(logrus.New())
	handler := NewUserHandler(mockService, logger)

	t.Run("Sucesso ao buscar usuários por nome parcial", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		users := []*model.User{
			{
				UID:      1,
				Username: "user1",
				Email:    "user1@example.com",
			},
			{
				UID:      2,
				Username: "user123",
				Email:    "user123@example.com",
			},
		}

		mockService.On("GetByName", mock.Anything, "user1").Return(users, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/users/name/user1", nil)
		req = mux.SetURLVars(req, map[string]string{"username": "user1"})
		rec := httptest.NewRecorder()

		handler.GetByName(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro usuário não encontrado", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByName", mock.Anything, "notfound").Return(nil, errors.New("usuário não encontrado")).Once()

		req := httptest.NewRequest(http.MethodGet, "/users/name/notfound", nil)
		req = mux.SetURLVars(req, map[string]string{"username": "notfound"})
		rec := httptest.NewRecorder()

		handler.GetByName(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro genérico ao buscar usuário por nome", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByName", mock.Anything, "error").Return(nil, errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodGet, "/users/name/error", nil)
		req = mux.SetURLVars(req, map[string]string{"username": "error"})
		rec := httptest.NewRecorder()

		handler.GetByName(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_Update(t *testing.T) {
	mockService := new(service.MockUserService)
	logger := logger.NewLoggerAdapter(logrus.New())
	handler := NewUserHandler(mockService, logger)

	t.Run("Sucesso ao atualizar usuário", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		userID := int64(1)
		userToUpdate := &model.User{
			UID:      userID,
			Username: "updatedUser",
			Email:    "updated@example.com",
			Version:  2,
		}

		requestBody := map[string]interface{}{
			"user": map[string]interface{}{
				"username": "updatedUser",
				"email":    "updated@example.com",
				"version":  2,
			},
		}

		body, _ := json.Marshal(requestBody)

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.UID == userID && u.Username == "updatedUser" && u.Version == 2
		})).Return(userToUpdate, nil).Once()

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
		requestBody := map[string]interface{}{
			"user": nil,
		}
		body, _ := json.Marshal(requestBody)

		req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro conflito de versão", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		userID := int64(1)

		requestBody := map[string]interface{}{
			"user": map[string]interface{}{
				"version": 2,
			},
		}

		body, _ := json.Marshal(requestBody)

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.UID == userID && u.Version == 2
		})).Return(nil, repo.ErrVersionConflict).Once()

		req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro genérico ao atualizar usuário", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		userID := int64(1)

		requestBody := map[string]interface{}{
			"user": map[string]interface{}{
				"version": 2,
			},
		}

		body, _ := json.Marshal(requestBody)

		mockService.On("Update", mock.Anything, mock.MatchedBy(func(u *model.User) bool {
			return u.UID == userID && u.Version == 2
		})).Return(nil, errors.New("erro interno")).Once()

		req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_Disable(t *testing.T) {
	mockService := new(service.MockUserService)
	logger := logger.NewLoggerAdapter(logrus.New())
	handler := NewUserHandler(mockService, logger)

	t.Run("Sucesso ao desabilitar usuário", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		// Mock do GetByID para buscar o usuário antes de atualizar
		mockService.On("GetByID", mock.Anything, int64(1)).Return(&model.User{
			UID:     1,
			Status:  true,
			Version: 5,
		}, nil).Once()

		// Mock do Update com usuário atualizado
		mockService.On("Update", mock.Anything, mock.MatchedBy(func(user *model.User) bool {
			return user.UID == 1 && user.Status == false && user.Version == 10
		})).Return(&model.User{
			UID:     1,
			Status:  false,
			Version: 10,
		}, nil).Once()

		body := `{"version": 10}`
		req := httptest.NewRequest(http.MethodPatch, "/users/1/disable", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.Empty(t, rec.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("Erro método não permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/1/disable", nil)
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("Erro ID inválido", func(t *testing.T) {
		body := `{"version": 1}`
		req := httptest.NewRequest(http.MethodPatch, "/users/abc/disable", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro versão inválida no corpo", func(t *testing.T) {
		body := `{}` // sem versão
		req := httptest.NewRequest(http.MethodPatch, "/users/1/disable", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro usuário não encontrado", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByID", mock.Anything, int64(999)).Return(nil, errors.New("usuário não encontrado")).Once()

		body := `{"version": 1}`
		req := httptest.NewRequest(http.MethodPatch, "/users/999/disable", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro conflito de versão ao desabilitar usuário", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByID", mock.Anything, int64(4)).Return(&model.User{
			UID:     4,
			Status:  true,
			Version: 2,
		}, nil).Once()

		mockService.On("Update", mock.Anything, mock.Anything).Return(nil, repo.ErrVersionConflict).Once()

		body := `{"version": 2}`
		req := httptest.NewRequest(http.MethodPatch, "/users/4/disable", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "4"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro genérico ao desabilitar usuário", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByID", mock.Anything, int64(2)).Return(&model.User{
			UID:     2,
			Status:  true,
			Version: 1,
		}, nil).Once()

		mockService.On("Update", mock.Anything, mock.Anything).Return(nil, errors.New("erro interno")).Once()

		body := `{"version": 1}`
		req := httptest.NewRequest(http.MethodPatch, "/users/2/disable", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro genérico ao buscar usuário", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByID", mock.Anything, int64(10)).Return(nil, errors.New("erro interno")).Once()

		body := `{"version": 1}` // JSON válido para passar na validação do decode
		req := httptest.NewRequest(http.MethodPatch, "/users/10/disable", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "10"})
		rec := httptest.NewRecorder()

		handler.Disable(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_Enable(t *testing.T) {
	mockService := new(service.MockUserService)
	logger := logger.NewLoggerAdapter(logrus.New())
	handler := NewUserHandler(mockService, logger)

	t.Run("Sucesso ao habilitar usuário", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		// Mock GetByID para buscar usuário atual
		mockService.On("GetByID", mock.Anything, int64(1)).Return(&model.User{
			UID:     1,
			Status:  false,
			Version: 3,
		}, nil).Once()

		// Mock Update com usuário atualizado
		mockService.On("Update", mock.Anything, mock.MatchedBy(func(user *model.User) bool {
			return user.UID == 1 && user.Status == true && user.Version == 4
		})).Return(&model.User{
			UID:     1,
			Status:  true,
			Version: 4,
		}, nil).Once()

		body := `{"version": 4}`
		req := httptest.NewRequest(http.MethodPatch, "/users/1/enable", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusNoContent, rec.Code)
		assert.Empty(t, rec.Body.String())
		mockService.AssertExpectations(t)
	})

	t.Run("Erro método não permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/1/enable", nil)
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("Erro ID inválido", func(t *testing.T) {
		body := `{"version": 1}`
		req := httptest.NewRequest(http.MethodPatch, "/users/abc/enable", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro conflito de versão ao habilitar usuário", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		// GetByID retorna usuário válido
		mockService.On("GetByID", mock.Anything, int64(3)).Return(&model.User{
			UID:     3,
			Status:  false,
			Version: 5,
		}, nil).Once()

		// Update retorna erro de conflito de versão
		mockService.On("Update", mock.Anything, mock.Anything).Return(nil, repo.ErrVersionConflict).Once()

		body := `{"version": 5}`
		req := httptest.NewRequest(http.MethodPatch, "/users/3/enable", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "3"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro versão inválida no corpo", func(t *testing.T) {
		body := `{}` // sem version
		req := httptest.NewRequest(http.MethodPatch, "/users/1/enable", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro usuário não encontrado", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByID", mock.Anything, int64(999)).Return(nil, errors.New("usuário não encontrado")).Once()

		body := `{"version": 1}`
		req := httptest.NewRequest(http.MethodPatch, "/users/999/enable", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro genérico ao habilitar usuário", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByID", mock.Anything, int64(2)).Return(&model.User{
			UID:     2,
			Status:  false,
			Version: 1,
		}, nil).Once()

		mockService.On("Update", mock.Anything, mock.Anything).Return(nil, errors.New("erro interno")).Once()

		body := `{"version": 1}`
		req := httptest.NewRequest(http.MethodPatch, "/users/2/enable", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "2"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro genérico ao buscar usuário", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("GetByID", mock.Anything, int64(10)).Return(nil, errors.New("erro interno")).Once()

		body := `{"version": 1}` // JSON válido para passar na validação do decode
		req := httptest.NewRequest(http.MethodPatch, "/users/10/enable", strings.NewReader(body))
		req = mux.SetURLVars(req, map[string]string{"id": "10"})
		rec := httptest.NewRecorder()

		handler.Enable(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})

}

func TestUserHandler_Delete(t *testing.T) {
	mockService := new(service.MockUserService)
	logger := logger.NewLoggerAdapter(logrus.New())
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
