package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	mockClient "github.com/WagaoCarvalho/backend_store_go/infra/mock/client"
	dto "github.com/WagaoCarvalho/backend_store_go/internal/dto/client/client"
	models "github.com/WagaoCarvalho/backend_store_go/internal/model/client/client"
	errMsg "github.com/WagaoCarvalho/backend_store_go/internal/pkg/err/message"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/logger"
	"github.com/WagaoCarvalho/backend_store_go/internal/pkg/utils"
	"github.com/gorilla/mux"
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
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logAdapter)

		clientID := int64(10)
		expectedClient := &models.Client{
			ID:      clientID,
			Name:    "Cliente Teste",
			Email:   utils.StrToPtr("teste@cliente.com"),
			CPF:     utils.StrToPtr("123.456.789-00"),
			CNPJ:    nil,
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
			Status  int           `json:"status"`
			Message string        `json:"message"`
			Data    dto.ClientDTO `json:"data"`
		}
		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Cliente encontrado", response.Message)
		assert.Equal(t, expectedClient.ID, *response.Data.ID)
		assert.Equal(t, expectedClient.Name, response.Data.Name)
		assert.Equal(t, *expectedClient.Email, *response.Data.Email)

		mockService.AssertExpectations(t)
	})

	t.Run("Client Not Found", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logAdapter)

		clientID := int64(42)
		mockService.On("GetByID", mock.Anything, clientID).Return((*models.Client)(nil), errMsg.ErrNotFound)

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
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logAdapter)

		req := newRequestWithVars(http.MethodGet, "/clients/invalid", nil, map[string]string{"id": "invalid"})
		w := httptest.NewRecorder()

		handler.GetByID(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestClientHandler_GetByName(t *testing.T) {
	baseLogger := logrus.New()
	baseLogger.Out = &bytes.Buffer{}
	logAdapter := logger.NewLoggerAdapter(baseLogger)

	t.Run("Success - Get Clients by Name", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logAdapter)

		name := "Cliente Teste"
		clientModel := &models.Client{
			ID:   1,
			Name: name,
		}

		// agora o mock retorna slice, não apenas 1 objeto
		mockService.On("GetByName", mock.Anything, "Cliente Teste").
			Return([]*models.Client{clientModel}, nil)

		req := httptest.NewRequest(http.MethodGet, "/clients/name/Cliente%20Teste", nil)
		req = mux.SetURLVars(req, map[string]string{"name": "Cliente Teste"})
		w := httptest.NewRecorder()

		handler.GetByName(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var response struct {
			Status  int             `json:"status"`
			Message string          `json:"message"`
			Data    []dto.ClientDTO `json:"data"`
		}

		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Clientes encontrados", response.Message)
		assert.Equal(t, clientModel.ID, *response.Data[0].ID)
		assert.Equal(t, clientModel.Name, response.Data[0].Name)

		mockService.AssertExpectations(t)
	})

	t.Run("Not Found - empty list", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logAdapter)

		mockService.On("GetByName", mock.Anything, "Inexistente").
			Return([]*models.Client{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/clients/name/Inexistente", nil)
		req = mux.SetURLVars(req, map[string]string{"name": "Inexistente"})
		w := httptest.NewRecorder()

		handler.GetByName(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode) // agora 200, não 404

		var response struct {
			Status  int             `json:"status"`
			Message string          `json:"message"`
			Data    []dto.ClientDTO `json:"data"`
		}

		err := json.NewDecoder(resp.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Equal(t, "Nenhum cliente encontrado", response.Message)
		assert.Empty(t, response.Data)

		mockService.AssertExpectations(t)
	})

	t.Run("Invalid param", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logAdapter)

		req := httptest.NewRequest(http.MethodGet, "/clients/name/", nil)
		req = mux.SetURLVars(req, map[string]string{"name": ""})
		w := httptest.NewRecorder()

		handler.GetByName(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Service Error", func(t *testing.T) {
		t.Parallel()
		mockService := new(mockClient.MockClient)
		handler := NewClientHandler(mockService, logAdapter)

		mockService.On("GetByName", mock.Anything, "Erro").
			Return(nil, assert.AnError)

		req := httptest.NewRequest(http.MethodGet, "/clients/name/Erro", nil)
		req = mux.SetURLVars(req, map[string]string{"name": "Erro"})
		w := httptest.NewRecorder()

		handler.GetByName(w, req)

		resp := w.Result()
		defer resp.Body.Close()

		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockService.AssertExpectations(t)
	})

}
