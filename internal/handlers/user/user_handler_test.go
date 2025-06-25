package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	models_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	models_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	models_user "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	repository "github.com/WagaoCarvalho/backend_store_go/internal/repositories/users"
	user_services "github.com/WagaoCarvalho/backend_store_go/internal/services/user/user_services_mock"
	"github.com/WagaoCarvalho/backend_store_go/internal/utils"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func muxSetVars(req *http.Request, vars map[string]string) *http.Request {
	return mux.SetURLVars(req, vars)
}

func TestUserHandler_Create(t *testing.T) {
	mockService := new(user_services.MockUserService)
	handler := NewUserHandler(mockService)

	t.Run("Sucesso ao criar usuário", func(t *testing.T) {
		expectedUser := models_user.User{
			UID:      1,
			Username: "testuser",
			Email:    "test@example.com",
		}

		address := &models_address.Address{Street: "Rua 1"}
		contact := &models_contact.Contact{Phone: "12345"}

		requestBody := map[string]interface{}{
			"user": map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
			},
			"category_id": []int64{1, 2},
			"address":     address,
			"contact":     contact,
		}

		body, _ := json.Marshal(requestBody)

		mockService.On("Create",
			mock.Anything,
			mock.MatchedBy(func(u *models_user.User) bool {
				return u.Username == "testuser" && u.Email == "test@example.com"
			}),
			[]int64{1, 2},
			address,
			contact,
		).Return(&expectedUser, nil)

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
		userData := &models_user.User{Username: "failuser", Email: "fail@example.com"}
		categoryIDs := []int64{1}
		address := &models_address.Address{Street: "Fail St"}
		contact := &models_contact.Contact{Phone: "99999"}

		requestBody := map[string]interface{}{
			"user":        userData,
			"category_id": categoryIDs,
			"address":     address,
			"contact":     contact,
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("Create", mock.Anything, userData, categoryIDs, address, contact).
			Return(nil, errors.New("erro ao criar usuário"))

		req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.Create(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_GetAll(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(user_services.MockUserService)
		handler := NewUserHandler(mockService)

		users := []*models_user.User{{UID: 1, Username: "João"}}
		mockService.On("GetAll", mock.Anything).Return(users, nil)

		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()

		handler.GetAll(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.NewDecoder(rec.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.Equal(t, "Usuários encontrados", resp.Message)
		assert.Equal(t, http.StatusOK, resp.Status)

		dataBytes, _ := json.Marshal(resp.Data)
		var result []*models_user.User
		err = json.Unmarshal(dataBytes, &result)
		assert.NoError(t, err)
		assert.Equal(t, users, result)

		mockService.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockService := new(user_services.MockUserService)
		handler := NewUserHandler(mockService)

		mockService.On("GetAll", mock.Anything).Return(nil, errors.New("erro de banco"))

		req := httptest.NewRequest(http.MethodGet, "/users", nil)
		rec := httptest.NewRecorder()

		handler.GetAll(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var resp utils.DefaultResponse
		err := json.NewDecoder(rec.Body).Decode(&resp)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusInternalServerError, resp.Status)
		assert.Contains(t, resp.Message, "erro ao buscar usuários")
		assert.Nil(t, resp.Data)

		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_GetByID(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(user_services.MockUserService)
		handler := NewUserHandler(mockService)

		expectedUser := models_user.User{UID: 10, Username: "Carlos"}
		mockService.On("GetByID", mock.Anything, int64(10)).Return(&expectedUser, nil) // <- importante: retornando *User

		req := httptest.NewRequest(http.MethodGet, "/users/10", nil)
		req = muxSetVars(req, map[string]string{"id": "10"})
		rec := httptest.NewRecorder()

		handler.GetByID(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.NewDecoder(rec.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Status)
		assert.Equal(t, "Usuário encontrado", resp.Message)

		// Type assertion correta
		actualUser, ok := resp.Data.(map[string]interface{})
		assert.True(t, ok)

		uid := int64(actualUser["uid"].(float64))
		assert.Equal(t, expectedUser.UID, uid)

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockService := new(user_services.MockUserService)
		handler := NewUserHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/users/abc", nil)
		req = muxSetVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()

		handler.GetByID(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		resp, err := utils.ParseErrorResponse(rec.Body.Bytes())
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "ID inválido")
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService := new(user_services.MockUserService)
		handler := NewUserHandler(mockService)

		mockService.On("GetByID", mock.Anything, int64(99)).Return((*models_user.User)(nil), fmt.Errorf("usuário não encontrado"))

		req := httptest.NewRequest(http.MethodGet, "/users/99", nil)
		req = muxSetVars(req, map[string]string{"id": "99"})
		rec := httptest.NewRecorder()

		handler.GetByID(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		resp, err := utils.ParseErrorResponse(rec.Body.Bytes())
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.Status)
		assert.Contains(t, resp.Message, "usuário não encontrado")

		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_GetVersionByID(t *testing.T) {
	mockService := new(user_services.MockUserService)
	handler := NewUserHandler(mockService)

	router := mux.NewRouter()
	router.HandleFunc("/users/{id}/version", handler.GetVersionByID).Methods(http.MethodGet)

	t.Run("deve retornar a versão com sucesso", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/1/version", nil)
		rec := httptest.NewRecorder()

		mockService.On("GetVersionByID", mock.Anything, int64(1)).
			Return(int64(5), nil).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), `"version":5`)
		assert.Contains(t, rec.Body.String(), `"Versão do usuário obtida com sucesso`)
		mockService.AssertExpectations(t)
	})

	t.Run("deve retornar erro 400 para ID inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/abc/version", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "ID inválido")
	})

	t.Run("deve retornar erro 404 se usuário não encontrado", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/2/version", nil)
		rec := httptest.NewRecorder()

		mockService.On("GetVersionByID", mock.Anything, int64(2)).
			Return(int64(0), repository.ErrUserNotFound).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Contains(t, rec.Body.String(), "usuário não encontrado")
		mockService.AssertExpectations(t)
	})

	t.Run("deve retornar erro 500 para erro genérico", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/3/version", nil)
		rec := httptest.NewRecorder()

		mockService.On("GetVersionByID", mock.Anything, int64(3)).
			Return(int64(0), fmt.Errorf("erro no banco")).Once()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "erro no banco")
		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_GetByEmail(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(user_services.MockUserService)
		handler := NewUserHandler(mockService)

		expectedUser := &models_user.User{UID: 20, Email: "carlos@email.com"}
		mockService.On("GetByEmail", mock.Anything, "carlos@email.com").Return(expectedUser, nil)

		req := httptest.NewRequest(http.MethodGet, "/users/email/carlos@email.com", nil)
		req = muxSetVars(req, map[string]string{"email": "carlos@email.com"})
		rec := httptest.NewRecorder()

		handler.GetByEmail(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.NewDecoder(rec.Body).Decode(&resp)
		assert.NoError(t, err)

		dataMap := resp.Data.(map[string]interface{})
		assert.Equal(t, float64(expectedUser.UID), dataMap["uid"])
		mockService.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		mockService := new(user_services.MockUserService)
		handler := NewUserHandler(mockService)

		mockService.On("GetByEmail", mock.Anything, "naoexiste@email.com").
			Return(nil, fmt.Errorf("usuário não encontrado"))

		req := httptest.NewRequest(http.MethodGet, "/users/email/naoexiste@email.com", nil)
		req = muxSetVars(req, map[string]string{"email": "naoexiste@email.com"})
		rec := httptest.NewRecorder()

		handler.GetByEmail(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		resp, err := utils.ParseErrorResponse(rec.Body.Bytes())
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.Status)
		assert.Contains(t, resp.Message, "usuário não encontrado")

		mockService.AssertExpectations(t)
	})

	t.Run("InternalError", func(t *testing.T) {
		mockService := new(user_services.MockUserService)
		handler := NewUserHandler(mockService)

		mockService.On("GetByEmail", mock.Anything, "erro@email.com").
			Return(nil, fmt.Errorf("erro inesperado no banco"))

		req := httptest.NewRequest(http.MethodGet, "/users/email/erro@email.com", nil)
		req = muxSetVars(req, map[string]string{"email": "erro@email.com"})
		rec := httptest.NewRecorder()

		handler.GetByEmail(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		resp, err := utils.ParseErrorResponse(rec.Body.Bytes())
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)
		assert.Contains(t, resp.Message, "erro inesperado")

		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_Update(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(user_services.MockUserService)
		handler := NewUserHandler(mockService)

		inputUser := &models_user.User{
			Username: "Atualizado",
			Address: &models_address.Address{
				Street: "Rua Nova",
			},
		}
		updatedUser := &models_user.User{UID: 1, Username: "Atualizado"}

		mockService.On("Update",
			mock.Anything,
			mock.MatchedBy(func(u *models_user.User) bool {
				return u.UID == 1 && u.Username == "Atualizado" && u.Address != nil && u.Address.Street == "Rua Nova"
			}),
			mock.Anything, // Address extraído dentro do handler do próprio u.Address
		).Return(updatedUser, nil)

		body := map[string]interface{}{
			"user": inputUser,
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(jsonBody))
		req = muxSetVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.NewDecoder(rec.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Status)
		assert.Equal(t, "Usuário atualizado com sucesso", resp.Message)

		mockService.AssertExpectations(t)
	})

	t.Run("UserDataNil", func(t *testing.T) {
		mockService := new(user_services.MockUserService)
		handler := NewUserHandler(mockService)

		body := map[string]interface{}{
			"user": nil, // user nil para simular o erro
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader(jsonBody))
		req = muxSetVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		resp, err := utils.ParseErrorResponse(rec.Body.Bytes())
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "dados do usuário são obrigatórios")

		// O serviço NÃO deve ser chamado
		mockService.AssertExpectations(t)
	})

	t.Run("InvalidMethod", func(t *testing.T) {
		mockService := new(user_services.MockUserService)
		handler := NewUserHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/users/1", nil) // GET ao invés de PUT
		req = muxSetVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)

		resp, err := utils.ParseErrorResponse(rec.Body.Bytes())
		assert.NoError(t, err)
		assert.Equal(t, http.StatusMethodNotAllowed, resp.Status)
		assert.Contains(t, resp.Message, "método GET não permitido")
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockService := new(user_services.MockUserService)
		handler := NewUserHandler(mockService)

		body := map[string]interface{}{"user": map[string]interface{}{"username": "teste"}}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPut, "/users/abc", bytes.NewReader(jsonBody))
		req = muxSetVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		resp, err := utils.ParseErrorResponse(rec.Body.Bytes())
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
	})

	t.Run("InvalidBody", func(t *testing.T) {
		mockService := new(user_services.MockUserService)
		handler := NewUserHandler(mockService)

		req := httptest.NewRequest(http.MethodPut, "/users/1", bytes.NewReader([]byte(`invalid-json`)))
		req = muxSetVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		resp, err := utils.ParseErrorResponse(rec.Body.Bytes())
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
	})

	t.Run("VersionConflict", func(t *testing.T) {
		mockService := new(user_services.MockUserService)
		handler := NewUserHandler(mockService)

		user := &models_user.User{Username: "Carlos"}

		mockService.On("Update",
			mock.Anything,
			mock.MatchedBy(func(u *models_user.User) bool {
				return u.UID == 2 && u.Username == "Carlos"
			}),
			mock.Anything,
		).Return(nil, repository.ErrVersionConflict)

		body := map[string]interface{}{
			"user": user,
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPut, "/users/2", bytes.NewReader(jsonBody))
		req = muxSetVars(req, map[string]string{"id": "2"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)

		resp, err := utils.ParseErrorResponse(rec.Body.Bytes())
		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, resp.Status)
		assert.Contains(t, resp.Message, "versão")

		mockService.AssertExpectations(t)
	})

	t.Run("InternalError", func(t *testing.T) {
		mockService := new(user_services.MockUserService)
		handler := NewUserHandler(mockService)

		user := &models_user.User{UID: 3, Username: "Carlos"}

		mockService.On("Update",
			mock.Anything,
			mock.MatchedBy(func(u *models_user.User) bool {
				return u.UID == 3 && u.Username == "Carlos"
			}),
			(*models_address.Address)(nil),
		).Return(nil, errors.New("erro inesperado"))

		body := map[string]interface{}{
			"user": user,
		}
		jsonBody, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPut, "/users/3", bytes.NewReader(jsonBody))
		req = muxSetVars(req, map[string]string{"id": "3"})
		rec := httptest.NewRecorder()

		handler.Update(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		resp, err := utils.ParseErrorResponse(rec.Body.Bytes())
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)
		assert.Contains(t, resp.Message, "erro inesperado")

		mockService.AssertExpectations(t)
	})
}

