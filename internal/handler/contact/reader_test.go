package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	mockContact "github.com/WagaoCarvalho/backend_store_go/infra/mock/contact"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestContactHandler_GetByID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("deve retornar contato com sucesso", func(t *testing.T) {
		mockService := new(mockContact.MockContact)
		handler := NewContact(mockService, logger)

		expectedID := int64(1)
		expectedModel := &models.Contact{
			ID:          expectedID,
			ContactName: "Contato Teste",
			Email:       "teste@email.com",
			Phone:       "123456789",
		}

		mockService.On("GetByID", mock.Anything, int64(1)).Return(expectedModel, nil)

		req := httptest.NewRequest(http.MethodGet, "/contacts/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response utils.DefaultResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Contato encontrado", response.Message)

		dataMap := response.Data.(map[string]interface{})
		assert.Equal(t, float64(expectedID), dataMap["id"])
		assert.Equal(t, expectedModel.ContactName, dataMap["contact_name"])
		assert.Equal(t, expectedModel.Email, dataMap["email"])
		assert.Equal(t, expectedModel.Phone, dataMap["phone"])

		mockService.AssertExpectations(t)
	})

	t.Run("deve retornar erro 400 se o ID for inválido", func(t *testing.T) {
		mockService := new(mockContact.MockContact)
		handler := NewContact(mockService, logger)

		req := httptest.NewRequest(http.MethodGet, "/contacts/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("deve retornar erro 404 se o serviço retornar erro", func(t *testing.T) {
		mockService := new(mockContact.MockContact)
		handler := NewContact(mockService, logger)

		mockService.On("GetByID", mock.Anything, int64(1)).Return((*models.Contact)(nil), assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/contacts/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})
}
