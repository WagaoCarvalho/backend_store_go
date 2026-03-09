package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	mockClient "github.com/WagaoCarvalho/backend_store_go/infra/mock/client_cpf"
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
	log := logger.NewLoggerAdapter(baseLogger)

	t.Run("erro - id inválido", func(t *testing.T) {
		mockService := new(mockClient.MockClientCpf)
		handler := NewClientCpfHandler(mockService, log)

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

	t.Run("erro - erro retornado pelo serviço", func(t *testing.T) {
		mockService := new(mockClient.MockClientCpf)
		handler := NewClientCpfHandler(mockService, log)

		mockService.
			On("GetVersionByID", mock.Anything, int64(123)).
			Return(0, errors.New("erro no banco"))

		req := httptest.NewRequest(http.MethodGet, "/clients/123/version", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "123"})
		rec := httptest.NewRecorder()

		handler.GetVersionByID(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)

		var resp utils.DefaultResponse
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("sucesso - versão encontrada", func(t *testing.T) {
		mockService := new(mockClient.MockClientCpf)
		handler := NewClientCpfHandler(mockService, log)

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

		data, ok := resp.Data.(map[string]any)
		assert.True(t, ok)
		assert.Equal(t, float64(123), data["client_id"])
		assert.Equal(t, float64(expectedVersion), data["version"])

		mockService.AssertExpectations(t)
	})
}
