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
	filter "github.com/WagaoCarvalho/backend_store_go/internal/model/user/filter"
	model "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserHandler_Filter(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	log := logger.NewLoggerAdapter(baseLogger)

	t.Run("sucesso - retorna lista de usuários", func(t *testing.T) {
		mockService := new(mockUser.MockUser)
		handler := NewUserFilterHandler(mockService, log)

		now := time.Now()
		description1 := "Usuário administrador"
		description2 := ""

		mockUsers := []*model.User{
			{
				UID:         1,
				Username:    "admin",
				Email:       "admin@example.com",
				Password:    "hashed_password_1",
				Description: description1,
				Status:      true,
				Version:     1,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
			{
				UID:         2,
				Username:    "guest",
				Email:       "guest@example.com",
				Password:    "hashed_password_2",
				Description: description2,
				Status:      false,
				Version:     1,
				CreatedAt:   now,
				UpdatedAt:   now,
			},
		}

		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return(mockUsers, nil).
			Once()

		req := httptest.NewRequest(
			http.MethodGet,
			"/users/filter?limit=10&offset=0",
			nil,
		)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.Status)
		assert.Equal(t, "Usuários listados com sucesso", resp.Message)

		data := resp.Data.(map[string]any)
		assert.Equal(t, float64(2), data["total"])

		items := data["items"].([]any)
		assert.Len(t, items, 2)

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - parse do status booleano true", func(t *testing.T) {
		mockService := new(mockUser.MockUser)
		handler := NewUserFilterHandler(mockService, log)

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f any) bool {
				filter, ok := f.(*filter.UserFilter)
				if !ok {
					return false
				}
				return filter.Status != nil && *filter.Status == true
			})).
			Return([]*model.User{}, nil).
			Once()

		req := httptest.NewRequest(
			http.MethodGet,
			"/users/filter?status=true&limit=10&offset=0",
			nil,
		)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - parse do status booleano false", func(t *testing.T) {
		mockService := new(mockUser.MockUser)
		handler := NewUserFilterHandler(mockService, log)

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f any) bool {
				filter, ok := f.(*filter.UserFilter)
				if !ok {
					return false
				}
				return filter.Status != nil && *filter.Status == false
			})).
			Return([]*model.User{}, nil).
			Once()

		req := httptest.NewRequest(
			http.MethodGet,
			"/users/filter?status=false&limit=10&offset=0",
			nil,
		)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - parse dos parâmetros de paginação padrão", func(t *testing.T) {
		mockService := new(mockUser.MockUser)
		handler := NewUserFilterHandler(mockService, log)

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f any) bool {
				filter, ok := f.(*filter.UserFilter)
				if !ok {
					return false
				}
				// Verifica se os valores padrão foram aplicados
				return filter.Limit == 10 && filter.Offset == 0
			})).
			Return([]*model.User{}, nil).
			Once()

		req := httptest.NewRequest(
			http.MethodGet,
			"/users/filter",
			nil,
		)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - filtro por username e email", func(t *testing.T) {
		mockService := new(mockUser.MockUser)
		handler := NewUserFilterHandler(mockService, log)

		username := "john"
		email := "john@example.com"

		mockService.
			On("Filter", mock.Anything, mock.MatchedBy(func(f any) bool {
				filter, ok := f.(*filter.UserFilter)
				if !ok {
					return false
				}
				return filter.Username == username && filter.Email == email
			})).
			Return([]*model.User{}, nil).
			Once()

		req := httptest.NewRequest(
			http.MethodGet,
			"/users/filter?username=john&email=john@example.com&limit=10&offset=0",
			nil,
		)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("erro - parâmetro desconhecido na query", func(t *testing.T) {
		mockService := new(mockUser.MockUser)
		handler := NewUserFilterHandler(mockService, log)

		req := httptest.NewRequest(
			http.MethodGet,
			"/users/filter?parametro_invalido=valor&limit=10",
			nil,
		)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "parâmetro de consulta inválido")

		mockService.AssertNotCalled(t, "Filter", mock.Anything, mock.Anything)
	})

	t.Run("erro - status com valor inválido", func(t *testing.T) {
		mockService := new(mockUser.MockUser)
		handler := NewUserFilterHandler(mockService, log)

		req := httptest.NewRequest(
			http.MethodGet,
			"/users/filter?status=invalido&limit=10",
			nil,
		)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)
		assert.Contains(t, resp.Message, "status deve ser true ou false")

		mockService.AssertNotCalled(t, "Filter", mock.Anything, mock.Anything)
	})

	t.Run("erro - falha ao converter filtro DTO para model", func(t *testing.T) {
		mockService := new(mockUser.MockUser)
		handler := NewUserFilterHandler(mockService, log)

		// Força erro no ToModel: email inválido
		req := httptest.NewRequest(
			http.MethodGet,
			"/users/filter?email=invalid-email",
			nil,
		)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)

		mockService.AssertNotCalled(t, "Filter", mock.Anything, mock.Anything)
	})

	t.Run("erro - filtro inválido retornado pelo serviço", func(t *testing.T) {
		mockService := new(mockUser.MockUser)
		handler := NewUserFilterHandler(mockService, log)

		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return(nil, errMsg.ErrInvalidFilter).
			Once()

		req := httptest.NewRequest(
			http.MethodGet,
			"/users/filter?limit=10&offset=0",
			nil,
		)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("erro - falha interna no serviço", func(t *testing.T) {
		mockService := new(mockUser.MockUser)
		handler := NewUserFilterHandler(mockService, log)

		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return(nil, errors.New("erro interno")).
			Once()

		req := httptest.NewRequest(
			http.MethodGet,
			"/users/filter?limit=10&offset=0",
			nil,
		)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("erro - username muito curto", func(t *testing.T) {
		mockService := new(mockUser.MockUser)
		handler := NewUserFilterHandler(mockService, log)

		req := httptest.NewRequest(
			http.MethodGet,
			"/users/filter?username=ab",
			nil,
		)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)

		mockService.AssertNotCalled(t, "Filter", mock.Anything, mock.Anything)
	})

	t.Run("sucesso - múltiplos parâmetros de data válidos", func(t *testing.T) {
		mockService := new(mockUser.MockUser)
		handler := NewUserFilterHandler(mockService, log)

		mockService.
			On("Filter", mock.Anything, mock.Anything).
			Return([]*model.User{}, nil).
			Once()

		req := httptest.NewRequest(
			http.MethodGet,
			"/users/filter?created_from=2024-01-01&created_to=2024-12-31&updated_from=2024-06-01&updated_to=2024-06-30",
			nil,
		)
		rec := httptest.NewRecorder()

		handler.Filter(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)
		mockService.AssertExpectations(t)
	})
}
