package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	service "github.com/WagaoCarvalho/backend_store_go/infra/mock/service/user"
	dtoAddress "github.com/WagaoCarvalho/backend_store_go/internal/dto/address"
	dtoContact "github.com/WagaoCarvalho/backend_store_go/internal/dto/contact"
	dtoUser "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/user"
	dtoUserCategories "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/user_category"
	dtoUserFull "github.com/WagaoCarvalho/backend_store_go/internal/dto/user/user_full"
	modelAddress "github.com/WagaoCarvalho/backend_store_go/internal/model/address"
	modelContact "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	modelUser "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user"
	modelCategories "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_categories"
	modelUserFull "github.com/WagaoCarvalho/backend_store_go/internal/model/user/user_full"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserHandler_CreateFull(t *testing.T) {
	mockService := new(service.MockUserFullService)
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)
	handler := NewUserFullHandler(mockService, logger)

	t.Run("Sucesso ao criar usuário completo", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		expectedUser := &modelUserFull.UserFull{
			User: &modelUser.User{
				UID:      1,
				Username: "testuser",
				Email:    "test@example.com",
			},
			Address: &modelAddress.Address{
				Street: "Rua A",
				City:   "Cidade B",
			},
			Contact: &modelContact.Contact{
				Phone: "123456789",
			},
			Categories: []modelCategories.UserCategory{
				{ID: 1},
			},
		}

		requestDTO := dtoUserFull.UserFullDTO{
			User: &dtoUser.UserDTO{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "senha123",
			},
			Address: &dtoAddress.AddressDTO{
				Street: "Rua A",
				City:   "Cidade B",
			},
			Contact: &dtoContact.ContactDTO{
				Phone: "123456789",
			},
			Categories: []dtoUserCategories.UserCategoryDTO{
				{ID: utils.UintPtr(1)},
			},
		}

		body, _ := json.Marshal(requestDTO)

		mockService.On("CreateFull",
			mock.Anything,
			mock.MatchedBy(func(u *modelUserFull.UserFull) bool {
				return u.User.Username == "testuser" &&
					u.User.Email == "test@example.com" &&
					u.Address != nil && u.Address.Street == "Rua A" &&
					u.Contact != nil && u.Contact.Phone == "123456789" &&
					len(u.Categories) == 1 && u.Categories[0].ID == 1
			}),
		).Return(expectedUser, nil).Once()

		req := httptest.NewRequest(http.MethodPost, "/users/full", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.CreateFull(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("Erro método não permitido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/full", nil)
		rec := httptest.NewRecorder()

		handler.CreateFull(rec, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})

	t.Run("Erro ao decodificar JSON inválido", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/users/full", bytes.NewReader([]byte("{invalid json")))
		rec := httptest.NewRecorder()

		handler.CreateFull(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro dados do usuário ausentes", func(t *testing.T) {
		requestDTO := dtoUserFull.UserFullDTO{
			User: nil, // obrigatório
		}
		body, _ := json.Marshal(requestDTO)

		req := httptest.NewRequest(http.MethodPost, "/users/full", bytes.NewReader(body))
		rec := httptest.NewRecorder()

		handler.CreateFull(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("Erro ao criar usuário completo no service", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		requestDTO := dtoUserFull.UserFullDTO{
			User: &dtoUser.UserDTO{
				Username: "failuser",
				Email:    "fail@example.com",
				Password: "senha123",
			},
		}
		body, _ := json.Marshal(requestDTO)

		mockService.On("CreateFull", mock.Anything, mock.Anything).
			Return(nil, errors.New("erro ao criar usuário completo")).Once()

		req := httptest.NewRequest(http.MethodPost, "/users/full", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.CreateFull(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}
