package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockContact "github.com/WagaoCarvalho/backend_store_go/infra/mock/contact"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/contact"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
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
	log := logger.NewLoggerAdapter(baseLogger)

	t.Run("sucesso ao buscar contato", func(t *testing.T) {
		mockService := new(mockContact.MockContact)
		handler := NewContactHandler(mockService, log)

		expected := &models.Contact{
			ID:          1,
			ContactName: "Contato Teste",
			Email:       "teste@email.com",
			Phone:       "123456789",
		}

		mockService.
			On("GetByID", mock.Anything, int64(1)).
			Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/contacts/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp utils.DefaultResponse
		err := json.NewDecoder(w.Body).Decode(&resp)
		assert.NoError(t, err)

		data := resp.Data.(map[string]any)
		assert.Equal(t, float64(1), data["id"])
		assert.Equal(t, expected.ContactName, data["contact_name"])

		mockService.AssertExpectations(t)
	})

	t.Run("retorna 400 quando id é inválido", func(t *testing.T) {
		mockService := new(mockContact.MockContact)
		handler := NewContactHandler(mockService, log)

		req := httptest.NewRequest(http.MethodGet, "/contacts/abc", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "abc"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockService.AssertNotCalled(t, "GetByID")
	})

	t.Run("retorna 404 quando contato não existe", func(t *testing.T) {
		mockService := new(mockContact.MockContact)
		handler := NewContactHandler(mockService, log)

		mockService.
			On("GetByID", mock.Anything, int64(1)).
			Return((*models.Contact)(nil), errMsg.ErrNotFound)

		req := httptest.NewRequest(http.MethodGet, "/contacts/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		mockService.AssertExpectations(t)
	})

	t.Run("retorna 500 para erro interno do serviço", func(t *testing.T) {
		mockService := new(mockContact.MockContact)
		handler := NewContactHandler(mockService, log)

		mockService.
			On("GetByID", mock.Anything, int64(1)).
			Return((*models.Contact)(nil), errors.New("erro interno"))

		req := httptest.NewRequest(http.MethodGet, "/contacts/1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "1"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockService.AssertExpectations(t)
	})
}
