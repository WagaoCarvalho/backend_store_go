package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WagaoCarvalho/backend_store_go/internal/logger"
	model_address "github.com/WagaoCarvalho/backend_store_go/internal/models/address"
	model_contact "github.com/WagaoCarvalho/backend_store_go/internal/models/contact"
	model_user "github.com/WagaoCarvalho/backend_store_go/internal/models/user"
	model_categories "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_categories"
	model_user_full "github.com/WagaoCarvalho/backend_store_go/internal/models/user/user_full"
	services "github.com/WagaoCarvalho/backend_store_go/internal/services/users/user_full_services/user_full_services_mock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUserHandler_CreateFull(t *testing.T) {
	mockService := new(services.MockUserFullService)
	logger := logger.NewLoggerAdapter(logrus.New())
	handler := NewUserFullHandler(mockService, logger)

	t.Run("Sucesso ao criar usuário completo", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		expectedUser := &model_user_full.UserFull{
			User: &model_user.User{
				UID:      1,
				Username: "testuser",
				Email:    "test@example.com",
			},
			Address: &model_address.Address{
				Street: "Rua A",
				City:   "Cidade B",
			},
			Contact: &model_contact.Contact{
				Phone: "123456789",
			},
			Categories: []model_categories.UserCategory{
				{ID: 1},
			},
		}

		requestBody := map[string]interface{}{
			"user": map[string]interface{}{
				"username": "testuser",
				"email":    "test@example.com",
				"password": "senha123",
			},
			"address": map[string]interface{}{
				"street": "Rua A",
				"city":   "Cidade B",
			},
			"contact": map[string]interface{}{
				"phone": "123456789",
			},
			"categories": []map[string]interface{}{
				{"id": 1},
			},
		}

		body, _ := json.Marshal(requestBody)

		mockService.On("CreateFull",
			mock.Anything,
			mock.MatchedBy(func(u *model_user_full.UserFull) bool {
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

	t.Run("Erro ao criar usuário completo no service", func(t *testing.T) {
		mockService.ExpectedCalls = nil

		requestBody := map[string]interface{}{
			"user": map[string]interface{}{
				"username": "failuser",
				"email":    "fail@example.com",
				"password": "464546465",
			},
		}
		body, _ := json.Marshal(requestBody)

		mockService.On("CreateFull", mock.Anything, mock.MatchedBy(func(u *model_user_full.UserFull) bool {
			return u.User.Username == "failuser" && u.User.Email == "fail@example.com"
		})).Return(nil, errors.New("erro ao criar usuário completo")).Once()

		req := httptest.NewRequest(http.MethodPost, "/users/full", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		handler.CreateFull(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		mockService.AssertExpectations(t)
	})
}