func TestUserHandler_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockService := new(user_services.MockUserService)
		handler := NewUserHandler(mockService)

		mockService.On("Delete", mock.Anything, int64(1)).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/users/1", nil)
		req = muxSetVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.NewDecoder(rec.Body).Decode(&resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Status)
		assert.Equal(t, "Usuário deletado com sucesso", resp.Message)
		assert.Nil(t, resp.Data)

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidID", func(t *testing.T) {
		mockService := new(user_services.MockUserService)
		handler := NewUserHandler(mockService)

		req := httptest.NewRequest(http.MethodDelete, "/users/abc", nil)
		req = muxSetVars(req, map[string]string{"id": "abc"})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		resp, err := utils.ParseErrorResponse(rec.Body.Bytes())
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "ID inválido")
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockService := new(user_services.MockUserService)
		handler := NewUserHandler(mockService)

		mockService.On("Delete", mock.Anything, int64(99)).
			Return(fmt.Errorf("usuário não encontrado"))

		req := httptest.NewRequest(http.MethodDelete, "/users/99", nil)
		req = muxSetVars(req, map[string]string{"id": "99"})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		resp, err := utils.ParseErrorResponse(rec.Body.Bytes())
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.Status)
		assert.Contains(t, resp.Message, "usuário não encontrado")

		mockService.AssertExpectations(t)
	})

	t.Run("InternalError", func(t *testing.T) {
		mockService := new(user_services.MockUserService)
		handler := NewUserHandler(mockService)

		mockService.On("Delete", mock.Anything, int64(2)).
			Return(errors.New("erro interno"))

		req := httptest.NewRequest(http.MethodDelete, "/users/2", nil)
		req = muxSetVars(req, map[string]string{"id": "2"})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		resp, err := utils.ParseErrorResponse(rec.Body.Bytes())
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)
		assert.Contains(t, resp.Message, "erro interno")

		mockService.AssertExpectations(t)
	})

	t.Run("InvalidMethod", func(t *testing.T) {
		mockService := new(user_services.MockUserService)
		handler := NewUserHandler(mockService)

		req := httptest.NewRequest(http.MethodGet, "/users/1", nil) // GET em vez de DELETE
		req = muxSetVars(req, map[string]string{"id": "1"})
		rec := httptest.NewRecorder()

		handler.Delete(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)

		resp, err := utils.ParseErrorResponse(rec.Body.Bytes())
		assert.NoError(t, err)
		assert.Equal(t, http.StatusMethodNotAllowed, resp.Status)
		assert.Contains(t, resp.Message, "método GET não permitido")
	})
}
