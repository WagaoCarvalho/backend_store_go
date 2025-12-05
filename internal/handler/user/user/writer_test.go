package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	mockUser "github.com/WagaoCarvalho/backend_store_go/infra/mock/user"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/user"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	errmsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
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

		userModel := &model.User{
			UID:         1,
			Username:    "testuser",
			Email:       "test@example.com",
			Password:    "hashedpassword",
			Description: "Test user",
			Status:      true,
			Version:     1,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		mockService.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).
			Run(func(args mock.Arguments) {
				userArg := args.Get(1).(*model.User)
				*userArg = *userModel
			}).
			Return(userModel, nil).Once()

		userDTO := dto.UserDTO{
			Username:    "testuser",
			Email:       "test@example.com",
			Password:    "Password123",
			Description: "Test user",
			Status:      true,
		}

		body, _ := json.Marshal(userDTO)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, float64(http.StatusCreated), response["status"])
		assert.Equal(t, "Usuário criado com sucesso", response["message"])
		assert.NotNil(t, response["data"])

		mockService.AssertExpectations(t)
	})

	t.Run("Erro método não permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("Erro ao parsear JSON inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader("{invalid json"))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro validação - usuário já existe no sistema", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).
			Return(nil, errors.New(errmsg.ErrInvalidData.Error()+": dados já existem no sistema")).Once()

		userDTO := dto.UserDTO{
			Username: "existinguser",
			Email:    "existing@example.com",
			Password: "Password123",
		}

		body, _ := json.Marshal(userDTO)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro serviço - erro genérico de criação", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).
			Return(nil, errors.New(errmsg.ErrCreate.Error()+": erro no banco de dados")).Once()

		userDTO := dto.UserDTO{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "Password123",
		}

		body, _ := json.Marshal(userDTO)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro serviço - usuário nulo retornado", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).
			Return(nil, nil).Once()

		userDTO := dto.UserDTO{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "Password123",
		}

		body, _ := json.Marshal(userDTO)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro serviço - erro interno no hashing", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).
			Return(nil, errors.New("erro interno do servidor: erro ao hashear senha")).Once()

		userDTO := dto.UserDTO{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "Password123",
		}

		body, _ := json.Marshal(userDTO)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Sucesso com campos opcionais omitidos", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		userModel := &model.User{
			UID:       2,
			Username:  "simpleuser",
			Email:     "simple@example.com",
			Password:  "hashedpassword",
			Status:    true,
			Version:   1,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		mockService.On("Create", mock.Anything, mock.AnythingOfType("*model.User")).
			Run(func(args mock.Arguments) {
				userArg := args.Get(1).(*model.User)
				*userArg = *userModel
			}).
			Return(userModel, nil).Once()

		userDTO := dto.UserDTO{
			Username: "simpleuser",
			Email:    "simple@example.com",
			Password: "Password123",
			Status:   true,
		}

		body, _ := json.Marshal(userDTO)
		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewBuffer(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, float64(http.StatusCreated), response["status"])
		assert.Equal(t, "Usuário criado com sucesso", response["message"])

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

		mockService.On("Update", mock.Anything, mock.AnythingOfType("*model.User")).
			Return(nil).Once()

		userDTO := dto.UserDTO{
			Username:    "updateduser",
			Email:       "updated@example.com",
			Description: "Updated description",
			Status:      true,
			Version:     1,
		}

		body, _ := json.Marshal(userDTO)
		req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, float64(http.StatusOK), response["status"])
		assert.Equal(t, "Usuário atualizado com sucesso", response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("Erro método não permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/users/1", nil)
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("Erro ID inválido", func(t *testing.T) {
		userDTO := dto.UserDTO{
			Username: "testuser",
			Email:    "test@example.com",
		}

		body, _ := json.Marshal(userDTO)
		req := httptest.NewRequest(http.MethodPut, "/users/abc", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro ao parsear JSON inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBufferString("{invalid json"))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro - conflito de versão", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		// CORREÇÃO: Usar fmt.Errorf com %w para que errors.Is() funcione
		err := fmt.Errorf("%w", errmsg.ErrVersionConflict)
		mockService.On("Update", mock.Anything, mock.AnythingOfType("*model.User")).
			Return(err).Once()

		userDTO := dto.UserDTO{
			Username: "testuser",
			Email:    "test@example.com",
			Version:  1,
		}

		body, _ := json.Marshal(userDTO)
		req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Contains(t, response["message"], "versão desatualizada")

		mockService.AssertExpectations(t)
	})

	t.Run("Erro - dados inválidos", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		// CORREÇÃO: Usar fmt.Errorf com %w
		err := fmt.Errorf("%w: %v", errmsg.ErrInvalidData, "email inválido")
		mockService.On("Update", mock.Anything, mock.AnythingOfType("*model.User")).
			Return(err).Once()

		userDTO := dto.UserDTO{
			Username: "", // Username vazio
			Email:    "invalid-email",
		}

		body, _ := json.Marshal(userDTO)
		req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Contains(t, response["message"], "email inválido")

		mockService.AssertExpectations(t)
	})

	t.Run("Erro - usuário não encontrado", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		// CORREÇÃO: Usar fmt.Errorf com %w
		err := fmt.Errorf("%w", errmsg.ErrNotFound)
		mockService.On("Update", mock.Anything, mock.AnythingOfType("*model.User")).
			Return(err).Once()

		userDTO := dto.UserDTO{
			Username: "nonexistent",
			Email:    "nonexistent@example.com",
		}

		body, _ := json.Marshal(userDTO)
		req := httptest.NewRequest(http.MethodPut, "/users/999", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "999"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Contains(t, response["message"], "usuário não encontrado")

		mockService.AssertExpectations(t)
	})

	t.Run("Erro - erro interno do servidor", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		// Este erro não deve ser um dos erros conhecidos (não usa %w com os erros definidos)
		err := errors.New("erro no banco de dados")
		mockService.On("Update", mock.Anything, mock.AnythingOfType("*model.User")).
			Return(err).Once()

		userDTO := dto.UserDTO{
			Username: "testuser",
			Email:    "test@example.com",
		}

		body, _ := json.Marshal(userDTO)
		req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Contains(t, response["message"], "erro interno do servidor")

		mockService.AssertExpectations(t)
	})

	t.Run("Sucesso - atualização sem password (opcional no update)", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("Update", mock.Anything, mock.AnythingOfType("*model.User")).
			Return(nil).Once()

		userDTO := dto.UserDTO{
			Username:    "testuser",
			Email:       "test@example.com",
			Description: "Updated without password",
			Status:      true,
			Version:     2,
			// Password não enviado - opcional no update
		}

		body, _ := json.Marshal(userDTO)
		req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, float64(http.StatusOK), response["status"])
		assert.Equal(t, "Usuário atualizado com sucesso", response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("Sucesso - atualização com password (se fornecido)", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		mockService.On("Update", mock.Anything, mock.AnythingOfType("*model.User")).
			Return(nil).Once()

		userDTO := dto.UserDTO{
			Username:    "testuser",
			Email:       "test@example.com",
			Password:    "NewPassword123", // Password fornecido
			Description: "Updated with new password",
			Status:      true,
			Version:     3,
		}

		body, _ := json.Marshal(userDTO)
		req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		assert.Equal(t, float64(http.StatusOK), response["status"])
		assert.Equal(t, "Usuário atualizado com sucesso", response["message"])

		mockService.AssertExpectations(t)
	})

	t.Run("Erro - dados inválidos sem mensagem adicional", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		// Erro sem mensagem adicional após o prefixo
		err := fmt.Errorf("%w", errmsg.ErrInvalidData)
		mockService.On("Update", mock.Anything, mock.AnythingOfType("*model.User")).
			Return(err).Once()

		userDTO := dto.UserDTO{
			Username: "testuser",
			Email:    "test@example.com",
		}

		body, _ := json.Marshal(userDTO)
		req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var response map[string]interface{}
		json.Unmarshal(rec.Body.Bytes(), &response)
		// Como não há mensagem após o prefixo, o handler retorna a string completa do erro
		assert.NotEmpty(t, response["message"])

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
