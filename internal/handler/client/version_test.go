package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockClient "github.com/WagaoCarvalho/backend_store_go/infra/mock/client"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClientHandler_GetVersionByID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logger := logger.NewLoggerAdapter(baseLogger)

	t.Run("erro - id inválido", func(t *testing.T) {
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logger)

		req := httptest.NewRequest(http.MethodGet, "/clients/invalid/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "invalid"})
		rec := httptest.NewRecorder()

		handler.GetVersionByID(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("erro - erro no serviço", func(t *testing.T) {
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logger)

		mockService.
			On("GetVersionByID", mock.Anything, int64(123)).
			Return(0, errors.New("erro no banco"))

		req := httptest.NewRequest(http.MethodGet, "/clients/123/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "123"})
		rec := httptest.NewRecorder()

		handler.GetVersionByID(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - versão encontrada", func(t *testing.T) {
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logger)

		expectedVersion := 5

		mockService.
			On("GetVersionByID", mock.Anything, int64(123)).
			Return(expectedVersion, nil)

		req := httptest.NewRequest(http.MethodGet, "/clients/123/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "123"})
		rec := httptest.NewRecorder()

		handler.GetVersionByID(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusOK, resp.Status)
		assert.Equal(t, "Versão do cliente recuperada com sucesso", resp.Message)

		dataMap, ok := resp.Data.(map[string]any)
		assert.True(t, ok)
		assert.Equal(t, float64(123), dataMap["client_id"]) // JSON converte int64 para float64
		assert.Equal(t, float64(expectedVersion), dataMap["version"])

		mockService.AssertExpectations(t)
	})
}
