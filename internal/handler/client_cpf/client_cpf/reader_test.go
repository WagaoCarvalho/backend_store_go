package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockClient "github.com/WagaoCarvalho/backend_store_go/infra/mock/client_cpf"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/client_cpf/client"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client_cpf/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestClientHandler_GetByID(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success - Get Client by ID", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClientCpf)
		handler := NewClientCpfHandler(mockService, logAdapter)

		clientID := int64(10)
		expectedClient := &models.ClientCpf{
			ID:      clientID,
			Name:    "Cliente Teste",
			Email:   *utils.StrToPtr("teste@cliente.com"),
			CPF:     *utils.StrToPtr("123.456.789-00"),
			Version: 1,
			Status:  true,
		}

		mockService.On("GetByID", mock.Anything, clientID).Return(expectedClient, nil)

		req := newRequestWithVars(http.MethodGet, "/clients/"+fmt.Sprint(clientID), nil, map[string]string{"id": fmt.Sprint(clientID)})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Status  int              `json:"status"`
			Message string           `json:"message"`
			Data    dto.ClientCpfDTO `json:"data"`
		}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Cliente encontrado", response.Message)
		assert.Equal(t, expectedClient.ID, response.Data.ID)
		assert.Equal(t, expectedClient.Name, response.Data.Name)
		assert.Equal(t, expectedClient.Email, response.Data.Email)

		mockService.AssertExpectations(t)
	})

	t.Run("Client Not Found", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClientCpf)
		handler := NewClientCpfHandler(mockService, logAdapter)

		clientID := int64(42)
		mockService.On("GetByID", mock.Anything, clientID).Return((*models.ClientCpf)(nil), errMsg.ErrNotFound)

		req := newRequestWithVars(http.MethodGet, "/clients/"+fmt.Sprint(clientID), nil, map[string]string{"id": fmt.Sprint(clientID)})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid ID Param", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClientCpf)
		handler := NewClientCpfHandler(mockService, logAdapter)

		req := newRequestWithVars(http.MethodGet, "/clients/invalid", nil, map[string]string{"id": "invalid"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}
